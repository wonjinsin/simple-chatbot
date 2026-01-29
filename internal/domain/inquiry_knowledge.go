package domain

import (
	"strings"
	"time"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
	"github.com/wonjinsin/simple-chatbot/pkg/utils"
)

// InquiryKnowledge represents a knowledge base entry for customer inquiries
type InquiryKnowledge struct {
	ID                   int
	Instruction          string
	InstructionEmbedding Embedding
	Response             string
	Category             string
	Intent               string
	Flags                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// NewInquiryKnowledge creates a new InquiryKnowledge instance with validation
func NewInquiryKnowledge(
	instruction, response, category, intent, flags string,
	embedding Embedding,
	now time.Time,
) (*InquiryKnowledge, error) {
	// Normalize and validate instruction
	instruction = strings.TrimSpace(instruction)
	if utils.IsEmptyOrWhitespace(instruction) {
		return nil, errors.New(constants.InvalidParameter, "instruction cannot be empty", nil)
	}

	// Normalize and validate response
	response = strings.TrimSpace(response)
	if utils.IsEmptyOrWhitespace(response) {
		return nil, errors.New(constants.InvalidParameter, "response cannot be empty", nil)
	}

	// Normalize optional fields
	category = strings.TrimSpace(category)
	intent = strings.TrimSpace(intent)
	flags = strings.TrimSpace(flags)

	return &InquiryKnowledge{
		Instruction:          instruction,
		InstructionEmbedding: embedding,
		Response:             response,
		Category:             category,
		Intent:               intent,
		Flags:                flags,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

// NewInquiryKnowledgeFromCSV creates InquiryKnowledge from CSV row data
func NewInquiryKnowledgeFromCSVs(csvRows []map[string]string) (InquiryKnowledges, error) {
	knowledgeItems := make(InquiryKnowledges, 0, len(csvRows))
	for _, csvRow := range csvRows {
		instruction := csvRow["instruction"]
		response := csvRow["response"]
		category := csvRow["category"]
		intent := csvRow["intent"]
		flags := csvRow["flags"]

		// Validate required fields
		instruction = strings.TrimSpace(instruction)
		if utils.IsEmptyOrWhitespace(instruction) {
			return nil, errors.New(constants.InvalidParameter, "instruction cannot be empty", nil)
		}

		response = strings.TrimSpace(response)
		if utils.IsEmptyOrWhitespace(response) {
			return nil, errors.New(constants.InvalidParameter, "response cannot be empty", nil)
		}

		knowledgeItems = append(knowledgeItems, &InquiryKnowledge{
			Instruction: instruction,
			Response:    response,
			Category:    strings.TrimSpace(category),
			Intent:      strings.TrimSpace(intent),
			Flags:       strings.TrimSpace(flags),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return knowledgeItems, nil
}

// InquiryKnowledges is a collection of InquiryKnowledge
type InquiryKnowledges []*InquiryKnowledge

// Instructions returns the instructions of the inquiry knowledge
func (is InquiryKnowledges) Instructions() []string {
	instructions := make([]string, len(is))
	for i, item := range is {
		instructions[i] = item.Instruction
	}
	return instructions
}

func (is InquiryKnowledges) SetEmbeddings(embeddings Embeddings) {
	for i, item := range is {
		item.InstructionEmbedding = embeddings[i]
	}
}

// InquirySimilarityResult represents an inquiry knowledge entry with its similarity score
type InquirySimilarityResult struct {
	Knowledge       *InquiryKnowledge
	SimilarityScore float64 // Cosine similarity score (0.0 to 1.0, higher is more similar)
}

// InquirySimilarityResults is a collection of InquirySimilarityResult
type InquirySimilarityResults []*InquirySimilarityResult
