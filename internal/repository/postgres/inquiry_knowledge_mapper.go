package postgres

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"github.com/wonjinsin/simple-chatbot/internal/domain"
	"github.com/wonjinsin/simple-chatbot/internal/repository/postgres/dao/ent"
)

// toEntInquiryKnowledge converts domain.InquiryKnowledge to ent.InquiryKnowledge
func toEntInquiryKnowledge(ik *domain.InquiryKnowledge) *ent.InquiryKnowledge {
	entIK := &ent.InquiryKnowledge{
		Instruction: ik.Instruction,
		Response:    ik.Response,
		Category:    ik.Category,
		Intent:      ik.Intent,
		Flags:       ik.Flags,
		CreatedAt:   ik.CreatedAt,
		UpdatedAt:   ik.UpdatedAt,
	}

	// Convert []float64 to pgvector.Vector
	if len(ik.InstructionEmbedding) > 0 {
		vec := make([]float32, len(ik.InstructionEmbedding))
		for i, v := range ik.InstructionEmbedding {
			vec[i] = float32(v)
		}
		entIK.InstructionEmbedding = pgvector.NewVector(vec)
	}

	if ik.ID != 0 {
		entIK.ID = ik.ID
	}

	return entIK
}

// toDomainInquiryKnowledge converts ent.InquiryKnowledge to domain.InquiryKnowledge
func toDomainInquiryKnowledge(entIK *ent.InquiryKnowledge) *domain.InquiryKnowledge {
	ik := &domain.InquiryKnowledge{
		ID:          entIK.ID,
		Instruction: entIK.Instruction,
		Response:    entIK.Response,
		Category:    entIK.Category,
		Intent:      entIK.Intent,
		Flags:       entIK.Flags,
		CreatedAt:   entIK.CreatedAt,
		UpdatedAt:   entIK.UpdatedAt,
	}

	// Convert pgvector.Vector to []float64
	if entIK.InstructionEmbedding.Slice() != nil {
		vec := entIK.InstructionEmbedding.Slice()
		ik.InstructionEmbedding = make([]float64, len(vec))
		for i, v := range vec {
			ik.InstructionEmbedding[i] = float64(v)
		}
	}

	return ik
}

// toEntInquiryKnowledgeData extracts data for bulk create
func toEntInquiryKnowledgeData(ik *domain.InquiryKnowledge) (
	instruction, response, category, intent, flags string,
	embedding pgvector.Vector,
	createdAt, updatedAt time.Time,
) {
	instruction = ik.Instruction
	response = ik.Response
	category = ik.Category
	intent = ik.Intent
	flags = ik.Flags
	createdAt = ik.CreatedAt
	updatedAt = ik.UpdatedAt

	// Convert []float64 to pgvector.Vector
	if len(ik.InstructionEmbedding) > 0 {
		vec := make([]float32, len(ik.InstructionEmbedding))
		for i, v := range ik.InstructionEmbedding {
			vec[i] = float32(v)
		}
		embedding = pgvector.NewVector(vec)
	}

	return
}
