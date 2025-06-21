package rag

type VectorRecord struct {
	Id               string    `json:"id"`
	Prompt           string    `json:"prompt"`
	Embedding        []float64 `json:"embedding"`
	CosineSimilarity float64
}


type VectorStore interface {
	GetAll() ([]VectorRecord, error)
	Save(vectorRecord VectorRecord) (VectorRecord, error)
	SearchSimilarities(embeddingFromQuestion VectorRecord, limit float64) ([]VectorRecord, error)
	SearchTopNSimilarities(embeddingFromQuestion VectorRecord, limit float64, max int) ([]VectorRecord, error)
}