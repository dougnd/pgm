package pgm

import (
	"flag"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	testImg := "P5\n2 1\n255\naa"

	i, err := Decode(strings.NewReader(testImg))
	assert.NotNil(t, i)
	assert.NoError(t, err)
	assert.Equal(t, image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 2, Y: 1},
	}, i.Bounds())
	assert.Equal(t, color.Gray{Y: 97}, i.At(0, 0))
	assert.Equal(t, color.Gray{Y: 97}, i.At(1, 0))
}

var testImageDir = flag.String("image-dir", "", "The test image directory")

func TestDecodeImageFiles(t *testing.T) {
	if testImageDir == nil || *testImageDir == "" {
		t.Skip("Skipping image decode test. To enable, supply the -image-dir flag.")
	}
	files, err := filepath.Glob(*testImageDir + "/*.pgm")
	if err != nil {
		t.Fatalf("cannot read files from directory %v", *testImageDir)
	}
	if len(files) == 0 {
		t.Errorf("no files in directory %v", *testImageDir)
	}

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			r, err := os.Open(f)
			require.NoError(t, err)
			img, err := Decode(r)
			assert.NoError(t, err)
			assert.NotNil(t, img)
		})
	}

}

func BenchmarkDecodeImageFiles(b *testing.B) {
	if testImageDir == nil || *testImageDir == "" {
		b.Skip("Skipping image decode test. To enable, supply the -image-dir flag.")
	}
	files, err := filepath.Glob(*testImageDir + "/*.pgm")
	if err != nil {
		b.Fatalf("cannot read files from directory %v", *testImageDir)
	}
	if len(files) == 0 {
		b.Errorf("no files in directory %v", *testImageDir)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := files[i%len(files)]
		b.Log(f)
		r, err := os.Open(f)
		require.NoError(b, err)
		img, err := Decode(r)
		assert.NoError(b, err)
		assert.NotNil(b, img)
	}
}
