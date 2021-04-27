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
	"image"
	"image/color"
	"image/png"
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
			boxHeight := math.Ceil((boxWidth) * 2)
			drawRefImg(gtx, 0, 0, width, boxHeight)
			drawGrayBar1(gtx, width, nbrBox, boxWidth, boxHeight, 0, boxHeight)
			drawGrayBar2(gtx, width, nbrBox, boxWidth, boxHeight, 0, 2*boxHeight)
			drawRefImg(gtx, 0, 3*boxHeight, width, boxHeight)
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
	op.Offset(toF32Pt(x, y)).Add(gtx.Ops)
	p.Begin(gtx.Ops)
	p.MoveTo(toF32Pt(0, 0))
	p.LineTo(toF32Pt(w, 0))
	p.LineTo(toF32Pt(w, h))
	p.LineTo(toF32Pt(0, h))
	p.Close()
	clip.Outline{Path: p.End()}.Op().Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func drawGrayBar2(gtx layout.Context, width, nbrBox, boxWidth, boxHeight, offsetX, offsetY float64) {
	for x := 0.; x < nbrBox; x++ {
		gray := uint8((x*255)/(nbrBox-1) + 0.5)
		drawRect(gtx, x*boxWidth+offsetX, offsetY, boxWidth, boxHeight, gray)
	}
}

var img image.Image

func drawRefImg(gtx layout.Context, x, y, w, h float64) {
	defer op.Save(gtx.Ops).Load()
	if img == nil {
		var err error
		img, err = loadPng("gamma-ramp32.png")
		if err != nil {
			log.Fatal("failed loading reference ramp:", err)
		}
	}
	imgSize := img.Bounds().Size()
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), toF32Pt(w/float64(imgSize.X), h/float64(imgSize.Y))).Offset(toF32Pt(x, y))).Add(gtx.Ops)
	paint.NewImageOp(img).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func saveAsPng(img image.Image, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		f.Close()
		os.Remove(fileName)
		return err
	}
	return nil
}

func loadPng(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}
