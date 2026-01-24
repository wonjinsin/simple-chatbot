package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/internal/domain"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
	"github.com/wonjinsin/simple-chatbot/pkg/file"
	"github.com/wonjinsin/simple-chatbot/pkg/logger"
)

const (
	batchSize = 50 // Batch size for embedding API calls
)

type InquiryServiceImpl struct {
	embeddingRepo repository.EmbeddingRepository
	knowledgeRepo repository.InquiryKnowledgeRepository
}

func NewInquiryServiceImpl(
	embeddingRepo repository.EmbeddingRepository,
	knowledgeRepo repository.InquiryKnowledgeRepository,
) *InquiryServiceImpl {
	return &InquiryServiceImpl{
		embeddingRepo: embeddingRepo,
		knowledgeRepo: knowledgeRepo,
	}
}

// EmbedInquiryOrigins reads CSV data, generates embeddings, and saves to database
func (s *InquiryServiceImpl) EmbedInquiryOrigins(ctx context.Context) error {
	// Step 1: Read CSV file
	csvRows, err := file.ReadCSVToMapArray("mock_data/data_set.csv")
	if err != nil {
		return errors.Wrap(err, "failed to read inquiry origins", constants.InternalError)
	}

	logger.LogInfo(ctx, fmt.Sprintf("read %d rows from CSV", len(csvRows)))

	// Step 2: Convert CSV rows to domain objects (without embeddings)
	knowledgeItems := make(domain.InquiryKnowledges, 0, len(csvRows))
	for _, row := range csvRows {
		item, err := domain.NewInquiryKnowledgeFromCSV(row)
		if err != nil {
			// Skip invalid rows
			fmt.Printf("Skipping invalid row: %v\n", err)
			continue
		}
		knowledgeItems = append(knowledgeItems, item)
	}

	fmt.Printf("Converted %d valid items from CSV\n", len(knowledgeItems))

	// Step 3: Process in batches
	totalBatches := (len(knowledgeItems) + batchSize - 1) / batchSize
	for i := 0; i < len(knowledgeItems); i += batchSize {
		end := min(i+batchSize, len(knowledgeItems))

		batch := knowledgeItems[i:end]
		batchNum := (i / batchSize) + 1

		fmt.Printf("Processing batch %d/%d (%d items)\n", batchNum, totalBatches, len(batch))

		if err := s.processBatch(ctx, batch); err != nil {
			return errors.Wrap(
				err,
				fmt.Sprintf("failed to process batch %d", batchNum),
				constants.InternalError,
			)
		}

		fmt.Printf("Batch %d/%d completed successfully\n", batchNum, totalBatches)
	}

	fmt.Printf("All %d items embedded and saved successfully\n", len(knowledgeItems))
	return nil
}

// processBatch embeds and saves a batch of inquiry knowledge items
func (s *InquiryServiceImpl) processBatch(
	ctx context.Context,
	batch domain.InquiryKnowledges,
) error {
	// Extract instructions for embedding
	instructions := make([]string, len(batch))
	for i, item := range batch {
		instructions[i] = item.Instruction
	}

	// Generate embeddings
	embeddings, err := s.embeddingRepo.EmbedStrings(ctx, instructions)
	if err != nil {
		return errors.Wrap(err, "failed to generate embeddings")
	}

	if len(embeddings) != len(batch) {
		return errors.New(
			constants.InternalError,
			fmt.Sprintf(
				"embedding count mismatch: expected %d, got %d",
				len(batch),
				len(embeddings),
			),
			nil,
		)
	}

	// Attach embeddings to items
	now := time.Now()
	for i, item := range batch {
		item.InstructionEmbedding = embeddings[i]
		item.CreatedAt = now
		item.UpdatedAt = now
	}

	// Save to database
	if err := s.knowledgeRepo.BatchSaveInquiryKnowledge(ctx, batch); err != nil {
		return errors.Wrap(err, "failed to save inquiry knowledge")
	}

	return nil
}
