// SPDX-License-Identifier: Unlicense OR MIT

package main

/*
Small Gio programm to test gamma correction.

The program draws two horizontal sequences of boxes colored from
black to white by using two different methods.

The first line of gray boxes test the effect of gamma correction
with anti-aliasing. Each box is drawn as a stack of black horizontal lines.
When the line thickness covers the full line height, is should be shown as
all black. When the line thickness is half the pixel height, it should be
drawn as mid gray, etc.

The second line, draws the boxes as simple rectangles filled with a specified
color. The color is gray and ranges from black to white.

With correct anti-aliasing and gamma correction, the gray scale must look
visually as uniformly increasing gray values.

The second line should look identical to the first line. If not, there is a
problem with gamma correction or anti-aliasing.
*/

import (
	"image/color"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	const nbrBox = 32.
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			width := float64(e.Size.X)
			boxWidth := width / nbrBox
			boxHeight := math.Ceil((boxWidth) * 2 * 1.61803398875) // make it nice by using the golden ratio
			drawGrayBar1(gtx, width, nbrBox, boxWidth, boxHeight, 0, 0)
			drawGrayBar2(gtx, width, nbrBox, boxWidth, boxHeight, 0, boxHeight)
			e.Frame(gtx.Ops)
		}
	}
}

func toF32Pt(x, y float64) f32.Point {
	return f32.Pt(float32(x), float32(y))
}

func drawGrayBar1(gtx layout.Context, width, nbrBox, boxWidth, boxHeight, offsetX, offsetY float64) {
	defer op.Save(gtx.Ops).Load()
	paint.ColorOp{Color: color.NRGBA{A: 0xFF}}.Add(gtx.Ops)
	var p clip.Path
	p.Begin(gtx.Ops)
	for y := 0.; y < boxHeight; y++ {
		for x := 0.; x < nbrBox; x++ {
			lineHeight := 1. - x/nbrBox
			p.MoveTo(toF32Pt(x*boxWidth+offsetX, y+offsetY))
			p.LineTo(toF32Pt((x+1)*boxWidth+offsetX, y+offsetY))
			p.LineTo(toF32Pt((x+1)*boxWidth+offsetX, lineHeight+y+offsetY))
			p.LineTo(toF32Pt(x*boxWidth+offsetX, lineHeight+y+offsetY))
			p.Close()
		}
	}
	clip.Outline{Path: p.End()}.Op().Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func drawRect(gtx layout.Context, x, y, w, h float64, gray uint8) {
	defer op.Save(gtx.Ops).Load()
	var p clip.Path
	paint.ColorOp{Color: color.NRGBA{R: gray, G: gray, B: gray, A: 0xFF}}.Add(gtx.Ops)
	p.Begin(gtx.Ops)
	p.MoveTo(toF32Pt(x, y))
	p.LineTo(toF32Pt(x+w, y))
	p.LineTo(toF32Pt(x+w, y+h))
	p.LineTo(toF32Pt(x, y+h))
	p.Close()
	clip.Outline{Path: p.End()}.Op().Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func drawGrayBar2(gtx layout.Context, width, nbrBox, boxWidth, boxHeight, offsetX, offsetY float64) {
	defer op.Save(gtx.Ops).Load()
	for x := 0.; x < nbrBox; x++ {
		gray := uint8((x*255)/(nbrBox-1) + 0.5)
		drawRect(gtx, x*boxWidth+offsetX, offsetY, boxWidth, boxHeight, gray)
	}
}
