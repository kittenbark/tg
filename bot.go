package tg

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"
)

type FilterFunc func(ctx context.Context, upd *Update) bool
type HandlerFunc func(ctx context.Context, upd *Update) error
type OnErrorFunc func(ctx context.Context, err error)

type Bot struct {
	context           context.Context
	contextTimeout    time.Duration
	contextCancelFunc context.CancelFunc

	pipeline       *pipe
	plugins        map[PluginHookType][]Plugin
	defaultHandler HandlerFunc

	syncHandling  bool
	pollTimeout   time.Duration
	updatesOffset int64
}

// Filter out unwanted updates — FilterFunc(unwanted_update) == false,
// every update for handler/branch below the filter will satisfy FilterFunc.
//
// Example:
//
//	bot.
//		Filter(tg.OnPrivateMessage).
//		Command("/start", ...). // Private messages only.
//		Branch(tg.OnText, ...) 		  // Private messages only.
func (bot *Bot) Filter(pred FilterFunc) *Bot {
	bot.pipeline.Last().Next = &pipe{Filter: pred}
	return bot
}

// Handle every update that comes through.
//
// Example:
//
//	bot.
//		Filter(tg.OnPrivateMessage).
//		Handle(tg.CommonTextReply("hii mom"))
func (bot *Bot) Handle(handler HandlerFunc) *Bot {
	bot.pipeline.Last().Next = &pipe{Handle: handler}
	return bot
}

// Default handles every filtered out update.
func (bot *Bot) Default(handler HandlerFunc) *Bot {
	bot.defaultHandler = handler
	return bot
}

// Command and only that command, if message do not start with "/" and the command's name,
// the update will be passed below the tree.
//
// Example:
//
//	bot.
//		Command("/start", ...).			// Handles only /start.
//		Handle(tg.CommonTextReply("hii mom")) 	// Handles everything but /start.
func (bot *Bot) Command(command string, handlerFunc HandlerFunc) *Bot {
	return bot.Branch(OnCommand(command), handlerFunc)
}

// Branch helps to separate and handle updates by predicates.
// If an update does not satisfy Branch's predicates, the updated is passed through to next branch/handler/filter below.
//
// Example:
//
//	bot.
//		Branch(tg.OnPhoto, tg.CommonTextReply("nice photo mom)).
//		Branch(tg.OnVideo, tg.CommonTextReply("cool video mom)).
//	    Handle(tg.CommonTextReply("love you mom"))	// This shall be sent for every non photo/video update.
func (bot *Bot) Branch(pred FilterFunc, handler HandlerFunc) *Bot {
	bot.complexBranch(
		Branch().Filter(pred).Handle(handler),
	)
	return bot
}

// complexBranch is pretty much alike Branch, but it *should* allow branching out pretty much alike bot declaring.
// todo: implement it for real, now its just filters.
func (bot *Bot) complexBranch(branch *BranchPipe) *Bot {
	bot.pipeline.Last().Next = &pipe{Branch: branch.pipeline}
	return bot
}

// Plugin system helps insert middleware.
// todo: document this better and probably add more plugin hooks.
func (bot *Bot) Plugin(plugin ...Plugin) *Bot {
	for _, plugin := range plugin {
		for _, hook := range plugin.Hooks() {
			bot.plugins[hook] = append(bot.plugins[hook], plugin)
		}
	}
	return bot
}

// Scheduler ensures no 429 Too Many Requests.
func (bot *Bot) Scheduler(scheduler ...Scheduler) *Bot {
	bot.context = context.WithValue(bot.context, ContextScheduler, at(scheduler, 0, NewScheduler()))
	return bot
}

// ContextWithCancel build new Context with a fresh timeout.
func (bot *Bot) ContextWithCancel() (ctx context.Context, cancel context.CancelFunc) {
	if bot.contextTimeout == 0 {
		return bot.context, func() {}
	}
	return context.WithTimeout(bot.context, bot.contextTimeout)
}

// Context is simplified version of ContextWithCancel, might temporarily leak resources if used with timeouts.
// Leaking recourses is OK if you do not care, but I encourage you to use ContextWithCancel.
func (bot *Bot) Context() context.Context {
	ctx, _ := bot.ContextWithCancel()
	return ctx
}

func (bot *Bot) OnError(fn OnErrorFunc) *Bot {
	return bot.Plugin(PluginOnError(fn))
}

func (bot *Bot) Ping() bool {
	_, err := GetMe(bot.Context())
	return err == nil
}

// Start locks the execution, interruptible with Stop.
// todo: implement webhooks.
func (bot *Bot) Start() {
	for {
		pollStart := time.Now()
		select {
		case <-bot.context.Done():
			return
		default:
			bot.longPollIteration()
		}
		time.Sleep(bot.pollTimeout - time.Since(pollStart))
	}
}

func (bot *Bot) Stop() {
	bot.contextCancelFunc()
}

func (bot *Bot) longPollIteration() {
	// Note: Telegram gives you 3s timeout if you have empty list if updates and poll for updates too often.
	updatesCtx, updatesCtxCancel := context.WithTimeout(bot.context, time.Second*4)
	updates, err := GetUpdates(updatesCtx, &OptGetUpdates{Offset: bot.updatesOffset})
	updatesCtxCancel()
	if err != nil {
		bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{updatesCtx, bot, err})
		return
	}
	if len(updates) == 0 {
		return
	}

	ctx, ctxCancel := bot.ContextWithCancel()
	ctxCancelWg := &sync.WaitGroup{}
	ctxCancelWg.Add(len(updates))
	go func() {
		ctxCancelWg.Wait()
		ctxCancel()
	}()

	for _, update := range slices.Backward(updates) {
		bot.pluginsHook(PluginHookOnUpdate, &PluginHookContextOnUpdate{ctx, bot, update})
		if bot.syncHandling {
			bot.handle(ctxCancelWg, ctx, update)
		} else {
			go bot.handle(ctxCancelWg, ctx, update)
		}
		bot.updatesOffset = max(bot.updatesOffset, update.UpdateId+1)
	}
}

func (bot *Bot) handle(updatesCancelContextWg *sync.WaitGroup, ctx context.Context, update *Update) {
	defer updatesCancelContextWg.Done()
	defer func() {
		if rec := recover(); rec != nil {
			bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{ctx, bot, fmt.Errorf("panic: %v", rec)})
		}
	}()

	if !bot.handlePipe(bot.pipeline, ctx, update) && bot.defaultHandler != nil {
		if err := bot.defaultHandler(ctx, update); err != nil {
			bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{ctx, bot, err})
		}
	}
}

func (bot *Bot) handlePipe(pipe *pipe, ctx context.Context, update *Update) bool {
	switch {
	case pipe == nil:
		return false

	case pipe.Filter != nil && !pipe.Filter(ctx, update):
		bot.pluginsHook(PluginHookOnFilter, &PluginHookContextOnFilter{ctx, bot, pipe.Filter})
		return false

	case pipe.Handle != nil:
		bot.pluginsHook(PluginHookOnHandleStart, &PluginHookContextOnHandleStart{ctx, bot, update, pipe.Handle})
		if err := pipe.Handle(ctx, update); err != nil {
			bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{ctx, bot, err})
		}
		bot.pluginsHook(PluginHookOnHandleFinish, &PluginHookContextOnHandleFinish{ctx, bot, update, pipe.Handle})
		return true

	default:
		return bot.handlePipe(pipe.Branch, ctx, update) || bot.handlePipe(pipe.Next, ctx, update)
	}
}

func (bot *Bot) pluginsHook(hook PluginHookType, ctx PluginHookContext) {
	plugins := bot.plugins[hook]
	if len(plugins) == 0 {
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(plugins))
	for _, plugin := range plugins {
		go func() {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					slog.Error("pluginsHook%panic", "err", r)
				}
			}()

			plugin.Apply(ctx)
		}()
	}
	wg.Wait()
}

// BranchPipe has similar interface to Handle/Filter/Branch Bot's pipeline configuring.
// The only difference making new branches is prohibited (it could cause infinite cycles).
type BranchPipe struct {
	pipeline *pipe
}

func Branch() *BranchPipe {
	return &BranchPipe{
		pipeline: &pipe{},
	}
}

func (branch *BranchPipe) Filter(pred FilterFunc) *BranchPipe {
	branch.pipeline.Last().Next = &pipe{Filter: pred}
	return branch
}

func (branch *BranchPipe) Handle(handler HandlerFunc) *BranchPipe {
	branch.pipeline.Last().Next = &pipe{Handle: handler}
	return branch
}

// pipe is, ugh, a part of pipeline.
// - Filter out update.
// - Handle update.
// - Branch out.
type pipe struct {
	Filter FilterFunc
	Handle HandlerFunc
	Branch *pipe
	Next   *pipe
}

func (p *pipe) Last() *pipe {
	if p == nil {
		return nil
	}
	pipe := p
	for pipe.Next != nil {
		pipe = pipe.Next
	}
	return pipe
}
