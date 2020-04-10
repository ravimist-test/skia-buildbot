package sobel

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.skia.org/infra/go/testutils/unittest"
	"go.skia.org/infra/golden/go/image/text"
)

func TestSobelFuzzyMatcher_IdenticalImages_ReturnsTrue(t *testing.T) {
	unittest.SmallTest(t)

	// TODO(lovisolo): Implement.
}

// TestSobel_SmallImages_Success tests various edge cases involving small images.
//
// Notes:
// - In practice the sobel() function will usually be invoked with much larger images.
// - While these edge cases are unlikely to occur in practice, they exposed a number of bugs during
//   development (e.g. indexing and off-by-one errors). Those same bugs could be exposed with
//   larger, more realistic images, but using small images makes it easier to debug failing tests.
func TestSobel_SmallImages_Success(t *testing.T) {
	unittest.SmallTest(t)

	tests := []struct {
		name           string
		input          *image.Gray
		expectedOutput *image.Gray
	}{
		{
			name: "empty image, returns empty image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			0 0`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			0 0`),
		},
		{
			name: "1x1 black image, returns 1x1 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			1 1
			0x00`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			1 1
			0x00`),
		},
		{
			name: "1x1 white image, returns 1x1 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			1 1
			0xFF`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			1 1
			0x00`),
		},
		{
			name: "2x2 black image, returns 2x2 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			2 2
			0x00 0x00
			0x00 0x00`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			2 2
			0x00 0x00
			0x00 0x00`),
		},
		{
			name: "2x2 white image, returns 2x2 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			2 2
			0xFF 0xFF
			0xFF 0xFF`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			2 2
			0x00 0x00
			0x00 0x00`),
		},
		{
			name: "2x2 image with an edge, returns 2x2 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			2 2
			0x00 0xFF
			0x00 0xFF`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			2 2
			0x00 0x00
			0x00 0x00`),
		},
		{
			name: "3x3 black image, returns 3x3 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			3 3
			0x00 0x00 0x00
			0x00 0x00 0x00
			0x00 0x00 0x00`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			3 3
			0x00 0x00 0x00
			0x00 0x00 0x00
			0x00 0x00 0x00`),
		},
		{
			name: "3x3 white image, returns 3x3 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			3 3
			0xFF 0xFF 0xFF
			0xFF 0xFF 0xFF
			0xFF 0xFF 0xFF`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			3 3
			0x00 0x00 0x00
			0x00 0x00 0x00
			0x00 0x00 0x00`),
		},
		{
			name: "3x3 image with an edge, returns 3x3 image with one high intensity pixel",
			input: text.MustToGray(`! SKTEXTSIMPLE
			3 3
			0x00 0xFF 0xFF
			0x00 0xFF 0xFF
			0x00 0xFF 0xFF`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			3 3
			0x00 0x00 0x00
			0x00 0xFF 0x00
			0x00 0x00 0x00`),
		},
		{
			name: "4x4 black image, returns 4x4 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			4 4
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			4 4
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00`),
		},
		{
			name: "4x4 white image, returns 4x4 black image",
			input: text.MustToGray(`! SKTEXTSIMPLE
			4 4
			0xFF 0xFF 0xFF 0xFF
			0xFF 0xFF 0xFF 0xFF
			0xFF 0xFF 0xFF 0xFF
			0xFF 0xFF 0xFF 0xFF`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			4 4
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00
			0x00 0x00 0x00 0x00`),
		},
		{
			name: "4x4 image with an edge, returns 4x4 image with 4 high intensity pixels",
			input: text.MustToGray(`! SKTEXTSIMPLE
			4 4
			0x00 0x00 0xFF 0xFF
			0x00 0x00 0xFF 0xFF
			0x00 0x00 0xFF 0xFF
			0x00 0x00 0xFF 0xFF`),
			expectedOutput: text.MustToGray(`! SKTEXTSIMPLE
			4 4
			0x00 0x00 0x00 0x00
			0x00 0xFF 0xFF 0x00
			0x00 0xFF 0xFF 0x00
			0x00 0x00 0x00 0x00`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertImagesEqual(t, tc.expectedOutput, sobel(tc.input))
		})
	}
}

func TestSobel_EdgesAtVariousAngles_Success(t *testing.T) {
	unittest.SmallTest(t)

	input0Degrees := text.MustToGray(`! SKTEXTSIMPLE
	10 10
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x44 0x44 0x44 0x44 0x44`)

	expectedOutput0Degrees := text.MustToGray(`! SKTEXTSIMPLE
	10 10
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0xFF 0xFF 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0xFF 0xFF 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0xFF 0xFF 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0xFF 0xFF 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0xFF 0xFF 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0xFF 0xFF 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0xFF 0xFF 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0xFF 0xFF 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00`)

	input30Degrees := text.MustToGray(`! SKTEXTSIMPLE
	10 10
	0x00 0x00 0x00 0x00 0x00 0x00 0x0E 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x00 0x36 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x0E 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x36 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x0E 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x36 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x0E 0x44 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x36 0x44 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x0E 0x44 0x44 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x36 0x44 0x44 0x44 0x44 0x44 0x44 0x44`)

	expectedOutput30Degrees := text.MustToGray(`! SKTEXTSIMPLE
	10 10
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x13 0xCE 0xFF 0x62 0x00 0x00
	0x00 0x00 0x00 0x00 0x62 0xFF 0xCE 0x13 0x00 0x00
	0x00 0x00 0x00 0x13 0xCE 0xFF 0x62 0x00 0x00 0x00
	0x00 0x00 0x00 0x62 0xFF 0xCE 0x13 0x00 0x00 0x00
	0x00 0x00 0x13 0xCE 0xFF 0x62 0x00 0x00 0x00 0x00
	0x00 0x00 0x62 0xFF 0xCE 0x13 0x00 0x00 0x00 0x00
	0x00 0x13 0xCE 0xFF 0x62 0x00 0x00 0x00 0x00 0x00
	0x00 0x62 0xFF 0xCE 0x13 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00`)

	input45Degrees := text.MustToGray(`! SKTEXTSIMPLE
	10 10
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x05 0x3F
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x05 0x3F 0x44
	0x00 0x00 0x00 0x00 0x00 0x00 0x05 0x3F 0x44 0x44
	0x00 0x00 0x00 0x00 0x00 0x05 0x3F 0x44 0x44 0x44
	0x00 0x00 0x00 0x00 0x05 0x3F 0x44 0x44 0x44 0x44
	0x00 0x00 0x00 0x05 0x3F 0x44 0x44 0x44 0x44 0x44
	0x00 0x00 0x05 0x3F 0x44 0x44 0x44 0x44 0x44 0x44
	0x00 0x05 0x3F 0x44 0x44 0x44 0x44 0x44 0x44 0x44
	0x05 0x3F 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44
	0x3F 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44`)

	expectedOutput45Degrees := text.MustToGray(`! SKTEXTSIMPLE
	10 10
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x07 0x67 0xFF 0xFF 0x00
	0x00 0x00 0x00 0x00 0x07 0x67 0xFF 0xFF 0x67 0x00
	0x00 0x00 0x00 0x07 0x67 0xFF 0xFF 0x67 0x07 0x00
	0x00 0x00 0x07 0x67 0xFF 0xFF 0x67 0x07 0x00 0x00
	0x00 0x07 0x67 0xFF 0xFF 0x67 0x07 0x00 0x00 0x00
	0x00 0x67 0xFF 0xFF 0x67 0x07 0x00 0x00 0x00 0x00
	0x00 0xFF 0xFF 0x67 0x07 0x00 0x00 0x00 0x00 0x00
	0x00 0xFF 0x67 0x07 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00`)

	input60Degrees := text.MustToGray(`! SKTEXTSIMPLE
	10 10
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x0E 0x36
	0x00 0x00 0x00 0x00 0x00 0x00 0x0E 0x36 0x44 0x44
	0x00 0x00 0x00 0x00 0x0E 0x36 0x44 0x44 0x44 0x44
	0x00 0x00 0x0E 0x36 0x44 0x44 0x44 0x44 0x44 0x44
	0x0E 0x36 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44
	0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44
	0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44
	0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44 0x44`)

	expectedOutput60Degrees := text.MustToGray(`! SKTEXTSIMPLE
	10 10
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x13 0x62 0x00
	0x00 0x00 0x00 0x00 0x00 0x13 0x62 0xCE 0xFF 0x00
	0x00 0x00 0x00 0x13 0x62 0xCE 0xFF 0xFF 0xCE 0x00
	0x00 0x13 0x62 0xCE 0xFF 0xFF 0xCE 0x62 0x13 0x00
	0x00 0xCE 0xFF 0xFF 0xCE 0x62 0x13 0x00 0x00 0x00
	0x00 0xFF 0xCE 0x62 0x13 0x00 0x00 0x00 0x00 0x00
	0x00 0x62 0x13 0x00 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
	0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00`)

	test := func(name string, input, expectedOutput *image.Gray) {
		t.Run(name, func(t *testing.T) {
			assertImagesEqual(t, expectedOutput, sobel(input))
		})
	}

	test("0 degrees", input0Degrees, expectedOutput0Degrees)
	test("30 degrees", input30Degrees, expectedOutput30Degrees)
	test("45 degrees", input45Degrees, expectedOutput45Degrees)
	test("60 degrees", input60Degrees, expectedOutput60Degrees)
}

func TestSobel_GoldenImage_Success(t *testing.T) {
	unittest.MediumTest(t)

	// Attribution for the test/input.png image used below:
	//
	//   Author: Simpsons contributor.
	//   License: CC BY-SA (https://creativecommons.org/licenses/by-sa/3.0).
	//   Source: https://en.wikipedia.org/wiki/File:Valve_original_%281%29.PNG.
	//   Modifications: PNG image was recoded using Golang's png.Decode() and png.Encode().

	input := readPngAsGray(t, "test/input.png")
	expectedOutput := readPngAsGray(t, "test/sobel-expected-output.png")
	assert.Equal(t, expectedOutput, sobel(input))
}

func TestZeroOutEdges_Success(t *testing.T) {
	unittest.SmallTest(t)

	tests := []struct {
		name           string
		input          image.Image
		edges          *image.Gray
		edgeThreshold  uint8
		expectedOutput image.Image
	}{
		{
			name: "empty image, returns empty image",
			input: text.MustToNRGBA(`! SKTEXTSIMPLE
			0 0`),
			edges: text.MustToGray(`! SKTEXTSIMPLE
			0 0`),
			edgeThreshold: 0,
			expectedOutput: text.MustToNRGBA(`! SKTEXTSIMPLE
			0 0`),
		},
		{
			name: "1x1 image with pixel below threshold, returns original image",
			input: text.MustToNRGBA(`! SKTEXTSIMPLE
			1 1
			0xAABBCCFF`),
			edges: text.MustToGray(`! SKTEXTSIMPLE
			1 1
			0x55`),
			edgeThreshold: 0xBB,
			expectedOutput: text.MustToNRGBA(`! SKTEXTSIMPLE
			1 1
			0xAABBCCFF`),
		},
		{
			name: "1x1 image with pixel above threshold, returns black image",
			input: text.MustToNRGBA(`! SKTEXTSIMPLE
			1 1
			0xAABBCCFF`),
			edges: text.MustToGray(`! SKTEXTSIMPLE
			1 1
			0xCC`),
			edgeThreshold: 0xBB,
			expectedOutput: text.MustToNRGBA(`! SKTEXTSIMPLE
			1 1
			0x000000FF`),
		},
		{
			name: "image with some pixels above threshold, returns image with zeroed out edges",
			input: text.MustToNRGBA(`! SKTEXTSIMPLE
			5 5
			0x111111FF 0x222222FF 0x333333FF 0x555555FF 0x888888FF
			0x111111FF 0x222222FF 0x333333FF 0x555555FF 0x888888FF
			0x111111FF 0x222222FF 0x333333FF 0x555555FF 0x888888FF
			0x111111FF 0x222222FF 0x333333FF 0x555555FF 0x888888FF
			0x111111FF 0x222222FF 0x333333FF 0x555555FF 0x888888FF`),
			edges: text.MustToGray(`! SKTEXTSIMPLE
			5 5
			0x00 0x00 0x00 0x00 0x00
			0x00 0x88 0xCC 0xFF 0x00
			0x00 0x88 0xCC 0xFF 0x00
			0x00 0x88 0xCC 0xFF 0x00
			0x00 0x00 0x00 0x00 0x00`),
			edgeThreshold: 0xAA,
			expectedOutput: text.MustToNRGBA(`! SKTEXTSIMPLE
			5 5
			0x111111FF 0x222222FF 0x333333FF 0x555555FF 0x888888FF
			0x111111FF 0x222222FF 0x000000FF 0x000000FF 0x888888FF
			0x111111FF 0x222222FF 0x000000FF 0x000000FF 0x888888FF
			0x111111FF 0x222222FF 0x000000FF 0x000000FF 0x888888FF
			0x111111FF 0x222222FF 0x333333FF 0x555555FF 0x888888FF`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertImagesEqual(t, tc.expectedOutput, zeroOutEdges(tc.input, tc.edges, tc.edgeThreshold))
		})
	}
}

func TestZeroOutEdges_GoldenImage_Success(t *testing.T) {
	unittest.MediumTest(t)

	input := readPng(t, "test/input.png")
	edges := readPngAsGray(t, "test/sobel-expected-output.png")
	expectedOutput := readPng(t, "test/zero-out-edges-expected-output.png")
	assert.Equal(t, expectedOutput, zeroOutEdges(input, edges, 0x55))
}

func TestZeroOutEdges_InputAndEdgesImagesHaveDifferentBounds_Panics(t *testing.T) {
	unittest.SmallTest(t)

	assert.Panics(t, func() {
		img := text.MustToNRGBA(`! SKTEXTSIMPLE
		2 1
		0x00 0x00`)
		edges := text.MustToGray(`! SKTEXTSIMPLE
		1 1
		0x00`)
		zeroOutEdges(img, edges, 0)
	})
}

// assertImagesEqual asserts that the two given images are equal, and prints out the actual image
// encoded as SKTEXT if the assertion is false.
func assertImagesEqual(t *testing.T, expected, actual image.Image) {
	assert.Equal(t, expected, actual, fmt.Sprintf("SKTEXT-encoded actual output:\n%s", imageToText(t, actual)))
}

// imageToText returns the given image as an SKTEXT-encoded string.
func imageToText(t *testing.T, img image.Image) string {
	// Convert image to NRGBA.
	nrgbaImg := image.NewNRGBA(img.Bounds())
	draw.Draw(nrgbaImg, img.Bounds(), img, img.Bounds().Min, draw.Src)

	// Encode and return image as SKTEXTSIMPLE.
	buf := &bytes.Buffer{}
	err := text.Encode(buf, nrgbaImg)
	require.NoError(t, err)
	return buf.String()
}

// readPngAsGray reads a PNG image from the file system, converts it to grayscale and returns it as
// an *image.Gray.
func readPngAsGray(t *testing.T, filename string) *image.Gray {
	// Read image.
	img := readPng(t, filename)

	// Convert to grayscale.
	grayImg := image.NewGray(img.Bounds())
	draw.Draw(grayImg, img.Bounds(), img, img.Bounds().Min, draw.Src)

	return grayImg
}

// readPng reads a PNG image from the file system and returns it as an *image.NRGBA.
func readPng(t *testing.T, filename string) *image.NRGBA {
	// Read image.
	imgBytes, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	// Decode image.
	img, err := png.Decode(bytes.NewReader(imgBytes))
	require.NoError(t, err)

	// Convert to NRGBA.
	nrgbaImg := image.NewNRGBA(img.Bounds())
	draw.Draw(nrgbaImg, img.Bounds(), img, img.Bounds().Min, draw.Src)

	return nrgbaImg
}
