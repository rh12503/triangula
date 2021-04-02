// Package draw contains the RunWindow function which you call from runners to open
// a window. You pass it a callback which is called at 60 frames per second. The
// callback gives you a Window which you use to do input handling, rendering and
// audio output. See the documentation for Window for more details on what you
// can do.
package draw

import "strconv"

// UpdateFunction is used as a callback when creating a window. It is called
// at 60Hz and you do all your event handling and drawing in it.
type UpdateFunction func(window Window)

// Window provides functions to draw simple primitives and images, handle
// keyboard and mouse events and play sounds.
// All drawing functions that have width and height as input expect those to be
// positive. Objects with negative width or height will silently be ignored and
// not drawn.
type Window interface {
	// Close closes the window which will stop the update loop after the current
	// frame (you will usually want to return from the update function after
	// calling Close or another frame will be displayed).
	Close()

	// Size returns the window's size in pixels.
	Size() (width, height int)

	// SetFullscreen toggles between the fixed-size window with title and border
	// and going full screen on the monitor that the window is placed on when
	// the call to SetFullscreen(true) occurs.
	// Use Window.Size to get the new size after this.
	// By default the window is not fullscreen. It always starts windowed.
	SetFullscreen(f bool)

	// ShowCursor set the OS' mouse cursor to visible or invisible. It defaults
	// to visible if you do not call ShowCursor.
	ShowCursor(show bool)

	// WasKeyPressed reports whether the specified key was pressed at any time
	// during the last frame. If the user presses a key and releases it in the
	// same frame, this function stores that information and will return true.
	// See the Key... constants for the available keys that can be queried.
	// NOTE do not use this for text input, use Characters instead.
	WasKeyPressed(key Key) bool

	// IsKeyDown reports whether the specified key is being held down at the
	// moment of calling this function.
	// See the Key... constants for the available keys that can be queried.
	IsKeyDown(key Key) bool

	// Characters returns all pressed keys translated to characters that
	// happened in the last frame. The runes in the string are ordered by the
	// time that the keys were entered.
	Characters() string

	// IsMouseDown reports whether the specified button is down at the time of
	// the function call
	IsMouseDown(button MouseButton) bool

	// Clicks returns all MouseClicks that occurred during the last frame.
	Clicks() []MouseClick

	// MousePositoin returns the current mouse position in pixels at the time of
	// the function call. It is relative to the drawing area of the window.
	MousePosition() (x, y int)

	// MouseWheelY returns the aggregate vertical mouse wheel rotation during
	// the last frame. A value of 1 typically corresponds to one tick of the
	// wheel. A positive value means the wheel was rotated forward, away from
	// the user, a negative value means the wheel was rotated backward towards
	// the user.
	MouseWheelY() float64

	// MouseWheelX returns the aggregate horizontal mouse wheel rotation during
	// the last frame. A value of 1 typically corresponds to one tick of the
	// wheel. A positive value means the wheel was rotated right, a negative
	// value means the wheel was rotated left.
	MouseWheelX() float64

	// DrawPoint draws a single point at the given screen position in pixels.
	DrawPoint(x, y int, color Color)

	// DrawLine draws a one pixel wide line from the first point to the second
	// (inclusive).
	DrawLine(fromX, fromY, toX, toY int, color Color)

	// DrawRect draws a one pixel wide rectangle outline.
	DrawRect(x, y, width, height int, color Color)

	// FillRect draws a filled rect.
	FillRect(x, y, width, height int, color Color)

	// DrawEllipse draws a one pixel wide ellipse. The top-left corner of the
	// surrounding rectangle is given by x and y, the horizontal and vertical
	// diameters are given by width and height.
	DrawEllipse(x, y, width, height int, color Color)

	// FillEllipse behaves like DrawEllipse but fills the ellipse with the color
	// instaed of only drawing the outline.
	FillEllipse(x, y, width, height int, color Color)

	DrawTriangle(x0, y0, x1, y1, x2, y2 int, color Color)

	FillTriangle(x0, y0, x1, y1, x2, y2 int, color Color)

	// DrawImageFile draws the untransformed image at the give position. If the
	// image file is not found or has the wrong format an error is returned.
	DrawImageFile(path string, x, y int) error

	// DrawImageFileTo draws the image to the given screen rectangle, possibly
	// scaling it in either direction, and rotates it around the rectangles
	// center point by the given angle. The rotation is clockwise.
	// If the image file is not found or has the wrong format an error is
	// returned.
	DrawImageFileTo(path string, x, y, w, h, rotationCWDeg int) error

	// DrawImageFileRotated draws the image with its top-left corner at the
	// given coordinates but roatated clockwise about the given angle in degrees
	// around its center. This means its top-left corner will only actually be
	// at the given location if the rotation is 0.
	// If the image file is not found or has the wrong format an error is
	// returned.
	DrawImageFileRotated(path string, x, y, rotationCWDeg int) error

	// DrawImageFilePart lets you specify the source and destination rectangle
	// for the image file to be drawn. The image is rotated about its center by
	// the given angle in degrees, clockwise. You may flip the image by
	// specifying a negative width or height. E.g. to flip in x direction,
	// instead of
	//
	//     DrawImageFilePart("x.png", 0, 0, 100, 100, 0, 0, 100, 100, 0)
	//
	// you would do
	//
	//     DrawImageFilePart("x.png", 100, 0, -100, 100, 0, 0, 100, 100, 0)
	//
	// swapping the source rectangle's left and right positions.
	//
	// If the image file is not found or has the wrong format an error is
	// returned.
	DrawImageFilePart(
		path string,
		sourceX, sourceY, sourceWidth, sourceHeight int,
		destX, destY, destWidth, destHeight int,
		rotationCWDeg int,
	) error

	// GetTextSize returns the size the given text would have when being drawn.
	GetTextSize(text string) (w, h int)

	// GetScaledTextSize returns the size the given text would have when being
	// drawn at the given scale.
	GetScaledTextSize(text string, scale float32) (w, h int)

	// DrawText draws a text string. New line characters ('\n') are not drawn
	// but force a line break and the next character is drawn on the line below
	// starting again at x.
	DrawText(text string, x, y int, color Color)

	// DrawScaledText behaves as DrawText, but the text is scaled. If scale = 1
	// this behaves exactly like DrawText, scales > 1 make the text bigger,
	// scales < 1 shrink it. Scales <= 0 will draw no text at all.
	DrawScaledText(text string, x, y int, scale float32, color Color)

	// PlaySoundFile only plays WAV sounds. If the file is not found or has the
	// wrong format an error is returned.
	PlaySoundFile(path string) error
}

// Color consists of four channels ranging from 0 to 1 each. A specifies the
// opacity, 1 being fully opaque and 0 being fully transparent.
type Color struct{ R, G, B, A float32 }

// RGB creates a color with full opacity. All values are in the range from 0 to
// 1.
func RGB(r, g, b float32) Color {
	return Color{r, g, b, 1}
}

// RGBA creates a color from the given channel values. All values are in the
// range from 0 to 1.
func RGBA(r, g, b, a float32) Color {
	return Color{r, g, b, a}
}

// These are predefined colors for intuitive use, no need to set color channels.
var (
	Black       = Color{0, 0, 0, 1}
	White       = Color{1, 1, 1, 1}
	Gray        = Color{0.5, 0.5, 0.5, 1}
	LightGray   = Color{0.75, 0.75, 0.75, 1}
	DarkGray    = Color{0.25, 0.25, 0.25, 1}
	Red         = Color{1, 0, 0, 1}
	LightRed    = Color{1, 0.5, 0.5, 1}
	DarkRed     = Color{0.5, 0, 0, 1}
	Green       = Color{0, 1, 0, 1}
	LightGreen  = Color{0.5, 1, 0.5, 1}
	DarkGreen   = Color{0, 0.5, 0, 1}
	Blue        = Color{0, 0, 1, 1}
	LightBlue   = Color{0.5, 0.5, 1, 1}
	DarkBlue    = Color{0, 0, 0.5, 1}
	Purple      = Color{1, 0, 1, 1}
	LightPurple = Color{1, 0.5, 1, 1}
	DarkPurple  = Color{0.5, 0, 0.5, 1}
	Yellow      = Color{1, 1, 0, 1}
	LightYellow = Color{1, 1, 0.5, 1}
	DarkYellow  = Color{0.5, 0.5, 0, 1}
	Cyan        = Color{0, 1, 1, 1}
	LightCyan   = Color{0.5, 1, 1, 1}
	DarkCyan    = Color{0, 0.5, 0.5, 1}
	Brown       = Color{0.5, 0.2, 0, 1}
	LightBrown  = Color{0.75, 0.3, 0, 1}
)

// MouseClick is used to store mouse click events.
type MouseClick struct {
	// X and Y are the screen position in pixels, relative to the drawing area.
	// X goes from left to right, starting at 0 and Y goes from top to bottom
	// starting at 0.
	// This means that pixel 0,0 is the top-left pixel in the drawing area (not
	// the title bar).
	X, Y   int
	Button MouseButton
}

// MouseButton is one of the three buttons typically present on a mouse.
type MouseButton int

// These are the possible values for MouseButton.
const (
	LeftButton MouseButton = iota
	MiddleButton
	RightButton

	// NOTE mouseButtonCount has to come last
	mouseButtonCount
)

// Key represents a key on the keyboard.
type Key int

// These are all available keyboard keys.
const (
	KeyA Key = 1 + iota
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyNum0
	KeyNum1
	KeyNum2
	KeyNum3
	KeyNum4
	KeyNum5
	KeyNum6
	KeyNum7
	KeyNum8
	KeyNum9
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyEnter
	KeyNumEnter
	KeyLeftControl
	KeyRightControl
	KeyLeftShift
	KeyRightShift
	KeyLeftAlt
	KeyRightAlt
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyEscape
	KeySpace
	KeyBackspace
	KeyTab
	KeyHome
	KeyEnd
	KeyPageDown
	KeyPageUp
	KeyDelete
	KeyInsert
	KeyNumAdd
	KeyNumSubtract
	KeyNumMultiply
	KeyNumDivide
	KeyCapslock
	KeyPrint
	KeyPause

	// NOTE keyCount has to come last
	keyCount
)

func (k Key) String() string {
	switch k {
	case KeyA:
		return "A"
	case KeyB:
		return "B"
	case KeyC:
		return "C"
	case KeyD:
		return "D"
	case KeyE:
		return "E"
	case KeyF:
		return "F"
	case KeyG:
		return "G"
	case KeyH:
		return "H"
	case KeyI:
		return "I"
	case KeyJ:
		return "J"
	case KeyK:
		return "K"
	case KeyL:
		return "L"
	case KeyM:
		return "M"
	case KeyN:
		return "N"
	case KeyO:
		return "O"
	case KeyP:
		return "P"
	case KeyQ:
		return "Q"
	case KeyR:
		return "R"
	case KeyS:
		return "S"
	case KeyT:
		return "T"
	case KeyU:
		return "U"
	case KeyV:
		return "V"
	case KeyW:
		return "W"
	case KeyX:
		return "X"
	case KeyY:
		return "Y"
	case KeyZ:
		return "Z"
	case Key0:
		return "0"
	case Key1:
		return "1"
	case Key2:
		return "2"
	case Key3:
		return "3"
	case Key4:
		return "4"
	case Key5:
		return "5"
	case Key6:
		return "6"
	case Key7:
		return "7"
	case Key8:
		return "8"
	case Key9:
		return "9"
	case KeyNum0:
		return "Num0"
	case KeyNum1:
		return "Num1"
	case KeyNum2:
		return "Num2"
	case KeyNum3:
		return "Num3"
	case KeyNum4:
		return "Num4"
	case KeyNum5:
		return "Num5"
	case KeyNum6:
		return "Num6"
	case KeyNum7:
		return "Num7"
	case KeyNum8:
		return "Num8"
	case KeyNum9:
		return "Num9"
	case KeyF1:
		return "F1"
	case KeyF2:
		return "F2"
	case KeyF3:
		return "F3"
	case KeyF4:
		return "F4"
	case KeyF5:
		return "F5"
	case KeyF6:
		return "F6"
	case KeyF7:
		return "F7"
	case KeyF8:
		return "F8"
	case KeyF9:
		return "F9"
	case KeyF10:
		return "F10"
	case KeyF11:
		return "F11"
	case KeyF12:
		return "F12"
	case KeyF13:
		return "F13"
	case KeyF14:
		return "F14"
	case KeyF15:
		return "F15"
	case KeyF16:
		return "F16"
	case KeyF17:
		return "F17"
	case KeyF18:
		return "F18"
	case KeyF19:
		return "F19"
	case KeyF20:
		return "F20"
	case KeyF21:
		return "F21"
	case KeyF22:
		return "F22"
	case KeyF23:
		return "F23"
	case KeyF24:
		return "F24"
	case KeyEnter:
		return "Enter"
	case KeyNumEnter:
		return "NumEnter"
	case KeyLeftControl:
		return "LeftControl"
	case KeyRightControl:
		return "RightControl"
	case KeyLeftShift:
		return "LeftShift"
	case KeyRightShift:
		return "RightShift"
	case KeyLeftAlt:
		return "LeftAlt"
	case KeyRightAlt:
		return "RightAlt"
	case KeyLeft:
		return "Left"
	case KeyRight:
		return "Right"
	case KeyUp:
		return "Up"
	case KeyDown:
		return "Down"
	case KeyEscape:
		return "Escape"
	case KeySpace:
		return "Space"
	case KeyBackspace:
		return "Backspace"
	case KeyTab:
		return "Tab"
	case KeyHome:
		return "Home"
	case KeyEnd:
		return "End"
	case KeyPageDown:
		return "PageDown"
	case KeyPageUp:
		return "PageUp"
	case KeyDelete:
		return "Delete"
	case KeyInsert:
		return "Insert"
	case KeyNumAdd:
		return "NumAdd"
	case KeyNumSubtract:
		return "NumSubtract"
	case KeyNumMultiply:
		return "NumMultiply"
	case KeyNumDivide:
		return "NumDivide"
	case KeyCapslock:
		return "Capslock"
	case KeyPrint:
		return "Print"
	case KeyPause:
		return "Pause"
	default:
		return "Unknown key " + strconv.Itoa(int(k))
	}
}
