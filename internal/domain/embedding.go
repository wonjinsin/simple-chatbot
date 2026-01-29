package domain

type Embedding []float64

func NewEmbedding(values []float64) Embedding {
	var embedding Embedding
	for _, value := range values {
		embedding = append(embedding, value)
	}
	return embedding
}

func (e *Embedding) IsEmpty() bool {
	return e == nil || len(*e) == 0
}

type Embeddings []Embedding

func NewEmbeddings(values [][]float64) Embeddings {
	embeddings := make(Embeddings, len(values))
	for i, value := range values {
		embeddings[i] = NewEmbedding(value)
	}
	return embeddings
}
