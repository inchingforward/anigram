package main

import (
	"fmt"
	"math"

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
	color          string
	mouseDown      bool
	currFrameIndex int
	anim           animation
)

type serverData struct {
	title string
	data  string
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

func prevFrame(event *dom.Event) {
	fmt.Println("prevFrame clicked")
}

func nextFrame(event *dom.Event) {
	fmt.Println("nextFrame clicked")
}

func clear(event *dom.Event) {
	fmt.Println("clear clicked")
}

func delete(event *dom.Event) {
	fmt.Println("delete clicked")
}

func play(event *dom.Event) {
	fmt.Println("play clicked")
}

func save(event *dom.Event) {
	fmt.Println("save clicked")
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

	// var boundingRect = canvas.getBoundingClientRect();
	// var x = e.clientX - boundingRect.left;
	// var y = e.clientY - boundingRect.top;

	// var row = Math.floor(y / PIXEL_SIZE);
	// var col = Math.floor(x / PIXEL_SIZE);

	// fillSquareAt(row, col, color, true);
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
	return animation{sd.title, stringToFrames(sd.data)}
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

func fillSquareAt(col int, row int, color string, updateAnimation bool) {
	ctx.BeginPath()
	ctx.FillStyle = color
	ctx.Rect(float64(row)*pixelSize, float64(col)*pixelSize, pixelSize, pixelSize)
	ctx.Fill()

	if updateAnimation {
		updateFrameSquare(row, col, color)
	}
}

func updateFrameSquare(row int, col int, colorStr string) {
	frame := anim.frames[currFrameIndex]
	frameIndex := coordToIndex(coord{row, col})
	color := colorToStr(color)
	frame = anim.frames[currFrameIndex]
	anim.frames[currFrameIndex] = string(frame[0:frameIndex]) + color + string(frame[frameIndex+1])
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

func main() {
	window := dom.Window()
	doc := dom.Document()
	cnvs := canvas.New(doc.GetElementById("canvas-grid").Object)
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
	setColor(medium)

	doc.GetElementById("prevFrame").AddEventListener(dom.EvtClick, prevFrame)
	doc.GetElementById("nextFrame").AddEventListener(dom.EvtClick, nextFrame)
	doc.GetElementById("clear").AddEventListener(dom.EvtClick, clear)
	doc.GetElementById("delete").AddEventListener(dom.EvtClick, delete)
	doc.GetElementById("play").AddEventListener(dom.EvtClick, play)
	doc.GetElementById("save").AddEventListener(dom.EvtClick, save)

	cnvs.AddEventListener(dom.EvtMousemove, canvasMouseMove)
	cnvs.AddEventListener(dom.EvtMousedown, canvasMouseDown)
	window.AddEventListener(dom.EvtMouseup, mouseUp)
	doc.AddEventListener(dom.EvtKeyup, keyUp)

	frameText = doc.GetElementById("frameText")

	initialData := serverData{
		title: "Testing",
		data:  "000000000000000000000000000000000010010000000000002002000000000000200200000000000020020100000000002002020000000000200200000000000021120100000000002222020000000000200202000000000020020200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000222220000000000222222200000002222000022000000020000002220000000000000222000000000000222200000000022222000000000222222000000000220000000000000022000000000000022200002222000002222222222000000022222222000000000000000000000000000000000000000000000000000000000022200000000000022222200000000222000220000000022200002200000002200000220000000000002222000000000002222000000000000000200000000000000022000000022200022200000002220002220000000222222222000000000022222000000000000000000000000000000000000000000000000000000000000022000000000000022200000000000022200000000000022220000000000022202000000000002200200000000002220020000000000222222222200000002222222220000000000020000000000000002000000000000002220000000000000222000000000000022200000000000000000000000000000000000000000022222222200000002222222220000002220000022000000222000000000000022222222200000000222222222000000022220002220000000000000222000022000000022200002220000002220000222000000220000022222002222000000222222222000000000000000000000000000000000000000000000000000000000002222200000000022222222000000022222220000000022000000000000002200000000000000220022222000000022022222220000002222000022200000222000000220000002200000222000000222002222200000002222222200000000222222200000000002222000000000000000000000000000000000000000000000000000000000022222222220000022222222222000002222200022000000022000002200000000000002200000000000002220000000000000220000000000000022000000000000022000000000000022200000000000002220000000000002222000000000000022200000000000000000000000000000000000000000000222200000000000222222000000000022222220000000002200022000000000220002200000000022222200000000000222220000000000222202200000000220000222000000022000022200000002200002220000000222222220000000022222220000000000000000000000000000000000000000000000000000000000022222200000000002222220000000002222222200000002220000220000000222000022000000022200002200000000220000220000000022200222000000002222222200000000022222020000002200000022000000222000022000000002222222200000000022222000000000000000000000000000000000000000000000000000000000002000022220000002200022222200002220002222220000022002200002200002200200000220000220020000022000022002000002200002200200000220000220022000220000022002222222000022220222222000002222002222000000222000000000000000000000000000"}

	currFrameIndex = 0
	anim = initialData.toAnimation()

	loadFrame(anim.frames[0])
	updateFrameText()
}
