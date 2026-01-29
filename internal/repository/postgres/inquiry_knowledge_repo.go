package postgres

import (
	"context"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/pgvector/pgvector-go"
	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/internal/domain"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/internal/repository/postgres/dao/ent"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
	"github.com/wonjinsin/simple-chatbot/pkg/utils"
)

type inquiryKnowledgeRepo struct {
	client *ent.Client
}

// NewInquiryKnowledgeRepository creates a new EntGo-based inquiry knowledge repository
func NewInquiryKnowledgeRepository(client *ent.Client) repository.InquiryKnowledgeRepository {
	return &inquiryKnowledgeRepo{client: client}
}

// BatchSaveInquiryKnowledge saves multiple inquiry knowledge entries to database using EntGo
func (r *inquiryKnowledgeRepo) BatchSaveInquiryKnowledge(
	ctx context.Context,
	items domain.InquiryKnowledges,
) error {
	if len(items) == 0 {
		return nil
	}

	// Use transaction for bulk operations
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	// Insert each item (EntGo doesn't have native bulk upsert, so we do it individually)
	for _, item := range items {
		instruction, response, category, intent, flags, embedding, createdAt, updatedAt :=
			toEntInquiryKnowledgeData(item)

		// Try to create, if conflict (already exists), update instead
		_, err := tx.InquiryKnowledge.Create().
			SetInstruction(instruction).
			SetNillableInstructionEmbedding(&embedding).
			SetResponse(response).
			SetNillableCategory(&category).
			SetNillableIntent(&intent).
			SetNillableFlags(&flags).
			SetCreatedAt(createdAt).
			SetUpdatedAt(updatedAt).
			Save(ctx)

		if err != nil {
			// If constraint error (duplicate instruction), update instead
			if ent.IsConstraintError(err) {
				// Find existing record and update
				_, updateErr := tx.InquiryKnowledge.Update().
					Where(func(s *entsql.Selector) {
						s.Where(entsql.EQ("instruction", instruction))
					}).
					SetNillableInstructionEmbedding(&embedding).
					SetResponse(response).
					SetNillableCategory(&category).
					SetNillableIntent(&intent).
					SetNillableFlags(&flags).
					SetUpdatedAt(updatedAt).
					Save(ctx)

				if updateErr != nil {
					if rollbackErr := tx.Rollback(); rollbackErr != nil {
						return errors.Wrap(rollbackErr, "failed to rollback after update error")
					}
					return errors.Wrap(updateErr, "failed to update inquiry knowledge")
				}
			} else {
				// Other errors, rollback
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return errors.Wrap(rollbackErr, "failed to rollback after insert error")
				}
				return errors.Wrap(err, "failed to insert inquiry knowledge")
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// FindSimilars finds inquiry knowledge entries similar to the given embedding vector with similarity
// scores
func (r *inquiryKnowledgeRepo) FindSimilars(
	ctx context.Context,
	embedding domain.Embedding,
	limit int,
) (domain.InquirySimilarityResults, error) {
	if limit <= 0 {
		return nil, errors.New(
			constants.InvalidParameter,
			"limit must be greater than 0",
			nil,
		)
	}

	// Convert []float64 to pgvector.Vector
	vec := make([]float32, len(embedding))
	for i, v := range embedding {
		vec[i] = float32(v)
	}
	queryVector := pgvector.NewVector(vec)

	// Query using pgvector's cosine distance operator (<=>)
	// First, get the most similar entries ordered by distance
	entResults, err := r.client.InquiryKnowledge.Query().
		Where(func(s *entsql.Selector) {
			s.Where(entsql.NotNull("instruction_embedding"))
		}).
		Order(func(s *entsql.Selector) {
			// Order by cosine distance (smaller distance = more similar)
			s.OrderExpr(entsql.Expr(fmt.Sprintf(
				"instruction_embedding <=> '%s'",
				queryVector.String(),
			)))
		}).
		Limit(limit).
		All(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to query similar inquiry knowledge")
	}

	if len(entResults) == 0 {
		return nil, errors.New(
			constants.NotFound,
			"no similar inquiry knowledge found",
			nil,
		)
	}

	// Calculate similarity scores for each result
	domainResults := make(domain.InquirySimilarityResults, len(entResults))

	for i, entIK := range entResults {
		// Convert pgvector.Vector to []float64
		vec2Slice := entIK.InstructionEmbedding.Slice()
		vec2 := make([]float64, len(vec2Slice))
		for j, v := range vec2Slice {
			vec2[j] = float64(v)
		}

		// Calculate normalized similarity score (0.0 to 1.0)
		similarity := utils.CalculateVectorSimilarity(embedding, vec2)

		domainResults[i] = &domain.InquirySimilarityResult{
			Knowledge:       toDomainInquiryKnowledge(entIK),
			SimilarityScore: similarity,
		}
	}

	return domainResults, nil
}
