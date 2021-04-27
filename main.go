// SPDX-License-Identifier: Unlicense OR MIT

package main

/*
Small Gio programm to test gamma correction.

The program draws four horizontal sequences of gray rectangles colored from
black (left) to white (right). Each raw is called a gray ramp. The top and
bottom rows are identical and are the reference ramp. They dispay the ideal
gray ramp. It is obtained from
https://blog.johnnovak.net/2016/09/21/what-every-coder-should-know-about-gamma/.

The gray ramp just below the top ramp test the effect of gamma correction
with anti-aliasing. Each box is drawn as a stack of black horizontal lines
not thicker than a pixel. When the line width covers the full pixel height,
the rectangle must be black. The percentage of pixel height covering
determines the gray value. When it is 50% the colour must be mid-gray, and
when it is 0%, the colour must be white. etc.

The second line, draws the boxes as simple rectangles filled with a specified
color. The color ramp increases linearly with a value in the range 0 to 255.

The bottom ramp is the again the reference ramp for easy comparison with the
gray ramp just above.

With correct anti-aliasing and gamma correction, all gray ramps must look
identical or very similar. If not, there is a problem.
*/

import (
	"bytes"
	_ "embed"
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

//go:embed gamma-ramp32.png
var gammaRamp32 []byte

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
	op.Offset(toF32Pt(offsetX, offsetY)).Add(gtx.Ops)
	var p clip.Path
	p.Begin(gtx.Ops)
	for y := 0.; y < boxHeight; y++ {
		for x := 0.; x < nbrBox; x++ {
			lineHeight := 1. - x/nbrBox
			p.MoveTo(toF32Pt(x*boxWidth, y))
			p.LineTo(toF32Pt((x+1)*boxWidth, y))
			p.LineTo(toF32Pt((x+1)*boxWidth, lineHeight+y))
			p.LineTo(toF32Pt(x*boxWidth, lineHeight+y))
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
		img, err = png.Decode(bytes.NewReader(gammaRamp32))
		if err != nil {
			log.Fatal("failed decoding reference ramp:", err)
		}
	}
	imgSize := img.Bounds().Size()
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), toF32Pt(w/float64(imgSize.X), h/float64(imgSize.Y))).Offset(toF32Pt(x, y))).Add(gtx.Ops)
	paint.NewImageOp(img).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
