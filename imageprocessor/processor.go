package imageprocessor

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"os"
)

func ResizeImage(inputPath string, width, height uint) ([]byte, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	resizedImage := resize.Resize(width, height, img, resize.Lanczos3)

	outFile, err := os.Create("resized_temp.jpg")
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, resizedImage, nil)
	if err != nil {
		return nil, err
	}

	resizedData, err := os.ReadFile("resized_temp.jpg")
	if err != nil {
		return nil, err
	}

	return resizedData, nil
}
