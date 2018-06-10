package pgm

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Decode decodes pgm binary formated images.
// Current limitations:
// * Only binary PGM (magic = P5).
// * Only one byte per pixel.
// * No comments allowed in the file.
func Decode(r io.Reader) (image.Image, error) {
	var magic string
	var w, h, max int
	n, err := fmt.Fscan(r, &magic, &w, &h, &max)
	if err != nil {
		return nil, errors.Wrap(err, "could decode all header information")
	}
	if n != 4 {
		return nil, errors.New("could decode all header information")
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not read image data")
	}
	// Skip whitespace.
	for b[0] == '\n' || b[0] == ' ' {
		b = b[1:]
	}
	if len(b) != w*h {
		return nil, errors.Errorf("with width of %v and height of %v I expect %v bytes but found %v", w, h, w*h, len(b))
	}
	return &image.Gray{
		Pix:    []uint8(b),
		Stride: w,
		Rect: image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: w, Y: h},
		},
	}, nil
}
