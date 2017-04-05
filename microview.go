/*

Package microview is used to remote control a MicroView - OLED Arduino Module

Installation

    go get -u github.com/peterhellberg/microview

Example

    package main

    import (
    	"flag"
    	"log"
    	"time"

    	. "github.com/peterhellberg/microview"
    )

    func main() {
    	name := flag.String("name", "/dev/cu.usbserial-DA00SSM3", "name of the serial port")

    	flag.Parse()

    	mv, err := OpenMicroView(*name, Delay(90*time.Millisecond))
    	if err != nil {
    		log.Fatal(err)
    	}

    	mv.Run(
    		RectFill(5, 5, 5, 15),
    		RectFill(25, 0, 30, 15),
    		Rect(1, 1, 20, 40),
    		Rect(40, 20, 20, 20),
    		Rect(40, 20, 15, 15),
    		Rect(40, 20, 10, 10),
    		Rect(40, 20, 5, 5),
    	)
    }

*/
package microview

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"time"

	term "github.com/pkg/term"
)

// Size
const (
	WIDTH  = 64
	HEIGHT = 48
)

// Colors
const (
	BLACK = 0
	WHITE = 1
)

// Draw modes
const (
	NORM = 0
	XOR  = 1
)

// Clear modes
const (
	PAGE = 0
	ALL  = 1
)

// Command identifiers
const (
	CLEAR        = 0
	INVERT       = 1
	CONTRAST     = 2
	DISPLAY      = 3
	SETCURSOR    = 4
	PIXEL        = 5
	LINE         = 6
	LINEH        = 7
	LINEV        = 8
	RECT         = 9
	RECTFILL     = 10
	CIRCLE       = 11
	CIRCLEFILL   = 12
	DRAWCHAR     = 13
	DRAWBITMAP   = 14 // Not implemented yet
	GETLCDWIDTH  = 15
	GETLCDHEIGHT = 16
	SETCOLOR     = 17
	SETDRAWMODE  = 18
)

// Command is a command to be sent to the MicroView
type Command []byte

// MicroView represents a remote MicroView micro controller
type MicroView struct {
	io.ReadWriteCloser
	delay time.Duration
}

// NewMicroView creates a new MicroView instance
func NewMicroView(rwc io.ReadWriteCloser, options ...func(*MicroView)) (*MicroView, error) {
	mv := &MicroView{
		ReadWriteCloser: rwc,
		delay:           50 * time.Millisecond,
	}

	for _, option := range options {
		option(mv)
	}

	return mv, nil
}

// OpenMicroView opens a terminal connection to a MicroView
func OpenMicroView(name string, options ...func(*MicroView)) (*MicroView, error) {
	term, err := term.Open(name, term.Speed(115200))
	if err != nil {
		return nil, err
	}

	// Read the welcome message from the MicroView
	term.Read([]byte("MicroView"))

	return NewMicroView(term, options...)
}

// Delay between each command sent to the MicroView (25ms seems to be the minimum delay)
func Delay(d time.Duration) func(*MicroView) {
	return func(mv *MicroView) {
		mv.delay = d
	}
}

// Run commands
func (mv *MicroView) Run(cmds ...Command) {
	for _, cmd := range cmds {
		mv.Write(cmd)

		time.Sleep(mv.delay)
	}
}

// DrawString draws the provided string at x,y
func (mv *MicroView) DrawString(x, y uint8, s string) {
	for i, r := range s {
		mv.Write(DrawChar(x+uint8(i*6), y, r))

		time.Sleep(mv.delay)
	}
}

// Set the pixel at x,y to WHITE or BLACK based on the provided color
func (mv *MicroView) Set(x, y int, c color.Color) {
	if r, g, b, _ := c.RGBA(); r+g+b > 0 {
		mv.Write(PixelWithColorAndMode(uint8(x), uint8(y), WHITE, NORM))
		time.Sleep(mv.delay)
		return
	}

	time.Sleep(5 * time.Millisecond)
}

// At returns the color at x,y (always black for now)
func (mv *MicroView) At(x, y int) color.Color {
	return color.Black
}

// ColorModel of the MicroView
func (mv *MicroView) ColorModel() color.Model {
	return color.GrayModel
}

// Bounds of the MicroView (it is 64x48)
func (mv *MicroView) Bounds() image.Rectangle {
	return image.Rect(0, 0, WIDTH, HEIGHT)
}

// Clear clears the display (ALL or PAGE)
// CMD_CLEAR, 0
//
// To clear GDRAM inside the LCD controller, pass in the variable mode = ALL
// and to clear screen page buffer, pass in the variable mode = PAGE
func Clear(mode uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d", CLEAR, mode))
}

// Invert inverts the display
// CMD_INVERT, 1
//
// The WHITE color of the display will turn to
// BLACK and the BLACK will turn to WHITE.
func Invert(inv bool) []byte {
	var i int

	if inv {
		i = 1
	}

	return []byte(fmt.Sprintf("%d,%d", INVERT, i))
}

// Contrast sets the contrast
// CMD_CONTRAST, 2
//
// OLED contrast value from 0 to 255.
//
// Note: Contrast level is not very obvious.
func Contrast(c uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d", CONTRAST, c))
}

// Display the buffer on the OLED
// CMD_DISPLAY, 3
//
// Bulk move the screen buffer to the SSD1306 controller's memory so that
// images/graphics drawn on the screen buffer will be displayed on the OLED.
func Display() []byte {
	return []byte(fmt.Sprintf("%d", DISPLAY))
}

// SetCursor sets the cursor position
// CMD_SETCURSOR,	4
//
// MicroView's cursor position to x,y.
func SetCursor(x, y uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d", SETCURSOR, x, y))
}

// Pixel draws a pixel at x,y
// CMD_PIXEL, 5
//
// Draw pixel using the current fore color and current draw mode
// in the screen buffer's x,y position.
func Pixel(x, y uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d", PIXEL, x, y))
}

// PixelWithColorAndMode draws pixel at x,y with color and mode
// CMD_PIXEL, 5
//
// Draw color pixel in the screen buffer's x,y position with NORM or XOR draw mode.
func PixelWithColorAndMode(x, y, color, mode uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d", PIXEL, x, y, color, mode))
}

// Line draws a line
// CMD_LINE, 6
//
// Draw line using current fore color and current draw mode
// from x0,y0 to x1,y1 of the screen buffer.
func Line(x0, y0, x1, y1 uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d", LINE, x0, y0, x1, y1))
}

// LineWithColorAndMode draws a line with color and mode
// CMD_LINE, 6
//
// Draw line using color and mode from x0,y0 to x1,y1 of the screen buffer.
func LineWithColorAndMode(x0, y0, x1, y1, color, mode uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d", LINE, x0, y0, x1, y1, color, mode))
}

// LineH draws a horizontal line
// CMD_LINEH, 7
//
// Draw horizontal line using current fore color and current draw mode
// from x,y to x+width,y of the screen buffer.
func LineH(x, y, w uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d", LINEH, x, y, w))
}

// LineHWithColorAndMode draws a horizontal line with color and mode.
// CMD_LINEH, 7
//
// Draw horizontal line using color and mode from x,y
// to x+width,y of the screen buffer.
func LineHWithColorAndMode(x, y, w, color, mode uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d,%d", LINEH, x, y, w, color, mode))
}

// LineV draws a vertical line
// CMD_LINEV, 8
//
// Draw vertical line using current fore color and current draw mode from x,y
// to x,y+height of the screen buffer.
func LineV(x, y, h uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d", LINEV, x, y, h))
}

// LineVWithColorAndMode draws a vertical line with color and mode
// CMD_LINEV, 8
//
// Draw vertical line using color and mode from x,y
// to x,y+height of the screen buffer.
func LineVWithColorAndMode(x, y, h, color, mode uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d,%d", LINEV, x, y, h, color, mode))
}

// Rect draws a rectangle
// CMD_RECT, 9
//
// Draw rectangle using current fore color and current draw mode from x,y
// to x+width,y+height of the screen buffer.
func Rect(x, y, w, h uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d", RECT, x, y, w, h))
}

// RectWithColorAndMode draws a rectangle with color and mode
// CMD_RECT, 9
//
// Draw rectangle using color and mode from x,y
// to x+width,y+height of the screen buffer.
func RectWithColorAndMode(x, y, w, h, color, mode uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d", RECT, x, y, w, h, color, mode))
}

// RectFill draws a filled rectangle.
// CMD_RECTFILL, 10
//
// Draw filled rectangle using current fore color and current draw mode from x,y
// to x+width,y+height of the screen buffer.
func RectFill(x, y, w, h uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d", RECTFILL, x, y, w, h))
}

// Fill the screen with the current fore color
// CMD_RECTFILL, 10
func Fill() []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d", RECTFILL, 0, 0, 64, 48))
}

// RectFillWithColorAndMode draws a filled rectangle with color and mode
// CMD_RECT, 10
//
// Draw filled rectangle using color and mode from x,y
// to x+width,y+height of the screen buffer.
func RectFillWithColorAndMode(x, y, w, h, color, mode uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d", RECTFILL, x, y, w, h, color, mode))
}

// Circle draws a circle
// CMD_CIRCLE, 11
//
// Draw circle with radius using current fore color
// and current draw mode at x,y of the screen buffer.
func Circle(x, y, r uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d", CIRCLE, x, y, r))
}

// CircleFill draws a filled circle
// CMD_CIRCLEFILL, 12
//
// Draw filled circle with radius using current fore color
// and current draw mode at x,y of the screen buffer.
func CircleFill(x, y, r uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%d", CIRCLEFILL, x, y, r))
}

// DrawChar draws a character
// CMD_DRAWCHAR, 13
//
// Draw character c using current color and current draw mode at x,y.
func DrawChar(x, y uint8, c rune) []byte {
	return []byte(fmt.Sprintf("%d,%d,%d,%v", DRAWCHAR, x, y, c))
}

// SetColor sets the color
// CMD_SETCOLOR, 17
//
// Set the current draw's color. Only WHITE and BLACK available.
func SetColor(c uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d", SETCOLOR, c))
}

// SetDrawMode sets the draw mode
// CMD_SETDRAWMODE, 18
//
// Set current draw mode with NORM or XOR.
func SetDrawMode(mode uint8) []byte {
	return []byte(fmt.Sprintf("%d,%d", SETDRAWMODE, mode))
}
