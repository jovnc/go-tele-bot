package handler

import (
	"context"
	"fmt"
	"strings"

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

var shareState = utils.NewManager()

// shareCallbackContext holds common data for share callback handlers
type shareCallbackContext struct {
	ctx         context.Context
	b           *bot.Bot
	callbackID  string
	userID      int64
	chatID      int64
	messageID   int
	state       *utils.UserState
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
		Text:        "📤 <b>Share</b>\n\nSelect what you want to share (tap to toggle):",
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
	state := shareState.Get(userID)
	if state == nil || state.Step != stepSelectOptions {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Please start with /share",
		})
		return
	}

	cc := &shareCallbackContext{
		ctx:        ctx,
		b:          b,
		callbackID: update.CallbackQuery.ID,
		userID:     userID,
		chatID:     update.CallbackQuery.Message.Message.Chat.ID,
		messageID:  update.CallbackQuery.Message.Message.ID,
		state:      state,
	}

	data := update.CallbackQuery.Data
	switch {
	case data == callbackShareDone:
		handleShareDone(cc)
	case data == callbackShareCancel:
		handleShareCancel(cc)
	case strings.HasPrefix(data, callbackShareOptPrefix):
		handleShareOptionToggle(cc, strings.TrimPrefix(data, callbackShareOptPrefix))
	}
}

// ShareEmailHandler handles email input
func ShareEmailHandler(ctx context.Context, b *bot.Bot, update *models.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
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
	
	// TODO: share the selected items to the email

	// TODO: send email to the user

	// Build summary
	var selectedItems []string
	for _, opt := range data.GetFolderNames() {
		if state.Selected[opt] {
			selectedItems = append(selectedItems, "• "+opt)
		}
	}

	summary := fmt.Sprintf(
		"🎉 <b>Share Summary</b>\n\n"+
			"<b>Selected items:</b>\n%s\n\n"+
			"<b>Email:</b> <code>%s</code>\n\n"+
			"Thank you! Your share request has been received.",
		strings.Join(selectedItems, "\n"),
		email,
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      summary,
		ParseMode: models.ParseModeHTML,
	})

	// Clean up state
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
		Text:      "✅ Options selected!\n\n📧 Now please enter your email address:",
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
		Text:      "❌ Share cancelled.",
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
