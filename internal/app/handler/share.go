package handler

import (
	"context"
	"fmt"
	"strings"

	"go-tele-bot/internal/data"
	"go-tele-bot/internal/state"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var shareState = state.NewManager()

// ShareCommandHandler handles the /share command
func ShareCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// Initialize state for this user
	shareState.Set(userID, &state.UserState{
		Step:     "select_options",
		Selected: make(map[string]bool),
	})

	// Send options keyboard
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "📤 <b>Share</b>\n\nSelect what you want to share (tap to toggle):",
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: buildOptionsKeyboard(make(map[string]bool)),
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
	data := update.CallbackQuery.Data

	state := shareState.Get(userID)
	if state == nil || state.Step != "select_options" {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Please start with /share",
		})
		return
	}

	// Handle "Done" button
	if data == "share_done" {
		if len(state.Selected) == 0 {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "Please select at least one option!",
				ShowAlert:       true,
			})
			return
		}

		// Move to email step
		state.Step = "enter_email"
		shareState.Set(userID, state)

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: messageID,
			Text:      "✅ Options selected!\n\n📧 Now please enter your email address:",
		})
		return
	}

	// Handle "Cancel" button
	if data == "share_cancel" {
		shareState.Delete(userID)

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Cancelled",
		})

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: messageID,
			Text:      "❌ Share cancelled.",
		})
		return
	}

	// Toggle option selection
	if strings.HasPrefix(data, "share_opt_") {
		option := strings.TrimPrefix(data, "share_opt_")
		
		if state.Selected[option] {
			delete(state.Selected, option)
		} else {
			state.Selected[option] = true
		}
		shareState.Set(userID, state)

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: buildOptionsKeyboard(state.Selected),
		})
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
	if state == nil || state.Step != "enter_email" {
		return false
	}

	email := update.Message.Text
	state.Email = email
	state.Step = "done"

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

// buildOptionsKeyboard creates the inline keyboard for option selection
func buildOptionsKeyboard(selected map[string]bool) *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton

	for _, option := range data.GetFolderNames() {
		label := option
		if selected[option] {
			label = "✅ " + option
		} else {
			label = "⬜ " + option
		}

		rows = append(rows, []models.InlineKeyboardButton{
			{Text: label, CallbackData: "share_opt_" + option},
		})
	}

	// Add Done and Cancel buttons
	rows = append(rows, []models.InlineKeyboardButton{
		{Text: "✓ Done", CallbackData: "share_done"},
		{Text: "✗ Cancel", CallbackData: "share_cancel"},
	})

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

