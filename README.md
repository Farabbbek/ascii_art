# ASCII Art Project

A collection of utilities for creating and manipulating ASCII art with different fonts and colors.

## Features
- Multiple ASCII art fonts:
  - Standard
  - Shadow
  - Thinkertoy
- Color support
- Text justification
- File system operations

## Supported Fonts

- `standard`: Classic ASCII art style
- `shadow`: Text with shadow effect
- `thinkertoy`: Minimalist ASCII art style

## Color 
```sh
cd color 
go run . --color=<color> <substring to be colored > "something"
go run . --color=red kit "a king kitten have kit"
```

## fs 
```sh
cd fs
Usage: go run . [STRING] [BANNER]
EX: go run . something standard
```

## justify
```sh
cd justify
Usage: go run . [OPTION] [STRING] [BANNER]
Example: go run . --align=right something standard
```
## Output 
```sh
cd output
Usage: go run . [OPTION] [STRING] [BANNER]
EX: go run . --output=<fileName.txt> something standard
```

