// Package color applies standard IRC color strings to text.
package color

import "fmt"

const (
	// ColorWhite is white.
	ColorWhite = 0

	// ColorBlack is black.
	ColorBlack = 1

	// ColorBlue is blue.
	ColorBlue = 2

	// ColorGreen is green.
	ColorGreen = 3

	// ColorRed is red.
	ColorRed = 4

	// ColorBrown is brown.
	ColorBrown = 5

	// ColorMagenta is magenta.
	ColorMagenta = 6

	// ColorOrange is orange.
	ColorOrange = 7

	// ColorYellow is yellow.
	ColorYellow = 8

	// ColorLightGreen is light green.
	ColorLightGreen = 9

	// ColorCyan is cyan.
	ColorCyan = 10

	// ColorLightCyan is light cyan.
	ColorLightCyan = 11

	// ColorLightBlue is light blue.
	ColorLightBlue = 12

	// ColorPink is pink.
	ColorPink = 13

	// ColorGrey is grey.
	ColorGrey = 14

	// ColorLightGrey is light grey.
	ColorLightGrey = 15
)

// Colors may be passed to the Color() function.
//
// To create a new Colors struct, you should use the New helper:
//
//   colors := &color.Colors{
//     Foreground: color.New(colors.ColorRed),
//     Background: color.New(colors.ColorBlack),
//   }
type Colors struct {
	Foreground *rune
	Background *rune
}

// New creates a new color value from an input color code.
func New(color rune) *rune {
	return &color
}

// Color can set both a foreground and background color on an input string.
//
// To set only a foreground color:
//   coloredText := color.Color("my text", &color.Color{Foreground: color.ColorRed})
//
// To set both foreground and background colors:
//
//   coloredText := color.Colour("my text", &color.Color{
//     Foreground: color.New(color.ColorRed),
//     Background: color.New(color.ColorBlack)
//   })
//
// The first color provided will be used for the foreground, and if a second is
// provided, it will be used for the background.
func Color(input string, colors *Colors) string {
	var code string

	if colors.Foreground != nil {
		code = fmt.Sprintf("\x03%d", *colors.Foreground)

		// A background color cannot be provided without a foreground one
		if colors.Background != nil {
			code = code + fmt.Sprintf(",%d", *colors.Background)
		}
	}

	return fmt.Sprintf("%s%s\x03", code, input)
}

// White makes text white.
func White(input string) string {
	return applyColor(input, ColorWhite)
}

// Black makes text black.
func Black(input string) string {
	return applyColor(input, ColorBlack)
}

// Blue makes text blue.
func Blue(input string) string {
	return applyColor(input, ColorBlue)
}

// Green makes text green.
func Green(input string) string {
	return applyColor(input, ColorGreen)
}

// Red makes text red.
func Red(input string) string {
	return applyColor(input, ColorRed)
}

// Brown makes text brown.
func Brown(input string) string {
	return applyColor(input, ColorBrown)
}

// Magenta makes text magenta.
func Magenta(input string) string {
	return applyColor(input, ColorMagenta)
}

// Orange makes text orange.
func Orange(input string) string {
	return applyColor(input, ColorOrange)
}

// Yellow makes text yellow.
func Yellow(input string) string {
	return applyColor(input, ColorYellow)
}

// LightGreen makes text light green.
func LightGreen(input string) string {
	return applyColor(input, ColorLightGreen)
}

// Cyan makes text cyan.
func Cyan(input string) string {
	return applyColor(input, ColorCyan)
}

// LightCyan makes text light cyan.
func LightCyan(input string) string {
	return applyColor(input, ColorLightCyan)
}

// LightBlue makes text light blue.
func LightBlue(input string) string {
	return applyColor(input, ColorLightBlue)
}

// Pink makes text pink.
func Pink(input string) string {
	return applyColor(input, ColorPink)
}

// Grey makes text grey.
func Grey(input string) string {
	return applyColor(input, ColorGrey)
}

// LightGrey makes text light grey.
func LightGrey(input string) string {
	return applyColor(input, ColorLightGrey)
}

func applyColor(input string, color rune) string {
	return fmt.Sprintf("\x03%d%s\x03", color, input)
}
