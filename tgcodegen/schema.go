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
	// RichText itself isn't a generated interface (see gen_types.go's Specials["RichText"]), but this
	// entry is still consulted by withDefaultByDiscriminators to auto-fill the "type" discriminant when
	// constructing a RichText* literal directly.
	"RichText": {
		property: "type",
		mapping: map[string]string{
			"anchor":                  "RichTextAnchor",
			"anchor_link":             "RichTextAnchorLink",
			"bank_card_number":        "RichTextBankCardNumber",
			"bold":                    "RichTextBold",
			"bot_command":             "RichTextBotCommand",
			"button":                  "RichTextButton",
			"cashtag":                 "RichTextCashtag",
			"code":                    "RichTextCode",
			"custom_emoji":            "RichTextCustomEmoji",
			"date_time":               "RichTextDateTime",
			"email_address":           "RichTextEmailAddress",
			"hashtag":                 "RichTextHashtag",
			"italic":                  "RichTextItalic",
			"marked":                  "RichTextMarked",
			"mathematical_expression": "RichTextMathematicalExpression",
			"mention":                 "RichTextMention",
			"phone_number":            "RichTextPhoneNumber",
			"reference":               "RichTextReference",
			"reference_link":          "RichTextReferenceLink",
			"spoiler":                 "RichTextSpoiler",
			"strikethrough":           "RichTextStrikethrough",
			"subscript":               "RichTextSubscript",
			"superscript":             "RichTextSuperscript",
			"text_mention":            "RichTextTextMention",
			"underline":               "RichTextUnderline",
			"url":                     "RichTextUrl",
		},
	},
	"RichBlock": {
		property: "type",
		mapping: map[string]string{
			"anchor":                  "RichBlockAnchor",
			"animation":               "RichBlockAnimation",
			"audio":                   "RichBlockAudio",
			"blockquote":              "RichBlockBlockQuotation",
			"buttons":                 "RichBlockButtons",
			"collage":                 "RichBlockCollage",
			"details":                 "RichBlockDetails",
			"divider":                 "RichBlockDivider",
			"document":                "RichBlockDocument",
			"expandable_blockquote":   "RichBlockExpandableBlockQuotation",
			"footer":                  "RichBlockFooter",
			"heading":                 "RichBlockSectionHeading",
			"list":                    "RichBlockList",
			"map":                     "RichBlockMap",
			"mathematical_expression": "RichBlockMathematicalExpression",
			"paragraph":               "RichBlockParagraph",
			"photo":                   "RichBlockPhoto",
			"pre":                     "RichBlockPreformatted",
			"pullquote":               "RichBlockPullQuotation",
			"slideshow":               "RichBlockSlideshow",
			"table":                   "RichBlockTable",
			"thinking":                "RichBlockThinking",
			"video":                   "RichBlockVideo",
			"voice_note":              "RichBlockVoiceNote",
		},
	},
	"InputRichBlock": {
		property: "type",
		mapping: map[string]string{
			"anchor":                  "InputRichBlockAnchor",
			"animation":               "InputRichBlockAnimation",
			"audio":                   "InputRichBlockAudio",
			"blockquote":              "InputRichBlockBlockQuotation",
			"buttons":                 "InputRichBlockButtons",
			"collage":                 "InputRichBlockCollage",
			"details":                 "InputRichBlockDetails",
			"divider":                 "InputRichBlockDivider",
			"document":                "InputRichBlockDocument",
			"expandable_blockquote":   "InputRichBlockExpandableBlockQuotation",
			"footer":                  "InputRichBlockFooter",
			"heading":                 "InputRichBlockSectionHeading",
			"list":                    "InputRichBlockList",
			"map":                     "InputRichBlockMap",
			"mathematical_expression": "InputRichBlockMathematicalExpression",
			"paragraph":               "InputRichBlockParagraph",
			"photo":                   "InputRichBlockPhoto",
			"pre":                     "InputRichBlockPreformatted",
			"pullquote":               "InputRichBlockPullQuotation",
			"slideshow":               "InputRichBlockSlideshow",
			"table":                   "InputRichBlockTable",
			"thinking":                "InputRichBlockThinking",
			"video":                   "InputRichBlockVideo",
			"voice_note":              "InputRichBlockVoiceNote",
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
