# Chunking
> 🚧 work in progress

## Chunking text

### Chunking text with overlap

```go
content := "# Star Trek: The Original Series"

chunks := rag.ChunkText(content, 10, 5)

fmt.Println("Chunks:")
for i, chunk := range chunks {
    fmt.Printf("Chunk %d: %s\n", i+1, chunk)
}
fmt.Println("Number of Chunks:", len(chunks))
```

**Output**:
```raw
Chunks:
Chunk 1: # Star Tre
Chunk 2: r Trek: Th
Chunk 3: k: The Ori
Chunk 4: e Original
Chunk 5: ginal Seri
Chunk 6:  Series
Chunk 7: es
Number of Chunks: 7
```

### Chunking text with delimiter

```go
content := "# Star Trek: The Original Series"

chunks := rag.SplitTextWithDelimiter(content, " ")

fmt.Println("Chunks:")
for i, chunk := range chunks {
    fmt.Printf("Chunk %d: %s\n", i+1, chunk)
}
fmt.Println("Number of Chunks:", len(chunks))
```

**Output**:
```raw
Chunks:
Chunk 1: #
Chunk 2: Star
Chunk 3: Trek:
Chunk 4: The
Chunk 5: Original
Chunk 6: Series
Number of Chunks: 6
```

## Chunking Markdown

### Chunking Markdown by Headers/Sections

```go
content := `
# Star Trek: The Original Series

"Star Trek: The Original Series" is an American science fiction television series 
created by Gene Roddenberry that follows the adventures of the starship USS Enterprise (NCC-1701) 
and its crew as they explore the galaxy.

## Season 1

Season 1 of "Star Trek: The Original Series" consists of 29 episodes, 
introducing viewers to the crew of the USS Enterprise and setting the stage 
for the series' exploration of complex themes such as morality, ethics, and the human condition. 

### Episode

> ...
`

chunks := rag.SplitMarkdownBySections(content)

fmt.Println("Chunks:")
for i, chunk := range chunks {
    fmt.Printf("Chunk %d: %s\n", i+1, chunk)
}
fmt.Println("Number of Chunks:", len(chunks))
```

**Output**:
```raw
Chunks:
Chunk 1: # Star Trek: The Original Series

        "Star Trek: The Original Series" is an American science fiction television series 
        created by Gene Roddenberry that follows the adventures of the starship USS Enterprise (NCC-1701) 
        and its crew as they explore the galaxy.
Chunk 2: ## Season 1

        Season 1 of "Star Trek: The Original Series" consists of 29 episodes, 
        introducing viewers to the crew of the USS Enterprise and setting the stage 
        for the series' exploration of complex themes such as morality, ethics, and the human condition.
Chunk 3: ### Episode

        > ...
Number of Chunks: 3
```

### Chunking Markdown with Hierarchy

```go
content := `
# Star Trek: The Original Series

"Star Trek: The Original Series" is an American science fiction television series 
created by Gene Roddenberry that follows the adventures of the starship USS Enterprise (NCC-1701) 
and its crew as they explore the galaxy.

## Season 1

Season 1 of "Star Trek: The Original Series" consists of 29 episodes, 
introducing viewers to the crew of the USS Enterprise and setting the stage 
for the series' exploration of complex themes such as morality, ethics, and the human condition. 

### Episode
> ...
`

// Trim leading whitespace (spaces and tabs) from each line
lines := strings.Split(content, "\n")
for i, line := range lines {
    lines[i] = strings.TrimLeft(line, " \t")
}
result := strings.Join(lines, "\n")


chunks := rag.ChunkWithMarkdownHierarchy(result)

fmt.Println("Chunks:")
for _, chunk := range chunks {
    fmt.Println(chunk)
    fmt.Println(strings.Repeat("-", 40))
}
fmt.Println("Number of Chunks:", len(chunks))
```

**Output**:
```raw
Chunks:
TITLE: # Star Trek: The Original Series
HIERARCHY: Star Trek: The Original Series
CONTENT: "Star Trek: The Original Series" is an American science fiction television series 
created by Gene Roddenberry that follows the adventures of the starship USS Enterprise (NCC-1701) 
and its crew as they explore the galaxy.
----------------------------------------
TITLE: ## Season 1
HIERARCHY: Star Trek: The Original Series > Season 1
CONTENT: Season 1 of "Star Trek: The Original Series" consists of 29 episodes, 
introducing viewers to the crew of the USS Enterprise and setting the stage 
for the series' exploration of complex themes such as morality, ethics, and the human condition.
----------------------------------------
TITLE: ### Episode
HIERARCHY: Star Trek: The Original Series > Season 1 > Episode
CONTENT: > ...
----------------------------------------
Number of Chunks: 3
```