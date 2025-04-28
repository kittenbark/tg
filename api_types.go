package tg

import (
	"context"
)

// AffiliateInfo Contains information about the affiliate that received a commission via this transaction.
type AffiliateInfo struct {
	// Optional. The bot or the user that received an affiliate commission if it was received by a bot or a user
	AffiliateUser *User `json:"affiliate_user,omitempty"`
	// Optional. The chat that received an affiliate commission if it was received by a chat
	AffiliateChat *Chat `json:"affiliate_chat,omitempty"`
	// The number of Telegram Stars received by the affiliate for each 1000 Telegram Stars received by the bot from referred users
	CommissionPerMille int64 `json:"commission_per_mille"`
	// Integer amount of Telegram Stars received by the affiliate from the transaction, rounded to 0; can be negative for refunds
	Amount int64 `json:"amount"`
	// Optional.
	// The number of 1/1000000000 shares of Telegram Stars received by the affiliate; from -999999999 to 999999999; can be negative for refunds
	NanostarAmount int64 `json:"nanostar_amount,omitempty"`
}

// Animation Represents an animation file (GIF or H.264/MPEG-4 AVC video without sound) to be sent.
type Animation struct {
	// Type of the result, must be animation
	Type string `json:"type" default:"animation"`
	// File to send. More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// Pass a file_id to send a file that exists on the Telegram servers (recommended), pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	Media InputFile `json:"media"`
	// Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size.
	// A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data.
	// Thumbnails can't be reused and can be only uploaded as a new file, so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	// More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// >> either: String
	Thumbnail InputFile `json:"thumbnail,omitempty"`
	// Optional. Caption of the animation to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the animation caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Animation width
	Width int64 `json:"width,omitempty"`
	// Optional. Animation height
	Height int64 `json:"height,omitempty"`
	// Optional. Animation duration in seconds
	Duration int64 `json:"duration,omitempty"`
	// Optional. Pass True if the animation needs to be covered with a spoiler animation
	HasSpoiler bool `json:"has_spoiler,omitempty"`
	// Used for uploading media.
	InputFile InputFile `json:"-"`
}

// Audio Represents an audio file to be treated as music to be sent.
type Audio struct {
	// Type of the result, must be audio
	Type string `json:"type" default:"audio"`
	// File to send. More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// Pass a file_id to send a file that exists on the Telegram servers (recommended), pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	Media InputFile `json:"media"`
	// Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size.
	// A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data.
	// Thumbnails can't be reused and can be only uploaded as a new file, so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	// More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// >> either: String
	Thumbnail InputFile `json:"thumbnail,omitempty"`
	// Optional. Caption of the audio to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the audio caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Duration of the audio in seconds
	Duration int64 `json:"duration,omitempty"`
	// Optional. Performer of the audio
	Performer string `json:"performer,omitempty"`
	// Optional. Title of the audio
	Title string `json:"title,omitempty"`
	// Used for uploading media.
	InputFile InputFile `json:"-"`
}

// BackgroundFill This object describes the way a background is filled based on the selected colors. Currently, it can be one of
// - BackgroundFillSolid
// - BackgroundFillGradient
// - BackgroundFillFreeformGradient
type BackgroundFill interface {
	OptSolid() *BackgroundFillSolid
	OptGradient() *BackgroundFillGradient
	OptFreeformGradient() *BackgroundFillFreeformGradient
}

var (
	_ BackgroundFill = &BackgroundFillSolid{}
	_ BackgroundFill = &BackgroundFillGradient{}
	_ BackgroundFill = &BackgroundFillFreeformGradient{}
)

func (impl *BackgroundFillSolid) OptSolid() *BackgroundFillSolid                       { return impl }
func (impl *BackgroundFillSolid) OptGradient() *BackgroundFillGradient                 { return nil }
func (impl *BackgroundFillSolid) OptFreeformGradient() *BackgroundFillFreeformGradient { return nil }

func (impl *BackgroundFillGradient) OptSolid() *BackgroundFillSolid                       { return nil }
func (impl *BackgroundFillGradient) OptGradient() *BackgroundFillGradient                 { return impl }
func (impl *BackgroundFillGradient) OptFreeformGradient() *BackgroundFillFreeformGradient { return nil }

func (impl *BackgroundFillFreeformGradient) OptSolid() *BackgroundFillSolid       { return nil }
func (impl *BackgroundFillFreeformGradient) OptGradient() *BackgroundFillGradient { return nil }
func (impl *BackgroundFillFreeformGradient) OptFreeformGradient() *BackgroundFillFreeformGradient {
	return impl
}

// BackgroundFillFreeformGradient The background is a freeform gradient that rotates after every message in the chat.
type BackgroundFillFreeformGradient struct {
	// Type of the background fill, always "freeform_gradient"
	Type string `json:"type" default:"freeform_gradient"`
	// A list of the 3 or 4 base colors that are used to generate the freeform gradient in the RGB24 format
	Colors []int64 `json:"colors"`
}

// BackgroundFillGradient The background is a gradient fill.
type BackgroundFillGradient struct {
	// Type of the background fill, always "gradient"
	Type string `json:"type" default:"gradient"`
	// Top color of the gradient in the RGB24 format
	TopColor int64 `json:"top_color"`
	// Bottom color of the gradient in the RGB24 format
	BottomColor int64 `json:"bottom_color"`
	// Clockwise rotation angle of the background fill in degrees; 0-359
	RotationAngle int64 `json:"rotation_angle"`
}

// BackgroundFillSolid The background is filled using the selected color.
type BackgroundFillSolid struct {
	// Type of the background fill, always "solid"
	Type string `json:"type" default:"solid"`
	// The color of the background fill in the RGB24 format
	Color int64 `json:"color"`
}

// BackgroundType This object describes the type of a background. Currently, it can be one of
// - BackgroundTypeFill
// - BackgroundTypeWallpaper
// - BackgroundTypePattern
// - BackgroundTypeChatTheme
type BackgroundType interface {
	OptFill() *BackgroundTypeFill
	OptWallpaper() *BackgroundTypeWallpaper
	OptPattern() *BackgroundTypePattern
	OptChatTheme() *BackgroundTypeChatTheme
}

var (
	_ BackgroundType = &BackgroundTypeFill{}
	_ BackgroundType = &BackgroundTypeWallpaper{}
	_ BackgroundType = &BackgroundTypePattern{}
	_ BackgroundType = &BackgroundTypeChatTheme{}
)

func (impl *BackgroundTypeFill) OptFill() *BackgroundTypeFill           { return impl }
func (impl *BackgroundTypeFill) OptWallpaper() *BackgroundTypeWallpaper { return nil }
func (impl *BackgroundTypeFill) OptPattern() *BackgroundTypePattern     { return nil }
func (impl *BackgroundTypeFill) OptChatTheme() *BackgroundTypeChatTheme { return nil }

func (impl *BackgroundTypeWallpaper) OptFill() *BackgroundTypeFill           { return nil }
func (impl *BackgroundTypeWallpaper) OptWallpaper() *BackgroundTypeWallpaper { return impl }
func (impl *BackgroundTypeWallpaper) OptPattern() *BackgroundTypePattern     { return nil }
func (impl *BackgroundTypeWallpaper) OptChatTheme() *BackgroundTypeChatTheme { return nil }

func (impl *BackgroundTypePattern) OptFill() *BackgroundTypeFill           { return nil }
func (impl *BackgroundTypePattern) OptWallpaper() *BackgroundTypeWallpaper { return nil }
func (impl *BackgroundTypePattern) OptPattern() *BackgroundTypePattern     { return impl }
func (impl *BackgroundTypePattern) OptChatTheme() *BackgroundTypeChatTheme { return nil }

func (impl *BackgroundTypeChatTheme) OptFill() *BackgroundTypeFill           { return nil }
func (impl *BackgroundTypeChatTheme) OptWallpaper() *BackgroundTypeWallpaper { return nil }
func (impl *BackgroundTypeChatTheme) OptPattern() *BackgroundTypePattern     { return nil }
func (impl *BackgroundTypeChatTheme) OptChatTheme() *BackgroundTypeChatTheme { return impl }

// BackgroundTypeChatTheme The background is taken directly from a built-in chat theme.
type BackgroundTypeChatTheme struct {
	// Type of the background, always "chat_theme"
	Type string `json:"type" default:"chat_theme"`
	// Name of the chat theme, which is usually an emoji
	ThemeName string `json:"theme_name"`
}

// BackgroundTypeFill The background is automatically filled based on the selected colors.
type BackgroundTypeFill struct {
	// Type of the background, always "fill"
	Type string `json:"type" default:"fill"`
	// The background fill
	Fill BackgroundFill `json:"fill"`
	// Dimming of the background in dark themes, as a percentage; 0-100
	DarkThemeDimming int64 `json:"dark_theme_dimming"`
}

// BackgroundTypePattern The background is a PNG or TGV (gzipped subset of SVG with MIME type "application/x-tgwallpattern") pattern to be combined with the background fill chosen by the user.
type BackgroundTypePattern struct {
	// Type of the background, always "pattern"
	Type string `json:"type" default:"pattern"`
	// Document with the pattern
	Document *TelegramDocument `json:"document"`
	// The background fill that is combined with the pattern
	Fill BackgroundFill `json:"fill"`
	// Intensity of the pattern when it is shown above the filled background; 0-100
	Intensity int64 `json:"intensity"`
	// Optional. True, if the background fill must be applied only to the pattern itself.
	// All other pixels are black in this case. For dark themes only
	IsInverted bool `json:"is_inverted,omitempty"`
	// Optional. True, if the background moves slightly when the device is tilted
	IsMoving bool `json:"is_moving,omitempty"`
}

// BackgroundTypeWallpaper The background is a wallpaper in the JPEG format.
type BackgroundTypeWallpaper struct {
	// Type of the background, always "wallpaper"
	Type string `json:"type" default:"wallpaper"`
	// Document with the wallpaper
	Document *TelegramDocument `json:"document"`
	// Dimming of the background in dark themes, as a percentage; 0-100
	DarkThemeDimming int64 `json:"dark_theme_dimming"`
	// Optional. True, if the wallpaper is downscaled to fit in a 450x450 square and then box-blurred with radius 12
	IsBlurred bool `json:"is_blurred,omitempty"`
	// Optional. True, if the background moves slightly when the device is tilted
	IsMoving bool `json:"is_moving,omitempty"`
}

// Birthdate Describes the birthdate of a user.
type Birthdate struct {
	// Day of the user's birth; 1-31
	Day int64 `json:"day"`
	// Month of the user's birth; 1-12
	Month int64 `json:"month"`
	// Optional. Year of the user's birth
	Year int64 `json:"year,omitempty"`
}

// BotCommand This object represents a bot command.
type BotCommand struct {
	// Text of the command; 1-32 characters. Can contain only lowercase English letters, digits and underscores.
	Command string `json:"command"`
	// Description of the command; 1-256 characters.
	Description string `json:"description"`
}

// BotCommandScope This object represents the scope to which bot commands are applied. Currently, the following 7 scopes are supported:
// - BotCommandScopeDefault
// - BotCommandScopeAllPrivateChats
// - BotCommandScopeAllGroupChats
// - BotCommandScopeAllChatAdministrators
// - BotCommandScopeChat
// - BotCommandScopeChatAdministrators
// - BotCommandScopeChatMember
type BotCommandScope interface {
	OptDefault() *BotCommandScopeDefault
	OptAllPrivateChats() *BotCommandScopeAllPrivateChats
	OptAllGroupChats() *BotCommandScopeAllGroupChats
	OptAllChatAdministrators() *BotCommandScopeAllChatAdministrators
	OptChat() *BotCommandScopeChat
	OptChatAdministrators() *BotCommandScopeChatAdministrators
	OptChatMember() *BotCommandScopeChatMember
}

var (
	_ BotCommandScope = &BotCommandScopeDefault{}
	_ BotCommandScope = &BotCommandScopeAllPrivateChats{}
	_ BotCommandScope = &BotCommandScopeAllGroupChats{}
	_ BotCommandScope = &BotCommandScopeAllChatAdministrators{}
	_ BotCommandScope = &BotCommandScopeChat{}
	_ BotCommandScope = &BotCommandScopeChatAdministrators{}
	_ BotCommandScope = &BotCommandScopeChatMember{}
)

func (impl *BotCommandScopeDefault) OptDefault() *BotCommandScopeDefault                 { return impl }
func (impl *BotCommandScopeDefault) OptAllPrivateChats() *BotCommandScopeAllPrivateChats { return nil }
func (impl *BotCommandScopeDefault) OptAllGroupChats() *BotCommandScopeAllGroupChats     { return nil }
func (impl *BotCommandScopeDefault) OptAllChatAdministrators() *BotCommandScopeAllChatAdministrators {
	return nil
}
func (impl *BotCommandScopeDefault) OptChat() *BotCommandScopeChat { return nil }
func (impl *BotCommandScopeDefault) OptChatAdministrators() *BotCommandScopeChatAdministrators {
	return nil
}
func (impl *BotCommandScopeDefault) OptChatMember() *BotCommandScopeChatMember { return nil }

func (impl *BotCommandScopeAllPrivateChats) OptDefault() *BotCommandScopeDefault { return nil }
func (impl *BotCommandScopeAllPrivateChats) OptAllPrivateChats() *BotCommandScopeAllPrivateChats {
	return impl
}
func (impl *BotCommandScopeAllPrivateChats) OptAllGroupChats() *BotCommandScopeAllGroupChats {
	return nil
}
func (impl *BotCommandScopeAllPrivateChats) OptAllChatAdministrators() *BotCommandScopeAllChatAdministrators {
	return nil
}
func (impl *BotCommandScopeAllPrivateChats) OptChat() *BotCommandScopeChat { return nil }
func (impl *BotCommandScopeAllPrivateChats) OptChatAdministrators() *BotCommandScopeChatAdministrators {
	return nil
}
func (impl *BotCommandScopeAllPrivateChats) OptChatMember() *BotCommandScopeChatMember { return nil }

func (impl *BotCommandScopeAllGroupChats) OptDefault() *BotCommandScopeDefault { return nil }
func (impl *BotCommandScopeAllGroupChats) OptAllPrivateChats() *BotCommandScopeAllPrivateChats {
	return nil
}
func (impl *BotCommandScopeAllGroupChats) OptAllGroupChats() *BotCommandScopeAllGroupChats {
	return impl
}
func (impl *BotCommandScopeAllGroupChats) OptAllChatAdministrators() *BotCommandScopeAllChatAdministrators {
	return nil
}
func (impl *BotCommandScopeAllGroupChats) OptChat() *BotCommandScopeChat { return nil }
func (impl *BotCommandScopeAllGroupChats) OptChatAdministrators() *BotCommandScopeChatAdministrators {
	return nil
}
func (impl *BotCommandScopeAllGroupChats) OptChatMember() *BotCommandScopeChatMember { return nil }

func (impl *BotCommandScopeAllChatAdministrators) OptDefault() *BotCommandScopeDefault { return nil }
func (impl *BotCommandScopeAllChatAdministrators) OptAllPrivateChats() *BotCommandScopeAllPrivateChats {
	return nil
}
func (impl *BotCommandScopeAllChatAdministrators) OptAllGroupChats() *BotCommandScopeAllGroupChats {
	return nil
}
func (impl *BotCommandScopeAllChatAdministrators) OptAllChatAdministrators() *BotCommandScopeAllChatAdministrators {
	return impl
}
func (impl *BotCommandScopeAllChatAdministrators) OptChat() *BotCommandScopeChat { return nil }
func (impl *BotCommandScopeAllChatAdministrators) OptChatAdministrators() *BotCommandScopeChatAdministrators {
	return nil
}
func (impl *BotCommandScopeAllChatAdministrators) OptChatMember() *BotCommandScopeChatMember {
	return nil
}

func (impl *BotCommandScopeChat) OptDefault() *BotCommandScopeDefault                 { return nil }
func (impl *BotCommandScopeChat) OptAllPrivateChats() *BotCommandScopeAllPrivateChats { return nil }
func (impl *BotCommandScopeChat) OptAllGroupChats() *BotCommandScopeAllGroupChats     { return nil }
func (impl *BotCommandScopeChat) OptAllChatAdministrators() *BotCommandScopeAllChatAdministrators {
	return nil
}
func (impl *BotCommandScopeChat) OptChat() *BotCommandScopeChat { return impl }
func (impl *BotCommandScopeChat) OptChatAdministrators() *BotCommandScopeChatAdministrators {
	return nil
}
func (impl *BotCommandScopeChat) OptChatMember() *BotCommandScopeChatMember { return nil }

func (impl *BotCommandScopeChatAdministrators) OptDefault() *BotCommandScopeDefault { return nil }
func (impl *BotCommandScopeChatAdministrators) OptAllPrivateChats() *BotCommandScopeAllPrivateChats {
	return nil
}
func (impl *BotCommandScopeChatAdministrators) OptAllGroupChats() *BotCommandScopeAllGroupChats {
	return nil
}
func (impl *BotCommandScopeChatAdministrators) OptAllChatAdministrators() *BotCommandScopeAllChatAdministrators {
	return nil
}
func (impl *BotCommandScopeChatAdministrators) OptChat() *BotCommandScopeChat { return nil }
func (impl *BotCommandScopeChatAdministrators) OptChatAdministrators() *BotCommandScopeChatAdministrators {
	return impl
}
func (impl *BotCommandScopeChatAdministrators) OptChatMember() *BotCommandScopeChatMember { return nil }

func (impl *BotCommandScopeChatMember) OptDefault() *BotCommandScopeDefault { return nil }
func (impl *BotCommandScopeChatMember) OptAllPrivateChats() *BotCommandScopeAllPrivateChats {
	return nil
}
func (impl *BotCommandScopeChatMember) OptAllGroupChats() *BotCommandScopeAllGroupChats { return nil }
func (impl *BotCommandScopeChatMember) OptAllChatAdministrators() *BotCommandScopeAllChatAdministrators {
	return nil
}
func (impl *BotCommandScopeChatMember) OptChat() *BotCommandScopeChat { return nil }
func (impl *BotCommandScopeChatMember) OptChatAdministrators() *BotCommandScopeChatAdministrators {
	return nil
}
func (impl *BotCommandScopeChatMember) OptChatMember() *BotCommandScopeChatMember { return impl }

// BotCommandScopeAllChatAdministrators Represents the scope of bot commands, covering all group and supergroup chat administrators.
type BotCommandScopeAllChatAdministrators struct {
	// Scope type, must be all_chat_administrators
	Type string `json:"type" default:"all_chat_administrators"`
}

// BotCommandScopeAllGroupChats Represents the scope of bot commands, covering all group and supergroup chats.
type BotCommandScopeAllGroupChats struct {
	// Scope type, must be all_group_chats
	Type string `json:"type" default:"all_group_chats"`
}

// BotCommandScopeAllPrivateChats Represents the scope of bot commands, covering all private chats.
type BotCommandScopeAllPrivateChats struct {
	// Scope type, must be all_private_chats
	Type string `json:"type" default:"all_private_chats"`
}

// BotCommandScopeChat Represents the scope of bot commands, covering a specific chat.
type BotCommandScopeChat struct {
	// Scope type, must be chat
	Type string `json:"type" default:"chat"`
	// Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	// >> either: String
	ChatId int64 `json:"chat_id"`
}

// BotCommandScopeChatAdministrators Represents the scope of bot commands, covering all administrators of a specific group or supergroup chat.
type BotCommandScopeChatAdministrators struct {
	// Scope type, must be chat_administrators
	Type string `json:"type" default:"chat_administrators"`
	// Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	// >> either: String
	ChatId int64 `json:"chat_id"`
}

// BotCommandScopeChatMember Represents the scope of bot commands, covering a specific member of a group or supergroup chat.
type BotCommandScopeChatMember struct {
	// Scope type, must be chat_member
	Type string `json:"type" default:"chat_member"`
	// Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	// >> either: String
	ChatId int64 `json:"chat_id"`
	// Unique identifier of the target user
	UserId int64 `json:"user_id"`
}

// BotCommandScopeDefault Represents the default scope of bot commands.
// Default commands are used if no commands with a narrower scope are specified for the user.
type BotCommandScopeDefault struct {
	// Scope type, must be default
	Type string `json:"type"`
}

// BotDescription This object represents the bot's description.
type BotDescription struct {
	// The bot's description
	Description string `json:"description"`
}

// BotName This object represents the bot's name.
type BotName struct {
	// The bot's name
	Name string `json:"name"`
}

// BotShortDescription This object represents the bot's short description.
type BotShortDescription struct {
	// The bot's short description
	ShortDescription string `json:"short_description"`
}

// BusinessConnection Describes the connection of the bot with a business account.
type BusinessConnection struct {
	// Unique identifier of the business connection
	Id string `json:"id"`
	// Business account user that created the business connection
	User *User `json:"user"`
	// Identifier of a private chat with the user who created the business connection.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier.
	UserChatId int64 `json:"user_chat_id"`
	// Date the connection was established in Unix time
	Date int64 `json:"date"`
	// True, if the bot can act on behalf of the business account in chats that were active in the last 24 hours
	CanReply bool `json:"can_reply"`
	// True, if the connection is active
	IsEnabled bool `json:"is_enabled"`
}

// BusinessIntro Contains information about the start page settings of a Telegram Business account.
type BusinessIntro struct {
	// Optional. Title text of the business intro
	Title string `json:"title,omitempty"`
	// Optional. Message text of the business intro
	Message string `json:"message,omitempty"`
	// Optional. Sticker of the business intro
	Sticker *Sticker `json:"sticker,omitempty"`
}

// BusinessLocation Contains information about the location of a Telegram Business account.
type BusinessLocation struct {
	// Address of the business
	Address string `json:"address"`
	// Optional. Location of the business
	Location *Location `json:"location,omitempty"`
}

// BusinessMessagesDeleted This object is received when messages are deleted from a connected business account.
type BusinessMessagesDeleted struct {
	// Unique identifier of the business connection
	BusinessConnectionId string `json:"business_connection_id"`
	// Information about a chat in the business account. The bot may not have access to the chat or the corresponding user.
	Chat *Chat `json:"chat"`
	// The list of identifiers of deleted messages in the chat of the business account
	MessageIds []int64 `json:"message_ids"`
}

// BusinessOpeningHours Describes the opening hours of a business.
type BusinessOpeningHours struct {
	// Unique name of the time zone for which the opening hours are defined
	TimeZoneName string `json:"time_zone_name"`
	// List of time intervals describing business opening hours
	OpeningHours []*BusinessOpeningHoursInterval `json:"opening_hours"`
}

// BusinessOpeningHoursInterval Describes an interval of time during which a business is open.
type BusinessOpeningHoursInterval struct {
	// The minute's sequence number in a week, starting on Monday, marking the start of the time interval during which the business is open; 0 - 7 * 24 * 60
	OpeningMinute int64 `json:"opening_minute"`
	// The minute's sequence number in a week, starting on Monday, marking the end of the time interval during which the business is open; 0 - 8 * 24 * 60
	ClosingMinute int64 `json:"closing_minute"`
}

// CallbackGame A placeholder, currently holds no information. Use BotFather to set up your game.
type CallbackGame struct {
}

// CallbackQuery This object represents an incoming callback query from a callback button in an inline keyboard.
// If the button that originated the query was attached to a message sent by the bot, the field message will be present.
// If the button was attached to a message sent via the bot (in inline mode), the field inline_message_id will be present.
// Exactly one of the fields data or game_short_name will be present.
type CallbackQuery struct {
	// Unique identifier for this query
	Id string `json:"id"`
	// Sender
	From *User `json:"from"`
	// Optional. Message sent by the bot with the callback button that originated the query
	Message *Message `json:"message,omitempty"`
	// Optional. Identifier of the message sent via the bot in inline mode, that originated the query.
	InlineMessageId string `json:"inline_message_id,omitempty"`
	// Global identifier, uniquely corresponding to the chat to which the message with the callback button was sent.
	// Useful for high scores in games.
	ChatInstance string `json:"chat_instance"`
	// Optional. Data associated with the callback button.
	// Be aware that the message originated the query can contain no callback buttons with this data.
	Data string `json:"data,omitempty"`
	// Optional. Short name of a Game to be returned, serves as the unique identifier for the game
	GameShortName string `json:"game_short_name,omitempty"`
}

// Chat This object represents a chat.
type Chat struct {
	// Unique identifier for this chat.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	Id int64 `json:"id"`
	// Type of the chat, can be either "private", "group", "supergroup" or "channel"
	Type string `json:"type"`
	// Optional. Title, for supergroups, channels and group chats
	Title string `json:"title,omitempty"`
	// Optional. Username, for private chats, supergroups and channels if available
	Username string `json:"username,omitempty"`
	// Optional. First name of the other party in a private chat
	FirstName string `json:"first_name,omitempty"`
	// Optional. Last name of the other party in a private chat
	LastName string `json:"last_name,omitempty"`
	// Optional. True, if the supergroup chat is a forum (has topics enabled)
	IsForum bool `json:"is_forum,omitempty"`
}

// ChatAdministratorRights Represents the rights of an administrator in a chat.
type ChatAdministratorRights struct {
	// True, if the user's presence in the chat is hidden
	IsAnonymous bool `json:"is_anonymous"`
	// True, if the administrator can access the chat event log, get boost list, see hidden supergroup and channel members, report spam messages and ignore slow mode.
	// Implied by any other administrator privilege.
	CanManageChat bool `json:"can_manage_chat"`
	// True, if the administrator can delete messages of other users
	CanDeleteMessages bool `json:"can_delete_messages"`
	// True, if the administrator can manage video chats
	CanManageVideoChats bool `json:"can_manage_video_chats"`
	// True, if the administrator can restrict, ban or unban chat members, or access supergroup statistics
	CanRestrictMembers bool `json:"can_restrict_members"`
	// True, if the administrator can add new administrators with a subset of their own privileges or demote administrators that they have promoted, directly or indirectly (promoted by administrators that were appointed by the user)
	CanPromoteMembers bool `json:"can_promote_members"`
	// True, if the user is allowed to change the chat title, photo and other settings
	CanChangeInfo bool `json:"can_change_info"`
	// True, if the user is allowed to invite new users to the chat
	CanInviteUsers bool `json:"can_invite_users"`
	// True, if the administrator can post stories to the chat
	CanPostStories bool `json:"can_post_stories"`
	// True, if the administrator can edit stories posted by other users, post stories to the chat page, pin chat stories, and access the chat's story archive
	CanEditStories bool `json:"can_edit_stories"`
	// True, if the administrator can delete stories posted by other users
	CanDeleteStories bool `json:"can_delete_stories"`
	// Optional.
	// True, if the administrator can post messages in the channel, or access channel statistics; for channels only
	CanPostMessages bool `json:"can_post_messages,omitempty"`
	// Optional. True, if the administrator can edit messages of other users and can pin messages; for channels only
	CanEditMessages bool `json:"can_edit_messages,omitempty"`
	// Optional. True, if the user is allowed to pin messages; for groups and supergroups only
	CanPinMessages bool `json:"can_pin_messages,omitempty"`
	// Optional. True, if the user is allowed to create, rename, close, and reopen forum topics; for supergroups only
	CanManageTopics bool `json:"can_manage_topics,omitempty"`
}

// ChatBackground This object represents a chat background.
type ChatBackground struct {
	// Type of the background
	Type BackgroundType `json:"type"`
}

// ChatBoost This object contains information about a chat boost.
type ChatBoost struct {
	// Unique identifier of the boost
	BoostId string `json:"boost_id"`
	// Point in time (Unix timestamp) when the chat was boosted
	AddDate int64 `json:"add_date"`
	// Point in time (Unix timestamp) when the boost will automatically expire, unless the booster's Telegram Premium subscription is prolonged
	ExpirationDate int64 `json:"expiration_date"`
	// Source of the added boost
	Source ChatBoostSource `json:"source"`
}

// ChatBoostAdded This object represents a service message about a user boosting a chat.
type ChatBoostAdded struct {
	// Number of boosts added by the user
	BoostCount int64 `json:"boost_count"`
}

// ChatBoostRemoved This object represents a boost removed from a chat.
type ChatBoostRemoved struct {
	// Chat which was boosted
	Chat *Chat `json:"chat"`
	// Unique identifier of the boost
	BoostId string `json:"boost_id"`
	// Point in time (Unix timestamp) when the boost was removed
	RemoveDate int64 `json:"remove_date"`
	// Source of the removed boost
	Source ChatBoostSource `json:"source"`
}

// ChatBoostSource This object describes the source of a chat boost. It can be one of
// - ChatBoostSourcePremium
// - ChatBoostSourceGiftCode
// - ChatBoostSourceGiveaway
type ChatBoostSource interface {
	OptPremium() *ChatBoostSourcePremium
	OptGiftCode() *ChatBoostSourceGiftCode
	OptGiveaway() *ChatBoostSourceGiveaway
}

var (
	_ ChatBoostSource = &ChatBoostSourcePremium{}
	_ ChatBoostSource = &ChatBoostSourceGiftCode{}
	_ ChatBoostSource = &ChatBoostSourceGiveaway{}
)

func (impl *ChatBoostSourcePremium) OptPremium() *ChatBoostSourcePremium   { return impl }
func (impl *ChatBoostSourcePremium) OptGiftCode() *ChatBoostSourceGiftCode { return nil }
func (impl *ChatBoostSourcePremium) OptGiveaway() *ChatBoostSourceGiveaway { return nil }

func (impl *ChatBoostSourceGiftCode) OptPremium() *ChatBoostSourcePremium   { return nil }
func (impl *ChatBoostSourceGiftCode) OptGiftCode() *ChatBoostSourceGiftCode { return impl }
func (impl *ChatBoostSourceGiftCode) OptGiveaway() *ChatBoostSourceGiveaway { return nil }

func (impl *ChatBoostSourceGiveaway) OptPremium() *ChatBoostSourcePremium   { return nil }
func (impl *ChatBoostSourceGiveaway) OptGiftCode() *ChatBoostSourceGiftCode { return nil }
func (impl *ChatBoostSourceGiveaway) OptGiveaway() *ChatBoostSourceGiveaway { return impl }

// ChatBoostSourceGiftCode The boost was obtained by the creation of Telegram Premium gift codes to boost a chat.
// Each such code boosts the chat 4 times for the duration of the corresponding Telegram Premium subscription.
type ChatBoostSourceGiftCode struct {
	// Source of the boost, always "gift_code"
	Source string `json:"source" default:"gift_code"`
	// User for which the gift code was created
	User *User `json:"user"`
}

// ChatBoostSourceGiveaway The boost was obtained by the creation of a Telegram Premium or a Telegram Star giveaway.
// This boosts the chat 4 times for the duration of the corresponding Telegram Premium subscription for Telegram Premium giveaways and prize_star_count / 500 times for one year for Telegram Star giveaways.
type ChatBoostSourceGiveaway struct {
	// Source of the boost, always "giveaway"
	Source string `json:"source" default:"giveaway"`
	// Identifier of a message in the chat with the giveaway; the message could have been deleted already.
	// May be 0 if the message isn't sent yet.
	GiveawayMessageId int64 `json:"giveaway_message_id"`
	// Optional. User that won the prize in the giveaway if any; for Telegram Premium giveaways only
	User *User `json:"user,omitempty"`
	// Optional. The number of Telegram Stars to be split between giveaway winners; for Telegram Star giveaways only
	PrizeStarCount int64 `json:"prize_star_count,omitempty"`
	// Optional. True, if the giveaway was completed, but there was no user to win the prize
	IsUnclaimed bool `json:"is_unclaimed,omitempty"`
}

// ChatBoostSourcePremium The boost was obtained by subscribing to Telegram Premium or by gifting a Telegram Premium subscription to another user.
type ChatBoostSourcePremium struct {
	// Source of the boost, always "premium"
	Source string `json:"source" default:"premium"`
	// User that boosted the chat
	User *User `json:"user"`
}

// ChatBoostUpdated This object represents a boost added to a chat or changed.
type ChatBoostUpdated struct {
	// Chat which was boosted
	Chat *Chat `json:"chat"`
	// Information about the chat boost
	Boost *ChatBoost `json:"boost"`
}

// ChatFullInfo This object contains full information about a chat.
type ChatFullInfo struct {
	// Unique identifier for this chat.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	Id int64 `json:"id"`
	// Type of the chat, can be either "private", "group", "supergroup" or "channel"
	Type string `json:"type"`
	// Optional. Title, for supergroups, channels and group chats
	Title string `json:"title,omitempty"`
	// Optional. Username, for private chats, supergroups and channels if available
	Username string `json:"username,omitempty"`
	// Optional. First name of the other party in a private chat
	FirstName string `json:"first_name,omitempty"`
	// Optional. Last name of the other party in a private chat
	LastName string `json:"last_name,omitempty"`
	// Optional. True, if the supergroup chat is a forum (has topics enabled)
	IsForum bool `json:"is_forum,omitempty"`
	// Identifier of the accent color for the chat name and backgrounds of the chat photo, reply header, and link preview.
	// See accent colors for more details.
	AccentColorId int64 `json:"accent_color_id"`
	// The maximum number of reactions that can be set on a message in the chat
	MaxReactionCount int64 `json:"max_reaction_count"`
	// Optional. Chat photo
	Photo *ChatPhoto `json:"photo,omitempty"`
	// Optional. If non-empty, the list of all active chat usernames; for private chats, supergroups and channels
	ActiveUsernames []string `json:"active_usernames,omitempty"`
	// Optional. For private chats, the date of birth of the user
	Birthdate *Birthdate `json:"birthdate,omitempty"`
	// Optional. For private chats with business accounts, the intro of the business
	BusinessIntro *BusinessIntro `json:"business_intro,omitempty"`
	// Optional. For private chats with business accounts, the location of the business
	BusinessLocation *BusinessLocation `json:"business_location,omitempty"`
	// Optional. For private chats with business accounts, the opening hours of the business
	BusinessOpeningHours *BusinessOpeningHours `json:"business_opening_hours,omitempty"`
	// Optional. For private chats, the personal channel of the user
	PersonalChat *Chat `json:"personal_chat,omitempty"`
	// Optional. List of available reactions allowed in the chat. If omitted, then all emoji reactions are allowed.
	AvailableReactions []ReactionType `json:"available_reactions,omitempty"`
	// Optional. Custom emoji identifier of the emoji chosen by the chat for the reply header and link preview background
	BackgroundCustomEmojiId string `json:"background_custom_emoji_id,omitempty"`
	// Optional. Identifier of the accent color for the chat's profile background.
	// See profile accent colors for more details.
	ProfileAccentColorId int64 `json:"profile_accent_color_id,omitempty"`
	// Optional. Custom emoji identifier of the emoji chosen by the chat for its profile background
	ProfileBackgroundCustomEmojiId string `json:"profile_background_custom_emoji_id,omitempty"`
	// Optional. Custom emoji identifier of the emoji status of the chat or the other party in a private chat
	EmojiStatusCustomEmojiId string `json:"emoji_status_custom_emoji_id,omitempty"`
	// Optional. Expiration date of the emoji status of the chat or the other party in a private chat, in Unix time, if any
	EmojiStatusExpirationDate int64 `json:"emoji_status_expiration_date,omitempty"`
	// Optional. Bio of the other party in a private chat
	Bio string `json:"bio,omitempty"`
	// Optional.
	// True, if privacy settings of the other party in the private chat allows to use tg://user?id=<user_id> links only in chats with the user
	HasPrivateForwards bool `json:"has_private_forwards,omitempty"`
	// Optional.
	// True, if the privacy settings of the other party restrict sending voice and video note messages in the private chat
	HasRestrictedVoiceAndVideoMessages bool `json:"has_restricted_voice_and_video_messages,omitempty"`
	// Optional. True, if users need to join the supergroup before they can send messages
	JoinToSendMessages bool `json:"join_to_send_messages,omitempty"`
	// Optional.
	// True, if all users directly joining the supergroup without using an invite link need to be approved by supergroup administrators
	JoinByRequest bool `json:"join_by_request,omitempty"`
	// Optional. Description, for groups, supergroups and channel chats
	Description string `json:"description,omitempty"`
	// Optional. Primary invite link, for groups, supergroups and channel chats
	InviteLink string `json:"invite_link,omitempty"`
	// Optional. The most recent pinned message (by sending date)
	PinnedMessage *Message `json:"pinned_message,omitempty"`
	// Optional. Default chat member permissions, for groups and supergroups
	Permissions *ChatPermissions `json:"permissions,omitempty"`
	// Optional. True, if paid media messages can be sent or forwarded to the channel chat.
	// The field is available only for channel chats.
	CanSendPaidMedia bool `json:"can_send_paid_media,omitempty"`
	// Optional.
	// For supergroups, the minimum allowed delay between consecutive messages sent by each unprivileged user; in seconds
	SlowModeDelay int64 `json:"slow_mode_delay,omitempty"`
	// Optional.
	// For supergroups, the minimum number of boosts that a non-administrator user needs to add in order to ignore slow mode and chat permissions
	UnrestrictBoostCount int64 `json:"unrestrict_boost_count,omitempty"`
	// Optional. The time after which all messages sent to the chat will be automatically deleted; in seconds
	MessageAutoDeleteTime int64 `json:"message_auto_delete_time,omitempty"`
	// Optional. True, if aggressive anti-spam checks are enabled in the supergroup.
	// The field is only available to chat administrators.
	HasAggressiveAntiSpamEnabled bool `json:"has_aggressive_anti_spam_enabled,omitempty"`
	// Optional. True, if non-administrators can only get the list of bots and administrators in the chat
	HasHiddenMembers bool `json:"has_hidden_members,omitempty"`
	// Optional. True, if messages from the chat can't be forwarded to other chats
	HasProtectedContent bool `json:"has_protected_content,omitempty"`
	// Optional. True, if new chat members will have access to old messages; available only to chat administrators
	HasVisibleHistory bool `json:"has_visible_history,omitempty"`
	// Optional. For supergroups, name of the group sticker set
	StickerSetName string `json:"sticker_set_name,omitempty"`
	// Optional. True, if the bot can change the group sticker set
	CanSetStickerSet bool `json:"can_set_sticker_set,omitempty"`
	// Optional. For supergroups, the name of the group's custom emoji sticker set.
	// Custom emoji from this set can be used by all users and bots in the group.
	CustomEmojiStickerSetName string `json:"custom_emoji_sticker_set_name,omitempty"`
	// Optional. Unique identifier for the linked chat, i.e.
	// the discussion group identifier for a channel and vice versa; for supergroups and channel chats.
	// This identifier may be greater than 32 bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it is smaller than 52 bits, so a signed 64 bit integer or double-precision float type are safe for storing this identifier.
	LinkedChatId int64 `json:"linked_chat_id,omitempty"`
	// Optional. For supergroups, the location to which the supergroup is connected
	Location *ChatLocation `json:"location,omitempty"`
}

// ChatInviteLink Represents an invite link for a chat.
type ChatInviteLink struct {
	// The invite link.
	// If the link was created by another chat administrator, then the second part of the link will be replaced with "...".
	InviteLink string `json:"invite_link"`
	// Creator of the link
	Creator *User `json:"creator"`
	// True, if users joining the chat via the link need to be approved by chat administrators
	CreatesJoinRequest bool `json:"creates_join_request"`
	// True, if the link is primary
	IsPrimary bool `json:"is_primary"`
	// True, if the link is revoked
	IsRevoked bool `json:"is_revoked"`
	// Optional. Invite link name
	Name string `json:"name,omitempty"`
	// Optional. Point in time (Unix timestamp) when the link will expire or has been expired
	ExpireDate int64 `json:"expire_date,omitempty"`
	// Optional.
	// The maximum number of users that can be members of the chat simultaneously after joining the chat via this invite link; 1-99999
	MemberLimit int64 `json:"member_limit,omitempty"`
	// Optional. Number of pending join requests created using this link
	PendingJoinRequestCount int64 `json:"pending_join_request_count,omitempty"`
	// Optional. The number of seconds the subscription will be active for before the next payment
	SubscriptionPeriod int64 `json:"subscription_period,omitempty"`
	// Optional.
	// The amount of Telegram Stars a user must pay initially and after each subsequent subscription period to be a member of the chat using the link
	SubscriptionPrice int64 `json:"subscription_price,omitempty"`
}

// ChatJoinRequest Represents a join request sent to a chat.
type ChatJoinRequest struct {
	// Chat to which the request was sent
	Chat *Chat `json:"chat"`
	// User that sent the join request
	From *User `json:"from"`
	// Identifier of a private chat with the user who sent the join request.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier.
	// The bot can use this identifier for 5 minutes to send messages until the join request is processed, assuming no other administrator contacted the user.
	UserChatId int64 `json:"user_chat_id"`
	// Date the request was sent in Unix time
	Date int64 `json:"date"`
	// Optional. Bio of the user.
	Bio string `json:"bio,omitempty"`
	// Optional. Chat invite link that was used by the user to send the join request
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

// ChatLocation Represents a location to which a chat is connected.
type ChatLocation struct {
	// The location to which the supergroup is connected. Can't be a live location.
	Location *Location `json:"location"`
	// Location address; 1-64 characters, as defined by the chat owner
	Address string `json:"address"`
}

// ChatMember This object contains information about one member of a chat. Currently, the following 6 types of chat members are supported:
// - ChatMemberOwner
// - ChatMemberAdministrator
// - ChatMemberMember
// - ChatMemberRestricted
// - ChatMemberLeft
// - ChatMemberBanned
type ChatMember interface {
	OptOwner() *ChatMemberOwner
	OptAdministrator() *ChatMemberAdministrator
	OptMember() *ChatMemberMember
	OptRestricted() *ChatMemberRestricted
	OptLeft() *ChatMemberLeft
	OptBanned() *ChatMemberBanned
}

var (
	_ ChatMember = &ChatMemberOwner{}
	_ ChatMember = &ChatMemberAdministrator{}
	_ ChatMember = &ChatMemberMember{}
	_ ChatMember = &ChatMemberRestricted{}
	_ ChatMember = &ChatMemberLeft{}
	_ ChatMember = &ChatMemberBanned{}
)

func (impl *ChatMemberOwner) OptOwner() *ChatMemberOwner                 { return impl }
func (impl *ChatMemberOwner) OptAdministrator() *ChatMemberAdministrator { return nil }
func (impl *ChatMemberOwner) OptMember() *ChatMemberMember               { return nil }
func (impl *ChatMemberOwner) OptRestricted() *ChatMemberRestricted       { return nil }
func (impl *ChatMemberOwner) OptLeft() *ChatMemberLeft                   { return nil }
func (impl *ChatMemberOwner) OptBanned() *ChatMemberBanned               { return nil }

func (impl *ChatMemberAdministrator) OptOwner() *ChatMemberOwner                 { return nil }
func (impl *ChatMemberAdministrator) OptAdministrator() *ChatMemberAdministrator { return impl }
func (impl *ChatMemberAdministrator) OptMember() *ChatMemberMember               { return nil }
func (impl *ChatMemberAdministrator) OptRestricted() *ChatMemberRestricted       { return nil }
func (impl *ChatMemberAdministrator) OptLeft() *ChatMemberLeft                   { return nil }
func (impl *ChatMemberAdministrator) OptBanned() *ChatMemberBanned               { return nil }

func (impl *ChatMemberMember) OptOwner() *ChatMemberOwner                 { return nil }
func (impl *ChatMemberMember) OptAdministrator() *ChatMemberAdministrator { return nil }
func (impl *ChatMemberMember) OptMember() *ChatMemberMember               { return impl }
func (impl *ChatMemberMember) OptRestricted() *ChatMemberRestricted       { return nil }
func (impl *ChatMemberMember) OptLeft() *ChatMemberLeft                   { return nil }
func (impl *ChatMemberMember) OptBanned() *ChatMemberBanned               { return nil }

func (impl *ChatMemberRestricted) OptOwner() *ChatMemberOwner                 { return nil }
func (impl *ChatMemberRestricted) OptAdministrator() *ChatMemberAdministrator { return nil }
func (impl *ChatMemberRestricted) OptMember() *ChatMemberMember               { return nil }
func (impl *ChatMemberRestricted) OptRestricted() *ChatMemberRestricted       { return impl }
func (impl *ChatMemberRestricted) OptLeft() *ChatMemberLeft                   { return nil }
func (impl *ChatMemberRestricted) OptBanned() *ChatMemberBanned               { return nil }

func (impl *ChatMemberLeft) OptOwner() *ChatMemberOwner                 { return nil }
func (impl *ChatMemberLeft) OptAdministrator() *ChatMemberAdministrator { return nil }
func (impl *ChatMemberLeft) OptMember() *ChatMemberMember               { return nil }
func (impl *ChatMemberLeft) OptRestricted() *ChatMemberRestricted       { return nil }
func (impl *ChatMemberLeft) OptLeft() *ChatMemberLeft                   { return impl }
func (impl *ChatMemberLeft) OptBanned() *ChatMemberBanned               { return nil }

func (impl *ChatMemberBanned) OptOwner() *ChatMemberOwner                 { return nil }
func (impl *ChatMemberBanned) OptAdministrator() *ChatMemberAdministrator { return nil }
func (impl *ChatMemberBanned) OptMember() *ChatMemberMember               { return nil }
func (impl *ChatMemberBanned) OptRestricted() *ChatMemberRestricted       { return nil }
func (impl *ChatMemberBanned) OptLeft() *ChatMemberLeft                   { return nil }
func (impl *ChatMemberBanned) OptBanned() *ChatMemberBanned               { return impl }

// ChatMemberAdministrator Represents a chat member that has some additional privileges.
type ChatMemberAdministrator struct {
	// The member's status in the chat, always "administrator"
	Status string `json:"status" default:"administrator"`
	// Information about the user
	User *User `json:"user"`
	// True, if the bot is allowed to edit administrator privileges of that user
	CanBeEdited bool `json:"can_be_edited"`
	// True, if the user's presence in the chat is hidden
	IsAnonymous bool `json:"is_anonymous"`
	// True, if the administrator can access the chat event log, get boost list, see hidden supergroup and channel members, report spam messages and ignore slow mode.
	// Implied by any other administrator privilege.
	CanManageChat bool `json:"can_manage_chat"`
	// True, if the administrator can delete messages of other users
	CanDeleteMessages bool `json:"can_delete_messages"`
	// True, if the administrator can manage video chats
	CanManageVideoChats bool `json:"can_manage_video_chats"`
	// True, if the administrator can restrict, ban or unban chat members, or access supergroup statistics
	CanRestrictMembers bool `json:"can_restrict_members"`
	// True, if the administrator can add new administrators with a subset of their own privileges or demote administrators that they have promoted, directly or indirectly (promoted by administrators that were appointed by the user)
	CanPromoteMembers bool `json:"can_promote_members"`
	// True, if the user is allowed to change the chat title, photo and other settings
	CanChangeInfo bool `json:"can_change_info"`
	// True, if the user is allowed to invite new users to the chat
	CanInviteUsers bool `json:"can_invite_users"`
	// True, if the administrator can post stories to the chat
	CanPostStories bool `json:"can_post_stories"`
	// True, if the administrator can edit stories posted by other users, post stories to the chat page, pin chat stories, and access the chat's story archive
	CanEditStories bool `json:"can_edit_stories"`
	// True, if the administrator can delete stories posted by other users
	CanDeleteStories bool `json:"can_delete_stories"`
	// Optional.
	// True, if the administrator can post messages in the channel, or access channel statistics; for channels only
	CanPostMessages bool `json:"can_post_messages,omitempty"`
	// Optional. True, if the administrator can edit messages of other users and can pin messages; for channels only
	CanEditMessages bool `json:"can_edit_messages,omitempty"`
	// Optional. True, if the user is allowed to pin messages; for groups and supergroups only
	CanPinMessages bool `json:"can_pin_messages,omitempty"`
	// Optional. True, if the user is allowed to create, rename, close, and reopen forum topics; for supergroups only
	CanManageTopics bool `json:"can_manage_topics,omitempty"`
	// Optional. Custom title for this user
	CustomTitle string `json:"custom_title,omitempty"`
}

// ChatMemberBanned Represents a chat member that was banned in the chat and can't return to the chat or view chat messages.
type ChatMemberBanned struct {
	// The member's status in the chat, always "kicked"
	Status string `json:"status" default:"kicked"`
	// Information about the user
	User *User `json:"user"`
	// Date when restrictions will be lifted for this user; Unix time. If 0, then the user is banned forever
	UntilDate int64 `json:"until_date"`
}

// ChatMemberLeft Represents a chat member that isn't currently a member of the chat, but may join it themselves.
type ChatMemberLeft struct {
	// The member's status in the chat, always "left"
	Status string `json:"status" default:"left"`
	// Information about the user
	User *User `json:"user"`
}

// ChatMemberMember Represents a chat member that has no additional privileges or restrictions.
type ChatMemberMember struct {
	// The member's status in the chat, always "member"
	Status string `json:"status" default:"member"`
	// Information about the user
	User *User `json:"user"`
	// Optional. Date when the user's subscription will expire; Unix time
	UntilDate int64 `json:"until_date,omitempty"`
}

// ChatMemberOwner Represents a chat member that owns the chat and has all administrator privileges.
type ChatMemberOwner struct {
	// The member's status in the chat, always "creator"
	Status string `json:"status" default:"creator"`
	// Information about the user
	User *User `json:"user"`
	// True, if the user's presence in the chat is hidden
	IsAnonymous bool `json:"is_anonymous"`
	// Optional. Custom title for this user
	CustomTitle string `json:"custom_title,omitempty"`
}

// ChatMemberRestricted Represents a chat member that is under certain restrictions in the chat. Supergroups only.
type ChatMemberRestricted struct {
	// The member's status in the chat, always "restricted"
	Status string `json:"status" default:"restricted"`
	// Information about the user
	User *User `json:"user"`
	// True, if the user is a member of the chat at the moment of the request
	IsMember bool `json:"is_member"`
	// True, if the user is allowed to send text messages, contacts, giveaways, giveaway winners, invoices, locations and venues
	CanSendMessages bool `json:"can_send_messages"`
	// True, if the user is allowed to send audios
	CanSendAudios bool `json:"can_send_audios"`
	// True, if the user is allowed to send documents
	CanSendDocuments bool `json:"can_send_documents"`
	// True, if the user is allowed to send photos
	CanSendPhotos bool `json:"can_send_photos"`
	// True, if the user is allowed to send videos
	CanSendVideos bool `json:"can_send_videos"`
	// True, if the user is allowed to send video notes
	CanSendVideoNotes bool `json:"can_send_video_notes"`
	// True, if the user is allowed to send voice notes
	CanSendVoiceNotes bool `json:"can_send_voice_notes"`
	// True, if the user is allowed to send polls
	CanSendPolls bool `json:"can_send_polls"`
	// True, if the user is allowed to send animations, games, stickers and use inline bots
	CanSendOtherMessages bool `json:"can_send_other_messages"`
	// True, if the user is allowed to add web page previews to their messages
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	// True, if the user is allowed to change the chat title, photo and other settings
	CanChangeInfo bool `json:"can_change_info"`
	// True, if the user is allowed to invite new users to the chat
	CanInviteUsers bool `json:"can_invite_users"`
	// True, if the user is allowed to pin messages
	CanPinMessages bool `json:"can_pin_messages"`
	// True, if the user is allowed to create forum topics
	CanManageTopics bool `json:"can_manage_topics"`
	// Date when restrictions will be lifted for this user; Unix time. If 0, then the user is restricted forever
	UntilDate int64 `json:"until_date"`
}

// ChatMemberUpdated This object represents changes in the status of a chat member.
type ChatMemberUpdated struct {
	// Chat the user belongs to
	Chat *Chat `json:"chat"`
	// Performer of the action, which resulted in the change
	From *User `json:"from"`
	// Date the change was done in Unix time
	Date int64 `json:"date"`
	// Previous information about the chat member
	OldChatMember ChatMember `json:"old_chat_member"`
	// New information about the chat member
	NewChatMember ChatMember `json:"new_chat_member"`
	// Optional. Chat invite link, which was used by the user to join the chat; for joining by invite link events only.
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
	// Optional.
	// True, if the user joined the chat after sending a direct join request without using an invite link and being approved by an administrator
	ViaJoinRequest bool `json:"via_join_request,omitempty"`
	// Optional. True, if the user joined the chat via a chat folder invite link
	ViaChatFolderInviteLink bool `json:"via_chat_folder_invite_link,omitempty"`
}

// ChatPermissions Describes actions that a non-administrator user is allowed to take in a chat.
type ChatPermissions struct {
	// Optional.
	// True, if the user is allowed to send text messages, contacts, giveaways, giveaway winners, invoices, locations and venues
	CanSendMessages bool `json:"can_send_messages,omitempty"`
	// Optional. True, if the user is allowed to send audios
	CanSendAudios bool `json:"can_send_audios,omitempty"`
	// Optional. True, if the user is allowed to send documents
	CanSendDocuments bool `json:"can_send_documents,omitempty"`
	// Optional. True, if the user is allowed to send photos
	CanSendPhotos bool `json:"can_send_photos,omitempty"`
	// Optional. True, if the user is allowed to send videos
	CanSendVideos bool `json:"can_send_videos,omitempty"`
	// Optional. True, if the user is allowed to send video notes
	CanSendVideoNotes bool `json:"can_send_video_notes,omitempty"`
	// Optional. True, if the user is allowed to send voice notes
	CanSendVoiceNotes bool `json:"can_send_voice_notes,omitempty"`
	// Optional. True, if the user is allowed to send polls
	CanSendPolls bool `json:"can_send_polls,omitempty"`
	// Optional. True, if the user is allowed to send animations, games, stickers and use inline bots
	CanSendOtherMessages bool `json:"can_send_other_messages,omitempty"`
	// Optional. True, if the user is allowed to add web page previews to their messages
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	// Optional. True, if the user is allowed to change the chat title, photo and other settings.
	// Ignored in public supergroups
	CanChangeInfo bool `json:"can_change_info,omitempty"`
	// Optional. True, if the user is allowed to invite new users to the chat
	CanInviteUsers bool `json:"can_invite_users,omitempty"`
	// Optional. True, if the user is allowed to pin messages. Ignored in public supergroups
	CanPinMessages bool `json:"can_pin_messages,omitempty"`
	// Optional. True, if the user is allowed to create forum topics. If omitted defaults to the value of can_pin_messages
	CanManageTopics bool `json:"can_manage_topics,omitempty"`
}

// ChatPhoto This object represents a chat photo.
type ChatPhoto struct {
	// File identifier of small (160x160) chat photo.
	// This file_id can be used only for photo download and only for as long as the photo is not changed.
	SmallFileId string `json:"small_file_id"`
	// Unique file identifier of small (160x160) chat photo, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	SmallFileUniqueId string `json:"small_file_unique_id"`
	// File identifier of big (640x640) chat photo.
	// This file_id can be used only for photo download and only for as long as the photo is not changed.
	BigFileId string `json:"big_file_id"`
	// Unique file identifier of big (640x640) chat photo, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	BigFileUniqueId string `json:"big_file_unique_id"`
}

func (impl *ChatPhoto) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.BigFileId)
}
func (impl *ChatPhoto) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.BigFileId, dirAndPattern...)
}

// ChatShared This object contains information about a chat that was shared with the bot using a KeyboardButtonRequestChat button.
type ChatShared struct {
	// Identifier of the request
	RequestId int64 `json:"request_id"`
	// Identifier of the shared chat.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier.
	// The bot may not have access to the chat and could be unable to use this identifier, unless the chat is already known to the bot by some other means.
	ChatId int64 `json:"chat_id"`
	// Optional. Title of the chat, if the title was requested by the bot.
	Title string `json:"title,omitempty"`
	// Optional. Username of the chat, if the username was requested by the bot and available.
	Username string `json:"username,omitempty"`
	// Optional. Available sizes of the chat photo, if the photo was requested by the bot
	Photo TelegramPhoto `json:"photo,omitempty"`
}

// ChosenInlineResult Represents a result of an inline query that was chosen by the user and sent to their chat partner.
// Note: It is necessary to enable inline feedback via @BotFather in order to receive these objects in updates.
type ChosenInlineResult struct {
	// The unique identifier for the result that was chosen
	ResultId string `json:"result_id"`
	// The user that chose the result
	From *User `json:"from"`
	// Optional. Sender location, only for bots that require user location
	Location *Location `json:"location,omitempty"`
	// Optional. Identifier of the sent inline message.
	// Available only if there is an inline keyboard attached to the message.
	// Will be also received in callback queries and can be used to edit the message.
	InlineMessageId string `json:"inline_message_id,omitempty"`
	// The query that was used to obtain the result
	Query string `json:"query"`
}

// Contact This object represents a phone contact.
type Contact struct {
	// Contact's phone number
	PhoneNumber string `json:"phone_number"`
	// Contact's first name
	FirstName string `json:"first_name"`
	// Optional. Contact's last name
	LastName string `json:"last_name,omitempty"`
	// Optional. Contact's user identifier in Telegram.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier.
	UserId int64 `json:"user_id,omitempty"`
	// Optional. Additional data about the contact in the form of a vCard
	Vcard string `json:"vcard,omitempty"`
}

// CopyTextButton This object represents an inline keyboard button that copies specified text to the clipboard.
type CopyTextButton struct {
	// The text to be copied to the clipboard; 1-256 characters
	Text string `json:"text"`
}

// Dice This object represents an animated emoji that displays a random value.
type Dice struct {
	// Emoji on which the dice throw animation is based
	Emoji string `json:"emoji"`
	// Value of the dice, 1-6 for "🎲", "🎯" and "🎳" base emoji, 1-5 for "🏀" and "⚽" base emoji, 1-64 for "🎰" base emoji
	Value int64 `json:"value"`
}

// Document Represents a general file to be sent.
type Document struct {
	// Type of the result, must be document
	Type string `json:"type" default:"document"`
	// File to send. More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// Pass a file_id to send a file that exists on the Telegram servers (recommended), pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	Media InputFile `json:"media"`
	// Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size.
	// A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data.
	// Thumbnails can't be reused and can be only uploaded as a new file, so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	// More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// >> either: String
	Thumbnail InputFile `json:"thumbnail,omitempty"`
	// Optional. Caption of the document to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the document caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Disables automatic server-side content type detection for files uploaded using multipart/form-data.
	// Always True, if the document is sent as part of an album.
	DisableContentTypeDetection bool `json:"disable_content_type_detection,omitempty"`
	// Used for uploading media.
	InputFile InputFile `json:"-"`
}

// EncryptedCredentials Describes data required for decrypting and authenticating EncryptedPassportElement.
// See the Telegram Passport Documentation for a complete description of the data decryption and authentication processes.
type EncryptedCredentials struct {
	// Base64-encoded encrypted JSON-serialized data with unique user's payload, data hashes and secrets required for EncryptedPassportElement decryption and authentication
	Data string `json:"data"`
	// Base64-encoded data hash for data authentication
	Hash string `json:"hash"`
	// Base64-encoded secret, encrypted with the bot's public RSA key, required for data decryption
	Secret string `json:"secret"`
}

// EncryptedPassportElement Describes documents or other Telegram Passport elements shared with the bot by the user.
type EncryptedPassportElement struct {
	// Element type.
	// One of "personal_details", "passport", "driver_license", "identity_card", "internal_passport", "address", "utility_bill", "bank_statement", "rental_agreement", "passport_registration", "temporary_registration", "phone_number", "email".
	Type string `json:"type"`
	// Optional. Can be decrypted and verified using the accompanying EncryptedCredentials.
	// Base64-encoded encrypted Telegram Passport element data provided by the user; available only for "personal_details", "passport", "driver_license", "identity_card", "internal_passport" and "address" types.
	Data string `json:"data,omitempty"`
	// Optional. User's verified phone number; available only for "phone_number" type
	PhoneNumber string `json:"phone_number,omitempty"`
	// Optional. User's verified email address; available only for "email" type
	Email string `json:"email,omitempty"`
	// Optional. Files can be decrypted and verified using the accompanying EncryptedCredentials.
	// Array of encrypted files with documents provided by the user; available only for "utility_bill", "bank_statement", "rental_agreement", "passport_registration" and "temporary_registration" types.
	Files []*PassportFile `json:"files,omitempty"`
	// Optional. The file can be decrypted and verified using the accompanying EncryptedCredentials.
	// Encrypted file with the front side of the document, provided by the user; available only for "passport", "driver_license", "identity_card" and "internal_passport".
	FrontSide *PassportFile `json:"front_side,omitempty"`
	// Optional. The file can be decrypted and verified using the accompanying EncryptedCredentials.
	// Encrypted file with the reverse side of the document, provided by the user; available only for "driver_license" and "identity_card".
	ReverseSide *PassportFile `json:"reverse_side,omitempty"`
	// Optional. The file can be decrypted and verified using the accompanying EncryptedCredentials.
	// Encrypted file with the selfie of the user holding a document, provided by the user; available if requested for "passport", "driver_license", "identity_card" and "internal_passport".
	Selfie *PassportFile `json:"selfie,omitempty"`
	// Optional. Files can be decrypted and verified using the accompanying EncryptedCredentials.
	// Array of encrypted files with translated versions of documents provided by the user; available if requested for "passport", "driver_license", "identity_card", "internal_passport", "utility_bill", "bank_statement", "rental_agreement", "passport_registration" and "temporary_registration" types.
	Translation []*PassportFile `json:"translation,omitempty"`
	// Base64-encoded element hash for using in PassportElementErrorUnspecified
	Hash string `json:"hash"`
}

// ExternalReplyInfo This object contains information about a message that is being replied to, which may come from another chat or forum topic.
type ExternalReplyInfo struct {
	// Origin of the message replied to by the given message
	Origin MessageOrigin `json:"origin"`
	// Optional. Chat the original message belongs to. Available only if the chat is a supergroup or a channel.
	Chat *Chat `json:"chat,omitempty"`
	// Optional. Unique message identifier inside the original chat.
	// Available only if the original chat is a supergroup or a channel.
	MessageId int64 `json:"message_id,omitempty"`
	// Optional. Options used for link preview generation for the original message, if it is a text message
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	// Optional. Message is an animation, information about the animation
	Animation *TelegramAnimation `json:"animation,omitempty"`
	// Optional. Message is an audio file, information about the file
	Audio *TelegramAudio `json:"audio,omitempty"`
	// Optional. Message is a general file, information about the file
	Document *TelegramDocument `json:"document,omitempty"`
	// Optional. Message contains paid media; information about the paid media
	PaidMedia *PaidMediaInfo `json:"paid_media,omitempty"`
	// Optional. Message is a photo, available sizes of the photo
	Photo TelegramPhoto `json:"photo,omitempty"`
	// Optional. Message is a sticker, information about the sticker
	Sticker *Sticker `json:"sticker,omitempty"`
	// Optional. Message is a forwarded story
	Story *Story `json:"story,omitempty"`
	// Optional. Message is a video, information about the video
	Video *TelegramVideo `json:"video,omitempty"`
	// Optional. Message is a video note, information about the video message
	VideoNote *VideoNote `json:"video_note,omitempty"`
	// Optional. Message is a voice message, information about the file
	Voice *Voice `json:"voice,omitempty"`
	// Optional. True, if the message media is covered by a spoiler animation
	HasMediaSpoiler bool `json:"has_media_spoiler,omitempty"`
	// Optional. Message is a shared contact, information about the contact
	Contact *Contact `json:"contact,omitempty"`
	// Optional. Message is a dice with random value
	Dice *Dice `json:"dice,omitempty"`
	// Optional. Message is a game, information about the game. More about games: https://core.telegram.org/bots/api#games
	Game *Game `json:"game,omitempty"`
	// Optional. Message is a scheduled giveaway, information about the giveaway
	Giveaway *Giveaway `json:"giveaway,omitempty"`
	// Optional. A giveaway with public winners was completed
	GiveawayWinners *GiveawayWinners `json:"giveaway_winners,omitempty"`
	// Optional. Message is an invoice for a payment, information about the invoice.
	// More about payments: https://core.telegram.org/bots/api#payments
	Invoice *Invoice `json:"invoice,omitempty"`
	// Optional. Message is a shared location, information about the location
	Location *Location `json:"location,omitempty"`
	// Optional. Message is a native poll, information about the poll
	Poll *Poll `json:"poll,omitempty"`
	// Optional. Message is a venue, information about the venue
	Venue *Venue `json:"venue,omitempty"`
}

// File This object represents a file ready to be downloaded.
// The file can be downloaded via the link https://api.telegram.org/file/bot<token>/<file_path>.
// It is guaranteed that the link will be valid for at least 1 hour.
// When the link expires, a new one can be requested by calling getFile.
type File struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Optional. File size in bytes.
	// It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
	FileSize int64 `json:"file_size,omitempty"`
	// Optional. File path. Use https://api.telegram.org/file/bot<token>/<file_path> to get the file.
	FilePath string `json:"file_path,omitempty"`
}

func (impl *File) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *File) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// ForceReply Upon receiving a message with this object, Telegram clients will display a reply interface to the user (act as if the user has selected the bot's message and tapped 'Reply').
// This can be extremely useful if you want to create user-friendly step-by-step interfaces without having to sacrifice privacy mode.
// Not supported in channels and for messages sent on behalf of a Telegram Business account.
type ForceReply struct {
	// Shows reply interface to the user, as if they manually selected the bot's message and tapped 'Reply'
	ForceReply bool `json:"force_reply"`
	// Optional. The placeholder to be shown in the input field when the reply is active; 1-64 characters
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	// Optional. Use this parameter if you want to force reply from specific users only.
	// Targets: 1) users that are @mentioned in the text of the Message object; 2) if the bot's message is a reply to a message in the same chat and forum topic, sender of the original message.
	Selective bool `json:"selective,omitempty"`
}

// ForumTopic This object represents a forum topic.
type ForumTopic struct {
	// Unique identifier of the forum topic
	MessageThreadId int64 `json:"message_thread_id"`
	// Name of the topic
	Name string `json:"name"`
	// Color of the topic icon in RGB format
	IconColor int64 `json:"icon_color"`
	// Optional. Unique identifier of the custom emoji shown as the topic icon
	IconCustomEmojiId string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicClosed This object represents a service message about a forum topic closed in the chat.
// Currently holds no information.
type ForumTopicClosed struct {
}

// ForumTopicCreated This object represents a service message about a new forum topic created in the chat.
type ForumTopicCreated struct {
	// Name of the topic
	Name string `json:"name"`
	// Color of the topic icon in RGB format
	IconColor int64 `json:"icon_color"`
	// Optional. Unique identifier of the custom emoji shown as the topic icon
	IconCustomEmojiId string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicEdited This object represents a service message about an edited forum topic.
type ForumTopicEdited struct {
	// Optional. New name of the topic, if it was edited
	Name string `json:"name,omitempty"`
	// Optional.
	// New identifier of the custom emoji shown as the topic icon, if it was edited; an empty string if the icon was removed
	IconCustomEmojiId string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicReopened This object represents a service message about a forum topic reopened in the chat.
// Currently holds no information.
type ForumTopicReopened struct {
}

// Game This object represents a game.
// Use BotFather to create and edit games, their short names will act as unique identifiers.
type Game struct {
	// Title of the game
	Title string `json:"title"`
	// Description of the game
	Description string `json:"description"`
	// Photo that will be displayed in the game message in chats.
	Photo TelegramPhoto `json:"photo"`
	// Optional. Brief description of the game or high scores included in the game message.
	// Can be automatically edited to include current high scores for the game when the bot calls setGameScore, or manually edited using editMessageText.
	// 0-4096 characters.
	Text string `json:"text,omitempty"`
	// Optional. Special entities that appear in text, such as usernames, URLs, bot commands, etc.
	TextEntities []*MessageEntity `json:"text_entities,omitempty"`
	// Optional. Animation that will be displayed in the game message in chats. Upload via BotFather
	Animation *TelegramAnimation `json:"animation,omitempty"`
}

// GameHighScore This object represents one row of the high scores table for a game.
type GameHighScore struct {
	// Position in high score table for the game
	Position int64 `json:"position"`
	// User
	User *User `json:"user"`
	// Score
	Score int64 `json:"score"`
}

// GeneralForumTopicHidden This object represents a service message about General forum topic hidden in the chat.
// Currently holds no information.
type GeneralForumTopicHidden struct {
}

// GeneralForumTopicUnhidden This object represents a service message about General forum topic unhidden in the chat.
// Currently holds no information.
type GeneralForumTopicUnhidden struct {
}

// Gift This object represents a gift that can be sent by the bot.
type Gift struct {
	// Unique identifier of the gift
	Id string `json:"id"`
	// The sticker that represents the gift
	Sticker *Sticker `json:"sticker"`
	// The number of Telegram Stars that must be paid to send the sticker
	StarCount int64 `json:"star_count"`
	// Optional. The total number of the gifts of this type that can be sent; for limited gifts only
	TotalCount int64 `json:"total_count,omitempty"`
	// Optional. The number of remaining gifts of this type that can be sent; for limited gifts only
	RemainingCount int64 `json:"remaining_count,omitempty"`
}

// Gifts This object represent a list of gifts.
type Gifts struct {
	// The list of gifts
	Gifts []*Gift `json:"gifts"`
}

// Giveaway This object represents a message about a scheduled giveaway.
type Giveaway struct {
	// The list of chats which the user must join to participate in the giveaway
	Chats []*Chat `json:"chats"`
	// Point in time (Unix timestamp) when winners of the giveaway will be selected
	WinnersSelectionDate int64 `json:"winners_selection_date"`
	// The number of users which are supposed to be selected as winners of the giveaway
	WinnerCount int64 `json:"winner_count"`
	// Optional. True, if only users who join the chats after the giveaway started should be eligible to win
	OnlyNewMembers bool `json:"only_new_members,omitempty"`
	// Optional. True, if the list of giveaway winners will be visible to everyone
	HasPublicWinners bool `json:"has_public_winners,omitempty"`
	// Optional. Description of additional giveaway prize
	PrizeDescription string `json:"prize_description,omitempty"`
	// Optional. If empty, then all users can participate in the giveaway.
	// A list of two-letter ISO 3166-1 alpha-2 country codes indicating the countries from which eligible users for the giveaway must come.
	// Users with a phone number that was bought on Fragment can always participate in giveaways.
	CountryCodes []string `json:"country_codes,omitempty"`
	// Optional. The number of Telegram Stars to be split between giveaway winners; for Telegram Star giveaways only
	PrizeStarCount int64 `json:"prize_star_count,omitempty"`
	// Optional.
	// The number of months the Telegram Premium subscription won from the giveaway will be active for; for Telegram Premium giveaways only
	PremiumSubscriptionMonthCount int64 `json:"premium_subscription_month_count,omitempty"`
}

// GiveawayCompleted This object represents a service message about the completion of a giveaway without public winners.
type GiveawayCompleted struct {
	// Number of winners in the giveaway
	WinnerCount int64 `json:"winner_count"`
	// Optional. Number of undistributed prizes
	UnclaimedPrizeCount int64 `json:"unclaimed_prize_count,omitempty"`
	// Optional. Message with the giveaway that was completed, if it wasn't deleted
	GiveawayMessage *Message `json:"giveaway_message,omitempty"`
	// Optional. True, if the giveaway is a Telegram Star giveaway.
	// Otherwise, currently, the giveaway is a Telegram Premium giveaway.
	IsStarGiveaway bool `json:"is_star_giveaway,omitempty"`
}

// GiveawayCreated This object represents a service message about the creation of a scheduled giveaway.
type GiveawayCreated struct {
	// Optional. The number of Telegram Stars to be split between giveaway winners; for Telegram Star giveaways only
	PrizeStarCount int64 `json:"prize_star_count,omitempty"`
}

// GiveawayWinners This object represents a message about the completion of a giveaway with public winners.
type GiveawayWinners struct {
	// The chat that created the giveaway
	Chat *Chat `json:"chat"`
	// Identifier of the message with the giveaway in the chat
	GiveawayMessageId int64 `json:"giveaway_message_id"`
	// Point in time (Unix timestamp) when winners of the giveaway were selected
	WinnersSelectionDate int64 `json:"winners_selection_date"`
	// Total number of winners in the giveaway
	WinnerCount int64 `json:"winner_count"`
	// List of up to 100 winners of the giveaway
	Winners []*User `json:"winners"`
	// Optional. The number of other chats the user had to join in order to be eligible for the giveaway
	AdditionalChatCount int64 `json:"additional_chat_count,omitempty"`
	// Optional. The number of Telegram Stars that were split between giveaway winners; for Telegram Star giveaways only
	PrizeStarCount int64 `json:"prize_star_count,omitempty"`
	// Optional.
	// The number of months the Telegram Premium subscription won from the giveaway will be active for; for Telegram Premium giveaways only
	PremiumSubscriptionMonthCount int64 `json:"premium_subscription_month_count,omitempty"`
	// Optional. Number of undistributed prizes
	UnclaimedPrizeCount int64 `json:"unclaimed_prize_count,omitempty"`
	// Optional. True, if only users who had joined the chats after the giveaway started were eligible to win
	OnlyNewMembers bool `json:"only_new_members,omitempty"`
	// Optional. True, if the giveaway was canceled because the payment for it was refunded
	WasRefunded bool `json:"was_refunded,omitempty"`
	// Optional. Description of additional giveaway prize
	PrizeDescription string `json:"prize_description,omitempty"`
}

// InlineKeyboardButton This object represents one button of an inline keyboard.
// Exactly one of the optional fields must be used to specify type of the button.
type InlineKeyboardButton struct {
	// Label text on the button
	Text string `json:"text"`
	// Optional. HTTP or tg:// URL to be opened when the button is pressed.
	// Links tg://user?id=<user_id> can be used to mention a user by their identifier without using a username, if this is allowed by their privacy settings.
	Url string `json:"url,omitempty"`
	// Optional. Data to be sent in a callback query to the bot when the button is pressed, 1-64 bytes
	CallbackData string `json:"callback_data,omitempty"`
	// Optional. Description of the Web App that will be launched when the user presses the button.
	// The Web App will be able to send an arbitrary message on behalf of the user using the method answerWebAppQuery.
	// Available only in private chats between a user and the bot.
	// Not supported for messages sent on behalf of a Telegram Business account.
	WebApp *WebAppInfo `json:"web_app,omitempty"`
	// Optional. An HTTPS URL used to automatically authorize the user.
	// Can be used as a replacement for the Telegram Login Widget.
	LoginUrl *LoginUrl `json:"login_url,omitempty"`
	// Optional. May be empty, in which case just the bot's username will be inserted.
	// If set, pressing the button will prompt the user to select one of their chats, open that chat and insert the bot's username and the specified inline query in the input field.
	// Not supported for messages sent on behalf of a Telegram Business account.
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`
	// Optional. May be empty, in which case only the bot's username will be inserted.
	// If set, pressing the button will insert the bot's username and the specified inline query in the current chat's input field.
	// This offers a quick way for the user to open your bot in inline mode in the same chat - good for selecting something from multiple options.
	// Not supported in channels and for messages sent on behalf of a Telegram Business account.
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
	// Optional. Not supported for messages sent on behalf of a Telegram Business account.
	// If set, pressing the button will prompt the user to select one of their chats of the specified type, open that chat and insert the bot's username and the specified inline query in the input field.
	SwitchInlineQueryChosenChat *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`
	// Optional. Description of the button that copies the specified text to the clipboard.
	CopyText *CopyTextButton `json:"copy_text,omitempty"`
	// Optional. Description of the game that will be launched when the user presses the button.
	// NOTE: This type of button must always be the first button in the first row.
	CallbackGame *CallbackGame `json:"callback_game,omitempty"`
	// Optional. Specify True, to send a Pay button.
	// Substrings "⭐" and "XTR" in the buttons's text will be replaced with a Telegram Star icon.
	// NOTE: This type of button must always be the first button in the first row and can only be used in invoice messages.
	Pay bool `json:"pay,omitempty"`
}

// InlineKeyboardMarkup This object represents an inline keyboard that appears right next to the message it belongs to.
type InlineKeyboardMarkup struct {
	// Array of button rows, each represented by an Array of InlineKeyboardButton objects
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineQuery This object represents an incoming inline query.
// When the user sends an empty query, your bot could return some default or trending results.
type InlineQuery struct {
	// Unique identifier for this query
	Id string `json:"id"`
	// Sender
	From *User `json:"from"`
	// Text of the query (up to 256 characters)
	Query string `json:"query"`
	// Offset of the results to be returned, can be controlled by the bot
	Offset string `json:"offset"`
	// Optional. Type of the chat from which the inline query was sent.
	// Can be either "sender" for a private chat with the inline query sender, "private", "group", "supergroup", or "channel".
	// The chat type should be always known for requests sent from official clients and most third-party clients, unless the request was sent from a secret chat
	ChatType string `json:"chat_type,omitempty"`
	// Optional. Sender location, only for bots that request user location
	Location *Location `json:"location,omitempty"`
}

// InlineQueryResult This object represents one result of an inline query. Telegram clients currently support results of the following 20 types:
// - InlineQueryResultCachedAudio
// - InlineQueryResultCachedDocument
// - InlineQueryResultCachedGif
// - InlineQueryResultCachedMpeg4Gif
// - InlineQueryResultCachedPhoto
// - InlineQueryResultCachedSticker
// - InlineQueryResultCachedVideo
// - InlineQueryResultCachedVoice
// - InlineQueryResultArticle
// - InlineQueryResultAudio
// - InlineQueryResultContact
// - InlineQueryResultGame
// - InlineQueryResultDocument
// - InlineQueryResultGif
// - InlineQueryResultLocation
// - InlineQueryResultMpeg4Gif
// - InlineQueryResultPhoto
// - InlineQueryResultVenue
// - InlineQueryResultVideo
// - InlineQueryResultVoice
// Note: All URLs passed in inline query results will be available to end users and therefore must be assumed to be public.
type InlineQueryResult interface {
	OptCachedAudio() *InlineQueryResultCachedAudio
	OptCachedDocument() *InlineQueryResultCachedDocument
	OptCachedGif() *InlineQueryResultCachedGif
	OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif
	OptCachedPhoto() *InlineQueryResultCachedPhoto
	OptCachedSticker() *InlineQueryResultCachedSticker
	OptCachedVideo() *InlineQueryResultCachedVideo
	OptCachedVoice() *InlineQueryResultCachedVoice
	OptArticle() *InlineQueryResultArticle
	OptAudio() *InlineQueryResultAudio
	OptContact() *InlineQueryResultContact
	OptGame() *InlineQueryResultGame
	OptDocument() *InlineQueryResultDocument
	OptGif() *InlineQueryResultGif
	OptLocation() *InlineQueryResultLocation
	OptMpeg4Gif() *InlineQueryResultMpeg4Gif
	OptPhoto() *InlineQueryResultPhoto
	OptVenue() *InlineQueryResultVenue
	OptVideo() *InlineQueryResultVideo
	OptVoice() *InlineQueryResultVoice
}

var (
	_ InlineQueryResult = &InlineQueryResultCachedAudio{}
	_ InlineQueryResult = &InlineQueryResultCachedDocument{}
	_ InlineQueryResult = &InlineQueryResultCachedGif{}
	_ InlineQueryResult = &InlineQueryResultCachedMpeg4Gif{}
	_ InlineQueryResult = &InlineQueryResultCachedPhoto{}
	_ InlineQueryResult = &InlineQueryResultCachedSticker{}
	_ InlineQueryResult = &InlineQueryResultCachedVideo{}
	_ InlineQueryResult = &InlineQueryResultCachedVoice{}
	_ InlineQueryResult = &InlineQueryResultArticle{}
	_ InlineQueryResult = &InlineQueryResultAudio{}
	_ InlineQueryResult = &InlineQueryResultContact{}
	_ InlineQueryResult = &InlineQueryResultGame{}
	_ InlineQueryResult = &InlineQueryResultDocument{}
	_ InlineQueryResult = &InlineQueryResultGif{}
	_ InlineQueryResult = &InlineQueryResultLocation{}
	_ InlineQueryResult = &InlineQueryResultMpeg4Gif{}
	_ InlineQueryResult = &InlineQueryResultPhoto{}
	_ InlineQueryResult = &InlineQueryResultVenue{}
	_ InlineQueryResult = &InlineQueryResultVideo{}
	_ InlineQueryResult = &InlineQueryResultVoice{}
)

func (impl *InlineQueryResultCachedAudio) OptCachedAudio() *InlineQueryResultCachedAudio { return impl }
func (impl *InlineQueryResultCachedAudio) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultCachedAudio) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultCachedAudio) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultCachedAudio) OptCachedPhoto() *InlineQueryResultCachedPhoto { return nil }
func (impl *InlineQueryResultCachedAudio) OptCachedSticker() *InlineQueryResultCachedSticker {
	return nil
}
func (impl *InlineQueryResultCachedAudio) OptCachedVideo() *InlineQueryResultCachedVideo { return nil }
func (impl *InlineQueryResultCachedAudio) OptCachedVoice() *InlineQueryResultCachedVoice { return nil }
func (impl *InlineQueryResultCachedAudio) OptArticle() *InlineQueryResultArticle         { return nil }
func (impl *InlineQueryResultCachedAudio) OptAudio() *InlineQueryResultAudio             { return nil }
func (impl *InlineQueryResultCachedAudio) OptContact() *InlineQueryResultContact         { return nil }
func (impl *InlineQueryResultCachedAudio) OptGame() *InlineQueryResultGame               { return nil }
func (impl *InlineQueryResultCachedAudio) OptDocument() *InlineQueryResultDocument       { return nil }
func (impl *InlineQueryResultCachedAudio) OptGif() *InlineQueryResultGif                 { return nil }
func (impl *InlineQueryResultCachedAudio) OptLocation() *InlineQueryResultLocation       { return nil }
func (impl *InlineQueryResultCachedAudio) OptMpeg4Gif() *InlineQueryResultMpeg4Gif       { return nil }
func (impl *InlineQueryResultCachedAudio) OptPhoto() *InlineQueryResultPhoto             { return nil }
func (impl *InlineQueryResultCachedAudio) OptVenue() *InlineQueryResultVenue             { return nil }
func (impl *InlineQueryResultCachedAudio) OptVideo() *InlineQueryResultVideo             { return nil }
func (impl *InlineQueryResultCachedAudio) OptVoice() *InlineQueryResultVoice             { return nil }

func (impl *InlineQueryResultCachedDocument) OptCachedAudio() *InlineQueryResultCachedAudio {
	return nil
}
func (impl *InlineQueryResultCachedDocument) OptCachedDocument() *InlineQueryResultCachedDocument {
	return impl
}
func (impl *InlineQueryResultCachedDocument) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultCachedDocument) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultCachedDocument) OptCachedPhoto() *InlineQueryResultCachedPhoto {
	return nil
}
func (impl *InlineQueryResultCachedDocument) OptCachedSticker() *InlineQueryResultCachedSticker {
	return nil
}
func (impl *InlineQueryResultCachedDocument) OptCachedVideo() *InlineQueryResultCachedVideo {
	return nil
}
func (impl *InlineQueryResultCachedDocument) OptCachedVoice() *InlineQueryResultCachedVoice {
	return nil
}
func (impl *InlineQueryResultCachedDocument) OptArticle() *InlineQueryResultArticle   { return nil }
func (impl *InlineQueryResultCachedDocument) OptAudio() *InlineQueryResultAudio       { return nil }
func (impl *InlineQueryResultCachedDocument) OptContact() *InlineQueryResultContact   { return nil }
func (impl *InlineQueryResultCachedDocument) OptGame() *InlineQueryResultGame         { return nil }
func (impl *InlineQueryResultCachedDocument) OptDocument() *InlineQueryResultDocument { return nil }
func (impl *InlineQueryResultCachedDocument) OptGif() *InlineQueryResultGif           { return nil }
func (impl *InlineQueryResultCachedDocument) OptLocation() *InlineQueryResultLocation { return nil }
func (impl *InlineQueryResultCachedDocument) OptMpeg4Gif() *InlineQueryResultMpeg4Gif { return nil }
func (impl *InlineQueryResultCachedDocument) OptPhoto() *InlineQueryResultPhoto       { return nil }
func (impl *InlineQueryResultCachedDocument) OptVenue() *InlineQueryResultVenue       { return nil }
func (impl *InlineQueryResultCachedDocument) OptVideo() *InlineQueryResultVideo       { return nil }
func (impl *InlineQueryResultCachedDocument) OptVoice() *InlineQueryResultVoice       { return nil }

func (impl *InlineQueryResultCachedGif) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultCachedGif) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultCachedGif) OptCachedGif() *InlineQueryResultCachedGif { return impl }
func (impl *InlineQueryResultCachedGif) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultCachedGif) OptCachedPhoto() *InlineQueryResultCachedPhoto { return nil }
func (impl *InlineQueryResultCachedGif) OptCachedSticker() *InlineQueryResultCachedSticker {
	return nil
}
func (impl *InlineQueryResultCachedGif) OptCachedVideo() *InlineQueryResultCachedVideo { return nil }
func (impl *InlineQueryResultCachedGif) OptCachedVoice() *InlineQueryResultCachedVoice { return nil }
func (impl *InlineQueryResultCachedGif) OptArticle() *InlineQueryResultArticle         { return nil }
func (impl *InlineQueryResultCachedGif) OptAudio() *InlineQueryResultAudio             { return nil }
func (impl *InlineQueryResultCachedGif) OptContact() *InlineQueryResultContact         { return nil }
func (impl *InlineQueryResultCachedGif) OptGame() *InlineQueryResultGame               { return nil }
func (impl *InlineQueryResultCachedGif) OptDocument() *InlineQueryResultDocument       { return nil }
func (impl *InlineQueryResultCachedGif) OptGif() *InlineQueryResultGif                 { return nil }
func (impl *InlineQueryResultCachedGif) OptLocation() *InlineQueryResultLocation       { return nil }
func (impl *InlineQueryResultCachedGif) OptMpeg4Gif() *InlineQueryResultMpeg4Gif       { return nil }
func (impl *InlineQueryResultCachedGif) OptPhoto() *InlineQueryResultPhoto             { return nil }
func (impl *InlineQueryResultCachedGif) OptVenue() *InlineQueryResultVenue             { return nil }
func (impl *InlineQueryResultCachedGif) OptVideo() *InlineQueryResultVideo             { return nil }
func (impl *InlineQueryResultCachedGif) OptVoice() *InlineQueryResultVoice             { return nil }

func (impl *InlineQueryResultCachedMpeg4Gif) OptCachedAudio() *InlineQueryResultCachedAudio {
	return nil
}
func (impl *InlineQueryResultCachedMpeg4Gif) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultCachedMpeg4Gif) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return impl
}
func (impl *InlineQueryResultCachedMpeg4Gif) OptCachedPhoto() *InlineQueryResultCachedPhoto {
	return nil
}
func (impl *InlineQueryResultCachedMpeg4Gif) OptCachedSticker() *InlineQueryResultCachedSticker {
	return nil
}
func (impl *InlineQueryResultCachedMpeg4Gif) OptCachedVideo() *InlineQueryResultCachedVideo {
	return nil
}
func (impl *InlineQueryResultCachedMpeg4Gif) OptCachedVoice() *InlineQueryResultCachedVoice {
	return nil
}
func (impl *InlineQueryResultCachedMpeg4Gif) OptArticle() *InlineQueryResultArticle   { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptAudio() *InlineQueryResultAudio       { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptContact() *InlineQueryResultContact   { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptGame() *InlineQueryResultGame         { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptDocument() *InlineQueryResultDocument { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptGif() *InlineQueryResultGif           { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptLocation() *InlineQueryResultLocation { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptMpeg4Gif() *InlineQueryResultMpeg4Gif { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptPhoto() *InlineQueryResultPhoto       { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptVenue() *InlineQueryResultVenue       { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptVideo() *InlineQueryResultVideo       { return nil }
func (impl *InlineQueryResultCachedMpeg4Gif) OptVoice() *InlineQueryResultVoice       { return nil }

func (impl *InlineQueryResultCachedPhoto) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultCachedPhoto) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultCachedPhoto) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultCachedPhoto) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultCachedPhoto) OptCachedPhoto() *InlineQueryResultCachedPhoto { return impl }
func (impl *InlineQueryResultCachedPhoto) OptCachedSticker() *InlineQueryResultCachedSticker {
	return nil
}
func (impl *InlineQueryResultCachedPhoto) OptCachedVideo() *InlineQueryResultCachedVideo { return nil }
func (impl *InlineQueryResultCachedPhoto) OptCachedVoice() *InlineQueryResultCachedVoice { return nil }
func (impl *InlineQueryResultCachedPhoto) OptArticle() *InlineQueryResultArticle         { return nil }
func (impl *InlineQueryResultCachedPhoto) OptAudio() *InlineQueryResultAudio             { return nil }
func (impl *InlineQueryResultCachedPhoto) OptContact() *InlineQueryResultContact         { return nil }
func (impl *InlineQueryResultCachedPhoto) OptGame() *InlineQueryResultGame               { return nil }
func (impl *InlineQueryResultCachedPhoto) OptDocument() *InlineQueryResultDocument       { return nil }
func (impl *InlineQueryResultCachedPhoto) OptGif() *InlineQueryResultGif                 { return nil }
func (impl *InlineQueryResultCachedPhoto) OptLocation() *InlineQueryResultLocation       { return nil }
func (impl *InlineQueryResultCachedPhoto) OptMpeg4Gif() *InlineQueryResultMpeg4Gif       { return nil }
func (impl *InlineQueryResultCachedPhoto) OptPhoto() *InlineQueryResultPhoto             { return nil }
func (impl *InlineQueryResultCachedPhoto) OptVenue() *InlineQueryResultVenue             { return nil }
func (impl *InlineQueryResultCachedPhoto) OptVideo() *InlineQueryResultVideo             { return nil }
func (impl *InlineQueryResultCachedPhoto) OptVoice() *InlineQueryResultVoice             { return nil }

func (impl *InlineQueryResultCachedSticker) OptCachedAudio() *InlineQueryResultCachedAudio {
	return nil
}
func (impl *InlineQueryResultCachedSticker) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultCachedSticker) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultCachedSticker) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultCachedSticker) OptCachedPhoto() *InlineQueryResultCachedPhoto {
	return nil
}
func (impl *InlineQueryResultCachedSticker) OptCachedSticker() *InlineQueryResultCachedSticker {
	return impl
}
func (impl *InlineQueryResultCachedSticker) OptCachedVideo() *InlineQueryResultCachedVideo {
	return nil
}
func (impl *InlineQueryResultCachedSticker) OptCachedVoice() *InlineQueryResultCachedVoice {
	return nil
}
func (impl *InlineQueryResultCachedSticker) OptArticle() *InlineQueryResultArticle   { return nil }
func (impl *InlineQueryResultCachedSticker) OptAudio() *InlineQueryResultAudio       { return nil }
func (impl *InlineQueryResultCachedSticker) OptContact() *InlineQueryResultContact   { return nil }
func (impl *InlineQueryResultCachedSticker) OptGame() *InlineQueryResultGame         { return nil }
func (impl *InlineQueryResultCachedSticker) OptDocument() *InlineQueryResultDocument { return nil }
func (impl *InlineQueryResultCachedSticker) OptGif() *InlineQueryResultGif           { return nil }
func (impl *InlineQueryResultCachedSticker) OptLocation() *InlineQueryResultLocation { return nil }
func (impl *InlineQueryResultCachedSticker) OptMpeg4Gif() *InlineQueryResultMpeg4Gif { return nil }
func (impl *InlineQueryResultCachedSticker) OptPhoto() *InlineQueryResultPhoto       { return nil }
func (impl *InlineQueryResultCachedSticker) OptVenue() *InlineQueryResultVenue       { return nil }
func (impl *InlineQueryResultCachedSticker) OptVideo() *InlineQueryResultVideo       { return nil }
func (impl *InlineQueryResultCachedSticker) OptVoice() *InlineQueryResultVoice       { return nil }

func (impl *InlineQueryResultCachedVideo) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultCachedVideo) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultCachedVideo) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultCachedVideo) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultCachedVideo) OptCachedPhoto() *InlineQueryResultCachedPhoto { return nil }
func (impl *InlineQueryResultCachedVideo) OptCachedSticker() *InlineQueryResultCachedSticker {
	return nil
}
func (impl *InlineQueryResultCachedVideo) OptCachedVideo() *InlineQueryResultCachedVideo { return impl }
func (impl *InlineQueryResultCachedVideo) OptCachedVoice() *InlineQueryResultCachedVoice { return nil }
func (impl *InlineQueryResultCachedVideo) OptArticle() *InlineQueryResultArticle         { return nil }
func (impl *InlineQueryResultCachedVideo) OptAudio() *InlineQueryResultAudio             { return nil }
func (impl *InlineQueryResultCachedVideo) OptContact() *InlineQueryResultContact         { return nil }
func (impl *InlineQueryResultCachedVideo) OptGame() *InlineQueryResultGame               { return nil }
func (impl *InlineQueryResultCachedVideo) OptDocument() *InlineQueryResultDocument       { return nil }
func (impl *InlineQueryResultCachedVideo) OptGif() *InlineQueryResultGif                 { return nil }
func (impl *InlineQueryResultCachedVideo) OptLocation() *InlineQueryResultLocation       { return nil }
func (impl *InlineQueryResultCachedVideo) OptMpeg4Gif() *InlineQueryResultMpeg4Gif       { return nil }
func (impl *InlineQueryResultCachedVideo) OptPhoto() *InlineQueryResultPhoto             { return nil }
func (impl *InlineQueryResultCachedVideo) OptVenue() *InlineQueryResultVenue             { return nil }
func (impl *InlineQueryResultCachedVideo) OptVideo() *InlineQueryResultVideo             { return nil }
func (impl *InlineQueryResultCachedVideo) OptVoice() *InlineQueryResultVoice             { return nil }

func (impl *InlineQueryResultCachedVoice) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultCachedVoice) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultCachedVoice) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultCachedVoice) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultCachedVoice) OptCachedPhoto() *InlineQueryResultCachedPhoto { return nil }
func (impl *InlineQueryResultCachedVoice) OptCachedSticker() *InlineQueryResultCachedSticker {
	return nil
}
func (impl *InlineQueryResultCachedVoice) OptCachedVideo() *InlineQueryResultCachedVideo { return nil }
func (impl *InlineQueryResultCachedVoice) OptCachedVoice() *InlineQueryResultCachedVoice { return impl }
func (impl *InlineQueryResultCachedVoice) OptArticle() *InlineQueryResultArticle         { return nil }
func (impl *InlineQueryResultCachedVoice) OptAudio() *InlineQueryResultAudio             { return nil }
func (impl *InlineQueryResultCachedVoice) OptContact() *InlineQueryResultContact         { return nil }
func (impl *InlineQueryResultCachedVoice) OptGame() *InlineQueryResultGame               { return nil }
func (impl *InlineQueryResultCachedVoice) OptDocument() *InlineQueryResultDocument       { return nil }
func (impl *InlineQueryResultCachedVoice) OptGif() *InlineQueryResultGif                 { return nil }
func (impl *InlineQueryResultCachedVoice) OptLocation() *InlineQueryResultLocation       { return nil }
func (impl *InlineQueryResultCachedVoice) OptMpeg4Gif() *InlineQueryResultMpeg4Gif       { return nil }
func (impl *InlineQueryResultCachedVoice) OptPhoto() *InlineQueryResultPhoto             { return nil }
func (impl *InlineQueryResultCachedVoice) OptVenue() *InlineQueryResultVenue             { return nil }
func (impl *InlineQueryResultCachedVoice) OptVideo() *InlineQueryResultVideo             { return nil }
func (impl *InlineQueryResultCachedVoice) OptVoice() *InlineQueryResultVoice             { return nil }

func (impl *InlineQueryResultArticle) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultArticle) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultArticle) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultArticle) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultArticle) OptCachedPhoto() *InlineQueryResultCachedPhoto     { return nil }
func (impl *InlineQueryResultArticle) OptCachedSticker() *InlineQueryResultCachedSticker { return nil }
func (impl *InlineQueryResultArticle) OptCachedVideo() *InlineQueryResultCachedVideo     { return nil }
func (impl *InlineQueryResultArticle) OptCachedVoice() *InlineQueryResultCachedVoice     { return nil }
func (impl *InlineQueryResultArticle) OptArticle() *InlineQueryResultArticle             { return impl }
func (impl *InlineQueryResultArticle) OptAudio() *InlineQueryResultAudio                 { return nil }
func (impl *InlineQueryResultArticle) OptContact() *InlineQueryResultContact             { return nil }
func (impl *InlineQueryResultArticle) OptGame() *InlineQueryResultGame                   { return nil }
func (impl *InlineQueryResultArticle) OptDocument() *InlineQueryResultDocument           { return nil }
func (impl *InlineQueryResultArticle) OptGif() *InlineQueryResultGif                     { return nil }
func (impl *InlineQueryResultArticle) OptLocation() *InlineQueryResultLocation           { return nil }
func (impl *InlineQueryResultArticle) OptMpeg4Gif() *InlineQueryResultMpeg4Gif           { return nil }
func (impl *InlineQueryResultArticle) OptPhoto() *InlineQueryResultPhoto                 { return nil }
func (impl *InlineQueryResultArticle) OptVenue() *InlineQueryResultVenue                 { return nil }
func (impl *InlineQueryResultArticle) OptVideo() *InlineQueryResultVideo                 { return nil }
func (impl *InlineQueryResultArticle) OptVoice() *InlineQueryResultVoice                 { return nil }

func (impl *InlineQueryResultAudio) OptCachedAudio() *InlineQueryResultCachedAudio       { return nil }
func (impl *InlineQueryResultAudio) OptCachedDocument() *InlineQueryResultCachedDocument { return nil }
func (impl *InlineQueryResultAudio) OptCachedGif() *InlineQueryResultCachedGif           { return nil }
func (impl *InlineQueryResultAudio) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif { return nil }
func (impl *InlineQueryResultAudio) OptCachedPhoto() *InlineQueryResultCachedPhoto       { return nil }
func (impl *InlineQueryResultAudio) OptCachedSticker() *InlineQueryResultCachedSticker   { return nil }
func (impl *InlineQueryResultAudio) OptCachedVideo() *InlineQueryResultCachedVideo       { return nil }
func (impl *InlineQueryResultAudio) OptCachedVoice() *InlineQueryResultCachedVoice       { return nil }
func (impl *InlineQueryResultAudio) OptArticle() *InlineQueryResultArticle               { return nil }
func (impl *InlineQueryResultAudio) OptAudio() *InlineQueryResultAudio                   { return impl }
func (impl *InlineQueryResultAudio) OptContact() *InlineQueryResultContact               { return nil }
func (impl *InlineQueryResultAudio) OptGame() *InlineQueryResultGame                     { return nil }
func (impl *InlineQueryResultAudio) OptDocument() *InlineQueryResultDocument             { return nil }
func (impl *InlineQueryResultAudio) OptGif() *InlineQueryResultGif                       { return nil }
func (impl *InlineQueryResultAudio) OptLocation() *InlineQueryResultLocation             { return nil }
func (impl *InlineQueryResultAudio) OptMpeg4Gif() *InlineQueryResultMpeg4Gif             { return nil }
func (impl *InlineQueryResultAudio) OptPhoto() *InlineQueryResultPhoto                   { return nil }
func (impl *InlineQueryResultAudio) OptVenue() *InlineQueryResultVenue                   { return nil }
func (impl *InlineQueryResultAudio) OptVideo() *InlineQueryResultVideo                   { return nil }
func (impl *InlineQueryResultAudio) OptVoice() *InlineQueryResultVoice                   { return nil }

func (impl *InlineQueryResultContact) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultContact) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultContact) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultContact) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultContact) OptCachedPhoto() *InlineQueryResultCachedPhoto     { return nil }
func (impl *InlineQueryResultContact) OptCachedSticker() *InlineQueryResultCachedSticker { return nil }
func (impl *InlineQueryResultContact) OptCachedVideo() *InlineQueryResultCachedVideo     { return nil }
func (impl *InlineQueryResultContact) OptCachedVoice() *InlineQueryResultCachedVoice     { return nil }
func (impl *InlineQueryResultContact) OptArticle() *InlineQueryResultArticle             { return nil }
func (impl *InlineQueryResultContact) OptAudio() *InlineQueryResultAudio                 { return nil }
func (impl *InlineQueryResultContact) OptContact() *InlineQueryResultContact             { return impl }
func (impl *InlineQueryResultContact) OptGame() *InlineQueryResultGame                   { return nil }
func (impl *InlineQueryResultContact) OptDocument() *InlineQueryResultDocument           { return nil }
func (impl *InlineQueryResultContact) OptGif() *InlineQueryResultGif                     { return nil }
func (impl *InlineQueryResultContact) OptLocation() *InlineQueryResultLocation           { return nil }
func (impl *InlineQueryResultContact) OptMpeg4Gif() *InlineQueryResultMpeg4Gif           { return nil }
func (impl *InlineQueryResultContact) OptPhoto() *InlineQueryResultPhoto                 { return nil }
func (impl *InlineQueryResultContact) OptVenue() *InlineQueryResultVenue                 { return nil }
func (impl *InlineQueryResultContact) OptVideo() *InlineQueryResultVideo                 { return nil }
func (impl *InlineQueryResultContact) OptVoice() *InlineQueryResultVoice                 { return nil }

func (impl *InlineQueryResultGame) OptCachedAudio() *InlineQueryResultCachedAudio       { return nil }
func (impl *InlineQueryResultGame) OptCachedDocument() *InlineQueryResultCachedDocument { return nil }
func (impl *InlineQueryResultGame) OptCachedGif() *InlineQueryResultCachedGif           { return nil }
func (impl *InlineQueryResultGame) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif { return nil }
func (impl *InlineQueryResultGame) OptCachedPhoto() *InlineQueryResultCachedPhoto       { return nil }
func (impl *InlineQueryResultGame) OptCachedSticker() *InlineQueryResultCachedSticker   { return nil }
func (impl *InlineQueryResultGame) OptCachedVideo() *InlineQueryResultCachedVideo       { return nil }
func (impl *InlineQueryResultGame) OptCachedVoice() *InlineQueryResultCachedVoice       { return nil }
func (impl *InlineQueryResultGame) OptArticle() *InlineQueryResultArticle               { return nil }
func (impl *InlineQueryResultGame) OptAudio() *InlineQueryResultAudio                   { return nil }
func (impl *InlineQueryResultGame) OptContact() *InlineQueryResultContact               { return nil }
func (impl *InlineQueryResultGame) OptGame() *InlineQueryResultGame                     { return impl }
func (impl *InlineQueryResultGame) OptDocument() *InlineQueryResultDocument             { return nil }
func (impl *InlineQueryResultGame) OptGif() *InlineQueryResultGif                       { return nil }
func (impl *InlineQueryResultGame) OptLocation() *InlineQueryResultLocation             { return nil }
func (impl *InlineQueryResultGame) OptMpeg4Gif() *InlineQueryResultMpeg4Gif             { return nil }
func (impl *InlineQueryResultGame) OptPhoto() *InlineQueryResultPhoto                   { return nil }
func (impl *InlineQueryResultGame) OptVenue() *InlineQueryResultVenue                   { return nil }
func (impl *InlineQueryResultGame) OptVideo() *InlineQueryResultVideo                   { return nil }
func (impl *InlineQueryResultGame) OptVoice() *InlineQueryResultVoice                   { return nil }

func (impl *InlineQueryResultDocument) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultDocument) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultDocument) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultDocument) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultDocument) OptCachedPhoto() *InlineQueryResultCachedPhoto     { return nil }
func (impl *InlineQueryResultDocument) OptCachedSticker() *InlineQueryResultCachedSticker { return nil }
func (impl *InlineQueryResultDocument) OptCachedVideo() *InlineQueryResultCachedVideo     { return nil }
func (impl *InlineQueryResultDocument) OptCachedVoice() *InlineQueryResultCachedVoice     { return nil }
func (impl *InlineQueryResultDocument) OptArticle() *InlineQueryResultArticle             { return nil }
func (impl *InlineQueryResultDocument) OptAudio() *InlineQueryResultAudio                 { return nil }
func (impl *InlineQueryResultDocument) OptContact() *InlineQueryResultContact             { return nil }
func (impl *InlineQueryResultDocument) OptGame() *InlineQueryResultGame                   { return nil }
func (impl *InlineQueryResultDocument) OptDocument() *InlineQueryResultDocument           { return impl }
func (impl *InlineQueryResultDocument) OptGif() *InlineQueryResultGif                     { return nil }
func (impl *InlineQueryResultDocument) OptLocation() *InlineQueryResultLocation           { return nil }
func (impl *InlineQueryResultDocument) OptMpeg4Gif() *InlineQueryResultMpeg4Gif           { return nil }
func (impl *InlineQueryResultDocument) OptPhoto() *InlineQueryResultPhoto                 { return nil }
func (impl *InlineQueryResultDocument) OptVenue() *InlineQueryResultVenue                 { return nil }
func (impl *InlineQueryResultDocument) OptVideo() *InlineQueryResultVideo                 { return nil }
func (impl *InlineQueryResultDocument) OptVoice() *InlineQueryResultVoice                 { return nil }

func (impl *InlineQueryResultGif) OptCachedAudio() *InlineQueryResultCachedAudio       { return nil }
func (impl *InlineQueryResultGif) OptCachedDocument() *InlineQueryResultCachedDocument { return nil }
func (impl *InlineQueryResultGif) OptCachedGif() *InlineQueryResultCachedGif           { return nil }
func (impl *InlineQueryResultGif) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif { return nil }
func (impl *InlineQueryResultGif) OptCachedPhoto() *InlineQueryResultCachedPhoto       { return nil }
func (impl *InlineQueryResultGif) OptCachedSticker() *InlineQueryResultCachedSticker   { return nil }
func (impl *InlineQueryResultGif) OptCachedVideo() *InlineQueryResultCachedVideo       { return nil }
func (impl *InlineQueryResultGif) OptCachedVoice() *InlineQueryResultCachedVoice       { return nil }
func (impl *InlineQueryResultGif) OptArticle() *InlineQueryResultArticle               { return nil }
func (impl *InlineQueryResultGif) OptAudio() *InlineQueryResultAudio                   { return nil }
func (impl *InlineQueryResultGif) OptContact() *InlineQueryResultContact               { return nil }
func (impl *InlineQueryResultGif) OptGame() *InlineQueryResultGame                     { return nil }
func (impl *InlineQueryResultGif) OptDocument() *InlineQueryResultDocument             { return nil }
func (impl *InlineQueryResultGif) OptGif() *InlineQueryResultGif                       { return impl }
func (impl *InlineQueryResultGif) OptLocation() *InlineQueryResultLocation             { return nil }
func (impl *InlineQueryResultGif) OptMpeg4Gif() *InlineQueryResultMpeg4Gif             { return nil }
func (impl *InlineQueryResultGif) OptPhoto() *InlineQueryResultPhoto                   { return nil }
func (impl *InlineQueryResultGif) OptVenue() *InlineQueryResultVenue                   { return nil }
func (impl *InlineQueryResultGif) OptVideo() *InlineQueryResultVideo                   { return nil }
func (impl *InlineQueryResultGif) OptVoice() *InlineQueryResultVoice                   { return nil }

func (impl *InlineQueryResultLocation) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultLocation) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultLocation) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultLocation) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultLocation) OptCachedPhoto() *InlineQueryResultCachedPhoto     { return nil }
func (impl *InlineQueryResultLocation) OptCachedSticker() *InlineQueryResultCachedSticker { return nil }
func (impl *InlineQueryResultLocation) OptCachedVideo() *InlineQueryResultCachedVideo     { return nil }
func (impl *InlineQueryResultLocation) OptCachedVoice() *InlineQueryResultCachedVoice     { return nil }
func (impl *InlineQueryResultLocation) OptArticle() *InlineQueryResultArticle             { return nil }
func (impl *InlineQueryResultLocation) OptAudio() *InlineQueryResultAudio                 { return nil }
func (impl *InlineQueryResultLocation) OptContact() *InlineQueryResultContact             { return nil }
func (impl *InlineQueryResultLocation) OptGame() *InlineQueryResultGame                   { return nil }
func (impl *InlineQueryResultLocation) OptDocument() *InlineQueryResultDocument           { return nil }
func (impl *InlineQueryResultLocation) OptGif() *InlineQueryResultGif                     { return nil }
func (impl *InlineQueryResultLocation) OptLocation() *InlineQueryResultLocation           { return impl }
func (impl *InlineQueryResultLocation) OptMpeg4Gif() *InlineQueryResultMpeg4Gif           { return nil }
func (impl *InlineQueryResultLocation) OptPhoto() *InlineQueryResultPhoto                 { return nil }
func (impl *InlineQueryResultLocation) OptVenue() *InlineQueryResultVenue                 { return nil }
func (impl *InlineQueryResultLocation) OptVideo() *InlineQueryResultVideo                 { return nil }
func (impl *InlineQueryResultLocation) OptVoice() *InlineQueryResultVoice                 { return nil }

func (impl *InlineQueryResultMpeg4Gif) OptCachedAudio() *InlineQueryResultCachedAudio { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptCachedDocument() *InlineQueryResultCachedDocument {
	return nil
}
func (impl *InlineQueryResultMpeg4Gif) OptCachedGif() *InlineQueryResultCachedGif { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif {
	return nil
}
func (impl *InlineQueryResultMpeg4Gif) OptCachedPhoto() *InlineQueryResultCachedPhoto     { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptCachedSticker() *InlineQueryResultCachedSticker { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptCachedVideo() *InlineQueryResultCachedVideo     { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptCachedVoice() *InlineQueryResultCachedVoice     { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptArticle() *InlineQueryResultArticle             { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptAudio() *InlineQueryResultAudio                 { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptContact() *InlineQueryResultContact             { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptGame() *InlineQueryResultGame                   { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptDocument() *InlineQueryResultDocument           { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptGif() *InlineQueryResultGif                     { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptLocation() *InlineQueryResultLocation           { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptMpeg4Gif() *InlineQueryResultMpeg4Gif           { return impl }
func (impl *InlineQueryResultMpeg4Gif) OptPhoto() *InlineQueryResultPhoto                 { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptVenue() *InlineQueryResultVenue                 { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptVideo() *InlineQueryResultVideo                 { return nil }
func (impl *InlineQueryResultMpeg4Gif) OptVoice() *InlineQueryResultVoice                 { return nil }

func (impl *InlineQueryResultPhoto) OptCachedAudio() *InlineQueryResultCachedAudio       { return nil }
func (impl *InlineQueryResultPhoto) OptCachedDocument() *InlineQueryResultCachedDocument { return nil }
func (impl *InlineQueryResultPhoto) OptCachedGif() *InlineQueryResultCachedGif           { return nil }
func (impl *InlineQueryResultPhoto) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif { return nil }
func (impl *InlineQueryResultPhoto) OptCachedPhoto() *InlineQueryResultCachedPhoto       { return nil }
func (impl *InlineQueryResultPhoto) OptCachedSticker() *InlineQueryResultCachedSticker   { return nil }
func (impl *InlineQueryResultPhoto) OptCachedVideo() *InlineQueryResultCachedVideo       { return nil }
func (impl *InlineQueryResultPhoto) OptCachedVoice() *InlineQueryResultCachedVoice       { return nil }
func (impl *InlineQueryResultPhoto) OptArticle() *InlineQueryResultArticle               { return nil }
func (impl *InlineQueryResultPhoto) OptAudio() *InlineQueryResultAudio                   { return nil }
func (impl *InlineQueryResultPhoto) OptContact() *InlineQueryResultContact               { return nil }
func (impl *InlineQueryResultPhoto) OptGame() *InlineQueryResultGame                     { return nil }
func (impl *InlineQueryResultPhoto) OptDocument() *InlineQueryResultDocument             { return nil }
func (impl *InlineQueryResultPhoto) OptGif() *InlineQueryResultGif                       { return nil }
func (impl *InlineQueryResultPhoto) OptLocation() *InlineQueryResultLocation             { return nil }
func (impl *InlineQueryResultPhoto) OptMpeg4Gif() *InlineQueryResultMpeg4Gif             { return nil }
func (impl *InlineQueryResultPhoto) OptPhoto() *InlineQueryResultPhoto                   { return impl }
func (impl *InlineQueryResultPhoto) OptVenue() *InlineQueryResultVenue                   { return nil }
func (impl *InlineQueryResultPhoto) OptVideo() *InlineQueryResultVideo                   { return nil }
func (impl *InlineQueryResultPhoto) OptVoice() *InlineQueryResultVoice                   { return nil }

func (impl *InlineQueryResultVenue) OptCachedAudio() *InlineQueryResultCachedAudio       { return nil }
func (impl *InlineQueryResultVenue) OptCachedDocument() *InlineQueryResultCachedDocument { return nil }
func (impl *InlineQueryResultVenue) OptCachedGif() *InlineQueryResultCachedGif           { return nil }
func (impl *InlineQueryResultVenue) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif { return nil }
func (impl *InlineQueryResultVenue) OptCachedPhoto() *InlineQueryResultCachedPhoto       { return nil }
func (impl *InlineQueryResultVenue) OptCachedSticker() *InlineQueryResultCachedSticker   { return nil }
func (impl *InlineQueryResultVenue) OptCachedVideo() *InlineQueryResultCachedVideo       { return nil }
func (impl *InlineQueryResultVenue) OptCachedVoice() *InlineQueryResultCachedVoice       { return nil }
func (impl *InlineQueryResultVenue) OptArticle() *InlineQueryResultArticle               { return nil }
func (impl *InlineQueryResultVenue) OptAudio() *InlineQueryResultAudio                   { return nil }
func (impl *InlineQueryResultVenue) OptContact() *InlineQueryResultContact               { return nil }
func (impl *InlineQueryResultVenue) OptGame() *InlineQueryResultGame                     { return nil }
func (impl *InlineQueryResultVenue) OptDocument() *InlineQueryResultDocument             { return nil }
func (impl *InlineQueryResultVenue) OptGif() *InlineQueryResultGif                       { return nil }
func (impl *InlineQueryResultVenue) OptLocation() *InlineQueryResultLocation             { return nil }
func (impl *InlineQueryResultVenue) OptMpeg4Gif() *InlineQueryResultMpeg4Gif             { return nil }
func (impl *InlineQueryResultVenue) OptPhoto() *InlineQueryResultPhoto                   { return nil }
func (impl *InlineQueryResultVenue) OptVenue() *InlineQueryResultVenue                   { return impl }
func (impl *InlineQueryResultVenue) OptVideo() *InlineQueryResultVideo                   { return nil }
func (impl *InlineQueryResultVenue) OptVoice() *InlineQueryResultVoice                   { return nil }

func (impl *InlineQueryResultVideo) OptCachedAudio() *InlineQueryResultCachedAudio       { return nil }
func (impl *InlineQueryResultVideo) OptCachedDocument() *InlineQueryResultCachedDocument { return nil }
func (impl *InlineQueryResultVideo) OptCachedGif() *InlineQueryResultCachedGif           { return nil }
func (impl *InlineQueryResultVideo) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif { return nil }
func (impl *InlineQueryResultVideo) OptCachedPhoto() *InlineQueryResultCachedPhoto       { return nil }
func (impl *InlineQueryResultVideo) OptCachedSticker() *InlineQueryResultCachedSticker   { return nil }
func (impl *InlineQueryResultVideo) OptCachedVideo() *InlineQueryResultCachedVideo       { return nil }
func (impl *InlineQueryResultVideo) OptCachedVoice() *InlineQueryResultCachedVoice       { return nil }
func (impl *InlineQueryResultVideo) OptArticle() *InlineQueryResultArticle               { return nil }
func (impl *InlineQueryResultVideo) OptAudio() *InlineQueryResultAudio                   { return nil }
func (impl *InlineQueryResultVideo) OptContact() *InlineQueryResultContact               { return nil }
func (impl *InlineQueryResultVideo) OptGame() *InlineQueryResultGame                     { return nil }
func (impl *InlineQueryResultVideo) OptDocument() *InlineQueryResultDocument             { return nil }
func (impl *InlineQueryResultVideo) OptGif() *InlineQueryResultGif                       { return nil }
func (impl *InlineQueryResultVideo) OptLocation() *InlineQueryResultLocation             { return nil }
func (impl *InlineQueryResultVideo) OptMpeg4Gif() *InlineQueryResultMpeg4Gif             { return nil }
func (impl *InlineQueryResultVideo) OptPhoto() *InlineQueryResultPhoto                   { return nil }
func (impl *InlineQueryResultVideo) OptVenue() *InlineQueryResultVenue                   { return nil }
func (impl *InlineQueryResultVideo) OptVideo() *InlineQueryResultVideo                   { return impl }
func (impl *InlineQueryResultVideo) OptVoice() *InlineQueryResultVoice                   { return nil }

func (impl *InlineQueryResultVoice) OptCachedAudio() *InlineQueryResultCachedAudio       { return nil }
func (impl *InlineQueryResultVoice) OptCachedDocument() *InlineQueryResultCachedDocument { return nil }
func (impl *InlineQueryResultVoice) OptCachedGif() *InlineQueryResultCachedGif           { return nil }
func (impl *InlineQueryResultVoice) OptCachedMpeg4Gif() *InlineQueryResultCachedMpeg4Gif { return nil }
func (impl *InlineQueryResultVoice) OptCachedPhoto() *InlineQueryResultCachedPhoto       { return nil }
func (impl *InlineQueryResultVoice) OptCachedSticker() *InlineQueryResultCachedSticker   { return nil }
func (impl *InlineQueryResultVoice) OptCachedVideo() *InlineQueryResultCachedVideo       { return nil }
func (impl *InlineQueryResultVoice) OptCachedVoice() *InlineQueryResultCachedVoice       { return nil }
func (impl *InlineQueryResultVoice) OptArticle() *InlineQueryResultArticle               { return nil }
func (impl *InlineQueryResultVoice) OptAudio() *InlineQueryResultAudio                   { return nil }
func (impl *InlineQueryResultVoice) OptContact() *InlineQueryResultContact               { return nil }
func (impl *InlineQueryResultVoice) OptGame() *InlineQueryResultGame                     { return nil }
func (impl *InlineQueryResultVoice) OptDocument() *InlineQueryResultDocument             { return nil }
func (impl *InlineQueryResultVoice) OptGif() *InlineQueryResultGif                       { return nil }
func (impl *InlineQueryResultVoice) OptLocation() *InlineQueryResultLocation             { return nil }
func (impl *InlineQueryResultVoice) OptMpeg4Gif() *InlineQueryResultMpeg4Gif             { return nil }
func (impl *InlineQueryResultVoice) OptPhoto() *InlineQueryResultPhoto                   { return nil }
func (impl *InlineQueryResultVoice) OptVenue() *InlineQueryResultVenue                   { return nil }
func (impl *InlineQueryResultVoice) OptVideo() *InlineQueryResultVideo                   { return nil }
func (impl *InlineQueryResultVoice) OptVoice() *InlineQueryResultVoice                   { return impl }

// InlineQueryResultArticle Represents a link to an article or web page.
type InlineQueryResultArticle struct {
	// Type of the result, must be article
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 Bytes
	Id string `json:"id"`
	// Title of the result
	Title string `json:"title"`
	// Content of the message to be sent
	InputMessageContent InputMessageContent `json:"input_message_content"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. URL of the result
	Url string `json:"url,omitempty"`
	// Optional. Pass True if you don't want the URL to be shown in the message
	HideUrl bool `json:"hide_url,omitempty"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
	// Optional. Url of the thumbnail for the result
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	// Optional. Thumbnail width
	ThumbnailWidth int64 `json:"thumbnail_width,omitempty"`
	// Optional. Thumbnail height
	ThumbnailHeight int64 `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultAudio Represents a link to an MP3 audio file. By default, this audio file will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
type InlineQueryResultAudio struct {
	// Type of the result, must be audio
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid URL for the audio file
	AudioUrl string `json:"audio_url"`
	// Title
	Title string `json:"title"`
	// Optional. Caption, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the audio caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Performer
	Performer string `json:"performer,omitempty"`
	// Optional. Audio duration in seconds
	AudioDuration int64 `json:"audio_duration,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the audio
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedAudio Represents a link to an MP3 audio file stored on the Telegram servers.
// By default, this audio file will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
type InlineQueryResultCachedAudio struct {
	// Type of the result, must be audio
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid file identifier for the audio file
	AudioFileId string `json:"audio_file_id"`
	// Optional. Caption, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the audio caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the audio
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (impl *InlineQueryResultCachedAudio) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.AudioFileId)
}
func (impl *InlineQueryResultCachedAudio) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.AudioFileId, dirAndPattern...)
}

// InlineQueryResultCachedDocument Represents a link to a file stored on the Telegram servers.
// By default, this file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the file.
type InlineQueryResultCachedDocument struct {
	// Type of the result, must be document
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// Title for the result
	Title string `json:"title"`
	// A valid file identifier for the file
	DocumentFileId string `json:"document_file_id"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
	// Optional. Caption of the document to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the document caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the file
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (impl *InlineQueryResultCachedDocument) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.DocumentFileId)
}
func (impl *InlineQueryResultCachedDocument) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.DocumentFileId, dirAndPattern...)
}

// InlineQueryResultCachedGif Represents a link to an animated GIF file stored on the Telegram servers.
// By default, this animated GIF file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with specified content instead of the animation.
type InlineQueryResultCachedGif struct {
	// Type of the result, must be gif
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid file identifier for the GIF file
	GifFileId string `json:"gif_file_id"`
	// Optional. Title for the result
	Title string `json:"title,omitempty"`
	// Optional. Caption of the GIF file to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the GIF animation
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (impl *InlineQueryResultCachedGif) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.GifFileId)
}
func (impl *InlineQueryResultCachedGif) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.GifFileId, dirAndPattern...)
}

// InlineQueryResultCachedMpeg4Gif Represents a link to a video animation (H.264/MPEG-4 AVC video without sound) stored on the Telegram servers.
// By default, this animated MPEG-4 file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultCachedMpeg4Gif struct {
	// Type of the result, must be mpeg4_gif
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid file identifier for the MPEG4 file
	Mpeg4FileId string `json:"mpeg4_file_id"`
	// Optional. Title for the result
	Title string `json:"title,omitempty"`
	// Optional. Caption of the MPEG-4 file to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the video animation
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (impl *InlineQueryResultCachedMpeg4Gif) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.Mpeg4FileId)
}
func (impl *InlineQueryResultCachedMpeg4Gif) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.Mpeg4FileId, dirAndPattern...)
}

// InlineQueryResultCachedPhoto Represents a link to a photo stored on the Telegram servers.
// By default, this photo will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.
type InlineQueryResultCachedPhoto struct {
	// Type of the result, must be photo
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid file identifier of the photo
	PhotoFileId string `json:"photo_file_id"`
	// Optional. Title for the result
	Title string `json:"title,omitempty"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
	// Optional. Caption of the photo to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the photo caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the photo
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (impl *InlineQueryResultCachedPhoto) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.PhotoFileId)
}
func (impl *InlineQueryResultCachedPhoto) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.PhotoFileId, dirAndPattern...)
}

// InlineQueryResultCachedSticker Represents a link to a sticker stored on the Telegram servers. By default, this sticker will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the sticker.
type InlineQueryResultCachedSticker struct {
	// Type of the result, must be sticker
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid file identifier of the sticker
	StickerFileId string `json:"sticker_file_id"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the sticker
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (impl *InlineQueryResultCachedSticker) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.StickerFileId)
}
func (impl *InlineQueryResultCachedSticker) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.StickerFileId, dirAndPattern...)
}

// InlineQueryResultCachedVideo Represents a link to a video file stored on the Telegram servers.
// By default, this video file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the video.
type InlineQueryResultCachedVideo struct {
	// Type of the result, must be video
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid file identifier for the video file
	VideoFileId string `json:"video_file_id"`
	// Title for the result
	Title string `json:"title"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
	// Optional. Caption of the video to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the video caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the video
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (impl *InlineQueryResultCachedVideo) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.VideoFileId)
}
func (impl *InlineQueryResultCachedVideo) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.VideoFileId, dirAndPattern...)
}

// InlineQueryResultCachedVoice Represents a link to a voice message stored on the Telegram servers.
// By default, this voice message will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the voice message.
type InlineQueryResultCachedVoice struct {
	// Type of the result, must be voice
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid file identifier for the voice message
	VoiceFileId string `json:"voice_file_id"`
	// Voice message title
	Title string `json:"title"`
	// Optional. Caption, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the voice message caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the voice message
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (impl *InlineQueryResultCachedVoice) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.VoiceFileId)
}
func (impl *InlineQueryResultCachedVoice) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.VoiceFileId, dirAndPattern...)
}

// InlineQueryResultContact Represents a contact with a phone number. By default, this contact will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the contact.
type InlineQueryResultContact struct {
	// Type of the result, must be contact
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 Bytes
	Id string `json:"id"`
	// Contact's phone number
	PhoneNumber string `json:"phone_number"`
	// Contact's first name
	FirstName string `json:"first_name"`
	// Optional. Contact's last name
	LastName string `json:"last_name,omitempty"`
	// Optional. Additional data about the contact in the form of a vCard, 0-2048 bytes
	Vcard string `json:"vcard,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the contact
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	// Optional. Url of the thumbnail for the result
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	// Optional. Thumbnail width
	ThumbnailWidth int64 `json:"thumbnail_width,omitempty"`
	// Optional. Thumbnail height
	ThumbnailHeight int64 `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultDocument Represents a link to a file. By default, this file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the file.
// Currently, only .PDF and .ZIP files can be sent using this method.
type InlineQueryResultDocument struct {
	// Type of the result, must be document
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// Title for the result
	Title string `json:"title"`
	// Optional. Caption of the document to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the document caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// A valid URL for the file
	DocumentUrl string `json:"document_url"`
	// MIME type of the content of the file, either "application/pdf" or "application/zip"
	MimeType string `json:"mime_type"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the file
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	// Optional. URL of the thumbnail (JPEG only) for the file
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	// Optional. Thumbnail width
	ThumbnailWidth int64 `json:"thumbnail_width,omitempty"`
	// Optional. Thumbnail height
	ThumbnailHeight int64 `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultGame Represents a Game.
type InlineQueryResultGame struct {
	// Type of the result, must be game
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// Short name of the game
	GameShortName string `json:"game_short_name"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// InlineQueryResultGif Represents a link to an animated GIF file.
// By default, this animated GIF file will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultGif struct {
	// Type of the result, must be gif
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid URL for the GIF file. File size must not exceed 1MB
	GifUrl string `json:"gif_url"`
	// Optional. Width of the GIF
	GifWidth int64 `json:"gif_width,omitempty"`
	// Optional. Height of the GIF
	GifHeight int64 `json:"gif_height,omitempty"`
	// Optional. Duration of the GIF in seconds
	GifDuration int64 `json:"gif_duration,omitempty"`
	// URL of the static (JPEG or GIF) or animated (MPEG4) thumbnail for the result
	ThumbnailUrl string `json:"thumbnail_url"`
	// Optional. MIME type of the thumbnail, must be one of "image/jpeg", "image/gif", or "video/mp4".
	// Defaults to "image/jpeg"
	ThumbnailMimeType string `json:"thumbnail_mime_type,omitempty"`
	// Optional. Title for the result
	Title string `json:"title,omitempty"`
	// Optional. Caption of the GIF file to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the GIF animation
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

// InlineQueryResultLocation Represents a location on a map. By default, the location will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the location.
type InlineQueryResultLocation struct {
	// Type of the result, must be location
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 Bytes
	Id string `json:"id"`
	// Location latitude in degrees
	Latitude float64 `json:"latitude"`
	// Location longitude in degrees
	Longitude float64 `json:"longitude"`
	// Location title
	Title string `json:"title"`
	// Optional. The radius of uncertainty for the location, measured in meters; 0-1500
	HorizontalAccuracy float64 `json:"horizontal_accuracy,omitempty"`
	// Optional.
	// Period in seconds during which the location can be updated, should be between 60 and 86400, or 0x7FFFFFFF for live locations that can be edited indefinitely.
	LivePeriod int64 `json:"live_period,omitempty"`
	// Optional. For live locations, a direction in which the user is moving, in degrees.
	// Must be between 1 and 360 if specified.
	Heading int64 `json:"heading,omitempty"`
	// Optional. Must be between 1 and 100000 if specified.
	// For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters.
	ProximityAlertRadius int64 `json:"proximity_alert_radius,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the location
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	// Optional. Url of the thumbnail for the result
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	// Optional. Thumbnail width
	ThumbnailWidth int64 `json:"thumbnail_width,omitempty"`
	// Optional. Thumbnail height
	ThumbnailHeight int64 `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultMpeg4Gif Represents a link to a video animation (H.264/MPEG-4 AVC video without sound).
// By default, this animated MPEG-4 file will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultMpeg4Gif struct {
	// Type of the result, must be mpeg4_gif
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid URL for the MPEG4 file. File size must not exceed 1MB
	Mpeg4Url string `json:"mpeg4_url"`
	// Optional. Video width
	Mpeg4Width int64 `json:"mpeg4_width,omitempty"`
	// Optional. Video height
	Mpeg4Height int64 `json:"mpeg4_height,omitempty"`
	// Optional. Video duration in seconds
	Mpeg4Duration int64 `json:"mpeg4_duration,omitempty"`
	// URL of the static (JPEG or GIF) or animated (MPEG4) thumbnail for the result
	ThumbnailUrl string `json:"thumbnail_url"`
	// Optional. MIME type of the thumbnail, must be one of "image/jpeg", "image/gif", or "video/mp4".
	// Defaults to "image/jpeg"
	ThumbnailMimeType string `json:"thumbnail_mime_type,omitempty"`
	// Optional. Title for the result
	Title string `json:"title,omitempty"`
	// Optional. Caption of the MPEG-4 file to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the video animation
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

// InlineQueryResultPhoto Represents a link to a photo. By default, this photo will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.
type InlineQueryResultPhoto struct {
	// Type of the result, must be photo
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid URL of the photo. Photo must be in JPEG format. Photo size must not exceed 5MB
	PhotoUrl string `json:"photo_url"`
	// URL of the thumbnail for the photo
	ThumbnailUrl string `json:"thumbnail_url"`
	// Optional. Width of the photo
	PhotoWidth int64 `json:"photo_width,omitempty"`
	// Optional. Height of the photo
	PhotoHeight int64 `json:"photo_height,omitempty"`
	// Optional. Title for the result
	Title string `json:"title,omitempty"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
	// Optional. Caption of the photo to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the photo caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the photo
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

// InlineQueryResultVenue Represents a venue. By default, the venue will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the venue.
type InlineQueryResultVenue struct {
	// Type of the result, must be venue
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 Bytes
	Id string `json:"id"`
	// Latitude of the venue location in degrees
	Latitude float64 `json:"latitude"`
	// Longitude of the venue location in degrees
	Longitude float64 `json:"longitude"`
	// Title of the venue
	Title string `json:"title"`
	// Address of the venue
	Address string `json:"address"`
	// Optional. Foursquare identifier of the venue if known
	FoursquareId string `json:"foursquare_id,omitempty"`
	// Optional. Foursquare type of the venue, if known.
	// (For example, "arts_entertainment/default", "arts_entertainment/aquarium" or "food/icecream".)
	FoursquareType string `json:"foursquare_type,omitempty"`
	// Optional. Google Places identifier of the venue
	GooglePlaceId string `json:"google_place_id,omitempty"`
	// Optional. Google Places type of the venue. (See supported types.)
	GooglePlaceType string `json:"google_place_type,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the venue
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	// Optional. Url of the thumbnail for the result
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	// Optional. Thumbnail width
	ThumbnailWidth int64 `json:"thumbnail_width,omitempty"`
	// Optional. Thumbnail height
	ThumbnailHeight int64 `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultVideo Represents a link to a page containing an embedded video player or a video file.
// By default, this video file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the video.
type InlineQueryResultVideo struct {
	// Type of the result, must be video
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid URL for the embedded video player or video file
	VideoUrl string `json:"video_url"`
	// MIME type of the content of the video URL, "text/html" or "video/mp4"
	MimeType string `json:"mime_type"`
	// URL of the thumbnail (JPEG only) for the video
	ThumbnailUrl string `json:"thumbnail_url"`
	// Title for the result
	Title string `json:"title"`
	// Optional. Caption of the video to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the video caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Video width
	VideoWidth int64 `json:"video_width,omitempty"`
	// Optional. Video height
	VideoHeight int64 `json:"video_height,omitempty"`
	// Optional. Video duration in seconds
	VideoDuration int64 `json:"video_duration,omitempty"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the video.
	// This field is required if InlineQueryResultVideo is used to send an HTML-page as a result (e.g., a YouTube video).
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

// InlineQueryResultVoice Represents a link to a voice recording in an .OGG container encoded with OPUS.
// By default, this voice recording will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the the voice message.
type InlineQueryResultVoice struct {
	// Type of the result, must be voice
	Type string `json:"type"`
	// Unique identifier for this result, 1-64 bytes
	Id string `json:"id"`
	// A valid URL for the voice recording
	VoiceUrl string `json:"voice_url"`
	// Recording title
	Title string `json:"title"`
	// Optional. Caption, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the voice message caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Recording duration in seconds
	VoiceDuration int64 `json:"voice_duration,omitempty"`
	// Optional. Inline keyboard attached to the message
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// Optional. Content of the message to be sent instead of the voice recording
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

// InlineQueryResultsButton This object represents a button to be shown above inline query results.
// You must use exactly one of the optional fields.
type InlineQueryResultsButton struct {
	// Label text on the button
	Text string `json:"text"`
	// Optional. Description of the Web App that will be launched when the user presses the button.
	// The Web App will be able to switch back to the inline mode using the method switchInlineQuery inside the Web App.
	WebApp *WebAppInfo `json:"web_app,omitempty"`
	// Optional. Deep-linking parameter for the /start message sent to the bot when a user presses the button.
	// 1-64 characters, only A-Z, a-z, 0-9, _ and - are allowed.
	// Example: An inline bot that sends YouTube videos can ask the user to connect the bot to their YouTube account to adapt search results accordingly.
	// To do this, it displays a 'Connect your YouTube account' button above the results, or even before showing any.
	// The user presses the button, switches to a private chat with the bot and, in doing so, passes a start parameter that instructs the bot to return an OAuth link.
	// Once done, the bot can offer a switch_inline button so that the user can easily return to the chat where they wanted to use the bot's inline capabilities.
	StartParameter string `json:"start_parameter,omitempty"`
}

// InputContactMessageContent Represents the content of a contact message to be sent as the result of an inline query.
type InputContactMessageContent struct {
	// Contact's phone number
	PhoneNumber string `json:"phone_number"`
	// Contact's first name
	FirstName string `json:"first_name"`
	// Optional. Contact's last name
	LastName string `json:"last_name,omitempty"`
	// Optional. Additional data about the contact in the form of a vCard, 0-2048 bytes
	Vcard string `json:"vcard,omitempty"`
}

// InputInvoiceMessageContent Represents the content of an invoice message to be sent as the result of an inline query.
type InputInvoiceMessageContent struct {
	// Product name, 1-32 characters
	Title string `json:"title"`
	// Product description, 1-255 characters
	Description string `json:"description"`
	// Bot-defined invoice payload, 1-128 bytes. This will not be displayed to the user, use it for your internal processes.
	Payload string `json:"payload"`
	// Optional. Payment provider token, obtained via @BotFather. Pass an empty string for payments in Telegram Stars.
	ProviderToken string `json:"provider_token,omitempty"`
	// Three-letter ISO 4217 currency code, see more on currencies. Pass "XTR" for payments in Telegram Stars.
	Currency string `json:"currency"`
	// Price breakdown, a JSON-serialized list of components (e.g.
	// product price, tax, discount, delivery cost, delivery tax, bonus, etc.).
	// Must contain exactly one item for payments in Telegram Stars.
	Prices []*LabeledPrice `json:"prices"`
	// Optional. The maximum accepted amount for tips in the smallest units of the currency (integer, not float/double).
	// For example, for a maximum tip of US$ 1.45 pass max_tip_amount = 145. Defaults to 0.
	// See the exp parameter in currencies.json, it shows the number of digits past the decimal point for each currency (2 for the majority of currencies).
	// Not supported for payments in Telegram Stars.
	MaxTipAmount int64 `json:"max_tip_amount,omitempty"`
	// Optional. At most 4 suggested tip amounts can be specified.
	// A JSON-serialized array of suggested amounts of tip in the smallest units of the currency (integer, not float/double).
	// The suggested tip amounts must be positive, passed in a strictly increased order and must not exceed max_tip_amount.
	SuggestedTipAmounts []int64 `json:"suggested_tip_amounts,omitempty"`
	// Optional. A JSON-serialized object for data about the invoice, which will be shared with the payment provider.
	// A detailed description of the required fields should be provided by the payment provider.
	ProviderData string `json:"provider_data,omitempty"`
	// Optional. URL of the product photo for the invoice. Can be a photo of the goods or a marketing image for a service.
	PhotoUrl string `json:"photo_url,omitempty"`
	// Optional. Photo size in bytes
	PhotoSize int64 `json:"photo_size,omitempty"`
	// Optional. Photo width
	PhotoWidth int64 `json:"photo_width,omitempty"`
	// Optional. Photo height
	PhotoHeight int64 `json:"photo_height,omitempty"`
	// Optional. Pass True if you require the user's full name to complete the order.
	// Ignored for payments in Telegram Stars.
	NeedName bool `json:"need_name,omitempty"`
	// Optional. Pass True if you require the user's phone number to complete the order.
	// Ignored for payments in Telegram Stars.
	NeedPhoneNumber bool `json:"need_phone_number,omitempty"`
	// Optional. Pass True if you require the user's email address to complete the order.
	// Ignored for payments in Telegram Stars.
	NeedEmail bool `json:"need_email,omitempty"`
	// Optional. Pass True if you require the user's shipping address to complete the order.
	// Ignored for payments in Telegram Stars.
	NeedShippingAddress bool `json:"need_shipping_address,omitempty"`
	// Optional. Pass True if the user's phone number should be sent to the provider.
	// Ignored for payments in Telegram Stars.
	SendPhoneNumberToProvider bool `json:"send_phone_number_to_provider,omitempty"`
	// Optional. Pass True if the user's email address should be sent to the provider.
	// Ignored for payments in Telegram Stars.
	SendEmailToProvider bool `json:"send_email_to_provider,omitempty"`
	// Optional. Pass True if the final price depends on the shipping method. Ignored for payments in Telegram Stars.
	IsFlexible bool `json:"is_flexible,omitempty"`
}

// InputLocationMessageContent Represents the content of a location message to be sent as the result of an inline query.
type InputLocationMessageContent struct {
	// Latitude of the location in degrees
	Latitude float64 `json:"latitude"`
	// Longitude of the location in degrees
	Longitude float64 `json:"longitude"`
	// Optional. The radius of uncertainty for the location, measured in meters; 0-1500
	HorizontalAccuracy float64 `json:"horizontal_accuracy,omitempty"`
	// Optional.
	// Period in seconds during which the location can be updated, should be between 60 and 86400, or 0x7FFFFFFF for live locations that can be edited indefinitely.
	LivePeriod int64 `json:"live_period,omitempty"`
	// Optional. For live locations, a direction in which the user is moving, in degrees.
	// Must be between 1 and 360 if specified.
	Heading int64 `json:"heading,omitempty"`
	// Optional. Must be between 1 and 100000 if specified.
	// For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters.
	ProximityAlertRadius int64 `json:"proximity_alert_radius,omitempty"`
}

// InputMedia This object represents the content of a media message to be sent. It should be one of
// - InputMediaAnimation
// - InputMediaDocument
// - InputMediaAudio
// - InputMediaPhoto
// - InputMediaVideo
type InputMedia interface {
	OptAnimation() *Animation
	OptDocument() *Document
	OptAudio() *Audio
	OptPhoto() *Photo
	OptVideo() *Video
}

var (
	_ InputMedia = &Animation{}
	_ InputMedia = &Document{}
	_ InputMedia = &Audio{}
	_ InputMedia = &Photo{}
	_ InputMedia = &Video{}
)

func (impl *Animation) OptAnimation() *Animation { return impl }
func (impl *Animation) OptDocument() *Document   { return nil }
func (impl *Animation) OptAudio() *Audio         { return nil }
func (impl *Animation) OptPhoto() *Photo         { return nil }
func (impl *Animation) OptVideo() *Video         { return nil }

func (impl *Document) OptAnimation() *Animation { return nil }
func (impl *Document) OptDocument() *Document   { return impl }
func (impl *Document) OptAudio() *Audio         { return nil }
func (impl *Document) OptPhoto() *Photo         { return nil }
func (impl *Document) OptVideo() *Video         { return nil }

func (impl *Audio) OptAnimation() *Animation { return nil }
func (impl *Audio) OptDocument() *Document   { return nil }
func (impl *Audio) OptAudio() *Audio         { return impl }
func (impl *Audio) OptPhoto() *Photo         { return nil }
func (impl *Audio) OptVideo() *Video         { return nil }

func (impl *Photo) OptAnimation() *Animation { return nil }
func (impl *Photo) OptDocument() *Document   { return nil }
func (impl *Photo) OptAudio() *Audio         { return nil }
func (impl *Photo) OptPhoto() *Photo         { return impl }
func (impl *Photo) OptVideo() *Video         { return nil }

func (impl *Video) OptAnimation() *Animation { return nil }
func (impl *Video) OptDocument() *Document   { return nil }
func (impl *Video) OptAudio() *Audio         { return nil }
func (impl *Video) OptPhoto() *Photo         { return nil }
func (impl *Video) OptVideo() *Video         { return impl }

// InputMessageContent This object represents the content of a message to be sent as a result of an inline query. Telegram clients currently support the following 5 types:
// - InputTextMessageContent
// - InputLocationMessageContent
// - InputVenueMessageContent
// - InputContactMessageContent
// - InputInvoiceMessageContent
type InputMessageContent interface {
	OptInputTextMessageContent() *InputTextMessageContent
	OptInputLocationMessageContent() *InputLocationMessageContent
	OptInputVenueMessageContent() *InputVenueMessageContent
	OptInputContactMessageContent() *InputContactMessageContent
	OptInputInvoiceMessageContent() *InputInvoiceMessageContent
}

var (
	_ InputMessageContent = &InputTextMessageContent{}
	_ InputMessageContent = &InputLocationMessageContent{}
	_ InputMessageContent = &InputVenueMessageContent{}
	_ InputMessageContent = &InputContactMessageContent{}
	_ InputMessageContent = &InputInvoiceMessageContent{}
)

func (impl *InputTextMessageContent) OptInputTextMessageContent() *InputTextMessageContent {
	return impl
}
func (impl *InputTextMessageContent) OptInputLocationMessageContent() *InputLocationMessageContent {
	return nil
}
func (impl *InputTextMessageContent) OptInputVenueMessageContent() *InputVenueMessageContent {
	return nil
}
func (impl *InputTextMessageContent) OptInputContactMessageContent() *InputContactMessageContent {
	return nil
}
func (impl *InputTextMessageContent) OptInputInvoiceMessageContent() *InputInvoiceMessageContent {
	return nil
}

func (impl *InputLocationMessageContent) OptInputTextMessageContent() *InputTextMessageContent {
	return nil
}
func (impl *InputLocationMessageContent) OptInputLocationMessageContent() *InputLocationMessageContent {
	return impl
}
func (impl *InputLocationMessageContent) OptInputVenueMessageContent() *InputVenueMessageContent {
	return nil
}
func (impl *InputLocationMessageContent) OptInputContactMessageContent() *InputContactMessageContent {
	return nil
}
func (impl *InputLocationMessageContent) OptInputInvoiceMessageContent() *InputInvoiceMessageContent {
	return nil
}

func (impl *InputVenueMessageContent) OptInputTextMessageContent() *InputTextMessageContent {
	return nil
}
func (impl *InputVenueMessageContent) OptInputLocationMessageContent() *InputLocationMessageContent {
	return nil
}
func (impl *InputVenueMessageContent) OptInputVenueMessageContent() *InputVenueMessageContent {
	return impl
}
func (impl *InputVenueMessageContent) OptInputContactMessageContent() *InputContactMessageContent {
	return nil
}
func (impl *InputVenueMessageContent) OptInputInvoiceMessageContent() *InputInvoiceMessageContent {
	return nil
}

func (impl *InputContactMessageContent) OptInputTextMessageContent() *InputTextMessageContent {
	return nil
}
func (impl *InputContactMessageContent) OptInputLocationMessageContent() *InputLocationMessageContent {
	return nil
}
func (impl *InputContactMessageContent) OptInputVenueMessageContent() *InputVenueMessageContent {
	return nil
}
func (impl *InputContactMessageContent) OptInputContactMessageContent() *InputContactMessageContent {
	return impl
}
func (impl *InputContactMessageContent) OptInputInvoiceMessageContent() *InputInvoiceMessageContent {
	return nil
}

func (impl *InputInvoiceMessageContent) OptInputTextMessageContent() *InputTextMessageContent {
	return nil
}
func (impl *InputInvoiceMessageContent) OptInputLocationMessageContent() *InputLocationMessageContent {
	return nil
}
func (impl *InputInvoiceMessageContent) OptInputVenueMessageContent() *InputVenueMessageContent {
	return nil
}
func (impl *InputInvoiceMessageContent) OptInputContactMessageContent() *InputContactMessageContent {
	return nil
}
func (impl *InputInvoiceMessageContent) OptInputInvoiceMessageContent() *InputInvoiceMessageContent {
	return impl
}

// InputPaidMedia This object describes the paid media to be sent. Currently, it can be one of
// - InputPaidMediaPhoto
// - InputPaidMediaVideo
type InputPaidMedia interface {
	OptPhoto() *InputPaidMediaPhoto
	OptVideo() *InputPaidMediaVideo
}

var (
	_ InputPaidMedia = &InputPaidMediaPhoto{}
	_ InputPaidMedia = &InputPaidMediaVideo{}
)

func (impl *InputPaidMediaPhoto) OptPhoto() *InputPaidMediaPhoto { return impl }
func (impl *InputPaidMediaPhoto) OptVideo() *InputPaidMediaVideo { return nil }

func (impl *InputPaidMediaVideo) OptPhoto() *InputPaidMediaPhoto { return nil }
func (impl *InputPaidMediaVideo) OptVideo() *InputPaidMediaVideo { return impl }

// InputPaidMediaPhoto The paid media to send is a photo.
type InputPaidMediaPhoto struct {
	// Type of the media, must be photo
	Type string `json:"type"`
	// File to send. More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// Pass a file_id to send a file that exists on the Telegram servers (recommended), pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	Media string `json:"media"`
}

// InputPaidMediaVideo The paid media to send is a video.
type InputPaidMediaVideo struct {
	// Type of the media, must be video
	Type string `json:"type"`
	// File to send. More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// Pass a file_id to send a file that exists on the Telegram servers (recommended), pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	Media string `json:"media"`
	// Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size.
	// A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data.
	// Thumbnails can't be reused and can be only uploaded as a new file, so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	// More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// >> either: String
	Thumbnail InputFile `json:"thumbnail,omitempty"`
	// Optional. Video width
	Width int64 `json:"width,omitempty"`
	// Optional. Video height
	Height int64 `json:"height,omitempty"`
	// Optional. Video duration in seconds
	Duration int64 `json:"duration,omitempty"`
	// Optional. Pass True if the uploaded video is suitable for streaming
	SupportsStreaming bool `json:"supports_streaming,omitempty"`
}

// InputPollOption This object contains information about one answer option in a poll to be sent.
type InputPollOption struct {
	// Option text, 1-100 characters
	Text string `json:"text"`
	// Optional. Mode for parsing entities in the text. See formatting options for more details.
	// Currently, only custom emoji entities are allowed
	TextParseMode string `json:"text_parse_mode,omitempty"`
	// Optional. A JSON-serialized list of special entities that appear in the poll option text.
	// It can be specified instead of text_parse_mode
	TextEntities []*MessageEntity `json:"text_entities,omitempty"`
}

// InputSticker This object describes a sticker to be added to a sticker set.
type InputSticker struct {
	// The added sticker. Animated and video stickers can't be uploaded via HTTP URL.
	// Pass a file_id as a String to send a file that already exists on the Telegram servers, pass an HTTP URL as a String for Telegram to get a file from the Internet, upload a new one using multipart/form-data, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	// More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// >> either: String
	Sticker InputFile `json:"sticker"`
	// Format of the added sticker, must be one of "static" for a .WEBP or .PNG image, "animated" for a .TGS animation, "video" for a WEBM video
	Format string `json:"format"`
	// List of 1-20 emoji associated with the sticker
	EmojiList []string `json:"emoji_list"`
	// Optional. Position where the mask should be placed on faces. For "mask" stickers only.
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	// Optional. List of 0-20 search keywords for the sticker with total length of up to 64 characters.
	// For "regular" and "custom_emoji" stickers only.
	Keywords []string `json:"keywords,omitempty"`
}

// InputTextMessageContent Represents the content of a text message to be sent as the result of an inline query.
type InputTextMessageContent struct {
	// Text of the message to be sent, 1-4096 characters
	MessageText string `json:"message_text"`
	// Optional. Mode for parsing entities in the message text. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in message text, which can be specified instead of parse_mode
	Entities []*MessageEntity `json:"entities,omitempty"`
	// Optional. Link preview generation options for the message
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

// InputVenueMessageContent Represents the content of a venue message to be sent as the result of an inline query.
type InputVenueMessageContent struct {
	// Latitude of the venue in degrees
	Latitude float64 `json:"latitude"`
	// Longitude of the venue in degrees
	Longitude float64 `json:"longitude"`
	// Name of the venue
	Title string `json:"title"`
	// Address of the venue
	Address string `json:"address"`
	// Optional. Foursquare identifier of the venue, if known
	FoursquareId string `json:"foursquare_id,omitempty"`
	// Optional. Foursquare type of the venue, if known.
	// (For example, "arts_entertainment/default", "arts_entertainment/aquarium" or "food/icecream".)
	FoursquareType string `json:"foursquare_type,omitempty"`
	// Optional. Google Places identifier of the venue
	GooglePlaceId string `json:"google_place_id,omitempty"`
	// Optional. Google Places type of the venue. (See supported types.)
	GooglePlaceType string `json:"google_place_type,omitempty"`
}

// Invoice This object contains basic information about an invoice.
type Invoice struct {
	// Product name
	Title string `json:"title"`
	// Product description
	Description string `json:"description"`
	// Unique bot deep-linking parameter that can be used to generate this invoice
	StartParameter string `json:"start_parameter"`
	// Three-letter ISO 4217 currency code, or "XTR" for payments in Telegram Stars
	Currency string `json:"currency"`
	// Total price in the smallest units of the currency (integer, not float/double).
	// For example, for a price of US$ 1.45 pass amount = 145.
	// See the exp parameter in currencies.json, it shows the number of digits past the decimal point for each currency (2 for the majority of currencies).
	TotalAmount int64 `json:"total_amount"`
}

// KeyboardButton This object represents one button of the reply keyboard. At most one of the optional fields must be used to specify type of the button. For simple text buttons, String can be used instead of this object to specify the button text.
// Note: request_users and request_chat options will only work in Telegram versions released after 3 February, 2023. Older clients will display unsupported message.
type KeyboardButton struct {
	// Text of the button. If none of the optional fields are used, it will be sent as a message when the button is pressed
	Text string `json:"text"`
	// Optional. If specified, pressing the button will open a list of suitable users. Available in private chats only.
	// Identifiers of selected users will be sent to the bot in a "users_shared" service message.
	RequestUsers *KeyboardButtonRequestUsers `json:"request_users,omitempty"`
	// Optional. If specified, pressing the button will open a list of suitable chats. Available in private chats only.
	// Tapping on a chat will send its identifier to the bot in a "chat_shared" service message.
	RequestChat *KeyboardButtonRequestChat `json:"request_chat,omitempty"`
	// Optional. If True, the user's phone number will be sent as a contact when the button is pressed.
	// Available in private chats only.
	RequestContact bool `json:"request_contact,omitempty"`
	// Optional. If True, the user's current location will be sent when the button is pressed.
	// Available in private chats only.
	RequestLocation bool `json:"request_location,omitempty"`
	// Optional. If specified, the user will be asked to create a poll and send it to the bot when the button is pressed.
	// Available in private chats only.
	RequestPoll *KeyboardButtonPollType `json:"request_poll,omitempty"`
	// Optional. If specified, the described Web App will be launched when the button is pressed.
	// The Web App will be able to send a "web_app_data" service message. Available in private chats only.
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

// KeyboardButtonPollType This object represents type of a poll, which is allowed to be created and sent when the corresponding button is pressed.
type KeyboardButtonPollType struct {
	// Optional. If quiz is passed, the user will be allowed to create only polls in the quiz mode.
	// If regular is passed, only regular polls will be allowed.
	// Otherwise, the user will be allowed to create a poll of any type.
	Type string `json:"type,omitempty"`
}

// KeyboardButtonRequestChat This object defines the criteria used to request a suitable chat.
// Information about the selected chat will be shared with the bot when the corresponding button is pressed.
// The bot will be granted requested rights in the chat if appropriate.
// More about requesting chats: https://core.telegram.org/bots/features#chat-and-user-selection.
type KeyboardButtonRequestChat struct {
	// Signed 32-bit identifier of the request, which will be received back in the ChatShared object.
	// Must be unique within the message
	RequestId int64 `json:"request_id"`
	// Pass True to request a channel chat, pass False to request a group or a supergroup chat.
	ChatIsChannel bool `json:"chat_is_channel"`
	// Optional. Pass True to request a forum supergroup, pass False to request a non-forum chat.
	// If not specified, no additional restrictions are applied.
	ChatIsForum bool `json:"chat_is_forum,omitempty"`
	// Optional. If not specified, no additional restrictions are applied.
	// Pass True to request a supergroup or a channel with a username, pass False to request a chat without a username.
	ChatHasUsername bool `json:"chat_has_username,omitempty"`
	// Optional. Pass True to request a chat owned by the user. Otherwise, no additional restrictions are applied.
	ChatIsCreated bool `json:"chat_is_created,omitempty"`
	// Optional. A JSON-serialized object listing the required administrator rights of the user in the chat.
	// The rights must be a superset of bot_administrator_rights. If not specified, no additional restrictions are applied.
	UserAdministratorRights *ChatAdministratorRights `json:"user_administrator_rights,omitempty"`
	// Optional. A JSON-serialized object listing the required administrator rights of the bot in the chat.
	// The rights must be a subset of user_administrator_rights. If not specified, no additional restrictions are applied.
	BotAdministratorRights *ChatAdministratorRights `json:"bot_administrator_rights,omitempty"`
	// Optional. Pass True to request a chat with the bot as a member. Otherwise, no additional restrictions are applied.
	BotIsMember bool `json:"bot_is_member,omitempty"`
	// Optional. Pass True to request the chat's title
	RequestTitle bool `json:"request_title,omitempty"`
	// Optional. Pass True to request the chat's username
	RequestUsername bool `json:"request_username,omitempty"`
	// Optional. Pass True to request the chat's photo
	RequestPhoto bool `json:"request_photo,omitempty"`
}

// KeyboardButtonRequestUsers This object defines the criteria used to request suitable users.
// Information about the selected users will be shared with the bot when the corresponding button is pressed.
// More about requesting users: https://core.telegram.org/bots/features#chat-and-user-selection
type KeyboardButtonRequestUsers struct {
	// Signed 32-bit identifier of the request that will be received back in the UsersShared object.
	// Must be unique within the message
	RequestId int64 `json:"request_id"`
	// Optional. Pass True to request bots, pass False to request regular users.
	// If not specified, no additional restrictions are applied.
	UserIsBot bool `json:"user_is_bot,omitempty"`
	// Optional. Pass True to request premium users, pass False to request non-premium users.
	// If not specified, no additional restrictions are applied.
	UserIsPremium bool `json:"user_is_premium,omitempty"`
	// Optional. The maximum number of users to be selected; 1-10. Defaults to 1.
	MaxQuantity int64 `json:"max_quantity,omitempty"`
	// Optional. Pass True to request the users' first and last names
	RequestName bool `json:"request_name,omitempty"`
	// Optional. Pass True to request the users' usernames
	RequestUsername bool `json:"request_username,omitempty"`
	// Optional. Pass True to request the users' photos
	RequestPhoto bool `json:"request_photo,omitempty"`
}

// LabeledPrice This object represents a portion of the price for goods or services.
type LabeledPrice struct {
	// Portion label
	Label string `json:"label"`
	// Price of the product in the smallest units of the currency (integer, not float/double).
	// For example, for a price of US$ 1.45 pass amount = 145.
	// See the exp parameter in currencies.json, it shows the number of digits past the decimal point for each currency (2 for the majority of currencies).
	Amount int64 `json:"amount"`
}

// LinkPreviewOptions Describes the options used for link preview generation.
type LinkPreviewOptions struct {
	// Optional. True, if the link preview is disabled
	IsDisabled bool `json:"is_disabled,omitempty"`
	// Optional. URL to use for the link preview. If empty, then the first URL found in the message text will be used
	Url string `json:"url,omitempty"`
	// Optional.
	// True, if the media in the link preview is supposed to be shrunk; ignored if the URL isn't explicitly specified or media size change isn't supported for the preview
	PreferSmallMedia bool `json:"prefer_small_media,omitempty"`
	// Optional.
	// True, if the media in the link preview is supposed to be enlarged; ignored if the URL isn't explicitly specified or media size change isn't supported for the preview
	PreferLargeMedia bool `json:"prefer_large_media,omitempty"`
	// Optional.
	// True, if the link preview must be shown above the message text; otherwise, the link preview will be shown below the message text
	ShowAboveText bool `json:"show_above_text,omitempty"`
}

// Location This object represents a point on the map.
type Location struct {
	// Latitude as defined by the sender
	Latitude float64 `json:"latitude"`
	// Longitude as defined by the sender
	Longitude float64 `json:"longitude"`
	// Optional. The radius of uncertainty for the location, measured in meters; 0-1500
	HorizontalAccuracy float64 `json:"horizontal_accuracy,omitempty"`
	// Optional. Time relative to the message sending date, during which the location can be updated; in seconds.
	// For active live locations only.
	LivePeriod int64 `json:"live_period,omitempty"`
	// Optional. The direction in which user is moving, in degrees; 1-360. For active live locations only.
	Heading int64 `json:"heading,omitempty"`
	// Optional. The maximum distance for proximity alerts about approaching another chat member, in meters.
	// For sent live locations only.
	ProximityAlertRadius int64 `json:"proximity_alert_radius,omitempty"`
}

// LoginUrl This object represents a parameter of the inline keyboard button used to automatically authorize a user. Serves as a great replacement for the Telegram Login Widget when the user is coming from Telegram. All the user needs to do is tap/click a button and confirm that they want to log in:
// Telegram apps support these buttons as of version 5.7.
type LoginUrl struct {
	// An HTTPS URL to be opened with user authorization data added to the query string when the button is pressed.
	// If the user refuses to provide authorization data, the original URL without information about the user will be opened.
	// The data added is the same as described in Receiving authorization data.
	// NOTE: You must always check the hash of the received data to verify the authentication and the integrity of the data as described in Checking authorization.
	Url string `json:"url"`
	// Optional. New text of the button in forwarded messages.
	ForwardText string `json:"forward_text,omitempty"`
	// Optional. Username of a bot, which will be used for user authorization. See Setting up a bot for more details.
	// If not specified, the current bot's username will be assumed. See Linking your domain to the bot for more details.
	// The url's domain must be the same as the domain linked with the bot.
	BotUsername string `json:"bot_username,omitempty"`
	// Optional. Pass True to request the permission for your bot to send messages to the user.
	RequestWriteAccess bool `json:"request_write_access,omitempty"`
}

// MaskPosition This object describes the position on faces where a mask should be placed by default.
type MaskPosition struct {
	// The part of the face relative to which the mask should be placed. One of "forehead", "eyes", "mouth", or "chin".
	Point string `json:"point"`
	// Shift by X-axis measured in widths of the mask scaled to the face size, from left to right.
	// For example, choosing -1.0 will place mask just to the left of the default mask position.
	XShift float64 `json:"x_shift"`
	// Shift by Y-axis measured in heights of the mask scaled to the face size, from top to bottom.
	// For example, 1.0 will place the mask just below the default mask position.
	YShift float64 `json:"y_shift"`
	// Mask scaling coefficient. For example, 2.0 means double size.
	Scale float64 `json:"scale"`
}

// MenuButton This object describes the bot's menu button in a private chat. It should be one of
// - MenuButtonCommands
// - MenuButtonWebApp
// - MenuButtonDefault
// If a menu button other than MenuButtonDefault is set for a private chat, then it is applied in the chat. Otherwise the default menu button is applied. By default, the menu button opens the list of bot commands.
type MenuButton interface {
	OptCommands() *MenuButtonCommands
	OptWebApp() *MenuButtonWebApp
	OptDefault() *MenuButtonDefault
}

var (
	_ MenuButton = &MenuButtonCommands{}
	_ MenuButton = &MenuButtonWebApp{}
	_ MenuButton = &MenuButtonDefault{}
)

func (impl *MenuButtonCommands) OptCommands() *MenuButtonCommands { return impl }
func (impl *MenuButtonCommands) OptWebApp() *MenuButtonWebApp     { return nil }
func (impl *MenuButtonCommands) OptDefault() *MenuButtonDefault   { return nil }

func (impl *MenuButtonWebApp) OptCommands() *MenuButtonCommands { return nil }
func (impl *MenuButtonWebApp) OptWebApp() *MenuButtonWebApp     { return impl }
func (impl *MenuButtonWebApp) OptDefault() *MenuButtonDefault   { return nil }

func (impl *MenuButtonDefault) OptCommands() *MenuButtonCommands { return nil }
func (impl *MenuButtonDefault) OptWebApp() *MenuButtonWebApp     { return nil }
func (impl *MenuButtonDefault) OptDefault() *MenuButtonDefault   { return impl }

// MenuButtonCommands Represents a menu button, which opens the bot's list of commands.
type MenuButtonCommands struct {
	// Type of the button, must be commands
	Type string `json:"type" default:"commands"`
}

// MenuButtonDefault Describes that no specific value for the menu button was set.
type MenuButtonDefault struct {
	// Type of the button, must be default
	Type string `json:"type" default:"default"`
}

// MenuButtonWebApp Represents a menu button, which launches a Web App.
type MenuButtonWebApp struct {
	// Type of the button, must be web_app
	Type string `json:"type" default:"web_app"`
	// Text on the button
	Text string `json:"text"`
	// Description of the Web App that will be launched when the user presses the button.
	// The Web App will be able to send an arbitrary message on behalf of the user using the method answerWebAppQuery.
	// Alternatively, a t.me link to a Web App of the bot can be specified in the object instead of the Web App's URL, in which case the Web App will be opened as if the user pressed the link.
	WebApp *WebAppInfo `json:"web_app"`
}

// Message This object represents a message.
type Message struct {
	// Unique message identifier inside this chat.
	// In specific instances (e.g., message containing a video sent to a big chat), the server might automatically schedule a message instead of sending it immediately.
	// In such cases, this field will be 0 and the relevant message will be unusable until it is actually sent
	MessageId int64 `json:"message_id"`
	// Optional. Unique identifier of a message thread to which the message belongs; for supergroups only
	MessageThreadId int64 `json:"message_thread_id,omitempty"`
	// Optional. Sender of the message; may be empty for messages sent to channels.
	// For backward compatibility, if the message was sent on behalf of a chat, the field contains a fake sender user in non-channel chats
	From *User `json:"from,omitempty"`
	// Optional. Sender of the message when sent on behalf of a chat.
	// For example, the supergroup itself for messages sent by its anonymous administrators or a linked channel for messages automatically forwarded to the channel's discussion group.
	// For backward compatibility, if the message was sent on behalf of a chat, the field from contains a fake sender user in non-channel chats.
	SenderChat *Chat `json:"sender_chat,omitempty"`
	// Optional. If the sender of the message boosted the chat, the number of boosts added by the user
	SenderBoostCount int64 `json:"sender_boost_count,omitempty"`
	// Optional. The bot that actually sent the message on behalf of the business account.
	// Available only for outgoing messages sent on behalf of the connected business account.
	SenderBusinessBot *User `json:"sender_business_bot,omitempty"`
	// Date the message was sent in Unix time. It is always a positive number, representing a valid date.
	Date int64 `json:"date"`
	// Optional. Unique identifier of the business connection from which the message was received.
	// If non-empty, the message belongs to a chat of the corresponding business account that is independent from any potential bot chat which might share the same identifier.
	BusinessConnectionId string `json:"business_connection_id,omitempty"`
	// Chat the message belongs to
	Chat *Chat `json:"chat"`
	// Optional. Information about the original message for forwarded messages
	ForwardOrigin MessageOrigin `json:"forward_origin,omitempty"`
	// Optional. True, if the message is sent to a forum topic
	IsTopicMessage bool `json:"is_topic_message,omitempty"`
	// Optional. True, if the message is a channel post that was automatically forwarded to the connected discussion group
	IsAutomaticForward bool `json:"is_automatic_forward,omitempty"`
	// Optional. For replies in the same chat and message thread, the original message.
	// Note that the Message object in this field will not contain further reply_to_message fields even if it itself is a reply.
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`
	// Optional. Information about the message that is being replied to, which may come from another chat or forum topic
	ExternalReply *ExternalReplyInfo `json:"external_reply,omitempty"`
	// Optional. For replies that quote part of the original message, the quoted part of the message
	Quote *TextQuote `json:"quote,omitempty"`
	// Optional. For replies to a story, the original story
	ReplyToStory *Story `json:"reply_to_story,omitempty"`
	// Optional. Bot through which the message was sent
	ViaBot *User `json:"via_bot,omitempty"`
	// Optional. Date the message was last edited in Unix time
	EditDate int64 `json:"edit_date,omitempty"`
	// Optional. True, if the message can't be forwarded
	HasProtectedContent bool `json:"has_protected_content,omitempty"`
	// Optional.
	// True, if the message was sent by an implicit action, for example, as an away or a greeting business message, or as a scheduled message
	IsFromOffline bool `json:"is_from_offline,omitempty"`
	// Optional. The unique identifier of a media message group this message belongs to
	MediaGroupId string `json:"media_group_id,omitempty"`
	// Optional.
	// Signature of the post author for messages in channels, or the custom title of an anonymous group administrator
	AuthorSignature string `json:"author_signature,omitempty"`
	// Optional. For text messages, the actual UTF-8 text of the message
	Text string `json:"text,omitempty"`
	// Optional. For text messages, special entities like usernames, URLs, bot commands, etc.
	// that appear in the text
	Entities []*MessageEntity `json:"entities,omitempty"`
	// Optional.
	// Options used for link preview generation for the message, if it is a text message and link preview options were changed
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	// Optional. Unique identifier of the message effect added to the message
	EffectId string `json:"effect_id,omitempty"`
	// Optional. Message is an animation, information about the animation.
	// For backward compatibility, when this field is set, the document field will also be set
	Animation *TelegramAnimation `json:"animation,omitempty"`
	// Optional. Message is an audio file, information about the file
	Audio *TelegramAudio `json:"audio,omitempty"`
	// Optional. Message is a general file, information about the file
	Document *TelegramDocument `json:"document,omitempty"`
	// Optional. Message contains paid media; information about the paid media
	PaidMedia *PaidMediaInfo `json:"paid_media,omitempty"`
	// Optional. Message is a photo, available sizes of the photo
	Photo TelegramPhoto `json:"photo,omitempty"`
	// Optional. Message is a sticker, information about the sticker
	Sticker *Sticker `json:"sticker,omitempty"`
	// Optional. Message is a forwarded story
	Story *Story `json:"story,omitempty"`
	// Optional. Message is a video, information about the video
	Video *TelegramVideo `json:"video,omitempty"`
	// Optional. Message is a video note, information about the video message
	VideoNote *VideoNote `json:"video_note,omitempty"`
	// Optional. Message is a voice message, information about the file
	Voice *Voice `json:"voice,omitempty"`
	// Optional. Caption for the animation, audio, document, paid media, photo, video or voice
	Caption string `json:"caption,omitempty"`
	// Optional. For messages with a caption, special entities like usernames, URLs, bot commands, etc.
	// that appear in the caption
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. True, if the message media is covered by a spoiler animation
	HasMediaSpoiler bool `json:"has_media_spoiler,omitempty"`
	// Optional. Message is a shared contact, information about the contact
	Contact *Contact `json:"contact,omitempty"`
	// Optional. Message is a dice with random value
	Dice *Dice `json:"dice,omitempty"`
	// Optional. Message is a game, information about the game. More about games: https://core.telegram.org/bots/api#games
	Game *Game `json:"game,omitempty"`
	// Optional. Message is a native poll, information about the poll
	Poll *Poll `json:"poll,omitempty"`
	// Optional. Message is a venue, information about the venue.
	// For backward compatibility, when this field is set, the location field will also be set
	Venue *Venue `json:"venue,omitempty"`
	// Optional. Message is a shared location, information about the location
	Location *Location `json:"location,omitempty"`
	// Optional.
	// New members that were added to the group or supergroup and information about them (the bot itself may be one of these members)
	NewChatMembers []*User `json:"new_chat_members,omitempty"`
	// Optional. A member was removed from the group, information about them (this member may be the bot itself)
	LeftChatMember *User `json:"left_chat_member,omitempty"`
	// Optional. A chat title was changed to this value
	NewChatTitle string `json:"new_chat_title,omitempty"`
	// Optional. A chat photo was change to this value
	NewChatPhoto TelegramPhoto `json:"new_chat_photo,omitempty"`
	// Optional. Service message: the chat photo was deleted
	DeleteChatPhoto bool `json:"delete_chat_photo,omitempty"`
	// Optional. Service message: the group has been created
	GroupChatCreated bool `json:"group_chat_created,omitempty"`
	// Optional. Service message: the supergroup has been created.
	// This field can't be received in a message coming through updates, because bot can't be a member of a supergroup when it is created.
	// It can only be found in reply_to_message if someone replies to a very first message in a directly created supergroup.
	SupergroupChatCreated bool `json:"supergroup_chat_created,omitempty"`
	// Optional. Service message: the channel has been created.
	// This field can't be received in a message coming through updates, because bot can't be a member of a channel when it is created.
	// It can only be found in reply_to_message if someone replies to a very first message in a channel.
	ChannelChatCreated bool `json:"channel_chat_created,omitempty"`
	// Optional. Service message: auto-delete timer settings changed in the chat
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	// Optional. The group has been migrated to a supergroup with the specified identifier.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	MigrateToChatId int64 `json:"migrate_to_chat_id,omitempty"`
	// Optional. The supergroup has been migrated from a group with the specified identifier.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	MigrateFromChatId int64 `json:"migrate_from_chat_id,omitempty"`
	// Optional. Specified message was pinned.
	// Note that the Message object in this field will not contain further reply_to_message fields even if it itself is a reply.
	PinnedMessage *Message `json:"pinned_message,omitempty"`
	// Optional. Message is an invoice for a payment, information about the invoice.
	// More about payments: https://core.telegram.org/bots/api#payments
	Invoice *Invoice `json:"invoice,omitempty"`
	// Optional. Message is a service message about a successful payment, information about the payment.
	// More about payments: https://core.telegram.org/bots/api#payments
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"`
	// Optional. Message is a service message about a refunded payment, information about the payment.
	// More about payments: https://core.telegram.org/bots/api#payments
	RefundedPayment *RefundedPayment `json:"refunded_payment,omitempty"`
	// Optional. Service message: users were shared with the bot
	UsersShared *UsersShared `json:"users_shared,omitempty"`
	// Optional. Service message: a chat was shared with the bot
	ChatShared *ChatShared `json:"chat_shared,omitempty"`
	// Optional. The domain name of the website on which the user has logged in.
	// More about Telegram Login: https://core.telegram.org/widgets/login
	ConnectedWebsite string `json:"connected_website,omitempty"`
	// Optional.
	// Service message: the user allowed the bot to write messages after adding it to the attachment or side menu, launching a Web App from a link, or accepting an explicit request from a Web App sent by the method requestWriteAccess
	WriteAccessAllowed *WriteAccessAllowed `json:"write_access_allowed,omitempty"`
	// Optional. Telegram Passport data
	PassportData *PassportData `json:"passport_data,omitempty"`
	// Optional. Service message. A user in the chat triggered another user's proximity alert while sharing Live Location.
	ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered,omitempty"`
	// Optional. Service message: user boosted the chat
	BoostAdded *ChatBoostAdded `json:"boost_added,omitempty"`
	// Optional. Service message: chat background set
	ChatBackgroundSet *ChatBackground `json:"chat_background_set,omitempty"`
	// Optional. Service message: forum topic created
	ForumTopicCreated *ForumTopicCreated `json:"forum_topic_created,omitempty"`
	// Optional. Service message: forum topic edited
	ForumTopicEdited *ForumTopicEdited `json:"forum_topic_edited,omitempty"`
	// Optional. Service message: forum topic closed
	ForumTopicClosed *ForumTopicClosed `json:"forum_topic_closed,omitempty"`
	// Optional. Service message: forum topic reopened
	ForumTopicReopened *ForumTopicReopened `json:"forum_topic_reopened,omitempty"`
	// Optional. Service message: the 'General' forum topic hidden
	GeneralForumTopicHidden *GeneralForumTopicHidden `json:"general_forum_topic_hidden,omitempty"`
	// Optional. Service message: the 'General' forum topic unhidden
	GeneralForumTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden,omitempty"`
	// Optional. Service message: a scheduled giveaway was created
	GiveawayCreated *GiveawayCreated `json:"giveaway_created,omitempty"`
	// Optional. The message is a scheduled giveaway message
	Giveaway *Giveaway `json:"giveaway,omitempty"`
	// Optional. A giveaway with public winners was completed
	GiveawayWinners *GiveawayWinners `json:"giveaway_winners,omitempty"`
	// Optional. Service message: a giveaway without public winners was completed
	GiveawayCompleted *GiveawayCompleted `json:"giveaway_completed,omitempty"`
	// Optional. Service message: video chat scheduled
	VideoChatScheduled *VideoChatScheduled `json:"video_chat_scheduled,omitempty"`
	// Optional. Service message: video chat started
	VideoChatStarted *VideoChatStarted `json:"video_chat_started,omitempty"`
	// Optional. Service message: video chat ended
	VideoChatEnded *VideoChatEnded `json:"video_chat_ended,omitempty"`
	// Optional. Service message: new participants invited to a video chat
	VideoChatParticipantsInvited *VideoChatParticipantsInvited `json:"video_chat_participants_invited,omitempty"`
	// Optional. Service message: data sent by a Web App
	WebAppData *WebAppData `json:"web_app_data,omitempty"`
	// Optional. Inline keyboard attached to the message. login_url buttons are represented as ordinary url buttons.
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// MessageAutoDeleteTimerChanged This object represents a service message about a change in auto-delete timer settings.
type MessageAutoDeleteTimerChanged struct {
	// New auto-delete time for messages in the chat; in seconds
	MessageAutoDeleteTime int64 `json:"message_auto_delete_time"`
}

// MessageEntity This object represents one special entity in a text message. For example, hashtags, usernames, URLs, etc.
type MessageEntity struct {
	// Type of the entity.
	// Currently, can be "mention" (@username), "hashtag" (#hashtag or #hashtag@chatusername), "cashtag" ($USD or $USD@chatusername), "bot_command" (/start@jobs_bot), "url" (https://telegram.org), "email" (do-not-reply@telegram.org), "phone_number" (+1-212-555-0123), "bold" (bold text), "italic" (italic text), "underline" (underlined text), "strikethrough" (strikethrough text), "spoiler" (spoiler message), "blockquote" (block quotation), "expandable_blockquote" (collapsed-by-default block quotation), "code" (monowidth string), "pre" (monowidth block), "text_link" (for clickable text URLs), "text_mention" (for users without usernames), "custom_emoji" (for inline custom emoji stickers)
	Type string `json:"type"`
	// Offset in UTF-16 code units to the start of the entity
	Offset int64 `json:"offset"`
	// Length of the entity in UTF-16 code units
	Length int64 `json:"length"`
	// Optional. For "text_link" only, URL that will be opened after user taps on the text
	Url string `json:"url,omitempty"`
	// Optional. For "text_mention" only, the mentioned user
	User *User `json:"user,omitempty"`
	// Optional. For "pre" only, the programming language of the entity text
	Language string `json:"language,omitempty"`
	// Optional. For "custom_emoji" only, unique identifier of the custom emoji.
	// Use getCustomEmojiStickers to get full information about the sticker
	CustomEmojiId string `json:"custom_emoji_id,omitempty"`
}

// MessageId This object represents a unique message identifier.
type MessageId struct {
	// Unique message identifier.
	// In specific instances (e.g., message containing a video sent to a big chat), the server might automatically schedule a message instead of sending it immediately.
	// In such cases, this field will be 0 and the relevant message will be unusable until it is actually sent
	MessageId int64 `json:"message_id"`
}

// MessageOrigin This object describes the origin of a message. It can be one of
// - MessageOriginUser
// - MessageOriginHiddenUser
// - MessageOriginChat
// - MessageOriginChannel
type MessageOrigin interface {
	OptUser() *MessageOriginUser
	OptHiddenUser() *MessageOriginHiddenUser
	OptChat() *MessageOriginChat
	OptChannel() *MessageOriginChannel
}

var (
	_ MessageOrigin = &MessageOriginUser{}
	_ MessageOrigin = &MessageOriginHiddenUser{}
	_ MessageOrigin = &MessageOriginChat{}
	_ MessageOrigin = &MessageOriginChannel{}
)

func (impl *MessageOriginUser) OptUser() *MessageOriginUser             { return impl }
func (impl *MessageOriginUser) OptHiddenUser() *MessageOriginHiddenUser { return nil }
func (impl *MessageOriginUser) OptChat() *MessageOriginChat             { return nil }
func (impl *MessageOriginUser) OptChannel() *MessageOriginChannel       { return nil }

func (impl *MessageOriginHiddenUser) OptUser() *MessageOriginUser             { return nil }
func (impl *MessageOriginHiddenUser) OptHiddenUser() *MessageOriginHiddenUser { return impl }
func (impl *MessageOriginHiddenUser) OptChat() *MessageOriginChat             { return nil }
func (impl *MessageOriginHiddenUser) OptChannel() *MessageOriginChannel       { return nil }

func (impl *MessageOriginChat) OptUser() *MessageOriginUser             { return nil }
func (impl *MessageOriginChat) OptHiddenUser() *MessageOriginHiddenUser { return nil }
func (impl *MessageOriginChat) OptChat() *MessageOriginChat             { return impl }
func (impl *MessageOriginChat) OptChannel() *MessageOriginChannel       { return nil }

func (impl *MessageOriginChannel) OptUser() *MessageOriginUser             { return nil }
func (impl *MessageOriginChannel) OptHiddenUser() *MessageOriginHiddenUser { return nil }
func (impl *MessageOriginChannel) OptChat() *MessageOriginChat             { return nil }
func (impl *MessageOriginChannel) OptChannel() *MessageOriginChannel       { return impl }

// MessageOriginChannel The message was originally sent to a channel chat.
type MessageOriginChannel struct {
	// Type of the message origin, always "channel"
	Type string `json:"type" default:"channel"`
	// Date the message was sent originally in Unix time
	Date int64 `json:"date"`
	// Channel chat to which the message was originally sent
	Chat *Chat `json:"chat"`
	// Unique message identifier inside the chat
	MessageId int64 `json:"message_id"`
	// Optional. Signature of the original post author
	AuthorSignature string `json:"author_signature,omitempty"`
}

// MessageOriginChat The message was originally sent on behalf of a chat to a group chat.
type MessageOriginChat struct {
	// Type of the message origin, always "chat"
	Type string `json:"type" default:"chat"`
	// Date the message was sent originally in Unix time
	Date int64 `json:"date"`
	// Chat that sent the message originally
	SenderChat *Chat `json:"sender_chat"`
	// Optional. For messages originally sent by an anonymous chat administrator, original message author signature
	AuthorSignature string `json:"author_signature,omitempty"`
}

// MessageOriginHiddenUser The message was originally sent by an unknown user.
type MessageOriginHiddenUser struct {
	// Type of the message origin, always "hidden_user"
	Type string `json:"type" default:"hidden_user"`
	// Date the message was sent originally in Unix time
	Date int64 `json:"date"`
	// Name of the user that sent the message originally
	SenderUserName string `json:"sender_user_name"`
}

// MessageOriginUser The message was originally sent by a known user.
type MessageOriginUser struct {
	// Type of the message origin, always "user"
	Type string `json:"type" default:"user"`
	// Date the message was sent originally in Unix time
	Date int64 `json:"date"`
	// User that sent the message originally
	SenderUser *User `json:"sender_user"`
}

// MessageReactionCountUpdated This object represents reaction changes on a message with anonymous reactions.
type MessageReactionCountUpdated struct {
	// The chat containing the message
	Chat *Chat `json:"chat"`
	// Unique message identifier inside the chat
	MessageId int64 `json:"message_id"`
	// Date of the change in Unix time
	Date int64 `json:"date"`
	// List of reactions that are present on the message
	Reactions []*ReactionCount `json:"reactions"`
}

// MessageReactionUpdated This object represents a change of a reaction on a message performed by a user.
type MessageReactionUpdated struct {
	// The chat containing the message the user reacted to
	Chat *Chat `json:"chat"`
	// Unique identifier of the message inside the chat
	MessageId int64 `json:"message_id"`
	// Optional. The user that changed the reaction, if the user isn't anonymous
	User *User `json:"user,omitempty"`
	// Optional. The chat on behalf of which the reaction was changed, if the user is anonymous
	ActorChat *Chat `json:"actor_chat,omitempty"`
	// Date of the change in Unix time
	Date int64 `json:"date"`
	// Previous list of reaction types that were set by the user
	OldReaction []ReactionType `json:"old_reaction"`
	// New list of reaction types that have been set by the user
	NewReaction []ReactionType `json:"new_reaction"`
}

// OrderInfo This object represents information about an order.
type OrderInfo struct {
	// Optional. User name
	Name string `json:"name,omitempty"`
	// Optional. User's phone number
	PhoneNumber string `json:"phone_number,omitempty"`
	// Optional. User email
	Email string `json:"email,omitempty"`
	// Optional. User shipping address
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// PaidMedia This object describes paid media. Currently, it can be one of
// - PaidMediaPreview
// - PaidMediaPhoto
// - PaidMediaVideo
type PaidMedia interface {
	OptPreview() *PaidMediaPreview
	OptPhoto() *PaidMediaPhoto
	OptVideo() *PaidMediaVideo
}

var (
	_ PaidMedia = &PaidMediaPreview{}
	_ PaidMedia = &PaidMediaPhoto{}
	_ PaidMedia = &PaidMediaVideo{}
)

func (impl *PaidMediaPreview) OptPreview() *PaidMediaPreview { return impl }
func (impl *PaidMediaPreview) OptPhoto() *PaidMediaPhoto     { return nil }
func (impl *PaidMediaPreview) OptVideo() *PaidMediaVideo     { return nil }

func (impl *PaidMediaPhoto) OptPreview() *PaidMediaPreview { return nil }
func (impl *PaidMediaPhoto) OptPhoto() *PaidMediaPhoto     { return impl }
func (impl *PaidMediaPhoto) OptVideo() *PaidMediaVideo     { return nil }

func (impl *PaidMediaVideo) OptPreview() *PaidMediaPreview { return nil }
func (impl *PaidMediaVideo) OptPhoto() *PaidMediaPhoto     { return nil }
func (impl *PaidMediaVideo) OptVideo() *PaidMediaVideo     { return impl }

// PaidMediaInfo Describes the paid media added to a message.
type PaidMediaInfo struct {
	// The number of Telegram Stars that must be paid to buy access to the media
	StarCount int64 `json:"star_count"`
	// Information about the paid media
	PaidMedia []PaidMedia `json:"paid_media"`
}

// PaidMediaPhoto The paid media is a photo.
type PaidMediaPhoto struct {
	// Type of the paid media, always "photo"
	Type string `json:"type"`
	// The photo
	Photo TelegramPhoto `json:"photo"`
}

// PaidMediaPreview The paid media isn't available before the payment.
type PaidMediaPreview struct {
	// Type of the paid media, always "preview"
	Type string `json:"type"`
	// Optional. Media width as defined by the sender
	Width int64 `json:"width,omitempty"`
	// Optional. Media height as defined by the sender
	Height int64 `json:"height,omitempty"`
	// Optional. Duration of the media in seconds as defined by the sender
	Duration int64 `json:"duration,omitempty"`
}

// PaidMediaPurchased This object contains information about a paid media purchase.
type PaidMediaPurchased struct {
	// User who purchased the media
	From *User `json:"from"`
	// Bot-specified paid media payload
	PaidMediaPayload string `json:"paid_media_payload"`
}

// PaidMediaVideo The paid media is a video.
type PaidMediaVideo struct {
	// Type of the paid media, always "video"
	Type string `json:"type"`
	// The video
	Video *TelegramVideo `json:"video"`
}

// PassportData Describes Telegram Passport data shared with the bot by the user.
type PassportData struct {
	// Array with information about documents and other Telegram Passport elements that was shared with the bot
	Data []*EncryptedPassportElement `json:"data"`
	// Encrypted credentials required to decrypt the data
	Credentials *EncryptedCredentials `json:"credentials"`
}

// PassportElementError This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user. It should be one of:
// - PassportElementErrorDataField
// - PassportElementErrorFrontSide
// - PassportElementErrorReverseSide
// - PassportElementErrorSelfie
// - PassportElementErrorFile
// - PassportElementErrorFiles
// - PassportElementErrorTranslationFile
// - PassportElementErrorTranslationFiles
// - PassportElementErrorUnspecified
type PassportElementError interface {
	OptDataField() *PassportElementErrorDataField
	OptFrontSide() *PassportElementErrorFrontSide
	OptReverseSide() *PassportElementErrorReverseSide
	OptSelfie() *PassportElementErrorSelfie
	OptFile() *PassportElementErrorFile
	OptFiles() *PassportElementErrorFiles
	OptTranslationFile() *PassportElementErrorTranslationFile
	OptTranslationFiles() *PassportElementErrorTranslationFiles
	OptUnspecified() *PassportElementErrorUnspecified
}

var (
	_ PassportElementError = &PassportElementErrorDataField{}
	_ PassportElementError = &PassportElementErrorFrontSide{}
	_ PassportElementError = &PassportElementErrorReverseSide{}
	_ PassportElementError = &PassportElementErrorSelfie{}
	_ PassportElementError = &PassportElementErrorFile{}
	_ PassportElementError = &PassportElementErrorFiles{}
	_ PassportElementError = &PassportElementErrorTranslationFile{}
	_ PassportElementError = &PassportElementErrorTranslationFiles{}
	_ PassportElementError = &PassportElementErrorUnspecified{}
)

func (impl *PassportElementErrorDataField) OptDataField() *PassportElementErrorDataField { return impl }
func (impl *PassportElementErrorDataField) OptFrontSide() *PassportElementErrorFrontSide { return nil }
func (impl *PassportElementErrorDataField) OptReverseSide() *PassportElementErrorReverseSide {
	return nil
}
func (impl *PassportElementErrorDataField) OptSelfie() *PassportElementErrorSelfie { return nil }
func (impl *PassportElementErrorDataField) OptFile() *PassportElementErrorFile     { return nil }
func (impl *PassportElementErrorDataField) OptFiles() *PassportElementErrorFiles   { return nil }
func (impl *PassportElementErrorDataField) OptTranslationFile() *PassportElementErrorTranslationFile {
	return nil
}
func (impl *PassportElementErrorDataField) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return nil
}
func (impl *PassportElementErrorDataField) OptUnspecified() *PassportElementErrorUnspecified {
	return nil
}

func (impl *PassportElementErrorFrontSide) OptDataField() *PassportElementErrorDataField { return nil }
func (impl *PassportElementErrorFrontSide) OptFrontSide() *PassportElementErrorFrontSide { return impl }
func (impl *PassportElementErrorFrontSide) OptReverseSide() *PassportElementErrorReverseSide {
	return nil
}
func (impl *PassportElementErrorFrontSide) OptSelfie() *PassportElementErrorSelfie { return nil }
func (impl *PassportElementErrorFrontSide) OptFile() *PassportElementErrorFile     { return nil }
func (impl *PassportElementErrorFrontSide) OptFiles() *PassportElementErrorFiles   { return nil }
func (impl *PassportElementErrorFrontSide) OptTranslationFile() *PassportElementErrorTranslationFile {
	return nil
}
func (impl *PassportElementErrorFrontSide) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return nil
}
func (impl *PassportElementErrorFrontSide) OptUnspecified() *PassportElementErrorUnspecified {
	return nil
}

func (impl *PassportElementErrorReverseSide) OptDataField() *PassportElementErrorDataField {
	return nil
}
func (impl *PassportElementErrorReverseSide) OptFrontSide() *PassportElementErrorFrontSide {
	return nil
}
func (impl *PassportElementErrorReverseSide) OptReverseSide() *PassportElementErrorReverseSide {
	return impl
}
func (impl *PassportElementErrorReverseSide) OptSelfie() *PassportElementErrorSelfie { return nil }
func (impl *PassportElementErrorReverseSide) OptFile() *PassportElementErrorFile     { return nil }
func (impl *PassportElementErrorReverseSide) OptFiles() *PassportElementErrorFiles   { return nil }
func (impl *PassportElementErrorReverseSide) OptTranslationFile() *PassportElementErrorTranslationFile {
	return nil
}
func (impl *PassportElementErrorReverseSide) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return nil
}
func (impl *PassportElementErrorReverseSide) OptUnspecified() *PassportElementErrorUnspecified {
	return nil
}

func (impl *PassportElementErrorSelfie) OptDataField() *PassportElementErrorDataField     { return nil }
func (impl *PassportElementErrorSelfie) OptFrontSide() *PassportElementErrorFrontSide     { return nil }
func (impl *PassportElementErrorSelfie) OptReverseSide() *PassportElementErrorReverseSide { return nil }
func (impl *PassportElementErrorSelfie) OptSelfie() *PassportElementErrorSelfie           { return impl }
func (impl *PassportElementErrorSelfie) OptFile() *PassportElementErrorFile               { return nil }
func (impl *PassportElementErrorSelfie) OptFiles() *PassportElementErrorFiles             { return nil }
func (impl *PassportElementErrorSelfie) OptTranslationFile() *PassportElementErrorTranslationFile {
	return nil
}
func (impl *PassportElementErrorSelfie) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return nil
}
func (impl *PassportElementErrorSelfie) OptUnspecified() *PassportElementErrorUnspecified { return nil }

func (impl *PassportElementErrorFile) OptDataField() *PassportElementErrorDataField     { return nil }
func (impl *PassportElementErrorFile) OptFrontSide() *PassportElementErrorFrontSide     { return nil }
func (impl *PassportElementErrorFile) OptReverseSide() *PassportElementErrorReverseSide { return nil }
func (impl *PassportElementErrorFile) OptSelfie() *PassportElementErrorSelfie           { return nil }
func (impl *PassportElementErrorFile) OptFile() *PassportElementErrorFile               { return impl }
func (impl *PassportElementErrorFile) OptFiles() *PassportElementErrorFiles             { return nil }
func (impl *PassportElementErrorFile) OptTranslationFile() *PassportElementErrorTranslationFile {
	return nil
}
func (impl *PassportElementErrorFile) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return nil
}
func (impl *PassportElementErrorFile) OptUnspecified() *PassportElementErrorUnspecified { return nil }

func (impl *PassportElementErrorFiles) OptDataField() *PassportElementErrorDataField     { return nil }
func (impl *PassportElementErrorFiles) OptFrontSide() *PassportElementErrorFrontSide     { return nil }
func (impl *PassportElementErrorFiles) OptReverseSide() *PassportElementErrorReverseSide { return nil }
func (impl *PassportElementErrorFiles) OptSelfie() *PassportElementErrorSelfie           { return nil }
func (impl *PassportElementErrorFiles) OptFile() *PassportElementErrorFile               { return nil }
func (impl *PassportElementErrorFiles) OptFiles() *PassportElementErrorFiles             { return impl }
func (impl *PassportElementErrorFiles) OptTranslationFile() *PassportElementErrorTranslationFile {
	return nil
}
func (impl *PassportElementErrorFiles) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return nil
}
func (impl *PassportElementErrorFiles) OptUnspecified() *PassportElementErrorUnspecified { return nil }

func (impl *PassportElementErrorTranslationFile) OptDataField() *PassportElementErrorDataField {
	return nil
}
func (impl *PassportElementErrorTranslationFile) OptFrontSide() *PassportElementErrorFrontSide {
	return nil
}
func (impl *PassportElementErrorTranslationFile) OptReverseSide() *PassportElementErrorReverseSide {
	return nil
}
func (impl *PassportElementErrorTranslationFile) OptSelfie() *PassportElementErrorSelfie { return nil }
func (impl *PassportElementErrorTranslationFile) OptFile() *PassportElementErrorFile     { return nil }
func (impl *PassportElementErrorTranslationFile) OptFiles() *PassportElementErrorFiles   { return nil }
func (impl *PassportElementErrorTranslationFile) OptTranslationFile() *PassportElementErrorTranslationFile {
	return impl
}
func (impl *PassportElementErrorTranslationFile) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return nil
}
func (impl *PassportElementErrorTranslationFile) OptUnspecified() *PassportElementErrorUnspecified {
	return nil
}

func (impl *PassportElementErrorTranslationFiles) OptDataField() *PassportElementErrorDataField {
	return nil
}
func (impl *PassportElementErrorTranslationFiles) OptFrontSide() *PassportElementErrorFrontSide {
	return nil
}
func (impl *PassportElementErrorTranslationFiles) OptReverseSide() *PassportElementErrorReverseSide {
	return nil
}
func (impl *PassportElementErrorTranslationFiles) OptSelfie() *PassportElementErrorSelfie { return nil }
func (impl *PassportElementErrorTranslationFiles) OptFile() *PassportElementErrorFile     { return nil }
func (impl *PassportElementErrorTranslationFiles) OptFiles() *PassportElementErrorFiles   { return nil }
func (impl *PassportElementErrorTranslationFiles) OptTranslationFile() *PassportElementErrorTranslationFile {
	return nil
}
func (impl *PassportElementErrorTranslationFiles) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return impl
}
func (impl *PassportElementErrorTranslationFiles) OptUnspecified() *PassportElementErrorUnspecified {
	return nil
}

func (impl *PassportElementErrorUnspecified) OptDataField() *PassportElementErrorDataField {
	return nil
}
func (impl *PassportElementErrorUnspecified) OptFrontSide() *PassportElementErrorFrontSide {
	return nil
}
func (impl *PassportElementErrorUnspecified) OptReverseSide() *PassportElementErrorReverseSide {
	return nil
}
func (impl *PassportElementErrorUnspecified) OptSelfie() *PassportElementErrorSelfie { return nil }
func (impl *PassportElementErrorUnspecified) OptFile() *PassportElementErrorFile     { return nil }
func (impl *PassportElementErrorUnspecified) OptFiles() *PassportElementErrorFiles   { return nil }
func (impl *PassportElementErrorUnspecified) OptTranslationFile() *PassportElementErrorTranslationFile {
	return nil
}
func (impl *PassportElementErrorUnspecified) OptTranslationFiles() *PassportElementErrorTranslationFiles {
	return nil
}
func (impl *PassportElementErrorUnspecified) OptUnspecified() *PassportElementErrorUnspecified {
	return impl
}

// PassportElementErrorDataField Represents an issue in one of the data fields that was provided by the user.
// The error is considered resolved when the field's value changes.
type PassportElementErrorDataField struct {
	// Error source, must be data
	Source string `json:"source" default:"data"`
	// The section of the user's Telegram Passport which has the error, one of "personal_details", "passport", "driver_license", "identity_card", "internal_passport", "address"
	Type string `json:"type"`
	// Name of the data field which has the error
	FieldName string `json:"field_name"`
	// Base64-encoded data hash
	DataHash string `json:"data_hash"`
	// Error message
	Message string `json:"message"`
}

// PassportElementErrorFile Represents an issue with a document scan.
// The error is considered resolved when the file with the document scan changes.
type PassportElementErrorFile struct {
	// Error source, must be file
	Source string `json:"source" default:"file"`
	// The section of the user's Telegram Passport which has the issue, one of "utility_bill", "bank_statement", "rental_agreement", "passport_registration", "temporary_registration"
	Type string `json:"type"`
	// Base64-encoded file hash
	FileHash string `json:"file_hash"`
	// Error message
	Message string `json:"message"`
}

// PassportElementErrorFiles Represents an issue with a list of scans.
// The error is considered resolved when the list of files containing the scans changes.
type PassportElementErrorFiles struct {
	// Error source, must be files
	Source string `json:"source" default:"files"`
	// The section of the user's Telegram Passport which has the issue, one of "utility_bill", "bank_statement", "rental_agreement", "passport_registration", "temporary_registration"
	Type string `json:"type"`
	// List of base64-encoded file hashes
	FileHashes []string `json:"file_hashes"`
	// Error message
	Message string `json:"message"`
}

// PassportElementErrorFrontSide Represents an issue with the front side of a document.
// The error is considered resolved when the file with the front side of the document changes.
type PassportElementErrorFrontSide struct {
	// Error source, must be front_side
	Source string `json:"source" default:"front_side"`
	// The section of the user's Telegram Passport which has the issue, one of "passport", "driver_license", "identity_card", "internal_passport"
	Type string `json:"type"`
	// Base64-encoded hash of the file with the front side of the document
	FileHash string `json:"file_hash"`
	// Error message
	Message string `json:"message"`
}

// PassportElementErrorReverseSide Represents an issue with the reverse side of a document.
// The error is considered resolved when the file with reverse side of the document changes.
type PassportElementErrorReverseSide struct {
	// Error source, must be reverse_side
	Source string `json:"source" default:"reverse_side"`
	// The section of the user's Telegram Passport which has the issue, one of "driver_license", "identity_card"
	Type string `json:"type"`
	// Base64-encoded hash of the file with the reverse side of the document
	FileHash string `json:"file_hash"`
	// Error message
	Message string `json:"message"`
}

// PassportElementErrorSelfie Represents an issue with the selfie with a document.
// The error is considered resolved when the file with the selfie changes.
type PassportElementErrorSelfie struct {
	// Error source, must be selfie
	Source string `json:"source" default:"selfie"`
	// The section of the user's Telegram Passport which has the issue, one of "passport", "driver_license", "identity_card", "internal_passport"
	Type string `json:"type"`
	// Base64-encoded hash of the file with the selfie
	FileHash string `json:"file_hash"`
	// Error message
	Message string `json:"message"`
}

// PassportElementErrorTranslationFile Represents an issue with one of the files that constitute the translation of a document.
// The error is considered resolved when the file changes.
type PassportElementErrorTranslationFile struct {
	// Error source, must be translation_file
	Source string `json:"source" default:"translation_file"`
	// Type of element of the user's Telegram Passport which has the issue, one of "passport", "driver_license", "identity_card", "internal_passport", "utility_bill", "bank_statement", "rental_agreement", "passport_registration", "temporary_registration"
	Type string `json:"type"`
	// Base64-encoded file hash
	FileHash string `json:"file_hash"`
	// Error message
	Message string `json:"message"`
}

// PassportElementErrorTranslationFiles Represents an issue with the translated version of a document.
// The error is considered resolved when a file with the document translation change.
type PassportElementErrorTranslationFiles struct {
	// Error source, must be translation_files
	Source string `json:"source" default:"translation_files"`
	// Type of element of the user's Telegram Passport which has the issue, one of "passport", "driver_license", "identity_card", "internal_passport", "utility_bill", "bank_statement", "rental_agreement", "passport_registration", "temporary_registration"
	Type string `json:"type"`
	// List of base64-encoded file hashes
	FileHashes []string `json:"file_hashes"`
	// Error message
	Message string `json:"message"`
}

// PassportElementErrorUnspecified Represents an issue in an unspecified place. The error is considered resolved when new data is added.
type PassportElementErrorUnspecified struct {
	// Error source, must be unspecified
	Source string `json:"source" default:"unspecified"`
	// Type of element of the user's Telegram Passport which has the issue
	Type string `json:"type"`
	// Base64-encoded element hash
	ElementHash string `json:"element_hash"`
	// Error message
	Message string `json:"message"`
}

// PassportFile This object represents a file uploaded to Telegram Passport.
// Currently all Telegram Passport files are in JPEG format when decrypted and don't exceed 10MB.
type PassportFile struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// File size in bytes
	FileSize int64 `json:"file_size"`
	// Unix time when the file was uploaded
	FileDate int64 `json:"file_date"`
}

func (impl *PassportFile) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *PassportFile) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// Photo Represents a photo to be sent.
type Photo struct {
	// Type of the result, must be photo
	Type string `json:"type" default:"photo"`
	// File to send. More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// Pass a file_id to send a file that exists on the Telegram servers (recommended), pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	Media InputFile `json:"media"`
	// Optional. Caption of the photo to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the photo caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Pass True if the photo needs to be covered with a spoiler animation
	HasSpoiler bool `json:"has_spoiler,omitempty"`
	// Used for uploading media.
	InputFile InputFile `json:"-"`
}

// PhotoSize This object represents one size of a photo or a file / sticker thumbnail.
type PhotoSize struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Photo width
	Width int64 `json:"width"`
	// Photo height
	Height int64 `json:"height"`
	// Optional. File size in bytes
	FileSize int64 `json:"file_size,omitempty"`
}

func (impl *PhotoSize) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *PhotoSize) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// Poll This object contains information about a poll.
type Poll struct {
	// Unique poll identifier
	Id string `json:"id"`
	// Poll question, 1-300 characters
	Question string `json:"question"`
	// Optional. Special entities that appear in the question.
	// Currently, only custom emoji entities are allowed in poll questions
	QuestionEntities []*MessageEntity `json:"question_entities,omitempty"`
	// List of poll options
	Options []*PollOption `json:"options"`
	// Total number of users that voted in the poll
	TotalVoterCount int64 `json:"total_voter_count"`
	// True, if the poll is closed
	IsClosed bool `json:"is_closed"`
	// True, if the poll is anonymous
	IsAnonymous bool `json:"is_anonymous"`
	// Poll type, currently can be "regular" or "quiz"
	Type string `json:"type"`
	// True, if the poll allows multiple answers
	AllowsMultipleAnswers bool `json:"allows_multiple_answers"`
	// Optional. 0-based identifier of the correct answer option.
	// Available only for polls in the quiz mode, which are closed, or was sent (not forwarded) by the bot or to the private chat with the bot.
	CorrectOptionId int64 `json:"correct_option_id,omitempty"`
	// Optional.
	// Text that is shown when a user chooses an incorrect answer or taps on the lamp icon in a quiz-style poll, 0-200 characters
	Explanation string `json:"explanation,omitempty"`
	// Optional. Special entities like usernames, URLs, bot commands, etc. that appear in the explanation
	ExplanationEntities []*MessageEntity `json:"explanation_entities,omitempty"`
	// Optional. Amount of time in seconds the poll will be active after creation
	OpenPeriod int64 `json:"open_period,omitempty"`
	// Optional. Point in time (Unix timestamp) when the poll will be automatically closed
	CloseDate int64 `json:"close_date,omitempty"`
}

// PollAnswer This object represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	// Unique poll identifier
	PollId string `json:"poll_id"`
	// Optional. The chat that changed the answer to the poll, if the voter is anonymous
	VoterChat *Chat `json:"voter_chat,omitempty"`
	// Optional. The user that changed the answer to the poll, if the voter isn't anonymous
	User *User `json:"user,omitempty"`
	// 0-based identifiers of chosen answer options. May be empty if the vote was retracted.
	OptionIds []int64 `json:"option_ids"`
}

// PollOption This object contains information about one answer option in a poll.
type PollOption struct {
	// Option text, 1-100 characters
	Text string `json:"text"`
	// Optional. Special entities that appear in the option text.
	// Currently, only custom emoji entities are allowed in poll option texts
	TextEntities []*MessageEntity `json:"text_entities,omitempty"`
	// Number of users that voted for this option
	VoterCount int64 `json:"voter_count"`
}

// PreCheckoutQuery This object contains information about an incoming pre-checkout query.
type PreCheckoutQuery struct {
	// Unique query identifier
	Id string `json:"id"`
	// User who sent the query
	From *User `json:"from"`
	// Three-letter ISO 4217 currency code, or "XTR" for payments in Telegram Stars
	Currency string `json:"currency"`
	// Total price in the smallest units of the currency (integer, not float/double).
	// For example, for a price of US$ 1.45 pass amount = 145.
	// See the exp parameter in currencies.json, it shows the number of digits past the decimal point for each currency (2 for the majority of currencies).
	TotalAmount int64 `json:"total_amount"`
	// Bot-specified invoice payload
	InvoicePayload string `json:"invoice_payload"`
	// Optional. Identifier of the shipping option chosen by the user
	ShippingOptionId string `json:"shipping_option_id,omitempty"`
	// Optional. Order information provided by the user
	OrderInfo *OrderInfo `json:"order_info,omitempty"`
}

// PreparedInlineMessage Describes an inline message to be sent by a user of a Mini App.
type PreparedInlineMessage struct {
	// Unique identifier of the prepared message
	Id string `json:"id"`
	// Expiration date of the prepared message, in Unix time. Expired prepared messages can no longer be used
	ExpirationDate int64 `json:"expiration_date"`
}

// ProximityAlertTriggered This object represents the content of a service message, sent whenever a user in the chat triggers a proximity alert set by another user.
type ProximityAlertTriggered struct {
	// User that triggered the alert
	Traveler *User `json:"traveler"`
	// User that set the alert
	Watcher *User `json:"watcher"`
	// The distance between the users
	Distance int64 `json:"distance"`
}

// ReactionCount Represents a reaction added to a message along with the number of times it was added.
type ReactionCount struct {
	// Type of the reaction
	Type ReactionType `json:"type"`
	// Number of times the reaction was added
	TotalCount int64 `json:"total_count"`
}

// ReactionType This object describes the type of a reaction. Currently, it can be one of
// - ReactionTypeEmoji
// - ReactionTypeCustomEmoji
// - ReactionTypePaid
type ReactionType interface {
	OptEmoji() *ReactionTypeEmoji
	OptCustomEmoji() *ReactionTypeCustomEmoji
	OptPaid() *ReactionTypePaid
}

var (
	_ ReactionType = &ReactionTypeEmoji{}
	_ ReactionType = &ReactionTypeCustomEmoji{}
	_ ReactionType = &ReactionTypePaid{}
)

func (impl *ReactionTypeEmoji) OptEmoji() *ReactionTypeEmoji             { return impl }
func (impl *ReactionTypeEmoji) OptCustomEmoji() *ReactionTypeCustomEmoji { return nil }
func (impl *ReactionTypeEmoji) OptPaid() *ReactionTypePaid               { return nil }

func (impl *ReactionTypeCustomEmoji) OptEmoji() *ReactionTypeEmoji             { return nil }
func (impl *ReactionTypeCustomEmoji) OptCustomEmoji() *ReactionTypeCustomEmoji { return impl }
func (impl *ReactionTypeCustomEmoji) OptPaid() *ReactionTypePaid               { return nil }

func (impl *ReactionTypePaid) OptEmoji() *ReactionTypeEmoji             { return nil }
func (impl *ReactionTypePaid) OptCustomEmoji() *ReactionTypeCustomEmoji { return nil }
func (impl *ReactionTypePaid) OptPaid() *ReactionTypePaid               { return impl }

// ReactionTypeCustomEmoji The reaction is based on a custom emoji.
type ReactionTypeCustomEmoji struct {
	// Type of the reaction, always "custom_emoji"
	Type string `json:"type" default:"custom_emoji"`
	// Custom emoji identifier
	CustomEmojiId string `json:"custom_emoji_id"`
}

// ReactionTypeEmoji The reaction is based on an emoji.
type ReactionTypeEmoji struct {
	// Type of the reaction, always "emoji"
	Type string `json:"type" default:"emoji"`
	// Reaction emoji.
	// Currently, it can be one of "👍", "👎", "❤", "🔥", "🥰", "👏", "😁", "🤔", "🤯", "😱", "🤬", "😢", "🎉", "🤩", "🤮", "💩", "🙏", "👌", "🕊", "🤡", "🥱", "🥴", "😍", "🐳", "❤‍🔥", "🌚", "🌭", "💯", "🤣", "⚡", "🍌", "🏆", "💔", "🤨", "😐", "🍓", "🍾", "💋", "🖕", "😈", "😴", "😭", "🤓", "👻", "👨‍💻", "👀", "🎃", "🙈", "😇", "😨", "🤝", "✍", "🤗", "🫡", "🎅", "🎄", "☃", "💅", "🤪", "🗿", "🆒", "💘", "🙉", "🦄", "😘", "💊", "🙊", "😎", "👾", "🤷‍♂", "🤷", "🤷‍♀", "😡"
	Emoji string `json:"emoji"`
}

// ReactionTypePaid The reaction is paid.
type ReactionTypePaid struct {
	// Type of the reaction, always "paid"
	Type string `json:"type" default:"paid"`
}

// RefundedPayment This object contains basic information about a refunded payment.
type RefundedPayment struct {
	// Three-letter ISO 4217 currency code, or "XTR" for payments in Telegram Stars. Currently, always "XTR"
	Currency string `json:"currency"`
	// Total refunded price in the smallest units of the currency (integer, not float/double).
	// For example, for a price of US$ 1.45, total_amount = 145.
	// See the exp parameter in currencies.json, it shows the number of digits past the decimal point for each currency (2 for the majority of currencies).
	TotalAmount int64 `json:"total_amount"`
	// Bot-specified invoice payload
	InvoicePayload string `json:"invoice_payload"`
	// Telegram payment identifier
	TelegramPaymentChargeId string `json:"telegram_payment_charge_id"`
	// Optional. Provider payment identifier
	ProviderPaymentChargeId string `json:"provider_payment_charge_id,omitempty"`
}

// ReplyKeyboardMarkup This object represents a custom keyboard with reply options (see Introduction to bots for details and examples).
// Not supported in channels and for messages sent on behalf of a Telegram Business account.
type ReplyKeyboardMarkup struct {
	// Array of button rows, each represented by an Array of KeyboardButton objects
	Keyboard [][]*KeyboardButton `json:"keyboard"`
	// Optional. Requests clients to always show the keyboard when the regular keyboard is hidden.
	// Defaults to false, in which case the custom keyboard can be hidden and opened with a keyboard icon.
	IsPersistent bool `json:"is_persistent,omitempty"`
	// Optional.
	// Requests clients to resize the keyboard vertically for optimal fit (e.g., make the keyboard smaller if there are just two rows of buttons).
	// Defaults to false, in which case the custom keyboard is always of the same height as the app's standard keyboard.
	ResizeKeyboard bool `json:"resize_keyboard,omitempty"`
	// Optional. Requests clients to hide the keyboard as soon as it's been used. Defaults to false.
	// The keyboard will still be available, but clients will automatically display the usual letter-keyboard in the chat - the user can press a special button in the input field to see the custom keyboard again.
	OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`
	// Optional. The placeholder to be shown in the input field when the keyboard is active; 1-64 characters
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	// Optional. Use this parameter if you want to show the keyboard to specific users only.
	// Targets: 1) users that are @mentioned in the text of the Message object; 2) if the bot's message is a reply to a message in the same chat and forum topic, sender of the original message.
	// Example: A user requests to change the bot's language, bot replies to the request with a keyboard to select the new language.
	// Other users in the group don't see the keyboard.
	Selective bool `json:"selective,omitempty"`
}

// ReplyKeyboardRemove Upon receiving a message with this object, Telegram clients will remove the current custom keyboard and display the default letter-keyboard.
// By default, custom keyboards are displayed until a new keyboard is sent by a bot.
// An exception is made for one-time keyboards that are hidden immediately after the user presses a button (see ReplyKeyboardMarkup).
// Not supported in channels and for messages sent on behalf of a Telegram Business account.
type ReplyKeyboardRemove struct {
	// Requests clients to remove the custom keyboard (user will not be able to summon this keyboard; if you want to hide the keyboard from sight but keep it accessible, use one_time_keyboard in ReplyKeyboardMarkup)
	RemoveKeyboard bool `json:"remove_keyboard"`
	// Optional. Use this parameter if you want to remove the keyboard for specific users only.
	// Targets: 1) users that are @mentioned in the text of the Message object; 2) if the bot's message is a reply to a message in the same chat and forum topic, sender of the original message.
	// Example: A user votes in a poll, bot returns confirmation message in reply to the vote and removes the keyboard for that user, while still showing the keyboard with poll options to users who haven't voted yet.
	Selective bool `json:"selective,omitempty"`
}

// ReplyParameters Describes reply parameters for the message that is being sent.
type ReplyParameters struct {
	// Identifier of the message that will be replied to in the current chat, or in the chat chat_id if it is specified
	MessageId int64 `json:"message_id"`
	// Optional. Not supported for messages sent on behalf of a business account.
	// If the message to be replied to is from a different chat, unique identifier for the chat or username of the channel (in the format @channelusername).
	// >> either: String
	ChatId int64 `json:"chat_id,omitempty"`
	// Optional. Pass True if the message should be sent even if the specified message to be replied to is not found.
	// Always False for replies in another chat or forum topic.
	// Always True for messages sent on behalf of a business account.
	AllowSendingWithoutReply bool `json:"allow_sending_without_reply,omitempty"`
	// Optional. Quoted part of the message to be replied to; 0-1024 characters after entities parsing.
	// The quote must be an exact substring of the message to be replied to, including bold, italic, underline, strikethrough, spoiler, and custom_emoji entities.
	// The message will fail to send if the quote isn't found in the original message.
	Quote string `json:"quote,omitempty"`
	// Optional. Mode for parsing entities in the quote. See formatting options for more details.
	QuoteParseMode string `json:"quote_parse_mode,omitempty"`
	// Optional. A JSON-serialized list of special entities that appear in the quote.
	// It can be specified instead of quote_parse_mode.
	QuoteEntities []*MessageEntity `json:"quote_entities,omitempty"`
	// Optional. Position of the quote in the original message in UTF-16 code units
	QuotePosition int64 `json:"quote_position,omitempty"`
}

// ResponseParameters Describes why a request was unsuccessful.
type ResponseParameters struct {
	// Optional. The group has been migrated to a supergroup with the specified identifier.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	MigrateToChatId int64 `json:"migrate_to_chat_id,omitempty"`
	// Optional. In case of exceeding flood control, the number of seconds left to wait before the request can be repeated
	RetryAfter int64 `json:"retry_after,omitempty"`
}

// RevenueWithdrawalState This object describes the state of a revenue withdrawal operation. Currently, it can be one of
// - RevenueWithdrawalStatePending
// - RevenueWithdrawalStateSucceeded
// - RevenueWithdrawalStateFailed
type RevenueWithdrawalState interface {
	OptPending() *RevenueWithdrawalStatePending
	OptSucceeded() *RevenueWithdrawalStateSucceeded
	OptFailed() *RevenueWithdrawalStateFailed
}

var (
	_ RevenueWithdrawalState = &RevenueWithdrawalStatePending{}
	_ RevenueWithdrawalState = &RevenueWithdrawalStateSucceeded{}
	_ RevenueWithdrawalState = &RevenueWithdrawalStateFailed{}
)

func (impl *RevenueWithdrawalStatePending) OptPending() *RevenueWithdrawalStatePending { return impl }
func (impl *RevenueWithdrawalStatePending) OptSucceeded() *RevenueWithdrawalStateSucceeded {
	return nil
}
func (impl *RevenueWithdrawalStatePending) OptFailed() *RevenueWithdrawalStateFailed { return nil }

func (impl *RevenueWithdrawalStateSucceeded) OptPending() *RevenueWithdrawalStatePending { return nil }
func (impl *RevenueWithdrawalStateSucceeded) OptSucceeded() *RevenueWithdrawalStateSucceeded {
	return impl
}
func (impl *RevenueWithdrawalStateSucceeded) OptFailed() *RevenueWithdrawalStateFailed { return nil }

func (impl *RevenueWithdrawalStateFailed) OptPending() *RevenueWithdrawalStatePending     { return nil }
func (impl *RevenueWithdrawalStateFailed) OptSucceeded() *RevenueWithdrawalStateSucceeded { return nil }
func (impl *RevenueWithdrawalStateFailed) OptFailed() *RevenueWithdrawalStateFailed       { return impl }

// RevenueWithdrawalStateFailed The withdrawal failed and the transaction was refunded.
type RevenueWithdrawalStateFailed struct {
	// Type of the state, always "failed"
	Type string `json:"type"`
}

// RevenueWithdrawalStatePending The withdrawal is in progress.
type RevenueWithdrawalStatePending struct {
	// Type of the state, always "pending"
	Type string `json:"type"`
}

// RevenueWithdrawalStateSucceeded The withdrawal succeeded.
type RevenueWithdrawalStateSucceeded struct {
	// Type of the state, always "succeeded"
	Type string `json:"type"`
	// Date the withdrawal was completed in Unix time
	Date int64 `json:"date"`
	// An HTTPS URL that can be used to see transaction details
	Url string `json:"url"`
}

// SentWebAppMessage Describes an inline message sent by a Web App on behalf of a user.
type SentWebAppMessage struct {
	// Optional. Identifier of the sent inline message.
	// Available only if there is an inline keyboard attached to the message.
	InlineMessageId string `json:"inline_message_id,omitempty"`
}

// SharedUser This object contains information about a user that was shared with the bot using a KeyboardButtonRequestUsers button.
type SharedUser struct {
	// Identifier of the shared user.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so 64-bit integers or double-precision float types are safe for storing these identifiers.
	// The bot may not have access to the user and could be unable to use this identifier, unless the user is already known to the bot by some other means.
	UserId int64 `json:"user_id"`
	// Optional. First name of the user, if the name was requested by the bot
	FirstName string `json:"first_name,omitempty"`
	// Optional. Last name of the user, if the name was requested by the bot
	LastName string `json:"last_name,omitempty"`
	// Optional. Username of the user, if the username was requested by the bot
	Username string `json:"username,omitempty"`
	// Optional. Available sizes of the chat photo, if the photo was requested by the bot
	Photo TelegramPhoto `json:"photo,omitempty"`
}

// ShippingAddress This object represents a shipping address.
type ShippingAddress struct {
	// Two-letter ISO 3166-1 alpha-2 country code
	CountryCode string `json:"country_code"`
	// State, if applicable
	State string `json:"state"`
	// City
	City string `json:"city"`
	// First line for the address
	StreetLine1 string `json:"street_line1"`
	// Second line for the address
	StreetLine2 string `json:"street_line2"`
	// Address post code
	PostCode string `json:"post_code"`
}

// ShippingOption This object represents one shipping option.
type ShippingOption struct {
	// Shipping option identifier
	Id string `json:"id"`
	// Option title
	Title string `json:"title"`
	// List of price portions
	Prices []*LabeledPrice `json:"prices"`
}

// ShippingQuery This object contains information about an incoming shipping query.
type ShippingQuery struct {
	// Unique query identifier
	Id string `json:"id"`
	// User who sent the query
	From *User `json:"from"`
	// Bot-specified invoice payload
	InvoicePayload string `json:"invoice_payload"`
	// User specified shipping address
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

// StarTransaction Describes a Telegram Star transaction.
type StarTransaction struct {
	// Unique identifier of the transaction.
	// Coincides with the identifier of the original transaction for refund transactions.
	// Coincides with SuccessfulPayment.telegram_payment_charge_id for successful incoming payments from users.
	Id string `json:"id"`
	// Integer amount of Telegram Stars transferred by the transaction
	Amount int64 `json:"amount"`
	// Optional. The number of 1/1000000000 shares of Telegram Stars transferred by the transaction; from 0 to 999999999
	NanostarAmount int64 `json:"nanostar_amount,omitempty"`
	// Date the transaction was created in Unix time
	Date int64 `json:"date"`
	// Optional. Only for incoming transactions
	// Source of an incoming transaction (e.g., a user purchasing goods or services, Fragment refunding a failed withdrawal).
	Source TransactionPartner `json:"source,omitempty"`
	// Optional. Receiver of an outgoing transaction (e.g., a user for a purchase refund, Fragment for a withdrawal).
	// Only for outgoing transactions
	Receiver TransactionPartner `json:"receiver,omitempty"`
}

// StarTransactions Contains a list of Telegram Star transactions.
type StarTransactions struct {
	// The list of transactions
	Transactions []*StarTransaction `json:"transactions"`
}

// Sticker This object represents a sticker.
type Sticker struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Type of the sticker, currently one of "regular", "mask", "custom_emoji".
	// The type of the sticker is independent from its format, which is determined by the fields is_animated and is_video.
	Type string `json:"type"`
	// Sticker width
	Width int64 `json:"width"`
	// Sticker height
	Height int64 `json:"height"`
	// True, if the sticker is animated
	IsAnimated bool `json:"is_animated"`
	// True, if the sticker is a video sticker
	IsVideo bool `json:"is_video"`
	// Optional. Sticker thumbnail in the .WEBP or .JPG format
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
	// Optional. Emoji associated with the sticker
	Emoji string `json:"emoji,omitempty"`
	// Optional. Name of the sticker set to which the sticker belongs
	SetName string `json:"set_name,omitempty"`
	// Optional. For premium regular stickers, premium animation for the sticker
	PremiumAnimation *File `json:"premium_animation,omitempty"`
	// Optional. For mask stickers, the position where the mask should be placed
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	// Optional. For custom emoji stickers, unique identifier of the custom emoji
	CustomEmojiId string `json:"custom_emoji_id,omitempty"`
	// Optional.
	// True, if the sticker must be repainted to a text color in messages, the color of the Telegram Premium badge in emoji status, white color on chat photos, or another appropriate color in other places
	NeedsRepainting bool `json:"needs_repainting,omitempty"`
	// Optional. File size in bytes
	FileSize int64 `json:"file_size,omitempty"`
}

func (impl *Sticker) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *Sticker) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// StickerSet This object represents a sticker set.
type StickerSet struct {
	// Sticker set name
	Name string `json:"name"`
	// Sticker set title
	Title string `json:"title"`
	// Type of stickers in the set, currently one of "regular", "mask", "custom_emoji"
	StickerType string `json:"sticker_type"`
	// List of all set stickers
	Stickers []*Sticker `json:"stickers"`
	// Optional. Sticker set thumbnail in the .WEBP, .TGS, or .WEBM format
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
}

// Story This object represents a story.
type Story struct {
	// Chat that posted the story
	Chat *Chat `json:"chat"`
	// Unique identifier for the story in the chat
	Id int64 `json:"id"`
}

// SuccessfulPayment This object contains basic information about a successful payment.
type SuccessfulPayment struct {
	// Three-letter ISO 4217 currency code, or "XTR" for payments in Telegram Stars
	Currency string `json:"currency"`
	// Total price in the smallest units of the currency (integer, not float/double).
	// For example, for a price of US$ 1.45 pass amount = 145.
	// See the exp parameter in currencies.json, it shows the number of digits past the decimal point for each currency (2 for the majority of currencies).
	TotalAmount int64 `json:"total_amount"`
	// Bot-specified invoice payload
	InvoicePayload string `json:"invoice_payload"`
	// Optional. Expiration date of the subscription, in Unix time; for recurring payments only
	SubscriptionExpirationDate int64 `json:"subscription_expiration_date,omitempty"`
	// Optional. True, if the payment is a recurring payment for a subscription
	IsRecurring bool `json:"is_recurring,omitempty"`
	// Optional. True, if the payment is the first payment for a subscription
	IsFirstRecurring bool `json:"is_first_recurring,omitempty"`
	// Optional. Identifier of the shipping option chosen by the user
	ShippingOptionId string `json:"shipping_option_id,omitempty"`
	// Optional. Order information provided by the user
	OrderInfo *OrderInfo `json:"order_info,omitempty"`
	// Telegram payment identifier
	TelegramPaymentChargeId string `json:"telegram_payment_charge_id"`
	// Provider payment identifier
	ProviderPaymentChargeId string `json:"provider_payment_charge_id"`
}

// SwitchInlineQueryChosenChat This object represents an inline button that switches the current user to inline mode in a chosen chat, with an optional default inline query.
type SwitchInlineQueryChosenChat struct {
	// Optional. The default inline query to be inserted in the input field.
	// If left empty, only the bot's username will be inserted
	Query string `json:"query,omitempty"`
	// Optional. True, if private chats with users can be chosen
	AllowUserChats bool `json:"allow_user_chats,omitempty"`
	// Optional. True, if private chats with bots can be chosen
	AllowBotChats bool `json:"allow_bot_chats,omitempty"`
	// Optional. True, if group and supergroup chats can be chosen
	AllowGroupChats bool `json:"allow_group_chats,omitempty"`
	// Optional. True, if channel chats can be chosen
	AllowChannelChats bool `json:"allow_channel_chats,omitempty"`
}

// TelegramAnimation This object represents an animation file (GIF or H.264/MPEG-4 AVC video without sound).
type TelegramAnimation struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Video width as defined by the sender
	Width int64 `json:"width"`
	// Video height as defined by the sender
	Height int64 `json:"height"`
	// Duration of the video in seconds as defined by the sender
	Duration int64 `json:"duration"`
	// Optional. Animation thumbnail as defined by the sender
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
	// Optional. Original animation filename as defined by the sender
	FileName string `json:"file_name,omitempty"`
	// Optional. MIME type of the file as defined by the sender
	MimeType string `json:"mime_type,omitempty"`
	// Optional. File size in bytes.
	// It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
	FileSize int64 `json:"file_size,omitempty"`
}

func (impl *TelegramAnimation) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *TelegramAnimation) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// TelegramAudio This object represents an audio file to be treated as music by the Telegram clients.
type TelegramAudio struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Duration of the audio in seconds as defined by the sender
	Duration int64 `json:"duration"`
	// Optional. Performer of the audio as defined by the sender or by audio tags
	Performer string `json:"performer,omitempty"`
	// Optional. Title of the audio as defined by the sender or by audio tags
	Title string `json:"title,omitempty"`
	// Optional. Original filename as defined by the sender
	FileName string `json:"file_name,omitempty"`
	// Optional. MIME type of the file as defined by the sender
	MimeType string `json:"mime_type,omitempty"`
	// Optional. File size in bytes.
	// It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
	FileSize int64 `json:"file_size,omitempty"`
	// Optional. Thumbnail of the album cover to which the music file belongs
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
}

func (impl *TelegramAudio) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *TelegramAudio) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// TelegramDocument This object represents a general file (as opposed to photos, voice messages and audio files).
type TelegramDocument struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Optional. Document thumbnail as defined by the sender
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
	// Optional. Original filename as defined by the sender
	FileName string `json:"file_name,omitempty"`
	// Optional. MIME type of the file as defined by the sender
	MimeType string `json:"mime_type,omitempty"`
	// Optional. File size in bytes.
	// It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
	FileSize int64 `json:"file_size,omitempty"`
}

func (impl *TelegramDocument) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *TelegramDocument) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// TelegramVideo This object represents a video file.
type TelegramVideo struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Video width as defined by the sender
	Width int64 `json:"width"`
	// Video height as defined by the sender
	Height int64 `json:"height"`
	// Duration of the video in seconds as defined by the sender
	Duration int64 `json:"duration"`
	// Optional. Video thumbnail
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
	// Optional. Original filename as defined by the sender
	FileName string `json:"file_name,omitempty"`
	// Optional. MIME type of the file as defined by the sender
	MimeType string `json:"mime_type,omitempty"`
	// Optional. File size in bytes.
	// It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
	FileSize int64 `json:"file_size,omitempty"`
}

func (impl *TelegramVideo) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *TelegramVideo) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// TextQuote This object contains information about the quoted part of a message that is replied to by the given message.
type TextQuote struct {
	// Text of the quoted part of a message that is replied to by the given message
	Text string `json:"text"`
	// Optional. Special entities that appear in the quote.
	// Currently, only bold, italic, underline, strikethrough, spoiler, and custom_emoji entities are kept in quotes.
	Entities []*MessageEntity `json:"entities,omitempty"`
	// Approximate quote position in the original message in UTF-16 code units as specified by the sender
	Position int64 `json:"position"`
	// Optional. True, if the quote was chosen manually by the message sender.
	// Otherwise, the quote was added automatically by the server.
	IsManual bool `json:"is_manual,omitempty"`
}

// TransactionPartner This object describes the source of a transaction, or its recipient for outgoing transactions. Currently, it can be one of
// - TransactionPartnerUser
// - TransactionPartnerAffiliateProgram
// - TransactionPartnerFragment
// - TransactionPartnerTelegramAds
// - TransactionPartnerTelegramApi
// - TransactionPartnerOther
type TransactionPartner interface {
	OptUser() *TransactionPartnerUser
	OptAffiliateProgram() *TransactionPartnerAffiliateProgram
	OptFragment() *TransactionPartnerFragment
	OptTelegramAds() *TransactionPartnerTelegramAds
	OptTelegramApi() *TransactionPartnerTelegramApi
	OptOther() *TransactionPartnerOther
}

var (
	_ TransactionPartner = &TransactionPartnerUser{}
	_ TransactionPartner = &TransactionPartnerAffiliateProgram{}
	_ TransactionPartner = &TransactionPartnerFragment{}
	_ TransactionPartner = &TransactionPartnerTelegramAds{}
	_ TransactionPartner = &TransactionPartnerTelegramApi{}
	_ TransactionPartner = &TransactionPartnerOther{}
)

func (impl *TransactionPartnerUser) OptUser() *TransactionPartnerUser { return impl }
func (impl *TransactionPartnerUser) OptAffiliateProgram() *TransactionPartnerAffiliateProgram {
	return nil
}
func (impl *TransactionPartnerUser) OptFragment() *TransactionPartnerFragment       { return nil }
func (impl *TransactionPartnerUser) OptTelegramAds() *TransactionPartnerTelegramAds { return nil }
func (impl *TransactionPartnerUser) OptTelegramApi() *TransactionPartnerTelegramApi { return nil }
func (impl *TransactionPartnerUser) OptOther() *TransactionPartnerOther             { return nil }

func (impl *TransactionPartnerAffiliateProgram) OptUser() *TransactionPartnerUser { return nil }
func (impl *TransactionPartnerAffiliateProgram) OptAffiliateProgram() *TransactionPartnerAffiliateProgram {
	return impl
}
func (impl *TransactionPartnerAffiliateProgram) OptFragment() *TransactionPartnerFragment { return nil }
func (impl *TransactionPartnerAffiliateProgram) OptTelegramAds() *TransactionPartnerTelegramAds {
	return nil
}
func (impl *TransactionPartnerAffiliateProgram) OptTelegramApi() *TransactionPartnerTelegramApi {
	return nil
}
func (impl *TransactionPartnerAffiliateProgram) OptOther() *TransactionPartnerOther { return nil }

func (impl *TransactionPartnerFragment) OptUser() *TransactionPartnerUser { return nil }
func (impl *TransactionPartnerFragment) OptAffiliateProgram() *TransactionPartnerAffiliateProgram {
	return nil
}
func (impl *TransactionPartnerFragment) OptFragment() *TransactionPartnerFragment       { return impl }
func (impl *TransactionPartnerFragment) OptTelegramAds() *TransactionPartnerTelegramAds { return nil }
func (impl *TransactionPartnerFragment) OptTelegramApi() *TransactionPartnerTelegramApi { return nil }
func (impl *TransactionPartnerFragment) OptOther() *TransactionPartnerOther             { return nil }

func (impl *TransactionPartnerTelegramAds) OptUser() *TransactionPartnerUser { return nil }
func (impl *TransactionPartnerTelegramAds) OptAffiliateProgram() *TransactionPartnerAffiliateProgram {
	return nil
}
func (impl *TransactionPartnerTelegramAds) OptFragment() *TransactionPartnerFragment { return nil }
func (impl *TransactionPartnerTelegramAds) OptTelegramAds() *TransactionPartnerTelegramAds {
	return impl
}
func (impl *TransactionPartnerTelegramAds) OptTelegramApi() *TransactionPartnerTelegramApi {
	return nil
}
func (impl *TransactionPartnerTelegramAds) OptOther() *TransactionPartnerOther { return nil }

func (impl *TransactionPartnerTelegramApi) OptUser() *TransactionPartnerUser { return nil }
func (impl *TransactionPartnerTelegramApi) OptAffiliateProgram() *TransactionPartnerAffiliateProgram {
	return nil
}
func (impl *TransactionPartnerTelegramApi) OptFragment() *TransactionPartnerFragment { return nil }
func (impl *TransactionPartnerTelegramApi) OptTelegramAds() *TransactionPartnerTelegramAds {
	return nil
}
func (impl *TransactionPartnerTelegramApi) OptTelegramApi() *TransactionPartnerTelegramApi {
	return impl
}
func (impl *TransactionPartnerTelegramApi) OptOther() *TransactionPartnerOther { return nil }

func (impl *TransactionPartnerOther) OptUser() *TransactionPartnerUser { return nil }
func (impl *TransactionPartnerOther) OptAffiliateProgram() *TransactionPartnerAffiliateProgram {
	return nil
}
func (impl *TransactionPartnerOther) OptFragment() *TransactionPartnerFragment       { return nil }
func (impl *TransactionPartnerOther) OptTelegramAds() *TransactionPartnerTelegramAds { return nil }
func (impl *TransactionPartnerOther) OptTelegramApi() *TransactionPartnerTelegramApi { return nil }
func (impl *TransactionPartnerOther) OptOther() *TransactionPartnerOther             { return impl }

// TransactionPartnerAffiliateProgram Describes the affiliate program that issued the affiliate commission received via this transaction.
type TransactionPartnerAffiliateProgram struct {
	// Type of the transaction partner, always "affiliate_program"
	Type string `json:"type"`
	// Optional. Information about the bot that sponsored the affiliate program
	SponsorUser *User `json:"sponsor_user,omitempty"`
	// The number of Telegram Stars received by the bot for each 1000 Telegram Stars received by the affiliate program sponsor from referred users
	CommissionPerMille int64 `json:"commission_per_mille"`
}

// TransactionPartnerFragment Describes a withdrawal transaction with Fragment.
type TransactionPartnerFragment struct {
	// Type of the transaction partner, always "fragment"
	Type string `json:"type"`
	// Optional. State of the transaction if the transaction is outgoing
	WithdrawalState RevenueWithdrawalState `json:"withdrawal_state,omitempty"`
}

// TransactionPartnerOther Describes a transaction with an unknown source or recipient.
type TransactionPartnerOther struct {
	// Type of the transaction partner, always "other"
	Type string `json:"type"`
}

// TransactionPartnerTelegramAds Describes a withdrawal transaction to the Telegram Ads platform.
type TransactionPartnerTelegramAds struct {
	// Type of the transaction partner, always "telegram_ads"
	Type string `json:"type"`
}

// TransactionPartnerTelegramApi Describes a transaction with payment for paid broadcasting.
type TransactionPartnerTelegramApi struct {
	// Type of the transaction partner, always "telegram_api"
	Type string `json:"type"`
	// The number of successful requests that exceeded regular limits and were therefore billed
	RequestCount int64 `json:"request_count"`
}

// TransactionPartnerUser Describes a transaction with a user.
type TransactionPartnerUser struct {
	// Type of the transaction partner, always "user"
	Type string `json:"type"`
	// Information about the user
	User *User `json:"user"`
	// Optional. Information about the affiliate that received a commission via this transaction
	Affiliate *AffiliateInfo `json:"affiliate,omitempty"`
	// Optional. Bot-specified invoice payload
	InvoicePayload string `json:"invoice_payload,omitempty"`
	// Optional. The duration of the paid subscription
	SubscriptionPeriod int64 `json:"subscription_period,omitempty"`
	// Optional. Information about the paid media bought by the user
	PaidMedia []PaidMedia `json:"paid_media,omitempty"`
	// Optional. Bot-specified paid media payload
	PaidMediaPayload string `json:"paid_media_payload,omitempty"`
	// Optional. The gift sent to the user by the bot
	Gift *Gift `json:"gift,omitempty"`
}

// Update This object represents an incoming update.
// At most one of the optional parameters can be present in any given update.
type Update struct {
	// The update's unique identifier. Update identifiers start from a certain positive number and increase sequentially.
	// This identifier becomes especially handy if you're using webhooks, since it allows you to ignore repeated updates or to restore the correct update sequence, should they get out of order.
	// If there are no new updates for at least a week, then identifier of the next update will be chosen randomly instead of sequentially.
	UpdateId int64 `json:"update_id"`
	// Optional. New incoming message of any kind - text, photo, sticker, etc.
	Message *Message `json:"message,omitempty"`
	// Optional. New version of a message that is known to the bot and was edited.
	// This update may at times be triggered by changes to message fields that are either unavailable or not actively used by your bot.
	EditedMessage *Message `json:"edited_message,omitempty"`
	// Optional. New incoming channel post of any kind - text, photo, sticker, etc.
	ChannelPost *Message `json:"channel_post,omitempty"`
	// Optional. New version of a channel post that is known to the bot and was edited.
	// This update may at times be triggered by changes to message fields that are either unavailable or not actively used by your bot.
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"`
	// Optional.
	// The bot was connected to or disconnected from a business account, or a user edited an existing connection with the bot
	BusinessConnection *BusinessConnection `json:"business_connection,omitempty"`
	// Optional. New message from a connected business account
	BusinessMessage *Message `json:"business_message,omitempty"`
	// Optional. New version of a message from a connected business account
	EditedBusinessMessage *Message `json:"edited_business_message,omitempty"`
	// Optional. Messages were deleted from a connected business account
	DeletedBusinessMessages *BusinessMessagesDeleted `json:"deleted_business_messages,omitempty"`
	// Optional. A reaction to a message was changed by a user. The update isn't received for reactions set by bots.
	// The bot must be an administrator in the chat and must explicitly specify "message_reaction" in the list of allowed_updates to receive these updates.
	MessageReaction *MessageReactionUpdated `json:"message_reaction,omitempty"`
	// Optional. Reactions to a message with anonymous reactions were changed.
	// The bot must be an administrator in the chat and must explicitly specify "message_reaction_count" in the list of allowed_updates to receive these updates.
	// The updates are grouped and can be sent with delay up to a few minutes.
	MessageReactionCount *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`
	// Optional. New incoming inline query
	InlineQuery *InlineQuery `json:"inline_query,omitempty"`
	// Optional. The result of an inline query that was chosen by a user and sent to their chat partner.
	// Please see our documentation on the feedback collecting for details on how to enable these updates for your bot.
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	// Optional. New incoming callback query
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
	// Optional. New incoming shipping query. Only for invoices with flexible price
	ShippingQuery *ShippingQuery `json:"shipping_query,omitempty"`
	// Optional. New incoming pre-checkout query. Contains full information about checkout
	PreCheckoutQuery *PreCheckoutQuery `json:"pre_checkout_query,omitempty"`
	// Optional. A user purchased paid media with a non-empty payload sent by the bot in a non-channel chat
	PurchasedPaidMedia *PaidMediaPurchased `json:"purchased_paid_media,omitempty"`
	// Optional. New poll state. Bots receive only updates about manually stopped polls and polls, which are sent by the bot
	Poll *Poll `json:"poll,omitempty"`
	// Optional. A user changed their answer in a non-anonymous poll.
	// Bots receive new votes only in polls that were sent by the bot itself.
	PollAnswer *PollAnswer `json:"poll_answer,omitempty"`
	// Optional. The bot's chat member status was updated in a chat.
	// For private chats, this update is received only when the bot is blocked or unblocked by the user.
	MyChatMember *ChatMemberUpdated `json:"my_chat_member,omitempty"`
	// Optional. A chat member's status was updated in a chat.
	// The bot must be an administrator in the chat and must explicitly specify "chat_member" in the list of allowed_updates to receive these updates.
	ChatMember *ChatMemberUpdated `json:"chat_member,omitempty"`
	// Optional. A request to join the chat has been sent.
	// The bot must have the can_invite_users administrator right in the chat to receive these updates.
	ChatJoinRequest *ChatJoinRequest `json:"chat_join_request,omitempty"`
	// Optional. A chat boost was added or changed. The bot must be an administrator in the chat to receive these updates.
	ChatBoost *ChatBoostUpdated `json:"chat_boost,omitempty"`
	// Optional. A boost was removed from a chat. The bot must be an administrator in the chat to receive these updates.
	RemovedChatBoost *ChatBoostRemoved `json:"removed_chat_boost,omitempty"`
}

// User This object represents a Telegram user or bot.
type User struct {
	// Unique identifier for this user or bot.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier.
	Id int64 `json:"id"`
	// True, if this user is a bot
	IsBot bool `json:"is_bot"`
	// User's or bot's first name
	FirstName string `json:"first_name"`
	// Optional. User's or bot's last name
	LastName string `json:"last_name,omitempty"`
	// Optional. User's or bot's username
	Username string `json:"username,omitempty"`
	// Optional. IETF language tag of the user's language
	LanguageCode string `json:"language_code,omitempty"`
	// Optional. True, if this user is a Telegram Premium user
	IsPremium bool `json:"is_premium,omitempty"`
	// Optional. True, if this user added the bot to the attachment menu
	AddedToAttachmentMenu bool `json:"added_to_attachment_menu,omitempty"`
	// Optional. True, if the bot can be invited to groups. Returned only in getMe.
	CanJoinGroups bool `json:"can_join_groups,omitempty"`
	// Optional. True, if privacy mode is disabled for the bot. Returned only in getMe.
	CanReadAllGroupMessages bool `json:"can_read_all_group_messages,omitempty"`
	// Optional. True, if the bot supports inline queries. Returned only in getMe.
	SupportsInlineQueries bool `json:"supports_inline_queries,omitempty"`
	// Optional. True, if the bot can be connected to a Telegram Business account to receive its messages.
	// Returned only in getMe.
	CanConnectToBusiness bool `json:"can_connect_to_business,omitempty"`
	// Optional. True, if the bot has a main Web App. Returned only in getMe.
	HasMainWebApp bool `json:"has_main_web_app,omitempty"`
}

// UserChatBoosts This object represents a list of boosts added to a chat by a user.
type UserChatBoosts struct {
	// The list of boosts added to the chat by the user
	Boosts []*ChatBoost `json:"boosts"`
}

// UserProfilePhotos This object represent a user's profile pictures.
type UserProfilePhotos struct {
	// Total number of profile pictures the target user has
	TotalCount int64 `json:"total_count"`
	// Requested profile pictures (in up to 4 sizes each)
	Photos []TelegramPhoto `json:"photos"`
}

// UsersShared This object contains information about the users whose identifiers were shared with the bot using a KeyboardButtonRequestUsers button.
type UsersShared struct {
	// Identifier of the request
	RequestId int64 `json:"request_id"`
	// Information about users shared with the bot.
	Users []*SharedUser `json:"users"`
}

// Venue This object represents a venue.
type Venue struct {
	// Venue location. Can't be a live location
	Location *Location `json:"location"`
	// Name of the venue
	Title string `json:"title"`
	// Address of the venue
	Address string `json:"address"`
	// Optional. Foursquare identifier of the venue
	FoursquareId string `json:"foursquare_id,omitempty"`
	// Optional. Foursquare type of the venue.
	// (For example, "arts_entertainment/default", "arts_entertainment/aquarium" or "food/icecream".)
	FoursquareType string `json:"foursquare_type,omitempty"`
	// Optional. Google Places identifier of the venue
	GooglePlaceId string `json:"google_place_id,omitempty"`
	// Optional. Google Places type of the venue. (See supported types.)
	GooglePlaceType string `json:"google_place_type,omitempty"`
}

// Video Represents a video to be sent.
type Video struct {
	// Type of the result, must be video
	Type string `json:"type" default:"video"`
	// File to send. More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// Pass a file_id to send a file that exists on the Telegram servers (recommended), pass an HTTP URL for Telegram to get a file from the Internet, or pass "attach://<file_attach_name>" to upload a new one using multipart/form-data under <file_attach_name> name.
	Media InputFile `json:"media"`
	// Optional. Thumbnail of the file sent; can be ignored if thumbnail generation for the file is supported server-side.
	// The thumbnail should be in JPEG format and less than 200 kB in size.
	// A thumbnail's width and height should not exceed 320.
	// Ignored if the file is not uploaded using multipart/form-data.
	// Thumbnails can't be reused and can be only uploaded as a new file, so you can pass "attach://<file_attach_name>" if the thumbnail was uploaded using multipart/form-data under <file_attach_name>.
	// More information on Sending Files: https://core.telegram.org/bots/api#sending-files
	// >> either: String
	Thumbnail InputFile `json:"thumbnail,omitempty"`
	// Optional. Caption of the video to be sent, 0-1024 characters after entities parsing
	Caption string `json:"caption,omitempty"`
	// Optional. Mode for parsing entities in the video caption. See formatting options for more details.
	ParseMode string `json:"parse_mode,omitempty"`
	// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`
	// Optional. Pass True, if the caption must be shown above the message media
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
	// Optional. Video width
	Width int64 `json:"width,omitempty"`
	// Optional. Video height
	Height int64 `json:"height,omitempty"`
	// Optional. Video duration in seconds
	Duration int64 `json:"duration,omitempty"`
	// Optional. Pass True if the uploaded video is suitable for streaming
	SupportsStreaming bool `json:"supports_streaming,omitempty"`
	// Optional. Pass True if the video needs to be covered with a spoiler animation
	HasSpoiler bool `json:"has_spoiler,omitempty"`
	// Used for uploading media.
	InputFile InputFile `json:"-"`
}

// VideoChatEnded This object represents a service message about a video chat ended in the chat.
type VideoChatEnded struct {
	// Video chat duration in seconds
	Duration int64 `json:"duration"`
}

// VideoChatParticipantsInvited This object represents a service message about new members invited to a video chat.
type VideoChatParticipantsInvited struct {
	// New members that were invited to the video chat
	Users []*User `json:"users"`
}

// VideoChatScheduled This object represents a service message about a video chat scheduled in the chat.
type VideoChatScheduled struct {
	// Point in time (Unix timestamp) when the video chat is supposed to be started by a chat administrator
	StartDate int64 `json:"start_date"`
}

// VideoChatStarted This object represents a service message about a video chat started in the chat.
// Currently holds no information.
type VideoChatStarted struct {
}

// VideoNote This object represents a video message (available in Telegram apps as of v.4.0).
type VideoNote struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Video width and height (diameter of the video message) as defined by the sender
	Length int64 `json:"length"`
	// Duration of the video in seconds as defined by the sender
	Duration int64 `json:"duration"`
	// Optional. Video thumbnail
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
	// Optional. File size in bytes
	FileSize int64 `json:"file_size,omitempty"`
}

func (impl *VideoNote) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *VideoNote) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// Voice This object represents a voice note.
type Voice struct {
	// Identifier for this file, which can be used to download or reuse the file
	FileId string `json:"file_id"`
	// Unique identifier for this file, which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueId string `json:"file_unique_id"`
	// Duration of the audio in seconds as defined by the sender
	Duration int64 `json:"duration"`
	// Optional. MIME type of the file as defined by the sender
	MimeType string `json:"mime_type,omitempty"`
	// Optional. File size in bytes.
	// It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
	FileSize int64 `json:"file_size,omitempty"`
}

func (impl *Voice) Download(ctx context.Context, path string) error {
	return GenericDownload(ctx, path, impl.FileId)
}
func (impl *Voice) DownloadTemp(ctx context.Context, dirAndPattern ...string) (filename string, err error) {
	return GenericDownloadTemp(ctx, impl.FileId, dirAndPattern...)
}

// WebAppData Describes data sent from a Web App to the bot.
type WebAppData struct {
	// The data. Be aware that a bad client can send arbitrary data in this field.
	Data string `json:"data"`
	// Text of the web_app keyboard button from which the Web App was opened.
	// Be aware that a bad client can send arbitrary data in this field.
	ButtonText string `json:"button_text"`
}

// WebAppInfo Describes a Web App.
type WebAppInfo struct {
	// An HTTPS URL of a Web App to be opened with additional data as specified in Initializing Web Apps
	Url string `json:"url"`
}

// WebhookInfo Describes the current status of a webhook.
type WebhookInfo struct {
	// Webhook URL, may be empty if webhook is not set up
	Url string `json:"url"`
	// True, if a custom certificate was provided for webhook certificate checks
	HasCustomCertificate bool `json:"has_custom_certificate"`
	// Number of updates awaiting delivery
	PendingUpdateCount int64 `json:"pending_update_count"`
	// Optional. Currently used webhook IP address
	IpAddress string `json:"ip_address,omitempty"`
	// Optional. Unix time for the most recent error that happened when trying to deliver an update via webhook
	LastErrorDate int64 `json:"last_error_date,omitempty"`
	// Optional.
	// Error message in human-readable format for the most recent error that happened when trying to deliver an update via webhook
	LastErrorMessage string `json:"last_error_message,omitempty"`
	// Optional.
	// Unix time of the most recent error that happened when trying to synchronize available updates with Telegram datacenters
	LastSynchronizationErrorDate int64 `json:"last_synchronization_error_date,omitempty"`
	// Optional. The maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery
	MaxConnections int64 `json:"max_connections,omitempty"`
	// Optional. A list of update types the bot is subscribed to. Defaults to all update types except chat_member
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

// WriteAccessAllowed This object represents a service message about a user allowing a bot to write messages after adding it to the attachment menu, launching a Web App from a link, or accepting an explicit request from a Web App sent by the method requestWriteAccess.
type WriteAccessAllowed struct {
	// Optional.
	// True, if the access was granted after the user accepted an explicit request from a Web App sent by the method requestWriteAccess
	FromRequest bool `json:"from_request,omitempty"`
	// Optional. Name of the Web App, if the access was granted when the Web App was launched from a link
	WebAppName string `json:"web_app_name,omitempty"`
	// Optional. True, if the access was granted when the bot was added to the attachment or side menu
	FromAttachmentMenu bool `json:"from_attachment_menu,omitempty"`
}
