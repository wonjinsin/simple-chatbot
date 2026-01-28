package utils

import "math"

// CalculateVectorSimilarity calculates the normalized similarity score between two vectors.
// Returns a value between 0 and 1, where:
//   - 1.0 means the vectors are identical (completely similar)
//   - 0.5 means the vectors are orthogonal (no similarity)
//   - 0.0 means the vectors are opposite (completely dissimilar)
func CalculateVectorSimilarity(vec1, vec2 []float64) float64 {
	if len(vec1) == 0 || len(vec2) == 0 {
		return 0.0
	}

	if len(vec1) != len(vec2) {
		return 0.0
	}

	var dotProduct, norm1, norm2 float64

	// Calculate dot product and norms
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		norm1 += vec1[i] * vec1[i]
		norm2 += vec2[i] * vec2[i]
	}

	// Avoid division by zero
	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	// Cosine similarity = dot product / (||vec1|| * ||vec2||)
	cosineSimilarity := dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))

	// Clamp to [-1, 1] range (in case of floating point errors)
	if cosineSimilarity > 1.0 {
		cosineSimilarity = 1.0
	} else if cosineSimilarity < -1.0 {
		cosineSimilarity = -1.0
	}

	// Normalize to [0, 1] range where 1 = identical, 0 = opposite
	normalizedSimilarity := (cosineSimilarity + 1.0) / 2.0

	return normalizedSimilarity
}
