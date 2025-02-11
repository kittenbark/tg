package main

import (
	"encoding/json"
	"io"
	"os"
)

var discriminators = map[string]*discriminator{
	"BackgroundFill": {
		property: "type",
		mapping: map[string]string{
			"freeform_gradient": "BackgroundFillFreeformGradient",
			"gradient":          "BackgroundFillGradient",
			"solid":             "BackgroundFillSolid",
		},
	},
	"BackgroundType": {
		property: "type",
		mapping: map[string]string{
			"chat_theme": "BackgroundTypeChatTheme",
			"fill":       "BackgroundTypeFill",
			"pattern":    "BackgroundTypePattern",
			"wallpaper":  "BackgroundTypeWallpaper",
		},
	},
	"BotCommandScope": {
		property: "type",
		mapping: map[string]string{
			"all_chat_administrators": "BotCommandScopeAllChatAdministrators",
			"all_group_chats":         "BotCommandScopeAllGroupChats",
			"all_private_chats":       "BotCommandScopeAllPrivateChats",
			"chat":                    "BotCommandScopeChat",
			"chat_administrators":     "BotCommandScopeChatAdministrators",
			"chat_member":             "BotCommandScopeChatMember",
		},
	},
	"ChatBoostSource": {
		property: "source",
		mapping: map[string]string{
			"gift_code": "ChatBoostSourceGiftCode",
			"giveaway":  "ChatBoostSourceGiveaway",
			"premium":   "ChatBoostSourcePremium",
		},
	},
	"ChatMember": {
		property: "status",
		mapping: map[string]string{
			"administrator": "ChatMemberAdministrator",
			"creator":       "ChatMemberOwner",
			"kicked":        "ChatMemberBanned",
			"left":          "ChatMemberLeft",
			"member":        "ChatMemberMember",
			"restricted":    "ChatMemberRestricted",
		},
	},
	"InputMedia": {
		property: "type",
		mapping: map[string]string{
			"animation": "InputMediaAnimation",
			"audio":     "InputMediaAudio",
			"document":  "InputMediaDocument",
			"photo":     "InputMediaPhoto",
			"video":     "InputMediaVideo",
		},
	},
	"MenuButton": {
		property: "type",
		mapping: map[string]string{
			"commands": "MenuButtonCommands",
			"default":  "MenuButtonDefault",
			"web_app":  "MenuButtonWebApp",
		},
	},
	"MessageOrigin": {
		property: "type",
		mapping: map[string]string{
			"channel":     "MessageOriginChannel",
			"chat":        "MessageOriginChat",
			"hidden_user": "MessageOriginHiddenUser",
			"user":        "MessageOriginUser",
		},
	},
	"PassportElementError": {
		property: "source",
		mapping: map[string]string{
			"data":              "PassportElementErrorDataField",
			"file":              "PassportElementErrorFile",
			"files":             "PassportElementErrorFiles",
			"front_side":        "PassportElementErrorFrontSide",
			"reverse_side":      "PassportElementErrorReverseSide",
			"selfie":            "PassportElementErrorSelfie",
			"translation_file":  "PassportElementErrorTranslationFile",
			"translation_files": "PassportElementErrorTranslationFiles",
			"unspecified":       "PassportElementErrorUnspecified",
		},
	},
	"ReactionType": {
		property: "type",
		mapping: map[string]string{
			"custom_emoji": "ReactionTypeCustomEmoji",
			"emoji":        "ReactionTypeEmoji",
			"paid":         "ReactionTypePaid",
		},
	},
}

func Read[T any](filename string) (*T, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var schema T
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}

	return &schema, nil
}

type Field struct {
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	Required    bool     `json:"required"`
	Description string   `json:"description"`
}

type DynamicSchema struct {
	Version     string             `json:"version"`
	ReleaseDate string             `json:"release_date"`
	Changelog   string             `json:"changelog"`
	Methods     map[string]*Method `json:"methods"`
	Types       map[string]*Type   `json:"types"`
}
