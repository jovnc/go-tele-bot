package utils

import (
	"go-tele-bot/internal/data"

	"github.com/go-telegram/bot/models"
)

// BuildOptionsKeyboard creates an inline keyboard for multi-select option selection
func BuildOptionsKeyboard(selected map[string]bool, callbackPrefix string) *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton

	for _, option := range data.GetFolderNames() {
		label := option
		if selected[option] {
			label = "✅ " + option
		} else {
			label = "⬜ " + option
		}

		rows = append(rows, []models.InlineKeyboardButton{
			{Text: label, CallbackData: callbackPrefix + option},
		})
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

// BuildActionButtons creates Done and Cancel action buttons
func BuildActionButtons(doneCallback, cancelCallback string) []models.InlineKeyboardButton {
	return []models.InlineKeyboardButton{
		{Text: "✓ Done", CallbackData: doneCallback},
		{Text: "✗ Cancel", CallbackData: cancelCallback},
	}
}

// BuildSelectKeyboard creates a complete multi-select keyboard with options and action buttons
func BuildSelectKeyboard(selected map[string]bool, optionPrefix, doneCallback, cancelCallback string) *models.InlineKeyboardMarkup {
	keyboard := BuildOptionsKeyboard(selected, optionPrefix)
	keyboard.InlineKeyboard = append(
		keyboard.InlineKeyboard,
		BuildActionButtons(doneCallback, cancelCallback),
	)
	return keyboard
}

