package rag

import (
	"testing"

	"github.com/budgies-nest/budgie/helpers"
)

// go test -v -run TestChunkWithMarkdownHierarchy

func TestChunkWithMarkdownHierarchy(t *testing.T) {
	markdownDocument, err := helpers.ReadTextFile("star-trek.md")
	if err != nil {
		t.Fatalf("Failed to read markdown file: %v", err)
	}
	chunks := ChunkWithMarkdownHierarchy(markdownDocument)
	if len(chunks) == 0 {
		t.Fatalf("Expected non-empty chunks, got %d", len(chunks))
	}
	if len(chunks) != 64 {
		t.Fatalf("Expected 64 chunks, got %d", len(chunks))
	}
}
