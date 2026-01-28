package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/pgvector/pgvector-go"
)

// InquiryKnowledge holds the schema definition for the InquiryKnowledge entity.
type InquiryKnowledge struct {
	ent.Schema
}

// Annotations of the InquiryKnowledge.
func (InquiryKnowledge) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "inquiry_knowledges"},
	}
}

// Fields of the InquiryKnowledge.
func (InquiryKnowledge) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.Text("instruction").
			NotEmpty(),
		field.Other("instruction_embedding", pgvector.Vector{}).
			SchemaType(map[string]string{
				dialect.Postgres: "vector(1536)",
			}).
			Optional(),
		field.Text("response").
			NotEmpty(),
		field.String("category").
			Optional(),
		field.String("intent").
			Optional(),
		field.String("flags").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes of the InquiryKnowledge.
func (InquiryKnowledge) Indexes() []ent.Index {
	return []ent.Index{
		// HNSW index for vector similarity search
		index.Fields("instruction_embedding").
			Annotations(
				entsql.IndexType("hnsw"),
				entsql.OpClass("vector_cosine_ops"),
			),
	}
}
