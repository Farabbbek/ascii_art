package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const defaultWidth = 200 // Increased default width to handle longer strings

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Parse alignment flag
	align := "left" // default alignment
	text := ""
	bannerName := "standard" // default banner

	if strings.HasPrefix(os.Args[1], "--align=") {
		align = strings.TrimPrefix(os.Args[1], "--align=")
		if !isValidAlignment(align) {
			printUsage()
			return
		}
		if len(os.Args) < 3 {
			printUsage()
			return
		}
		text = os.Args[2]
		if len(os.Args) > 3 {
			bannerName = os.Args[3]
		}
	} else {
		text = os.Args[1]
		if len(os.Args) > 2 {
			bannerName = os.Args[2]
		}
	}

	// Load banner
	banner, err := loadBanner(bannerName)
	if err != nil {
		fmt.Println("Error loading banner:", err)
		return
	}

	// Generate and print ASCII art
	asciiArt := generateAsciiArt(text, banner, align, defaultWidth)
	printAligned(asciiArt, align)
}

func isValidAlignment(align string) bool {
	return align == "left" || align == "right" || align == "center" || align == "justify"
}

func loadBanner(name string) (map[rune][]string, error) {
	filename := fmt.Sprintf("banner/%s.txt", name)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	banner := make(map[rune][]string)
	scanner := bufio.NewScanner(file)
	currentChar := ' '
	var currentLines []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			if len(currentLines) > 0 {
				banner[currentChar] = currentLines
				currentLines = []string{}
				currentChar++
			}
		} else {
			currentLines = append(currentLines, line)
		}
	}

	if len(currentLines) > 0 {
		banner[currentChar] = currentLines
	}

	return banner, nil
}

func generateAsciiArt(text string, banner map[rune][]string, align string, width int) []string {
	if text == "" {
		return []string{}
	}

	height := len(banner[' ']) // Assuming all characters have same height
	result := make([]string, height)

	// If not justifying, generate as before
	if align != "justify" {
		for _, char := range text {
			charLines, exists := banner[char]
			if !exists {
				charLines = banner[' '] // Use space for unknown characters
			}

			for i := 0; i < height; i++ {
				result[i] += charLines[i]
			}
		}
		return result
	}

	// For justify: calculate spacing between characters
	// First, collect all character blocks
	var charBlocks [][]string
	for _, char := range text {
		charLines, exists := banner[char]
		if !exists {
			charLines = banner[' '] // Use space for unknown characters
		}
		charBlocks = append(charBlocks, charLines)
	}

	// Calculate total width of all blocks without extra spacing
	blockWidth := 0
	if len(charBlocks) > 0 {
		blockWidth = len(charBlocks[0][0]) // Width of one character block
	}
	totalBlockWidth := blockWidth * len(charBlocks)

	// If the total width is already >= target width, no need to justify
	if totalBlockWidth >= width || len(charBlocks) <= 1 {
		for _, charLines := range charBlocks {
			for i := 0; i < height; i++ {
				result[i] += charLines[i]
			}
		}
		return result
	}

	// Calculate spaces to distribute
	extraSpaces := width - totalBlockWidth
	gaps := len(charBlocks) - 1
	if gaps <= 0 {
		for _, charLines := range charBlocks {
			for i := 0; i < height; i++ {
				result[i] += charLines[i]
			}
		}
		return result
	}

	spacesPerGap := extraSpaces / gaps
	extraSpacesForLastGaps := extraSpaces % gaps

	// Build the justified ASCII art
	for charIdx, charLines := range charBlocks {
		for i := 0; i < height; i++ {
			result[i] += charLines[i]
		}

		// Add spaces after this character, unless it's the last one
		if charIdx < len(charBlocks)-1 {
			spaceCount := spacesPerGap
			if charIdx < extraSpacesForLastGaps {
				spaceCount++
			}
			spacePadding := strings.Repeat(" ", spaceCount)
			for i := 0; i < height; i++ {
				result[i] += spacePadding
			}
		}
	}

	return result
}

func printAligned(lines []string, align string) {
	if len(lines) == 0 {
		return
	}

	// Get maximum line length
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Use the larger of defaultWidth and maxLen
	width := defaultWidth
	if maxLen > width {
		width = maxLen
	}

	for _, line := range lines {
		switch align {
		case "right":
			if len(line) < width {
				padding := strings.Repeat(" ", width-len(line))
				fmt.Println(padding + line)
			} else {
				fmt.Println(line)
			}
		case "center":
			if len(line) < width {
				padding := strings.Repeat(" ", (width-len(line))/2)
				fmt.Println(padding + line)
			} else {
				fmt.Println(line)
			}
		case "justify":
			// Justification already handled in generateAsciiArt
			fmt.Println(line)
		default: // left alignment
			fmt.Println(line)
		}
	}
}

func printUsage() {
	fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
	fmt.Println("\nExample: go run . --align=right something standard")
}
