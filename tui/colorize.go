package tui

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func Colorize(content *string) (string, error) {
	style := getStyle()
	formatter, err := getFormatter()
	if formatter == nil {
		return "", fmt.Errorf("terminal source code formatter not found")
	}

	lexer := lexers.Analyse(*content)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := tokenizeInput(lexer, content)
	if err != nil {
		return "", err
	}

	output, err := formatTokens(formatter, style, iterator)
	if err != nil {
		return "", err
	}

	return output, nil
}

func getStyle() *chroma.Style {
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	return style
}

func getFormatter() (chroma.Formatter, error) {
	formatter := formatters.Get("terminal")
	if formatter == nil {
		return nil, fmt.Errorf("terminal formatter not found")
	}
	return formatter, nil
}

func tokenizeInput(lexer chroma.Lexer, content *string) (*chroma.Iterator, error) {
	iterator, err := lexer.Tokenise(nil, *content)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize content: %v", err)
	}
	return &iterator, nil
}

func formatTokens(formatter chroma.Formatter, style *chroma.Style, iterator *chroma.Iterator) (string, error) {
	var b bytes.Buffer
	err := formatter.Format(&b, style, *iterator)
	if err != nil {
		return "", fmt.Errorf("failed to format tokens: %v", err)
	}
	return b.String(), nil
}
