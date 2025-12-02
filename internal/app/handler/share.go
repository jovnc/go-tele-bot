package handler

import (
	"context"
	"fmt"
	"strings"

	"go-tele-bot/internal/app/service"
	"go-tele-bot/internal/data"
	"go-tele-bot/internal/utils"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Callback data prefixes
const (
	callbackShareOptPrefix = "share_opt_"
	callbackShareDone      = "share_done"
	callbackShareCancel    = "share_cancel"
)

// State steps
const (
	stepSelectOptions = "select_options"
	stepEnterEmail    = "enter_email"
	stepDone          = "done"
)

// Message text
const (
	messageShare       = "📤 <b>Share Folders</b>\n\n" +
		"<b>Step 1 of 2:</b> Select items to share\n" +
		"Tap to toggle selection:"
	messageShareDone   = "✅ <b>Great choice!</b>\n\n" +
		"<b>Step 2 of 2:</b> Enter recipient email (with or without @gmail.com)\n\n" +
		"📧 Type the email address below:"
	messageShareCancel = "❌ <b>Cancelled</b>\n\n" +
		"No folders were shared.\n" +
		"Use /share to start again."
	messageShareFailed = "⚠️ <b>Something went wrong</b>\n\n" +
		"We couldn't share the folders.\n" +
		"Please try again with /share"
	messageShareSummary = "🎉 <b>Shared Successfully!</b>\n\n" +
		"<b>Items shared:</b>\n%s\n\n" +
		"<b>Shared with:</b> <code>%s</code>\n\n" +
		"✅ The recipient will receive access shortly."
)

var shareState = utils.NewManager()

// shareCallbackContext holds common data for share callback handlers
type shareCallbackContext struct {
	ctx        context.Context
	b          *bot.Bot
	callbackID string
	userID     int64
	chatID     int64
	messageID  int
	state      *utils.UserState
}

// ShareCommandHandler handles the /share command
func ShareCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// Initialize state for this user
	state := &utils.UserState{
		Step:     stepSelectOptions,
		Selected: make(map[string]bool),
	}
	shareState.Set(userID, state)

	// Send options keyboard
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        messageShare,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: utils.BuildSelectKeyboard(make(map[string]bool), callbackShareOptPrefix, callbackShareDone, callbackShareCancel),
	})
}

// ShareCallbackHandler handles button clicks
func ShareCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	messageID := update.CallbackQuery.Message.Message.ID
	callbackID := update.CallbackQuery.ID
	data := update.CallbackQuery.Data

	state := shareState.Get(userID)
	if state == nil || state.Step != stepSelectOptions {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackID,
			Text:            "Please start with /share",
		})
		return
	}

	cc := &shareCallbackContext{
		ctx:        ctx,
		b:          b,
		callbackID: callbackID,
		userID:     userID,
		chatID:     chatID,
		messageID:  messageID,
		state:      state,
	}

	switch {
	case data == callbackShareDone:
		handleShareDone(cc)
	case data == callbackShareCancel:
		handleShareCancel(cc)
	case strings.HasPrefix(data, callbackShareOptPrefix):
		handleShareOptionToggle(cc, strings.TrimPrefix(data, callbackShareOptPrefix))
	}
}

// ShareEmailInputHandler handles email input for the share flow, returns true if handled
func ShareEmailInputHandler(ctx context.Context, b *bot.Bot, update *models.Update) bool {
	if update == nil || update.Message == nil || update.Message.Text == "" {
		return false
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	state := shareState.Get(userID)
	if state == nil || state.Step != stepEnterEmail {
		return false
	}

	email := update.Message.Text
	state.Email = email
	state.Step = stepDone

	// Share the selected items to the email
	shareService, err := service.NewShareService(ctx)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   messageShareFailed,
		})
		shareState.Delete(userID)
		return true
	}

	err = shareService.ShareFolders(ctx, state.Selected, email)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   messageShareFailed,
		})
		shareState.Delete(userID)
		return true
	}

	// Build summary
	var selectedItems []string
	for _, opt := range data.GetFolderNames() {
		if state.Selected[opt] {
			selectedItems = append(selectedItems, "• "+opt)
		}
	}

	summary := fmt.Sprintf(messageShareSummary, strings.Join(selectedItems, "\n"), email)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      summary,
		ParseMode: models.ParseModeHTML,
	})

	shareState.Delete(userID)
	return true
}

// Helper functions

// handleShareDone handles the "Done" button click
func handleShareDone(cc *shareCallbackContext) {
	if len(cc.state.Selected) == 0 {
		cc.b.AnswerCallbackQuery(cc.ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cc.callbackID,
			Text:            "Please select at least one option!",
			ShowAlert:       true,
		})
		return
	}

	cc.state.Step = stepEnterEmail
	shareState.Set(cc.userID, cc.state)

	go cc.b.AnswerCallbackQuery(cc.ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: cc.callbackID,
	})

	cc.b.EditMessageText(cc.ctx, &bot.EditMessageTextParams{
		ChatID:    cc.chatID,
		MessageID: cc.messageID,
		Text:      messageShareDone,
		ParseMode: models.ParseModeHTML,
	})
}

// handleShareCancel handles the "Cancel" button click
func handleShareCancel(cc *shareCallbackContext) {
	shareState.Delete(cc.userID)

	go cc.b.AnswerCallbackQuery(cc.ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: cc.callbackID,
		Text:            "Cancelled",
	})

	cc.b.EditMessageText(cc.ctx, &bot.EditMessageTextParams{
		ChatID:    cc.chatID,
		MessageID: cc.messageID,
		Text:      messageShareCancel,
		ParseMode: models.ParseModeHTML,
	})
}

// handleShareOptionToggle handles the option toggle button click
func handleShareOptionToggle(cc *shareCallbackContext, option string) {
	if cc.state.Selected[option] {
		delete(cc.state.Selected, option)
	} else {
		cc.state.Selected[option] = true
	}
	shareState.Set(cc.userID, cc.state)

	go cc.b.AnswerCallbackQuery(cc.ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: cc.callbackID,
	})

	cc.b.EditMessageReplyMarkup(cc.ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      cc.chatID,
		MessageID:   cc.messageID,
		ReplyMarkup: utils.BuildSelectKeyboard(cc.state.Selected, callbackShareOptPrefix, callbackShareDone, callbackShareCancel),
	})
}
