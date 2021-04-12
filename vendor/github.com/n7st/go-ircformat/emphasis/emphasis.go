// Package emphasis applies emphasis (bold, italic etc) to text.
package emphasis

import "fmt"

// Bold makes text bold.
func Bold(input string) string {
	return fmt.Sprintf("\x02%s\x02", input)
}

// Italic makes text italic.
func Italic(input string) string {
	return fmt.Sprintf("\x1D%s\x1D", input)
}

// Underline underlines text.
func Underline(input string) string {
	return fmt.Sprintf("\x1F%s\x1F", input)
}

// Strikethrough strikes text out.
func Strikethrough(input string) string {
	return fmt.Sprintf("\x1E%s\x1E", input)
}

// Monospace sets text in a monospace font.
func Monospace(input string) string {
	return fmt.Sprintf("\x11%s\x11", input)
}
