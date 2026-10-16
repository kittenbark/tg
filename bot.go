package tg

import (
	"context"
	"errors"
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
	stopUpdates       chan bool
	contextTimeout    time.Duration
	contextCancelFunc context.CancelFunc

	// lifecycleLock guards contextCancelFunc, stopUpdates and ready —
	// Start writes them, Stop/StopImmediately/Ready read them from other
	// goroutines, so plain field access races (confirmed with -race).
	lifecycleLock sync.Mutex
	ready         chan struct{}

	pipelineLock   sync.Mutex
	pipeline       pipe
	plugins        map[PluginHookType][]Plugin
	defaultHandler HandlerFunc

	syncHandling  bool
	syncStart     sync.Mutex
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
func (bot *Bot) Filter(pred ...FilterFunc) *Bot {
	bot.pipelineLock.Lock()
	defer bot.pipelineLock.Unlock()
	for _, fn := range pred {
		bot.pipeline.Last().Next = &pipe{Filter: fn}
	}
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
	bot.pipelineLock.Lock()
	defer bot.pipelineLock.Unlock()
	bot.pipeline.Last().Next = &pipe{Handle: handler}
	return bot
}

// Default handles every filtered out update.
func (bot *Bot) Default(handler HandlerFunc) *Bot {
	bot.defaultHandler = handler
	return bot
}

// Command and only that command, if message do not start with "/" and the command's Name,
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
	bot.pipelineLock.Lock()
	defer bot.pipelineLock.Unlock()
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

// RedactToken controls whether the bot's own token is scrubbed out of any text
// (message text, captions, button labels, ...) before it's sent to Telegram —
// enabled by default.
func (bot *Bot) RedactToken(enabled ...bool) *Bot {
	bot.context = context.WithValue(bot.context, ContextRedactToken, at(enabled, 0, true))
	return bot
}

// MediaAutofill controls whether sending a local video/document/animation/
// audio/video-note file automatically generates a thumbnail and fills
// width/height/duration when the caller didn't already set them —
// enabled by default.
func (bot *Bot) MediaAutofill(enabled ...bool) *Bot {
	bot.context = context.WithValue(bot.context, ContextMediaAutofill, at(enabled, 0, true))
	return bot
}

// VideoFrameDecode controls whether media autofill best-effort decodes a
// local H.264 video's first frame into a real thumbnail, independent of
// MediaAutofill as a whole (e.g. to keep cheap width/height/duration
// autofill while disabling this more expensive, more experimental decode
// path) — enabled by default.
func (bot *Bot) VideoFrameDecode(enabled ...bool) *Bot {
	bot.context = context.WithValue(bot.context, ContextVideoFrameDecode, at(enabled, 0, true))
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

func (bot *Bot) ExtraContext(fn ...ExtraContext) *Bot {
	for _, f := range fn {
		bot.context = f(bot.context)
	}
	return bot
}

func (bot *Bot) OnError(fn OnErrorFunc) *Bot {
	return bot.Plugin(PluginOnError(fn))
}

func (bot *Bot) Help(commands ...string) *Bot {
	return bot.HelpScoped(&BotCommandScopeDefault{Type: "default"}, commands...)
}

func (bot *Bot) HelpScoped(scope BotCommandScope, commands ...string) *Bot {
	result := []*BotCommand{}
	for info := range slices.Chunk(commands, 2) {
		result = append(result, &BotCommand{
			Command:     info[0],
			Description: at(info, 1, info[0]),
		})
	}

	ctx, cancel := bot.ContextWithCancel()
	defer cancel()
	if _, err := SetMyCommands(ctx, result, &OptSetMyCommands{Scope: scope}); err != nil {
		bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{ctx, bot, err})
	}

	return bot
}

// Start locks the execution, interruptible with Stop.
// todo: implement webhooks.
func (bot *Bot) Start(updates ...*Update) {
	bot.syncStart.Lock()
	defer bot.syncStart.Unlock()

	ctx, cancel := context.WithCancel(bot.context)
	defer cancel()

	bot.lifecycleLock.Lock()
	bot.contextCancelFunc = cancel
	stopUpdates := make(chan bool, 1)
	bot.stopUpdates = stopUpdates
	// Reuse whatever channel Ready() is currently handing out (the one from
	// construction, or the one a prior Start left behind — see below) rather
	// than replacing it: a caller may have already fetched it via Ready()
	// before this goroutine got scheduled, and swapping in a new channel
	// here would leave that caller waiting on one that never closes.
	ready := bot.ready
	bot.lifecycleLock.Unlock()

	// Leave a fresh, open channel behind for the next Start call, once this
	// one is done with its own ready channel (already closed by then).
	defer func() {
		bot.lifecycleLock.Lock()
		bot.ready = make(chan struct{})
		bot.lifecycleLock.Unlock()
	}()

	bot.handleUpdates(ctx, updates)
	close(ready)

	for {
		pollStart := time.Now()
		select {
		case <-ctx.Done():
			return
		case <-stopUpdates:
			return
		default:
			bot.longPollIteration(ctx)
		}
		time.Sleep(bot.pollTimeout - time.Since(pollStart))
	}
}

// Stop lets handlers finish their jobs, do not check for any more updates.
// Safe to call before Start (a no-op) and concurrently with Start/StopImmediately.
func (bot *Bot) Stop() {
	bot.lifecycleLock.Lock()
	stopUpdates := bot.stopUpdates
	bot.lifecycleLock.Unlock()

	if stopUpdates == nil {
		return
	}
	select {
	case stopUpdates <- true:
	default:
	}
}

// StopImmediately stops polling, by canceling the context.
// Safe to call before Start (a no-op) and concurrently with Start/Stop.
func (bot *Bot) StopImmediately() {
	bot.lifecycleLock.Lock()
	cancel := bot.contextCancelFunc
	bot.lifecycleLock.Unlock()

	if cancel != nil {
		cancel()
	}
}

// Ready returns a channel that closes once Start has processed any
// pre-loaded updates and is about to begin polling. Safe to call before
// Start — the channel it returns is the one that specific next Start call
// will close. After Start returns (Stop/StopImmediately), a fresh, not yet
// closed channel takes its place for the next Start call.
func (bot *Bot) Ready() <-chan struct{} {
	bot.lifecycleLock.Lock()
	defer bot.lifecycleLock.Unlock()
	return bot.ready
}

func (bot *Bot) longPollIteration(ctx context.Context) {
	updatesCtx, updatesCtxCancel := bot.ContextWithCancel()
	updates, err := GetUpdates(updatesCtx, &OptGetUpdates{Offset: bot.updatesOffset})
	updatesCtxCancel()
	// Note: Telegram gives you 3s+ timeout if you have empty list if updates and poll for updates too often.
	if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{updatesCtx, bot, err})
		if IsUnauthorized(err) {
			bot.StopImmediately()
		}
		return
	}
	if len(updates) == 0 {
		return
	}

	bot.handleUpdates(ctx, updates)
}

func (bot *Bot) handleUpdates(ctx context.Context, updates []*Update) {
	for _, update := range updates {
		bot.pluginsHook(PluginHookOnUpdate, &PluginHookContextOnUpdate{ctx, bot, update})

		// Each update gets its own independent timeout so a slow handler in one
		// update does not eat into the timeout budget of subsequent updates.
		updCtx, updCancel := bot.ContextWithCancel()

		if bot.syncHandling {
			bot.handle(updCtx, update)
			updCancel()
		} else {
			go func(u *Update, cancel context.CancelFunc) {
				defer cancel()
				bot.handle(updCtx, u)
			}(update, updCancel)
		}
		bot.updatesOffset = max(bot.updatesOffset, update.UpdateId+1)
	}
}

func (bot *Bot) handle(ctx context.Context, update *Update) {
	defer func() {
		if rec := recover(); rec != nil {
			bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{ctx, bot, fmt.Errorf("panic: %v", rec)})
		}
	}()

	if !bot.handlePipe(&bot.pipeline, ctx, update) && bot.defaultHandler != nil {
		if err := bot.defaultHandler(ctx, update); err != nil {
			bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{ctx, bot, err})
		}
	}
}

func (bot *Bot) handlePipe(pipe *pipe, ctx context.Context, update *Update) (handled bool) {
	switch {
	case pipe == nil:
		return false

	case pipe.Filter != nil && !filterWrappedPanics(pipe.Filter, ctx, update):
		bot.pluginsHook(PluginHookOnFilter, &PluginHookContextOnFilter{ctx, bot, pipe.Filter})
		return false

	case pipe.Handle != nil:
		bot.pluginsHook(PluginHookOnHandleStart, &PluginHookContextOnHandleStart{ctx, bot, update, pipe.Handle})
		err := pipe.Handle(ctx, update)
		if err != nil {
			bot.pluginsHook(PluginHookOnError, &PluginHookContextOnError{ctx, bot, err})
		}
		bot.pluginsHook(PluginHookOnHandleFinish, &PluginHookContextOnHandleFinish{ctx, bot, update, pipe.Handle, err})
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
					slog.Error("pluginsHook#panic", "err", r)
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

func filterWrappedPanics(filter FilterFunc, ctx context.Context, upd *Update) (result bool) {
	defer func() {
		if rec := recover(); rec != nil {
			result = false
		}
	}()
	return filter(ctx, upd)
}
