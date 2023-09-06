package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"
)

const (
	bytesPerPixel = 3
)

// Trims the extension of the file.
func trimFilenameExtension(filename string) string {
	pos := strings.LastIndexByte(filename, '.')
	if pos != -1 {
		return filename[:pos]
	}
	return filename
}

// Reads binary data from the file.
func readDataFromFile(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

// Generates an RGBA image from the data.
func createImage(data []byte, imageSize int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

	x, y := 0, 0
	for i := 0; i < len(data)-bytesPerPixel; i += bytesPerPixel {
		r, g, b := data[i], data[i+1], data[i+2]
		img.Set(x, y, color.NRGBA{R: r, G: g, B: b, A: 255})

		x++
		if x == imageSize {
			x = 0
			y++
			if y == imageSize {
				break
			}
		}
	}

	return img
}

// Encodes and saves the provided image to an output file.
func saveImageToFile(img *image.RGBA, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

// Reads data from the input file, creates an image, and saves it.
func process(file string) error {
	data, err := readDataFromFile(file)
	if err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	dataLength := int((len(data)))
	imageSize := int(math.Sqrt(float64(dataLength / bytesPerPixel)))

	img := createImage(data, imageSize)

	outputImage := trimFilenameExtension(file) + "_image.png"
	if err := saveImageToFile(img, outputImage); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	fmt.Println("Image successfully created:", outputImage)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	for _, filename := range os.Args[1:] {
		if err := process(filename); err != nil {
			fmt.Printf("Error processing %s: %v\n", filename, err)
		}
	}
}
