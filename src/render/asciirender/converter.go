package asciirender

import (
	"fmt"
	"image"
	"image/color"
)

type AsciiPixel struct {
	charDepth      uint32
	grayscaleValue [3]uint32
	rgbValue       [3]uint32
}

type AsciiChar struct {
	OriginalColor string
	SetColor      string
	Simple        string
	RgbValue      [3]uint32
}

var _asciiTableSimple = " .:-=+*#%@"
var _asciiTableDetailed = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

type AsciiArt struct {
	termColorLevel Level
	negative       bool
	colored        bool
	grayscale      bool
	complex        bool
	colorBg        bool
	customMap      string
	fontColor      [3]int
}

func NewAsciiArt(negative bool, colored bool, grayscale bool, complex bool, colorBg bool, customMap string, fontColor [3]int) *AsciiArt {
	return &AsciiArt{
		negative:       negative,
		colored:        colored,
		grayscale:      grayscale,
		complex:        complex,
		colorBg:        colorBg,
		customMap:      customMap,
		termColorLevel: _colorLevel,
		fontColor:      fontColor,
	}
}

func (a *AsciiArt) ConvertToAsciiPixels(img image.Image, dimensions []int, width, height int, flipX bool, flipY bool, full bool) ([][]AsciiPixel, error) {
	//smallImg, err := a.resizeImage(img, full, dimensions, width, height)
	//TODO RESIZE
	smallImg := img
	//if err != nil {
	//	return nil, err
	//}
	var imgSet [][]AsciiPixel
	b := smallImg.Bounds()
	// These nested loops iterate through each pixel of resized image and get an AsciiPixel instance
	for y := b.Min.Y; y < b.Max.Y; y++ {
		var temp []AsciiPixel
		for x := b.Min.X; x < b.Max.X; x++ {
			oldPixel := smallImg.At(x, y)
			grayPixel := color.GrayModel.Convert(oldPixel)
			r1, g1, b1, _ := grayPixel.RGBA()
			charDepth := r1 / 257 // Only Red is needed from RGB for charDepth in AsciiPixel since they have the same value for grayscale images
			r1 = r1 / 257
			g1 = g1 / 257
			b1 = b1 / 257
			r2, g2, b2, _ := oldPixel.RGBA()
			r2 = r2 / 257
			g2 = g2 / 257
			b2 = b2 / 257
			temp = append(temp, AsciiPixel{
				charDepth:      charDepth,
				grayscaleValue: [3]uint32{r1, g1, b1},
				rgbValue:       [3]uint32{r2, g2, b2},
			})
		}
		imgSet = append(imgSet, temp)
	}
	if flipX || flipY {
		imgSet = a.reverse(imgSet, flipX, flipY)
	}
	return imgSet, nil
}

func (a *AsciiArt) Convert(imgSet [][]AsciiPixel) ([][]AsciiChar, error) {
	const maxVal float64 = 255
	height := len(imgSet)
	width := len(imgSet[0])
	chosenTable := make(map[int]string)
	if len(a.customMap) == 0 {
		var charSet string
		if a.complex {
			charSet = _asciiTableDetailed
		} else {
			charSet = _asciiTableSimple
		}
		for index, char := range charSet {
			chosenTable[index] = string(char)
		}
	} else {
		for index, char := range a.customMap {
			chosenTable[index] = string(char)
		}
	}
	var result [][]AsciiChar
	for i := 0; i < height; i++ {
		var tempSlice []AsciiChar
		for j := 0; j < width; j++ {
			value := float64(imgSet[i][j].charDepth)
			tempFloat := (value / maxVal) * float64(len(chosenTable))
			if value == maxVal {
				tempFloat = float64(len(chosenTable) - 1)
			}
			tempInt := int(tempFloat)
			var r, g, b int
			if a.colored {
				r = int(imgSet[i][j].rgbValue[0])
				g = int(imgSet[i][j].rgbValue[1])
				b = int(imgSet[i][j].rgbValue[2])
			} else {
				r = int(imgSet[i][j].grayscaleValue[0])
				g = int(imgSet[i][j].grayscaleValue[1])
				b = int(imgSet[i][j].grayscaleValue[2])
			}
			if a.negative {
				r = 255 - r
				g = 255 - g
				b = 255 - b
				if a.colored {
					imgSet[i][j].rgbValue = [3]uint32{uint32(r), uint32(g), uint32(b)}
				} else {
					imgSet[i][j].grayscaleValue = [3]uint32{uint32(r), uint32(g), uint32(b)}
				}
				tempInt = (len(chosenTable) - 1) - tempInt
			}
			var char AsciiChar
			asciiChar := chosenTable[tempInt]
			char.Simple = asciiChar
			var err error
			if a.colorBg {
				char.OriginalColor, err = a.getColoredCharForTerm(uint8(r), uint8(g), uint8(b), asciiChar, true)
			} else {
				char.OriginalColor, err = a.getColoredCharForTerm(uint8(r), uint8(g), uint8(b), asciiChar, false)
			}
			if (a.colored || a.grayscale) && err != nil {
				return nil, err
			}
			if a.fontColor != [3]int{255, 255, 255} {
				fcR := a.fontColor[0]
				fcG := a.fontColor[1]
				fcB := a.fontColor[2]
				if a.colorBg {
					char.SetColor, err = a.getColoredCharForTerm(uint8(fcR), uint8(fcG), uint8(fcB), asciiChar, true)
				} else {
					char.SetColor, err = a.getColoredCharForTerm(uint8(fcR), uint8(fcG), uint8(fcB), asciiChar, false)
				}
				if err != nil {
					return nil, err
				}
			}
			if a.colored {
				char.RgbValue = imgSet[i][j].rgbValue
			} else {
				char.RgbValue = imgSet[i][j].grayscaleValue
			}
			tempSlice = append(tempSlice, char)
		}
		result = append(result, tempSlice)
	}
	return result, nil
}

func (a *AsciiArt) getColoredCharForTerm(r uint8, g uint8, b uint8, char string, background bool) (string, error) {
	var coloredChar string
	if a.termColorLevel == ColorLevelMillions {
		colorRenderer := RGB(r, g, b, background)
		coloredChar = colorRenderer.Sprintf("%v", char)
	} else if a.termColorLevel == ColorLevelHundreds {
		colorRenderer := RGB(r, g, b, background).C256()
		coloredChar = colorRenderer.Sprintf("%v", char)
	} else {
		return "", fmt.Errorf("your terminal supports neither 24-bit nor 8-bit colors. Other coloring options aren't available")
	}
	return coloredChar, nil
}

/*
func (a *AsciiArt) resizeImage(img image.Image, full bool, dimensions []int, width, height int) (image.Image, error) {

	var asciiWidth, asciiHeight int
	var smallImg image.Image

	imgWidth := float64(img.Bounds().Dx())
	imgHeight := float64(img.Bounds().Dy())
	aspectRatio := imgWidth / imgHeight

	if full {
		terminalWidth, _, err := winsize.GetTerminalSize()
		if err != nil {
			return nil, err
		}
		asciiWidth = terminalWidth - 1
		asciiHeight = int(float64(asciiWidth) / aspectRatio)
		asciiHeight = int(0.5 * float64(asciiHeight))
	} else if (width != 0 || height != 0) && len(dimensions) == 0 {
		if width != 0 && height == 0 {
			asciiWidth = width
			asciiHeight = int(float64(asciiWidth) / aspectRatio)
			asciiHeight = int(0.5 * float64(asciiHeight))
			if asciiHeight == 0 {
				asciiHeight = 1
			}
		} else if height != 0 && width == 0 {
			asciiHeight = height
			asciiWidth = int(float64(asciiHeight) * aspectRatio)
			asciiWidth = int(2 * float64(asciiWidth))
			if asciiWidth == 0 {
				asciiWidth = 1
			}
		} else {
			return nil, fmt.Errorf("error: both width and height can't be set. Use dimensions instead")
		}
	} else if len(dimensions) == 0 {
		terminalWidth, terminalHeight, err := winsize.GetTerminalSize()
		if err != nil {
			return nil, err
		}
		asciiHeight = terminalHeight - 1
		asciiWidth = int(float64(asciiHeight) * aspectRatio)
		asciiWidth = int(2 * float64(asciiWidth))
		if asciiWidth >= terminalWidth {
			asciiWidth = terminalWidth - 1
			asciiHeight = int(float64(asciiWidth) / aspectRatio)
			asciiHeight = int(0.5 * float64(asciiHeight))
		}
	} else {
		asciiWidth = dimensions[0]
		asciiHeight = dimensions[1]
	}

	//smallImg = imaging.Resize(img, asciiWidth, asciiHeight, imaging.Lanczos)
	//return smallImg, nil
}
*/

func (a *AsciiArt) reverse(imgSet [][]AsciiPixel, flipX, flipY bool) [][]AsciiPixel {
	if flipX {
		for _, row := range imgSet {
			for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
				row[i], row[j] = row[j], row[i]
			}
		}
	}
	if flipY {
		for i, j := 0, len(imgSet)-1; i < j; i, j = i+1, j-1 {
			imgSet[i], imgSet[j] = imgSet[j], imgSet[i]
		}
	}
	return imgSet
}
