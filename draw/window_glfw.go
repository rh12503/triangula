// +build glfw !windows,!sdl2

package draw

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	_ "image/png" // We allow loading PNGs by default.
	"io"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gonutz/gl/v2.1/gl"
	"github.com/gonutz/glfw/v3.1/glfw"
)

func init() {
	runtime.LockOSThread()
}

type window struct {
	running        bool
	pressed        []Key
	typed          []rune
	window         *glfw.Window
	width, height  float64
	textures       map[string]texture
	clicks         []MouseClick
	mouseX, mouseY int
	wheelX, wheelY float64
}

// RunWindow creates a new window and calls update 60 times per second.
func RunWindow(title string, width, height int, update UpdateFunction) error {
	if err := initSound(); err != nil {
		return err
	}
	defer closeSound()

	err := glfw.Init()
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 1)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	win, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return err
	}
	win.MakeContextCurrent()
	// center the window on the screen (omitting the window border)
	screen := glfw.GetMonitors()[0].GetVideoMode()
	win.SetPos((screen.Width-width)/2, (screen.Height-height)/2)

	err = gl.Init()
	if err != nil {
		return err
	}
	gl.MatrixMode(gl.PROJECTION)
	gl.Ortho(0, float64(width), float64(height), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	w := &window{
		running:  true,
		window:   win,
		width:    float64(width),
		height:   float64(height),
		textures: make(map[string]texture),
	}
	win.SetKeyCallback(w.keyPress)
	win.SetCharCallback(w.charTyped)
	win.SetMouseButtonCallback(w.mouseButtonEvent)
	win.SetCursorPosCallback(w.mousePositionChanged)
	win.SetScrollCallback(func(_ *glfw.Window, dx, dy float64) {
		w.wheelX += dx
		w.wheelY += dy
	})
	win.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		w.width, w.height = float64(width), float64(height)
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(0, w.width, w.height, 0, -1, 1)
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.MatrixMode(gl.MODELVIEW)
	})

	lastUpdateTime := time.Now().Add(-time.Hour)
	const updateInterval = 1.0 / 60.0
	for w.running && !win.ShouldClose() {
		glfw.PollEvents()

		now := time.Now()
		if now.Sub(lastUpdateTime).Seconds() > updateInterval {
			gl.ClearColor(0, 0, 0, 1)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			update(w)

			w.pressed = w.pressed[0:0]
			w.typed = w.typed[0:0]
			w.clicks = w.clicks[0:0]
			w.wheelX = 0
			w.wheelY = 0

			lastUpdateTime = now
			win.SwapBuffers()
		} else {
			time.Sleep(time.Millisecond)
		}
	}

	w.cleanUp()

	return nil
}

func (w *window) Close() {
	w.running = false
}

func (w *window) Size() (int, int) {
	return int(w.width + 0.5), int(w.height + 0.5)
}

func (w *window) SetFullscreen(f bool) {
	// TODO Find out how to toggle full screen in GLFW 3.1 and tell OpenGL about
	// it.
}

func (w *window) ShowCursor(show bool) {
	if show {
		w.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		w.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
}

func (w *window) keyPress(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		w.pressed = append(w.pressed, tokey(key))
	}
}

func (w *window) WasKeyPressed(key Key) bool {
	for _, pressed := range w.pressed {
		if pressed == key {
			return true
		}
	}
	return false
}

func (w *window) WasCharTyped(char rune) bool {
	for _, typed := range w.typed {
		if char == typed {
			return true
		}
	}
	return false
}

func (w *window) charTyped(_ *glfw.Window, char rune) {
	w.typed = append(w.typed, char)
}

func (w *window) IsKeyDown(key Key) bool {
	k := toGlfwKey(key)
	if k == glfw.KeyUnknown {
		return false
	}
	return w.window.GetKey(k) == glfw.Press
}

func (w *window) DrawPoint(x, y int, color Color) {
	gl.Begin(gl.POINTS)
	gl.Color4f(color.R, color.G, color.B, color.A)
	gl.Vertex2f(float32(x)+0.5, float32(y)+0.5)
	gl.End()
}

func (w *window) FillRect(x, y, width, height int, color Color) {
	if width <= 0 || height <= 0 {
		return
	}
	gl.Begin(gl.QUADS)
	gl.Color4f(color.R, color.G, color.B, color.A)
	gl.Vertex2i(int32(x), int32(y))
	gl.Vertex2i(int32(x+width), int32(y))
	gl.Vertex2i(int32(x+width), int32(y+height))
	gl.Vertex2i(int32(x), int32(y+height))
	gl.End()
}

func (w *window) DrawRect(x, y, width, height int, color Color) {
	if width <= 0 || height <= 0 {
		return
	}
	if width == 1 && height == 1 {
		w.DrawPoint(x, y, color)
		return
	}
	gl.Begin(gl.LINE_STRIP)
	gl.Color4f(color.R, color.G, color.B, color.A)
	gl.Vertex2f(float32(x)+0.5, float32(y)+0.5)
	gl.Vertex2f(float32(x+width)-0.5, float32(y)+0.5)
	gl.Vertex2f(float32(x+width)-0.5, float32(y+height)-0.5)
	gl.Vertex2f(float32(x)+0.5, float32(y+height)-0.5)
	gl.Vertex2f(float32(x)+0.5, float32(y)+0.5)
	gl.End()
}

func (w *window) DrawTriangle(x0, y0, x1, y1, x2, y2 int, color Color) {
	gl.Begin(gl.LINE_STRIP)
	gl.Color4f(color.R, color.G, color.B, color.A)
	gl.Vertex2f(float32(x0)+0.5, float32(y0)+0.5)
	gl.Vertex2f(float32(x1)-0.5, float32(y1)+0.5)
	gl.Vertex2f(float32(x2)-0.5, float32(y2)-0.5)
	gl.Vertex2f(float32(x0)+0.5, float32(y0)-0.5)
	gl.End()
}

func (w *window) FillTriangle(x0, y0, x1, y1, x2, y2 int, color Color) {
	gl.Begin(gl.TRIANGLES)
	gl.Color4f(color.R, color.G, color.B, color.A)
	gl.Vertex2i(int32(x0), int32(y0))
	gl.Vertex2i(int32(x1), int32(y1))
	gl.Vertex2i(int32(x2), int32(y2))
	gl.End()
}

func (w *window) DrawLine(x, y, x2, y2 int, color Color) {
	if x == x2 && y == y2 {
		w.DrawPoint(x, y, color)
		return
	}
	gl.Begin(gl.LINES)
	gl.Color4f(color.R, color.G, color.B, color.A)
	gl.Vertex2f(float32(x)+0.5, float32(y)+0.5)
	gl.Vertex2f(float32(x2+sign(x2-x))+0.5, float32(y2+sign(y2-y))+0.5)
	gl.End()
}

func sign(x int) int {
	if x == 0 {
		return 0
	}
	if x > 0 {
		return 1
	}
	return -1
}

type texture struct {
	id   uint32
	w, h int
}

func (w *window) loadTexture(r io.Reader, name string) (texture, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return texture{}, err
	}

	var rgba *image.RGBA
	if asRGBA, ok := img.(*image.RGBA); ok {
		rgba = asRGBA
	} else {
		rgba = image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			return texture{}, errors.New("unsupported stride")
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	}

	var tex uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Bounds().Dx()),
		int32(rgba.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	gl.Disable(gl.TEXTURE_2D)

	w.textures[name] = texture{
		id: tex,
		w:  rgba.Bounds().Dx(),
		h:  rgba.Bounds().Dy(),
	}

	return w.textures[name], nil
}

func (w *window) getOrLoadTexture(path string) (texture, error) {
	if tex, ok := w.textures[path]; ok {
		return tex, nil
	}

	imgFile, err := os.Open(path)
	if err != nil {
		return texture{}, err
	}
	defer imgFile.Close()

	return w.loadTexture(imgFile, path)
}

func (w *window) cleanUp() {
	for _, tex := range w.textures {
		gl.DeleteTextures(1, &tex.id)
	}
	w.textures = nil
}

func (w *window) Clicks() []MouseClick {
	return w.clicks
}

func (w *window) Characters() string {
	return string(w.typed)
}

func (w *window) mouseButtonEvent(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		b := toMouseButton(button)
		x, y := w.window.GetCursorPos()
		w.clicks = append(w.clicks, MouseClick{X: int(x), Y: int(y), Button: b})
	}
}

func (w *window) mousePositionChanged(_ *glfw.Window, x, y float64) {
	w.mouseX, w.mouseY = int(x+0.5), int(y+0.5)
}

func (w *window) MousePosition() (int, int) {
	return w.mouseX, w.mouseY
}

func (w *window) MouseWheelY() float64 {
	return w.wheelY
}

func (w *window) MouseWheelX() float64 {
	return w.wheelX
}

func toMouseButton(b glfw.MouseButton) MouseButton {
	if b == glfw.MouseButtonRight {
		return RightButton
	}
	if b == glfw.MouseButtonMiddle {
		return MiddleButton
	}
	return LeftButton
}

func (w *window) IsMouseDown(button MouseButton) bool {
	return w.window.GetMouseButton(toGlfwButton(button)) == glfw.Press
}

func toGlfwButton(b MouseButton) glfw.MouseButton {
	if b == RightButton {
		return glfw.MouseButtonRight
	}
	if b == MiddleButton {
		return glfw.MouseButtonMiddle
	}
	return glfw.MouseButtonLeft
}

func (w *window) DrawEllipse(x, y, width, height int, color Color) {
	outline := ellipseOutline(x, y, width, height)
	if len(outline) == 0 {
		return
	}
	gl.Begin(gl.POINTS)
	gl.Color4f(color.R, color.G, color.B, color.A)
	for _, p := range outline {
		gl.Vertex2f(float32(p.x)+0.5, float32(p.y)+0.5)
	}
	gl.End()
}

func (w *window) FillEllipse(x, y, width, height int, color Color) {
	area := ellipseArea(x, y, width, height)
	if len(area) == 0 {
		return
	}
	gl.Begin(gl.LINES)
	gl.Color4f(color.R, color.G, color.B, color.A)
	for i := 0; i < len(area); i += 2 {
		gl.Vertex2f(float32(area[i].x)+0.5, float32(area[i].y)+0.5)
		gl.Vertex2f(float32(area[i+1].x)+1.0, float32(area[i+1].y)+1.0)
	}
	gl.End()
}

func (w *window) DrawImageFile(path string, x, y int) error {
	tex, err := w.getOrLoadTexture(path)
	if err != nil {
		return err
	}

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.id)
	gl.Begin(gl.QUADS)

	gl.Color4f(1, 1, 1, 1)

	gl.TexCoord2i(0, 0)
	gl.Vertex2i(int32(x), int32(y))

	gl.TexCoord2i(1, 0)
	gl.Vertex2i(int32(x+tex.w), int32(y))

	gl.TexCoord2i(1, 1)
	gl.Vertex2i(int32(x+tex.w), int32(y+tex.h))

	gl.TexCoord2i(0, 1)
	gl.Vertex2i(int32(x), int32(y+tex.h))

	gl.End()
	gl.Disable(gl.TEXTURE_2D)

	return nil
}

func (w *window) DrawImageFileRotated(path string, x, y, degrees int) error {
	return w.DrawImageFileTo(path, x, y, -1, -1, degrees)
}

func (w *window) DrawImageFileTo(path string, x, y, width, height, degrees int) error {
	tex, err := w.getOrLoadTexture(path)
	if err != nil {
		return err
	}

	if width == -1 && height == -1 {
		width, height = tex.w, tex.h
	}

	x1, y1 := float32(x), float32(y)
	x2, y2 := float32(x+width-0), float32(y+height-0)
	cx, cy := x1+float32(width)/2, y1+float32(height)/2
	sin, cos := math.Sincos(float64(degrees) / 180 * math.Pi)
	sin32, cos32 := float32(sin), float32(cos)
	p := [4]pointf{
		{x1, y1},
		{x2, y1},
		{x2, y2},
		{x1, y2},
	}
	for i := range p {
		p[i].x, p[i].y = p[i].x-cx, p[i].y-cy
		p[i].x, p[i].y = cos32*p[i].x-sin32*p[i].y, sin32*p[i].x+cos32*p[i].y
		p[i].x, p[i].y = p[i].x+cx, p[i].y+cy
	}

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.id)
	gl.Begin(gl.QUADS)

	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2i(0, 0)
	gl.Vertex2f(p[0].x, p[0].y)

	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2i(1, 0)
	gl.Vertex2f(p[1].x, p[1].y)

	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2i(1, 1)
	gl.Vertex2f(p[2].x, p[2].y)

	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2i(0, 1)
	gl.Vertex2f(p[3].x, p[3].y)

	gl.End()
	gl.Disable(gl.TEXTURE_2D)

	return nil
}

func (w *window) DrawImageFilePart(
	path string,
	sourceX, sourceY, sourceWidth, sourceHeight int,
	destX, destY, destWidth, destHeight int,
	rotationCWDeg int,
) error {
	tex, err := w.getOrLoadTexture(path)
	if err != nil {
		return err
	}

	x1, y1 := float32(destX), float32(destY)
	x2, y2 := float32(destX+destWidth), float32(destY+destHeight)
	cx, cy := x1+float32(destWidth)/2, y1+float32(destHeight)/2
	sin, cos := math.Sincos(float64(rotationCWDeg) / 180 * math.Pi)
	sin32, cos32 := float32(sin), float32(cos)
	p := [4]pointf{
		{x1, y1},
		{x2, y1},
		{x2, y2},
		{x1, y2},
	}
	for i := range p {
		p[i].x, p[i].y = p[i].x-cx, p[i].y-cy
		p[i].x, p[i].y = cos32*p[i].x-sin32*p[i].y, sin32*p[i].x+cos32*p[i].y
		p[i].x, p[i].y = p[i].x+cx, p[i].y+cy
	}

	u0 := float32(sourceX) / float32(tex.w)
	u1 := float32(sourceX+sourceWidth) / float32(tex.w)
	v0 := float32(sourceY) / float32(tex.h)
	v1 := float32(sourceY+sourceHeight) / float32(tex.h)

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.id)
	gl.Begin(gl.QUADS)

	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2f(u0, v0)
	gl.Vertex2f(p[0].x, p[0].y)

	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2f(u1, v0)
	gl.Vertex2f(p[1].x, p[1].y)

	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2f(u1, v1)
	gl.Vertex2f(p[2].x, p[2].y)

	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2f(u0, v1)
	gl.Vertex2f(p[3].x, p[3].y)

	gl.End()
	gl.Disable(gl.TEXTURE_2D)

	return nil

	return nil
}

type pointf struct{ x, y float32 }

func (w *window) GetTextSize(text string) (width, height int) {
	return w.GetScaledTextSize(text, 1.0)
}

func (w *window) GetScaledTextSize(text string, scale float32) (width, height int) {
	fontTexture, ok := w.textures[fontTextureID]
	if !ok {
		return 0, 0
	}
	width = int(float32(fontTexture.w/16)*scale + 0.5)
	height = int(float32(fontTexture.h/16)*scale + 0.5)
	lines := strings.Split(text, "\n")
	maxLineW := 0
	for _, line := range lines {
		w := utf8.RuneCountInString(line)
		if w > maxLineW {
			maxLineW = w
		}
	}
	return width * maxLineW, height * len(lines)
}

func (w *window) DrawText(text string, x, y int, color Color) {
	w.DrawScaledText(text, x, y, 1.0, color)
}

const fontTextureID = "///font_texture"

func (w *window) DrawScaledText(text string, x, y int, scale float32, color Color) {
	fontTexture, ok := w.textures[fontTextureID]
	if !ok {
		var err error
		fontTexture, err = w.loadTexture(bytes.NewReader(bitmapFontWhitePng[:]), fontTextureID)
		if err != nil {
			panic(err)
		}
	}

	width, height := int32(fontTexture.w/16), int32(fontTexture.h/16)
	width = int32(float32(width)*scale + 0.5)
	height = int32(float32(height)*scale + 0.5)

	var srcX, srcY float32
	destX, destY := int32(x), int32(y)

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, fontTexture.id)

	gl.Begin(gl.QUADS)
	for _, r := range text {
		if r == '\n' {
			destX = int32(x)
			destY += height
			continue
		}
		r = runeToFont(r)

		srcX = float32(r%16) / 16
		srcY = float32(r/16) / 16

		gl.Color4f(color.R, color.G, color.B, color.A)
		gl.TexCoord2f(srcX, srcY)
		gl.Vertex2i(destX, destY)

		gl.Color4f(color.R, color.G, color.B, color.A)
		gl.TexCoord2f(srcX+1.0/16, srcY)
		gl.Vertex2i(destX+width, destY)

		gl.Color4f(color.R, color.G, color.B, color.A)
		gl.TexCoord2f(srcX+1.0/16, srcY+1.0/16)
		gl.Vertex2i(destX+width, destY+height)

		gl.Color4f(color.R, color.G, color.B, color.A)
		gl.TexCoord2f(srcX, srcY+1.0/16)
		gl.Vertex2i(destX, destY+height)

		destX += width
	}
	gl.End()
	gl.Disable(gl.TEXTURE_2D)
}

func (w *window) PlaySoundFile(path string) error {
	return playSoundFile(path)
}

func toGlfwKey(key Key) glfw.Key {
	switch key {
	case KeyA:
		return glfw.KeyA
	case KeyB:
		return glfw.KeyB
	case KeyC:
		return glfw.KeyC
	case KeyD:
		return glfw.KeyD
	case KeyE:
		return glfw.KeyE
	case KeyF:
		return glfw.KeyF
	case KeyG:
		return glfw.KeyG
	case KeyH:
		return glfw.KeyH
	case KeyI:
		return glfw.KeyI
	case KeyJ:
		return glfw.KeyJ
	case KeyK:
		return glfw.KeyK
	case KeyL:
		return glfw.KeyL
	case KeyM:
		return glfw.KeyM
	case KeyN:
		return glfw.KeyN
	case KeyO:
		return glfw.KeyO
	case KeyP:
		return glfw.KeyP
	case KeyQ:
		return glfw.KeyQ
	case KeyR:
		return glfw.KeyR
	case KeyS:
		return glfw.KeyS
	case KeyT:
		return glfw.KeyT
	case KeyU:
		return glfw.KeyU
	case KeyV:
		return glfw.KeyV
	case KeyW:
		return glfw.KeyW
	case KeyX:
		return glfw.KeyX
	case KeyY:
		return glfw.KeyY
	case KeyZ:
		return glfw.KeyZ
	case Key0:
		return glfw.Key0
	case Key1:
		return glfw.Key1
	case Key2:
		return glfw.Key2
	case Key3:
		return glfw.Key3
	case Key4:
		return glfw.Key4
	case Key5:
		return glfw.Key5
	case Key6:
		return glfw.Key6
	case Key7:
		return glfw.Key7
	case Key8:
		return glfw.Key8
	case Key9:
		return glfw.Key9
	case KeyNum0:
		return glfw.KeyKP0
	case KeyNum1:
		return glfw.KeyKP1
	case KeyNum2:
		return glfw.KeyKP2
	case KeyNum3:
		return glfw.KeyKP3
	case KeyNum4:
		return glfw.KeyKP4
	case KeyNum5:
		return glfw.KeyKP5
	case KeyNum6:
		return glfw.KeyKP6
	case KeyNum7:
		return glfw.KeyKP7
	case KeyNum8:
		return glfw.KeyKP8
	case KeyNum9:
		return glfw.KeyKP9
	case KeyF1:
		return glfw.KeyF1
	case KeyF2:
		return glfw.KeyF2
	case KeyF3:
		return glfw.KeyF3
	case KeyF4:
		return glfw.KeyF4
	case KeyF5:
		return glfw.KeyF5
	case KeyF6:
		return glfw.KeyF6
	case KeyF7:
		return glfw.KeyF7
	case KeyF8:
		return glfw.KeyF8
	case KeyF9:
		return glfw.KeyF9
	case KeyF10:
		return glfw.KeyF10
	case KeyF11:
		return glfw.KeyF11
	case KeyF12:
		return glfw.KeyF12
	case KeyF13:
		return glfw.KeyF13
	case KeyF14:
		return glfw.KeyF14
	case KeyF15:
		return glfw.KeyF15
	case KeyF16:
		return glfw.KeyF16
	case KeyF17:
		return glfw.KeyF17
	case KeyF18:
		return glfw.KeyF18
	case KeyF19:
		return glfw.KeyF19
	case KeyF20:
		return glfw.KeyF20
	case KeyF21:
		return glfw.KeyF21
	case KeyF22:
		return glfw.KeyF22
	case KeyF23:
		return glfw.KeyF23
	case KeyF24:
		return glfw.KeyF24
	case KeyEnter:
		return glfw.KeyEnter
	case KeyNumEnter:
		return glfw.KeyKPEnter
	case KeyLeftControl:
		return glfw.KeyLeftControl
	case KeyRightControl:
		return glfw.KeyRightControl
	case KeyLeftShift:
		return glfw.KeyLeftShift
	case KeyRightShift:
		return glfw.KeyRightShift
	case KeyLeftAlt:
		return glfw.KeyLeftAlt
	case KeyRightAlt:
		return glfw.KeyRightAlt
	case KeyLeft:
		return glfw.KeyLeft
	case KeyRight:
		return glfw.KeyRight
	case KeyUp:
		return glfw.KeyUp
	case KeyDown:
		return glfw.KeyDown
	case KeyEscape:
		return glfw.KeyEscape
	case KeySpace:
		return glfw.KeySpace
	case KeyBackspace:
		return glfw.KeyBackspace
	case KeyTab:
		return glfw.KeyTab
	case KeyHome:
		return glfw.KeyHome
	case KeyEnd:
		return glfw.KeyEnd
	case KeyPageDown:
		return glfw.KeyPageDown
	case KeyPageUp:
		return glfw.KeyPageUp
	case KeyDelete:
		return glfw.KeyDelete
	case KeyInsert:
		return glfw.KeyInsert
	case KeyNumAdd:
		return glfw.KeyKPAdd
	case KeyNumSubtract:
		return glfw.KeyKPSubtract
	case KeyNumMultiply:
		return glfw.KeyKPMultiply
	case KeyNumDivide:
		return glfw.KeyKPDivide
	case KeyCapslock:
		return glfw.KeyCapsLock
	case KeyPrint:
		return glfw.KeyPrintScreen
	case KeyPause:
		return glfw.KeyPause
	}

	return glfw.KeyUnknown
}

func tokey(k glfw.Key) Key {
	switch k {
	case glfw.KeyA:
		return KeyA
	case glfw.KeyB:
		return KeyB
	case glfw.KeyC:
		return KeyC
	case glfw.KeyD:
		return KeyD
	case glfw.KeyE:
		return KeyE
	case glfw.KeyF:
		return KeyF
	case glfw.KeyG:
		return KeyG
	case glfw.KeyH:
		return KeyH
	case glfw.KeyI:
		return KeyI
	case glfw.KeyJ:
		return KeyJ
	case glfw.KeyK:
		return KeyK
	case glfw.KeyL:
		return KeyL
	case glfw.KeyM:
		return KeyM
	case glfw.KeyN:
		return KeyN
	case glfw.KeyO:
		return KeyO
	case glfw.KeyP:
		return KeyP
	case glfw.KeyQ:
		return KeyQ
	case glfw.KeyR:
		return KeyR
	case glfw.KeyS:
		return KeyS
	case glfw.KeyT:
		return KeyT
	case glfw.KeyU:
		return KeyU
	case glfw.KeyV:
		return KeyV
	case glfw.KeyW:
		return KeyW
	case glfw.KeyX:
		return KeyX
	case glfw.KeyY:
		return KeyY
	case glfw.KeyZ:
		return KeyZ
	case glfw.Key0:
		return Key0
	case glfw.Key1:
		return Key1
	case glfw.Key2:
		return Key2
	case glfw.Key3:
		return Key3
	case glfw.Key4:
		return Key4
	case glfw.Key5:
		return Key5
	case glfw.Key6:
		return Key6
	case glfw.Key7:
		return Key7
	case glfw.Key8:
		return Key8
	case glfw.Key9:
		return Key9
	case glfw.KeyKP0:
		return KeyNum0
	case glfw.KeyKP1:
		return KeyNum1
	case glfw.KeyKP2:
		return KeyNum2
	case glfw.KeyKP3:
		return KeyNum3
	case glfw.KeyKP4:
		return KeyNum4
	case glfw.KeyKP5:
		return KeyNum5
	case glfw.KeyKP6:
		return KeyNum6
	case glfw.KeyKP7:
		return KeyNum7
	case glfw.KeyKP8:
		return KeyNum8
	case glfw.KeyKP9:
		return KeyNum9
	case glfw.KeyF1:
		return KeyF1
	case glfw.KeyF2:
		return KeyF2
	case glfw.KeyF3:
		return KeyF3
	case glfw.KeyF4:
		return KeyF4
	case glfw.KeyF5:
		return KeyF5
	case glfw.KeyF6:
		return KeyF6
	case glfw.KeyF7:
		return KeyF7
	case glfw.KeyF8:
		return KeyF8
	case glfw.KeyF9:
		return KeyF9
	case glfw.KeyF10:
		return KeyF10
	case glfw.KeyF11:
		return KeyF11
	case glfw.KeyF12:
		return KeyF12
	case glfw.KeyF13:
		return KeyF13
	case glfw.KeyF14:
		return KeyF14
	case glfw.KeyF15:
		return KeyF15
	case glfw.KeyF16:
		return KeyF16
	case glfw.KeyF17:
		return KeyF17
	case glfw.KeyF18:
		return KeyF18
	case glfw.KeyF19:
		return KeyF19
	case glfw.KeyF20:
		return KeyF20
	case glfw.KeyF21:
		return KeyF21
	case glfw.KeyF22:
		return KeyF22
	case glfw.KeyF23:
		return KeyF23
	case glfw.KeyF24:
		return KeyF24
	case glfw.KeyEnter:
		return KeyEnter
	case glfw.KeyKPEnter:
		return KeyNumEnter
	case glfw.KeyLeftControl:
		return KeyLeftControl
	case glfw.KeyRightControl:
		return KeyRightControl
	case glfw.KeyLeftShift:
		return KeyLeftShift
	case glfw.KeyRightShift:
		return KeyRightShift
	case glfw.KeyLeftAlt:
		return KeyLeftAlt
	case glfw.KeyRightAlt:
		return KeyRightAlt
	case glfw.KeyLeft:
		return KeyLeft
	case glfw.KeyRight:
		return KeyRight
	case glfw.KeyUp:
		return KeyUp
	case glfw.KeyDown:
		return KeyDown
	case glfw.KeyEscape:
		return KeyEscape
	case glfw.KeySpace:
		return KeySpace
	case glfw.KeyBackspace:
		return KeyBackspace
	case glfw.KeyTab:
		return KeyTab
	case glfw.KeyHome:
		return KeyHome
	case glfw.KeyEnd:
		return KeyEnd
	case glfw.KeyPageDown:
		return KeyPageDown
	case glfw.KeyPageUp:
		return KeyPageUp
	case glfw.KeyDelete:
		return KeyDelete
	case glfw.KeyInsert:
		return KeyInsert
	case glfw.KeyKPAdd:
		return KeyNumAdd
	case glfw.KeyKPSubtract:
		return KeyNumSubtract
	case glfw.KeyKPMultiply:
		return KeyNumMultiply
	case glfw.KeyKPDivide:
		return KeyNumDivide
	case glfw.KeyCapsLock:
		return KeyCapslock
	case glfw.KeyPrintScreen:
		return KeyPrint
	case glfw.KeyPause:
		return KeyPause
	}

	return Key(0)
}
