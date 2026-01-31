package langchain

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/wonjinsin/simple-chatbot/internal/repository/langchain/shared"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

type AnswerRefineRepo struct {
	llm *openai.ChatModel
}

// NewAnswerRefineRepo creates a new answer refine repository
func NewAnswerRefineRepo(llm *openai.ChatModel) *AnswerRefineRepo {
	return &AnswerRefineRepo{llm: llm}
}

func (r *AnswerRefineRepo) RefineAnswer(
	ctx context.Context,
	contextStr string,
) (string, error) {
	// Create a prompt template
	template := prompt.FromMessages(
		schema.GoTemplate,
		schema.SystemMessage(
			"You are a JSON-only response assistant. You MUST respond with ONLY valid JSON. The response must be a single JSON object with an 'answer' field containing a plain string value. Do NOT use markdown code blocks, backticks, or any formatting. Do NOT nest JSON objects. Return ONLY the raw JSON object.",
		),
		schema.SystemMessage(
			`You are a helpful assistant that answers questions based on the provided context.
			Use the context information to provide accurate and relevant answers.
			If the context doesn't contain enough information to answer the question, say so honestly.`,
		),
		schema.UserMessage(
			`Context information:
			{{.context}}

			Please answer the question based on the context provided above.
			Return your response as a JSON object with this exact structure: {"answer": "your answer here"}.
			The answer field must contain a plain string, not nested JSON.`,
		),
	)

	// Render the template with data
	variables := map[string]any{
		"context": contextStr,
	}

	// Create parser for json with markdown cleaning
	type JSONResponse struct {
		Answer string `json:"answer"`
	}

	// JSON parser that cleans markdown before parsing
	jsonParserLambda := shared.NewJSONParserLambda[*JSONResponse]()

	chain, err := compose.NewChain[map[string]any, *JSONResponse]().
		AppendChatTemplate(template).
		AppendChatModel(r.llm).
		AppendLambda(jsonParserLambda).
		Compile(ctx)

	if err != nil {
		return "", errors.Wrap(err, "failed to compile chain")
	}

	result, err := chain.Invoke(ctx, variables)
	if err != nil {
		return "", errors.Wrap(err, "failed to invoke chain")
	}
	return result.Answer, nil
}
