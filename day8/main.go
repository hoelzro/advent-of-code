package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	width            = 25
	height           = 6
	blackColor       = '0'
	whiteColor       = '1'
	transparentColor = '2'
)

type layer struct {
	layerBytes []byte
}

func splitIntoLayers(rawImage []byte) []*layer {
	layerSize := width * height
	result := []*layer{}

	for len(rawImage) > 0 {
		layerBytes := rawImage[:layerSize]
		layer := &layer{
			layerBytes,
		}
		result = append(result, layer)
		rawImage = rawImage[layerSize:]
	}

	return result
}

func overlayLayers(layers []*layer) *layer {
	result := &layer{
		layerBytes: make([]byte, len(layers[0].layerBytes)),
	}

	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]
		for j, b := range layer.layerBytes {
			if b != transparentColor {
				result.layerBytes[j] = b
			}
		}
	}

	return result
}

func drawImage(layer *layer) {
	rawBytes := layer.layerBytes

	for len(rawBytes) > 0 {
		row := &strings.Builder{}

		for _, b := range rawBytes[:width] {
			if b == blackColor {
				row.WriteRune('█')
			} else if b == whiteColor {
				row.WriteRune(' ')
			}
		}

		fmt.Println(row.String())

		rawBytes = rawBytes[width:]
	}
}

func main() {
	imageBytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic("couldn't read file")
	}

	imageBytes = bytes.TrimSpace(imageBytes)

	layers := splitIntoLayers(imageBytes)

	// Part 1

	var layerWithFewestZeroes *layer
	fewestZeroes := width * height

	for _, layer := range layers {
		numZeroes := 0

		for _, b := range layer.layerBytes {
			if b == '0' {
				numZeroes++
			}
		}

		if numZeroes < fewestZeroes {
			fewestZeroes = numZeroes
			layerWithFewestZeroes = layer
		}
	}

	oneCount := 0
	twoCount := 0

	for _, b := range layerWithFewestZeroes.layerBytes {
		if b == '1' {
			oneCount++
		} else if b == '2' {
			twoCount++
		}
	}

	fmt.Println(oneCount * twoCount)

	// Part 2
	result := overlayLayers(layers)
	drawImage(result)
}
