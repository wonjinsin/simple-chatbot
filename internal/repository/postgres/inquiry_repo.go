package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/wonjinsin/simple-chatbot/internal/domain"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

type inquiryKnowledgeRepo struct {
	db *sql.DB
}

// NewInquiryKnowledgeRepository creates a new PostgreSQL-based inquiry knowledge repository
func NewInquiryKnowledgeRepository(db *sql.DB) repository.InquiryKnowledgeRepository {
	return &inquiryKnowledgeRepo{db: db}
}

// Close closes the database connection
func (r *inquiryKnowledgeRepo) Close() error {
	if err := r.db.Close(); err != nil {
		return errors.Wrap(err, "failed to close database connection")
	}
	return nil
}

// BatchSaveInquiryKnowledge saves multiple inquiry knowledge entries to database
func (r *inquiryKnowledgeRepo) BatchSaveInquiryKnowledge(
	ctx context.Context,
	items domain.InquiryKnowledges,
) error {
	if len(items) == 0 {
		return nil
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		_ = tx.Rollback() // Rollback is safe to call even after Commit
	}()

	// Prepare INSERT statement
	query := `
		INSERT INTO inquiry_knowledge_base (
			instruction,
			instruction_embedding,
			response,
			category,
			intent,
			flags,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (instruction) 
		DO UPDATE SET
			instruction_embedding = EXCLUDED.instruction_embedding,
			response = EXCLUDED.response,
			category = EXCLUDED.category,
			intent = EXCLUDED.intent,
			flags = EXCLUDED.flags,
			updated_at = EXCLUDED.updated_at
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement")
	}
	defer stmt.Close()

	// Insert each item
	for _, item := range items {
		// Convert embedding slice to pgvector format string: '[1.0, 2.0, 3.0]'
		embeddingStr := formatVectorString(item.InstructionEmbedding)

		_, err := stmt.ExecContext(
			ctx,
			item.Instruction,
			embeddingStr,
			item.Response,
			nullStringIfEmpty(item.Category),
			nullStringIfEmpty(item.Intent),
			nullStringIfEmpty(item.Flags),
			item.CreatedAt,
			item.UpdatedAt,
		)
		if err != nil {
			return errors.Wrap(err, "failed to insert inquiry knowledge")
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// formatVectorString converts float64 slice to pgvector format: '[1.0, 2.0, 3.0]'
func formatVectorString(vector []float64) string {
	if len(vector) == 0 {
		return "[]"
	}

	parts := make([]string, len(vector))
	for i, v := range vector {
		parts[i] = fmt.Sprintf("%f", v)
	}

	return "[" + strings.Join(parts, ",") + "]"
}

// nullStringIfEmpty returns sql.NullString for optional string fields
func nullStringIfEmpty(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
