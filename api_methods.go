package tg

import (
	"context"
)

// AddStickerToSet Use this method to add a new sticker to a set created by the bot. Emoji sticker sets can have up to 200 stickers.
// Other sticker sets can have up to 120 stickers. Returns True on success.
func AddStickerToSet(ctx context.Context, userId int64, name string, sticker *InputSticker) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId  int64         `json:"user_id"`
		Name    string        `json:"name"`
		Sticker *InputSticker `json:"sticker"`
	}
	request := &Request{
		UserId:  userId,
		Name:    name,
		Sticker: sticker,
	}
	return GenericRequest[Request, bool](ctx, "addStickerToSet", request)
}

// AnswerCallbackQuery Use this method to send answers to callback queries sent from inline keyboards. On success, True is returned.
// The answer will be displayed to the user as a notification at the top of the chat screen or as an alert.
func AnswerCallbackQuery(ctx context.Context, callbackQueryId string, opts ...*OptAnswerCallbackQuery) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		CallbackQueryId string `json:"callback_query_id"`
		Text            string `json:"text,omitempty"`
		ShowAlert       bool   `json:"show_alert,omitempty"`
		Url             string `json:"url,omitempty"`
		CacheTime       int64  `json:"cache_time,omitempty"`
	}
	request := &Request{
		CallbackQueryId: callbackQueryId,
	}
	for _, opt := range opts {
		if opt.Text != "" {
			request.Text = opt.Text
		}
		if opt.ShowAlert {
			request.ShowAlert = opt.ShowAlert
		}
		if opt.Url != "" {
			request.Url = opt.Url
		}
		if opt.CacheTime != 0 {
			request.CacheTime = opt.CacheTime
		}
	}
	return GenericRequest[Request, bool](ctx, "answerCallbackQuery", request)
}

type OptAnswerCallbackQuery struct {
	Text      string
	ShowAlert bool
	Url       string
	CacheTime int64
}

// AnswerInlineQuery Use this method to send answers to an inline query. On success, True is returned.
// No more than 50 results per query are allowed.
func AnswerInlineQuery(ctx context.Context, inlineQueryId string, results []InlineQueryResult, opts ...*OptAnswerInlineQuery) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		InlineQueryId string                    `json:"inline_query_id"`
		Results       []InlineQueryResult       `json:"results"`
		CacheTime     int64                     `json:"cache_time,omitempty"`
		IsPersonal    bool                      `json:"is_personal,omitempty"`
		NextOffset    string                    `json:"next_offset,omitempty"`
		Button        *InlineQueryResultsButton `json:"button,omitempty"`
	}
	request := &Request{
		InlineQueryId: inlineQueryId,
		Results:       results,
	}
	for _, opt := range opts {
		if opt.CacheTime != 0 {
			request.CacheTime = opt.CacheTime
		}
		if opt.IsPersonal {
			request.IsPersonal = opt.IsPersonal
		}
		if opt.NextOffset != "" {
			request.NextOffset = opt.NextOffset
		}
		if opt.Button != nil {
			request.Button = opt.Button
		}
	}
	return GenericRequest[Request, bool](ctx, "answerInlineQuery", request)
}

type OptAnswerInlineQuery struct {
	CacheTime  int64
	IsPersonal bool
	NextOffset string
	Button     *InlineQueryResultsButton
}

// AnswerPreCheckoutQuery Once the user has confirmed their payment and shipping details, the Bot API sends the final confirmation in the form of an Update with the field pre_checkout_query.
// Use this method to respond to such pre-checkout queries. On success, True is returned.
// Note: The Bot API must receive an answer within 10 seconds after the pre-checkout query was sent.
func AnswerPreCheckoutQuery(ctx context.Context, preCheckoutQueryId string, ok bool, opts ...*OptAnswerPreCheckoutQuery) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		PreCheckoutQueryId string `json:"pre_checkout_query_id"`
		Ok                 bool   `json:"ok"`
		ErrorMessage       string `json:"error_message,omitempty"`
	}
	request := &Request{
		PreCheckoutQueryId: preCheckoutQueryId,
		Ok:                 ok,
	}
	for _, opt := range opts {
		if opt.ErrorMessage != "" {
			request.ErrorMessage = opt.ErrorMessage
		}
	}
	return GenericRequest[Request, bool](ctx, "answerPreCheckoutQuery", request)
}

type OptAnswerPreCheckoutQuery struct {
	ErrorMessage string
}

// AnswerShippingQuery If you sent an invoice requesting a shipping address and the parameter is_flexible was specified, the Bot API will send an Update with a shipping_query field to the bot.
// Use this method to reply to shipping queries. On success, True is returned.
func AnswerShippingQuery(ctx context.Context, shippingQueryId string, ok bool, opts ...*OptAnswerShippingQuery) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		ShippingQueryId string            `json:"shipping_query_id"`
		Ok              bool              `json:"ok"`
		ShippingOptions []*ShippingOption `json:"shipping_options,omitempty"`
		ErrorMessage    string            `json:"error_message,omitempty"`
	}
	request := &Request{
		ShippingQueryId: shippingQueryId,
		Ok:              ok,
	}
	for _, opt := range opts {
		if opt.ShippingOptions != nil {
			request.ShippingOptions = opt.ShippingOptions
		}
		if opt.ErrorMessage != "" {
			request.ErrorMessage = opt.ErrorMessage
		}
	}
	return GenericRequest[Request, bool](ctx, "answerShippingQuery", request)
}

type OptAnswerShippingQuery struct {
	ShippingOptions []*ShippingOption
	ErrorMessage    string
}

// AnswerWebAppQuery Use this method to set the result of an interaction with a Web App and send a corresponding message on behalf of the user to the chat from which the query originated.
// On success, a SentWebAppMessage object is returned.
func AnswerWebAppQuery(ctx context.Context, webAppQueryId string, result InlineQueryResult) (*SentWebAppMessage, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		WebAppQueryId string            `json:"web_app_query_id"`
		Result        InlineQueryResult `json:"result"`
	}
	request := &Request{
		WebAppQueryId: webAppQueryId,
		Result:        result,
	}
	return GenericRequest[Request, *SentWebAppMessage](ctx, "answerWebAppQuery", request)
}

// ApproveChatJoinRequest Use this method to approve a chat join request. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_invite_users administrator right.
func ApproveChatJoinRequest(ctx context.Context, chatId int64, userId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
		UserId int64 `json:"user_id"`
	}
	request := &Request{
		ChatId: chatId,
		UserId: userId,
	}
	return GenericRequest[Request, bool](ctx, "approveChatJoinRequest", request)
}

// BanChatMember Use this method to ban a user in a group, a supergroup or a channel. Returns True on success.
// In the case of supergroups and channels, the user will not be able to return to the chat on their own using invite links, etc., unless unbanned first.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
func BanChatMember(ctx context.Context, chatId int64, userId int64, opts ...*OptBanChatMember) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId         int64 `json:"chat_id"`
		UserId         int64 `json:"user_id"`
		UntilDate      int64 `json:"until_date,omitempty"`
		RevokeMessages bool  `json:"revoke_messages,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		UserId: userId,
	}
	for _, opt := range opts {
		if opt.UntilDate != 0 {
			request.UntilDate = opt.UntilDate
		}
		if opt.RevokeMessages {
			request.RevokeMessages = opt.RevokeMessages
		}
	}
	return GenericRequest[Request, bool](ctx, "banChatMember", request)
}

type OptBanChatMember struct {
	UntilDate      int64
	RevokeMessages bool
}

// BanChatSenderChat Use this method to ban a channel chat in a supergroup or a channel. Returns True on success.
// Until the chat is unbanned, the owner of the banned chat won't be able to send messages on behalf of any of their channels.
// The bot must be an administrator in the supergroup or channel for this to work and must have the appropriate administrator rights.
func BanChatSenderChat(ctx context.Context, chatId int64, senderChatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId       int64 `json:"chat_id"`
		SenderChatId int64 `json:"sender_chat_id"`
	}
	request := &Request{
		ChatId:       chatId,
		SenderChatId: senderChatId,
	}
	return GenericRequest[Request, bool](ctx, "banChatSenderChat", request)
}

// Close Use this method to close the bot instance before moving it from one local server to another.
// You need to delete the webhook before calling this method to ensure that the bot isn't launched again after server restart.
// The method will return error 429 in the first 10 minutes after the bot is launched.
// Returns True on success. Requires no parameters.
func Close(ctx context.Context) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
	}
	request := &Request{}
	return GenericRequest[Request, bool](ctx, "close", request)
}

// CloseForumTopic Use this method to close an open topic in a forum supergroup chat. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights, unless it is the creator of the topic.
func CloseForumTopic(ctx context.Context, chatId int64, messageThreadId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId          int64 `json:"chat_id"`
		MessageThreadId int64 `json:"message_thread_id"`
	}
	request := &Request{
		ChatId:          chatId,
		MessageThreadId: messageThreadId,
	}
	return GenericRequest[Request, bool](ctx, "closeForumTopic", request)
}

// CloseGeneralForumTopic Use this method to close an open 'General' topic in a forum supergroup chat. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights.
func CloseGeneralForumTopic(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "closeGeneralForumTopic", request)
}

// ConvertGiftToStars Converts a given regular gift to Telegram Stars. Requires the can_convert_gifts_to_stars business bot right.
// Returns True on success.
func ConvertGiftToStars(ctx context.Context, businessConnectionId string, ownedGiftId string) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		OwnedGiftId          string `json:"owned_gift_id"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		OwnedGiftId:          ownedGiftId,
	}
	return GenericRequest[Request, bool](ctx, "convertGiftToStars", request)
}

// CopyMessage Use this method to copy messages of any kind. Returns the MessageId of the sent message on success.
// Service messages, paid media messages, giveaway messages, giveaway winners messages, and invoice messages can't be copied.
// A quiz poll can be copied only if the value of the field correct_option_id is known to the bot.
// The method is analogous to the method forwardMessage, but the copied message doesn't have a link to the original message.
func CopyMessage(ctx context.Context, chatId int64, fromChatId int64, messageId int64, opts ...*OptCopyMessage) (*MessageId, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId                int64                                                                       `json:"chat_id"`
		MessageThreadId       int64                                                                       `json:"message_thread_id,omitempty"`
		FromChatId            int64                                                                       `json:"from_chat_id"`
		MessageId             int64                                                                       `json:"message_id"`
		VideoStartTimestamp   int64                                                                       `json:"video_start_timestamp,omitempty"`
		Caption               string                                                                      `json:"caption,omitempty"`
		ParseMode             string                                                                      `json:"parse_mode,omitempty"`
		CaptionEntities       []*MessageEntity                                                            `json:"caption_entities,omitempty"`
		ShowCaptionAboveMedia bool                                                                        `json:"show_caption_above_media,omitempty"`
		DisableNotification   bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent        bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast    bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		ReplyParameters       *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageId:  messageId,
	}
	for _, opt := range opts {
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.VideoStartTimestamp != 0 {
			request.VideoStartTimestamp = opt.VideoStartTimestamp
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.ShowCaptionAboveMedia {
			request.ShowCaptionAboveMedia = opt.ShowCaptionAboveMedia
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *MessageId](ctx, "copyMessage", request)
}

type OptCopyMessage struct {
	MessageThreadId       int64
	VideoStartTimestamp   int64
	Caption               string
	ParseMode             string
	CaptionEntities       []*MessageEntity
	ShowCaptionAboveMedia bool
	DisableNotification   bool
	ProtectContent        bool
	AllowPaidBroadcast    bool
	ReplyParameters       *ReplyParameters
	ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

type VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply interface {
	variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyInlineKeyboardMarkup() *InlineKeyboardMarkup
	variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardMarkup() *ReplyKeyboardMarkup
	variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardRemove() *ReplyKeyboardRemove
	variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyForceReply() *ForceReply
}

var (
	_ VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply = &InlineKeyboardMarkup{}
	_ VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply = &ReplyKeyboardMarkup{}
	_ VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply = &ReplyKeyboardRemove{}
	_ VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply = &ForceReply{}
)

func (impl *InlineKeyboardMarkup) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyInlineKeyboardMarkup() *InlineKeyboardMarkup {
	return impl
}
func (impl *InlineKeyboardMarkup) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardMarkup() *ReplyKeyboardMarkup {
	return nil
}
func (impl *InlineKeyboardMarkup) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardRemove() *ReplyKeyboardRemove {
	return nil
}
func (impl *InlineKeyboardMarkup) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyForceReply() *ForceReply {
	return nil
}

func (impl *ReplyKeyboardMarkup) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyInlineKeyboardMarkup() *InlineKeyboardMarkup {
	return nil
}
func (impl *ReplyKeyboardMarkup) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardMarkup() *ReplyKeyboardMarkup {
	return impl
}
func (impl *ReplyKeyboardMarkup) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardRemove() *ReplyKeyboardRemove {
	return nil
}
func (impl *ReplyKeyboardMarkup) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyForceReply() *ForceReply {
	return nil
}

func (impl *ReplyKeyboardRemove) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyInlineKeyboardMarkup() *InlineKeyboardMarkup {
	return nil
}
func (impl *ReplyKeyboardRemove) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardMarkup() *ReplyKeyboardMarkup {
	return nil
}
func (impl *ReplyKeyboardRemove) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardRemove() *ReplyKeyboardRemove {
	return impl
}
func (impl *ReplyKeyboardRemove) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyForceReply() *ForceReply {
	return nil
}

func (impl *ForceReply) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyInlineKeyboardMarkup() *InlineKeyboardMarkup {
	return nil
}
func (impl *ForceReply) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardMarkup() *ReplyKeyboardMarkup {
	return nil
}
func (impl *ForceReply) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyReplyKeyboardRemove() *ReplyKeyboardRemove {
	return nil
}
func (impl *ForceReply) variantinlinekeyboardmarkupreplykeyboardmarkupreplykeyboardremoveforcereplyForceReply() *ForceReply {
	return impl
}

// CopyMessages Use this method to copy messages of any kind. Album grouping is kept for copied messages.
// If some of the specified messages can't be found or copied, they are skipped.
// Service messages, paid media messages, giveaway messages, giveaway winners messages, and invoice messages can't be copied.
// A quiz poll can be copied only if the value of the field correct_option_id is known to the bot.
// The method is analogous to the method forwardMessages, but the copied messages don't have a link to the original message.
// On success, an array of MessageId of the sent messages is returned.
func CopyMessages(ctx context.Context, chatId int64, fromChatId int64, messageIds []int64, opts ...*OptCopyMessages) ([]*MessageId, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId              int64   `json:"chat_id"`
		MessageThreadId     int64   `json:"message_thread_id,omitempty"`
		FromChatId          int64   `json:"from_chat_id"`
		MessageIds          []int64 `json:"message_ids"`
		DisableNotification bool    `json:"disable_notification,omitempty"`
		ProtectContent      bool    `json:"protect_content,omitempty"`
		RemoveCaption       bool    `json:"remove_caption,omitempty"`
	}
	request := &Request{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageIds: messageIds,
	}
	for _, opt := range opts {
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.RemoveCaption {
			request.RemoveCaption = opt.RemoveCaption
		}
	}
	return GenericRequest[Request, []*MessageId](ctx, "copyMessages", request)
}

type OptCopyMessages struct {
	MessageThreadId     int64
	DisableNotification bool
	ProtectContent      bool
	RemoveCaption       bool
}

// CreateChatInviteLink Use this method to create an additional invite link for a chat. Returns the new invite link as ChatInviteLink object.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// The link can be revoked using the method revokeChatInviteLink.
func CreateChatInviteLink(ctx context.Context, chatId int64, opts ...*OptCreateChatInviteLink) (*ChatInviteLink, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId             int64  `json:"chat_id"`
		Name               string `json:"name,omitempty"`
		ExpireDate         int64  `json:"expire_date,omitempty"`
		MemberLimit        int64  `json:"member_limit,omitempty"`
		CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
	}
	for _, opt := range opts {
		if opt.Name != "" {
			request.Name = opt.Name
		}
		if opt.ExpireDate != 0 {
			request.ExpireDate = opt.ExpireDate
		}
		if opt.MemberLimit != 0 {
			request.MemberLimit = opt.MemberLimit
		}
		if opt.CreatesJoinRequest {
			request.CreatesJoinRequest = opt.CreatesJoinRequest
		}
	}
	return GenericRequest[Request, *ChatInviteLink](ctx, "createChatInviteLink", request)
}

type OptCreateChatInviteLink struct {
	Name               string
	ExpireDate         int64
	MemberLimit        int64
	CreatesJoinRequest bool
}

// CreateChatSubscriptionInviteLink Use this method to create a subscription invite link for a channel chat.
// The bot must have the can_invite_users administrator rights.
// The link can be edited using the method editChatSubscriptionInviteLink or revoked using the method revokeChatInviteLink.
// Returns the new invite link as a ChatInviteLink object.
func CreateChatSubscriptionInviteLink(ctx context.Context, chatId int64, subscriptionPeriod int64, subscriptionPrice int64, opts ...*OptCreateChatSubscriptionInviteLink) (*ChatInviteLink, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId             int64  `json:"chat_id"`
		Name               string `json:"name,omitempty"`
		SubscriptionPeriod int64  `json:"subscription_period"`
		SubscriptionPrice  int64  `json:"subscription_price"`
	}
	request := &Request{
		ChatId:             chatId,
		SubscriptionPeriod: subscriptionPeriod,
		SubscriptionPrice:  subscriptionPrice,
	}
	for _, opt := range opts {
		if opt.Name != "" {
			request.Name = opt.Name
		}
	}
	return GenericRequest[Request, *ChatInviteLink](ctx, "createChatSubscriptionInviteLink", request)
}

type OptCreateChatSubscriptionInviteLink struct {
	Name string
}

// CreateForumTopic Use this method to create a topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights.
// Returns information about the created topic as a ForumTopic object.
func CreateForumTopic(ctx context.Context, chatId int64, name string, opts ...*OptCreateForumTopic) (*ForumTopic, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId            int64  `json:"chat_id"`
		Name              string `json:"name"`
		IconColor         int64  `json:"icon_color,omitempty"`
		IconCustomEmojiId string `json:"icon_custom_emoji_id,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		Name:   name,
	}
	for _, opt := range opts {
		if opt.IconColor != 0 {
			request.IconColor = opt.IconColor
		}
		if opt.IconCustomEmojiId != "" {
			request.IconCustomEmojiId = opt.IconCustomEmojiId
		}
	}
	return GenericRequest[Request, *ForumTopic](ctx, "createForumTopic", request)
}

type OptCreateForumTopic struct {
	IconColor         int64
	IconCustomEmojiId string
}

// CreateInvoiceLink Use this method to create a link for an invoice. Returns the created invoice link as String on success.
func CreateInvoiceLink(ctx context.Context, title string, description string, payload string, currency string, prices []*LabeledPrice, opts ...*OptCreateInvoiceLink) (string, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId      string          `json:"business_connection_id,omitempty"`
		Title                     string          `json:"title"`
		Description               string          `json:"description"`
		Payload                   string          `json:"payload"`
		ProviderToken             string          `json:"provider_token,omitempty"`
		Currency                  string          `json:"currency"`
		Prices                    []*LabeledPrice `json:"prices"`
		SubscriptionPeriod        int64           `json:"subscription_period,omitempty"`
		MaxTipAmount              int64           `json:"max_tip_amount,omitempty"`
		SuggestedTipAmounts       []int64         `json:"suggested_tip_amounts,omitempty"`
		ProviderData              string          `json:"provider_data,omitempty"`
		PhotoUrl                  string          `json:"photo_url,omitempty"`
		PhotoSize                 int64           `json:"photo_size,omitempty"`
		PhotoWidth                int64           `json:"photo_width,omitempty"`
		PhotoHeight               int64           `json:"photo_height,omitempty"`
		NeedName                  bool            `json:"need_name,omitempty"`
		NeedPhoneNumber           bool            `json:"need_phone_number,omitempty"`
		NeedEmail                 bool            `json:"need_email,omitempty"`
		NeedShippingAddress       bool            `json:"need_shipping_address,omitempty"`
		SendPhoneNumberToProvider bool            `json:"send_phone_number_to_provider,omitempty"`
		SendEmailToProvider       bool            `json:"send_email_to_provider,omitempty"`
		IsFlexible                bool            `json:"is_flexible,omitempty"`
	}
	request := &Request{
		Title:       title,
		Description: description,
		Payload:     payload,
		Currency:    currency,
		Prices:      prices,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.ProviderToken != "" {
			request.ProviderToken = opt.ProviderToken
		}
		if opt.SubscriptionPeriod != 0 {
			request.SubscriptionPeriod = opt.SubscriptionPeriod
		}
		if opt.MaxTipAmount != 0 {
			request.MaxTipAmount = opt.MaxTipAmount
		}
		if opt.SuggestedTipAmounts != nil {
			request.SuggestedTipAmounts = opt.SuggestedTipAmounts
		}
		if opt.ProviderData != "" {
			request.ProviderData = opt.ProviderData
		}
		if opt.PhotoUrl != "" {
			request.PhotoUrl = opt.PhotoUrl
		}
		if opt.PhotoSize != 0 {
			request.PhotoSize = opt.PhotoSize
		}
		if opt.PhotoWidth != 0 {
			request.PhotoWidth = opt.PhotoWidth
		}
		if opt.PhotoHeight != 0 {
			request.PhotoHeight = opt.PhotoHeight
		}
		if opt.NeedName {
			request.NeedName = opt.NeedName
		}
		if opt.NeedPhoneNumber {
			request.NeedPhoneNumber = opt.NeedPhoneNumber
		}
		if opt.NeedEmail {
			request.NeedEmail = opt.NeedEmail
		}
		if opt.NeedShippingAddress {
			request.NeedShippingAddress = opt.NeedShippingAddress
		}
		if opt.SendPhoneNumberToProvider {
			request.SendPhoneNumberToProvider = opt.SendPhoneNumberToProvider
		}
		if opt.SendEmailToProvider {
			request.SendEmailToProvider = opt.SendEmailToProvider
		}
		if opt.IsFlexible {
			request.IsFlexible = opt.IsFlexible
		}
	}
	return GenericRequest[Request, string](ctx, "createInvoiceLink", request)
}

type OptCreateInvoiceLink struct {
	BusinessConnectionId      string
	ProviderToken             string
	SubscriptionPeriod        int64
	MaxTipAmount              int64
	SuggestedTipAmounts       []int64
	ProviderData              string
	PhotoUrl                  string
	PhotoSize                 int64
	PhotoWidth                int64
	PhotoHeight               int64
	NeedName                  bool
	NeedPhoneNumber           bool
	NeedEmail                 bool
	NeedShippingAddress       bool
	SendPhoneNumberToProvider bool
	SendEmailToProvider       bool
	IsFlexible                bool
}

// CreateNewStickerSet Use this method to create a new sticker set owned by a user. Returns True on success.
// The bot will be able to edit the sticker set thus created.
func CreateNewStickerSet(ctx context.Context, userId int64, name string, title string, stickers []*InputSticker, opts ...*OptCreateNewStickerSet) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId          int64           `json:"user_id"`
		Name            string          `json:"name"`
		Title           string          `json:"title"`
		Stickers        []*InputSticker `json:"stickers"`
		StickerType     string          `json:"sticker_type,omitempty"`
		NeedsRepainting bool            `json:"needs_repainting,omitempty"`
	}
	request := &Request{
		UserId:   userId,
		Name:     name,
		Title:    title,
		Stickers: stickers,
	}
	for _, opt := range opts {
		if opt.StickerType != "" {
			request.StickerType = opt.StickerType
		}
		if opt.NeedsRepainting {
			request.NeedsRepainting = opt.NeedsRepainting
		}
	}
	return GenericRequest[Request, bool](ctx, "createNewStickerSet", request)
}

type OptCreateNewStickerSet struct {
	StickerType     string
	NeedsRepainting bool
}

// DeclineChatJoinRequest Use this method to decline a chat join request. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_invite_users administrator right.
func DeclineChatJoinRequest(ctx context.Context, chatId int64, userId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
		UserId int64 `json:"user_id"`
	}
	request := &Request{
		ChatId: chatId,
		UserId: userId,
	}
	return GenericRequest[Request, bool](ctx, "declineChatJoinRequest", request)
}

// DeleteBusinessMessages Delete messages on behalf of a business account. Returns True on success.
// Requires the can_delete_outgoing_messages business bot right to delete messages sent by the bot itself, or the can_delete_all_messages business bot right to delete any message.
func DeleteBusinessMessages(ctx context.Context, businessConnectionId string, messageIds []int64) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string  `json:"business_connection_id"`
		MessageIds           []int64 `json:"message_ids"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		MessageIds:           messageIds,
	}
	return GenericRequest[Request, bool](ctx, "deleteBusinessMessages", request)
}

// DeleteChatPhoto Use this method to delete a chat photo. Photos can't be changed for private chats.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Returns True on success.
func DeleteChatPhoto(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "deleteChatPhoto", request)
}

// DeleteChatStickerSet Use this method to delete a group sticker set from a supergroup. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Use the field can_set_sticker_set optionally returned in getChat requests to check if the bot can use this method.
func DeleteChatStickerSet(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "deleteChatStickerSet", request)
}

// DeleteForumTopic Use this method to delete a forum topic along with all its messages in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have the can_delete_messages administrator rights.
// Returns True on success.
func DeleteForumTopic(ctx context.Context, chatId int64, messageThreadId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId          int64 `json:"chat_id"`
		MessageThreadId int64 `json:"message_thread_id"`
	}
	request := &Request{
		ChatId:          chatId,
		MessageThreadId: messageThreadId,
	}
	return GenericRequest[Request, bool](ctx, "deleteForumTopic", request)
}

// DeleteMessage Use this method to delete a message, including service messages, with the following limitations:
// - A message can only be deleted if it was sent less than 48 hours ago.
// - Service messages about a supergroup, channel, or forum topic creation can't be deleted.
// - A dice message in a private chat can only be deleted if it was sent more than 24 hours ago.
// - Bots can delete outgoing messages in private chats, groups, and supergroups.
// - Bots can delete incoming messages in private chats.
// - Bots granted can_post_messages permissions can delete outgoing messages in channels.
// - If the bot is an administrator of a group, it can delete any message there.
// - If the bot has can_delete_messages permission in a supergroup or a channel, it can delete any message there.
// Returns True on success.
func DeleteMessage(ctx context.Context, chatId int64, messageId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId    int64 `json:"chat_id"`
		MessageId int64 `json:"message_id"`
	}
	request := &Request{
		ChatId:    chatId,
		MessageId: messageId,
	}
	return GenericRequest[Request, bool](ctx, "deleteMessage", request)
}

// DeleteMessages Use this method to delete multiple messages simultaneously. Returns True on success.
// If some of the specified messages can't be found, they are skipped.
func DeleteMessages(ctx context.Context, chatId int64, messageIds []int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId     int64   `json:"chat_id"`
		MessageIds []int64 `json:"message_ids"`
	}
	request := &Request{
		ChatId:     chatId,
		MessageIds: messageIds,
	}
	return GenericRequest[Request, bool](ctx, "deleteMessages", request)
}

// DeleteMyCommands Use this method to delete the list of the bot's commands for the given scope and user language.
// After deletion, higher level commands will be shown to affected users. Returns True on success.
func DeleteMyCommands(ctx context.Context, opts ...*OptDeleteMyCommands) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Scope        BotCommandScope `json:"scope,omitempty"`
		LanguageCode string          `json:"language_code,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.Scope != nil {
			request.Scope = opt.Scope
		}
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, bool](ctx, "deleteMyCommands", request)
}

type OptDeleteMyCommands struct {
	Scope        BotCommandScope
	LanguageCode string
}

// DeleteStickerFromSet Use this method to delete a sticker from a set created by the bot. Returns True on success.
func DeleteStickerFromSet(ctx context.Context, sticker string) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Sticker string `json:"sticker"`
	}
	request := &Request{
		Sticker: sticker,
	}
	return GenericRequest[Request, bool](ctx, "deleteStickerFromSet", request)
}

// DeleteStickerSet Use this method to delete a sticker set that was created by the bot. Returns True on success.
func DeleteStickerSet(ctx context.Context, name string) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Name string `json:"name"`
	}
	request := &Request{
		Name: name,
	}
	return GenericRequest[Request, bool](ctx, "deleteStickerSet", request)
}

// DeleteStory Deletes a story previously posted by the bot on behalf of a managed business account.
// Requires the can_manage_stories business bot right. Returns True on success.
func DeleteStory(ctx context.Context, businessConnectionId string, storyId int64) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		StoryId              int64  `json:"story_id"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		StoryId:              storyId,
	}
	return GenericRequest[Request, bool](ctx, "deleteStory", request)
}

// DeleteWebhook Use this method to remove webhook integration if you decide to switch back to getUpdates.
// Returns True on success.
func DeleteWebhook(ctx context.Context, opts ...*OptDeleteWebhook) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.DropPendingUpdates {
			request.DropPendingUpdates = opt.DropPendingUpdates
		}
	}
	return GenericRequest[Request, bool](ctx, "deleteWebhook", request)
}

type OptDeleteWebhook struct {
	DropPendingUpdates bool
}

// EditChatInviteLink Use this method to edit a non-primary invite link created by the bot.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Returns the edited invite link as a ChatInviteLink object.
func EditChatInviteLink(ctx context.Context, chatId int64, inviteLink string, opts ...*OptEditChatInviteLink) (*ChatInviteLink, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId             int64  `json:"chat_id"`
		InviteLink         string `json:"invite_link"`
		Name               string `json:"name,omitempty"`
		ExpireDate         int64  `json:"expire_date,omitempty"`
		MemberLimit        int64  `json:"member_limit,omitempty"`
		CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
	}
	request := &Request{
		ChatId:     chatId,
		InviteLink: inviteLink,
	}
	for _, opt := range opts {
		if opt.Name != "" {
			request.Name = opt.Name
		}
		if opt.ExpireDate != 0 {
			request.ExpireDate = opt.ExpireDate
		}
		if opt.MemberLimit != 0 {
			request.MemberLimit = opt.MemberLimit
		}
		if opt.CreatesJoinRequest {
			request.CreatesJoinRequest = opt.CreatesJoinRequest
		}
	}
	return GenericRequest[Request, *ChatInviteLink](ctx, "editChatInviteLink", request)
}

type OptEditChatInviteLink struct {
	Name               string
	ExpireDate         int64
	MemberLimit        int64
	CreatesJoinRequest bool
}

// EditChatSubscriptionInviteLink Use this method to edit a subscription invite link created by the bot.
// The bot must have the can_invite_users administrator rights.
// Returns the edited invite link as a ChatInviteLink object.
func EditChatSubscriptionInviteLink(ctx context.Context, chatId int64, inviteLink string, opts ...*OptEditChatSubscriptionInviteLink) (*ChatInviteLink, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId     int64  `json:"chat_id"`
		InviteLink string `json:"invite_link"`
		Name       string `json:"name,omitempty"`
	}
	request := &Request{
		ChatId:     chatId,
		InviteLink: inviteLink,
	}
	for _, opt := range opts {
		if opt.Name != "" {
			request.Name = opt.Name
		}
	}
	return GenericRequest[Request, *ChatInviteLink](ctx, "editChatSubscriptionInviteLink", request)
}

type OptEditChatSubscriptionInviteLink struct {
	Name string
}

// EditForumTopic Use this method to edit name and icon of a topic in a forum supergroup chat. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights, unless it is the creator of the topic.
func EditForumTopic(ctx context.Context, chatId int64, messageThreadId int64, opts ...*OptEditForumTopic) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId            int64  `json:"chat_id"`
		MessageThreadId   int64  `json:"message_thread_id"`
		Name              string `json:"name,omitempty"`
		IconCustomEmojiId string `json:"icon_custom_emoji_id,omitempty"`
	}
	request := &Request{
		ChatId:          chatId,
		MessageThreadId: messageThreadId,
	}
	for _, opt := range opts {
		if opt.Name != "" {
			request.Name = opt.Name
		}
		if opt.IconCustomEmojiId != "" {
			request.IconCustomEmojiId = opt.IconCustomEmojiId
		}
	}
	return GenericRequest[Request, bool](ctx, "editForumTopic", request)
}

type OptEditForumTopic struct {
	Name              string
	IconCustomEmojiId string
}

// EditGeneralForumTopic Use this method to edit the name of the 'General' topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights.
// Returns True on success.
func EditGeneralForumTopic(ctx context.Context, chatId int64, name string) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64  `json:"chat_id"`
		Name   string `json:"name"`
	}
	request := &Request{
		ChatId: chatId,
		Name:   name,
	}
	return GenericRequest[Request, bool](ctx, "editGeneralForumTopic", request)
}

// EditMessageCaption Use this method to edit captions of messages.
// On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
// Note that business messages that were not sent by the bot and do not contain an inline keyboard can only be edited within 48 hours from the time they were sent.
func EditMessageCaption(ctx context.Context, opts ...*OptEditMessageCaption) (*Message, error) /* >> either: [bool] */ {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId  string                `json:"business_connection_id,omitempty"`
		ChatId                int64                 `json:"chat_id,omitempty"`
		MessageId             int64                 `json:"message_id,omitempty"`
		InlineMessageId       string                `json:"inline_message_id,omitempty"`
		Caption               string                `json:"caption,omitempty"`
		ParseMode             string                `json:"parse_mode,omitempty"`
		CaptionEntities       []*MessageEntity      `json:"caption_entities,omitempty"`
		ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
		ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
		if opt.InlineMessageId != "" {
			request.InlineMessageId = opt.InlineMessageId
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.ShowCaptionAboveMedia {
			request.ShowCaptionAboveMedia = opt.ShowCaptionAboveMedia
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "editMessageCaption", request)
}

type OptEditMessageCaption struct {
	BusinessConnectionId  string
	ChatId                int64
	MessageId             int64
	InlineMessageId       string
	Caption               string
	ParseMode             string
	CaptionEntities       []*MessageEntity
	ShowCaptionAboveMedia bool
	ReplyMarkup           *InlineKeyboardMarkup
}

// EditMessageLiveLocation Use this method to edit live location messages.
// A location can be edited until its live_period expires or editing is explicitly disabled by a call to stopMessageLiveLocation.
// On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
func EditMessageLiveLocation(ctx context.Context, latitude float64, longitude float64, opts ...*OptEditMessageLiveLocation) (*Message, error) /* >> either: [bool] */ {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string                `json:"business_connection_id,omitempty"`
		ChatId               int64                 `json:"chat_id,omitempty"`
		MessageId            int64                 `json:"message_id,omitempty"`
		InlineMessageId      string                `json:"inline_message_id,omitempty"`
		Latitude             float64               `json:"latitude"`
		Longitude            float64               `json:"longitude"`
		LivePeriod           int64                 `json:"live_period,omitempty"`
		HorizontalAccuracy   float64               `json:"horizontal_accuracy,omitempty"`
		Heading              int64                 `json:"heading,omitempty"`
		ProximityAlertRadius int64                 `json:"proximity_alert_radius,omitempty"`
		ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{
		Latitude:  latitude,
		Longitude: longitude,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
		if opt.InlineMessageId != "" {
			request.InlineMessageId = opt.InlineMessageId
		}
		if opt.LivePeriod != 0 {
			request.LivePeriod = opt.LivePeriod
		}
		if opt.HorizontalAccuracy != 0.0 {
			request.HorizontalAccuracy = opt.HorizontalAccuracy
		}
		if opt.Heading != 0 {
			request.Heading = opt.Heading
		}
		if opt.ProximityAlertRadius != 0 {
			request.ProximityAlertRadius = opt.ProximityAlertRadius
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "editMessageLiveLocation", request)
}

type OptEditMessageLiveLocation struct {
	BusinessConnectionId string
	ChatId               int64
	MessageId            int64
	InlineMessageId      string
	LivePeriod           int64
	HorizontalAccuracy   float64
	Heading              int64
	ProximityAlertRadius int64
	ReplyMarkup          *InlineKeyboardMarkup
}

// EditMessageMedia Use this method to edit animation, audio, document, photo, or video messages, or to add media to text messages.
// If a message is part of a message album, then it can be edited only to an audio for audio albums, only to a document for document albums and to a photo or a video otherwise.
// When an inline message is edited, a new file can't be uploaded; use a previously uploaded file via its file_id or specify a URL.
// On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
// Note that business messages that were not sent by the bot and do not contain an inline keyboard can only be edited within 48 hours from the time they were sent.
func EditMessageMedia(ctx context.Context, media InputMedia, opts ...*OptEditMessageMedia) (*Message, error) /* >> either: [bool] */ {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string                `json:"business_connection_id,omitempty"`
		ChatId               int64                 `json:"chat_id,omitempty"`
		MessageId            int64                 `json:"message_id,omitempty"`
		InlineMessageId      string                `json:"inline_message_id,omitempty"`
		Media                InputMedia            `json:"media"`
		ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{
		Media: media,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
		if opt.InlineMessageId != "" {
			request.InlineMessageId = opt.InlineMessageId
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "editMessageMedia", request)
}

type OptEditMessageMedia struct {
	BusinessConnectionId string
	ChatId               int64
	MessageId            int64
	InlineMessageId      string
	ReplyMarkup          *InlineKeyboardMarkup
}

// EditMessageReplyMarkup Use this method to edit only the reply markup of messages.
// On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
// Note that business messages that were not sent by the bot and do not contain an inline keyboard can only be edited within 48 hours from the time they were sent.
func EditMessageReplyMarkup(ctx context.Context, opts ...*OptEditMessageReplyMarkup) (*Message, error) /* >> either: [bool] */ {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string                `json:"business_connection_id,omitempty"`
		ChatId               int64                 `json:"chat_id,omitempty"`
		MessageId            int64                 `json:"message_id,omitempty"`
		InlineMessageId      string                `json:"inline_message_id,omitempty"`
		ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
		if opt.InlineMessageId != "" {
			request.InlineMessageId = opt.InlineMessageId
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "editMessageReplyMarkup", request)
}

type OptEditMessageReplyMarkup struct {
	BusinessConnectionId string
	ChatId               int64
	MessageId            int64
	InlineMessageId      string
	ReplyMarkup          *InlineKeyboardMarkup
}

// EditMessageText Use this method to edit text and game messages.
// On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
// Note that business messages that were not sent by the bot and do not contain an inline keyboard can only be edited within 48 hours from the time they were sent.
func EditMessageText(ctx context.Context, text string, opts ...*OptEditMessageText) (*Message, error) /* >> either: [bool] */ {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string                `json:"business_connection_id,omitempty"`
		ChatId               int64                 `json:"chat_id,omitempty"`
		MessageId            int64                 `json:"message_id,omitempty"`
		InlineMessageId      string                `json:"inline_message_id,omitempty"`
		Text                 string                `json:"text"`
		ParseMode            string                `json:"parse_mode,omitempty"`
		Entities             []*MessageEntity      `json:"entities,omitempty"`
		LinkPreviewOptions   *LinkPreviewOptions   `json:"link_preview_options,omitempty"`
		ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{
		Text: text,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
		if opt.InlineMessageId != "" {
			request.InlineMessageId = opt.InlineMessageId
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.Entities != nil {
			request.Entities = opt.Entities
		}
		if opt.LinkPreviewOptions != nil {
			request.LinkPreviewOptions = opt.LinkPreviewOptions
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "editMessageText", request)
}

type OptEditMessageText struct {
	BusinessConnectionId string
	ChatId               int64
	MessageId            int64
	InlineMessageId      string
	ParseMode            string
	Entities             []*MessageEntity
	LinkPreviewOptions   *LinkPreviewOptions
	ReplyMarkup          *InlineKeyboardMarkup
}

// EditStory Edits a story previously posted by the bot on behalf of a managed business account.
// Requires the can_manage_stories business bot right. Returns Story on success.
func EditStory(ctx context.Context, businessConnectionId string, storyId int64, content InputStoryContent, opts ...*OptEditStory) (*Story, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string            `json:"business_connection_id"`
		StoryId              int64             `json:"story_id"`
		Content              InputStoryContent `json:"content"`
		Caption              string            `json:"caption,omitempty"`
		ParseMode            string            `json:"parse_mode,omitempty"`
		CaptionEntities      []*MessageEntity  `json:"caption_entities,omitempty"`
		Areas                []*StoryArea      `json:"areas,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		StoryId:              storyId,
		Content:              content,
	}
	for _, opt := range opts {
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.Areas != nil {
			request.Areas = opt.Areas
		}
	}
	return GenericRequest[Request, *Story](ctx, "editStory", request)
}

type OptEditStory struct {
	Caption         string
	ParseMode       string
	CaptionEntities []*MessageEntity
	Areas           []*StoryArea
}

// EditUserStarSubscription Allows the bot to cancel or re-enable extension of a subscription paid in Telegram Stars.
// Returns True on success.
func EditUserStarSubscription(ctx context.Context, userId int64, telegramPaymentChargeId string, isCanceled bool) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId                  int64  `json:"user_id"`
		TelegramPaymentChargeId string `json:"telegram_payment_charge_id"`
		IsCanceled              bool   `json:"is_canceled"`
	}
	request := &Request{
		UserId:                  userId,
		TelegramPaymentChargeId: telegramPaymentChargeId,
		IsCanceled:              isCanceled,
	}
	return GenericRequest[Request, bool](ctx, "editUserStarSubscription", request)
}

// ExportChatInviteLink Use this method to generate a new primary invite link for a chat; any previously generated primary link is revoked.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Returns the new invite link as String on success.
func ExportChatInviteLink(ctx context.Context, chatId int64) (string, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, string](ctx, "exportChatInviteLink", request)
}

// ForwardMessage Use this method to forward messages of any kind. On success, the sent Message is returned.
// Service messages and messages with protected content can't be forwarded.
func ForwardMessage(ctx context.Context, chatId int64, fromChatId int64, messageId int64, opts ...*OptForwardMessage) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId              int64 `json:"chat_id"`
		MessageThreadId     int64 `json:"message_thread_id,omitempty"`
		FromChatId          int64 `json:"from_chat_id"`
		VideoStartTimestamp int64 `json:"video_start_timestamp,omitempty"`
		DisableNotification bool  `json:"disable_notification,omitempty"`
		ProtectContent      bool  `json:"protect_content,omitempty"`
		MessageId           int64 `json:"message_id"`
	}
	request := &Request{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageId:  messageId,
	}
	for _, opt := range opts {
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.VideoStartTimestamp != 0 {
			request.VideoStartTimestamp = opt.VideoStartTimestamp
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
	}
	return GenericRequest[Request, *Message](ctx, "forwardMessage", request)
}

type OptForwardMessage struct {
	MessageThreadId     int64
	VideoStartTimestamp int64
	DisableNotification bool
	ProtectContent      bool
}

// ForwardMessages Use this method to forward multiple messages of any kind. Album grouping is kept for forwarded messages.
// If some of the specified messages can't be found or forwarded, they are skipped.
// Service messages and messages with protected content can't be forwarded.
// On success, an array of MessageId of the sent messages is returned.
func ForwardMessages(ctx context.Context, chatId int64, fromChatId int64, messageIds []int64, opts ...*OptForwardMessages) ([]*MessageId, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId              int64   `json:"chat_id"`
		MessageThreadId     int64   `json:"message_thread_id,omitempty"`
		FromChatId          int64   `json:"from_chat_id"`
		MessageIds          []int64 `json:"message_ids"`
		DisableNotification bool    `json:"disable_notification,omitempty"`
		ProtectContent      bool    `json:"protect_content,omitempty"`
	}
	request := &Request{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageIds: messageIds,
	}
	for _, opt := range opts {
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
	}
	return GenericRequest[Request, []*MessageId](ctx, "forwardMessages", request)
}

type OptForwardMessages struct {
	MessageThreadId     int64
	DisableNotification bool
	ProtectContent      bool
}

// GetAvailableGifts Returns the list of gifts that can be sent by the bot to users and channel chats.
// Requires no parameters. Returns a Gifts object.
func GetAvailableGifts(ctx context.Context) (*Gifts, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
	}
	request := &Request{}
	return GenericRequest[Request, *Gifts](ctx, "getAvailableGifts", request)
}

// GetBusinessAccountGifts Returns the gifts received and owned by a managed business account. Returns OwnedGifts on success.
// Requires the can_view_gifts_and_stars business bot right.
func GetBusinessAccountGifts(ctx context.Context, businessConnectionId string, opts ...*OptGetBusinessAccountGifts) (*OwnedGifts, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		ExcludeUnsaved       bool   `json:"exclude_unsaved,omitempty"`
		ExcludeSaved         bool   `json:"exclude_saved,omitempty"`
		ExcludeUnlimited     bool   `json:"exclude_unlimited,omitempty"`
		ExcludeLimited       bool   `json:"exclude_limited,omitempty"`
		ExcludeUnique        bool   `json:"exclude_unique,omitempty"`
		SortByPrice          bool   `json:"sort_by_price,omitempty"`
		Offset               string `json:"offset,omitempty"`
		Limit                int64  `json:"limit,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
	}
	for _, opt := range opts {
		if opt.ExcludeUnsaved {
			request.ExcludeUnsaved = opt.ExcludeUnsaved
		}
		if opt.ExcludeSaved {
			request.ExcludeSaved = opt.ExcludeSaved
		}
		if opt.ExcludeUnlimited {
			request.ExcludeUnlimited = opt.ExcludeUnlimited
		}
		if opt.ExcludeLimited {
			request.ExcludeLimited = opt.ExcludeLimited
		}
		if opt.ExcludeUnique {
			request.ExcludeUnique = opt.ExcludeUnique
		}
		if opt.SortByPrice {
			request.SortByPrice = opt.SortByPrice
		}
		if opt.Offset != "" {
			request.Offset = opt.Offset
		}
		if opt.Limit != 0 {
			request.Limit = opt.Limit
		}
	}
	return GenericRequest[Request, *OwnedGifts](ctx, "getBusinessAccountGifts", request)
}

type OptGetBusinessAccountGifts struct {
	ExcludeUnsaved   bool
	ExcludeSaved     bool
	ExcludeUnlimited bool
	ExcludeLimited   bool
	ExcludeUnique    bool
	SortByPrice      bool
	Offset           string
	Limit            int64
}

// GetBusinessAccountStarBalance Returns the amount of Telegram Stars owned by a managed business account. Returns StarAmount on success.
// Requires the can_view_gifts_and_stars business bot right.
func GetBusinessAccountStarBalance(ctx context.Context, businessConnectionId string) (*StarAmount, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
	}
	return GenericRequest[Request, *StarAmount](ctx, "getBusinessAccountStarBalance", request)
}

// GetBusinessConnection Use this method to get information about the connection of the bot with a business account.
// Returns a BusinessConnection object on success.
func GetBusinessConnection(ctx context.Context, businessConnectionId string) (*BusinessConnection, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
	}
	return GenericRequest[Request, *BusinessConnection](ctx, "getBusinessConnection", request)
}

// GetChat Use this method to get up-to-date information about the chat. Returns a ChatFullInfo object on success.
func GetChat(ctx context.Context, chatId int64) (*ChatFullInfo, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, *ChatFullInfo](ctx, "getChat", request)
}

// GetChatAdministrators Use this method to get a list of administrators in a chat, which aren't bots. Returns an Array of ChatMember objects.
func GetChatAdministrators(ctx context.Context, chatId int64) ([]ChatMember, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, []ChatMember](ctx, "getChatAdministrators", request)
}

// GetChatMember Use this method to get information about a member of a chat. Returns a ChatMember object on success.
// The method is only guaranteed to work for other users if the bot is an administrator in the chat.
func GetChatMember(ctx context.Context, chatId int64, userId int64) (ChatMember, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
		UserId int64 `json:"user_id"`
	}
	request := &Request{
		ChatId: chatId,
		UserId: userId,
	}
	return GenericRequest[Request, ChatMember](ctx, "getChatMember", request)
}

// GetChatMemberCount Use this method to get the number of members in a chat. Returns Int on success.
func GetChatMemberCount(ctx context.Context, chatId int64) (int64, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, int64](ctx, "getChatMemberCount", request)
}

// GetChatMenuButton Use this method to get the current value of the bot's menu button in a private chat, or the default menu button.
// Returns MenuButton on success.
func GetChatMenuButton(ctx context.Context, opts ...*OptGetChatMenuButton) (MenuButton, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		ChatId int64 `json:"chat_id,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
	}
	return GenericRequest[Request, MenuButton](ctx, "getChatMenuButton", request)
}

type OptGetChatMenuButton struct {
	ChatId int64
}

// GetCustomEmojiStickers Use this method to get information about custom emoji stickers by their identifiers.
// Returns an Array of Sticker objects.
func GetCustomEmojiStickers(ctx context.Context, customEmojiIds []string) ([]*Sticker, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		CustomEmojiIds []string `json:"custom_emoji_ids"`
	}
	request := &Request{
		CustomEmojiIds: customEmojiIds,
	}
	return GenericRequest[Request, []*Sticker](ctx, "getCustomEmojiStickers", request)
}

// GetFile Use this method to get basic information about a file and prepare it for downloading. For the moment, bots can download files of up to 20MB in size. On success, a File object is returned. The file can then be downloaded via the link https://api.telegram.org/file/bot<token>/<file_path>, where <file_path> is taken from the response. It is guaranteed that the link will be valid for at least 1 hour. When the link expires, a new one can be requested by calling getFile again.
// Note: This function may not preserve the original file name and MIME type. You should save the file's MIME type and name (if available) when the File object is received.
func GetFile(ctx context.Context, fileId string) (*File, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		FileId string `json:"file_id"`
	}
	request := &Request{
		FileId: fileId,
	}
	return GenericRequest[Request, *File](ctx, "getFile", request)
}

// GetForumTopicIconStickers Use this method to get custom emoji stickers, which can be used as a forum topic icon by any user.
// Requires no parameters. Returns an Array of Sticker objects.
func GetForumTopicIconStickers(ctx context.Context) ([]*Sticker, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
	}
	request := &Request{}
	return GenericRequest[Request, []*Sticker](ctx, "getForumTopicIconStickers", request)
}

// GetGameHighScores Use this method to get data for high score tables. Returns an Array of GameHighScore objects.
// Will return the score of the specified user and several of their neighbors in a game.
func GetGameHighScores(ctx context.Context, userId int64, opts ...*OptGetGameHighScores) ([]*GameHighScore, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId          int64  `json:"user_id"`
		ChatId          int64  `json:"chat_id,omitempty"`
		MessageId       int64  `json:"message_id,omitempty"`
		InlineMessageId string `json:"inline_message_id,omitempty"`
	}
	request := &Request{
		UserId: userId,
	}
	for _, opt := range opts {
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
		if opt.InlineMessageId != "" {
			request.InlineMessageId = opt.InlineMessageId
		}
	}
	return GenericRequest[Request, []*GameHighScore](ctx, "getGameHighScores", request)
}

type OptGetGameHighScores struct {
	ChatId          int64
	MessageId       int64
	InlineMessageId string
}

// GetMe A simple method for testing your bot's authentication token. Requires no parameters.
// Returns basic information about the bot in form of a User object.
func GetMe(ctx context.Context) (*User, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
	}
	request := &Request{}
	return GenericRequest[Request, *User](ctx, "getMe", request)
}

// GetMyCommands Use this method to get the current list of the bot's commands for the given scope and user language.
// Returns an Array of BotCommand objects. If commands aren't set, an empty list is returned.
func GetMyCommands(ctx context.Context, opts ...*OptGetMyCommands) ([]*BotCommand, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Scope        BotCommandScope `json:"scope,omitempty"`
		LanguageCode string          `json:"language_code,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.Scope != nil {
			request.Scope = opt.Scope
		}
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, []*BotCommand](ctx, "getMyCommands", request)
}

type OptGetMyCommands struct {
	Scope        BotCommandScope
	LanguageCode string
}

// GetMyDefaultAdministratorRights Use this method to get the current default administrator rights of the bot.
// Returns ChatAdministratorRights on success.
func GetMyDefaultAdministratorRights(ctx context.Context, opts ...*OptGetMyDefaultAdministratorRights) (*ChatAdministratorRights, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		ForChannels bool `json:"for_channels,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.ForChannels {
			request.ForChannels = opt.ForChannels
		}
	}
	return GenericRequest[Request, *ChatAdministratorRights](ctx, "getMyDefaultAdministratorRights", request)
}

type OptGetMyDefaultAdministratorRights struct {
	ForChannels bool
}

// GetMyDescription Use this method to get the current bot description for the given user language. Returns BotDescription on success.
func GetMyDescription(ctx context.Context, opts ...*OptGetMyDescription) (*BotDescription, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		LanguageCode string `json:"language_code,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, *BotDescription](ctx, "getMyDescription", request)
}

type OptGetMyDescription struct {
	LanguageCode string
}

// GetMyName Use this method to get the current bot name for the given user language. Returns BotName on success.
func GetMyName(ctx context.Context, opts ...*OptGetMyName) (*BotName, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		LanguageCode string `json:"language_code,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, *BotName](ctx, "getMyName", request)
}

type OptGetMyName struct {
	LanguageCode string
}

// GetMyShortDescription Use this method to get the current bot short description for the given user language.
// Returns BotShortDescription on success.
func GetMyShortDescription(ctx context.Context, opts ...*OptGetMyShortDescription) (*BotShortDescription, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		LanguageCode string `json:"language_code,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, *BotShortDescription](ctx, "getMyShortDescription", request)
}

type OptGetMyShortDescription struct {
	LanguageCode string
}

// GetStarTransactions Returns the bot's Telegram Star transactions in chronological order. On success, returns a StarTransactions object.
func GetStarTransactions(ctx context.Context, opts ...*OptGetStarTransactions) (*StarTransactions, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Offset int64 `json:"offset,omitempty"`
		Limit  int64 `json:"limit,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.Offset != 0 {
			request.Offset = opt.Offset
		}
		if opt.Limit != 0 {
			request.Limit = opt.Limit
		}
	}
	return GenericRequest[Request, *StarTransactions](ctx, "getStarTransactions", request)
}

type OptGetStarTransactions struct {
	Offset int64
	Limit  int64
}

// GetStickerSet Use this method to get a sticker set. On success, a StickerSet object is returned.
func GetStickerSet(ctx context.Context, name string) (*StickerSet, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Name string `json:"name"`
	}
	request := &Request{
		Name: name,
	}
	return GenericRequest[Request, *StickerSet](ctx, "getStickerSet", request)
}

// GetUpdates Use this method to receive incoming updates using long polling (wiki). Returns an Array of Update objects.
func GetUpdates(ctx context.Context, opts ...*OptGetUpdates) ([]*Update, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Offset         int64    `json:"offset,omitempty"`
		Limit          int64    `json:"limit,omitempty"`
		Timeout        int64    `json:"timeout,omitempty"`
		AllowedUpdates []string `json:"allowed_updates,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.Offset != 0 {
			request.Offset = opt.Offset
		}
		if opt.Limit != 0 {
			request.Limit = opt.Limit
		}
		if opt.Timeout != 0 {
			request.Timeout = opt.Timeout
		}
		if opt.AllowedUpdates != nil {
			request.AllowedUpdates = opt.AllowedUpdates
		}
	}
	return GenericRequest[Request, []*Update](ctx, "getUpdates", request)
}

type OptGetUpdates struct {
	Offset         int64
	Limit          int64
	Timeout        int64
	AllowedUpdates []string
}

// GetUserChatBoosts Use this method to get the list of boosts added to a chat by a user. Requires administrator rights in the chat.
// Returns a UserChatBoosts object.
func GetUserChatBoosts(ctx context.Context, chatId int64, userId int64) (*UserChatBoosts, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
		UserId int64 `json:"user_id"`
	}
	request := &Request{
		ChatId: chatId,
		UserId: userId,
	}
	return GenericRequest[Request, *UserChatBoosts](ctx, "getUserChatBoosts", request)
}

// GetUserProfilePhotos Use this method to get a list of profile pictures for a user. Returns a UserProfilePhotos object.
func GetUserProfilePhotos(ctx context.Context, userId int64, opts ...*OptGetUserProfilePhotos) (*UserProfilePhotos, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId int64 `json:"user_id"`
		Offset int64 `json:"offset,omitempty"`
		Limit  int64 `json:"limit,omitempty"`
	}
	request := &Request{
		UserId: userId,
	}
	for _, opt := range opts {
		if opt.Offset != 0 {
			request.Offset = opt.Offset
		}
		if opt.Limit != 0 {
			request.Limit = opt.Limit
		}
	}
	return GenericRequest[Request, *UserProfilePhotos](ctx, "getUserProfilePhotos", request)
}

type OptGetUserProfilePhotos struct {
	Offset int64
	Limit  int64
}

// GetWebhookInfo Use this method to get current webhook status. Requires no parameters. On success, returns a WebhookInfo object.
// If the bot is using getUpdates, will return an object with the url field empty.
func GetWebhookInfo(ctx context.Context) (*WebhookInfo, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
	}
	request := &Request{}
	return GenericRequest[Request, *WebhookInfo](ctx, "getWebhookInfo", request)
}

// GiftPremiumSubscription Gifts a Telegram Premium subscription to the given user. Returns True on success.
func GiftPremiumSubscription(ctx context.Context, userId int64, monthCount int64, starCount int64, opts ...*OptGiftPremiumSubscription) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId        int64            `json:"user_id"`
		MonthCount    int64            `json:"month_count"`
		StarCount     int64            `json:"star_count"`
		Text          string           `json:"text,omitempty"`
		TextParseMode string           `json:"text_parse_mode,omitempty"`
		TextEntities  []*MessageEntity `json:"text_entities,omitempty"`
	}
	request := &Request{
		UserId:     userId,
		MonthCount: monthCount,
		StarCount:  starCount,
	}
	for _, opt := range opts {
		if opt.Text != "" {
			request.Text = opt.Text
		}
		if opt.TextParseMode != "" {
			request.TextParseMode = opt.TextParseMode
		}
		if opt.TextEntities != nil {
			request.TextEntities = opt.TextEntities
		}
	}
	return GenericRequest[Request, bool](ctx, "giftPremiumSubscription", request)
}

type OptGiftPremiumSubscription struct {
	Text          string
	TextParseMode string
	TextEntities  []*MessageEntity
}

// HideGeneralForumTopic Use this method to hide the 'General' topic in a forum supergroup chat. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights.
// The topic will be automatically closed if it was open.
func HideGeneralForumTopic(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "hideGeneralForumTopic", request)
}

// LeaveChat Use this method for your bot to leave a group, supergroup or channel. Returns True on success.
func LeaveChat(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "leaveChat", request)
}

// LogOut Use this method to log out from the cloud Bot API server before launching the bot locally.
// You must log out the bot before running it locally, otherwise there is no guarantee that the bot will receive updates.
// After a successful call, you can immediately log in on a local server, but will not be able to log in back to the cloud Bot API server for 10 minutes.
// Returns True on success. Requires no parameters.
func LogOut(ctx context.Context) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
	}
	request := &Request{}
	return GenericRequest[Request, bool](ctx, "logOut", request)
}

// PinChatMessage Use this method to add a message to the list of pinned messages in a chat. Returns True on success.
// If the chat is not a private chat, the bot must be an administrator in the chat for this to work and must have the 'can_pin_messages' administrator right in a supergroup or 'can_edit_messages' administrator right in a channel.
func PinChatMessage(ctx context.Context, chatId int64, messageId int64, opts ...*OptPinChatMessage) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id,omitempty"`
		ChatId               int64  `json:"chat_id"`
		MessageId            int64  `json:"message_id"`
		DisableNotification  bool   `json:"disable_notification,omitempty"`
	}
	request := &Request{
		ChatId:    chatId,
		MessageId: messageId,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
	}
	return GenericRequest[Request, bool](ctx, "pinChatMessage", request)
}

type OptPinChatMessage struct {
	BusinessConnectionId string
	DisableNotification  bool
}

// PostStory Posts a story on behalf of a managed business account. Requires the can_manage_stories business bot right.
// Returns Story on success.
func PostStory(ctx context.Context, businessConnectionId string, content InputStoryContent, activePeriod int64, opts ...*OptPostStory) (*Story, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string            `json:"business_connection_id"`
		Content              InputStoryContent `json:"content"`
		ActivePeriod         int64             `json:"active_period"`
		Caption              string            `json:"caption,omitempty"`
		ParseMode            string            `json:"parse_mode,omitempty"`
		CaptionEntities      []*MessageEntity  `json:"caption_entities,omitempty"`
		Areas                []*StoryArea      `json:"areas,omitempty"`
		PostToChatPage       bool              `json:"post_to_chat_page,omitempty"`
		ProtectContent       bool              `json:"protect_content,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		Content:              content,
		ActivePeriod:         activePeriod,
	}
	for _, opt := range opts {
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.Areas != nil {
			request.Areas = opt.Areas
		}
		if opt.PostToChatPage {
			request.PostToChatPage = opt.PostToChatPage
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
	}
	return GenericRequest[Request, *Story](ctx, "postStory", request)
}

type OptPostStory struct {
	Caption         string
	ParseMode       string
	CaptionEntities []*MessageEntity
	Areas           []*StoryArea
	PostToChatPage  bool
	ProtectContent  bool
}

// PromoteChatMember Use this method to promote or demote a user in a supergroup or a channel. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Pass False for all boolean parameters to demote a user.
func PromoteChatMember(ctx context.Context, chatId int64, userId int64, opts ...*OptPromoteChatMember) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId              int64 `json:"chat_id"`
		UserId              int64 `json:"user_id"`
		IsAnonymous         bool  `json:"is_anonymous,omitempty"`
		CanManageChat       bool  `json:"can_manage_chat,omitempty"`
		CanDeleteMessages   bool  `json:"can_delete_messages,omitempty"`
		CanManageVideoChats bool  `json:"can_manage_video_chats,omitempty"`
		CanRestrictMembers  bool  `json:"can_restrict_members,omitempty"`
		CanPromoteMembers   bool  `json:"can_promote_members,omitempty"`
		CanChangeInfo       bool  `json:"can_change_info,omitempty"`
		CanInviteUsers      bool  `json:"can_invite_users,omitempty"`
		CanPostStories      bool  `json:"can_post_stories,omitempty"`
		CanEditStories      bool  `json:"can_edit_stories,omitempty"`
		CanDeleteStories    bool  `json:"can_delete_stories,omitempty"`
		CanPostMessages     bool  `json:"can_post_messages,omitempty"`
		CanEditMessages     bool  `json:"can_edit_messages,omitempty"`
		CanPinMessages      bool  `json:"can_pin_messages,omitempty"`
		CanManageTopics     bool  `json:"can_manage_topics,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		UserId: userId,
	}
	for _, opt := range opts {
		if opt.IsAnonymous {
			request.IsAnonymous = opt.IsAnonymous
		}
		if opt.CanManageChat {
			request.CanManageChat = opt.CanManageChat
		}
		if opt.CanDeleteMessages {
			request.CanDeleteMessages = opt.CanDeleteMessages
		}
		if opt.CanManageVideoChats {
			request.CanManageVideoChats = opt.CanManageVideoChats
		}
		if opt.CanRestrictMembers {
			request.CanRestrictMembers = opt.CanRestrictMembers
		}
		if opt.CanPromoteMembers {
			request.CanPromoteMembers = opt.CanPromoteMembers
		}
		if opt.CanChangeInfo {
			request.CanChangeInfo = opt.CanChangeInfo
		}
		if opt.CanInviteUsers {
			request.CanInviteUsers = opt.CanInviteUsers
		}
		if opt.CanPostStories {
			request.CanPostStories = opt.CanPostStories
		}
		if opt.CanEditStories {
			request.CanEditStories = opt.CanEditStories
		}
		if opt.CanDeleteStories {
			request.CanDeleteStories = opt.CanDeleteStories
		}
		if opt.CanPostMessages {
			request.CanPostMessages = opt.CanPostMessages
		}
		if opt.CanEditMessages {
			request.CanEditMessages = opt.CanEditMessages
		}
		if opt.CanPinMessages {
			request.CanPinMessages = opt.CanPinMessages
		}
		if opt.CanManageTopics {
			request.CanManageTopics = opt.CanManageTopics
		}
	}
	return GenericRequest[Request, bool](ctx, "promoteChatMember", request)
}

type OptPromoteChatMember struct {
	IsAnonymous         bool
	CanManageChat       bool
	CanDeleteMessages   bool
	CanManageVideoChats bool
	CanRestrictMembers  bool
	CanPromoteMembers   bool
	CanChangeInfo       bool
	CanInviteUsers      bool
	CanPostStories      bool
	CanEditStories      bool
	CanDeleteStories    bool
	CanPostMessages     bool
	CanEditMessages     bool
	CanPinMessages      bool
	CanManageTopics     bool
}

// ReadBusinessMessage Marks incoming message as read on behalf of a business account. Requires the can_read_messages business bot right.
// Returns True on success.
func ReadBusinessMessage(ctx context.Context, businessConnectionId string, chatId int64, messageId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		ChatId               int64  `json:"chat_id"`
		MessageId            int64  `json:"message_id"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		ChatId:               chatId,
		MessageId:            messageId,
	}
	return GenericRequest[Request, bool](ctx, "readBusinessMessage", request)
}

// RefundStarPayment Refunds a successful payment in Telegram Stars. Returns True on success.
func RefundStarPayment(ctx context.Context, userId int64, telegramPaymentChargeId string) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId                  int64  `json:"user_id"`
		TelegramPaymentChargeId string `json:"telegram_payment_charge_id"`
	}
	request := &Request{
		UserId:                  userId,
		TelegramPaymentChargeId: telegramPaymentChargeId,
	}
	return GenericRequest[Request, bool](ctx, "refundStarPayment", request)
}

// RemoveBusinessAccountProfilePhoto Removes the current profile photo of a managed business account. Returns True on success.
// Requires the can_edit_profile_photo business bot right.
func RemoveBusinessAccountProfilePhoto(ctx context.Context, businessConnectionId string, opts ...*OptRemoveBusinessAccountProfilePhoto) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		IsPublic             bool   `json:"is_public,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
	}
	for _, opt := range opts {
		if opt.IsPublic {
			request.IsPublic = opt.IsPublic
		}
	}
	return GenericRequest[Request, bool](ctx, "removeBusinessAccountProfilePhoto", request)
}

type OptRemoveBusinessAccountProfilePhoto struct {
	IsPublic bool
}

// RemoveChatVerification Removes verification from a chat that is currently verified on behalf of the organization represented by the bot.
// Returns True on success.
func RemoveChatVerification(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "removeChatVerification", request)
}

// RemoveUserVerification Removes verification from a user who is currently verified on behalf of the organization represented by the bot.
// Returns True on success.
func RemoveUserVerification(ctx context.Context, userId int64) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId int64 `json:"user_id"`
	}
	request := &Request{
		UserId: userId,
	}
	return GenericRequest[Request, bool](ctx, "removeUserVerification", request)
}

// ReopenForumTopic Use this method to reopen a closed topic in a forum supergroup chat. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights, unless it is the creator of the topic.
func ReopenForumTopic(ctx context.Context, chatId int64, messageThreadId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId          int64 `json:"chat_id"`
		MessageThreadId int64 `json:"message_thread_id"`
	}
	request := &Request{
		ChatId:          chatId,
		MessageThreadId: messageThreadId,
	}
	return GenericRequest[Request, bool](ctx, "reopenForumTopic", request)
}

// ReopenGeneralForumTopic Use this method to reopen a closed 'General' topic in a forum supergroup chat. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights.
// The topic will be automatically unhidden if it was hidden.
func ReopenGeneralForumTopic(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "reopenGeneralForumTopic", request)
}

// ReplaceStickerInSet Use this method to replace an existing sticker in a sticker set with a new one. Returns True on success.
// The method is equivalent to calling deleteStickerFromSet, then addStickerToSet, then setStickerPositionInSet.
func ReplaceStickerInSet(ctx context.Context, userId int64, name string, oldSticker string, sticker *InputSticker) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId     int64         `json:"user_id"`
		Name       string        `json:"name"`
		OldSticker string        `json:"old_sticker"`
		Sticker    *InputSticker `json:"sticker"`
	}
	request := &Request{
		UserId:     userId,
		Name:       name,
		OldSticker: oldSticker,
		Sticker:    sticker,
	}
	return GenericRequest[Request, bool](ctx, "replaceStickerInSet", request)
}

// RestrictChatMember Use this method to restrict a user in a supergroup. Pass True for all permissions to lift restrictions from a user.
// The bot must be an administrator in the supergroup for this to work and must have the appropriate administrator rights.
// Returns True on success.
func RestrictChatMember(ctx context.Context, chatId int64, userId int64, permissions *ChatPermissions, opts ...*OptRestrictChatMember) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId                        int64            `json:"chat_id"`
		UserId                        int64            `json:"user_id"`
		Permissions                   *ChatPermissions `json:"permissions"`
		UseIndependentChatPermissions bool             `json:"use_independent_chat_permissions,omitempty"`
		UntilDate                     int64            `json:"until_date,omitempty"`
	}
	request := &Request{
		ChatId:      chatId,
		UserId:      userId,
		Permissions: permissions,
	}
	for _, opt := range opts {
		if opt.UseIndependentChatPermissions {
			request.UseIndependentChatPermissions = opt.UseIndependentChatPermissions
		}
		if opt.UntilDate != 0 {
			request.UntilDate = opt.UntilDate
		}
	}
	return GenericRequest[Request, bool](ctx, "restrictChatMember", request)
}

type OptRestrictChatMember struct {
	UseIndependentChatPermissions bool
	UntilDate                     int64
}

// RevokeChatInviteLink Use this method to revoke an invite link created by the bot.
// If the primary link is revoked, a new link is automatically generated.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Returns the revoked invite link as ChatInviteLink object.
func RevokeChatInviteLink(ctx context.Context, chatId int64, inviteLink string) (*ChatInviteLink, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId     int64  `json:"chat_id"`
		InviteLink string `json:"invite_link"`
	}
	request := &Request{
		ChatId:     chatId,
		InviteLink: inviteLink,
	}
	return GenericRequest[Request, *ChatInviteLink](ctx, "revokeChatInviteLink", request)
}

// SavePreparedInlineMessage Stores a message that can be sent by a user of a Mini App. Returns a PreparedInlineMessage object.
func SavePreparedInlineMessage(ctx context.Context, userId int64, result InlineQueryResult, opts ...*OptSavePreparedInlineMessage) (*PreparedInlineMessage, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId            int64             `json:"user_id"`
		Result            InlineQueryResult `json:"result"`
		AllowUserChats    bool              `json:"allow_user_chats,omitempty"`
		AllowBotChats     bool              `json:"allow_bot_chats,omitempty"`
		AllowGroupChats   bool              `json:"allow_group_chats,omitempty"`
		AllowChannelChats bool              `json:"allow_channel_chats,omitempty"`
	}
	request := &Request{
		UserId: userId,
		Result: result,
	}
	for _, opt := range opts {
		if opt.AllowUserChats {
			request.AllowUserChats = opt.AllowUserChats
		}
		if opt.AllowBotChats {
			request.AllowBotChats = opt.AllowBotChats
		}
		if opt.AllowGroupChats {
			request.AllowGroupChats = opt.AllowGroupChats
		}
		if opt.AllowChannelChats {
			request.AllowChannelChats = opt.AllowChannelChats
		}
	}
	return GenericRequest[Request, *PreparedInlineMessage](ctx, "savePreparedInlineMessage", request)
}

type OptSavePreparedInlineMessage struct {
	AllowUserChats    bool
	AllowBotChats     bool
	AllowGroupChats   bool
	AllowChannelChats bool
}

// SendAnimation Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound).
// On success, the sent Message is returned.
// Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.
func SendAnimation(ctx context.Context, chatId int64, animation InputFile, opts ...*OptSendAnimation) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId  string                                                                      `json:"business_connection_id,omitempty"`
		ChatId                int64                                                                       `json:"chat_id"`
		MessageThreadId       int64                                                                       `json:"message_thread_id,omitempty"`
		Animation             InputFile                                                                   `json:"animation"`
		Duration              int64                                                                       `json:"duration,omitempty"`
		Width                 int64                                                                       `json:"width,omitempty"`
		Height                int64                                                                       `json:"height,omitempty"`
		Thumbnail             InputFile                                                                   `json:"thumbnail,omitempty"`
		Caption               string                                                                      `json:"caption,omitempty"`
		ParseMode             string                                                                      `json:"parse_mode,omitempty"`
		CaptionEntities       []*MessageEntity                                                            `json:"caption_entities,omitempty"`
		ShowCaptionAboveMedia bool                                                                        `json:"show_caption_above_media,omitempty"`
		HasSpoiler            bool                                                                        `json:"has_spoiler,omitempty"`
		DisableNotification   bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent        bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast    bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId       string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters       *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:    chatId,
		Animation: animation,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Duration != 0 {
			request.Duration = opt.Duration
		}
		if opt.Width != 0 {
			request.Width = opt.Width
		}
		if opt.Height != 0 {
			request.Height = opt.Height
		}
		if opt.Thumbnail != nil {
			request.Thumbnail = opt.Thumbnail
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.ShowCaptionAboveMedia {
			request.ShowCaptionAboveMedia = opt.ShowCaptionAboveMedia
		}
		if opt.HasSpoiler {
			request.HasSpoiler = opt.HasSpoiler
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "sendAnimation", request)
}

type OptSendAnimation struct {
	BusinessConnectionId  string
	MessageThreadId       int64
	Duration              int64
	Width                 int64
	Height                int64
	Thumbnail             InputFile
	Caption               string
	ParseMode             string
	CaptionEntities       []*MessageEntity
	ShowCaptionAboveMedia bool
	HasSpoiler            bool
	DisableNotification   bool
	ProtectContent        bool
	AllowPaidBroadcast    bool
	MessageEffectId       string
	ReplyParameters       *ReplyParameters
	ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendAudio Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.
// For sending voice messages, use the sendVoice method instead.
func SendAudio(ctx context.Context, chatId int64, audio InputFile, opts ...*OptSendAudio) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		Audio                InputFile                                                                   `json:"audio"`
		Caption              string                                                                      `json:"caption,omitempty"`
		ParseMode            string                                                                      `json:"parse_mode,omitempty"`
		CaptionEntities      []*MessageEntity                                                            `json:"caption_entities,omitempty"`
		Duration             int64                                                                       `json:"duration,omitempty"`
		Performer            string                                                                      `json:"performer,omitempty"`
		Title                string                                                                      `json:"title,omitempty"`
		Thumbnail            InputFile                                                                   `json:"thumbnail,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		Audio:  audio,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.Duration != 0 {
			request.Duration = opt.Duration
		}
		if opt.Performer != "" {
			request.Performer = opt.Performer
		}
		if opt.Title != "" {
			request.Title = opt.Title
		}
		if opt.Thumbnail != nil {
			request.Thumbnail = opt.Thumbnail
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "sendAudio", request)
}

type OptSendAudio struct {
	BusinessConnectionId string
	MessageThreadId      int64
	Caption              string
	ParseMode            string
	CaptionEntities      []*MessageEntity
	Duration             int64
	Performer            string
	Title                string
	Thumbnail            InputFile
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendChatAction Use this method when you need to tell the user that something is happening on the bot's side. The status is set for 5 seconds or less (when a message arrives from your bot, Telegram clients clear its typing status). Returns True on success.
// We only recommend using this method when a response from the bot will take a noticeable amount of time to arrive.
func SendChatAction(ctx context.Context, chatId int64, action string, opts ...*OptSendChatAction) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id,omitempty"`
		ChatId               int64  `json:"chat_id"`
		MessageThreadId      int64  `json:"message_thread_id,omitempty"`
		Action               string `json:"action"`
	}
	request := &Request{
		ChatId: chatId,
		Action: action,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
	}
	return GenericRequest[Request, bool](ctx, "sendChatAction", request)
}

type OptSendChatAction struct {
	BusinessConnectionId string
	MessageThreadId      int64
}

// SendContact Use this method to send phone contacts. On success, the sent Message is returned.
func SendContact(ctx context.Context, chatId int64, phoneNumber string, firstName string, opts ...*OptSendContact) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		PhoneNumber          string                                                                      `json:"phone_number"`
		FirstName            string                                                                      `json:"first_name"`
		LastName             string                                                                      `json:"last_name,omitempty"`
		Vcard                string                                                                      `json:"vcard,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:      chatId,
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.LastName != "" {
			request.LastName = opt.LastName
		}
		if opt.Vcard != "" {
			request.Vcard = opt.Vcard
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendContact", request)
}

type OptSendContact struct {
	BusinessConnectionId string
	MessageThreadId      int64
	LastName             string
	Vcard                string
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendDice Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned.
func SendDice(ctx context.Context, chatId int64, opts ...*OptSendDice) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		Emoji                string                                                                      `json:"emoji,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Emoji != "" {
			request.Emoji = opt.Emoji
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendDice", request)
}

type OptSendDice struct {
	BusinessConnectionId string
	MessageThreadId      int64
	Emoji                string
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendDocument Use this method to send general files. On success, the sent Message is returned.
// Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.
func SendDocument(ctx context.Context, chatId int64, document InputFile, opts ...*OptSendDocument) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId        string                                                                      `json:"business_connection_id,omitempty"`
		ChatId                      int64                                                                       `json:"chat_id"`
		MessageThreadId             int64                                                                       `json:"message_thread_id,omitempty"`
		Document                    InputFile                                                                   `json:"document"`
		Thumbnail                   InputFile                                                                   `json:"thumbnail,omitempty"`
		Caption                     string                                                                      `json:"caption,omitempty"`
		ParseMode                   string                                                                      `json:"parse_mode,omitempty"`
		CaptionEntities             []*MessageEntity                                                            `json:"caption_entities,omitempty"`
		DisableContentTypeDetection bool                                                                        `json:"disable_content_type_detection,omitempty"`
		DisableNotification         bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent              bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast          bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId             string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters             *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup                 VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:   chatId,
		Document: document,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Thumbnail != nil {
			request.Thumbnail = opt.Thumbnail
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.DisableContentTypeDetection {
			request.DisableContentTypeDetection = opt.DisableContentTypeDetection
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "sendDocument", request)
}

type OptSendDocument struct {
	BusinessConnectionId        string
	MessageThreadId             int64
	Thumbnail                   InputFile
	Caption                     string
	ParseMode                   string
	CaptionEntities             []*MessageEntity
	DisableContentTypeDetection bool
	DisableNotification         bool
	ProtectContent              bool
	AllowPaidBroadcast          bool
	MessageEffectId             string
	ReplyParameters             *ReplyParameters
	ReplyMarkup                 VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendGame Use this method to send a game. On success, the sent Message is returned.
func SendGame(ctx context.Context, chatId int64, gameShortName string, opts ...*OptSendGame) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                `json:"business_connection_id,omitempty"`
		ChatId               int64                 `json:"chat_id"`
		MessageThreadId      int64                 `json:"message_thread_id,omitempty"`
		GameShortName        string                `json:"game_short_name"`
		DisableNotification  bool                  `json:"disable_notification,omitempty"`
		ProtectContent       bool                  `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                  `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters      `json:"reply_parameters,omitempty"`
		ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:        chatId,
		GameShortName: gameShortName,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendGame", request)
}

type OptSendGame struct {
	BusinessConnectionId string
	MessageThreadId      int64
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          *InlineKeyboardMarkup
}

// SendGift Sends a gift to the given user or channel chat. The gift can't be converted to Telegram Stars by the receiver.
// Returns True on success.
func SendGift(ctx context.Context, giftId string, opts ...*OptSendGift) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId        int64            `json:"user_id,omitempty"`
		ChatId        int64            `json:"chat_id,omitempty"`
		GiftId        string           `json:"gift_id"`
		PayForUpgrade bool             `json:"pay_for_upgrade,omitempty"`
		Text          string           `json:"text,omitempty"`
		TextParseMode string           `json:"text_parse_mode,omitempty"`
		TextEntities  []*MessageEntity `json:"text_entities,omitempty"`
	}
	request := &Request{
		GiftId: giftId,
	}
	for _, opt := range opts {
		if opt.UserId != 0 {
			request.UserId = opt.UserId
		}
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.PayForUpgrade {
			request.PayForUpgrade = opt.PayForUpgrade
		}
		if opt.Text != "" {
			request.Text = opt.Text
		}
		if opt.TextParseMode != "" {
			request.TextParseMode = opt.TextParseMode
		}
		if opt.TextEntities != nil {
			request.TextEntities = opt.TextEntities
		}
	}
	return GenericRequest[Request, bool](ctx, "sendGift", request)
}

type OptSendGift struct {
	UserId        int64
	ChatId        int64
	PayForUpgrade bool
	Text          string
	TextParseMode string
	TextEntities  []*MessageEntity
}

// SendInvoice Use this method to send invoices. On success, the sent Message is returned.
func SendInvoice(ctx context.Context, chatId int64, title string, description string, payload string, currency string, prices []*LabeledPrice, opts ...*OptSendInvoice) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId                    int64                 `json:"chat_id"`
		MessageThreadId           int64                 `json:"message_thread_id,omitempty"`
		Title                     string                `json:"title"`
		Description               string                `json:"description"`
		Payload                   string                `json:"payload"`
		ProviderToken             string                `json:"provider_token,omitempty"`
		Currency                  string                `json:"currency"`
		Prices                    []*LabeledPrice       `json:"prices"`
		MaxTipAmount              int64                 `json:"max_tip_amount,omitempty"`
		SuggestedTipAmounts       []int64               `json:"suggested_tip_amounts,omitempty"`
		StartParameter            string                `json:"start_parameter,omitempty"`
		ProviderData              string                `json:"provider_data,omitempty"`
		PhotoUrl                  string                `json:"photo_url,omitempty"`
		PhotoSize                 int64                 `json:"photo_size,omitempty"`
		PhotoWidth                int64                 `json:"photo_width,omitempty"`
		PhotoHeight               int64                 `json:"photo_height,omitempty"`
		NeedName                  bool                  `json:"need_name,omitempty"`
		NeedPhoneNumber           bool                  `json:"need_phone_number,omitempty"`
		NeedEmail                 bool                  `json:"need_email,omitempty"`
		NeedShippingAddress       bool                  `json:"need_shipping_address,omitempty"`
		SendPhoneNumberToProvider bool                  `json:"send_phone_number_to_provider,omitempty"`
		SendEmailToProvider       bool                  `json:"send_email_to_provider,omitempty"`
		IsFlexible                bool                  `json:"is_flexible,omitempty"`
		DisableNotification       bool                  `json:"disable_notification,omitempty"`
		ProtectContent            bool                  `json:"protect_content,omitempty"`
		AllowPaidBroadcast        bool                  `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId           string                `json:"message_effect_id,omitempty"`
		ReplyParameters           *ReplyParameters      `json:"reply_parameters,omitempty"`
		ReplyMarkup               *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:      chatId,
		Title:       title,
		Description: description,
		Payload:     payload,
		Currency:    currency,
		Prices:      prices,
	}
	for _, opt := range opts {
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.ProviderToken != "" {
			request.ProviderToken = opt.ProviderToken
		}
		if opt.MaxTipAmount != 0 {
			request.MaxTipAmount = opt.MaxTipAmount
		}
		if opt.SuggestedTipAmounts != nil {
			request.SuggestedTipAmounts = opt.SuggestedTipAmounts
		}
		if opt.StartParameter != "" {
			request.StartParameter = opt.StartParameter
		}
		if opt.ProviderData != "" {
			request.ProviderData = opt.ProviderData
		}
		if opt.PhotoUrl != "" {
			request.PhotoUrl = opt.PhotoUrl
		}
		if opt.PhotoSize != 0 {
			request.PhotoSize = opt.PhotoSize
		}
		if opt.PhotoWidth != 0 {
			request.PhotoWidth = opt.PhotoWidth
		}
		if opt.PhotoHeight != 0 {
			request.PhotoHeight = opt.PhotoHeight
		}
		if opt.NeedName {
			request.NeedName = opt.NeedName
		}
		if opt.NeedPhoneNumber {
			request.NeedPhoneNumber = opt.NeedPhoneNumber
		}
		if opt.NeedEmail {
			request.NeedEmail = opt.NeedEmail
		}
		if opt.NeedShippingAddress {
			request.NeedShippingAddress = opt.NeedShippingAddress
		}
		if opt.SendPhoneNumberToProvider {
			request.SendPhoneNumberToProvider = opt.SendPhoneNumberToProvider
		}
		if opt.SendEmailToProvider {
			request.SendEmailToProvider = opt.SendEmailToProvider
		}
		if opt.IsFlexible {
			request.IsFlexible = opt.IsFlexible
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendInvoice", request)
}

type OptSendInvoice struct {
	MessageThreadId           int64
	ProviderToken             string
	MaxTipAmount              int64
	SuggestedTipAmounts       []int64
	StartParameter            string
	ProviderData              string
	PhotoUrl                  string
	PhotoSize                 int64
	PhotoWidth                int64
	PhotoHeight               int64
	NeedName                  bool
	NeedPhoneNumber           bool
	NeedEmail                 bool
	NeedShippingAddress       bool
	SendPhoneNumberToProvider bool
	SendEmailToProvider       bool
	IsFlexible                bool
	DisableNotification       bool
	ProtectContent            bool
	AllowPaidBroadcast        bool
	MessageEffectId           string
	ReplyParameters           *ReplyParameters
	ReplyMarkup               *InlineKeyboardMarkup
}

// SendLocation Use this method to send point on the map. On success, the sent Message is returned.
func SendLocation(ctx context.Context, chatId int64, latitude float64, longitude float64, opts ...*OptSendLocation) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		Latitude             float64                                                                     `json:"latitude"`
		Longitude            float64                                                                     `json:"longitude"`
		HorizontalAccuracy   float64                                                                     `json:"horizontal_accuracy,omitempty"`
		LivePeriod           int64                                                                       `json:"live_period,omitempty"`
		Heading              int64                                                                       `json:"heading,omitempty"`
		ProximityAlertRadius int64                                                                       `json:"proximity_alert_radius,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:    chatId,
		Latitude:  latitude,
		Longitude: longitude,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.HorizontalAccuracy != 0.0 {
			request.HorizontalAccuracy = opt.HorizontalAccuracy
		}
		if opt.LivePeriod != 0 {
			request.LivePeriod = opt.LivePeriod
		}
		if opt.Heading != 0 {
			request.Heading = opt.Heading
		}
		if opt.ProximityAlertRadius != 0 {
			request.ProximityAlertRadius = opt.ProximityAlertRadius
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendLocation", request)
}

type OptSendLocation struct {
	BusinessConnectionId string
	MessageThreadId      int64
	HorizontalAccuracy   float64
	LivePeriod           int64
	Heading              int64
	ProximityAlertRadius int64
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendMediaGroup Use this method to send a group of photos, videos, documents or audios as an album.
// Documents and audio files can be only grouped in an album with messages of the same type.
// On success, an array of Messages that were sent is returned.
func SendMediaGroup(ctx context.Context, chatId int64, media Album, opts ...*OptSendMediaGroup) ([]*Message, error) {
	schedule(ctx, chatId, len(media))
	defer scheduleDone(ctx, chatId, len(media))
	type Request struct {
		BusinessConnectionId string           `json:"business_connection_id,omitempty"`
		ChatId               int64            `json:"chat_id"`
		MessageThreadId      int64            `json:"message_thread_id,omitempty"`
		Media                Album            `json:"media"`
		DisableNotification  bool             `json:"disable_notification,omitempty"`
		ProtectContent       bool             `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool             `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string           `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters `json:"reply_parameters,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		Media:  media,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
	}
	return GenericRequestMultipart[Request, []*Message](ctx, "sendMediaGroup", request)
}

type OptSendMediaGroup struct {
	BusinessConnectionId string
	MessageThreadId      int64
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
}

// SendMessage Use this method to send text messages. On success, the sent Message is returned.
func SendMessage(ctx context.Context, chatId int64, text string, opts ...*OptSendMessage) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		Text                 string                                                                      `json:"text"`
		ParseMode            string                                                                      `json:"parse_mode,omitempty"`
		Entities             []*MessageEntity                                                            `json:"entities,omitempty"`
		LinkPreviewOptions   *LinkPreviewOptions                                                         `json:"link_preview_options,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		Text:   text,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.Entities != nil {
			request.Entities = opt.Entities
		}
		if opt.LinkPreviewOptions != nil {
			request.LinkPreviewOptions = opt.LinkPreviewOptions
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendMessage", request)
}

type OptSendMessage struct {
	BusinessConnectionId string
	MessageThreadId      int64
	ParseMode            string
	Entities             []*MessageEntity
	LinkPreviewOptions   *LinkPreviewOptions
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendPaidMedia Use this method to send paid media. On success, the sent Message is returned.
func SendPaidMedia(ctx context.Context, chatId int64, starCount int64, media []InputPaidMedia, opts ...*OptSendPaidMedia) (*Message, error) {
	schedule(ctx, chatId, len(media))
	defer scheduleDone(ctx, chatId, len(media))
	type Request struct {
		BusinessConnectionId  string                                                                      `json:"business_connection_id,omitempty"`
		ChatId                int64                                                                       `json:"chat_id"`
		StarCount             int64                                                                       `json:"star_count"`
		Media                 []InputPaidMedia                                                            `json:"media"`
		Payload               string                                                                      `json:"payload,omitempty"`
		Caption               string                                                                      `json:"caption,omitempty"`
		ParseMode             string                                                                      `json:"parse_mode,omitempty"`
		CaptionEntities       []*MessageEntity                                                            `json:"caption_entities,omitempty"`
		ShowCaptionAboveMedia bool                                                                        `json:"show_caption_above_media,omitempty"`
		DisableNotification   bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent        bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast    bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		ReplyParameters       *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:    chatId,
		StarCount: starCount,
		Media:     media,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.Payload != "" {
			request.Payload = opt.Payload
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.ShowCaptionAboveMedia {
			request.ShowCaptionAboveMedia = opt.ShowCaptionAboveMedia
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendPaidMedia", request)
}

type OptSendPaidMedia struct {
	BusinessConnectionId  string
	Payload               string
	Caption               string
	ParseMode             string
	CaptionEntities       []*MessageEntity
	ShowCaptionAboveMedia bool
	DisableNotification   bool
	ProtectContent        bool
	AllowPaidBroadcast    bool
	ReplyParameters       *ReplyParameters
	ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendPhoto Use this method to send photos. On success, the sent Message is returned.
func SendPhoto(ctx context.Context, chatId int64, photo InputFile, opts ...*OptSendPhoto) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId  string                                                                      `json:"business_connection_id,omitempty"`
		ChatId                int64                                                                       `json:"chat_id"`
		MessageThreadId       int64                                                                       `json:"message_thread_id,omitempty"`
		Photo                 InputFile                                                                   `json:"photo"`
		Caption               string                                                                      `json:"caption,omitempty"`
		ParseMode             string                                                                      `json:"parse_mode,omitempty"`
		CaptionEntities       []*MessageEntity                                                            `json:"caption_entities,omitempty"`
		ShowCaptionAboveMedia bool                                                                        `json:"show_caption_above_media,omitempty"`
		HasSpoiler            bool                                                                        `json:"has_spoiler,omitempty"`
		DisableNotification   bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent        bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast    bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId       string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters       *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		Photo:  photo,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.ShowCaptionAboveMedia {
			request.ShowCaptionAboveMedia = opt.ShowCaptionAboveMedia
		}
		if opt.HasSpoiler {
			request.HasSpoiler = opt.HasSpoiler
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "sendPhoto", request)
}

type OptSendPhoto struct {
	BusinessConnectionId  string
	MessageThreadId       int64
	Caption               string
	ParseMode             string
	CaptionEntities       []*MessageEntity
	ShowCaptionAboveMedia bool
	HasSpoiler            bool
	DisableNotification   bool
	ProtectContent        bool
	AllowPaidBroadcast    bool
	MessageEffectId       string
	ReplyParameters       *ReplyParameters
	ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendPoll Use this method to send a native poll. On success, the sent Message is returned.
func SendPoll(ctx context.Context, chatId int64, question string, options []*InputPollOption, opts ...*OptSendPoll) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId  string                                                                      `json:"business_connection_id,omitempty"`
		ChatId                int64                                                                       `json:"chat_id"`
		MessageThreadId       int64                                                                       `json:"message_thread_id,omitempty"`
		Question              string                                                                      `json:"question"`
		QuestionParseMode     string                                                                      `json:"question_parse_mode,omitempty"`
		QuestionEntities      []*MessageEntity                                                            `json:"question_entities,omitempty"`
		Options               []*InputPollOption                                                          `json:"options"`
		IsAnonymous           bool                                                                        `json:"is_anonymous,omitempty"`
		Type                  string                                                                      `json:"type,omitempty"`
		AllowsMultipleAnswers bool                                                                        `json:"allows_multiple_answers,omitempty"`
		CorrectOptionId       int64                                                                       `json:"correct_option_id,omitempty"`
		Explanation           string                                                                      `json:"explanation,omitempty"`
		ExplanationParseMode  string                                                                      `json:"explanation_parse_mode,omitempty"`
		ExplanationEntities   []*MessageEntity                                                            `json:"explanation_entities,omitempty"`
		OpenPeriod            int64                                                                       `json:"open_period,omitempty"`
		CloseDate             int64                                                                       `json:"close_date,omitempty"`
		IsClosed              bool                                                                        `json:"is_closed,omitempty"`
		DisableNotification   bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent        bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast    bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId       string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters       *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:   chatId,
		Question: question,
		Options:  options,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.QuestionParseMode != "" {
			request.QuestionParseMode = opt.QuestionParseMode
		}
		if opt.QuestionEntities != nil {
			request.QuestionEntities = opt.QuestionEntities
		}
		if opt.IsAnonymous {
			request.IsAnonymous = opt.IsAnonymous
		}
		if opt.Type != "" {
			request.Type = opt.Type
		}
		if opt.AllowsMultipleAnswers {
			request.AllowsMultipleAnswers = opt.AllowsMultipleAnswers
		}
		if opt.CorrectOptionId != 0 {
			request.CorrectOptionId = opt.CorrectOptionId
		}
		if opt.Explanation != "" {
			request.Explanation = opt.Explanation
		}
		if opt.ExplanationParseMode != "" {
			request.ExplanationParseMode = opt.ExplanationParseMode
		}
		if opt.ExplanationEntities != nil {
			request.ExplanationEntities = opt.ExplanationEntities
		}
		if opt.OpenPeriod != 0 {
			request.OpenPeriod = opt.OpenPeriod
		}
		if opt.CloseDate != 0 {
			request.CloseDate = opt.CloseDate
		}
		if opt.IsClosed {
			request.IsClosed = opt.IsClosed
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendPoll", request)
}

type OptSendPoll struct {
	BusinessConnectionId  string
	MessageThreadId       int64
	QuestionParseMode     string
	QuestionEntities      []*MessageEntity
	IsAnonymous           bool
	Type                  string
	AllowsMultipleAnswers bool
	CorrectOptionId       int64
	Explanation           string
	ExplanationParseMode  string
	ExplanationEntities   []*MessageEntity
	OpenPeriod            int64
	CloseDate             int64
	IsClosed              bool
	DisableNotification   bool
	ProtectContent        bool
	AllowPaidBroadcast    bool
	MessageEffectId       string
	ReplyParameters       *ReplyParameters
	ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendSticker Use this method to send static .WEBP, animated .TGS, or video .WEBM stickers.
// On success, the sent Message is returned.
func SendSticker(ctx context.Context, chatId int64, sticker InputFile, opts ...*OptSendSticker) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		Sticker              InputFile                                                                   `json:"sticker"`
		Emoji                string                                                                      `json:"emoji,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:  chatId,
		Sticker: sticker,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Emoji != "" {
			request.Emoji = opt.Emoji
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "sendSticker", request)
}

type OptSendSticker struct {
	BusinessConnectionId string
	MessageThreadId      int64
	Emoji                string
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendVenue Use this method to send information about a venue. On success, the sent Message is returned.
func SendVenue(ctx context.Context, chatId int64, latitude float64, longitude float64, title string, address string, opts ...*OptSendVenue) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		Latitude             float64                                                                     `json:"latitude"`
		Longitude            float64                                                                     `json:"longitude"`
		Title                string                                                                      `json:"title"`
		Address              string                                                                      `json:"address"`
		FoursquareId         string                                                                      `json:"foursquare_id,omitempty"`
		FoursquareType       string                                                                      `json:"foursquare_type,omitempty"`
		GooglePlaceId        string                                                                      `json:"google_place_id,omitempty"`
		GooglePlaceType      string                                                                      `json:"google_place_type,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:    chatId,
		Latitude:  latitude,
		Longitude: longitude,
		Title:     title,
		Address:   address,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.FoursquareId != "" {
			request.FoursquareId = opt.FoursquareId
		}
		if opt.FoursquareType != "" {
			request.FoursquareType = opt.FoursquareType
		}
		if opt.GooglePlaceId != "" {
			request.GooglePlaceId = opt.GooglePlaceId
		}
		if opt.GooglePlaceType != "" {
			request.GooglePlaceType = opt.GooglePlaceType
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "sendVenue", request)
}

type OptSendVenue struct {
	BusinessConnectionId string
	MessageThreadId      int64
	FoursquareId         string
	FoursquareType       string
	GooglePlaceId        string
	GooglePlaceType      string
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendVideo Use this method to send video files, Telegram clients support MPEG4 videos (other formats may be sent as Document).
// On success, the sent Message is returned.
// Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.
func SendVideo(ctx context.Context, chatId int64, video InputFile, opts ...*OptSendVideo) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId  string                                                                      `json:"business_connection_id,omitempty"`
		ChatId                int64                                                                       `json:"chat_id"`
		MessageThreadId       int64                                                                       `json:"message_thread_id,omitempty"`
		Video                 InputFile                                                                   `json:"video"`
		Duration              int64                                                                       `json:"duration,omitempty"`
		Width                 int64                                                                       `json:"width,omitempty"`
		Height                int64                                                                       `json:"height,omitempty"`
		Thumbnail             InputFile                                                                   `json:"thumbnail,omitempty"`
		Cover                 InputFile                                                                   `json:"cover,omitempty"`
		StartTimestamp        int64                                                                       `json:"start_timestamp,omitempty"`
		Caption               string                                                                      `json:"caption,omitempty"`
		ParseMode             string                                                                      `json:"parse_mode,omitempty"`
		CaptionEntities       []*MessageEntity                                                            `json:"caption_entities,omitempty"`
		ShowCaptionAboveMedia bool                                                                        `json:"show_caption_above_media,omitempty"`
		HasSpoiler            bool                                                                        `json:"has_spoiler,omitempty"`
		SupportsStreaming     bool                                                                        `json:"supports_streaming,omitempty"`
		DisableNotification   bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent        bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast    bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId       string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters       *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		Video:  video,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Duration != 0 {
			request.Duration = opt.Duration
		}
		if opt.Width != 0 {
			request.Width = opt.Width
		}
		if opt.Height != 0 {
			request.Height = opt.Height
		}
		if opt.Thumbnail != nil {
			request.Thumbnail = opt.Thumbnail
		}
		if opt.Cover != nil {
			request.Cover = opt.Cover
		}
		if opt.StartTimestamp != 0 {
			request.StartTimestamp = opt.StartTimestamp
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.ShowCaptionAboveMedia {
			request.ShowCaptionAboveMedia = opt.ShowCaptionAboveMedia
		}
		if opt.HasSpoiler {
			request.HasSpoiler = opt.HasSpoiler
		}
		if opt.SupportsStreaming {
			request.SupportsStreaming = opt.SupportsStreaming
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "sendVideo", request)
}

type OptSendVideo struct {
	BusinessConnectionId  string
	MessageThreadId       int64
	Duration              int64
	Width                 int64
	Height                int64
	Thumbnail             InputFile
	Cover                 InputFile
	StartTimestamp        int64
	Caption               string
	ParseMode             string
	CaptionEntities       []*MessageEntity
	ShowCaptionAboveMedia bool
	HasSpoiler            bool
	SupportsStreaming     bool
	DisableNotification   bool
	ProtectContent        bool
	AllowPaidBroadcast    bool
	MessageEffectId       string
	ReplyParameters       *ReplyParameters
	ReplyMarkup           VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendVideoNote As of v.4.0, Telegram clients support rounded square MPEG4 videos of up to 1 minute long.
// Use this method to send video messages. On success, the sent Message is returned.
func SendVideoNote(ctx context.Context, chatId int64, videoNote InputFile, opts ...*OptSendVideoNote) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		VideoNote            InputFile                                                                   `json:"video_note"`
		Duration             int64                                                                       `json:"duration,omitempty"`
		Length               int64                                                                       `json:"length,omitempty"`
		Thumbnail            InputFile                                                                   `json:"thumbnail,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:    chatId,
		VideoNote: videoNote,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Duration != 0 {
			request.Duration = opt.Duration
		}
		if opt.Length != 0 {
			request.Length = opt.Length
		}
		if opt.Thumbnail != nil {
			request.Thumbnail = opt.Thumbnail
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "sendVideoNote", request)
}

type OptSendVideoNote struct {
	BusinessConnectionId string
	MessageThreadId      int64
	Duration             int64
	Length               int64
	Thumbnail            InputFile
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SendVoice Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message.
// For this to work, your audio must be in an .OGG file encoded with OPUS, or in .MP3 format, or in .M4A format (other formats may be sent as Audio or Document).
// On success, the sent Message is returned.
// Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.
func SendVoice(ctx context.Context, chatId int64, voice InputFile, opts ...*OptSendVoice) (*Message, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                                                                      `json:"business_connection_id,omitempty"`
		ChatId               int64                                                                       `json:"chat_id"`
		MessageThreadId      int64                                                                       `json:"message_thread_id,omitempty"`
		Voice                InputFile                                                                   `json:"voice"`
		Caption              string                                                                      `json:"caption,omitempty"`
		ParseMode            string                                                                      `json:"parse_mode,omitempty"`
		CaptionEntities      []*MessageEntity                                                            `json:"caption_entities,omitempty"`
		Duration             int64                                                                       `json:"duration,omitempty"`
		DisableNotification  bool                                                                        `json:"disable_notification,omitempty"`
		ProtectContent       bool                                                                        `json:"protect_content,omitempty"`
		AllowPaidBroadcast   bool                                                                        `json:"allow_paid_broadcast,omitempty"`
		MessageEffectId      string                                                                      `json:"message_effect_id,omitempty"`
		ReplyParameters      *ReplyParameters                                                            `json:"reply_parameters,omitempty"`
		ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		Voice:  voice,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageThreadId != 0 {
			request.MessageThreadId = opt.MessageThreadId
		}
		if opt.Caption != "" {
			request.Caption = opt.Caption
		}
		if opt.ParseMode != "" {
			request.ParseMode = opt.ParseMode
		}
		if opt.CaptionEntities != nil {
			request.CaptionEntities = opt.CaptionEntities
		}
		if opt.Duration != 0 {
			request.Duration = opt.Duration
		}
		if opt.DisableNotification {
			request.DisableNotification = opt.DisableNotification
		}
		if opt.ProtectContent {
			request.ProtectContent = opt.ProtectContent
		}
		if opt.AllowPaidBroadcast {
			request.AllowPaidBroadcast = opt.AllowPaidBroadcast
		}
		if opt.MessageEffectId != "" {
			request.MessageEffectId = opt.MessageEffectId
		}
		if opt.ReplyParameters != nil {
			request.ReplyParameters = opt.ReplyParameters
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequestMultipart[Request, *Message](ctx, "sendVoice", request)
}

type OptSendVoice struct {
	BusinessConnectionId string
	MessageThreadId      int64
	Caption              string
	ParseMode            string
	CaptionEntities      []*MessageEntity
	Duration             int64
	DisableNotification  bool
	ProtectContent       bool
	AllowPaidBroadcast   bool
	MessageEffectId      string
	ReplyParameters      *ReplyParameters
	ReplyMarkup          VariantInlineKeyboardMarkupReplyKeyboardMarkupReplyKeyboardRemoveForceReply
}

// SetBusinessAccountBio Changes the bio of a managed business account. Requires the can_change_bio business bot right.
// Returns True on success.
func SetBusinessAccountBio(ctx context.Context, businessConnectionId string, opts ...*OptSetBusinessAccountBio) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		Bio                  string `json:"bio,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
	}
	for _, opt := range opts {
		if opt.Bio != "" {
			request.Bio = opt.Bio
		}
	}
	return GenericRequest[Request, bool](ctx, "setBusinessAccountBio", request)
}

type OptSetBusinessAccountBio struct {
	Bio string
}

// SetBusinessAccountGiftSettings Changes the privacy settings pertaining to incoming gifts in a managed business account.
// Requires the can_change_gift_settings business bot right. Returns True on success.
func SetBusinessAccountGiftSettings(ctx context.Context, businessConnectionId string, showGiftButton bool, acceptedGiftTypes *AcceptedGiftTypes) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string             `json:"business_connection_id"`
		ShowGiftButton       bool               `json:"show_gift_button"`
		AcceptedGiftTypes    *AcceptedGiftTypes `json:"accepted_gift_types"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		ShowGiftButton:       showGiftButton,
		AcceptedGiftTypes:    acceptedGiftTypes,
	}
	return GenericRequest[Request, bool](ctx, "setBusinessAccountGiftSettings", request)
}

// SetBusinessAccountName Changes the first and last name of a managed business account. Requires the can_change_name business bot right.
// Returns True on success.
func SetBusinessAccountName(ctx context.Context, businessConnectionId string, firstName string, opts ...*OptSetBusinessAccountName) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		FirstName            string `json:"first_name"`
		LastName             string `json:"last_name,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		FirstName:            firstName,
	}
	for _, opt := range opts {
		if opt.LastName != "" {
			request.LastName = opt.LastName
		}
	}
	return GenericRequest[Request, bool](ctx, "setBusinessAccountName", request)
}

type OptSetBusinessAccountName struct {
	LastName string
}

// SetBusinessAccountProfilePhoto Changes the profile photo of a managed business account. Requires the can_edit_profile_photo business bot right.
// Returns True on success.
func SetBusinessAccountProfilePhoto(ctx context.Context, businessConnectionId string, photo InputProfilePhoto, opts ...*OptSetBusinessAccountProfilePhoto) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string            `json:"business_connection_id"`
		Photo                InputProfilePhoto `json:"photo"`
		IsPublic             bool              `json:"is_public,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		Photo:                photo,
	}
	for _, opt := range opts {
		if opt.IsPublic {
			request.IsPublic = opt.IsPublic
		}
	}
	return GenericRequest[Request, bool](ctx, "setBusinessAccountProfilePhoto", request)
}

type OptSetBusinessAccountProfilePhoto struct {
	IsPublic bool
}

// SetBusinessAccountUsername Changes the username of a managed business account. Requires the can_change_username business bot right.
// Returns True on success.
func SetBusinessAccountUsername(ctx context.Context, businessConnectionId string, opts ...*OptSetBusinessAccountUsername) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		Username             string `json:"username,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
	}
	for _, opt := range opts {
		if opt.Username != "" {
			request.Username = opt.Username
		}
	}
	return GenericRequest[Request, bool](ctx, "setBusinessAccountUsername", request)
}

type OptSetBusinessAccountUsername struct {
	Username string
}

// SetChatAdministratorCustomTitle Use this method to set a custom title for an administrator in a supergroup promoted by the bot.
// Returns True on success.
func SetChatAdministratorCustomTitle(ctx context.Context, chatId int64, userId int64, customTitle string) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId      int64  `json:"chat_id"`
		UserId      int64  `json:"user_id"`
		CustomTitle string `json:"custom_title"`
	}
	request := &Request{
		ChatId:      chatId,
		UserId:      userId,
		CustomTitle: customTitle,
	}
	return GenericRequest[Request, bool](ctx, "setChatAdministratorCustomTitle", request)
}

// SetChatDescription Use this method to change the description of a group, a supergroup or a channel.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Returns True on success.
func SetChatDescription(ctx context.Context, chatId int64, opts ...*OptSetChatDescription) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId      int64  `json:"chat_id"`
		Description string `json:"description,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
	}
	for _, opt := range opts {
		if opt.Description != "" {
			request.Description = opt.Description
		}
	}
	return GenericRequest[Request, bool](ctx, "setChatDescription", request)
}

type OptSetChatDescription struct {
	Description string
}

// SetChatMenuButton Use this method to change the bot's menu button in a private chat, or the default menu button.
// Returns True on success.
func SetChatMenuButton(ctx context.Context, opts ...*OptSetChatMenuButton) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		ChatId     int64      `json:"chat_id,omitempty"`
		MenuButton MenuButton `json:"menu_button,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MenuButton != nil {
			request.MenuButton = opt.MenuButton
		}
	}
	return GenericRequest[Request, bool](ctx, "setChatMenuButton", request)
}

type OptSetChatMenuButton struct {
	ChatId     int64
	MenuButton MenuButton
}

// SetChatPermissions Use this method to set default chat permissions for all members. Returns True on success.
// The bot must be an administrator in the group or a supergroup for this to work and must have the can_restrict_members administrator rights.
func SetChatPermissions(ctx context.Context, chatId int64, permissions *ChatPermissions, opts ...*OptSetChatPermissions) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId                        int64            `json:"chat_id"`
		Permissions                   *ChatPermissions `json:"permissions"`
		UseIndependentChatPermissions bool             `json:"use_independent_chat_permissions,omitempty"`
	}
	request := &Request{
		ChatId:      chatId,
		Permissions: permissions,
	}
	for _, opt := range opts {
		if opt.UseIndependentChatPermissions {
			request.UseIndependentChatPermissions = opt.UseIndependentChatPermissions
		}
	}
	return GenericRequest[Request, bool](ctx, "setChatPermissions", request)
}

type OptSetChatPermissions struct {
	UseIndependentChatPermissions bool
}

// SetChatPhoto Use this method to set a new profile photo for the chat. Photos can't be changed for private chats.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Returns True on success.
func SetChatPhoto(ctx context.Context, chatId int64, photo *LocalFile) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64      `json:"chat_id"`
		Photo  *LocalFile `json:"photo"`
	}
	request := &Request{
		ChatId: chatId,
		Photo:  photo,
	}
	return GenericRequestMultipart[Request, bool](ctx, "setChatPhoto", request)
}

// SetChatStickerSet Use this method to set a new group sticker set for a supergroup. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Use the field can_set_sticker_set optionally returned in getChat requests to check if the bot can use this method.
func SetChatStickerSet(ctx context.Context, chatId int64, stickerSetName string) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId         int64  `json:"chat_id"`
		StickerSetName string `json:"sticker_set_name"`
	}
	request := &Request{
		ChatId:         chatId,
		StickerSetName: stickerSetName,
	}
	return GenericRequest[Request, bool](ctx, "setChatStickerSet", request)
}

// SetChatTitle Use this method to change the title of a chat. Titles can't be changed for private chats.
// The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.
// Returns True on success.
func SetChatTitle(ctx context.Context, chatId int64, title string) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64  `json:"chat_id"`
		Title  string `json:"title"`
	}
	request := &Request{
		ChatId: chatId,
		Title:  title,
	}
	return GenericRequest[Request, bool](ctx, "setChatTitle", request)
}

// SetCustomEmojiStickerSetThumbnail Use this method to set the thumbnail of a custom emoji sticker set. Returns True on success.
func SetCustomEmojiStickerSetThumbnail(ctx context.Context, name string, opts ...*OptSetCustomEmojiStickerSetThumbnail) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Name          string `json:"name"`
		CustomEmojiId string `json:"custom_emoji_id,omitempty"`
	}
	request := &Request{
		Name: name,
	}
	for _, opt := range opts {
		if opt.CustomEmojiId != "" {
			request.CustomEmojiId = opt.CustomEmojiId
		}
	}
	return GenericRequest[Request, bool](ctx, "setCustomEmojiStickerSetThumbnail", request)
}

type OptSetCustomEmojiStickerSetThumbnail struct {
	CustomEmojiId string
}

// SetGameScore Use this method to set the score of the specified user in a game message.
// On success, if the message is not an inline message, the Message is returned, otherwise True is returned.
// Returns an error, if the new score is not greater than the user's current score in the chat and force is False.
func SetGameScore(ctx context.Context, userId int64, score int64, opts ...*OptSetGameScore) (*Message, error) /* >> either: [bool] */ {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId             int64  `json:"user_id"`
		Score              int64  `json:"score"`
		Force              bool   `json:"force,omitempty"`
		DisableEditMessage bool   `json:"disable_edit_message,omitempty"`
		ChatId             int64  `json:"chat_id,omitempty"`
		MessageId          int64  `json:"message_id,omitempty"`
		InlineMessageId    string `json:"inline_message_id,omitempty"`
	}
	request := &Request{
		UserId: userId,
		Score:  score,
	}
	for _, opt := range opts {
		if opt.Force {
			request.Force = opt.Force
		}
		if opt.DisableEditMessage {
			request.DisableEditMessage = opt.DisableEditMessage
		}
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
		if opt.InlineMessageId != "" {
			request.InlineMessageId = opt.InlineMessageId
		}
	}
	return GenericRequest[Request, *Message](ctx, "setGameScore", request)
}

type OptSetGameScore struct {
	Force              bool
	DisableEditMessage bool
	ChatId             int64
	MessageId          int64
	InlineMessageId    string
}

// SetMessageReaction Use this method to change the chosen reactions on a message. Service messages of some types can't be reacted to.
// Automatically forwarded messages from a channel to its discussion group have the same available reactions as messages in the channel.
// Bots can't use paid reactions. Returns True on success.
func SetMessageReaction(ctx context.Context, chatId int64, messageId int64, opts ...*OptSetMessageReaction) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId    int64          `json:"chat_id"`
		MessageId int64          `json:"message_id"`
		Reaction  []ReactionType `json:"reaction,omitempty"`
		IsBig     bool           `json:"is_big,omitempty"`
	}
	request := &Request{
		ChatId:    chatId,
		MessageId: messageId,
	}
	for _, opt := range opts {
		if opt.Reaction != nil {
			request.Reaction = opt.Reaction
		}
		if opt.IsBig {
			request.IsBig = opt.IsBig
		}
	}
	return GenericRequest[Request, bool](ctx, "setMessageReaction", request)
}

type OptSetMessageReaction struct {
	Reaction []ReactionType
	IsBig    bool
}

// SetMyCommands Use this method to change the list of the bot's commands. See this manual for more details about bot commands.
// Returns True on success.
func SetMyCommands(ctx context.Context, commands []*BotCommand, opts ...*OptSetMyCommands) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Commands     []*BotCommand   `json:"commands"`
		Scope        BotCommandScope `json:"scope,omitempty"`
		LanguageCode string          `json:"language_code,omitempty"`
	}
	request := &Request{
		Commands: commands,
	}
	for _, opt := range opts {
		if opt.Scope != nil {
			request.Scope = opt.Scope
		}
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, bool](ctx, "setMyCommands", request)
}

type OptSetMyCommands struct {
	Scope        BotCommandScope
	LanguageCode string
}

// SetMyDefaultAdministratorRights Use this method to change the default administrator rights requested by the bot when it's added as an administrator to groups or channels.
// These rights will be suggested to users, but they are free to modify the list before adding the bot.
// Returns True on success.
func SetMyDefaultAdministratorRights(ctx context.Context, opts ...*OptSetMyDefaultAdministratorRights) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Rights      *ChatAdministratorRights `json:"rights,omitempty"`
		ForChannels bool                     `json:"for_channels,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.Rights != nil {
			request.Rights = opt.Rights
		}
		if opt.ForChannels {
			request.ForChannels = opt.ForChannels
		}
	}
	return GenericRequest[Request, bool](ctx, "setMyDefaultAdministratorRights", request)
}

type OptSetMyDefaultAdministratorRights struct {
	Rights      *ChatAdministratorRights
	ForChannels bool
}

// SetMyDescription Use this method to change the bot's description, which is shown in the chat with the bot if the chat is empty.
// Returns True on success.
func SetMyDescription(ctx context.Context, opts ...*OptSetMyDescription) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Description  string `json:"description,omitempty"`
		LanguageCode string `json:"language_code,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.Description != "" {
			request.Description = opt.Description
		}
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, bool](ctx, "setMyDescription", request)
}

type OptSetMyDescription struct {
	Description  string
	LanguageCode string
}

// SetMyName Use this method to change the bot's name. Returns True on success.
func SetMyName(ctx context.Context, opts ...*OptSetMyName) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Name         string `json:"name,omitempty"`
		LanguageCode string `json:"language_code,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.Name != "" {
			request.Name = opt.Name
		}
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, bool](ctx, "setMyName", request)
}

type OptSetMyName struct {
	Name         string
	LanguageCode string
}

// SetMyShortDescription Use this method to change the bot's short description, which is shown on the bot's profile page and is sent together with the link when users share the bot.
// Returns True on success.
func SetMyShortDescription(ctx context.Context, opts ...*OptSetMyShortDescription) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		ShortDescription string `json:"short_description,omitempty"`
		LanguageCode     string `json:"language_code,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.ShortDescription != "" {
			request.ShortDescription = opt.ShortDescription
		}
		if opt.LanguageCode != "" {
			request.LanguageCode = opt.LanguageCode
		}
	}
	return GenericRequest[Request, bool](ctx, "setMyShortDescription", request)
}

type OptSetMyShortDescription struct {
	ShortDescription string
	LanguageCode     string
}

// SetPassportDataErrors Informs a user that some of the Telegram Passport elements they provided contains errors. The user will not be able to re-submit their Passport to you until the errors are fixed (the contents of the field for which you returned the error must change). Returns True on success.
// Use this if the data submitted by the user doesn't satisfy the standards your service requires for any reason. For example, if a birthday date seems invalid, a submitted document is blurry, a scan shows evidence of tampering, etc. Supply some details in the error message to make sure the user knows how to correct the issues.
func SetPassportDataErrors(ctx context.Context, userId int64, errors []PassportElementError) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId int64                  `json:"user_id"`
		Errors []PassportElementError `json:"errors"`
	}
	request := &Request{
		UserId: userId,
		Errors: errors,
	}
	return GenericRequest[Request, bool](ctx, "setPassportDataErrors", request)
}

// SetStickerEmojiList Use this method to change the list of emoji assigned to a regular or custom emoji sticker.
// The sticker must belong to a sticker set created by the bot. Returns True on success.
func SetStickerEmojiList(ctx context.Context, sticker string, emojiList []string) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Sticker   string   `json:"sticker"`
		EmojiList []string `json:"emoji_list"`
	}
	request := &Request{
		Sticker:   sticker,
		EmojiList: emojiList,
	}
	return GenericRequest[Request, bool](ctx, "setStickerEmojiList", request)
}

// SetStickerKeywords Use this method to change search keywords assigned to a regular or custom emoji sticker.
// The sticker must belong to a sticker set created by the bot. Returns True on success.
func SetStickerKeywords(ctx context.Context, sticker string, opts ...*OptSetStickerKeywords) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Sticker  string   `json:"sticker"`
		Keywords []string `json:"keywords,omitempty"`
	}
	request := &Request{
		Sticker: sticker,
	}
	for _, opt := range opts {
		if opt.Keywords != nil {
			request.Keywords = opt.Keywords
		}
	}
	return GenericRequest[Request, bool](ctx, "setStickerKeywords", request)
}

type OptSetStickerKeywords struct {
	Keywords []string
}

// SetStickerMaskPosition Use this method to change the mask position of a mask sticker. Returns True on success.
// The sticker must belong to a sticker set that was created by the bot.
func SetStickerMaskPosition(ctx context.Context, sticker string, opts ...*OptSetStickerMaskPosition) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Sticker      string        `json:"sticker"`
		MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	}
	request := &Request{
		Sticker: sticker,
	}
	for _, opt := range opts {
		if opt.MaskPosition != nil {
			request.MaskPosition = opt.MaskPosition
		}
	}
	return GenericRequest[Request, bool](ctx, "setStickerMaskPosition", request)
}

type OptSetStickerMaskPosition struct {
	MaskPosition *MaskPosition
}

// SetStickerPositionInSet Use this method to move a sticker in a set created by the bot to a specific position.
// Returns True on success.
func SetStickerPositionInSet(ctx context.Context, sticker string, position int64) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Sticker  string `json:"sticker"`
		Position int64  `json:"position"`
	}
	request := &Request{
		Sticker:  sticker,
		Position: position,
	}
	return GenericRequest[Request, bool](ctx, "setStickerPositionInSet", request)
}

// SetStickerSetThumbnail Use this method to set the thumbnail of a regular or mask sticker set. Returns True on success.
// The format of the thumbnail file must match the format of the stickers in the set.
func SetStickerSetThumbnail(ctx context.Context, name string, userId int64, format string, opts ...*OptSetStickerSetThumbnail) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Name      string    `json:"name"`
		UserId    int64     `json:"user_id"`
		Thumbnail InputFile `json:"thumbnail,omitempty"`
		Format    string    `json:"format"`
	}
	request := &Request{
		Name:   name,
		UserId: userId,
		Format: format,
	}
	for _, opt := range opts {
		if opt.Thumbnail != nil {
			request.Thumbnail = opt.Thumbnail
		}
	}
	return GenericRequestMultipart[Request, bool](ctx, "setStickerSetThumbnail", request)
}

type OptSetStickerSetThumbnail struct {
	Thumbnail InputFile
}

// SetStickerSetTitle Use this method to set the title of a created sticker set. Returns True on success.
func SetStickerSetTitle(ctx context.Context, name string, title string) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	request := &Request{
		Name:  name,
		Title: title,
	}
	return GenericRequest[Request, bool](ctx, "setStickerSetTitle", request)
}

// SetUserEmojiStatus Changes the emoji status for a given user that previously allowed the bot to manage their emoji status via the Mini App method requestEmojiStatusAccess.
// Returns True on success.
func SetUserEmojiStatus(ctx context.Context, userId int64, opts ...*OptSetUserEmojiStatus) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId                    int64  `json:"user_id"`
		EmojiStatusCustomEmojiId  string `json:"emoji_status_custom_emoji_id,omitempty"`
		EmojiStatusExpirationDate int64  `json:"emoji_status_expiration_date,omitempty"`
	}
	request := &Request{
		UserId: userId,
	}
	for _, opt := range opts {
		if opt.EmojiStatusCustomEmojiId != "" {
			request.EmojiStatusCustomEmojiId = opt.EmojiStatusCustomEmojiId
		}
		if opt.EmojiStatusExpirationDate != 0 {
			request.EmojiStatusExpirationDate = opt.EmojiStatusExpirationDate
		}
	}
	return GenericRequest[Request, bool](ctx, "setUserEmojiStatus", request)
}

type OptSetUserEmojiStatus struct {
	EmojiStatusCustomEmojiId  string
	EmojiStatusExpirationDate int64
}

// SetWebhook Use this method to specify a URL and receive incoming updates via an outgoing webhook. Whenever there is an update for the bot, we will send an HTTPS POST request to the specified URL, containing a JSON-serialized Update. In case of an unsuccessful request (a request with response HTTP status code different from 2XY), we will repeat the request and give up after a reasonable amount of attempts. Returns True on success.
// If you'd like to make sure that the webhook was set by you, you can specify secret data in the parameter secret_token. If specified, the request will contain a header "X-Telegram-Bot-Api-Secret-Token" with the secret token as content.
func SetWebhook(ctx context.Context, url string, opts ...*OptSetWebhook) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		Url                string     `json:"url"`
		Certificate        *LocalFile `json:"certificate,omitempty"`
		IpAddress          string     `json:"ip_address,omitempty"`
		MaxConnections     int64      `json:"max_connections,omitempty"`
		AllowedUpdates     []string   `json:"allowed_updates,omitempty"`
		DropPendingUpdates bool       `json:"drop_pending_updates,omitempty"`
		SecretToken        string     `json:"secret_token,omitempty"`
	}
	request := &Request{
		Url: url,
	}
	for _, opt := range opts {
		if opt.Certificate != nil {
			request.Certificate = opt.Certificate
		}
		if opt.IpAddress != "" {
			request.IpAddress = opt.IpAddress
		}
		if opt.MaxConnections != 0 {
			request.MaxConnections = opt.MaxConnections
		}
		if opt.AllowedUpdates != nil {
			request.AllowedUpdates = opt.AllowedUpdates
		}
		if opt.DropPendingUpdates {
			request.DropPendingUpdates = opt.DropPendingUpdates
		}
		if opt.SecretToken != "" {
			request.SecretToken = opt.SecretToken
		}
	}
	return GenericRequestMultipart[Request, bool](ctx, "setWebhook", request)
}

type OptSetWebhook struct {
	Certificate        *LocalFile
	IpAddress          string
	MaxConnections     int64
	AllowedUpdates     []string
	DropPendingUpdates bool
	SecretToken        string
}

// StopMessageLiveLocation Use this method to stop updating a live location message before live_period expires.
// On success, if the message is not an inline message, the edited Message is returned, otherwise True is returned.
func StopMessageLiveLocation(ctx context.Context, opts ...*OptStopMessageLiveLocation) (*Message, error) /* >> either: [bool] */ {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string                `json:"business_connection_id,omitempty"`
		ChatId               int64                 `json:"chat_id,omitempty"`
		MessageId            int64                 `json:"message_id,omitempty"`
		InlineMessageId      string                `json:"inline_message_id,omitempty"`
		ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.ChatId != 0 {
			request.ChatId = opt.ChatId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
		if opt.InlineMessageId != "" {
			request.InlineMessageId = opt.InlineMessageId
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Message](ctx, "stopMessageLiveLocation", request)
}

type OptStopMessageLiveLocation struct {
	BusinessConnectionId string
	ChatId               int64
	MessageId            int64
	InlineMessageId      string
	ReplyMarkup          *InlineKeyboardMarkup
}

// StopPoll Use this method to stop a poll which was sent by the bot. On success, the stopped Poll is returned.
func StopPoll(ctx context.Context, chatId int64, messageId int64, opts ...*OptStopPoll) (*Poll, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string                `json:"business_connection_id,omitempty"`
		ChatId               int64                 `json:"chat_id"`
		MessageId            int64                 `json:"message_id"`
		ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}
	request := &Request{
		ChatId:    chatId,
		MessageId: messageId,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.ReplyMarkup != nil {
			request.ReplyMarkup = opt.ReplyMarkup
		}
	}
	return GenericRequest[Request, *Poll](ctx, "stopPoll", request)
}

type OptStopPoll struct {
	BusinessConnectionId string
	ReplyMarkup          *InlineKeyboardMarkup
}

// TransferBusinessAccountStars Transfers Telegram Stars from the business account balance to the bot's balance.
// Requires the can_transfer_stars business bot right. Returns True on success.
func TransferBusinessAccountStars(ctx context.Context, businessConnectionId string, starCount int64) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		StarCount            int64  `json:"star_count"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		StarCount:            starCount,
	}
	return GenericRequest[Request, bool](ctx, "transferBusinessAccountStars", request)
}

// TransferGift Transfers an owned unique gift to another user. Requires the can_transfer_and_upgrade_gifts business bot right.
// Requires can_transfer_stars business bot right if the transfer is paid. Returns True on success.
func TransferGift(ctx context.Context, businessConnectionId string, ownedGiftId string, newOwnerChatId int64, opts ...*OptTransferGift) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		OwnedGiftId          string `json:"owned_gift_id"`
		NewOwnerChatId       int64  `json:"new_owner_chat_id"`
		StarCount            int64  `json:"star_count,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		OwnedGiftId:          ownedGiftId,
		NewOwnerChatId:       newOwnerChatId,
	}
	for _, opt := range opts {
		if opt.StarCount != 0 {
			request.StarCount = opt.StarCount
		}
	}
	return GenericRequest[Request, bool](ctx, "transferGift", request)
}

type OptTransferGift struct {
	StarCount int64
}

// UnbanChatMember Use this method to unban a previously banned user in a supergroup or channel. Returns True on success.
// The user will not return to the group or channel automatically, but will be able to join via link, etc.
// The bot must be an administrator for this to work.
// By default, this method guarantees that after the call the user is not a member of the chat, but will be able to join it.
// So if the user is a member of the chat they will also be removed from the chat.
// If you don't want this, use the parameter only_if_banned.
func UnbanChatMember(ctx context.Context, chatId int64, userId int64, opts ...*OptUnbanChatMember) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId       int64 `json:"chat_id"`
		UserId       int64 `json:"user_id"`
		OnlyIfBanned bool  `json:"only_if_banned,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
		UserId: userId,
	}
	for _, opt := range opts {
		if opt.OnlyIfBanned {
			request.OnlyIfBanned = opt.OnlyIfBanned
		}
	}
	return GenericRequest[Request, bool](ctx, "unbanChatMember", request)
}

type OptUnbanChatMember struct {
	OnlyIfBanned bool
}

// UnbanChatSenderChat Use this method to unban a previously banned channel chat in a supergroup or channel.
// The bot must be an administrator for this to work and must have the appropriate administrator rights.
// Returns True on success.
func UnbanChatSenderChat(ctx context.Context, chatId int64, senderChatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId       int64 `json:"chat_id"`
		SenderChatId int64 `json:"sender_chat_id"`
	}
	request := &Request{
		ChatId:       chatId,
		SenderChatId: senderChatId,
	}
	return GenericRequest[Request, bool](ctx, "unbanChatSenderChat", request)
}

// UnhideGeneralForumTopic Use this method to unhide the 'General' topic in a forum supergroup chat. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights.
func UnhideGeneralForumTopic(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "unhideGeneralForumTopic", request)
}

// UnpinAllChatMessages Use this method to clear the list of pinned messages in a chat. Returns True on success.
// If the chat is not a private chat, the bot must be an administrator in the chat for this to work and must have the 'can_pin_messages' administrator right in a supergroup or 'can_edit_messages' administrator right in a channel.
func UnpinAllChatMessages(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "unpinAllChatMessages", request)
}

// UnpinAllForumTopicMessages Use this method to clear the list of pinned messages in a forum topic. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_pin_messages administrator right in the supergroup.
func UnpinAllForumTopicMessages(ctx context.Context, chatId int64, messageThreadId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId          int64 `json:"chat_id"`
		MessageThreadId int64 `json:"message_thread_id"`
	}
	request := &Request{
		ChatId:          chatId,
		MessageThreadId: messageThreadId,
	}
	return GenericRequest[Request, bool](ctx, "unpinAllForumTopicMessages", request)
}

// UnpinAllGeneralForumTopicMessages Use this method to clear the list of pinned messages in a General forum topic. Returns True on success.
// The bot must be an administrator in the chat for this to work and must have the can_pin_messages administrator right in the supergroup.
func UnpinAllGeneralForumTopicMessages(ctx context.Context, chatId int64) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId int64 `json:"chat_id"`
	}
	request := &Request{
		ChatId: chatId,
	}
	return GenericRequest[Request, bool](ctx, "unpinAllGeneralForumTopicMessages", request)
}

// UnpinChatMessage Use this method to remove a message from the list of pinned messages in a chat. Returns True on success.
// If the chat is not a private chat, the bot must be an administrator in the chat for this to work and must have the 'can_pin_messages' administrator right in a supergroup or 'can_edit_messages' administrator right in a channel.
func UnpinChatMessage(ctx context.Context, chatId int64, opts ...*OptUnpinChatMessage) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id,omitempty"`
		ChatId               int64  `json:"chat_id"`
		MessageId            int64  `json:"message_id,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
	}
	for _, opt := range opts {
		if opt.BusinessConnectionId != "" {
			request.BusinessConnectionId = opt.BusinessConnectionId
		}
		if opt.MessageId != 0 {
			request.MessageId = opt.MessageId
		}
	}
	return GenericRequest[Request, bool](ctx, "unpinChatMessage", request)
}

type OptUnpinChatMessage struct {
	BusinessConnectionId string
	MessageId            int64
}

// UpgradeGift Upgrades a given regular gift to a unique gift. Requires the can_transfer_and_upgrade_gifts business bot right.
// Additionally requires the can_transfer_stars business bot right if the upgrade is paid.
// Returns True on success.
func UpgradeGift(ctx context.Context, businessConnectionId string, ownedGiftId string, opts ...*OptUpgradeGift) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		BusinessConnectionId string `json:"business_connection_id"`
		OwnedGiftId          string `json:"owned_gift_id"`
		KeepOriginalDetails  bool   `json:"keep_original_details,omitempty"`
		StarCount            int64  `json:"star_count,omitempty"`
	}
	request := &Request{
		BusinessConnectionId: businessConnectionId,
		OwnedGiftId:          ownedGiftId,
	}
	for _, opt := range opts {
		if opt.KeepOriginalDetails {
			request.KeepOriginalDetails = opt.KeepOriginalDetails
		}
		if opt.StarCount != 0 {
			request.StarCount = opt.StarCount
		}
	}
	return GenericRequest[Request, bool](ctx, "upgradeGift", request)
}

type OptUpgradeGift struct {
	KeepOriginalDetails bool
	StarCount           int64
}

// UploadStickerFile Use this method to upload a file with a sticker for later use in the createNewStickerSet, addStickerToSet, or replaceStickerInSet methods (the file can be used multiple times).
// Returns the uploaded File on success.
func UploadStickerFile(ctx context.Context, userId int64, sticker *LocalFile, stickerFormat string) (*File, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId        int64      `json:"user_id"`
		Sticker       *LocalFile `json:"sticker"`
		StickerFormat string     `json:"sticker_format"`
	}
	request := &Request{
		UserId:        userId,
		Sticker:       sticker,
		StickerFormat: stickerFormat,
	}
	return GenericRequestMultipart[Request, *File](ctx, "uploadStickerFile", request)
}

// VerifyChat Verifies a chat on behalf of the organization which is represented by the bot. Returns True on success.
func VerifyChat(ctx context.Context, chatId int64, opts ...*OptVerifyChat) (bool, error) {
	schedule(ctx, chatId, 1)
	defer scheduleDone(ctx, chatId, 1)
	type Request struct {
		ChatId            int64  `json:"chat_id"`
		CustomDescription string `json:"custom_description,omitempty"`
	}
	request := &Request{
		ChatId: chatId,
	}
	for _, opt := range opts {
		if opt.CustomDescription != "" {
			request.CustomDescription = opt.CustomDescription
		}
	}
	return GenericRequest[Request, bool](ctx, "verifyChat", request)
}

type OptVerifyChat struct {
	CustomDescription string
}

// VerifyUser Verifies a user on behalf of the organization which is represented by the bot. Returns True on success.
func VerifyUser(ctx context.Context, userId int64, opts ...*OptVerifyUser) (bool, error) {
	schedule(ctx, 0, 1)
	defer scheduleDone(ctx, 0, 1)
	type Request struct {
		UserId            int64  `json:"user_id"`
		CustomDescription string `json:"custom_description,omitempty"`
	}
	request := &Request{
		UserId: userId,
	}
	for _, opt := range opts {
		if opt.CustomDescription != "" {
			request.CustomDescription = opt.CustomDescription
		}
	}
	return GenericRequest[Request, bool](ctx, "verifyUser", request)
}

type OptVerifyUser struct {
	CustomDescription string
}
