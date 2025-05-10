package main

import (
	"bufio"   // Для построчного чтения файла
	"fmt"     // Для форматированного ввода-вывода
	"os"      // Для работы с файлами и аргументами командной строки
	"strings" // Для работы со строками
)

// Структуры данных:“
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
// filename: имя файла шрифта для загрузки (например, "standard.txt")
func (a *ASCIIArt) LoadFont(filename string) error {
	// Открываем файл шрифта
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close() // Гарантируем закрытие файла после завершения функции
	// Список всех поддерживаемых символов в порядке ASCII
	// Этот массив определяет порядок, в котором символы читаются из файла шрифта
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
	var currentLines [8]string // Хранит текущий обрабатываемый символ
	lineIndex := 0             // Текущая строка внутри символа (0-7)
	charIndex := 0             // Текущий обрабатываемый символ
	// Читаем файл построчно
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // Пустая строка означает конец текущего символа
			if lineIndex > 0 && charIndex < len(supportedChars) {
				a.chars[supportedChars[charIndex]] = ASCIIChar{lines: currentLines}
				charIndex++                // Переходим к следующему символу
				lineIndex = 0              // Сбрасываем счётчик строк
				currentLines = [8]string{} // Очищаем строки для следующего символа
			}
		} else {
			if lineIndex < 8 { // Сохраняем строку, если не превышен лимит в 8 строк
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
// Возвращает полный ASCII-арт в виде строки
func (a *ASCIIArt) RenderText(input string) string {
	var result strings.Builder // Эффективно строит результирующую строку
	// Обработка особых случаев
	if input == "" {
		return "" // Пустой ввод
	}
	// if input "\n" then return "$"
	if input == "\\n" {
		return "$" // Только перенос строки
	}

	// Разделяем ввод на строки по символам \n
	lines := strings.Split(strings.ReplaceAll(input, "\\n", "\n"), "\n")
	// Обрабатываем каждую строку ввода
	for i, line := range lines {
		if line == "" { // Обработка пустых строк
			result.WriteString("\n")
			continue
		}
		// Преобразуем текущую строку в ASCII-арт
		artLines := [8]string{} // Хранит 8 строк текущего блока
		for _, char := range line {
			if art, exists := a.chars[char]; exists {
				// Строим каждую строку ASCII-арта
				for j := 0; j < 8; j++ {
					artLines[j] += art.lines[j]
				}
			}
		}
		// Добавляем $ в конец каждой строки и соединяем переносами строк
		for j := 0; j < 8; j++ {
			result.WriteString(artLines[j])
			result.WriteString("\n")
		}
		// Обрабатываем пробелы между блоками текста
		// Добавляем дополнительный пробел только если не в конце или если ввод заканчивается на \n
		if i < len(lines)-1 || (len(input) >= 2 && input[len(input)-2:] == "\\n") {
			result.WriteString("")
		}
	}
	return result.String()
}

func main() {
	if len(os.Args) < 2 { // Проверяем правильное количество аргументов
		return
	}

	// Получаем имя файла шрифта из аргументов командной строки
	input := os.Args[1] // Получаем входной текст из аргументов командной строки
	fontFile := "standard.txt"

	// Проверяем, какой тип шрифта был указан
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "standard":
			fontFile = "standard.txt"
		case "shadow":
			fontFile = "shadow.txt"
		case "thinkertoy":
			fontFile = "thinkertoy.txt"
		default:
			fmt.Printf("Error: Unknown font type '%s'. Supported types are: standard, shadow, thinkertoy.\n", os.Args[2])
			return
		}
	}

	// Создаём новый обработчик ASCII-арта и загружаем шрифт
	ascii := NewASCIIArt()
	if err := ascii.LoadFont(fontFile); err != nil {
		fmt.Printf("Error loading font file '%s': %v\n", fontFile, err)
		return
	}
	// Преобразуем входной текст в ASCII-арт и выводим
	output := ascii.RenderText(input)
	fmt.Print(output)
}
