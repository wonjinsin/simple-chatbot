package shared

import (
	"context"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

type CleanMarkdownJSONParser[T any] struct {
	baseParser schema.MessageParser[T]
}

func NewJSONParserLambda[T any]() *compose.Lambda {
	parser := &CleanMarkdownJSONParser[T]{
		baseParser: schema.NewMessageJSONParser[T](&schema.MessageJSONParseConfig{
			ParseFrom: schema.MessageParseFromContent,
		}),
	}
	return compose.MessageParser(parser)
}

// Parse cleans markdown code blocks from message content and then parses it
func (p *CleanMarkdownJSONParser[T]) Parse(ctx context.Context, msg *schema.Message) (T, error) {
	var result T
	if msg == nil {
		return result, nil
	}

	// Clean markdown code blocks from content
	content := p.cleanMarkdown(msg.Content)

	// Create a temporary message with cleaned content
	cleanedMsg := &schema.Message{
		Role:      msg.Role,
		Content:   content,
		ToolCalls: msg.ToolCalls,
	}

	// Use base parser to parse the cleaned message
	parsed, err := p.baseParser.Parse(ctx, cleanedMsg)
	if err != nil {
		return result, errors.Wrap(err, "failed to parse message")
	}
	return parsed, nil
}

// cleanMarkdown removes markdown code blocks and extracts JSON
func (p *CleanMarkdownJSONParser[T]) cleanMarkdown(content string) string {
	// Remove markdown code blocks (```json ... ``` or ``` ... ```)
	markdownCodeBlockRegex := regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")
	matches := markdownCodeBlockRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		// Extract JSON from code block
		return strings.TrimSpace(matches[1])
	}

	// If no code block, try to find JSON object in the content
	// Find the first { and match until the corresponding }
	startIdx := strings.Index(content, "{")
	if startIdx != -1 {
		braceCount := 0
		for i := startIdx; i < len(content); i++ {
			switch content[i] {
			case '{':
				braceCount++
			case '}':
				braceCount--
				if braceCount == 0 {
					return strings.TrimSpace(content[startIdx : i+1])
				}
			}
		}
	}

	// Return trimmed content if no JSON found
	return strings.TrimSpace(content)
}
