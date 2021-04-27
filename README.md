# Testing Gio’s gamma correction

This small gio program crates a window showing four gray ramps.
The top and bottom ramps are identical and the reference ramps. 
They display what an ideal gray ramp must look like. 

The second gray ramp from the top creates the gray by using 
anti-aliasing. A 50% covered pixel should be perfectly mid-gray.

The third ramp from the top is simply drawing rectangles filled 
with linarly increasing values of gray. 0 is black and 255 is white. 

All gray ramps should be identical. The reference ramp is from 
["what every coder should know about gamma"](https://blog.johnnovak.net/2016/09/21/what-every-coder-should-know-about-gamma/) is visible below

## Usage

To install the program execute `go install github.com/chmike/testGioGamma@latest`.

To run it, execute `testGioGamma` provided that $GOPATH/bin is in your $PATH variable.