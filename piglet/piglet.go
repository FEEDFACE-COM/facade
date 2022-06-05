//go:build (linux && arm) || DARWIN_GUI
// +build linux,arm DARWIN_GUI

package piglet

import "C"
import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"fmt"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	glfw "github.com/go-gl/glfw/v3.3/glfw"

	"unsafe"
)

/*
	go get -u -tags=gles3 github.com/go-gl/glfw/v3.3/glfw
*/

var monitor *glfw.Monitor
var window *glfw.Window
var vidmode *glfw.VidMode
var winpos *gfx.Frame

const WINDOW_WIDTH, WINDOW_HEIGHT = 864, 540
const MONITOR_WIDTH, MONITOR_HEIGHT = 1728, 1117

// Create a new EGL rendering context
func CreateContext() error {
	log.Debug("%s init", Desc())

	var err error = glfw.Init()
	if err != nil {
		return log.NewError("fail to initialize renderer: %s", err)
	}
	monitor = glfw.GetPrimaryMonitor()
	vidmode = monitor.GetVideoMode()
	log.Debug("%s mode %dx%d @%d fps", Desc(), vidmode.Width, vidmode.Height, vidmode.RefreshRate)
	{
		for _, mode := range monitor.GetVideoModes() {
			w, h, fps := mode.Width, mode.Height, mode.RefreshRate
			log.Debug("%s mode %dx%d @%d fps", Desc(), w, h, fps)
		}
	}

	window, err = glfw.CreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "FACADE by FEEDFACE.COM", nil, nil)
	if err != nil {
		glfw.Terminate()
		return log.NewError("fail to glfw create window: %s", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	window.SetAspectRatio(WINDOW_WIDTH, WINDOW_HEIGHT)
	window.SetSizeLimits(WINDOW_WIDTH/2., WINDOW_WIDTH/2., gl.DONT_CARE, gl.DONT_CARE)
	//window.SetSizeCallback(func(win *glfw.Window, w int, h int) { SizeFun(w, h) })
	//window.SetFramebufferSizeCallback(func(win *glfw.Window, w int, h int) { renderer.FramebufferSizeFun(w, h) })
	window.SetKeyCallback(func(win *glfw.Window, k glfw.Key, c int, a glfw.Action, m glfw.ModifierKey) { KeyFun(k, a, m) })
	//window.SetRefreshCallback( func(win *glfw.Window) { renderer.RefreshFun() } )

	window.MakeContextCurrent()
	return nil
}

func Loop() bool {
	return !window.ShouldClose()
}

func Desc() string {
	w, h := 0, 0
	if window != nil {
		w, h = window.GetFramebufferSize()
	}
	return fmt.Sprintf("piglet[%dx%d]", w, h)
}

// Attach the EGL rendering context to the EGL surface
func MakeCurrent() error {
	window.MakeContextCurrent()
	return nil
}

// Post the EGL surface color buffer to the native display
func SwapBuffers() error {
	window.SwapBuffers()
	glfw.PollEvents()
	return nil
}

// Return the size of the native display, in pixels
func GetDisplaySize() (int32, int32) {
	w, h := window.GetFramebufferSize()
	return int32(w), int32(h)
}

// Return a GL or an EGL extension function
func GetProcAddress(name string) unsafe.Pointer {
	return glfw.GetProcAddress(name)
}

// Destroy an EGL rendering context
func DestroyContext() error {
	log.Debug("%s terminate", Desc())
	glfw.Terminate()
	return nil
}

func ErrorString(e uint32) string {
	return glfw.ErrorCode(e).String()
}

func ToggleFullScreen() {
	if winpos == nil {
		x, y := window.GetPos()
		w, h := window.GetSize()
		fps := vidmode.RefreshRate
		frame := gfx.Frame{P: gfx.Point{X: float32(x), Y: float32(y)}, S: gfx.Size{W: float32(w), H: float32(h)}}
		winpos = &frame
		log.Info("%s fullscreen %dx%d @%d fps", Desc(), MONITOR_WIDTH, MONITOR_HEIGHT, fps)
		window.SetMonitor(monitor, 0, 0, MONITOR_WIDTH, MONITOR_HEIGHT, fps)
	} else {
		x, y := int(winpos.P.X), int(winpos.P.Y)
		w, h := int(winpos.S.W), int(winpos.S.H)
		fps := vidmode.RefreshRate
		winpos = nil
		log.Info("%s window %dx%d @%d fps", Desc(), w, h, fps)
		window.SetMonitor(nil, x, y, w, h, fps)
	}
}

func KeyFun(key glfw.Key, action glfw.Action, mod glfw.ModifierKey) {

	if mod == 0x2 && action == 0x1 && key == 0x43 {
		log.Notice("%s key ctrl-c", Desc())
		window.SetShouldClose(true)
		return
	}

	if key == 0x20 && action == 0x1 {
		log.Notice("%s key space", Desc())
		ToggleFullScreen()
		return
	}

	if key == 0x100 && action == 0x1 {
		log.Notice("%s key escape", Desc())
		window.SetShouldClose(true)
		return
	}

	log.Debug("%s key 0x%02x action %d mod %x", Desc(), key, action, mod)

}
