package tg

import (
	"context"
	"sync"
	"time"
)

type Scheduler interface {
	Schedule(ctx context.Context, chat int64, weight int)
	Done(ctx context.Context, chat int64, weight int)
}

func NewScheduler(clauses ...SchedulerClause) Scheduler {
	return NewSchedulerVerbose(time.Millisecond*200, clauses...)
}

func NewSchedulerVerbose(pollingRate time.Duration, clauses ...SchedulerClause) Scheduler {
	if clauses == nil {
		clauses = []SchedulerClause{
			SchedulerClauseGlobal(30, time.Millisecond*1_500),
			SchedulerClauseChat(20, time.Millisecond*60_500),
			SchedulerClauseChat(10, time.Millisecond*10_500),
			SchedulerClauseUser(100, time.Millisecond*30_500),
			SchedulerClauseUser(30, time.Millisecond*5_500),
		}
	}
	return &clauseScheduler{
		clauses:     clauses,
		pollingRate: pollingRate,
		mutex:       &sync.Mutex{},
	}
}

type clauseScheduler struct {
	clauses     []SchedulerClause
	pollingRate time.Duration
	mutex       *sync.Mutex
}

func (scheduler *clauseScheduler) Schedule(ctx context.Context, chat int64, weight int) {
	for {
		if scheduler.trySchedule(chat, weight) {
			return
		}
		time.Sleep(scheduler.pollingRate)
	}
}

func (scheduler *clauseScheduler) Done(ctx context.Context, chat int64, weight int) {
	for _, clause := range scheduler.clauses {
		clause.Done(scheduler.mutex, chat, weight)
	}
}

func (scheduler *clauseScheduler) trySchedule(chat int64, weight int) bool {
	scheduler.mutex.Lock()
	defer scheduler.mutex.Unlock()

	for _, clause := range scheduler.clauses {
		if !clause.TrySchedule(chat, weight) {
			return false
		}
	}

	for _, clause := range scheduler.clauses {
		clause.Schedule(chat, weight)
	}
	return true
}

type SchedulerClause interface {
	TrySchedule(chat int64, weight int) bool
	Schedule(chat int64, weight int)
	Done(mutex *sync.Mutex, chat int64, weight int)
}

var (
	_ SchedulerClause = (*schedulerClauseCounter)(nil)
	_ SchedulerClause = (*schedulerClauseChat)(nil)
)

func SchedulerClauseGlobal(quota int, timeout time.Duration) SchedulerClause {
	return &schedulerClauseCounter{
		quota:   quota,
		timeout: timeout,
		state:   0,
	}
}

func SchedulerClauseUser(quota int, timeout time.Duration) SchedulerClause {
	return &schedulerClauseChat{
		quota:   quota,
		timeout: timeout,
		pred:    func(chat int64) bool { return chat > 0 },
		state:   map[int64]int{},
	}
}

func SchedulerClauseChat(quota int, timeout time.Duration) SchedulerClause {
	return &schedulerClauseChat{
		quota:   quota,
		timeout: timeout,
		pred:    func(chat int64) bool { return chat < 0 },
		state:   map[int64]int{},
	}
}

type schedulerClauseCounter struct {
	quota   int
	timeout time.Duration

	state int
}

func (clause *schedulerClauseCounter) TrySchedule(chat int64, weight int) bool {
	return clause.state+weight <= clause.quota
}

func (clause *schedulerClauseCounter) Schedule(chat int64, weight int) {
	clause.state += weight
}

func (clause *schedulerClauseCounter) Done(mutex *sync.Mutex, chat int64, weight int) {
	time.AfterFunc(clause.timeout, func() {
		mutex.Lock()
		defer mutex.Unlock()

		clause.state -= weight
	})
}

type schedulerClauseChat struct {
	quota   int
	timeout time.Duration

	pred  func(chat int64) bool
	state map[int64]int
}

func (clause *schedulerClauseChat) TrySchedule(chat int64, weight int) bool {
	if !clause.pred(chat) {
		return true
	}
	if state, ok := clause.state[chat]; ok {
		return state+weight <= clause.quota
	}
	return true
}

func (clause *schedulerClauseChat) Schedule(chat int64, weight int) {
	if !clause.pred(chat) {
		return
	}
	clause.state[chat] += weight
}

func (clause *schedulerClauseChat) Done(mutex *sync.Mutex, chat int64, weight int) {
	if !clause.pred(chat) {
		return
	}
	time.AfterFunc(clause.timeout, func() {
		mutex.Lock()
		defer mutex.Unlock()

		if state, ok := clause.state[chat]; ok {
			if weight >= state {
				delete(clause.state, chat)
			} else {
				clause.state[chat] -= weight
			}
		}
	})
}

func schedule(ctx context.Context, chat int64, weight int) {
	if scheduler, ok := ctx.Value(ContextScheduler).(Scheduler); ok {
		scheduler.Schedule(ctx, chat, weight)
	}
}

func scheduleDone(ctx context.Context, chat int64, weight int) {
	if scheduler, ok := ctx.Value(ContextScheduler).(Scheduler); ok {
		scheduler.Done(ctx, chat, weight)
	}
}
