# Testing Gio’s gamma correction

This small gio program crates a window showing two gray ramps.

The top ramp creates the gray by using anti-aliasing.
A 50% covered pixel should be perfectly mid-gray.

The bottom ramp is simply drawing rectangles filled with linarly
increasing values of gray. 0 is black and 255 is white. 

Both ramps should be identical and match the reference ramp below. The 
reference ramp is from ["what every coder should know about gamma"](https://blog.johnnovak.net/2016/09/21/what-every-coder-should-know-about-gamma/) is visible below

![Reference](gamma-ramp32.png)

A pull request to show the reference ramp in the window next to the
two ramps would be great. I don’t know how to show a pixelmap image 
with Gio. 