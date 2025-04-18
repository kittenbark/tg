package tg

import "encoding/json"

func (impl *BackgroundTypeFill) UnmarshalJSON(data []byte) error {
	type BackgroundFillUnmarshalJoinedFill struct {
		Type          *string  `json:"type"`
		Color         *int64   `json:"color"`
		TopColor      *int64   `json:"top_color"`
		BottomColor   *int64   `json:"bottom_color"`
		RotationAngle *int64   `json:"rotation_angle"`
		Colors        *[]int64 `json:"colors"`
	}
	type BaseInstance struct {
		// Type of the background, always "fill"
		Type string `json:"type"`
		// Dimming of the background in dark themes, as a percentage; 0-100
		DarkThemeDimming int64 `json:"dark_theme_dimming"`
		// Joint of structs, used for parsing variant interfaces.
		Fill *BackgroundFillUnmarshalJoinedFill `json:"fill"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.DarkThemeDimming = inst.DarkThemeDimming
	if inst.Fill != nil && inst.Fill.Type == nil {
		switch *inst.Fill.Type {
		case "freeform_gradient":
			impl.Fill = &BackgroundFillFreeformGradient{
				Type:   deref(inst.Fill.Type),
				Colors: deref(inst.Fill.Colors),
			}
		case "gradient":
			impl.Fill = &BackgroundFillGradient{
				Type:          deref(inst.Fill.Type),
				TopColor:      deref(inst.Fill.TopColor),
				BottomColor:   deref(inst.Fill.BottomColor),
				RotationAngle: deref(inst.Fill.RotationAngle),
			}
		case "solid":
			impl.Fill = &BackgroundFillSolid{
				Type:  deref(inst.Fill.Type),
				Color: deref(inst.Fill.Color),
			}
		}
	}
	return nil
}

func (impl *BackgroundTypePattern) UnmarshalJSON(data []byte) error {
	type BackgroundFillUnmarshalJoinedFill struct {
		Type          *string  `json:"type"`
		Color         *int64   `json:"color"`
		TopColor      *int64   `json:"top_color"`
		BottomColor   *int64   `json:"bottom_color"`
		RotationAngle *int64   `json:"rotation_angle"`
		Colors        *[]int64 `json:"colors"`
	}
	type BaseInstance struct {
		// Type of the background, always "pattern"
		Type string `json:"type"`
		// Document with the pattern
		Document *TelegramDocument `json:"document"`
		// Intensity of the pattern when it is shown above the filled background; 0-100
		Intensity int64 `json:"intensity"`
		// Optional. True, if the background fill must be applied only to the pattern itself.
		// All other pixels are black in this case. For dark themes only
		IsInverted bool `json:"is_inverted"`
		// Optional. True, if the background moves slightly when the device is tilted
		IsMoving bool `json:"is_moving"`
		// Joint of structs, used for parsing variant interfaces.
		Fill *BackgroundFillUnmarshalJoinedFill `json:"fill"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Document = inst.Document
	impl.Intensity = inst.Intensity
	impl.IsInverted = inst.IsInverted
	impl.IsMoving = inst.IsMoving
	if inst.Fill != nil && inst.Fill.Type == nil {
		switch *inst.Fill.Type {
		case "freeform_gradient":
			impl.Fill = &BackgroundFillFreeformGradient{
				Type:   deref(inst.Fill.Type),
				Colors: deref(inst.Fill.Colors),
			}
		case "gradient":
			impl.Fill = &BackgroundFillGradient{
				Type:          deref(inst.Fill.Type),
				TopColor:      deref(inst.Fill.TopColor),
				BottomColor:   deref(inst.Fill.BottomColor),
				RotationAngle: deref(inst.Fill.RotationAngle),
			}
		case "solid":
			impl.Fill = &BackgroundFillSolid{
				Type:  deref(inst.Fill.Type),
				Color: deref(inst.Fill.Color),
			}
		}
	}
	return nil
}

func (impl *ChatBackground) UnmarshalJSON(data []byte) error {
	type BackgroundTypeUnmarshalJoinedType struct {
		Type             *string            `json:"type"`
		Fill             *BackgroundFill    `json:"fill"`
		DarkThemeDimming *int64             `json:"dark_theme_dimming"`
		Document         **TelegramDocument `json:"document"`
		IsBlurred        *bool              `json:"is_blurred"`
		IsMoving         *bool              `json:"is_moving"`
		Intensity        *int64             `json:"intensity"`
		IsInverted       *bool              `json:"is_inverted"`
		ThemeName        *string            `json:"theme_name"`
	}
	type BaseInstance struct {

		// Joint of structs, used for parsing variant interfaces.
		Type *BackgroundTypeUnmarshalJoinedType `json:"type"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	if inst.Type != nil && inst.Type.Type == nil {
		switch *inst.Type.Type {
		case "pattern":
			impl.Type = &BackgroundTypePattern{
				Type:       deref(inst.Type.Type),
				Document:   deref(inst.Type.Document),
				Fill:       deref(inst.Type.Fill),
				Intensity:  deref(inst.Type.Intensity),
				IsInverted: deref(inst.Type.IsInverted),
				IsMoving:   deref(inst.Type.IsMoving),
			}
		case "wallpaper":
			impl.Type = &BackgroundTypeWallpaper{
				Type:             deref(inst.Type.Type),
				Document:         deref(inst.Type.Document),
				DarkThemeDimming: deref(inst.Type.DarkThemeDimming),
				IsBlurred:        deref(inst.Type.IsBlurred),
				IsMoving:         deref(inst.Type.IsMoving),
			}
		case "chat_theme":
			impl.Type = &BackgroundTypeChatTheme{
				Type:      deref(inst.Type.Type),
				ThemeName: deref(inst.Type.ThemeName),
			}
		case "fill":
			impl.Type = &BackgroundTypeFill{
				Type:             deref(inst.Type.Type),
				Fill:             deref(inst.Type.Fill),
				DarkThemeDimming: deref(inst.Type.DarkThemeDimming),
			}
		}
	}
	return nil
}

func (impl *ChatBoost) UnmarshalJSON(data []byte) error {
	type ChatBoostSourceUnmarshalJoinedSource struct {
		Source            *string `json:"source"`
		User              **User  `json:"user"`
		GiveawayMessageId *int64  `json:"giveaway_message_id"`
		PrizeStarCount    *int64  `json:"prize_star_count"`
		IsUnclaimed       *bool   `json:"is_unclaimed"`
	}
	type BaseInstance struct {
		// Unique identifier of the boost
		BoostId string `json:"boost_id"`
		// Point in time (Unix timestamp) when the chat was boosted
		AddDate int64 `json:"add_date"`
		// Point in time (Unix timestamp) when the boost will automatically expire, unless the booster's Telegram Premium subscription is prolonged
		ExpirationDate int64 `json:"expiration_date"`
		// Joint of structs, used for parsing variant interfaces.
		Source *ChatBoostSourceUnmarshalJoinedSource `json:"source"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.BoostId = inst.BoostId
	impl.AddDate = inst.AddDate
	impl.ExpirationDate = inst.ExpirationDate
	if inst.Source != nil && inst.Source.Source == nil {
		switch *inst.Source.Source {
		case "gift_code":
			impl.Source = &ChatBoostSourceGiftCode{
				Source: deref(inst.Source.Source),
				User:   deref(inst.Source.User),
			}
		case "giveaway":
			impl.Source = &ChatBoostSourceGiveaway{
				Source:            deref(inst.Source.Source),
				GiveawayMessageId: deref(inst.Source.GiveawayMessageId),
				User:              deref(inst.Source.User),
				PrizeStarCount:    deref(inst.Source.PrizeStarCount),
				IsUnclaimed:       deref(inst.Source.IsUnclaimed),
			}
		case "premium":
			impl.Source = &ChatBoostSourcePremium{
				Source: deref(inst.Source.Source),
				User:   deref(inst.Source.User),
			}
		}
	}
	return nil
}

func (impl *ChatBoostRemoved) UnmarshalJSON(data []byte) error {
	type ChatBoostSourceUnmarshalJoinedSource struct {
		Source            *string `json:"source"`
		User              **User  `json:"user"`
		GiveawayMessageId *int64  `json:"giveaway_message_id"`
		PrizeStarCount    *int64  `json:"prize_star_count"`
		IsUnclaimed       *bool   `json:"is_unclaimed"`
	}
	type BaseInstance struct {
		// Chat which was boosted
		Chat *Chat `json:"chat"`
		// Unique identifier of the boost
		BoostId string `json:"boost_id"`
		// Point in time (Unix timestamp) when the boost was removed
		RemoveDate int64 `json:"remove_date"`
		// Joint of structs, used for parsing variant interfaces.
		Source *ChatBoostSourceUnmarshalJoinedSource `json:"source"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Chat = inst.Chat
	impl.BoostId = inst.BoostId
	impl.RemoveDate = inst.RemoveDate
	if inst.Source != nil && inst.Source.Source == nil {
		switch *inst.Source.Source {
		case "gift_code":
			impl.Source = &ChatBoostSourceGiftCode{
				Source: deref(inst.Source.Source),
				User:   deref(inst.Source.User),
			}
		case "giveaway":
			impl.Source = &ChatBoostSourceGiveaway{
				Source:            deref(inst.Source.Source),
				GiveawayMessageId: deref(inst.Source.GiveawayMessageId),
				User:              deref(inst.Source.User),
				PrizeStarCount:    deref(inst.Source.PrizeStarCount),
				IsUnclaimed:       deref(inst.Source.IsUnclaimed),
			}
		case "premium":
			impl.Source = &ChatBoostSourcePremium{
				Source: deref(inst.Source.Source),
				User:   deref(inst.Source.User),
			}
		}
	}
	return nil
}

func (impl *ChatFullInfo) UnmarshalJSON(data []byte) error {
	type ReactionTypeUnmarshalJoinedAvailableReactions struct {
		Type          *string `json:"type"`
		Emoji         *string `json:"emoji"`
		CustomEmojiId *string `json:"custom_emoji_id"`
	}
	type BaseInstance struct {
		// Unique identifier for this chat.
		// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
		// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
		Id int64 `json:"id"`
		// Type of the chat, can be either "private", "group", "supergroup" or "channel"
		Type string `json:"type"`
		// Optional. Title, for supergroups, channels and group chats
		Title string `json:"title"`
		// Optional. Username, for private chats, supergroups and channels if available
		Username string `json:"username"`
		// Optional. First name of the other party in a private chat
		FirstName string `json:"first_name"`
		// Optional. Last name of the other party in a private chat
		LastName string `json:"last_name"`
		// Optional. True, if the supergroup chat is a forum (has topics enabled)
		IsForum bool `json:"is_forum"`
		// Identifier of the accent color for the chat name and backgrounds of the chat photo, reply header, and link preview.
		// See accent colors for more details.
		AccentColorId int64 `json:"accent_color_id"`
		// The maximum number of reactions that can be set on a message in the chat
		MaxReactionCount int64 `json:"max_reaction_count"`
		// Optional. Chat photo
		Photo *ChatPhoto `json:"photo"`
		// Optional. If non-empty, the list of all active chat usernames; for private chats, supergroups and channels
		ActiveUsernames []string `json:"active_usernames"`
		// Optional. For private chats, the date of birth of the user
		Birthdate *Birthdate `json:"birthdate"`
		// Optional. For private chats with business accounts, the intro of the business
		BusinessIntro *BusinessIntro `json:"business_intro"`
		// Optional. For private chats with business accounts, the location of the business
		BusinessLocation *BusinessLocation `json:"business_location"`
		// Optional. For private chats with business accounts, the opening hours of the business
		BusinessOpeningHours *BusinessOpeningHours `json:"business_opening_hours"`
		// Optional. For private chats, the personal channel of the user
		PersonalChat *Chat `json:"personal_chat"`
		// Optional. Custom emoji identifier of the emoji chosen by the chat for the reply header and link preview background
		BackgroundCustomEmojiId string `json:"background_custom_emoji_id"`
		// Optional. Identifier of the accent color for the chat's profile background.
		// See profile accent colors for more details.
		ProfileAccentColorId int64 `json:"profile_accent_color_id"`
		// Optional. Custom emoji identifier of the emoji chosen by the chat for its profile background
		ProfileBackgroundCustomEmojiId string `json:"profile_background_custom_emoji_id"`
		// Optional. Custom emoji identifier of the emoji status of the chat or the other party in a private chat
		EmojiStatusCustomEmojiId string `json:"emoji_status_custom_emoji_id"`
		// Optional. Expiration date of the emoji status of the chat or the other party in a private chat, in Unix time, if any
		EmojiStatusExpirationDate int64 `json:"emoji_status_expiration_date"`
		// Optional. Bio of the other party in a private chat
		Bio string `json:"bio"`
		// Optional.
		// True, if privacy settings of the other party in the private chat allows to use tg://user?id=<user_id> links only in chats with the user
		HasPrivateForwards bool `json:"has_private_forwards"`
		// Optional.
		// True, if the privacy settings of the other party restrict sending voice and video note messages in the private chat
		HasRestrictedVoiceAndVideoMessages bool `json:"has_restricted_voice_and_video_messages"`
		// Optional. True, if users need to join the supergroup before they can send messages
		JoinToSendMessages bool `json:"join_to_send_messages"`
		// Optional.
		// True, if all users directly joining the supergroup without using an invite link need to be approved by supergroup administrators
		JoinByRequest bool `json:"join_by_request"`
		// Optional. Description, for groups, supergroups and channel chats
		Description string `json:"description"`
		// Optional. Primary invite link, for groups, supergroups and channel chats
		InviteLink string `json:"invite_link"`
		// Optional. The most recent pinned message (by sending date)
		PinnedMessage *Message `json:"pinned_message"`
		// Optional. Default chat member permissions, for groups and supergroups
		Permissions *ChatPermissions `json:"permissions"`
		// Optional. True, if paid media messages can be sent or forwarded to the channel chat.
		// The field is available only for channel chats.
		CanSendPaidMedia bool `json:"can_send_paid_media"`
		// Optional.
		// For supergroups, the minimum allowed delay between consecutive messages sent by each unprivileged user; in seconds
		SlowModeDelay int64 `json:"slow_mode_delay"`
		// Optional.
		// For supergroups, the minimum number of boosts that a non-administrator user needs to add in order to ignore slow mode and chat permissions
		UnrestrictBoostCount int64 `json:"unrestrict_boost_count"`
		// Optional. The time after which all messages sent to the chat will be automatically deleted; in seconds
		MessageAutoDeleteTime int64 `json:"message_auto_delete_time"`
		// Optional. True, if aggressive anti-spam checks are enabled in the supergroup.
		// The field is only available to chat administrators.
		HasAggressiveAntiSpamEnabled bool `json:"has_aggressive_anti_spam_enabled"`
		// Optional. True, if non-administrators can only get the list of bots and administrators in the chat
		HasHiddenMembers bool `json:"has_hidden_members"`
		// Optional. True, if messages from the chat can't be forwarded to other chats
		HasProtectedContent bool `json:"has_protected_content"`
		// Optional. True, if new chat members will have access to old messages; available only to chat administrators
		HasVisibleHistory bool `json:"has_visible_history"`
		// Optional. For supergroups, name of the group sticker set
		StickerSetName string `json:"sticker_set_name"`
		// Optional. True, if the bot can change the group sticker set
		CanSetStickerSet bool `json:"can_set_sticker_set"`
		// Optional. For supergroups, the name of the group's custom emoji sticker set.
		// Custom emoji from this set can be used by all users and bots in the group.
		CustomEmojiStickerSetName string `json:"custom_emoji_sticker_set_name"`
		// Optional. Unique identifier for the linked chat, i.e.
		// the discussion group identifier for a channel and vice versa; for supergroups and channel chats.
		// This identifier may be greater than 32 bits and some programming languages may have difficulty/silent defects in interpreting it.
		// But it is smaller than 52 bits, so a signed 64 bit integer or double-precision float type are safe for storing this identifier.
		LinkedChatId int64 `json:"linked_chat_id"`
		// Optional. For supergroups, the location to which the supergroup is connected
		Location *ChatLocation `json:"location"`
		// Joint of structs, used for parsing variant interfaces.
		AvailableReactions []*ReactionTypeUnmarshalJoinedAvailableReactions `json:"available_reactions,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Id = inst.Id
	impl.Type = inst.Type
	impl.Title = inst.Title
	impl.Username = inst.Username
	impl.FirstName = inst.FirstName
	impl.LastName = inst.LastName
	impl.IsForum = inst.IsForum
	impl.AccentColorId = inst.AccentColorId
	impl.MaxReactionCount = inst.MaxReactionCount
	impl.Photo = inst.Photo
	impl.ActiveUsernames = inst.ActiveUsernames
	impl.Birthdate = inst.Birthdate
	impl.BusinessIntro = inst.BusinessIntro
	impl.BusinessLocation = inst.BusinessLocation
	impl.BusinessOpeningHours = inst.BusinessOpeningHours
	impl.PersonalChat = inst.PersonalChat
	impl.BackgroundCustomEmojiId = inst.BackgroundCustomEmojiId
	impl.ProfileAccentColorId = inst.ProfileAccentColorId
	impl.ProfileBackgroundCustomEmojiId = inst.ProfileBackgroundCustomEmojiId
	impl.EmojiStatusCustomEmojiId = inst.EmojiStatusCustomEmojiId
	impl.EmojiStatusExpirationDate = inst.EmojiStatusExpirationDate
	impl.Bio = inst.Bio
	impl.HasPrivateForwards = inst.HasPrivateForwards
	impl.HasRestrictedVoiceAndVideoMessages = inst.HasRestrictedVoiceAndVideoMessages
	impl.JoinToSendMessages = inst.JoinToSendMessages
	impl.JoinByRequest = inst.JoinByRequest
	impl.Description = inst.Description
	impl.InviteLink = inst.InviteLink
	impl.PinnedMessage = inst.PinnedMessage
	impl.Permissions = inst.Permissions
	impl.CanSendPaidMedia = inst.CanSendPaidMedia
	impl.SlowModeDelay = inst.SlowModeDelay
	impl.UnrestrictBoostCount = inst.UnrestrictBoostCount
	impl.MessageAutoDeleteTime = inst.MessageAutoDeleteTime
	impl.HasAggressiveAntiSpamEnabled = inst.HasAggressiveAntiSpamEnabled
	impl.HasHiddenMembers = inst.HasHiddenMembers
	impl.HasProtectedContent = inst.HasProtectedContent
	impl.HasVisibleHistory = inst.HasVisibleHistory
	impl.StickerSetName = inst.StickerSetName
	impl.CanSetStickerSet = inst.CanSetStickerSet
	impl.CustomEmojiStickerSetName = inst.CustomEmojiStickerSetName
	impl.LinkedChatId = inst.LinkedChatId
	impl.Location = inst.Location
	if len(inst.AvailableReactions) != 0 {
		impl.AvailableReactions = []ReactionType{}
		for _, item := range inst.AvailableReactions {
			if item == nil || item.Type == nil {
				continue
			}
			switch *item.Type {
			case "custom_emoji":
				impl.AvailableReactions = append(impl.AvailableReactions, &ReactionTypeCustomEmoji{
					Type:          deref(item.Type),
					CustomEmojiId: deref(item.CustomEmojiId),
				})
			case "emoji":
				impl.AvailableReactions = append(impl.AvailableReactions, &ReactionTypeEmoji{
					Type:  deref(item.Type),
					Emoji: deref(item.Emoji),
				})
			case "paid":
				impl.AvailableReactions = append(impl.AvailableReactions, &ReactionTypePaid{
					Type: deref(item.Type),
				})
			}
		}
	}
	return nil
}

func (impl *ChatMemberUpdated) UnmarshalJSON(data []byte) error {
	type ChatMemberUnmarshalJoinedOldChatMember struct {
		Status                *string `json:"status"`
		User                  **User  `json:"user"`
		IsAnonymous           *bool   `json:"is_anonymous"`
		CustomTitle           *string `json:"custom_title"`
		CanBeEdited           *bool   `json:"can_be_edited"`
		CanManageChat         *bool   `json:"can_manage_chat"`
		CanDeleteMessages     *bool   `json:"can_delete_messages"`
		CanManageVideoChats   *bool   `json:"can_manage_video_chats"`
		CanRestrictMembers    *bool   `json:"can_restrict_members"`
		CanPromoteMembers     *bool   `json:"can_promote_members"`
		CanChangeInfo         *bool   `json:"can_change_info"`
		CanInviteUsers        *bool   `json:"can_invite_users"`
		CanPostStories        *bool   `json:"can_post_stories"`
		CanEditStories        *bool   `json:"can_edit_stories"`
		CanDeleteStories      *bool   `json:"can_delete_stories"`
		CanPostMessages       *bool   `json:"can_post_messages"`
		CanEditMessages       *bool   `json:"can_edit_messages"`
		CanPinMessages        *bool   `json:"can_pin_messages"`
		CanManageTopics       *bool   `json:"can_manage_topics"`
		UntilDate             *int64  `json:"until_date"`
		IsMember              *bool   `json:"is_member"`
		CanSendMessages       *bool   `json:"can_send_messages"`
		CanSendAudios         *bool   `json:"can_send_audios"`
		CanSendDocuments      *bool   `json:"can_send_documents"`
		CanSendPhotos         *bool   `json:"can_send_photos"`
		CanSendVideos         *bool   `json:"can_send_videos"`
		CanSendVideoNotes     *bool   `json:"can_send_video_notes"`
		CanSendVoiceNotes     *bool   `json:"can_send_voice_notes"`
		CanSendPolls          *bool   `json:"can_send_polls"`
		CanSendOtherMessages  *bool   `json:"can_send_other_messages"`
		CanAddWebPagePreviews *bool   `json:"can_add_web_page_previews"`
	}
	type ChatMemberUnmarshalJoinedNewChatMember struct {
		Status                *string `json:"status"`
		User                  **User  `json:"user"`
		IsAnonymous           *bool   `json:"is_anonymous"`
		CustomTitle           *string `json:"custom_title"`
		CanBeEdited           *bool   `json:"can_be_edited"`
		CanManageChat         *bool   `json:"can_manage_chat"`
		CanDeleteMessages     *bool   `json:"can_delete_messages"`
		CanManageVideoChats   *bool   `json:"can_manage_video_chats"`
		CanRestrictMembers    *bool   `json:"can_restrict_members"`
		CanPromoteMembers     *bool   `json:"can_promote_members"`
		CanChangeInfo         *bool   `json:"can_change_info"`
		CanInviteUsers        *bool   `json:"can_invite_users"`
		CanPostStories        *bool   `json:"can_post_stories"`
		CanEditStories        *bool   `json:"can_edit_stories"`
		CanDeleteStories      *bool   `json:"can_delete_stories"`
		CanPostMessages       *bool   `json:"can_post_messages"`
		CanEditMessages       *bool   `json:"can_edit_messages"`
		CanPinMessages        *bool   `json:"can_pin_messages"`
		CanManageTopics       *bool   `json:"can_manage_topics"`
		UntilDate             *int64  `json:"until_date"`
		IsMember              *bool   `json:"is_member"`
		CanSendMessages       *bool   `json:"can_send_messages"`
		CanSendAudios         *bool   `json:"can_send_audios"`
		CanSendDocuments      *bool   `json:"can_send_documents"`
		CanSendPhotos         *bool   `json:"can_send_photos"`
		CanSendVideos         *bool   `json:"can_send_videos"`
		CanSendVideoNotes     *bool   `json:"can_send_video_notes"`
		CanSendVoiceNotes     *bool   `json:"can_send_voice_notes"`
		CanSendPolls          *bool   `json:"can_send_polls"`
		CanSendOtherMessages  *bool   `json:"can_send_other_messages"`
		CanAddWebPagePreviews *bool   `json:"can_add_web_page_previews"`
	}
	type BaseInstance struct {
		// Chat the user belongs to
		Chat *Chat `json:"chat"`
		// Performer of the action, which resulted in the change
		From *User `json:"from"`
		// Date the change was done in Unix time
		Date int64 `json:"date"`
		// Optional. Chat invite link, which was used by the user to join the chat; for joining by invite link events only.
		InviteLink *ChatInviteLink `json:"invite_link"`
		// Optional.
		// True, if the user joined the chat after sending a direct join request without using an invite link and being approved by an administrator
		ViaJoinRequest bool `json:"via_join_request"`
		// Optional. True, if the user joined the chat via a chat folder invite link
		ViaChatFolderInviteLink bool `json:"via_chat_folder_invite_link"`
		// Joint of structs, used for parsing variant interfaces.
		OldChatMember *ChatMemberUnmarshalJoinedOldChatMember `json:"old_chat_member"`
		NewChatMember *ChatMemberUnmarshalJoinedNewChatMember `json:"new_chat_member"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Chat = inst.Chat
	impl.From = inst.From
	impl.Date = inst.Date
	impl.InviteLink = inst.InviteLink
	impl.ViaJoinRequest = inst.ViaJoinRequest
	impl.ViaChatFolderInviteLink = inst.ViaChatFolderInviteLink
	if inst.OldChatMember != nil && inst.OldChatMember.Status == nil {
		switch *inst.OldChatMember.Status {
		case "administrator":
			impl.OldChatMember = &ChatMemberAdministrator{
				Status:              deref(inst.OldChatMember.Status),
				User:                deref(inst.OldChatMember.User),
				CanBeEdited:         deref(inst.OldChatMember.CanBeEdited),
				IsAnonymous:         deref(inst.OldChatMember.IsAnonymous),
				CanManageChat:       deref(inst.OldChatMember.CanManageChat),
				CanDeleteMessages:   deref(inst.OldChatMember.CanDeleteMessages),
				CanManageVideoChats: deref(inst.OldChatMember.CanManageVideoChats),
				CanRestrictMembers:  deref(inst.OldChatMember.CanRestrictMembers),
				CanPromoteMembers:   deref(inst.OldChatMember.CanPromoteMembers),
				CanChangeInfo:       deref(inst.OldChatMember.CanChangeInfo),
				CanInviteUsers:      deref(inst.OldChatMember.CanInviteUsers),
				CanPostStories:      deref(inst.OldChatMember.CanPostStories),
				CanEditStories:      deref(inst.OldChatMember.CanEditStories),
				CanDeleteStories:    deref(inst.OldChatMember.CanDeleteStories),
				CanPostMessages:     deref(inst.OldChatMember.CanPostMessages),
				CanEditMessages:     deref(inst.OldChatMember.CanEditMessages),
				CanPinMessages:      deref(inst.OldChatMember.CanPinMessages),
				CanManageTopics:     deref(inst.OldChatMember.CanManageTopics),
				CustomTitle:         deref(inst.OldChatMember.CustomTitle),
			}
		case "creator":
			impl.OldChatMember = &ChatMemberOwner{
				Status:      deref(inst.OldChatMember.Status),
				User:        deref(inst.OldChatMember.User),
				IsAnonymous: deref(inst.OldChatMember.IsAnonymous),
				CustomTitle: deref(inst.OldChatMember.CustomTitle),
			}
		case "kicked":
			impl.OldChatMember = &ChatMemberBanned{
				Status:    deref(inst.OldChatMember.Status),
				User:      deref(inst.OldChatMember.User),
				UntilDate: deref(inst.OldChatMember.UntilDate),
			}
		case "left":
			impl.OldChatMember = &ChatMemberLeft{
				Status: deref(inst.OldChatMember.Status),
				User:   deref(inst.OldChatMember.User),
			}
		case "member":
			impl.OldChatMember = &ChatMemberMember{
				Status:    deref(inst.OldChatMember.Status),
				User:      deref(inst.OldChatMember.User),
				UntilDate: deref(inst.OldChatMember.UntilDate),
			}
		case "restricted":
			impl.OldChatMember = &ChatMemberRestricted{
				Status:                deref(inst.OldChatMember.Status),
				User:                  deref(inst.OldChatMember.User),
				IsMember:              deref(inst.OldChatMember.IsMember),
				CanSendMessages:       deref(inst.OldChatMember.CanSendMessages),
				CanSendAudios:         deref(inst.OldChatMember.CanSendAudios),
				CanSendDocuments:      deref(inst.OldChatMember.CanSendDocuments),
				CanSendPhotos:         deref(inst.OldChatMember.CanSendPhotos),
				CanSendVideos:         deref(inst.OldChatMember.CanSendVideos),
				CanSendVideoNotes:     deref(inst.OldChatMember.CanSendVideoNotes),
				CanSendVoiceNotes:     deref(inst.OldChatMember.CanSendVoiceNotes),
				CanSendPolls:          deref(inst.OldChatMember.CanSendPolls),
				CanSendOtherMessages:  deref(inst.OldChatMember.CanSendOtherMessages),
				CanAddWebPagePreviews: deref(inst.OldChatMember.CanAddWebPagePreviews),
				CanChangeInfo:         deref(inst.OldChatMember.CanChangeInfo),
				CanInviteUsers:        deref(inst.OldChatMember.CanInviteUsers),
				CanPinMessages:        deref(inst.OldChatMember.CanPinMessages),
				CanManageTopics:       deref(inst.OldChatMember.CanManageTopics),
				UntilDate:             deref(inst.OldChatMember.UntilDate),
			}
		}
	}
	if inst.NewChatMember != nil && inst.NewChatMember.Status == nil {
		switch *inst.NewChatMember.Status {
		case "creator":
			impl.NewChatMember = &ChatMemberOwner{
				Status:      deref(inst.NewChatMember.Status),
				User:        deref(inst.NewChatMember.User),
				IsAnonymous: deref(inst.NewChatMember.IsAnonymous),
				CustomTitle: deref(inst.NewChatMember.CustomTitle),
			}
		case "kicked":
			impl.NewChatMember = &ChatMemberBanned{
				Status:    deref(inst.NewChatMember.Status),
				User:      deref(inst.NewChatMember.User),
				UntilDate: deref(inst.NewChatMember.UntilDate),
			}
		case "left":
			impl.NewChatMember = &ChatMemberLeft{
				Status: deref(inst.NewChatMember.Status),
				User:   deref(inst.NewChatMember.User),
			}
		case "member":
			impl.NewChatMember = &ChatMemberMember{
				Status:    deref(inst.NewChatMember.Status),
				User:      deref(inst.NewChatMember.User),
				UntilDate: deref(inst.NewChatMember.UntilDate),
			}
		case "restricted":
			impl.NewChatMember = &ChatMemberRestricted{
				Status:                deref(inst.NewChatMember.Status),
				User:                  deref(inst.NewChatMember.User),
				IsMember:              deref(inst.NewChatMember.IsMember),
				CanSendMessages:       deref(inst.NewChatMember.CanSendMessages),
				CanSendAudios:         deref(inst.NewChatMember.CanSendAudios),
				CanSendDocuments:      deref(inst.NewChatMember.CanSendDocuments),
				CanSendPhotos:         deref(inst.NewChatMember.CanSendPhotos),
				CanSendVideos:         deref(inst.NewChatMember.CanSendVideos),
				CanSendVideoNotes:     deref(inst.NewChatMember.CanSendVideoNotes),
				CanSendVoiceNotes:     deref(inst.NewChatMember.CanSendVoiceNotes),
				CanSendPolls:          deref(inst.NewChatMember.CanSendPolls),
				CanSendOtherMessages:  deref(inst.NewChatMember.CanSendOtherMessages),
				CanAddWebPagePreviews: deref(inst.NewChatMember.CanAddWebPagePreviews),
				CanChangeInfo:         deref(inst.NewChatMember.CanChangeInfo),
				CanInviteUsers:        deref(inst.NewChatMember.CanInviteUsers),
				CanPinMessages:        deref(inst.NewChatMember.CanPinMessages),
				CanManageTopics:       deref(inst.NewChatMember.CanManageTopics),
				UntilDate:             deref(inst.NewChatMember.UntilDate),
			}
		case "administrator":
			impl.NewChatMember = &ChatMemberAdministrator{
				Status:              deref(inst.NewChatMember.Status),
				User:                deref(inst.NewChatMember.User),
				CanBeEdited:         deref(inst.NewChatMember.CanBeEdited),
				IsAnonymous:         deref(inst.NewChatMember.IsAnonymous),
				CanManageChat:       deref(inst.NewChatMember.CanManageChat),
				CanDeleteMessages:   deref(inst.NewChatMember.CanDeleteMessages),
				CanManageVideoChats: deref(inst.NewChatMember.CanManageVideoChats),
				CanRestrictMembers:  deref(inst.NewChatMember.CanRestrictMembers),
				CanPromoteMembers:   deref(inst.NewChatMember.CanPromoteMembers),
				CanChangeInfo:       deref(inst.NewChatMember.CanChangeInfo),
				CanInviteUsers:      deref(inst.NewChatMember.CanInviteUsers),
				CanPostStories:      deref(inst.NewChatMember.CanPostStories),
				CanEditStories:      deref(inst.NewChatMember.CanEditStories),
				CanDeleteStories:    deref(inst.NewChatMember.CanDeleteStories),
				CanPostMessages:     deref(inst.NewChatMember.CanPostMessages),
				CanEditMessages:     deref(inst.NewChatMember.CanEditMessages),
				CanPinMessages:      deref(inst.NewChatMember.CanPinMessages),
				CanManageTopics:     deref(inst.NewChatMember.CanManageTopics),
				CustomTitle:         deref(inst.NewChatMember.CustomTitle),
			}
		}
	}
	return nil
}

func (impl *ExternalReplyInfo) UnmarshalJSON(data []byte) error {
	type MessageOriginUnmarshalJoinedOrigin struct {
		Type            *string `json:"type"`
		Date            *int64  `json:"date"`
		SenderUser      **User  `json:"sender_user"`
		SenderUserName  *string `json:"sender_user_name"`
		SenderChat      **Chat  `json:"sender_chat"`
		AuthorSignature *string `json:"author_signature"`
		Chat            **Chat  `json:"chat"`
		MessageId       *int64  `json:"message_id"`
	}
	type BaseInstance struct {
		// Optional. Chat the original message belongs to. Available only if the chat is a supergroup or a channel.
		Chat *Chat `json:"chat"`
		// Optional. Unique message identifier inside the original chat.
		// Available only if the original chat is a supergroup or a channel.
		MessageId int64 `json:"message_id"`
		// Optional. Options used for link preview generation for the original message, if it is a text message
		LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options"`
		// Optional. Message is an animation, information about the animation
		Animation *TelegramAnimation `json:"animation"`
		// Optional. Message is an audio file, information about the file
		Audio *TelegramAudio `json:"audio"`
		// Optional. Message is a general file, information about the file
		Document *TelegramDocument `json:"document"`
		// Optional. Message contains paid media; information about the paid media
		PaidMedia *PaidMediaInfo `json:"paid_media"`
		// Optional. Message is a photo, available sizes of the photo
		Photo TelegramPhoto `json:"photo"`
		// Optional. Message is a sticker, information about the sticker
		Sticker *Sticker `json:"sticker"`
		// Optional. Message is a forwarded story
		Story *Story `json:"story"`
		// Optional. Message is a video, information about the video
		Video *TelegramVideo `json:"video"`
		// Optional. Message is a video note, information about the video message
		VideoNote *VideoNote `json:"video_note"`
		// Optional. Message is a voice message, information about the file
		Voice *Voice `json:"voice"`
		// Optional. True, if the message media is covered by a spoiler animation
		HasMediaSpoiler bool `json:"has_media_spoiler"`
		// Optional. Message is a shared contact, information about the contact
		Contact *Contact `json:"contact"`
		// Optional. Message is a dice with random value
		Dice *Dice `json:"dice"`
		// Optional. Message is a game, information about the game. More about games: https://core.telegram.org/bots/api#games
		Game *Game `json:"game"`
		// Optional. Message is a scheduled giveaway, information about the giveaway
		Giveaway *Giveaway `json:"giveaway"`
		// Optional. A giveaway with public winners was completed
		GiveawayWinners *GiveawayWinners `json:"giveaway_winners"`
		// Optional. Message is an invoice for a payment, information about the invoice.
		// More about payments: https://core.telegram.org/bots/api#payments
		Invoice *Invoice `json:"invoice"`
		// Optional. Message is a shared location, information about the location
		Location *Location `json:"location"`
		// Optional. Message is a native poll, information about the poll
		Poll *Poll `json:"poll"`
		// Optional. Message is a venue, information about the venue
		Venue *Venue `json:"venue"`
		// Joint of structs, used for parsing variant interfaces.
		Origin *MessageOriginUnmarshalJoinedOrigin `json:"origin"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Chat = inst.Chat
	impl.MessageId = inst.MessageId
	impl.LinkPreviewOptions = inst.LinkPreviewOptions
	impl.Animation = inst.Animation
	impl.Audio = inst.Audio
	impl.Document = inst.Document
	impl.PaidMedia = inst.PaidMedia
	impl.Photo = inst.Photo
	impl.Sticker = inst.Sticker
	impl.Story = inst.Story
	impl.Video = inst.Video
	impl.VideoNote = inst.VideoNote
	impl.Voice = inst.Voice
	impl.HasMediaSpoiler = inst.HasMediaSpoiler
	impl.Contact = inst.Contact
	impl.Dice = inst.Dice
	impl.Game = inst.Game
	impl.Giveaway = inst.Giveaway
	impl.GiveawayWinners = inst.GiveawayWinners
	impl.Invoice = inst.Invoice
	impl.Location = inst.Location
	impl.Poll = inst.Poll
	impl.Venue = inst.Venue
	if inst.Origin != nil && inst.Origin.Type == nil {
		switch *inst.Origin.Type {
		case "channel":
			impl.Origin = &MessageOriginChannel{
				Type:            deref(inst.Origin.Type),
				Date:            deref(inst.Origin.Date),
				Chat:            deref(inst.Origin.Chat),
				MessageId:       deref(inst.Origin.MessageId),
				AuthorSignature: deref(inst.Origin.AuthorSignature),
			}
		case "chat":
			impl.Origin = &MessageOriginChat{
				Type:            deref(inst.Origin.Type),
				Date:            deref(inst.Origin.Date),
				SenderChat:      deref(inst.Origin.SenderChat),
				AuthorSignature: deref(inst.Origin.AuthorSignature),
			}
		case "hidden_user":
			impl.Origin = &MessageOriginHiddenUser{
				Type:           deref(inst.Origin.Type),
				Date:           deref(inst.Origin.Date),
				SenderUserName: deref(inst.Origin.SenderUserName),
			}
		case "user":
			impl.Origin = &MessageOriginUser{
				Type:       deref(inst.Origin.Type),
				Date:       deref(inst.Origin.Date),
				SenderUser: deref(inst.Origin.SenderUser),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultArticle) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be article
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 Bytes
		Id string `json:"id"`
		// Title of the result
		Title string `json:"title"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Optional. URL of the result
		Url string `json:"url"`
		// Optional. Pass True if you don't want the URL to be shown in the message
		HideUrl bool `json:"hide_url"`
		// Optional. Short description of the result
		Description string `json:"description"`
		// Optional. Url of the thumbnail for the result
		ThumbnailUrl string `json:"thumbnail_url"`
		// Optional. Thumbnail width
		ThumbnailWidth int64 `json:"thumbnail_width"`
		// Optional. Thumbnail height
		ThumbnailHeight int64 `json:"thumbnail_height"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.Title = inst.Title
	impl.ReplyMarkup = inst.ReplyMarkup
	impl.Url = inst.Url
	impl.HideUrl = inst.HideUrl
	impl.Description = inst.Description
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.ThumbnailWidth = inst.ThumbnailWidth
	impl.ThumbnailHeight = inst.ThumbnailHeight
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultAudio) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be audio
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid URL for the audio file
		AudioUrl string `json:"audio_url"`
		// Title
		Title string `json:"title"`
		// Optional. Caption, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the audio caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Performer
		Performer string `json:"performer"`
		// Optional. Audio duration in seconds
		AudioDuration int64 `json:"audio_duration"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.AudioUrl = inst.AudioUrl
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.Performer = inst.Performer
	impl.AudioDuration = inst.AudioDuration
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultCachedAudio) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be audio
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid file identifier for the audio file
		AudioFileId string `json:"audio_file_id"`
		// Optional. Caption, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the audio caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.AudioFileId = inst.AudioFileId
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultCachedDocument) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be document
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// Title for the result
		Title string `json:"title"`
		// A valid file identifier for the file
		DocumentFileId string `json:"document_file_id"`
		// Optional. Short description of the result
		Description string `json:"description"`
		// Optional. Caption of the document to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the document caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.Title = inst.Title
	impl.DocumentFileId = inst.DocumentFileId
	impl.Description = inst.Description
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultCachedGif) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be gif
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid file identifier for the GIF file
		GifFileId string `json:"gif_file_id"`
		// Optional. Title for the result
		Title string `json:"title"`
		// Optional. Caption of the GIF file to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Pass True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.GifFileId = inst.GifFileId
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultCachedMpeg4Gif) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be mpeg4_gif
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid file identifier for the MPEG4 file
		Mpeg4FileId string `json:"mpeg4_file_id"`
		// Optional. Title for the result
		Title string `json:"title"`
		// Optional. Caption of the MPEG-4 file to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Pass True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.Mpeg4FileId = inst.Mpeg4FileId
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultCachedPhoto) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be photo
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid file identifier of the photo
		PhotoFileId string `json:"photo_file_id"`
		// Optional. Title for the result
		Title string `json:"title"`
		// Optional. Short description of the result
		Description string `json:"description"`
		// Optional. Caption of the photo to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the photo caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Pass True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.PhotoFileId = inst.PhotoFileId
	impl.Title = inst.Title
	impl.Description = inst.Description
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultCachedSticker) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be sticker
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid file identifier of the sticker
		StickerFileId string `json:"sticker_file_id"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.StickerFileId = inst.StickerFileId
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultCachedVideo) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be video
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid file identifier for the video file
		VideoFileId string `json:"video_file_id"`
		// Title for the result
		Title string `json:"title"`
		// Optional. Short description of the result
		Description string `json:"description"`
		// Optional. Caption of the video to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the video caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Pass True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.VideoFileId = inst.VideoFileId
	impl.Title = inst.Title
	impl.Description = inst.Description
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultCachedVoice) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be voice
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid file identifier for the voice message
		VoiceFileId string `json:"voice_file_id"`
		// Voice message title
		Title string `json:"title"`
		// Optional. Caption, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the voice message caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.VoiceFileId = inst.VoiceFileId
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultContact) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be contact
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 Bytes
		Id string `json:"id"`
		// Contact's phone number
		PhoneNumber string `json:"phone_number"`
		// Contact's first name
		FirstName string `json:"first_name"`
		// Optional. Contact's last name
		LastName string `json:"last_name"`
		// Optional. Additional data about the contact in the form of a vCard, 0-2048 bytes
		Vcard string `json:"vcard"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Optional. Url of the thumbnail for the result
		ThumbnailUrl string `json:"thumbnail_url"`
		// Optional. Thumbnail width
		ThumbnailWidth int64 `json:"thumbnail_width"`
		// Optional. Thumbnail height
		ThumbnailHeight int64 `json:"thumbnail_height"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.PhoneNumber = inst.PhoneNumber
	impl.FirstName = inst.FirstName
	impl.LastName = inst.LastName
	impl.Vcard = inst.Vcard
	impl.ReplyMarkup = inst.ReplyMarkup
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.ThumbnailWidth = inst.ThumbnailWidth
	impl.ThumbnailHeight = inst.ThumbnailHeight
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultDocument) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be document
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// Title for the result
		Title string `json:"title"`
		// Optional. Caption of the document to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the document caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// A valid URL for the file
		DocumentUrl string `json:"document_url"`
		// MIME type of the content of the file, either "application/pdf" or "application/zip"
		MimeType string `json:"mime_type"`
		// Optional. Short description of the result
		Description string `json:"description"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Optional. URL of the thumbnail (JPEG only) for the file
		ThumbnailUrl string `json:"thumbnail_url"`
		// Optional. Thumbnail width
		ThumbnailWidth int64 `json:"thumbnail_width"`
		// Optional. Thumbnail height
		ThumbnailHeight int64 `json:"thumbnail_height"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.DocumentUrl = inst.DocumentUrl
	impl.MimeType = inst.MimeType
	impl.Description = inst.Description
	impl.ReplyMarkup = inst.ReplyMarkup
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.ThumbnailWidth = inst.ThumbnailWidth
	impl.ThumbnailHeight = inst.ThumbnailHeight
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultGif) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be gif
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid URL for the GIF file. File size must not exceed 1MB
		GifUrl string `json:"gif_url"`
		// Optional. Width of the GIF
		GifWidth int64 `json:"gif_width"`
		// Optional. Height of the GIF
		GifHeight int64 `json:"gif_height"`
		// Optional. Duration of the GIF in seconds
		GifDuration int64 `json:"gif_duration"`
		// URL of the static (JPEG or GIF) or animated (MPEG4) thumbnail for the result
		ThumbnailUrl string `json:"thumbnail_url"`
		// Optional. MIME type of the thumbnail, must be one of "image/jpeg", "image/gif", or "video/mp4".
		// Defaults to "image/jpeg"
		ThumbnailMimeType string `json:"thumbnail_mime_type"`
		// Optional. Title for the result
		Title string `json:"title"`
		// Optional. Caption of the GIF file to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Pass True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.GifUrl = inst.GifUrl
	impl.GifWidth = inst.GifWidth
	impl.GifHeight = inst.GifHeight
	impl.GifDuration = inst.GifDuration
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.ThumbnailMimeType = inst.ThumbnailMimeType
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultLocation) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
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
		HorizontalAccuracy float64 `json:"horizontal_accuracy"`
		// Optional.
		// Period in seconds during which the location can be updated, should be between 60 and 86400, or 0x7FFFFFFF for live locations that can be edited indefinitely.
		LivePeriod int64 `json:"live_period"`
		// Optional. For live locations, a direction in which the user is moving, in degrees.
		// Must be between 1 and 360 if specified.
		Heading int64 `json:"heading"`
		// Optional. Must be between 1 and 100000 if specified.
		// For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters.
		ProximityAlertRadius int64 `json:"proximity_alert_radius"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Optional. Url of the thumbnail for the result
		ThumbnailUrl string `json:"thumbnail_url"`
		// Optional. Thumbnail width
		ThumbnailWidth int64 `json:"thumbnail_width"`
		// Optional. Thumbnail height
		ThumbnailHeight int64 `json:"thumbnail_height"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.Latitude = inst.Latitude
	impl.Longitude = inst.Longitude
	impl.Title = inst.Title
	impl.HorizontalAccuracy = inst.HorizontalAccuracy
	impl.LivePeriod = inst.LivePeriod
	impl.Heading = inst.Heading
	impl.ProximityAlertRadius = inst.ProximityAlertRadius
	impl.ReplyMarkup = inst.ReplyMarkup
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.ThumbnailWidth = inst.ThumbnailWidth
	impl.ThumbnailHeight = inst.ThumbnailHeight
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultMpeg4Gif) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be mpeg4_gif
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid URL for the MPEG4 file. File size must not exceed 1MB
		Mpeg4Url string `json:"mpeg4_url"`
		// Optional. Video width
		Mpeg4Width int64 `json:"mpeg4_width"`
		// Optional. Video height
		Mpeg4Height int64 `json:"mpeg4_height"`
		// Optional. Video duration in seconds
		Mpeg4Duration int64 `json:"mpeg4_duration"`
		// URL of the static (JPEG or GIF) or animated (MPEG4) thumbnail for the result
		ThumbnailUrl string `json:"thumbnail_url"`
		// Optional. MIME type of the thumbnail, must be one of "image/jpeg", "image/gif", or "video/mp4".
		// Defaults to "image/jpeg"
		ThumbnailMimeType string `json:"thumbnail_mime_type"`
		// Optional. Title for the result
		Title string `json:"title"`
		// Optional. Caption of the MPEG-4 file to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Pass True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.Mpeg4Url = inst.Mpeg4Url
	impl.Mpeg4Width = inst.Mpeg4Width
	impl.Mpeg4Height = inst.Mpeg4Height
	impl.Mpeg4Duration = inst.Mpeg4Duration
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.ThumbnailMimeType = inst.ThumbnailMimeType
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultPhoto) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be photo
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid URL of the photo. Photo must be in JPEG format. Photo size must not exceed 5MB
		PhotoUrl string `json:"photo_url"`
		// URL of the thumbnail for the photo
		ThumbnailUrl string `json:"thumbnail_url"`
		// Optional. Width of the photo
		PhotoWidth int64 `json:"photo_width"`
		// Optional. Height of the photo
		PhotoHeight int64 `json:"photo_height"`
		// Optional. Title for the result
		Title string `json:"title"`
		// Optional. Short description of the result
		Description string `json:"description"`
		// Optional. Caption of the photo to be sent, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the photo caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Pass True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.PhotoUrl = inst.PhotoUrl
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.PhotoWidth = inst.PhotoWidth
	impl.PhotoHeight = inst.PhotoHeight
	impl.Title = inst.Title
	impl.Description = inst.Description
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultVenue) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
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
		FoursquareId string `json:"foursquare_id"`
		// Optional. Foursquare type of the venue, if known.
		// (For example, "arts_entertainment/default", "arts_entertainment/aquarium" or "food/icecream".)
		FoursquareType string `json:"foursquare_type"`
		// Optional. Google Places identifier of the venue
		GooglePlaceId string `json:"google_place_id"`
		// Optional. Google Places type of the venue. (See supported types.)
		GooglePlaceType string `json:"google_place_type"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Optional. Url of the thumbnail for the result
		ThumbnailUrl string `json:"thumbnail_url"`
		// Optional. Thumbnail width
		ThumbnailWidth int64 `json:"thumbnail_width"`
		// Optional. Thumbnail height
		ThumbnailHeight int64 `json:"thumbnail_height"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.Latitude = inst.Latitude
	impl.Longitude = inst.Longitude
	impl.Title = inst.Title
	impl.Address = inst.Address
	impl.FoursquareId = inst.FoursquareId
	impl.FoursquareType = inst.FoursquareType
	impl.GooglePlaceId = inst.GooglePlaceId
	impl.GooglePlaceType = inst.GooglePlaceType
	impl.ReplyMarkup = inst.ReplyMarkup
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.ThumbnailWidth = inst.ThumbnailWidth
	impl.ThumbnailHeight = inst.ThumbnailHeight
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultVideo) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
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
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the video caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Pass True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. Video width
		VideoWidth int64 `json:"video_width"`
		// Optional. Video height
		VideoHeight int64 `json:"video_height"`
		// Optional. Video duration in seconds
		VideoDuration int64 `json:"video_duration"`
		// Optional. Short description of the result
		Description string `json:"description"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.VideoUrl = inst.VideoUrl
	impl.MimeType = inst.MimeType
	impl.ThumbnailUrl = inst.ThumbnailUrl
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.VideoWidth = inst.VideoWidth
	impl.VideoHeight = inst.VideoHeight
	impl.VideoDuration = inst.VideoDuration
	impl.Description = inst.Description
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *InlineQueryResultVoice) UnmarshalJSON(data []byte) error {
	type InputMessageContentUnmarshalJoinedInputMessageContent struct {
		MessageText               *string              `json:"message_text"`
		ParseMode                 *string              `json:"parse_mode"`
		Entities                  *[]*MessageEntity    `json:"entities"`
		LinkPreviewOptions        **LinkPreviewOptions `json:"link_preview_options"`
		Latitude                  *float64             `json:"latitude"`
		Longitude                 *float64             `json:"longitude"`
		HorizontalAccuracy        *float64             `json:"horizontal_accuracy"`
		LivePeriod                *int64               `json:"live_period"`
		Heading                   *int64               `json:"heading"`
		ProximityAlertRadius      *int64               `json:"proximity_alert_radius"`
		Title                     *string              `json:"title"`
		Address                   *string              `json:"address"`
		FoursquareId              *string              `json:"foursquare_id"`
		FoursquareType            *string              `json:"foursquare_type"`
		GooglePlaceId             *string              `json:"google_place_id"`
		GooglePlaceType           *string              `json:"google_place_type"`
		PhoneNumber               *string              `json:"phone_number"`
		FirstName                 *string              `json:"first_name"`
		LastName                  *string              `json:"last_name"`
		Vcard                     *string              `json:"vcard"`
		Description               *string              `json:"description"`
		Payload                   *string              `json:"payload"`
		ProviderToken             *string              `json:"provider_token"`
		Currency                  *string              `json:"currency"`
		Prices                    *[]*LabeledPrice     `json:"prices"`
		MaxTipAmount              *int64               `json:"max_tip_amount"`
		SuggestedTipAmounts       *[]int64             `json:"suggested_tip_amounts"`
		ProviderData              *string              `json:"provider_data"`
		PhotoUrl                  *string              `json:"photo_url"`
		PhotoSize                 *int64               `json:"photo_size"`
		PhotoWidth                *int64               `json:"photo_width"`
		PhotoHeight               *int64               `json:"photo_height"`
		NeedName                  *bool                `json:"need_name"`
		NeedPhoneNumber           *bool                `json:"need_phone_number"`
		NeedEmail                 *bool                `json:"need_email"`
		NeedShippingAddress       *bool                `json:"need_shipping_address"`
		SendPhoneNumberToProvider *bool                `json:"send_phone_number_to_provider"`
		SendEmailToProvider       *bool                `json:"send_email_to_provider"`
		IsFlexible                *bool                `json:"is_flexible"`
	}
	type BaseInstance struct {
		// Type of the result, must be voice
		Type string `json:"type"`
		// Unique identifier for this result, 1-64 bytes
		Id string `json:"id"`
		// A valid URL for the voice recording
		VoiceUrl string `json:"voice_url"`
		// Recording title
		Title string `json:"title"`
		// Optional. Caption, 0-1024 characters after entities parsing
		Caption string `json:"caption"`
		// Optional. Mode for parsing entities in the voice message caption. See formatting options for more details.
		ParseMode string `json:"parse_mode"`
		// Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. Recording duration in seconds
		VoiceDuration int64 `json:"voice_duration"`
		// Optional. Inline keyboard attached to the message
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		InputMessageContent *InputMessageContentUnmarshalJoinedInputMessageContent `json:"input_message_content,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.Id = inst.Id
	impl.VoiceUrl = inst.VoiceUrl
	impl.Title = inst.Title
	impl.Caption = inst.Caption
	impl.ParseMode = inst.ParseMode
	impl.CaptionEntities = inst.CaptionEntities
	impl.VoiceDuration = inst.VoiceDuration
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.InputMessageContent != nil {
		nonEmptyFields := []string{}
		if inst.InputMessageContent.MessageText != nil {
			nonEmptyFields = append(nonEmptyFields, "MessageText")
		}
		if inst.InputMessageContent.ParseMode != nil {
			nonEmptyFields = append(nonEmptyFields, "ParseMode")
		}
		if inst.InputMessageContent.Entities != nil {
			nonEmptyFields = append(nonEmptyFields, "Entities")
		}
		if inst.InputMessageContent.LinkPreviewOptions != nil {
			nonEmptyFields = append(nonEmptyFields, "LinkPreviewOptions")
		}
		if inst.InputMessageContent.PhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "PhoneNumber")
		}
		if inst.InputMessageContent.FirstName != nil {
			nonEmptyFields = append(nonEmptyFields, "FirstName")
		}
		if inst.InputMessageContent.LastName != nil {
			nonEmptyFields = append(nonEmptyFields, "LastName")
		}
		if inst.InputMessageContent.Vcard != nil {
			nonEmptyFields = append(nonEmptyFields, "Vcard")
		}
		if inst.InputMessageContent.Latitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Latitude")
		}
		if inst.InputMessageContent.Longitude != nil {
			nonEmptyFields = append(nonEmptyFields, "Longitude")
		}
		if inst.InputMessageContent.HorizontalAccuracy != nil {
			nonEmptyFields = append(nonEmptyFields, "HorizontalAccuracy")
		}
		if inst.InputMessageContent.LivePeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "LivePeriod")
		}
		if inst.InputMessageContent.Heading != nil {
			nonEmptyFields = append(nonEmptyFields, "Heading")
		}
		if inst.InputMessageContent.ProximityAlertRadius != nil {
			nonEmptyFields = append(nonEmptyFields, "ProximityAlertRadius")
		}
		if inst.InputMessageContent.Title != nil {
			nonEmptyFields = append(nonEmptyFields, "Title")
		}
		if inst.InputMessageContent.Address != nil {
			nonEmptyFields = append(nonEmptyFields, "Address")
		}
		if inst.InputMessageContent.FoursquareId != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareId")
		}
		if inst.InputMessageContent.FoursquareType != nil {
			nonEmptyFields = append(nonEmptyFields, "FoursquareType")
		}
		if inst.InputMessageContent.GooglePlaceId != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceId")
		}
		if inst.InputMessageContent.GooglePlaceType != nil {
			nonEmptyFields = append(nonEmptyFields, "GooglePlaceType")
		}
		if inst.InputMessageContent.Description != nil {
			nonEmptyFields = append(nonEmptyFields, "Description")
		}
		if inst.InputMessageContent.Payload != nil {
			nonEmptyFields = append(nonEmptyFields, "Payload")
		}
		if inst.InputMessageContent.ProviderToken != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderToken")
		}
		if inst.InputMessageContent.Currency != nil {
			nonEmptyFields = append(nonEmptyFields, "Currency")
		}
		if inst.InputMessageContent.Prices != nil {
			nonEmptyFields = append(nonEmptyFields, "Prices")
		}
		if inst.InputMessageContent.MaxTipAmount != nil {
			nonEmptyFields = append(nonEmptyFields, "MaxTipAmount")
		}
		if inst.InputMessageContent.SuggestedTipAmounts != nil {
			nonEmptyFields = append(nonEmptyFields, "SuggestedTipAmounts")
		}
		if inst.InputMessageContent.ProviderData != nil {
			nonEmptyFields = append(nonEmptyFields, "ProviderData")
		}
		if inst.InputMessageContent.PhotoUrl != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoUrl")
		}
		if inst.InputMessageContent.PhotoSize != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoSize")
		}
		if inst.InputMessageContent.PhotoWidth != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoWidth")
		}
		if inst.InputMessageContent.PhotoHeight != nil {
			nonEmptyFields = append(nonEmptyFields, "PhotoHeight")
		}
		if inst.InputMessageContent.NeedName != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedName")
		}
		if inst.InputMessageContent.NeedPhoneNumber != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedPhoneNumber")
		}
		if inst.InputMessageContent.NeedEmail != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedEmail")
		}
		if inst.InputMessageContent.NeedShippingAddress != nil {
			nonEmptyFields = append(nonEmptyFields, "NeedShippingAddress")
		}
		if inst.InputMessageContent.SendPhoneNumberToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendPhoneNumberToProvider")
		}
		if inst.InputMessageContent.SendEmailToProvider != nil {
			nonEmptyFields = append(nonEmptyFields, "SendEmailToProvider")
		}
		if inst.InputMessageContent.IsFlexible != nil {
			nonEmptyFields = append(nonEmptyFields, "IsFlexible")
		}
		switch {
		case containsAll([]string{"MessageText", "ParseMode", "Entities", "LinkPreviewOptions"}, nonEmptyFields):
			impl.InputMessageContent = &InputTextMessageContent{
				MessageText:        deref(inst.InputMessageContent.MessageText),
				ParseMode:          deref(inst.InputMessageContent.ParseMode),
				Entities:           deref(inst.InputMessageContent.Entities),
				LinkPreviewOptions: deref(inst.InputMessageContent.LinkPreviewOptions),
			}
		case containsAll([]string{"PhoneNumber", "FirstName", "LastName", "Vcard"}, nonEmptyFields):
			impl.InputMessageContent = &InputContactMessageContent{
				PhoneNumber: deref(inst.InputMessageContent.PhoneNumber),
				FirstName:   deref(inst.InputMessageContent.FirstName),
				LastName:    deref(inst.InputMessageContent.LastName),
				Vcard:       deref(inst.InputMessageContent.Vcard),
			}
		case containsAll([]string{"Latitude", "Longitude", "HorizontalAccuracy", "LivePeriod", "Heading", "ProximityAlertRadius"}, nonEmptyFields):
			impl.InputMessageContent = &InputLocationMessageContent{
				Latitude:             deref(inst.InputMessageContent.Latitude),
				Longitude:            deref(inst.InputMessageContent.Longitude),
				HorizontalAccuracy:   deref(inst.InputMessageContent.HorizontalAccuracy),
				LivePeriod:           deref(inst.InputMessageContent.LivePeriod),
				Heading:              deref(inst.InputMessageContent.Heading),
				ProximityAlertRadius: deref(inst.InputMessageContent.ProximityAlertRadius),
			}
		case containsAll([]string{"Latitude", "Longitude", "Title", "Address", "FoursquareId", "FoursquareType", "GooglePlaceId", "GooglePlaceType"}, nonEmptyFields):
			impl.InputMessageContent = &InputVenueMessageContent{
				Latitude:        deref(inst.InputMessageContent.Latitude),
				Longitude:       deref(inst.InputMessageContent.Longitude),
				Title:           deref(inst.InputMessageContent.Title),
				Address:         deref(inst.InputMessageContent.Address),
				FoursquareId:    deref(inst.InputMessageContent.FoursquareId),
				FoursquareType:  deref(inst.InputMessageContent.FoursquareType),
				GooglePlaceId:   deref(inst.InputMessageContent.GooglePlaceId),
				GooglePlaceType: deref(inst.InputMessageContent.GooglePlaceType),
			}
		case containsAll([]string{"Title", "Description", "Payload", "ProviderToken", "Currency", "Prices", "MaxTipAmount", "SuggestedTipAmounts", "ProviderData", "PhotoUrl", "PhotoSize", "PhotoWidth", "PhotoHeight", "NeedName", "NeedPhoneNumber", "NeedEmail", "NeedShippingAddress", "SendPhoneNumberToProvider", "SendEmailToProvider", "IsFlexible"}, nonEmptyFields):
			impl.InputMessageContent = &InputInvoiceMessageContent{
				Title:                     deref(inst.InputMessageContent.Title),
				Description:               deref(inst.InputMessageContent.Description),
				Payload:                   deref(inst.InputMessageContent.Payload),
				ProviderToken:             deref(inst.InputMessageContent.ProviderToken),
				Currency:                  deref(inst.InputMessageContent.Currency),
				Prices:                    deref(inst.InputMessageContent.Prices),
				MaxTipAmount:              deref(inst.InputMessageContent.MaxTipAmount),
				SuggestedTipAmounts:       deref(inst.InputMessageContent.SuggestedTipAmounts),
				ProviderData:              deref(inst.InputMessageContent.ProviderData),
				PhotoUrl:                  deref(inst.InputMessageContent.PhotoUrl),
				PhotoSize:                 deref(inst.InputMessageContent.PhotoSize),
				PhotoWidth:                deref(inst.InputMessageContent.PhotoWidth),
				PhotoHeight:               deref(inst.InputMessageContent.PhotoHeight),
				NeedName:                  deref(inst.InputMessageContent.NeedName),
				NeedPhoneNumber:           deref(inst.InputMessageContent.NeedPhoneNumber),
				NeedEmail:                 deref(inst.InputMessageContent.NeedEmail),
				NeedShippingAddress:       deref(inst.InputMessageContent.NeedShippingAddress),
				SendPhoneNumberToProvider: deref(inst.InputMessageContent.SendPhoneNumberToProvider),
				SendEmailToProvider:       deref(inst.InputMessageContent.SendEmailToProvider),
				IsFlexible:                deref(inst.InputMessageContent.IsFlexible),
			}
		}
	}
	return nil
}

func (impl *Message) UnmarshalJSON(data []byte) error {
	type MessageOriginUnmarshalJoinedForwardOrigin struct {
		Type            *string `json:"type"`
		Date            *int64  `json:"date"`
		SenderUser      **User  `json:"sender_user"`
		SenderUserName  *string `json:"sender_user_name"`
		SenderChat      **Chat  `json:"sender_chat"`
		AuthorSignature *string `json:"author_signature"`
		Chat            **Chat  `json:"chat"`
		MessageId       *int64  `json:"message_id"`
	}
	type BaseInstance struct {
		// Unique message identifier inside this chat.
		// In specific instances (e.g., message containing a video sent to a big chat), the server might automatically schedule a message instead of sending it immediately.
		// In such cases, this field will be 0 and the relevant message will be unusable until it is actually sent
		MessageId int64 `json:"message_id"`
		// Optional. Unique identifier of a message thread to which the message belongs; for supergroups only
		MessageThreadId int64 `json:"message_thread_id"`
		// Optional. Sender of the message; may be empty for messages sent to channels.
		// For backward compatibility, if the message was sent on behalf of a chat, the field contains a fake sender user in non-channel chats
		From *User `json:"from"`
		// Optional. Sender of the message when sent on behalf of a chat.
		// For example, the supergroup itself for messages sent by its anonymous administrators or a linked channel for messages automatically forwarded to the channel's discussion group.
		// For backward compatibility, if the message was sent on behalf of a chat, the field from contains a fake sender user in non-channel chats.
		SenderChat *Chat `json:"sender_chat"`
		// Optional. If the sender of the message boosted the chat, the number of boosts added by the user
		SenderBoostCount int64 `json:"sender_boost_count"`
		// Optional. The bot that actually sent the message on behalf of the business account.
		// Available only for outgoing messages sent on behalf of the connected business account.
		SenderBusinessBot *User `json:"sender_business_bot"`
		// Date the message was sent in Unix time. It is always a positive number, representing a valid date.
		Date int64 `json:"date"`
		// Optional. Unique identifier of the business connection from which the message was received.
		// If non-empty, the message belongs to a chat of the corresponding business account that is independent from any potential bot chat which might share the same identifier.
		BusinessConnectionId string `json:"business_connection_id"`
		// Chat the message belongs to
		Chat *Chat `json:"chat"`
		// Optional. True, if the message is sent to a forum topic
		IsTopicMessage bool `json:"is_topic_message"`
		// Optional. True, if the message is a channel post that was automatically forwarded to the connected discussion group
		IsAutomaticForward bool `json:"is_automatic_forward"`
		// Optional. For replies in the same chat and message thread, the original message.
		// Note that the Message object in this field will not contain further reply_to_message fields even if it itself is a reply.
		ReplyToMessage *Message `json:"reply_to_message"`
		// Optional. Information about the message that is being replied to, which may come from another chat or forum topic
		ExternalReply *ExternalReplyInfo `json:"external_reply"`
		// Optional. For replies that quote part of the original message, the quoted part of the message
		Quote *TextQuote `json:"quote"`
		// Optional. For replies to a story, the original story
		ReplyToStory *Story `json:"reply_to_story"`
		// Optional. Bot through which the message was sent
		ViaBot *User `json:"via_bot"`
		// Optional. Date the message was last edited in Unix time
		EditDate int64 `json:"edit_date"`
		// Optional. True, if the message can't be forwarded
		HasProtectedContent bool `json:"has_protected_content"`
		// Optional.
		// True, if the message was sent by an implicit action, for example, as an away or a greeting business message, or as a scheduled message
		IsFromOffline bool `json:"is_from_offline"`
		// Optional. The unique identifier of a media message group this message belongs to
		MediaGroupId string `json:"media_group_id"`
		// Optional.
		// Signature of the post author for messages in channels, or the custom title of an anonymous group administrator
		AuthorSignature string `json:"author_signature"`
		// Optional. For text messages, the actual UTF-8 text of the message
		Text string `json:"text"`
		// Optional. For text messages, special entities like usernames, URLs, bot commands, etc.
		// that appear in the text
		Entities []*MessageEntity `json:"entities"`
		// Optional.
		// Options used for link preview generation for the message, if it is a text message and link preview options were changed
		LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options"`
		// Optional. Unique identifier of the message effect added to the message
		EffectId string `json:"effect_id"`
		// Optional. Message is an animation, information about the animation.
		// For backward compatibility, when this field is set, the document field will also be set
		Animation *TelegramAnimation `json:"animation"`
		// Optional. Message is an audio file, information about the file
		Audio *TelegramAudio `json:"audio"`
		// Optional. Message is a general file, information about the file
		Document *TelegramDocument `json:"document"`
		// Optional. Message contains paid media; information about the paid media
		PaidMedia *PaidMediaInfo `json:"paid_media"`
		// Optional. Message is a photo, available sizes of the photo
		Photo TelegramPhoto `json:"photo"`
		// Optional. Message is a sticker, information about the sticker
		Sticker *Sticker `json:"sticker"`
		// Optional. Message is a forwarded story
		Story *Story `json:"story"`
		// Optional. Message is a video, information about the video
		Video *TelegramVideo `json:"video"`
		// Optional. Message is a video note, information about the video message
		VideoNote *VideoNote `json:"video_note"`
		// Optional. Message is a voice message, information about the file
		Voice *Voice `json:"voice"`
		// Optional. Caption for the animation, audio, document, paid media, photo, video or voice
		Caption string `json:"caption"`
		// Optional. For messages with a caption, special entities like usernames, URLs, bot commands, etc.
		// that appear in the caption
		CaptionEntities []*MessageEntity `json:"caption_entities"`
		// Optional. True, if the caption must be shown above the message media
		ShowCaptionAboveMedia bool `json:"show_caption_above_media"`
		// Optional. True, if the message media is covered by a spoiler animation
		HasMediaSpoiler bool `json:"has_media_spoiler"`
		// Optional. Message is a shared contact, information about the contact
		Contact *Contact `json:"contact"`
		// Optional. Message is a dice with random value
		Dice *Dice `json:"dice"`
		// Optional. Message is a game, information about the game. More about games: https://core.telegram.org/bots/api#games
		Game *Game `json:"game"`
		// Optional. Message is a native poll, information about the poll
		Poll *Poll `json:"poll"`
		// Optional. Message is a venue, information about the venue.
		// For backward compatibility, when this field is set, the location field will also be set
		Venue *Venue `json:"venue"`
		// Optional. Message is a shared location, information about the location
		Location *Location `json:"location"`
		// Optional.
		// New members that were added to the group or supergroup and information about them (the bot itself may be one of these members)
		NewChatMembers []*User `json:"new_chat_members"`
		// Optional. A member was removed from the group, information about them (this member may be the bot itself)
		LeftChatMember *User `json:"left_chat_member"`
		// Optional. A chat title was changed to this value
		NewChatTitle string `json:"new_chat_title"`
		// Optional. A chat photo was change to this value
		NewChatPhoto TelegramPhoto `json:"new_chat_photo"`
		// Optional. Service message: the chat photo was deleted
		DeleteChatPhoto bool `json:"delete_chat_photo"`
		// Optional. Service message: the group has been created
		GroupChatCreated bool `json:"group_chat_created"`
		// Optional. Service message: the supergroup has been created.
		// This field can't be received in a message coming through updates, because bot can't be a member of a supergroup when it is created.
		// It can only be found in reply_to_message if someone replies to a very first message in a directly created supergroup.
		SupergroupChatCreated bool `json:"supergroup_chat_created"`
		// Optional. Service message: the channel has been created.
		// This field can't be received in a message coming through updates, because bot can't be a member of a channel when it is created.
		// It can only be found in reply_to_message if someone replies to a very first message in a channel.
		ChannelChatCreated bool `json:"channel_chat_created"`
		// Optional. Service message: auto-delete timer settings changed in the chat
		MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed"`
		// Optional. The group has been migrated to a supergroup with the specified identifier.
		// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
		// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
		MigrateToChatId int64 `json:"migrate_to_chat_id"`
		// Optional. The supergroup has been migrated from a group with the specified identifier.
		// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
		// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
		MigrateFromChatId int64 `json:"migrate_from_chat_id"`
		// Optional. Specified message was pinned.
		// Note that the Message object in this field will not contain further reply_to_message fields even if it itself is a reply.
		PinnedMessage *Message `json:"pinned_message"`
		// Optional. Message is an invoice for a payment, information about the invoice.
		// More about payments: https://core.telegram.org/bots/api#payments
		Invoice *Invoice `json:"invoice"`
		// Optional. Message is a service message about a successful payment, information about the payment.
		// More about payments: https://core.telegram.org/bots/api#payments
		SuccessfulPayment *SuccessfulPayment `json:"successful_payment"`
		// Optional. Message is a service message about a refunded payment, information about the payment.
		// More about payments: https://core.telegram.org/bots/api#payments
		RefundedPayment *RefundedPayment `json:"refunded_payment"`
		// Optional. Service message: users were shared with the bot
		UsersShared *UsersShared `json:"users_shared"`
		// Optional. Service message: a chat was shared with the bot
		ChatShared *ChatShared `json:"chat_shared"`
		// Optional. The domain name of the website on which the user has logged in.
		// More about Telegram Login: https://core.telegram.org/widgets/login
		ConnectedWebsite string `json:"connected_website"`
		// Optional.
		// Service message: the user allowed the bot to write messages after adding it to the attachment or side menu, launching a Web App from a link, or accepting an explicit request from a Web App sent by the method requestWriteAccess
		WriteAccessAllowed *WriteAccessAllowed `json:"write_access_allowed"`
		// Optional. Telegram Passport data
		PassportData *PassportData `json:"passport_data"`
		// Optional. Service message. A user in the chat triggered another user's proximity alert while sharing Live Location.
		ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered"`
		// Optional. Service message: user boosted the chat
		BoostAdded *ChatBoostAdded `json:"boost_added"`
		// Optional. Service message: chat background set
		ChatBackgroundSet *ChatBackground `json:"chat_background_set"`
		// Optional. Service message: forum topic created
		ForumTopicCreated *ForumTopicCreated `json:"forum_topic_created"`
		// Optional. Service message: forum topic edited
		ForumTopicEdited *ForumTopicEdited `json:"forum_topic_edited"`
		// Optional. Service message: forum topic closed
		ForumTopicClosed *ForumTopicClosed `json:"forum_topic_closed"`
		// Optional. Service message: forum topic reopened
		ForumTopicReopened *ForumTopicReopened `json:"forum_topic_reopened"`
		// Optional. Service message: the 'General' forum topic hidden
		GeneralForumTopicHidden *GeneralForumTopicHidden `json:"general_forum_topic_hidden"`
		// Optional. Service message: the 'General' forum topic unhidden
		GeneralForumTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden"`
		// Optional. Service message: a scheduled giveaway was created
		GiveawayCreated *GiveawayCreated `json:"giveaway_created"`
		// Optional. The message is a scheduled giveaway message
		Giveaway *Giveaway `json:"giveaway"`
		// Optional. A giveaway with public winners was completed
		GiveawayWinners *GiveawayWinners `json:"giveaway_winners"`
		// Optional. Service message: a giveaway without public winners was completed
		GiveawayCompleted *GiveawayCompleted `json:"giveaway_completed"`
		// Optional. Service message: video chat scheduled
		VideoChatScheduled *VideoChatScheduled `json:"video_chat_scheduled"`
		// Optional. Service message: video chat started
		VideoChatStarted *VideoChatStarted `json:"video_chat_started"`
		// Optional. Service message: video chat ended
		VideoChatEnded *VideoChatEnded `json:"video_chat_ended"`
		// Optional. Service message: new participants invited to a video chat
		VideoChatParticipantsInvited *VideoChatParticipantsInvited `json:"video_chat_participants_invited"`
		// Optional. Service message: data sent by a Web App
		WebAppData *WebAppData `json:"web_app_data"`
		// Optional. Inline keyboard attached to the message. login_url buttons are represented as ordinary url buttons.
		ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
		// Joint of structs, used for parsing variant interfaces.
		ForwardOrigin *MessageOriginUnmarshalJoinedForwardOrigin `json:"forward_origin,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.MessageId = inst.MessageId
	impl.MessageThreadId = inst.MessageThreadId
	impl.From = inst.From
	impl.SenderChat = inst.SenderChat
	impl.SenderBoostCount = inst.SenderBoostCount
	impl.SenderBusinessBot = inst.SenderBusinessBot
	impl.Date = inst.Date
	impl.BusinessConnectionId = inst.BusinessConnectionId
	impl.Chat = inst.Chat
	impl.IsTopicMessage = inst.IsTopicMessage
	impl.IsAutomaticForward = inst.IsAutomaticForward
	impl.ReplyToMessage = inst.ReplyToMessage
	impl.ExternalReply = inst.ExternalReply
	impl.Quote = inst.Quote
	impl.ReplyToStory = inst.ReplyToStory
	impl.ViaBot = inst.ViaBot
	impl.EditDate = inst.EditDate
	impl.HasProtectedContent = inst.HasProtectedContent
	impl.IsFromOffline = inst.IsFromOffline
	impl.MediaGroupId = inst.MediaGroupId
	impl.AuthorSignature = inst.AuthorSignature
	impl.Text = inst.Text
	impl.Entities = inst.Entities
	impl.LinkPreviewOptions = inst.LinkPreviewOptions
	impl.EffectId = inst.EffectId
	impl.Animation = inst.Animation
	impl.Audio = inst.Audio
	impl.Document = inst.Document
	impl.PaidMedia = inst.PaidMedia
	impl.Photo = inst.Photo
	impl.Sticker = inst.Sticker
	impl.Story = inst.Story
	impl.Video = inst.Video
	impl.VideoNote = inst.VideoNote
	impl.Voice = inst.Voice
	impl.Caption = inst.Caption
	impl.CaptionEntities = inst.CaptionEntities
	impl.ShowCaptionAboveMedia = inst.ShowCaptionAboveMedia
	impl.HasMediaSpoiler = inst.HasMediaSpoiler
	impl.Contact = inst.Contact
	impl.Dice = inst.Dice
	impl.Game = inst.Game
	impl.Poll = inst.Poll
	impl.Venue = inst.Venue
	impl.Location = inst.Location
	impl.NewChatMembers = inst.NewChatMembers
	impl.LeftChatMember = inst.LeftChatMember
	impl.NewChatTitle = inst.NewChatTitle
	impl.NewChatPhoto = inst.NewChatPhoto
	impl.DeleteChatPhoto = inst.DeleteChatPhoto
	impl.GroupChatCreated = inst.GroupChatCreated
	impl.SupergroupChatCreated = inst.SupergroupChatCreated
	impl.ChannelChatCreated = inst.ChannelChatCreated
	impl.MessageAutoDeleteTimerChanged = inst.MessageAutoDeleteTimerChanged
	impl.MigrateToChatId = inst.MigrateToChatId
	impl.MigrateFromChatId = inst.MigrateFromChatId
	impl.PinnedMessage = inst.PinnedMessage
	impl.Invoice = inst.Invoice
	impl.SuccessfulPayment = inst.SuccessfulPayment
	impl.RefundedPayment = inst.RefundedPayment
	impl.UsersShared = inst.UsersShared
	impl.ChatShared = inst.ChatShared
	impl.ConnectedWebsite = inst.ConnectedWebsite
	impl.WriteAccessAllowed = inst.WriteAccessAllowed
	impl.PassportData = inst.PassportData
	impl.ProximityAlertTriggered = inst.ProximityAlertTriggered
	impl.BoostAdded = inst.BoostAdded
	impl.ChatBackgroundSet = inst.ChatBackgroundSet
	impl.ForumTopicCreated = inst.ForumTopicCreated
	impl.ForumTopicEdited = inst.ForumTopicEdited
	impl.ForumTopicClosed = inst.ForumTopicClosed
	impl.ForumTopicReopened = inst.ForumTopicReopened
	impl.GeneralForumTopicHidden = inst.GeneralForumTopicHidden
	impl.GeneralForumTopicUnhidden = inst.GeneralForumTopicUnhidden
	impl.GiveawayCreated = inst.GiveawayCreated
	impl.Giveaway = inst.Giveaway
	impl.GiveawayWinners = inst.GiveawayWinners
	impl.GiveawayCompleted = inst.GiveawayCompleted
	impl.VideoChatScheduled = inst.VideoChatScheduled
	impl.VideoChatStarted = inst.VideoChatStarted
	impl.VideoChatEnded = inst.VideoChatEnded
	impl.VideoChatParticipantsInvited = inst.VideoChatParticipantsInvited
	impl.WebAppData = inst.WebAppData
	impl.ReplyMarkup = inst.ReplyMarkup
	if inst.ForwardOrigin != nil && inst.ForwardOrigin.Type == nil {
		switch *inst.ForwardOrigin.Type {
		case "channel":
			impl.ForwardOrigin = &MessageOriginChannel{
				Type:            deref(inst.ForwardOrigin.Type),
				Date:            deref(inst.ForwardOrigin.Date),
				Chat:            deref(inst.ForwardOrigin.Chat),
				MessageId:       deref(inst.ForwardOrigin.MessageId),
				AuthorSignature: deref(inst.ForwardOrigin.AuthorSignature),
			}
		case "chat":
			impl.ForwardOrigin = &MessageOriginChat{
				Type:            deref(inst.ForwardOrigin.Type),
				Date:            deref(inst.ForwardOrigin.Date),
				SenderChat:      deref(inst.ForwardOrigin.SenderChat),
				AuthorSignature: deref(inst.ForwardOrigin.AuthorSignature),
			}
		case "hidden_user":
			impl.ForwardOrigin = &MessageOriginHiddenUser{
				Type:           deref(inst.ForwardOrigin.Type),
				Date:           deref(inst.ForwardOrigin.Date),
				SenderUserName: deref(inst.ForwardOrigin.SenderUserName),
			}
		case "user":
			impl.ForwardOrigin = &MessageOriginUser{
				Type:       deref(inst.ForwardOrigin.Type),
				Date:       deref(inst.ForwardOrigin.Date),
				SenderUser: deref(inst.ForwardOrigin.SenderUser),
			}
		}
	}
	return nil
}

func (impl *MessageReactionUpdated) UnmarshalJSON(data []byte) error {
	type ReactionTypeUnmarshalJoinedOldReaction struct {
		Type          *string `json:"type"`
		Emoji         *string `json:"emoji"`
		CustomEmojiId *string `json:"custom_emoji_id"`
	}
	type ReactionTypeUnmarshalJoinedNewReaction struct {
		Type          *string `json:"type"`
		Emoji         *string `json:"emoji"`
		CustomEmojiId *string `json:"custom_emoji_id"`
	}
	type BaseInstance struct {
		// The chat containing the message the user reacted to
		Chat *Chat `json:"chat"`
		// Unique identifier of the message inside the chat
		MessageId int64 `json:"message_id"`
		// Optional. The user that changed the reaction, if the user isn't anonymous
		User *User `json:"user"`
		// Optional. The chat on behalf of which the reaction was changed, if the user is anonymous
		ActorChat *Chat `json:"actor_chat"`
		// Date of the change in Unix time
		Date int64 `json:"date"`
		// Joint of structs, used for parsing variant interfaces.
		OldReaction []*ReactionTypeUnmarshalJoinedOldReaction `json:"old_reaction"`
		NewReaction []*ReactionTypeUnmarshalJoinedNewReaction `json:"new_reaction"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Chat = inst.Chat
	impl.MessageId = inst.MessageId
	impl.User = inst.User
	impl.ActorChat = inst.ActorChat
	impl.Date = inst.Date
	if len(inst.OldReaction) != 0 {
		impl.OldReaction = []ReactionType{}
		for _, item := range inst.OldReaction {
			if item == nil || item.Type == nil {
				continue
			}
			switch *item.Type {
			case "custom_emoji":
				impl.OldReaction = append(impl.OldReaction, &ReactionTypeCustomEmoji{
					Type:          deref(item.Type),
					CustomEmojiId: deref(item.CustomEmojiId),
				})
			case "emoji":
				impl.OldReaction = append(impl.OldReaction, &ReactionTypeEmoji{
					Type:  deref(item.Type),
					Emoji: deref(item.Emoji),
				})
			case "paid":
				impl.OldReaction = append(impl.OldReaction, &ReactionTypePaid{
					Type: deref(item.Type),
				})
			}
		}
	}
	if len(inst.NewReaction) != 0 {
		impl.NewReaction = []ReactionType{}
		for _, item := range inst.NewReaction {
			if item == nil || item.Type == nil {
				continue
			}
			switch *item.Type {
			case "custom_emoji":
				impl.NewReaction = append(impl.NewReaction, &ReactionTypeCustomEmoji{
					Type:          deref(item.Type),
					CustomEmojiId: deref(item.CustomEmojiId),
				})
			case "emoji":
				impl.NewReaction = append(impl.NewReaction, &ReactionTypeEmoji{
					Type:  deref(item.Type),
					Emoji: deref(item.Emoji),
				})
			case "paid":
				impl.NewReaction = append(impl.NewReaction, &ReactionTypePaid{
					Type: deref(item.Type),
				})
			}
		}
	}
	return nil
}

func (impl *PaidMediaInfo) UnmarshalJSON(data []byte) error {
	type PaidMediaUnmarshalJoinedPaidMedia struct {
		Type     *string         `json:"type"`
		Width    *int64          `json:"width"`
		Height   *int64          `json:"height"`
		Duration *int64          `json:"duration"`
		Photo    *TelegramPhoto  `json:"photo"`
		Video    **TelegramVideo `json:"video"`
	}
	type BaseInstance struct {
		// The number of Telegram Stars that must be paid to buy access to the media
		StarCount int64 `json:"star_count"`
		// Joint of structs, used for parsing variant interfaces.
		PaidMedia []*PaidMediaUnmarshalJoinedPaidMedia `json:"paid_media"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.StarCount = inst.StarCount
	if len(inst.PaidMedia) != 0 {
		impl.PaidMedia = []PaidMedia{}
		for _, item := range inst.PaidMedia {
			if item == nil {
				continue
			}
			nonEmptyFields := []string{}
			if item.Type != nil {
				nonEmptyFields = append(nonEmptyFields, "Type")
			}
			if item.Photo != nil {
				nonEmptyFields = append(nonEmptyFields, "Photo")
			}
			if item.Video != nil {
				nonEmptyFields = append(nonEmptyFields, "Video")
			}
			if item.Width != nil {
				nonEmptyFields = append(nonEmptyFields, "Width")
			}
			if item.Height != nil {
				nonEmptyFields = append(nonEmptyFields, "Height")
			}
			if item.Duration != nil {
				nonEmptyFields = append(nonEmptyFields, "Duration")
			}
			switch {
			case containsAll([]string{"Type", "Photo"}, nonEmptyFields):
				impl.PaidMedia = append(impl.PaidMedia, &PaidMediaPhoto{
					Type:  deref(item.Type),
					Photo: deref(item.Photo),
				})
			case containsAll([]string{"Type", "Video"}, nonEmptyFields):
				impl.PaidMedia = append(impl.PaidMedia, &PaidMediaVideo{
					Type:  deref(item.Type),
					Video: deref(item.Video),
				})
			case containsAll([]string{"Type", "Width", "Height", "Duration"}, nonEmptyFields):
				impl.PaidMedia = append(impl.PaidMedia, &PaidMediaPreview{
					Type:     deref(item.Type),
					Width:    deref(item.Width),
					Height:   deref(item.Height),
					Duration: deref(item.Duration),
				})
			}
		}
	}
	return nil
}

func (impl *ReactionCount) UnmarshalJSON(data []byte) error {
	type ReactionTypeUnmarshalJoinedType struct {
		Type          *string `json:"type"`
		Emoji         *string `json:"emoji"`
		CustomEmojiId *string `json:"custom_emoji_id"`
	}
	type BaseInstance struct {
		// Number of times the reaction was added
		TotalCount int64 `json:"total_count"`
		// Joint of structs, used for parsing variant interfaces.
		Type *ReactionTypeUnmarshalJoinedType `json:"type"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.TotalCount = inst.TotalCount
	if inst.Type != nil && inst.Type.Type == nil {
		switch *inst.Type.Type {
		case "custom_emoji":
			impl.Type = &ReactionTypeCustomEmoji{
				Type:          deref(inst.Type.Type),
				CustomEmojiId: deref(inst.Type.CustomEmojiId),
			}
		case "emoji":
			impl.Type = &ReactionTypeEmoji{
				Type:  deref(inst.Type.Type),
				Emoji: deref(inst.Type.Emoji),
			}
		case "paid":
			impl.Type = &ReactionTypePaid{
				Type: deref(inst.Type.Type),
			}
		}
	}
	return nil
}

func (impl *StarTransaction) UnmarshalJSON(data []byte) error {
	type TransactionPartnerUnmarshalJoinedSource struct {
		Type               *string                 `json:"type"`
		User               **User                  `json:"user"`
		Affiliate          **AffiliateInfo         `json:"affiliate"`
		InvoicePayload     *string                 `json:"invoice_payload"`
		SubscriptionPeriod *int64                  `json:"subscription_period"`
		PaidMedia          *[]PaidMedia            `json:"paid_media"`
		PaidMediaPayload   *string                 `json:"paid_media_payload"`
		Gift               **Gift                  `json:"gift"`
		SponsorUser        **User                  `json:"sponsor_user"`
		CommissionPerMille *int64                  `json:"commission_per_mille"`
		WithdrawalState    *RevenueWithdrawalState `json:"withdrawal_state"`
		RequestCount       *int64                  `json:"request_count"`
	}
	type TransactionPartnerUnmarshalJoinedReceiver struct {
		Type               *string                 `json:"type"`
		User               **User                  `json:"user"`
		Affiliate          **AffiliateInfo         `json:"affiliate"`
		InvoicePayload     *string                 `json:"invoice_payload"`
		SubscriptionPeriod *int64                  `json:"subscription_period"`
		PaidMedia          *[]PaidMedia            `json:"paid_media"`
		PaidMediaPayload   *string                 `json:"paid_media_payload"`
		Gift               **Gift                  `json:"gift"`
		SponsorUser        **User                  `json:"sponsor_user"`
		CommissionPerMille *int64                  `json:"commission_per_mille"`
		WithdrawalState    *RevenueWithdrawalState `json:"withdrawal_state"`
		RequestCount       *int64                  `json:"request_count"`
	}
	type BaseInstance struct {
		// Unique identifier of the transaction.
		// Coincides with the identifier of the original transaction for refund transactions.
		// Coincides with SuccessfulPayment.telegram_payment_charge_id for successful incoming payments from users.
		Id string `json:"id"`
		// Integer amount of Telegram Stars transferred by the transaction
		Amount int64 `json:"amount"`
		// Optional. The number of 1/1000000000 shares of Telegram Stars transferred by the transaction; from 0 to 999999999
		NanostarAmount int64 `json:"nanostar_amount"`
		// Date the transaction was created in Unix time
		Date int64 `json:"date"`
		// Joint of structs, used for parsing variant interfaces.
		Source   *TransactionPartnerUnmarshalJoinedSource   `json:"source,omitempty"`
		Receiver *TransactionPartnerUnmarshalJoinedReceiver `json:"receiver,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Id = inst.Id
	impl.Amount = inst.Amount
	impl.NanostarAmount = inst.NanostarAmount
	impl.Date = inst.Date
	if inst.Source != nil {
		nonEmptyFields := []string{}
		if inst.Source.Type != nil {
			nonEmptyFields = append(nonEmptyFields, "Type")
		}
		if inst.Source.WithdrawalState != nil {
			nonEmptyFields = append(nonEmptyFields, "WithdrawalState")
		}
		if inst.Source.RequestCount != nil {
			nonEmptyFields = append(nonEmptyFields, "RequestCount")
		}
		if inst.Source.SponsorUser != nil {
			nonEmptyFields = append(nonEmptyFields, "SponsorUser")
		}
		if inst.Source.CommissionPerMille != nil {
			nonEmptyFields = append(nonEmptyFields, "CommissionPerMille")
		}
		if inst.Source.User != nil {
			nonEmptyFields = append(nonEmptyFields, "User")
		}
		if inst.Source.Affiliate != nil {
			nonEmptyFields = append(nonEmptyFields, "Affiliate")
		}
		if inst.Source.InvoicePayload != nil {
			nonEmptyFields = append(nonEmptyFields, "InvoicePayload")
		}
		if inst.Source.SubscriptionPeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "SubscriptionPeriod")
		}
		if inst.Source.PaidMedia != nil {
			nonEmptyFields = append(nonEmptyFields, "PaidMedia")
		}
		if inst.Source.PaidMediaPayload != nil {
			nonEmptyFields = append(nonEmptyFields, "PaidMediaPayload")
		}
		if inst.Source.Gift != nil {
			nonEmptyFields = append(nonEmptyFields, "Gift")
		}
		switch {
		case containsAll([]string{"Type"}, nonEmptyFields):
			impl.Source = &TransactionPartnerTelegramAds{
				Type: deref(inst.Source.Type),
			}
		case containsAll([]string{"Type"}, nonEmptyFields):
			impl.Source = &TransactionPartnerOther{
				Type: deref(inst.Source.Type),
			}
		case containsAll([]string{"Type", "WithdrawalState"}, nonEmptyFields):
			impl.Source = &TransactionPartnerFragment{
				Type:            deref(inst.Source.Type),
				WithdrawalState: deref(inst.Source.WithdrawalState),
			}
		case containsAll([]string{"Type", "RequestCount"}, nonEmptyFields):
			impl.Source = &TransactionPartnerTelegramApi{
				Type:         deref(inst.Source.Type),
				RequestCount: deref(inst.Source.RequestCount),
			}
		case containsAll([]string{"Type", "SponsorUser", "CommissionPerMille"}, nonEmptyFields):
			impl.Source = &TransactionPartnerAffiliateProgram{
				Type:               deref(inst.Source.Type),
				SponsorUser:        deref(inst.Source.SponsorUser),
				CommissionPerMille: deref(inst.Source.CommissionPerMille),
			}
		case containsAll([]string{"Type", "User", "Affiliate", "InvoicePayload", "SubscriptionPeriod", "PaidMedia", "PaidMediaPayload", "Gift"}, nonEmptyFields):
			impl.Source = &TransactionPartnerUser{
				Type:               deref(inst.Source.Type),
				User:               deref(inst.Source.User),
				Affiliate:          deref(inst.Source.Affiliate),
				InvoicePayload:     deref(inst.Source.InvoicePayload),
				SubscriptionPeriod: deref(inst.Source.SubscriptionPeriod),
				PaidMedia:          deref(inst.Source.PaidMedia),
				PaidMediaPayload:   deref(inst.Source.PaidMediaPayload),
				Gift:               deref(inst.Source.Gift),
			}
		}
	}
	if inst.Receiver != nil {
		nonEmptyFields := []string{}
		if inst.Receiver.Type != nil {
			nonEmptyFields = append(nonEmptyFields, "Type")
		}
		if inst.Receiver.WithdrawalState != nil {
			nonEmptyFields = append(nonEmptyFields, "WithdrawalState")
		}
		if inst.Receiver.RequestCount != nil {
			nonEmptyFields = append(nonEmptyFields, "RequestCount")
		}
		if inst.Receiver.SponsorUser != nil {
			nonEmptyFields = append(nonEmptyFields, "SponsorUser")
		}
		if inst.Receiver.CommissionPerMille != nil {
			nonEmptyFields = append(nonEmptyFields, "CommissionPerMille")
		}
		if inst.Receiver.User != nil {
			nonEmptyFields = append(nonEmptyFields, "User")
		}
		if inst.Receiver.Affiliate != nil {
			nonEmptyFields = append(nonEmptyFields, "Affiliate")
		}
		if inst.Receiver.InvoicePayload != nil {
			nonEmptyFields = append(nonEmptyFields, "InvoicePayload")
		}
		if inst.Receiver.SubscriptionPeriod != nil {
			nonEmptyFields = append(nonEmptyFields, "SubscriptionPeriod")
		}
		if inst.Receiver.PaidMedia != nil {
			nonEmptyFields = append(nonEmptyFields, "PaidMedia")
		}
		if inst.Receiver.PaidMediaPayload != nil {
			nonEmptyFields = append(nonEmptyFields, "PaidMediaPayload")
		}
		if inst.Receiver.Gift != nil {
			nonEmptyFields = append(nonEmptyFields, "Gift")
		}
		switch {
		case containsAll([]string{"Type"}, nonEmptyFields):
			impl.Receiver = &TransactionPartnerTelegramAds{
				Type: deref(inst.Receiver.Type),
			}
		case containsAll([]string{"Type"}, nonEmptyFields):
			impl.Receiver = &TransactionPartnerOther{
				Type: deref(inst.Receiver.Type),
			}
		case containsAll([]string{"Type", "WithdrawalState"}, nonEmptyFields):
			impl.Receiver = &TransactionPartnerFragment{
				Type:            deref(inst.Receiver.Type),
				WithdrawalState: deref(inst.Receiver.WithdrawalState),
			}
		case containsAll([]string{"Type", "RequestCount"}, nonEmptyFields):
			impl.Receiver = &TransactionPartnerTelegramApi{
				Type:         deref(inst.Receiver.Type),
				RequestCount: deref(inst.Receiver.RequestCount),
			}
		case containsAll([]string{"Type", "SponsorUser", "CommissionPerMille"}, nonEmptyFields):
			impl.Receiver = &TransactionPartnerAffiliateProgram{
				Type:               deref(inst.Receiver.Type),
				SponsorUser:        deref(inst.Receiver.SponsorUser),
				CommissionPerMille: deref(inst.Receiver.CommissionPerMille),
			}
		case containsAll([]string{"Type", "User", "Affiliate", "InvoicePayload", "SubscriptionPeriod", "PaidMedia", "PaidMediaPayload", "Gift"}, nonEmptyFields):
			impl.Receiver = &TransactionPartnerUser{
				Type:               deref(inst.Receiver.Type),
				User:               deref(inst.Receiver.User),
				Affiliate:          deref(inst.Receiver.Affiliate),
				InvoicePayload:     deref(inst.Receiver.InvoicePayload),
				SubscriptionPeriod: deref(inst.Receiver.SubscriptionPeriod),
				PaidMedia:          deref(inst.Receiver.PaidMedia),
				PaidMediaPayload:   deref(inst.Receiver.PaidMediaPayload),
				Gift:               deref(inst.Receiver.Gift),
			}
		}
	}
	return nil
}

func (impl *TransactionPartnerFragment) UnmarshalJSON(data []byte) error {
	type RevenueWithdrawalStateUnmarshalJoinedWithdrawalState struct {
		Type *string `json:"type"`
		Date *int64  `json:"date"`
		Url  *string `json:"url"`
	}
	type BaseInstance struct {
		// Type of the transaction partner, always "fragment"
		Type string `json:"type"`
		// Joint of structs, used for parsing variant interfaces.
		WithdrawalState *RevenueWithdrawalStateUnmarshalJoinedWithdrawalState `json:"withdrawal_state,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	if inst.WithdrawalState != nil {
		nonEmptyFields := []string{}
		if inst.WithdrawalState.Type != nil {
			nonEmptyFields = append(nonEmptyFields, "Type")
		}
		if inst.WithdrawalState.Date != nil {
			nonEmptyFields = append(nonEmptyFields, "Date")
		}
		if inst.WithdrawalState.Url != nil {
			nonEmptyFields = append(nonEmptyFields, "Url")
		}
		switch {
		case containsAll([]string{"Type"}, nonEmptyFields):
			impl.WithdrawalState = &RevenueWithdrawalStatePending{
				Type: deref(inst.WithdrawalState.Type),
			}
		case containsAll([]string{"Type"}, nonEmptyFields):
			impl.WithdrawalState = &RevenueWithdrawalStateFailed{
				Type: deref(inst.WithdrawalState.Type),
			}
		case containsAll([]string{"Type", "Date", "Url"}, nonEmptyFields):
			impl.WithdrawalState = &RevenueWithdrawalStateSucceeded{
				Type: deref(inst.WithdrawalState.Type),
				Date: deref(inst.WithdrawalState.Date),
				Url:  deref(inst.WithdrawalState.Url),
			}
		}
	}
	return nil
}

func (impl *TransactionPartnerUser) UnmarshalJSON(data []byte) error {
	type PaidMediaUnmarshalJoinedPaidMedia struct {
		Type     *string         `json:"type"`
		Width    *int64          `json:"width"`
		Height   *int64          `json:"height"`
		Duration *int64          `json:"duration"`
		Photo    *TelegramPhoto  `json:"photo"`
		Video    **TelegramVideo `json:"video"`
	}
	type BaseInstance struct {
		// Type of the transaction partner, always "user"
		Type string `json:"type"`
		// Information about the user
		User *User `json:"user"`
		// Optional. Information about the affiliate that received a commission via this transaction
		Affiliate *AffiliateInfo `json:"affiliate"`
		// Optional. Bot-specified invoice payload
		InvoicePayload string `json:"invoice_payload"`
		// Optional. The duration of the paid subscription
		SubscriptionPeriod int64 `json:"subscription_period"`
		// Optional. Bot-specified paid media payload
		PaidMediaPayload string `json:"paid_media_payload"`
		// Optional. The gift sent to the user by the bot
		Gift *Gift `json:"gift"`
		// Joint of structs, used for parsing variant interfaces.
		PaidMedia []*PaidMediaUnmarshalJoinedPaidMedia `json:"paid_media,omitempty"`
	}
	var inst BaseInstance
	if err := json.Unmarshal(data, &inst); err != nil {
		return err
	}
	impl.Type = inst.Type
	impl.User = inst.User
	impl.Affiliate = inst.Affiliate
	impl.InvoicePayload = inst.InvoicePayload
	impl.SubscriptionPeriod = inst.SubscriptionPeriod
	impl.PaidMediaPayload = inst.PaidMediaPayload
	impl.Gift = inst.Gift
	if len(inst.PaidMedia) != 0 {
		impl.PaidMedia = []PaidMedia{}
		for _, item := range inst.PaidMedia {
			if item == nil {
				continue
			}
			nonEmptyFields := []string{}
			if item.Type != nil {
				nonEmptyFields = append(nonEmptyFields, "Type")
			}
			if item.Photo != nil {
				nonEmptyFields = append(nonEmptyFields, "Photo")
			}
			if item.Video != nil {
				nonEmptyFields = append(nonEmptyFields, "Video")
			}
			if item.Width != nil {
				nonEmptyFields = append(nonEmptyFields, "Width")
			}
			if item.Height != nil {
				nonEmptyFields = append(nonEmptyFields, "Height")
			}
			if item.Duration != nil {
				nonEmptyFields = append(nonEmptyFields, "Duration")
			}
			switch {
			case containsAll([]string{"Type", "Photo"}, nonEmptyFields):
				impl.PaidMedia = append(impl.PaidMedia, &PaidMediaPhoto{
					Type:  deref(item.Type),
					Photo: deref(item.Photo),
				})
			case containsAll([]string{"Type", "Video"}, nonEmptyFields):
				impl.PaidMedia = append(impl.PaidMedia, &PaidMediaVideo{
					Type:  deref(item.Type),
					Video: deref(item.Video),
				})
			case containsAll([]string{"Type", "Width", "Height", "Duration"}, nonEmptyFields):
				impl.PaidMedia = append(impl.PaidMedia, &PaidMediaPreview{
					Type:     deref(item.Type),
					Width:    deref(item.Width),
					Height:   deref(item.Height),
					Duration: deref(item.Duration),
				})
			}
		}
	}
	return nil
}
