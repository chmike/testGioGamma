# Testing Gio’s gamma correction

This small gio program crates a window showing multiple gray ramps
to test gamma correction in Gio. Ramps are numbered from 1 to 6, from
top toward bottom. At the bottom I also display slightly slanted lines
to see the effect of AA on lines rendering. 

## Explanation of what each graw ramp shows.

1. This is an sRGB color ramp. The sRGB color value is linear
but the brightness is not linear. See  
["what every coder should know about gamma"](https://blog.johnnovak.net/2016/09/21/what-every-coder-should-know-about-gamma/) is visible below
for a more detailed explanation.
2. This gray ramp is linear in brightness if your screen is 
correctly calibrated. It is made by using anti-aliasing
to generate the gray. Each gray rectangle is drawn as a 
stack of thin black rectangles over white background with 
a height not exceeding a pixel. The thinner the rectangle
the whiter the gray. Anti-aliasing is performed in linear
color space (not sRGB), and the result is converted to
sRGB space before display. 
3. This gray ramp is made by simply filling the rectangles
with a gray color value in sRGB space. It is thus converted 
to linear color space before drawing. Once drawn, the 
resulting image is converted back to sRGB space for display. 
The color change from and to sRGB cancel out. 
4. This gray ramp is the same as gray ramp 1 to ease comparison.
5. This gray ramp is also made with anti-aliasing, but it is
made of thin white rectangles over a black background.
6. This is the same as gray ramp 2 to ease comparison. It should
appear identical to ramp 5.

The content of the GPU frame buffer is in sRGB color. Thus
a screen color picker will show colors in sRGB color space, 
and a screen capture will get colors in sRGB color space. 
That is fine since sRGB is an optimal encoding for 8 bit
values. 

## White strip artefact

In a previous program version an artefact appeared in gray 
ramp 3 when the right and left rectangle borders did not 
match exactly a pixel boundary. This is due to alpha blending.

When multiple shapes partially cover a pixel, the resulting 
color varies a lot by the way the shape cover the pixel. The 
alpha channel value can’t carry that information and artifact 
may appear. 

In the previous version, the artifact where bright vertical 
lines. Changing the background color to black made these 
artifacts less visible. But they are still visible in light 
gray rectangles as darker vertical lines.

These artifacts are much less visible with a high resolution 
screen because the artifacts are thinner. 

With gpu rendering and the computation power that comes with 
it, we could compute the exact color of the pixel and get rid
of these artifacts.


## Usage

To install the program execute `go install github.com/chmike/testGioGamma@latest`.

To run it, execute `testGioGamma` provided that $GOPATH/bin is in your $PATH variable.