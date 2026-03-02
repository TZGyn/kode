package color

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
)

type Color color.RGBA

func NewHex(hex string) Color {
	return hexColor(hex)
}

func NewRGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

func NewRGBA(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

func NewRGBFloat(r, g, b float64) Color {
	R := uint8(math.Min(255, (255 * r)))
	G := uint8(math.Min(255, (255 * g)))
	B := uint8(math.Min(255, (255 * b)))

	return Color{R: R, G: G, B: B}
}

func NewRGBAFloat(r, g, b, a float64) Color {

	R := uint8(math.Min(255, (255 * r)))
	G := uint8(math.Min(255, (255 * g)))
	B := uint8(math.Min(255, (255 * b)))
	A := uint8(math.Min(255, (255 * a)))

	return Color{R: R, G: G, B: B, A: A}
}

// HexColor converts hex color to color.RGBA with "#FFFFFF" format
func hexColor(hex string) Color {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return Color{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func (c Color) ToHex(bg Color) string {
	alpha := float64(c.A) / 255.0
	invAlpha := 1 - alpha

	r := uint8(math.Round(float64(c.R)*alpha + float64(bg.R)*invAlpha))
	g := uint8(math.Round(float64(c.G)*alpha + float64(bg.G)*invAlpha))
	b := uint8(math.Round(float64(c.B)*alpha + float64(bg.B)*invAlpha))

	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
