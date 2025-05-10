# ASCII Art Justify

A command-line tool that generates ASCII art from text with different alignment options (left, right, center, justify).

## Description

This program takes a string as input and converts it into ASCII art using predefined banner templates. It supports different text alignment options:
- Left (default)
- Right
- Center
- Justify

## Installation

```bash
git clone https://github.com/atusup/ascii-art-justify.git
cd ascii-art-justify
```

## Usage

Basic usage:
```bash
go run . [OPTION] [STRING] [BANNER]
```

### Parameters:
- `OPTION`: Alignment option (--align=left|right|center|justify)
- `STRING`: The text to convert to ASCII art
- `BANNER`: Banner style (standard, shadow, thinkertoy)

### Examples:

Left alignment (default):
```bash
go run main.go "Hello"
```

Right alignment:
```bash
go run main.go --align=right "Hello" standard
```

Center alignment:
```bash
go run main.go  --align=center "Hello" standard
```

Justify alignment:
```bash
go run main.go --align=justify "Hello World" standard
```

## Project Structure

```
ascii-art-justify/
├─main.go
└── banner/
    └── standard.txt
    └── shadow.txt
    └── thinkertoy.txt
```

## Features

- Multiple text alignment options
- Support for different banner styles
- Clean and modular code structure
- Error handling for invalid inputs
- Default width of 225 characters

## Technical Details

- Written in Go
- Uses file I/O for banner templates
- Efficient string manipulation

## Error Handling

The program handles various error cases:
- Invalid alignment options
- Missing arguments
- Invalid banner files
- File reading errors