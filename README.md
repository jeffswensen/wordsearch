# Word Search Puzzle Generator

A Go-based command-line tool that generates customizable word search puzzles as PNG images. Perfect for educators, parents, or anyone who enjoys word puzzles!

## Features

- **10x10 Grid Generation**: Creates standard-sized word search puzzles
- **8-Directional Word Placement**: Words can be placed horizontally, vertically, and diagonally in all directions
- **First-Grade Vocabulary**: Built-in vocabulary list with 90+ age-appropriate words
- **Random Selection**: Automatically selects 10 random words for each puzzle
- **Visual Output**: Generates high-quality PNG images with clear, readable text
- **Integrated Word List**: Each puzzle includes a "Find these words" section at the top
- **Smart Placement Algorithm**: Intelligently places words to avoid conflicts
- **Random Fill**: Empty cells are filled with random letters to increase difficulty

## Requirements

- Go 1.24.4 or higher
- No additional system dependencies required

## Installation

1. Clone or download this repository
2. Navigate to the project directory
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Basic Usage

Generate a new word search puzzle:

```bash
go run .
```

This will create a new `output.png` file in the current directory containing your word search puzzle.

### Output

The generated puzzle includes:
- A header with "Find these words:" followed by the list of hidden words
- A 10x10 letter grid containing the hidden words
- Console output showing:
  - The generated grid
  - Word placement details (position and direction)
  - Verification of successful word placement

### Example Output

```
Generating word search puzzle...
Word search puzzle created successfully: output.png (400x285 pixels)

Word search puzzle grid:
A B C D E F G H I J
K L M N O P Q R S T
...

Words placed in the puzzle:
- ENORMOUS: Row 2, Col 1, Direction DOWN_RIGHT
- GRATEFUL: Row 0, Col 3, Direction RIGHT
...

Total words placed: 10
```

## Technical Details

### Grid Specifications
- **Grid Size**: 10x10 cells
- **Cell Size**: 40x40 pixels
- **Image Dimensions**: 400px width, variable height (depends on word count)
- **Font**: Basic 7x13 monospace font for clear readability

### Word Placement Algorithm
- Attempts up to 100 placements per word for optimal distribution
- Supports 8 directions: RIGHT, LEFT, DOWN, UP, DOWN_RIGHT, DOWN_LEFT, UP_RIGHT, UP_LEFT
- Prevents word conflicts by checking existing letter placements
- Falls back gracefully if a word cannot be placed

### Vocabulary Source
The default vocabulary consists of 90+ first-grade level words including:
- Basic emotions (glad, upset, grateful)
- Action words (grab, flee, rescue)
- Descriptive words (enormous, tiny, rough)
- Nature words (forest, mountain, thunder)

## Customization

### Modifying the Vocabulary

To use your own word list, edit the `firstGradeVocab` slice in `basic.go`:

```go
var firstGradeVocab = []string{
    "your",
    "custom",
    "words",
    "here",
}
```

### Adjusting Grid Size

Modify the `gridSize` constant to change the puzzle dimensions:

```go
const gridSize = 15  // Creates a 15x15 grid
```

### Changing Word Count

Modify the word selection loop in the `main()` function:

```go
for i := 0; i < 15 && len(temp) > 0; i++ {  // Select 15 words instead of 10
```

## File Structure

```
wordsearch/
├── basic.go          # Main application code
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── README.md        # This file
├── .gitignore       # Git ignore rules
├── output.png       # Generated puzzle (created after running)
└── example.png      # Sample output
```

## Code Structure

- **Main Generation**: `main()` function orchestrates the entire process
- **Grid Creation**: `generateWordSearchGrid()` handles word placement logic
- **Word Placement**: `canPlaceWord()` and `placeWord()` manage individual word positioning
- **Image Rendering**: `drawLetterGrid()` and `drawWordList()` create the visual output
- **Verification**: Built-in verification ensures all words are correctly placed

## Dependencies

- `golang.org/x/image`: Extended image processing capabilities for font rendering

## Contributing

Feel free to submit issues and enhancement requests! Some ideas for contributions:

- **Difficulty Levels**: Add easy/medium/hard modes with different grid sizes
- **Custom Themes**: Support for different visual themes and colors
- **Multiple Formats**: Export to PDF, SVG, or other formats
- **Interactive Mode**: Command-line interface for custom word input
- **Batch Generation**: Generate multiple puzzles at once

## License

This project is open source. Feel free to use, modify, and distribute as needed.

---

**Happy Puzzle Solving!** 🧩
