package main

import (
	"fmt"

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
	cnvs          *canvas.Canvas
	ctx           *canvas.Context2D
	backgroundBox *dom.Element
	lightBox      *dom.Element
	mediumBox     *dom.Element
	darkBox       *dom.Element
	currentBox    *dom.Element
	frameText     *dom.Element
	color         string
	mouseDown     bool
)

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
	backgroundBox.Style.SetProperty("backgroundColor", background)
	backgroundBox.AddEventListener(dom.EvtClick, func(event *dom.Event) {
		fmt.Println("clicked the backgroundColor")
	})

	lightBox = doc.GetElementById("light-box")
	lightBox.Style.SetProperty("backgroundColor", light)
	lightBox.AddEventListener(dom.EvtClick, func(event *dom.Event) {
		fmt.Println("clicked the light color")
	})

	mediumBox = doc.GetElementById("medium-box")
	mediumBox.Style.SetProperty("backgroundColor", medium)
	mediumBox.AddEventListener(dom.EvtClick, func(event *dom.Event) {
		fmt.Println("clicked the medium color")
	})

	darkBox = doc.GetElementById("dark-box")
	darkBox.Style.SetProperty("backgroundColor", dark)
	darkBox.AddEventListener(dom.EvtClick, func(event *dom.Event) {
		fmt.Println("clicked the dark color")
	})

	currentBox = doc.GetElementById("current-box")
	color = medium

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
}
