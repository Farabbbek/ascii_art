package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Создаем структуру данных:
type ASCIIChar struct {
	lines [8]string // Хранит 8 строк ASCII-символа
}

// ASCIIArt - основная структура, которая хранит весь шрифт
type ASCIIArt struct {
	chars map[rune]ASCIIChar // Карта, связывающая каждый символ (руну) с его ASCII-представлением
}

// NewASCIIArt - создает новый экземпляр структуры
func NewASCIIArt() *ASCIIArt {
	return &ASCIIArt{
		chars: make(map[rune]ASCIIChar), // Инициализируем пустую карту для хранения символов
	}
}

// LoadFont читает файл шрифта и загружает все символы в память
func (a *ASCIIArt) LoadFont(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close()
	// Список всех поддерживаемых символов в порядке ASCII
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
	scanner := bufio.NewScanner(file)
	var currentLines [8]string
	lineIndex := 0
	charIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if lineIndex > 0 && charIndex < len(supportedChars) {
				a.chars[supportedChars[charIndex]] = ASCIIChar{lines: currentLines}
				charIndex++
				lineIndex = 0
				currentLines = [8]string{}
			}
		} else {
			if lineIndex < 8 {
				currentLines[lineIndex] = line
				lineIndex++
			}
		}
	}
	// Обрабатываем последний символ в файле
	if lineIndex > 0 && charIndex < len(supportedChars) {
		a.chars[supportedChars[charIndex]] = ASCIIChar{lines: currentLines}
	}
	return nil
}

// RenderText преобразует входной текст в ASCII-арт представление
func (a *ASCIIArt) RenderText(input string) string {
	var result strings.Builder
	// Обработка особых случаев
	if input == "" {
		return ""
	}

	if input == "\\n" {
		return ""
	}

	lines := strings.Split(input, "\\n")

	for i, line := range lines {
		if line == "" { // Обработка пустых строк
			result.WriteString("\n")
			continue
		}

		artLines := [8]string{}
		for _, char := range line {
			if art, exists := a.chars[char]; exists {
				for j := 0; j < 8; j++ {
					artLines[j] += art.lines[j]
				}
			}
		}
		// Добавляем строки ASCII-арта, соединяя переносами строк
		for j := 0; j < 8; j++ {
			result.WriteString(artLines[j])
			result.WriteString("\n")
		}
		// Обрабатываем пробелы между блоками текста
		if i < len(lines)-1 || (len(input) >= 2 && input[len(input)-2:] == "\\n") {
			result.WriteString("")
		}
	}
	return result.String()
}

func main() {
	// Проверяем есть ли необходимое нам число аргументов (2-4)
	if len(os.Args) < 2 || len(os.Args) > 4 {
		printUsage()
		return
	}

	var outputFile string
	var text, bannerType string

	// Разбор аргументов командной строки в зависимости от их количества
	switch len(os.Args) {
	case 2:
		// Случай с двумя аргументами (только строка)
		text = os.Args[1]
		bannerType = "standard"

	case 3:
		// Случай с тремя аргументами (может быть либо флаг --output + строка, либо строка + тип баннера)
		if strings.HasPrefix(os.Args[1], "--output=") {
			if !isValidOutputFlag(os.Args[1]) {
				printUsage() // Если флаг некорректный, выводим инструкцию по использованию
				return
			}
			// Извлекаем имя файла для вывода (удаляем префикс "--output=")
			outputFile = strings.TrimPrefix(os.Args[1], "--output=")
			text = os.Args[2]
			// Тип баннера по умолчанию
			bannerType = "standard"
		} else {
			// Если флаг не присутствует, то первый аргумент — это текст, второй — тип баннера
			text = os.Args[1]
			bannerType = os.Args[2]
		}

	case 4:
		// Случай с четырьмя аргументами (--output флаг + строка + тип баннера)
		if !isValidOutputFlag(os.Args[1]) {
			printUsage()
			return
		}
		// Извлекаем имя файла для вывода
		outputFile = strings.TrimPrefix(os.Args[1], "--output=")
		text = os.Args[2]
		bannerType = os.Args[3]
	}

	// Создаем новый процессор для ASCII-арта и загружаем шрифт
	ascii := NewASCIIArt()
	if err := ascii.LoadFont(bannerType + ".txt"); err != nil {
		fmt.Printf("Ошибка при загрузке шрифта: %v\n", err)
		return
	}

	// Генерируем ASCII-арт для заданного текста
	output := ascii.RenderText(text)

	// Обрабатываем вывод в зависимости от того, указан ли файл для вывода
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Printf("Ошибка при записи в файл: %v\n", err)
			return
		}
	} else {
		fmt.Print(output)
	}
}

// Вспомогательная функция для вывода инструкции по использованию
func printUsage() {
	fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
	fmt.Println("\nEX: go run . --output=<fileName.txt> something standard")
}

// Вспомогательная функция для проверки формата флага --output
func isValidOutputFlag(flag string) bool {
	// Проверка, что флаг начинается с "--output=" и имеет длину больше 9 (для наличия самого пути)
	return strings.HasPrefix(flag, "--output=") && len(flag) > 9
}
