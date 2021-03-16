package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	canvas "github.com/oskca/gopherjs-canvas"
	dom "github.com/oskca/gopherjs-dom"
)

const (
	dark       = "#000000"
	medium     = "#999999"
	background = "#d2d2d2"
	light      = "#ffffff"
	pixelSize  = 16
	side       = pixelSize * pixelSize
)

var (
	cnvs           *canvas.Canvas
	ctx            *canvas.Context2D
	backgroundBox  *dom.Element
	lightBox       *dom.Element
	mediumBox      *dom.Element
	darkBox        *dom.Element
	currentBox     *dom.Element
	frameText      *dom.Element
	timeInput      *dom.Element
	loopCbx        *dom.Element
	color          string
	mouseDown      bool
	currFrameIndex int
	anim           animation
	emptyFrame     string
	copyBuffer     string
)

type serverData struct {
	Title string `json:"title"`
	Data  string `json:"data"`
}

type animation struct {
	title  string
	frames []string
}

type coord struct {
	x int
	y int
}

func setColor(nextColor string) {
	color = nextColor
	currentBox.Style.SetProperty("background-color", color)
}

func firstFrame(event *dom.Event) {
	currFrameIndex = 0
	loadFrame(anim.frames[currFrameIndex])
	updateFrameText()
}

func prevFrame(event *dom.Event) {
	nextFrameIndex := currFrameIndex - 1
	if nextFrameIndex < 0 {
		return
	}

	currFrameIndex = nextFrameIndex

	loadFrame(anim.frames[currFrameIndex])
	updateFrameText()
}

func nextFrame(event *dom.Event) {
	if currFrameIndex >= len(anim.frames)-1 {
		anim.frames = append(anim.frames, emptyFrame)
	}

	currFrameIndex++

	loadFrame(anim.frames[currFrameIndex])
	updateFrameText()
}

func lastFrame(event *dom.Event) {
	currFrameIndex = len(anim.frames) - 1
	loadFrame(anim.frames[currFrameIndex])
	updateFrameText()
}

func clear(event *dom.Event) {
	ctx.BeginPath()
	ctx.FillStyle = background
	ctx.Rect(0, 0, side, side)
	ctx.Fill()

	anim.frames[currFrameIndex] = emptyFrame
}

func copy(event *dom.Event) {
	copyBuffer = anim.frames[currFrameIndex]
}

func paste(event *dom.Event) {
	anim.frames[currFrameIndex] = copyBuffer
	loadFrame(anim.frames[currFrameIndex])
}

func delete(event *dom.Event) {
	if currFrameIndex == 0 {
		anim.frames = anim.frames[1:]
	} else if currFrameIndex == len(anim.frames)-1 {
		anim.frames = anim.frames[0:currFrameIndex]
		currFrameIndex--
	} else {
		anim.frames = append(anim.frames[0:currFrameIndex], anim.frames[currFrameIndex+1:]...)
	}

	loadFrame(anim.frames[currFrameIndex])
	updateFrameText()
}

func animate() {
	loadFrame(anim.frames[currFrameIndex])
	updateFrameText()

	if currFrameIndex < (len(anim.frames) - 1) {
		currFrameIndex++
		t := timeInput.Get("value").Int()
		time.AfterFunc(time.Duration(t)*time.Millisecond, animate)
	} else {
		loop := loopCbx.Get("checked").Bool()
		if loop {
			currFrameIndex = 0
			t := timeInput.Get("value").Int()
			time.AfterFunc(time.Duration(t)*time.Millisecond, animate)
		}
	}
}

func play(event *dom.Event) {
	currFrameIndex = 0
	animate()
}

func save(event *dom.Event) {
	animStr := strings.Join(anim.frames, "")
	fmt.Println(animStr)
}

func canvasMouseDown(event *dom.Event) {
	mouseDown = true
	canvasMouseOver(event)
}

func canvasMouseMove(event *dom.Event) {
	canvasMouseOver(event)
}

func canvasMouseOver(event *dom.Event) {
	if !mouseDown {
		return
	}

	handlePaintMovement(event.ClientX, event.ClientY)
}

func touchStart(event *dom.Event) {
	mouseDown = true
}

func touchEnd(event *dom.Event) {
	mouseDown = false
}

func touchMove(event *dom.Event) {
	event.PreventDefault()
	touch := event.Get("touches").Index(0)

	x := touch.Get("clientX").Int()
	y := touch.Get("clientY").Int()

	handlePaintMovement(x, y)
}

func handlePaintMovement(x, y int) {
	x = x - cnvs.Call("getBoundingClientRect").Get("left").Int()
	y = y - cnvs.Call("getBoundingClientRect").Get("top").Int()

	row := int(math.Floor(float64(y) / pixelSize))
	col := int(math.Floor(float64(x) / pixelSize))

	fillSquareAt(row, col, color, true)
}

func keyUp(event *dom.Event) {
	if event.Key == "ArrowLeft" {
		prevFrame(event)
	} else if event.Key == "ArrowRight" {
		nextFrame(event)
	} else if event.Key == "KeyC" {
		clear(event)
	} else if event.Key == "KeyP" {
		play(event)
	}
}

func mouseUp(event *dom.Event) {
	mouseDown = false
}

func stringToFrames(data string) []string {
	if data == "" {
		return make([]string, 0)
	}

	var s []string
	for i := 0; i < len(data); i += side {
		if i+side > len(data) {
			s = append(s, data[i:])
		} else {
			s = append(s, data[i:i+side])
		}
	}

	return s
}

func (sd serverData) toAnimation() animation {
	return animation{sd.Title, stringToFrames(sd.Data)}
}

func strToColor(char string) string {
	if char == "0" {
		return background
	} else if char == "1" {
		return light
	} else if char == "2" {
		return medium
	} else if char == "3" {
		return dark
	} else {
		return background
	}
}

func colorToStr(color string) string {
	if color == background {
		return "0"
	} else if color == light {
		return "1"
	} else if color == medium {
		return "2"
	} else if color == dark {
		return "3"
	} else {
		return "0"
	}
}

func indexToCoord(index int) coord {
	return coord{index % pixelSize, int(math.Floor(float64(index) / pixelSize))}
}

func coordToIndex(c coord) int {
	return c.y*pixelSize + c.x
}

func fillSquareAt(row int, col int, color string, updateAnimation bool) {
	ctx.BeginPath()
	ctx.FillStyle = color
	ctx.Rect(float64(col)*pixelSize, float64(row)*pixelSize, pixelSize, pixelSize)
	ctx.Fill()

	if updateAnimation {
		updateFrameSquare(row, col, color)
	}
}

func updateFrameSquare(row int, col int, colorStr string) {
	frameIndex := coordToIndex(coord{col, row})
	color := colorToStr(color)
	frame := anim.frames[currFrameIndex]
	newFrame := string(frame[0:frameIndex]) + color + string(frame[frameIndex+1:])
	anim.frames[currFrameIndex] = newFrame
}

func loadFrame(frame string) {
	for i := 0; i < len(frame); i++ {
		ch := string(frame[i])
		color = strToColor(ch)
		coord := indexToCoord(i)
		fillSquareAt(coord.y, coord.x, color, false)
	}
}

func updateFrameText() {
	frameText.InnerHTML = fmt.Sprintf("Frame %d of %d", (currFrameIndex + 1), len(anim.frames))
}

func getAnimation() (animation, error) {
	resp, err := http.Get("/api/animations/1")
	if err != nil {
		fmt.Println(err)
		return animation{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return animation{}, err
	}

	var initialData serverData
	err = json.Unmarshal(body, &initialData)
	if err != nil {
		return animation{}, err
	}

	anim = initialData.toAnimation()

	return anim, nil
}

func main() {
	window := dom.Window()
	doc := dom.Document()
	cnvs = canvas.New(doc.GetElementById("canvas-grid").Object)
	ctx = cnvs.GetContext2D()

	cnvs.Width = side
	cnvs.Height = side

	ctx.FillStyle = background
	ctx.BeginPath()
	ctx.Rect(0, 0, side, side)
	ctx.Fill()

	backgroundBox = doc.GetElementById("background-box")
	backgroundBox.Style.SetProperty("background-color", background)
	backgroundBox.AddEventListener(dom.EvtClick, func(event *dom.Event) {
		setColor(background)
	})

	lightBox = doc.GetElementById("light-box")
	lightBox.Style.SetProperty("background-color", light)
	lightBox.AddEventListener(dom.EvtClick, func(event *dom.Event) {
		setColor(light)
	})

	mediumBox = doc.GetElementById("medium-box")
	mediumBox.Style.SetProperty("background-color", medium)
	mediumBox.AddEventListener(dom.EvtClick, func(event *dom.Event) {
		setColor(medium)
	})

	darkBox = doc.GetElementById("dark-box")
	darkBox.Style.SetProperty("background-color", dark)
	darkBox.AddEventListener(dom.EvtClick, func(event *dom.Event) {
		setColor(dark)
	})

	currentBox = doc.GetElementById("current-box")
	setColor(background)

	doc.GetElementById("firstFrame").AddEventListener(dom.EvtClick, firstFrame)
	doc.GetElementById("prevFrame").AddEventListener(dom.EvtClick, prevFrame)
	doc.GetElementById("nextFrame").AddEventListener(dom.EvtClick, nextFrame)
	doc.GetElementById("lastFrame").AddEventListener(dom.EvtClick, lastFrame)
	doc.GetElementById("clear").AddEventListener(dom.EvtClick, clear)
	doc.GetElementById("copy").AddEventListener(dom.EvtClick, copy)
	doc.GetElementById("paste").AddEventListener(dom.EvtClick, paste)
	doc.GetElementById("delete").AddEventListener(dom.EvtClick, delete)
	doc.GetElementById("play").AddEventListener(dom.EvtClick, play)
	doc.GetElementById("save").AddEventListener(dom.EvtClick, save)

	cnvs.AddEventListener(dom.EvtMousemove, canvasMouseMove)
	cnvs.AddEventListener(dom.EvtMousedown, canvasMouseDown)
	cnvs.AddEventListener("touchstart", touchStart)
	cnvs.AddEventListener("touchend", touchEnd)
	cnvs.AddEventListener("touchmove", touchMove)

	window.AddEventListener(dom.EvtMouseup, mouseUp)
	doc.AddEventListener(dom.EvtKeyup, keyUp)

	emptyFrame = strings.Repeat("0", side)
	copyBuffer = emptyFrame

	frameText = doc.GetElementById("frameText")
	timeInput = doc.GetElementById("time")
	loopCbx = doc.GetElementById("loop")

	currFrameIndex = 0
	anim, err := getAnimation()
	if err != nil {
		fmt.Printf("Error getting animation: %v", err)
		return
	}

	loadFrame(anim.frames[0])
	updateFrameText()
}
