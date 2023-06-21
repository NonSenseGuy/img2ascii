package img2ascii

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"sort"

	"github.com/nfnt/resize"
	"golang.org/x/term"
)

const ASCII_CHARS = "`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

type ImageBuffer struct {
	Width           int
	Height          int
	IntensityMatrix [][]uint8
	AsciiMatrix     [][]rune
}

func NewImageBuffer(width, height int, intensityMatrix [][]uint8, asciiMatrix [][]rune) *ImageBuffer {
	return &ImageBuffer{
		Width:           width,
		Height:          height,
		IntensityMatrix: intensityMatrix,
		AsciiMatrix:     asciiMatrix,
	}
}

func ReadImageFromFile(filepath string) (image.Image, string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	return image.Decode(file)
}

func ImageBufferFromFile(filepath string) (*ImageBuffer, error) {
	img, _, err := ReadImageFromFile(filepath)
	if err != nil {
		return nil, err
	}

	width, height, err := term.GetSize(0)

	img, _ = resizeImage(img, uint(width), uint(height))

	return ImageBufferFromImage(img), nil
}

func ImageBufferFromImage(img image.Image) *ImageBuffer {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	intensityMatrix := getIntensityMatrix(img, width, height)
	intensityMatrix = normalizeIntensityMatrix(intensityMatrix)
	asciiMatrix := getAsciiChars(intensityMatrix, width, height)

	return NewImageBuffer(width, height, intensityMatrix, asciiMatrix)
}

func getIntensityMatrix(img image.Image, width, height int) [][]uint8 {
	intensityMatrix := make([][]uint8, height)
	for y := 0; y < height; y++ {
		intensityMatrix[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			intensityMatrix[y][x] = gray.Y
		}
	}

	return intensityMatrix
}

func normalizeIntensityMatrix(intensityMatrix [][]uint8) [][]uint8 {
	min, max := getMinMaxPixel(intensityMatrix)
	fmt.Println(min, max)
	for i := 0; i < len(intensityMatrix); i++ {
		for j := 0; j < len(intensityMatrix[i]); j++ {
			denominator := float64(max) - float64(min)
			num := float64(intensityMatrix[i][j] - min)
			normalizedPixel := 255 * num / denominator
			intensityMatrix[i][j] = uint8(normalizedPixel)
		}
	}

	return intensityMatrix
}

func getAsciiChars(intensityMatrix [][]uint8, width, height int) [][]rune {
	asciiMatrix := make([][]rune, height)
	for y := 0; y < height; y++ {
		asciiMatrix[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			pixel := float64(intensityMatrix[y][x])
			// fit 0..255 into 0..65
			lenAscii := float64(len(ASCII_CHARS))
			asciiIndex := int(pixel / float64(255) * (lenAscii - 1))
			asciiChar := ASCII_CHARS[asciiIndex]
			asciiMatrix[y][x] = rune(asciiChar)
		}
	}

	return asciiMatrix
}

func resizeImage(image image.Image, maxW, maxH uint) (image.Image, error) {
	resizedImage := resize.Thumbnail(maxW, maxH, image, resize.Lanczos3)

	return resizedImage, nil
}

func (ib *ImageBuffer) PrintAsciiImage() {
	for i := 0; i < ib.Height; i++ {
		fmt.Print("\n")
		for j := 0; j < ib.Width; j++ {
			fmt.Printf("%c", ib.AsciiMatrix[i][j])
			fmt.Printf("%c", ib.AsciiMatrix[i][j])
			fmt.Printf("%c", ib.AsciiMatrix[i][j])
		}
	}
}

func (ib *ImageBuffer) WriteImageFromIntensityMatrix() {
	img := image.NewGray(image.Rect(0, 0, ib.Width, ib.Height))
	for y := 0; y < ib.Height; y++ {
		for x := 0; x < ib.Width; x++ {
			img.SetGray(x, y, color.Gray{Y: ib.IntensityMatrix[y][x]})
		}
	}

	f, err := os.Create("files/output.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
}

func getMinMaxPixel(intensityMatrix [][]uint8) (uint8, uint8) {
	flat_img := flatten(intensityMatrix)
	sort.Slice(flat_img, func(i, j int) bool {
		return flat_img[i] < flat_img[j]
	})

	return flat_img[0], flat_img[len(flat_img)-1]
}

func flatten(matrix [][]uint8) []uint8 {
	flat_arr := make([]uint8, len(matrix[0])*len(matrix))
	for _, row := range matrix {
		flat_arr = append(flat_arr, row...)
	}

	fmt.Println(len(flat_arr))

	return flat_arr
}
