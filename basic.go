package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	gridSize       = 10
	cellSize       = 40
	imgWidth       = gridSize * cellSize
	titleHeight    = 25 // Space for "Find these words:" title
	wordRowSpacing = 20 // Vertical spacing between word rows
	wordsPerRow    = 4  // Number of words per row
	headerPadding  = 15 // Padding between word list and grid
)

// Direction constants for word placement
const (
	RIGHT      = iota // →
	LEFT              // ←
	DOWN              // ↓
	UP                // ↑
	DOWN_RIGHT        // ↘
	DOWN_LEFT         // ↙
	UP_RIGHT          // ↗
	UP_LEFT           // ↖
)

// WordPlacement stores information about a placed word
type WordPlacement struct {
	Word      string
	Row       int
	Col       int
	Direction int
}

var firstGradeVocab = []string{
	"ache",
	"enormous",
	"equal",
	"exclaim",
	"exhausted",
	"expensive",
	"fancy",
	"fasten",
	"filthy",
	"flat",
	"flee",
	"fog",
	"footprint",
	"forest",
	"freezing",
	"gather",
	"giant",
	"glad",
	"gleaming",
	"glum",
	"grab",
	"grateful",
	"grin",
	"grip",
	"groan",
	"hatch",
	"heap",
	"hide",
	"hobby",
	"honest",
	"howl",
	"illustrator",
	"injury",
	"jealous",
	"knob",
	"lively",
	"loosen",
	"mask",
	"misty",
	"modern",
	"mountain",
	"narrow",
	"obey",
	"pain",
	"passenger",
	"pattern",
	"pest",
	"polish",
	"pretend",
	"promise",
	"rapid",
	"remove",
	"repeat",
	"rescue",
	"restart",
	"return",
	"ripe",
	"rise",
	"roar",
	"rough",
	"rusty",
	"scold",
	"scratch",
	"seed",
	"selfish",
	"serious",
	"shell",
	"shovel",
	"shriek",
	"sibling",
	"silent",
	"simple",
	"slippery",
	"sly",
	"sneaky",
	"sob",
	"spiral",
	"splendid",
	"sprinkle",
	"squirm",
	"startle",
	"steep",
	"stormy",
	"striped",
	"surround",
	"switch",
	"terrified",
	"thick",
	"thunder",
	"ticket",
	"timid",
	"transportation",
	"travel",
	"trust",
	"upset",
	"weed",
	"whimper",
	"whirl",
	"wicked",
	"yank",
}

// Command line flag for custom vocabulary file
var vocabFile = flag.String("vocab", "", "Path to custom vocabulary file (one word per line)")

// loadCustomVocab reads a vocabulary file and returns a slice of words
func loadCustomVocab(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open vocabulary file: %v", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" { // Skip empty lines
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading vocabulary file: %v", err)
	}

	if len(words) < 10 {
		return nil, fmt.Errorf("vocabulary file must contain at least 10 words, found %d", len(words))
	}

	return words, nil
}

// filterWordsByLength removes words that are too long to fit in the grid
func filterWordsByLength(words []string, maxLength int) []string {
	var filteredWords []string
	var removedWords []string

	for _, word := range words {
		if len(word) <= maxLength {
			filteredWords = append(filteredWords, word)
		} else {
			removedWords = append(removedWords, word)
		}
	}

	// Print information about filtered words
	if len(removedWords) > 0 {
		fmt.Printf("Removed %d words that are too long for %dx%d grid:\n", len(removedWords), maxLength, maxLength)
		for _, word := range removedWords {
			fmt.Printf("  - %s (%d characters)\n", word, len(word))
		}
	}

	return filteredWords
}

// calculateHeaderHeight calculates the required header space based on number of words
func calculateHeaderHeight(wordCount int) int {
	wordRows := (wordCount + wordsPerRow - 1) / wordsPerRow // Ceiling division
	return titleHeight + (wordRows * wordRowSpacing) + headerPadding
}

func main() {
	// Parse command line flags
	flag.Parse()

	fmt.Println("Generating word search puzzle...")

	// Seed random number generator
	rand.Seed(int64(os.Getpid())) // Use process ID for randomness

	// Determine which vocabulary to use
	var vocabulary []string
	if *vocabFile != "" {
		fmt.Printf("Loading custom vocabulary from: %s\n", *vocabFile)
		customVocab, err := loadCustomVocab(*vocabFile)
		if err != nil {
			log.Fatalf("Error loading custom vocabulary: %v", err)
		}
		vocabulary = customVocab
		fmt.Printf("Loaded %d words from custom vocabulary file\n", len(vocabulary))
	} else {
		vocabulary = firstGradeVocab
		fmt.Printf("Using default first-grade vocabulary (%d words)\n", len(vocabulary))
	}

	// Filter out words that are too long for the grid
	vocabulary = filterWordsByLength(vocabulary, gridSize)
	fmt.Printf("After filtering: %d words available for puzzle generation\n", len(vocabulary))

	// Check if we have enough words after filtering
	if len(vocabulary) < 10 {
		log.Fatalf("Not enough words available after filtering. Need at least 10 words, but only %d words fit in a %dx%d grid", len(vocabulary), gridSize, gridSize)
	}

	// Words to place in the puzzle - select 10 random words from chosen vocabulary
	words := make([]string, 0, 10)
	temp := make([]string, len(vocabulary))
	copy(temp, vocabulary)

	for i := 0; i < 10 && len(temp) > 0; i++ {
		randomIndex := rand.Intn(len(temp))
		words = append(words, strings.ToUpper(temp[randomIndex]))

		temp[randomIndex] = temp[len(temp)-1]
		temp = temp[:len(temp)-1]
	}

	// Calculate dynamic header height and total image height
	headerHeight := calculateHeaderHeight(len(words))
	imgHeight := gridSize*cellSize + headerHeight

	// Create image
	rect := image.Rect(0, 0, imgWidth, imgHeight)
	img := image.NewRGBA(rect)

	// Fill with white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// Draw the word list at the top
	drawWordList(img, words)

	// Generate and draw the word search grid
	letterGrid, placements := generateWordSearchGrid(words)
	drawLetterGrid(img, letterGrid, headerHeight)

	// Save to file
	outFile, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}

	fmt.Printf("Word search puzzle created successfully: output.png (%dx%d pixels)\n", imgWidth, imgHeight)
	printWordSearchGrid(letterGrid)
	printWordPlacements(placements)
	verifyWordPlacements(letterGrid, placements)
}

// generateWordSearchGrid creates a word search puzzle with the given words
func generateWordSearchGrid(words []string) ([gridSize][gridSize]rune, []WordPlacement) {
	var grid [gridSize][gridSize]rune
	var placements []WordPlacement

	// Initialize grid with empty spaces
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			grid[row][col] = 0 // 0 represents empty space
		}
	}

	// Try to place each word
	for _, word := range words {
		placed := false
		attempts := 0
		maxAttempts := 100

		for !placed && attempts < maxAttempts {
			// Random starting position
			row := rand.Intn(gridSize)
			col := rand.Intn(gridSize)
			direction := rand.Intn(8) // 8 directions

			if canPlaceWord(grid, word, row, col, direction) {
				placeWord(&grid, word, row, col, direction)
				placements = append(placements, WordPlacement{
					Word:      word,
					Row:       row,
					Col:       col,
					Direction: direction,
				})
				placed = true
			}
			attempts++
		}

		if !placed {
			fmt.Printf("Warning: Could not place word '%s' after %d attempts\n", word, maxAttempts)
		}
	}

	// Fill empty spaces with random letters
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if grid[row][col] == 0 {
				grid[row][col] = 'A' + rune(rand.Intn(26))
			}
		}
	}

	return grid, placements
}

// canPlaceWord checks if a word can be placed at the given position and direction
func canPlaceWord(grid [gridSize][gridSize]rune, word string, row, col, direction int) bool {
	wordLen := len(word)

	// Get direction deltas
	dRow, dCol := getDirectionDeltas(direction)

	// Check if word fits within grid bounds
	endRow := row + dRow*(wordLen-1)
	endCol := col + dCol*(wordLen-1)

	if endRow < 0 || endRow >= gridSize || endCol < 0 || endCol >= gridSize {
		return false
	}

	// Check for conflicts with existing letters
	for i, letter := range word {
		currentRow := row + dRow*i
		currentCol := col + dCol*i
		existingLetter := grid[currentRow][currentCol]

		// If there's already a letter, it must match
		if existingLetter != 0 && existingLetter != rune(letter) {
			return false
		}
	}

	return true
}

// placeWord places a word in the grid at the specified position and direction
func placeWord(grid *[gridSize][gridSize]rune, word string, row, col, direction int) {
	dRow, dCol := getDirectionDeltas(direction)

	for i, letter := range word {
		currentRow := row + dRow*i
		currentCol := col + dCol*i
		grid[currentRow][currentCol] = rune(letter)
	}
}

// getDirectionDeltas returns the row and column deltas for a given direction
func getDirectionDeltas(direction int) (int, int) {
	switch direction {
	case RIGHT:
		return 0, 1
	case LEFT:
		return 0, -1
	case DOWN:
		return 1, 0
	case UP:
		return -1, 0
	case DOWN_RIGHT:
		return 1, 1
	case DOWN_LEFT:
		return 1, -1
	case UP_RIGHT:
		return -1, 1
	case UP_LEFT:
		return -1, -1
	default:
		return 0, 1 // Default to RIGHT
	}
}

// drawWordList draws the list of words to find at the top of the image
func drawWordList(img *image.RGBA, words []string) {
	face := basicfont.Face7x13

	// Draw title "Find these words:"
	titleX := 10
	titleY := 20
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}), // Black color
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(titleX), Y: fixed.I(titleY)},
	}
	drawer.DrawString("Find these words:")

	// Draw words in rows (4 words per row)
	wordsPerRow := 4
	wordSpacing := 80 // Horizontal spacing between words
	rowSpacing := 20  // Vertical spacing between rows
	startX := 10
	startY := 45

	for i, word := range words {
		row := i / wordsPerRow
		col := i % wordsPerRow

		x := startX + col*wordSpacing
		y := startY + row*rowSpacing

		drawer.Dot = fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
		drawer.DrawString(word)
	}
}

// drawLetterGrid draws all letters in the grid onto the image
func drawLetterGrid(img *image.RGBA, grid [gridSize][gridSize]rune, headerHeight int) {
	face := basicfont.Face7x13

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			letter := grid[row][col]
			drawLetterInCell(img, face, letter, row, col, headerHeight)
		}
	}
}

// drawLetterInCell draws a single letter centered in the specified grid cell
func drawLetterInCell(img *image.RGBA, face font.Face, letter rune, row, col, headerHeight int) {
	// Calculate cell position (offset by header height)
	cellX := col * cellSize
	cellY := row*cellSize + headerHeight

	// Calculate letter position to center it in the cell
	// basicfont.Face7x13 is 7 pixels wide, 13 pixels tall
	letterWidth := 7
	letterHeight := 13

	// Center the letter in the cell
	x := cellX + (cellSize-letterWidth)/2
	y := cellY + (cellSize-letterHeight)/2 + 10 // +10 for baseline adjustment

	// Create a font drawer
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}), // Black color
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}

	// Draw the letter
	drawer.DrawString(string(letter))
}

// printWordSearchGrid prints the generated grid to console for reference
func printWordSearchGrid(grid [gridSize][gridSize]rune) {
	fmt.Println("\nWord search puzzle grid:")
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			fmt.Printf("%c ", grid[row][col])
		}
		fmt.Println()
	}
}

// printWordPlacements prints information about where each word was placed
func printWordPlacements(placements []WordPlacement) {
	fmt.Println("\nWords placed in the puzzle:")
	directionNames := []string{"RIGHT", "LEFT", "DOWN", "UP", "DOWN_RIGHT", "DOWN_LEFT", "UP_RIGHT", "UP_LEFT"}

	for _, placement := range placements {
		directionName := "UNKNOWN"
		if placement.Direction >= 0 && placement.Direction < len(directionNames) {
			directionName = directionNames[placement.Direction]
		}
		fmt.Printf("- %s: Row %d, Col %d, Direction %s\n",
			placement.Word, placement.Row, placement.Col, directionName)
	}

	fmt.Printf("\nTotal words placed: %d\n", len(placements))
}

// verifyWordPlacements checks if the placed words actually exist in the grid
func verifyWordPlacements(grid [gridSize][gridSize]rune, placements []WordPlacement) {
	fmt.Println("\nVerifying word placements:")

	for _, placement := range placements {
		word := placement.Word
		row := placement.Row
		col := placement.Col
		direction := placement.Direction

		dRow, dCol := getDirectionDeltas(direction)

		// Extract the word from the grid
		var extractedWord string
		for i := 0; i < len(word); i++ {
			currentRow := row + dRow*i
			currentCol := col + dCol*i

			if currentRow >= 0 && currentRow < gridSize && currentCol >= 0 && currentCol < gridSize {
				extractedWord += string(grid[currentRow][currentCol])
			} else {
				extractedWord += "?"
			}
		}

		if extractedWord == word {
			fmt.Printf("✓ %s found correctly at (%d,%d)\n", word, row, col)
		} else {
			fmt.Printf("✗ %s NOT found at (%d,%d) - found '%s' instead\n", word, row, col, extractedWord)
		}
	}
}
