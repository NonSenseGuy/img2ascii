package main

import (
	"ascii_img/img2ascii"
	"fmt"
	"os"
)

func main() {
	err := img2ascii.ParseWebcamCapture(os.Args[1])
	if err != nil {
		panic(err)
	}

	ib, err := img2ascii.ImageBufferFromFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	fmt.Println(ib.Width, ib.Height)
	ib.WriteImageFromIntensityMatrix()
	ib.PrintAsciiImage()
}
