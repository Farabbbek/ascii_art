package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Константы ANSI для цветов
const (
	reset  = "\033[0m"
	red    = "\033[31m"
	orange = "\033[38;5;208m"
	yellow = "\033[33m"
	green  = "\033[32m"
	blue   = "\033[34m"
	indigo = "\033[38;5;54m"
	violet = "\033[35m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	white  = "\033[37m"
)

// ColorConfig хранит данные о цветах и строках, которые нужно окрасить
type ColorConfig struct {
	color     string
	substring string
	enabled   bool //флаг, включено ли раскрашивание.
}

// ASCIIChar представляет один символ ASCII, состоящий из 8 строк.
type ASCIIChar struct {
	lines [8]string
}

// ASCIIArt хранит карту символов и их ASCII-представления
type ASCIIArt struct {
	chars map[rune]ASCIIChar
}

// NewASCIIArt создаёт и возвращает новый экземпляр структуры ASCIIArt
func NewASCIIArt() *ASCIIArt {
	return &ASCIIArt{
		chars: make(map[rune]ASCIIChar),
	}
}

// getColorCode возвращает соответствующий ANSI-код для этого цвета
func getColorCode(color string) string {
	colorMap := map[string]string{
		"red":    red,
		"orange": orange,
		"yellow": yellow,
		"green":  green,
		"blue":   blue,
		"indigo": indigo,
		"violet": violet,
		"purple": purple,
		"cyan":   cyan,
		"white":  white,
	}
	if code, exists := colorMap[strings.ToLower(color)]; exists {
		return code
	}
	return white // цвет по умолчанию
}

// parseArgs парсинг аргументов командной строки
func parseArgs(args []string) (ColorConfig, string, string, error) {
	var colorConfig ColorConfig
	var text, banner string

	if len(args) < 2 {
		return colorConfig, "", "", fmt.Errorf("insufficient arguments")
	}

	// Проверка наличия опции цвета
	if strings.HasPrefix(args[1], "--color=") {
		colorConfig.enabled = true
		colorConfig.color = strings.TrimPrefix(args[1], "--color=")

		if len(args) == 3 {
			// Пример: go run . --color=blue "hello"
			text = args[2]
			banner = "standard" // баннер по умолчанию
		} else if len(args) == 4 {
			if strings.HasSuffix(args[3], ".txt") || args[3] == "standard" || args[3] == "shadow" || args[3] == "thinkertoy" {
				// Пример: go run . --color=green "hello" thinkertoy
				text = args[2]
				banner = args[3]
			} else {
				// Пример: go run . --color=red kit "a king kitten have kit"
				colorConfig.substring = args[2]
				text = args[3]
				banner = "standard" // Default banner
			}
		} else if len(args) == 5 {
			// Пример: go run . --color=red h "hello" standard
			colorConfig.substring = args[2]
			text = args[3]
			banner = args[4]
		} else {
			return colorConfig, "", "", fmt.Errorf("invalid number of arguments")
		}
	} else {
		// Пример: go run . "hello" standard
		text = args[1]
		if len(args) >= 3 {
			banner = args[2]
		} else {
			banner = "standard" // баннер по умолчанию
		}
	}

	return colorConfig, text, banner, nil
}

// LoadFont загружает шрифты из файла
func (a *ASCIIArt) LoadFont(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close()

	// Список символов, которые поддерживаются шрифтом
	supportedChars := []rune{
		' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		':', ';', '<', '=', '>', '?', '@',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'[', '\\', ']', '^', '_', '`',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'{', '|', '}', '~',
	}

	// Сканируем файл построчно
	scanner := bufio.NewScanner(file)
	var currentLines [8]string
	lineIndex := 0
	charIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// Если символ уже полностью прочитан, сохраняем его в map
			if lineIndex > 0 && charIndex < len(supportedChars) {
				a.chars[supportedChars[charIndex]] = ASCIIChar{lines: currentLines}
				charIndex++
				lineIndex = 0
				currentLines = [8]string{}
			}
		} else {
			// Если строка не пустая, добавляем ее в массив для текущего символа
			if lineIndex < 8 {
				currentLines[lineIndex] = line
				lineIndex++
			}
		}
	}

	// Добавляем последний символ, если файл не завершен пустой строкой
	if lineIndex > 0 && charIndex < len(supportedChars) {
		a.chars[supportedChars[charIndex]] = ASCIIChar{lines: currentLines}
	}

	return nil
}

// RenderText генерирует ASCII-арт из переданного текста с учетом цветовой конфигурации
func (a *ASCIIArt) RenderText(input string, colorConfig ColorConfig) string {
	var result strings.Builder

	// Если текст пустой, просто возвращаем пустую строку.
	if input == "" {
		return ""
	}

	// Если введен специальный символ для новой строки, заменяем его на стандартный символ новой строки
	if input == "\\n" {
		return "\n"
	}

	// Разбиваем входной текст на строки, разделенные символами \n
	lines := strings.Split(input, "\\n")
	colorCode := getColorCode(colorConfig.color)

	for lineNum, line := range lines {
		if line == "" {
			result.WriteString("\n")
			continue
		}

		// Массив для хранения 8 строк ASCII-арта для каждой строки текста
		var artLines [8]strings.Builder
		var colorPositions []int

		// Если цветная подсветка включена и подстрока задана, ищем все позиции для окраски
		if colorConfig.enabled && colorConfig.substring != "" {
			startIndex := 0
			for {
				index := strings.Index(strings.ToLower(line[startIndex:]), strings.ToLower(colorConfig.substring))
				if index == -1 {
					break
				}
				colorPositions = append(colorPositions, startIndex+index)
				startIndex += index + 1
			}
		}

		for charIdx, char := range line {
			// Проверяем, есть ли ASCII-арт для текущего символа
			if art, exists := a.chars[char]; exists {
				// Проверяем, нужно ли применить цвет для этого символа
				shouldColor := colorConfig.enabled && (colorConfig.substring == "" ||
					(func() bool {
						// Проверяем, попадает ли текущий символ в область подсветки
						for _, pos := range colorPositions {
							if charIdx >= pos && charIdx < pos+len(colorConfig.substring) {
								return true
							}
						}
						return false
					})())

				// Добавляем все строки символа в финальный результат
				for lineIdx := 0; lineIdx < 8; lineIdx++ {
					if shouldColor {
						artLines[lineIdx].WriteString(colorCode + art.lines[lineIdx] + reset)
					} else {
						artLines[lineIdx].WriteString(art.lines[lineIdx])
					}
				}
			}
		}

		// Если в конце строки имеется специальный символ для новой строки, добавляем пустую строку
		for lineIdx := 0; lineIdx < 8; lineIdx++ {
			result.WriteString(artLines[lineIdx].String())
			result.WriteString("\n")
		}

		// Возвращаем сгенерированный ASCII-арт
		if lineNum < len(lines)-1 || (len(input) >= 2 && input[len(input)-2:] == "\\n") {
			result.WriteString("")
		}
	}

	return result.String()
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  1. go run . [STRING] [BANNER]")
	fmt.Println("  2. go run . [OPTION] [STRING]")
	fmt.Println("  3. go run . [OPTION] [STRING] [BANNER]")
	fmt.Println("\nExamples:")
	fmt.Println("  go run . \"hello\" standard")
	fmt.Println("  go run . --color=red kit \"a king kitten have kit\"")
	fmt.Println("  go run . --color=red h \"hello\" standard")
	fmt.Println("  go run . --color=green \"hello\" thinkertoy")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	colorConfig, text, bannerType, err := parseArgs(os.Args)
	if err != nil {
		printUsage()
		return
	}

	var fontFile string
	switch bannerType {
	case "standard", "standard.txt":
		fontFile = "standard.txt"
	case "shadow", "shadow.txt":
		fontFile = "shadow.txt"
	case "thinkertoy", "thinkertoy.txt":
		fontFile = "thinkertoy.txt"
	default:
		fmt.Printf("Error: Unknown font type '%s'. Supported types are: standard, shadow, thinkertoy.\n", bannerType)
		return
	}

	ascii := NewASCIIArt()
	if err := ascii.LoadFont(fontFile); err != nil {
		fmt.Printf("Error loading font file '%s': %v\n", fontFile, err)
		return
	}

	output := ascii.RenderText(text, colorConfig)
	fmt.Print(output)
}
