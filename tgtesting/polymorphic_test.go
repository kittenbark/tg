package tgtesting

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kittenbark/tg"
)

// Regression tests for two related decode bugs:
//  1. GetChatMember/GetChatAdministrators/GetChatMenuButton return a bare
//     polymorphic interface (ChatMember/[]ChatMember/MenuButton); encoding/json
//     can't pick a concrete type for a bare interface on its own, so these
//     calls used to fail with "json: cannot unmarshal object into Go value of
//     type tg.ChatMember".
//  2. The discriminator switch for a *singular* polymorphic interface field
//     (e.g. ChatMemberUpdated.OldChatMember/NewChatMember) had an inverted
//     nil-check, so it never ran for a real payload and silently left the
//     field nil.

func TestGetChatMemberPolymorphic(t *testing.T) {
	t.Parallel()

	MakeTestOK(
		chatMemberPtr(&tg.ChatMemberAdministrator{
			Status:        "administrator",
			User:          &tg.User{Id: 42, FirstName: "Ann"},
			CanBeEdited:   true,
			CanManageChat: true,
		}),
		func(ctx context.Context) (*tg.ChatMember, error) {
			result, err := tg.GetChatMember(ctx, 123456, 42)
			return &result, err
		},
		Stub{Url: "/getChatMember", Result: StubResultOK(http.StatusOK, &tg.ChatMemberAdministrator{
			Status:        "administrator",
			User:          &tg.User{Id: 42, FirstName: "Ann"},
			CanBeEdited:   true,
			CanManageChat: true,
		})},
	)(t)
}

func TestGetChatAdministratorsPolymorphic(t *testing.T) {
	t.Parallel()

	expected := []tg.ChatMember{
		&tg.ChatMemberOwner{Status: "creator", User: &tg.User{Id: 1, FirstName: "Owner"}},
		&tg.ChatMemberMember{Status: "member", User: &tg.User{Id: 2, FirstName: "Member"}},
	}

	MakeTestOK(
		&expected,
		func(ctx context.Context) (*[]tg.ChatMember, error) {
			result, err := tg.GetChatAdministrators(ctx, 123456)
			return &result, err
		},
		Stub{Url: "/getChatAdministrators", Result: StubResultOK(http.StatusOK, []tg.ChatMember{
			&tg.ChatMemberOwner{Status: "creator", User: &tg.User{Id: 1, FirstName: "Owner"}},
			&tg.ChatMemberMember{Status: "member", User: &tg.User{Id: 2, FirstName: "Member"}},
		})},
	)(t)
}

func TestGetChatMenuButtonPolymorphic(t *testing.T) {
	t.Parallel()

	MakeTestOK(
		menuButtonPtr(&tg.MenuButtonCommands{Type: "commands"}),
		func(ctx context.Context) (*tg.MenuButton, error) {
			result, err := tg.GetChatMenuButton(ctx)
			return &result, err
		},
		Stub{Url: "/getChatMenuButton", Result: StubResultOK(http.StatusOK, &tg.MenuButtonCommands{Type: "commands"})},
	)(t)
}

// TestChatMemberUpdatedConcreteTypes guards the inverted-nil-check bug
// directly: before the fix, OldChatMember/NewChatMember always came back nil
// for a real payload (status present), regardless of this test's stub setup.
func TestChatMemberUpdatedConcreteTypes(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"chat": {"id": 123456, "type": "group"},
		"from": {"id": 1, "first_name": "Admin"},
		"date": 1700000000,
		"old_chat_member": {"status": "member", "user": {"id": 2, "first_name": "Bob"}},
		"new_chat_member": {"status": "administrator", "user": {"id": 2, "first_name": "Bob"}, "can_be_edited": true, "can_manage_chat": true}
	}`)

	var upd tg.ChatMemberUpdated
	require.NoError(t, json.Unmarshal(data, &upd))

	oldMember, ok := upd.OldChatMember.(*tg.ChatMemberMember)
	require.True(t, ok, "OldChatMember did not decode to *tg.ChatMemberMember")
	require.Equal(t, int64(2), oldMember.User.Id)

	newMember, ok := upd.NewChatMember.(*tg.ChatMemberAdministrator)
	require.True(t, ok, "NewChatMember did not decode to *tg.ChatMemberAdministrator")
	require.Equal(t, int64(2), newMember.User.Id)
	require.True(t, newMember.CanManageChat)
}

func chatMemberPtr(v tg.ChatMember) *tg.ChatMember { return &v }
func menuButtonPtr(v tg.MenuButton) *tg.MenuButton { return &v }
