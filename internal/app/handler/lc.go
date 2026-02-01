package handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go-tele-bot/internal/app/service"
	"go-tele-bot/internal/config"
	"go-tele-bot/internal/data"
	"go-tele-bot/internal/utils"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Message text
const (
	messageLcError = "⚠️ Something went wrong"
	messageLcExit = "👋 Exited LC Mode"
)

var lcState = utils.NewManager()
var openAIService *service.OpenAIService

// GetLcState returns the LC state for a user
func GetLcState(userID int64) *utils.UserState {
	return lcState.Get(userID)
}

// getShareState is a helper to check if user is in share flow
func getShareState(userID int64) *utils.UserState {
	return GetShareState(userID)
}

// InitLcHandler initializes the LC handler with OpenAI service
func InitLcHandler() error {
	sysPrompt := data.GetLcSystemPrompt()

	var err error
	openAIService, err = service.NewOpenAIService(config.GlobalConfig.OpenAIAPIKey, sysPrompt)
	if err != nil {
		return fmt.Errorf("failed to initialize OpenAI service: %w", err)
	}
	return nil
}

// LcHandler handles the /lc command
func LcHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	topic := ""
	if update.Message.Text != "" {
		text := update.Message.Text
		if len(text) > 3 { // "/lc" is 3 characters
			topic = strings.TrimSpace(text[3:])
		}
	}

	// Check if user is already in share flow
	if shareState := getShareState(userID); shareState != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      "⚠️ Please complete or cancel your current /share flow first before starting LC mode.",
		})
		return
	}

	// Send immediate feedback to user
	statusMsg := "🤖 Generating a LeetCode problem for you..."
	if topic != "" {
		statusMsg = fmt.Sprintf("🤖 Generating a LeetCode problem about %s for you...", topic)
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      statusMsg,
	})
	if err != nil {
		log.Printf("Error sending status message: %v", err)
	}

	// Get LeetCode problem from OpenAI
	problem, conversationHistory, err := openAIService.StartLCSession(ctx, topic)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      messageLcError,
		})
		return
	}

	log.Printf("Generated problem for user %d (length: %d chars)", userID, len(problem))

	// Convert service.Message to utils.ConversationMessage
	utilsHistory := make([]utils.ConversationMessage, len(conversationHistory))
	for i, msg := range conversationHistory {
		utilsHistory[i] = utils.ConversationMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Initialize LC mode state for this user
	state := &utils.UserState{
		ConversationHistory: utilsHistory,
	}
	lcState.Set(userID, state)

	// Send the problem to user
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      problem,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("Error sending problem with Markdown to user %d: %v", userID, err)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   problem,
		})
		if err != nil {
			log.Printf("Error sending problem as plain text to user %d: %v", userID, err)
		} else {
			log.Printf("Successfully sent problem as plain text to user %d", userID)
		}
	} else {
		log.Printf("Successfully sent problem with Markdown to user %d", userID)
	}
}

// LcMessageHandler processes messages when user is in LC mode
func LcMessageHandler(ctx context.Context, b *bot.Bot, update *models.Update) bool {
	if update == nil || update.Message == nil || update.Message.Text == "" {
		return false
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	state := lcState.Get(userID)
	if state == nil {
		return false
	}

	// Convert utils.ConversationMessage to service.Message
	serviceHistory := make([]service.Message, len(state.ConversationHistory))
	for i, msg := range state.ConversationHistory {
		serviceHistory[i] = service.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Get response from OpenAI
	answer, updatedHistory, err := openAIService.AskQuestion(ctx, serviceHistory, update.Message.Text)
	if err != nil {
		log.Printf("Error asking question: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      messageLcError,
			ParseMode: models.ParseModeMarkdown,
		})
		return true
	}

	log.Printf("Generated answer for user %d (length: %d chars)", userID, len(answer))

	// Convert back to utils.ConversationMessage
	utilsHistory := make([]utils.ConversationMessage, len(updatedHistory))
	for i, msg := range updatedHistory {
		utilsHistory[i] = utils.ConversationMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Update conversation history in state
	state.ConversationHistory = utilsHistory
	lcState.Set(userID, state)

	// Send response to user
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      answer,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("Error sending answer with Markdown to user %d: %v", userID, err)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   answer,
		})
		if err != nil {
			log.Printf("Error sending answer as plain text to user %d: %v", userID, err)
		} else {
			log.Printf("Successfully sent answer as plain text to user %d", userID)
		}
	} else {
		log.Printf("Successfully sent answer with Markdown to user %d", userID)
	}

	return true
}

func ExitLcHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// Check if user is in LC mode
	state := lcState.Get(userID)
	if state == nil {
		return
	}

	// Clear LC mode state
	lcState.Delete(userID)

	// Send confirmation message
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      messageLcExit,
	})
}
