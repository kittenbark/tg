package tg

import "context"

type Album = []InputMedia

var (
	_ InputMedia = &Photo{}
)

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
