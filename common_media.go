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
// Note: handling is happening after a delay, which could be adjusted with ConfigHandleAlbum (500ms by default).
func HandleAlbum(fn func(ctx context.Context, updates []*Update) error, cfg ...*ConfigHandleAlbum) HandlerFunc {
	config := at(cfg, 0, &ConfigHandleAlbum{
		HandlingTimeout: parseFromEnvDurationMust(EnvTimeoutPolling, defaultPollingTimeout*2+time.Millisecond*50),
	})
	cacheMutex := &sync.Mutex{}
	cache := map[string][]*Update{}
	return func(ctx context.Context, upd *Update) error {
		if upd == nil || upd.Message == nil || upd.Message.MediaGroupId == "" {
			return fn(ctx, []*Update{upd})
		}

		cacheMutex.Lock()
		if _, ok := cache[upd.Message.MediaGroupId]; ok {
			cache[upd.Message.MediaGroupId] = append(cache[upd.Message.MediaGroupId], upd)
			cacheMutex.Unlock()
			return nil
		}

		cache[upd.Message.MediaGroupId] = []*Update{upd}
		cacheMutex.Unlock()
		defer func() {
			cacheMutex.Lock()
			defer cacheMutex.Unlock()
			delete(cache, upd.Message.MediaGroupId)
		}()

		time.Sleep(config.HandlingTimeout)
		cacheMutex.Lock()
		album := cache[upd.Message.MediaGroupId]
		slices.SortFunc(album, func(a, b *Update) int { return cmp.Compare(a.Message.MessageId, b.Message.MessageId) })
		cacheMutex.Unlock()
		return fn(ctx, album)
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
