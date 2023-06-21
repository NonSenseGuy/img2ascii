package img2ascii

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func IsFfmpegInstalled() error {
	var errBuffer bytes.Buffer
	cmd := exec.Command("ffmpeg", "-version")
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg not found: %w", err)
	}

	return nil
}

func ParseWebcamCapture(filepath string) error {
	var inBuffer, outBuffer, errBuffer bytes.Buffer
	cmd := exec.Command(
		"ffmpeg",
		"-f", "avfoundation",
		"-framerate", "30",
		"-i", "FaceTime HD Camera",
		"-video_size", "1920x1080",
		"-vframes", "1",
		filepath,
	)
	cmd.Stdin = &inBuffer
	cmd.Stderr = &errBuffer
	cmd.Stdout = &outBuffer

	err := cmd.Start()
	if err != nil {
		return err
	}

	fmt.Println(inBuffer.String())

	cmd.Wait()

	fmt.Println(inBuffer.String())

	// fmt.Println(outBuffer.String(), errBuffer.String())
	return nil
}

func InitWebcam() error {
	cmd := exec.Command("ffmpeg",
		"-f", "avfoundation",
		"-framerate", "30",
		"-pixel_format", "0rgb",
		"-i", "FaceTime HD Camera",
		"-video_size", "640*480",
		"-vcodec", "rawvideo",
		"-an", "-sn", "-fps_mode", "vfr",
		"-f", "image2pipe", "-",
	)

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	builder := bytes.Buffer{}
	buffer := make([]byte, 1024)
	for {
		n, err := pipe.Read(buffer)
		builder.Write(buffer[:n])
		if err == io.EOF {
			fmt.Println("eof")
			break
		}
	}

	fmt.Println(builder.String())

	cmd.Wait()

	return nil
}
