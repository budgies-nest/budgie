# Advanced Markdown Chunking in Go with `ParseMarkdownWithLineage`

This Go function takes markdown text and breaks it down into structured chunks based on headers, while keeping track of the hierarchical relationship between sections.

Here's what it does step by step:

**Main Purpose**: Parse markdown content and create a hierarchy of sections with their parent-child relationships.

**How it works**:

1. **Finds Headers**: Uses a regex pattern to identify markdown headers (lines starting with `#`, `##`, `###`, etc.)

2. **Extracts Content**: For each header found, it collects all the text that follows until the next header appears

3. **Tracks Hierarchy**: Uses a "stack" to keep track of parent headers:
   - When it finds a header at the same or higher level, it removes deeper headers from the stack
   - The current parent is whatever header is on top of the stack

4. **Builds Lineage**: Creates a breadcrumb trail showing the path from the top-level header down to the current one (like "Introduction > Getting Started > Installation")

5. **Creates Chunks**: Each chunk contains:
   - The header text and its level (`#` = level 1, `##` = level 2, etc.)
   - The content under that header
   - Information about its parent header
   - The full lineage path

**Example**: If you have markdown like:
```markdown
# Chapter 1
Some intro text
## Section A
Details about A
### Subsection A.1
More details
```

It would create chunks where "Subsection A.1" knows its parent is "Section A", and its lineage is "Chapter 1 > Section A > Subsection A.1".

This is commonly used in RAG (Retrieval-Augmented Generation) systems to maintain context when searching through documentation.