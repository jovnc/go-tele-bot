package service

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// OpenAIService handles OpenAI API interactions
type OpenAIService struct {
	client       openai.Client
	systemPrompt string
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(apiKey string, systemPrompt string) (*OpenAIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	return &OpenAIService{
		client:       client,
		systemPrompt: systemPrompt,
	}, nil
}

// Message represents a conversation message
type Message struct {
	Role    string
	Content string
}

// StartLCSession initiates a new LC session and returns a random LeetCode problem
func (s *OpenAIService) StartLCSession(ctx context.Context, topic string) (string, []Message, error) {
	var userMessage string
	if topic != "" {
		userMessage = fmt.Sprintf("Give me a LeetCode problem about %s to practice.", topic)
	} else {
		userMessage = "Give me a random LeetCode problem to practice."
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(s.systemPrompt),
		openai.UserMessage(userMessage),
	}

	params := openai.ChatCompletionNewParams{
		Model:       openai.ChatModelGPT5_2_2025_12_11,
		Messages:    messages,
		Temperature: openai.Float(1.0),
	}

	response, err := s.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get LeetCode problem: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", nil, fmt.Errorf("no response from OpenAI")
	}

	problem := response.Choices[0].Message.Content

	// Initialize conversation history
	conversationHistory := []Message{
		{Role: "system", Content: s.systemPrompt},
		{Role: "user", Content: userMessage},
		{Role: "assistant", Content: problem},
	}

	return problem, conversationHistory, nil
}

// AskQuestion sends a question and returns the assistant's response
func (s *OpenAIService) AskQuestion(ctx context.Context, conversationHistory []Message, question string) (string, []Message, error) {
	// Convert conversation history to OpenAI format
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(conversationHistory)+1)
	for _, msg := range conversationHistory {
		switch msg.Role {
		case "system":
			messages = append(messages, openai.SystemMessage(msg.Content))
		case "user":
			messages = append(messages, openai.UserMessage(msg.Content))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(msg.Content))
		}
	}

	// Add new user question
	messages = append(messages, openai.UserMessage(question))

	params := openai.ChatCompletionNewParams{
		Model:       openai.ChatModelGPT5_2ChatLatest,
		Messages:    messages,
		Temperature: openai.Float(1.0),
	}

	response, err := s.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", conversationHistory, fmt.Errorf("failed to get response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", conversationHistory, fmt.Errorf("no response from OpenAI")
	}

	answer := response.Choices[0].Message.Content

	// Update conversation history
	updatedHistory := append(conversationHistory, Message{Role: "user", Content: question})
	updatedHistory = append(updatedHistory, Message{Role: "assistant", Content: answer})

	return answer, updatedHistory, nil
}
