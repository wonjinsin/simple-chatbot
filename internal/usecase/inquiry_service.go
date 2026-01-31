package usecase

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/internal/domain"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
	"github.com/wonjinsin/simple-chatbot/pkg/file"
	"github.com/wonjinsin/simple-chatbot/pkg/utils"
)

const (
	batchSize       = 50 // Batch size for embedding API calls
	similarityLimit = 3  // Number of similar entries to retrieve
)

type InquiryServiceImpl struct {
	embeddingRepo    repository.EmbeddingRepository
	knowledgeRepo    repository.InquiryKnowledgeRepository
	answerRefineRepo repository.AnswerRefineRepository
}

func NewInquiryServiceImpl(
	embeddingRepo repository.EmbeddingRepository,
	knowledgeRepo repository.InquiryKnowledgeRepository,
	answerRefineRepo repository.AnswerRefineRepository,
) *InquiryServiceImpl {
	return &InquiryServiceImpl{
		embeddingRepo:    embeddingRepo,
		knowledgeRepo:    knowledgeRepo,
		answerRefineRepo: answerRefineRepo,
	}
}

// EmbedInquiryOrigins reads CSV data, generates embeddings, and saves to database
func (s *InquiryServiceImpl) EmbedInquiryOrigins(ctx context.Context) error {
	// Step 1: Read CSV file
	csvRows, err := file.ReadCSVToMapArray("mock_data/data_set.csv")
	if err != nil {
		return errors.Wrap(err, "failed to read inquiry origins", constants.InternalError)
	}

	// Step 2: Convert CSV rows to domain objects (without embeddings)
	knowledgeItems, err := domain.NewInquiryKnowledgeFromCSVs(csvRows)
	if err != nil {
		return errors.Wrap(
			err,
			"failed to convert CSV rows to domain objects",
			constants.InternalError,
		)
	}

	// Step 3: Process in batches
	i := 0
	for batch := range slices.Chunk(knowledgeItems, batchSize) {
		instructions := batch.Instructions()

		embeddings, err := s.embeddingRepo.EmbedStrings(ctx, instructions)
		if err != nil {
			return errors.Wrap(
				err,
				fmt.Sprintf("failed to generate embeddings for batch %d", i),
				constants.InternalError,
			)
		}
		batch.SetEmbeddings(embeddings)

		if err := s.knowledgeRepo.BatchSaveInquiryKnowledge(ctx, batch); err != nil {
			return errors.Wrap(
				err,
				fmt.Sprintf("failed to save inquiry knowledge for batch %d", i),
				constants.InternalError,
			)
		}

		i++
	}

	return nil
}

// Ask answers a user question by finding similar inquiry knowledge and refining the answer
func (s *InquiryServiceImpl) Ask(
	ctx context.Context,
	msg string,
) (string, error) {
	// Step 1: Validate input message
	msg = strings.TrimSpace(msg)
	if utils.IsEmptyOrWhitespace(msg) {
		return "", errors.New(
			constants.InvalidParameter,
			"question cannot be empty",
			nil,
		)
	}

	// Step 2: Generate embedding for the user's question
	embedding, err := s.embeddingRepo.EmbedString(ctx, msg)
	if err != nil {
		return "", errors.Wrap(
			err,
			"failed to generate embedding for question",
			constants.InternalError,
		)
	}

	if embedding.IsEmpty() {
		return "", errors.New(
			constants.InternalError,
			"embedding generation returned empty result",
			nil,
		)
	}

	// Step 3: Find similar inquiry knowledge entries with similarity scores
	similarEntries, err := s.knowledgeRepo.FindSimilars(ctx, embedding, similarityLimit)
	if err != nil {
		return "", errors.Wrap(
			err,
			"failed to find similar inquiry knowledge",
			constants.InternalError,
		)
	}

	// Step 4: Build context from similar entries
	var contextBuilder strings.Builder
	for _, entry := range similarEntries {
		contextBuilder.WriteString(fmt.Sprintf("Question: %s\n", entry.Knowledge.Instruction))
		contextBuilder.WriteString(fmt.Sprintf("Answer: %s\n", entry.Knowledge.Response))
		contextBuilder.WriteString(fmt.Sprintf("Similarity: %.4f\n\n", entry.SimilarityScore))
	}

	// Step 5: Refine answer using LLM with context
	refinedAnswer, err := s.answerRefineRepo.RefineAnswer(ctx, contextBuilder.String())
	if err != nil {
		return "", errors.Wrap(
			err,
			"failed to refine answer",
			constants.InternalError,
		)
	}

	return refinedAnswer, nil
}
