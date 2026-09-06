package tg

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"time"
)

type Album = []InputMedia

var (
	_ InputMedia = &Photo{}
)

// HandleAlbum groups updates corresponding to the same MediaGroupId.
//
// Uses a sliding-window timer: the window resets on each new part that arrives,
// so the handler fires as soon as the album is quiet for HandlingTimeout rather
// than after a fixed worst-case sleep. The window defaults to 300ms.
//
// The first update's goroutine blocks until the window closes (context-aware),
// which means errors from fn propagate normally through the bot's error hooks.
func HandleAlbum(fn func(ctx context.Context, updates []*Update) error, cfg ...*ConfigHandleAlbum) HandlerFunc {
	config := at(cfg, 0, &ConfigHandleAlbum{HandlingTimeout: 300 * time.Millisecond})

	type group struct {
		updates []*Update
		timer   *time.Timer
		done    chan struct{}
		err     error
	}

	var mu sync.Mutex
	groups := map[string]*group{}

	return func(ctx context.Context, upd *Update) error {
		if upd == nil || upd.Message == nil || upd.Message.MediaGroupId == "" {
			return fn(ctx, []*Update{upd})
		}
		id := upd.Message.MediaGroupId

		mu.Lock()
		g, exists := groups[id]
		if exists {
			g.updates = append(g.updates, upd)
			g.timer.Reset(config.HandlingTimeout)
			mu.Unlock()
			return nil
		}

		done := make(chan struct{})
		g = &group{updates: []*Update{upd}, done: done}
		groups[id] = g
		g.timer = time.AfterFunc(config.HandlingTimeout, func() {
			mu.Lock()
			album := g.updates
			delete(groups, id)
			mu.Unlock()
			slices.SortFunc(album, func(a, b *Update) int { return cmp.Compare(a.Message.MessageId, b.Message.MessageId) })
			g.err = fn(ctx, album)
			close(done)
		})
		mu.Unlock()

		select {
		case <-done:
			return g.err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type ConfigHandleAlbum struct {
	HandlingTimeout time.Duration
}

// TelegramPhoto is wrapper around list of PhotoSize, the last in the list is the biggest picture.
type TelegramPhoto []*PhotoSize

func (photo TelegramPhoto) FileId() string       { return photo.biggest().FileId }
func (photo TelegramPhoto) FileUniqueId() string { return photo.biggest().FileUniqueId }
func (photo TelegramPhoto) Width() int64         { return photo.biggest().Width }
func (photo TelegramPhoto) Height() int64        { return photo.biggest().Height }
func (photo TelegramPhoto) FileSize() int64      { return photo.biggest().FileSize }
func (photo TelegramPhoto) biggest() *PhotoSize  { return photo[len(photo)-1] }

func (photo TelegramPhoto) Download(ctx context.Context, path string) error {
	if len(photo) == 0 {
		return &Error{Description: "no photo to download"}
	}
	return (photo)[len(photo)-1].Download(ctx, path)
}

func (photo TelegramPhoto) DownloadTemp(ctx context.Context, dirAndPattern ...string) (string, error) {
	if len(photo) == 0 {
		return "", &Error{Description: "no photo to download"}
	}
	return (photo)[len(photo)-1].DownloadTemp(ctx, dirAndPattern...)
}
