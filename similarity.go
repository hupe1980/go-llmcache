package llmcache

import "math"

// Similarity calculates the cosine similarity between two embeddings.
func Similarity(a, b []float64) float64 {
	dotProduct := float64(0)
	magnitudeA := float64(0)
	magnitudeB := float64(0)

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)

	if magnitudeA > 0 && magnitudeB > 0 {
		return dotProduct / (magnitudeA * magnitudeB)
	}

	return 0
}
