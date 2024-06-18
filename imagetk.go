package imagetk

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	tk "github.com/topxeq/tkc"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var versionG = "0.9a"

type ImageTK struct {
	Version string
}

var ITKX = &ImageTK{Version: versionG}

func NewImageTK() *ImageTK {
	return &ImageTK{Version: versionG}
}

func (p *ImageTK) GetVersion() string {
	return p.Version
}

func (p *ImageTK) NewRGBA(r, g, b, a uint8) color.RGBA {
	return color.RGBA{r, g, b, a}
}

func (p *ImageTK) NewRGBAP(r, g, b, a uint8) *color.RGBA {
	return &color.RGBA{r, g, b, a}
}

func (p *ImageTK) NewNRGBAFromHex(strA string) color.NRGBA {
	r, g, b, a := tk.ParseHexColor(strA)

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (p *ImageTK) NewNRGBAPFromHex(strA string) *color.NRGBA {
	r, g, b, a := tk.ParseHexColor(strA)

	return &color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (p *ImageTK) NewRGBAFromHex(strA string) color.RGBA {
	r, g, b, a := tk.ParseHexColor(strA)

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (p *ImageTK) NewRGBAPFromHex(strA string) *color.RGBA {
	r, g, b, a := tk.ParseHexColor(strA)

	return &color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (p *ImageTK) LoadRGBAFromImage(imageA image.Image) (*image.RGBA, error) {
	switch imageT := imageA.(type) {
	case *image.RGBA:
		return imageT, nil
	default:
		rgba := image.NewRGBA(imageT.Bounds())
		draw.Draw(rgba, imageT.Bounds(), imageT, image.Pt(0, 0), draw.Src)
		return rgba, nil
	}

}

func (p *ImageTK) LoadPlotImage(plt *plot.Plot, w vg.Length, h vg.Length) (*image.RGBA, error) {

	var bufT bytes.Buffer

	writerT, errT := plt.WriterTo(w, h, "png")

	if errT != nil {
		return nil, errT
	}

	_, errT = writerT.WriteTo(&bufT)

	if errT != nil {
		return nil, errT
	}

	readerT := bytes.NewReader(bufT.Bytes())

	// defer readerT.Close()

	// imgFile, err := os.Open(imgPath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer imgFile.Close()

	img, err := png.Decode(readerT)
	if err != nil {
		return nil, err
	}

	switch trueImg := img.(type) {
	case *image.RGBA:
		return trueImg, nil
	default:
		rgba := image.NewRGBA(trueImg.Bounds())
		draw.Draw(rgba, trueImg.Bounds(), trueImg, image.Pt(0, 0), draw.Src)
		return rgba, nil
	}
}

func (p *ImageTK) NewPlotXY(xA, yA float64) plotter.XY {
	return plotter.XY{X: xA, Y: yA}
}

// ParseHexColor inspired by gg
func (p *ImageTK) ParseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255

	if len(x) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}

	if len(x) == 6 {
		format := "%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b)
	}

	if len(x) == 8 {
		format := "%02x%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b, &a)
	}

	return
}

func (p *ImageTK) EncodePNG(imgA image.Image) {
}

// LoadPlotImageInMemory formatA support png, jpg...
func (p *ImageTK) LoadPlotImageInMemory(plotA *plot.Plot, w vg.Length, h vg.Length, formatA string) (*bytes.Buffer, error) {

	var bufT bytes.Buffer

	writerT, errT := plotA.WriterTo(w, h, formatA)

	if errT != nil {
		return nil, errT
	}

	_, errT = writerT.WriteTo(&bufT)

	if errT != nil {
		return nil, errT
	}

	return &bufT, nil

}

func (p *ImageTK) SaveImageAs(imageA image.Image, filePathA string, formatA ...string) error {
	fileT, errT := os.Create(filePathA)
	if errT != nil {
		return errT
	}
	defer fileT.Close()

	var formatT string

	if formatA == nil || len(formatA) < 1 {
		formatT = ".png"
	} else {
		formatT = tk.ToLower(formatA[0])
	}

	switch formatT {
	case "", ".png":
		errT = png.Encode(fileT, imageA)
	case ".jpg", ".jpeg":
		errT = jpeg.Encode(fileT, imageA, nil)
	case ".gif":
		errT = gif.Encode(fileT, imageA, nil)
	default:
		errT = png.Encode(fileT, imageA)
	}

	return errT

}

func (p *ImageTK) GetImageFileContent(fileNameA string, imageTypeA string) image.Image {
	file, err := os.Open(fileNameA)
	if err != nil {
		return nil
	}

	fileExt := strings.ToLower(filepath.Ext(fileNameA))
	var img image.Image

	if imageTypeA != "" {
		fileExt = imageTypeA
	}

	// tk.Pl("fileExt: %v", fileExt)

	switch fileExt {
	case ".jpg", ".jpeg", "jpg", "jpeg":
		img, err = jpeg.Decode(file)
	case ".png", "png":
		img, err = png.Decode(file)
	case ".gif", "gif":
		img, err = gif.Decode(file)
	default:
		img, err = jpeg.Decode(file)
	}

	if err != nil {
		// tk.Pl("err: %v", err)
		return nil
	}
	defer file.Close()

	return img
}

func (p *ImageTK) LoadImage(fileNameA string) (image.Image, error) {
	file, err := os.Open(fileNameA)
	if err != nil {
		return nil, err
	}

	fileExt := strings.ToLower(filepath.Ext(fileNameA))
	var img image.Image

	switch fileExt {
	case ".jpg", ".jpeg", "jpg", "jpeg":
		img, err = jpeg.Decode(file)
	case ".png", "png":
		img, err = png.Decode(file)
	case ".gif", "gif":
		img, err = gif.Decode(file)
	default:
		img, err = jpeg.Decode(file)
	}

	if err != nil {
		// tk.Pl("err: %v", err)
		return nil, err
	}
	defer file.Close()

	return img, nil
}

func (p *ImageTK) GetImageFileContentAndThumb(fileNameA string, maxWidthA uint, maxHeightA uint, imageTypeA string) image.Image {
	img := p.GetImageFileContent(fileNameA, imageTypeA)
	if img == nil {
		return nil
	}

	m := p.Thumbnail(maxWidthA, maxHeightA, img, MitchellNetravali)

	return m
}

func (p *ImageTK) Thumbnail(maxWidth, maxHeight uint, img image.Image, interp InterpolationFunction) image.Image {
	origBounds := img.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())
	newWidth, newHeight := origWidth, origHeight

	// Return original image if it have same or smaller size as constraints
	if maxWidth >= origWidth && maxHeight >= origHeight {
		return img
	}

	// Preserve aspect ratio
	if origWidth > maxWidth {
		newHeight = uint(origHeight * maxWidth / origWidth)
		if newHeight < 1 {
			newHeight = 1
		}
		newWidth = maxWidth
	}

	if newHeight > maxHeight {
		newWidth = uint(newWidth * maxHeight / newHeight)
		if newWidth < 1 {
			newWidth = 1
		}
		newHeight = maxHeight
	}
	return p.ResizeImage(int(newWidth), int(newHeight), img, interp)
}

func resizeNearest(width, height uint, scaleX, scaleY float64, img image.Image, interp InterpolationFunction) image.Image {
	taps, _ := interp.kernel()
	cpus := runtime.NumCPU()
	wg := sync.WaitGroup{}

	switch input := img.(type) {
	case *image.RGBA:
		// 8-bit precision
		temp := image.NewRGBA(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				nearestRGBA(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				nearestRGBA(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.YCbCr:
		// 8-bit precision
		// accessing the YCbCr arrays in a tight loop is slow.
		// converting the image to ycc increases performance by 2x.
		temp := newYCC(image.Rect(0, 0, input.Bounds().Dy(), int(width)), input.SubsampleRatio)
		result := newYCC(image.Rect(0, 0, int(width), int(height)), input.SubsampleRatio)

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		in := imageYCbCrToYCC(input)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				nearestYCbCr(in, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				nearestYCbCr(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result.YCbCr()
	case *image.RGBA64:
		// 16-bit precision
		temp := image.NewRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestRGBA64(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestGeneric(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray:
		// 8-bit precision
		temp := image.NewGray(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				nearestGray(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				nearestGray(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray16:
		// 16-bit precision
		temp := image.NewGray16(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray16(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				nearestGray16(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				nearestGray16(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	default:
		// 16-bit precision
		temp := image.NewRGBA64(image.Rect(0, 0, img.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestGeneric(img, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	}

}

func (p *ImageTK) ResizeImage(widthA, heightA int, img image.Image, interpA ...InterpolationFunction) image.Image {
	width := uint(widthA)
	height := uint(heightA)
	scaleX, scaleY := calcFactors(width, height, float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	if width == 0 {
		width = uint(0.7 + float64(img.Bounds().Dx())/scaleX)
	}
	if height == 0 {
		height = uint(0.7 + float64(img.Bounds().Dy())/scaleY)
	}

	var interp InterpolationFunction

	if interpA == nil || len(interpA) < 1 {
		interp = Lanczos3
	} else {
		interp = interpA[0]
	}

	if interp == NearestNeighbor {
		return resizeNearest(width, height, scaleX, scaleY, img, interp)
	}

	taps, kernel := interp.kernel()
	cpus := runtime.NumCPU()
	wg := sync.WaitGroup{}

	// Generic access to image.Image is slow in tight loops.
	// The optimal access has to be determined from the concrete image type.
	switch input := img.(type) {
	case *image.RGBA:
		// 8-bit precision
		temp := image.NewRGBA(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				resizeRGBA(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				resizeRGBA(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.YCbCr:
		// 8-bit precision
		// accessing the YCbCr arrays in a tight loop is slow.
		// converting the image to ycc increases performance by 2x.
		temp := newYCC(image.Rect(0, 0, input.Bounds().Dy(), int(width)), input.SubsampleRatio)
		result := newYCC(image.Rect(0, 0, int(width), int(height)), input.SubsampleRatio)

		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		in := imageYCbCrToYCC(input)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				resizeYCbCr(in, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				resizeYCbCr(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result.YCbCr()
	case *image.RGBA64:
		// 16-bit precision
		temp := image.NewRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeRGBA64(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeGeneric(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray:
		// 8-bit precision
		temp := image.NewGray(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				resizeGray(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				resizeGray(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray16:
		// 16-bit precision
		temp := image.NewGray16(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray16(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				resizeGray16(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				resizeGray16(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	default:
		// 16-bit precision
		temp := image.NewRGBA64(image.Rect(0, 0, img.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		// horizontal filter, results in transposed temporary image
		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeGeneric(img, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		// horizontal filter on transposed image, result is not transposed
		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	}
}

type InterpolationFunction int

// InterpolationFunction constants
const (
	// Nearest-neighbor interpolation
	NearestNeighbor InterpolationFunction = iota
	// Bilinear interpolation
	Bilinear
	// Bicubic interpolation (with cubic hermite spline)
	Bicubic
	// Mitchell-Netravali interpolation
	MitchellNetravali
	// Lanczos interpolation (a=2)
	Lanczos2
	// Lanczos interpolation (a=3)
	Lanczos3
)

func nearest(in float64) float64 {
	if in >= -0.5 && in < 0.5 {
		return 1
	}
	return 0
}

func linear(in float64) float64 {
	in = math.Abs(in)
	if in <= 1 {
		return 1 - in
	}
	return 0
}

func cubic(in float64) float64 {
	in = math.Abs(in)
	if in <= 1 {
		return in*in*(1.5*in-2.5) + 1.0
	}
	if in <= 2 {
		return in*(in*(2.5-0.5*in)-4.0) + 2.0
	}
	return 0
}

func mitchellnetravali(in float64) float64 {
	in = math.Abs(in)
	if in <= 1 {
		return (7.0*in*in*in - 12.0*in*in + 5.33333333333) * 0.16666666666
	}
	if in <= 2 {
		return (-2.33333333333*in*in*in + 12.0*in*in - 20.0*in + 10.6666666667) * 0.16666666666
	}
	return 0
}

func sinc(x float64) float64 {
	x = math.Abs(x) * math.Pi
	if x >= 1.220703e-4 {
		return math.Sin(x) / x
	}
	return 1
}

func lanczos2(in float64) float64 {
	if in > -2 && in < 2 {
		return sinc(in) * sinc(in*0.5)
	}
	return 0
}

func lanczos3(in float64) float64 {
	if in > -3 && in < 3 {
		return sinc(in) * sinc(in*0.3333333333333333)
	}
	return 0
}

func createWeights8(dy, filterLength int, blur, scale float64, kernel func(float64) float64) ([]int16, []int, int) {
	filterLength = filterLength * int(math.Max(math.Ceil(blur*scale), 1))
	filterFactor := math.Min(1./(blur*scale), 1)

	coeffs := make([]int16, dy*filterLength)
	start := make([]int, dy)
	for y := 0; y < dy; y++ {
		interpX := scale * (float64(y) + 0.5)
		start[y] = int(interpX) - filterLength/2 + 1
		interpX -= float64(start[y])
		for i := 0; i < filterLength; i++ {
			in := (interpX - float64(i)) * filterFactor
			coeffs[y*filterLength+i] = int16(kernel(in) * 256)
		}
	}

	return coeffs, start, filterLength
}

// range [-65536,65536]
func createWeights16(dy, filterLength int, blur, scale float64, kernel func(float64) float64) ([]int32, []int, int) {
	filterLength = filterLength * int(math.Max(math.Ceil(blur*scale), 1))
	filterFactor := math.Min(1./(blur*scale), 1)

	coeffs := make([]int32, dy*filterLength)
	start := make([]int, dy)
	for y := 0; y < dy; y++ {
		interpX := scale * (float64(y) + 0.5)
		start[y] = int(interpX) - filterLength/2 + 1
		interpX -= float64(start[y])
		for i := 0; i < filterLength; i++ {
			in := (interpX - float64(i)) * filterFactor
			coeffs[y*filterLength+i] = int32(kernel(in) * 65536)
		}
	}

	return coeffs, start, filterLength
}

func createWeightsNearest(dy, filterLength int, blur, scale float64) ([]bool, []int, int) {
	filterLength = filterLength * int(math.Max(math.Ceil(blur*scale), 1))
	filterFactor := math.Min(1./(blur*scale), 1)

	coeffs := make([]bool, dy*filterLength)
	start := make([]int, dy)
	for y := 0; y < dy; y++ {
		interpX := scale * (float64(y) + 0.5)
		start[y] = int(interpX) - filterLength/2 + 1
		interpX -= float64(start[y])
		for i := 0; i < filterLength; i++ {
			in := (interpX - float64(i)) * filterFactor
			if in >= -0.5 && in < 0.5 {
				coeffs[y*filterLength+i] = true
			} else {
				coeffs[y*filterLength+i] = false
			}
		}
	}

	return coeffs, start, filterLength
}

// kernal, returns an InterpolationFunctions taps and kernel.
func (i InterpolationFunction) kernel() (int, func(float64) float64) {
	switch i {
	case Bilinear:
		return 2, linear
	case Bicubic:
		return 4, cubic
	case MitchellNetravali:
		return 4, mitchellnetravali
	case Lanczos2:
		return 4, lanczos2
	case Lanczos3:
		return 6, lanczos3
	default:
		// Default to NearestNeighbor.
		return 2, nearest
	}
}

// values <1 will sharpen the image
var blur = 1.0

// Calculates scaling factors using old and new image dimensions.
func calcFactors(width, height uint, oldWidth, oldHeight float64) (scaleX, scaleY float64) {
	if width == 0 {
		if height == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = oldHeight / float64(height)
			scaleX = scaleY
		}
	} else {
		scaleX = oldWidth / float64(width)
		if height == 0 {
			scaleY = scaleX
		} else {
			scaleY = oldHeight / float64(height)
		}
	}
	return
}

type imageWithSubImage interface {
	image.Image
	SubImage(image.Rectangle) image.Image
}

func makeSlice(img imageWithSubImage, i, n int) image.Image {
	return img.SubImage(image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y+i*img.Bounds().Dy()/n, img.Bounds().Max.X, img.Bounds().Min.Y+(i+1)*img.Bounds().Dy()/n))
}

func floatToUint8(x float32) uint8 {
	// Nearest-neighbor values are always
	// positive no need to check lower-bound.
	if x > 0xfe {
		return 0xff
	}
	return uint8(x)
}

func floatToUint16(x float32) uint16 {
	if x > 0xfffe {
		return 0xffff
	}
	return uint16(x)
}

func nearestGeneric(in image.Image, out *image.RGBA64, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					r, g, b, a := in.At(xi+in.Bounds().Min.X, x+in.Bounds().Min.Y).RGBA()
					rgba[0] += float32(r)
					rgba[1] += float32(g)
					rgba[2] += float32(b)
					rgba[3] += float32(a)
					sum++
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := floatToUint16(rgba[0] / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
			value = floatToUint16(rgba[1] / sum)
			out.Pix[offset+2] = uint8(value >> 8)
			out.Pix[offset+3] = uint8(value)
			value = floatToUint16(rgba[2] / sum)
			out.Pix[offset+4] = uint8(value >> 8)
			out.Pix[offset+5] = uint8(value)
			value = floatToUint16(rgba[3] / sum)
			out.Pix[offset+6] = uint8(value >> 8)
			out.Pix[offset+7] = uint8(value)
		}
	}
}

func nearestRGBA(in *image.RGBA, out *image.RGBA, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 4 * maxX
					default:
						xi *= 4
					}
					rgba[0] += float32(row[xi+0])
					rgba[1] += float32(row[xi+1])
					rgba[2] += float32(row[xi+2])
					rgba[3] += float32(row[xi+3])
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*4
			out.Pix[xo+0] = floatToUint8(rgba[0] / sum)
			out.Pix[xo+1] = floatToUint8(rgba[1] / sum)
			out.Pix[xo+2] = floatToUint8(rgba[2] / sum)
			out.Pix[xo+3] = floatToUint8(rgba[3] / sum)
		}
	}
}

func nearestRGBA64(in *image.RGBA64, out *image.RGBA64, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 8 * maxX
					default:
						xi *= 8
					}
					rgba[0] += float32(uint16(row[xi+0])<<8 | uint16(row[xi+1]))
					rgba[1] += float32(uint16(row[xi+2])<<8 | uint16(row[xi+3]))
					rgba[2] += float32(uint16(row[xi+4])<<8 | uint16(row[xi+5]))
					rgba[3] += float32(uint16(row[xi+6])<<8 | uint16(row[xi+7]))
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := floatToUint16(rgba[0] / sum)
			out.Pix[xo+0] = uint8(value >> 8)
			out.Pix[xo+1] = uint8(value)
			value = floatToUint16(rgba[1] / sum)
			out.Pix[xo+2] = uint8(value >> 8)
			out.Pix[xo+3] = uint8(value)
			value = floatToUint16(rgba[2] / sum)
			out.Pix[xo+4] = uint8(value >> 8)
			out.Pix[xo+5] = uint8(value)
			value = floatToUint16(rgba[3] / sum)
			out.Pix[xo+6] = uint8(value >> 8)
			out.Pix[xo+7] = uint8(value)
		}
	}
}

func nearestGray(in *image.Gray, out *image.Gray, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					gray += float32(row[xi])
					sum++
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x - newBounds.Min.X)
			out.Pix[offset] = floatToUint8(gray / sum)
		}
	}
}

// ycc is an in memory YCbCr image.  The Y, Cb and Cr samples are held in a
// single slice to increase resizing performance.
type ycc struct {
	// Pix holds the image's pixels, in Y, Cb, Cr order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*3].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// SubsampleRatio is the subsample ratio of the original YCbCr image.
	SubsampleRatio image.YCbCrSubsampleRatio
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *ycc) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *ycc) Bounds() image.Rectangle {
	return p.Rect
}

func (p *ycc) ColorModel() color.Model {
	return color.YCbCrModel
}

func (p *ycc) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.YCbCr{}
	}
	i := p.PixOffset(x, y)
	return color.YCbCr{
		p.Pix[i+0],
		p.Pix[i+1],
		p.Pix[i+2],
	}
}

func (p *ycc) Opaque() bool {
	return true
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ycc) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &ycc{SubsampleRatio: p.SubsampleRatio}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &ycc{
		Pix:            p.Pix[i:],
		Stride:         p.Stride,
		Rect:           r,
		SubsampleRatio: p.SubsampleRatio,
	}
}

// newYCC returns a new ycc with the given bounds and subsample ratio.
func newYCC(r image.Rectangle, s image.YCbCrSubsampleRatio) *ycc {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 3*w*h)
	return &ycc{Pix: buf, Stride: 3 * w, Rect: r, SubsampleRatio: s}
}

// YCbCr converts ycc to a YCbCr image with the same subsample ratio
// as the YCbCr image that ycc was generated from.
func (p *ycc) YCbCr() *image.YCbCr {
	ycbcr := image.NewYCbCr(p.Rect, p.SubsampleRatio)
	var off int

	switch ycbcr.SubsampleRatio {
	case image.YCbCrSubsampleRatio422:
		for y := ycbcr.Rect.Min.Y; y < ycbcr.Rect.Max.Y; y++ {
			yy := (y - ycbcr.Rect.Min.Y) * ycbcr.YStride
			cy := (y - ycbcr.Rect.Min.Y) * ycbcr.CStride
			for x := ycbcr.Rect.Min.X; x < ycbcr.Rect.Max.X; x++ {
				xx := (x - ycbcr.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx/2
				ycbcr.Y[yi] = p.Pix[off+0]
				ycbcr.Cb[ci] = p.Pix[off+1]
				ycbcr.Cr[ci] = p.Pix[off+2]
				off += 3
			}
		}
	case image.YCbCrSubsampleRatio420:
		for y := ycbcr.Rect.Min.Y; y < ycbcr.Rect.Max.Y; y++ {
			yy := (y - ycbcr.Rect.Min.Y) * ycbcr.YStride
			cy := (y/2 - ycbcr.Rect.Min.Y/2) * ycbcr.CStride
			for x := ycbcr.Rect.Min.X; x < ycbcr.Rect.Max.X; x++ {
				xx := (x - ycbcr.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx/2
				ycbcr.Y[yi] = p.Pix[off+0]
				ycbcr.Cb[ci] = p.Pix[off+1]
				ycbcr.Cr[ci] = p.Pix[off+2]
				off += 3
			}
		}
	case image.YCbCrSubsampleRatio440:
		for y := ycbcr.Rect.Min.Y; y < ycbcr.Rect.Max.Y; y++ {
			yy := (y - ycbcr.Rect.Min.Y) * ycbcr.YStride
			cy := (y/2 - ycbcr.Rect.Min.Y/2) * ycbcr.CStride
			for x := ycbcr.Rect.Min.X; x < ycbcr.Rect.Max.X; x++ {
				xx := (x - ycbcr.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx
				ycbcr.Y[yi] = p.Pix[off+0]
				ycbcr.Cb[ci] = p.Pix[off+1]
				ycbcr.Cr[ci] = p.Pix[off+2]
				off += 3
			}
		}
	default:
		// Default to 4:4:4 subsampling.
		for y := ycbcr.Rect.Min.Y; y < ycbcr.Rect.Max.Y; y++ {
			yy := (y - ycbcr.Rect.Min.Y) * ycbcr.YStride
			cy := (y - ycbcr.Rect.Min.Y) * ycbcr.CStride
			for x := ycbcr.Rect.Min.X; x < ycbcr.Rect.Max.X; x++ {
				xx := (x - ycbcr.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx
				ycbcr.Y[yi] = p.Pix[off+0]
				ycbcr.Cb[ci] = p.Pix[off+1]
				ycbcr.Cr[ci] = p.Pix[off+2]
				off += 3
			}
		}
	}
	return ycbcr
}

// imageYCbCrToYCC converts a YCbCr image to a ycc image for resizing.
func imageYCbCrToYCC(in *image.YCbCr) *ycc {
	w, h := in.Rect.Dx(), in.Rect.Dy()
	r := image.Rect(0, 0, w, h)
	buf := make([]uint8, 3*w*h)
	p := ycc{Pix: buf, Stride: 3 * w, Rect: r, SubsampleRatio: in.SubsampleRatio}
	var off int

	switch in.SubsampleRatio {
	case image.YCbCrSubsampleRatio422:
		for y := in.Rect.Min.Y; y < in.Rect.Max.Y; y++ {
			yy := (y - in.Rect.Min.Y) * in.YStride
			cy := (y - in.Rect.Min.Y) * in.CStride
			for x := in.Rect.Min.X; x < in.Rect.Max.X; x++ {
				xx := (x - in.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx/2
				p.Pix[off+0] = in.Y[yi]
				p.Pix[off+1] = in.Cb[ci]
				p.Pix[off+2] = in.Cr[ci]
				off += 3
			}
		}
	case image.YCbCrSubsampleRatio420:
		for y := in.Rect.Min.Y; y < in.Rect.Max.Y; y++ {
			yy := (y - in.Rect.Min.Y) * in.YStride
			cy := (y/2 - in.Rect.Min.Y/2) * in.CStride
			for x := in.Rect.Min.X; x < in.Rect.Max.X; x++ {
				xx := (x - in.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx/2
				p.Pix[off+0] = in.Y[yi]
				p.Pix[off+1] = in.Cb[ci]
				p.Pix[off+2] = in.Cr[ci]
				off += 3
			}
		}
	case image.YCbCrSubsampleRatio440:
		for y := in.Rect.Min.Y; y < in.Rect.Max.Y; y++ {
			yy := (y - in.Rect.Min.Y) * in.YStride
			cy := (y/2 - in.Rect.Min.Y/2) * in.CStride
			for x := in.Rect.Min.X; x < in.Rect.Max.X; x++ {
				xx := (x - in.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx
				p.Pix[off+0] = in.Y[yi]
				p.Pix[off+1] = in.Cb[ci]
				p.Pix[off+2] = in.Cr[ci]
				off += 3
			}
		}
	default:
		// Default to 4:4:4 subsampling.
		for y := in.Rect.Min.Y; y < in.Rect.Max.Y; y++ {
			yy := (y - in.Rect.Min.Y) * in.YStride
			cy := (y - in.Rect.Min.Y) * in.CStride
			for x := in.Rect.Min.X; x < in.Rect.Max.X; x++ {
				xx := (x - in.Rect.Min.X)
				yi := yy + xx
				ci := cy + xx
				p.Pix[off+0] = in.Y[yi]
				p.Pix[off+1] = in.Cb[ci]
				p.Pix[off+2] = in.Cr[ci]
				off += 3
			}
		}
	}
	return &p
}

func nearestGray16(in *image.Gray16, out *image.Gray16, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 2 * maxX
					default:
						xi *= 2
					}
					gray += float32(uint16(row[xi+0])<<8 | uint16(row[xi+1]))
					sum++
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*2
			value := floatToUint16(gray / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
		}
	}
}

// Keep value in [0,255] range.
func clampUint8(in int32) uint8 {
	if in < 0 {
		return 0
	}
	if in > 255 {
		return 255
	}
	return uint8(in)
}

// Keep value in [0,65535] range.
func clampUint16(in int64) uint16 {
	if in < 0 {
		return 0
	}
	if in > 65535 {
		return 65535
	}
	return uint16(in)
}

func resizeGeneric(in image.Image, out *image.RGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					r, g, b, a := in.At(xi+in.Bounds().Min.X, x+in.Bounds().Min.Y).RGBA()
					rgba[0] += int64(coeff) * int64(r)
					rgba[1] += int64(coeff) * int64(g)
					rgba[2] += int64(coeff) * int64(b)
					rgba[3] += int64(coeff) * int64(a)
					sum += int64(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[offset+2] = uint8(value >> 8)
			out.Pix[offset+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[offset+4] = uint8(value >> 8)
			out.Pix[offset+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[offset+6] = uint8(value >> 8)
			out.Pix[offset+7] = uint8(value)
		}
	}
}

func resizeRGBA(in *image.RGBA, out *image.RGBA, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 4 * maxX
					default:
						xi *= 4
					}
					rgba[0] += int32(coeff) * int32(row[xi+0])
					rgba[1] += int32(coeff) * int32(row[xi+1])
					rgba[2] += int32(coeff) * int32(row[xi+2])
					rgba[3] += int32(coeff) * int32(row[xi+3])
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*4
			out.Pix[xo+0] = clampUint8(rgba[0] / sum)
			out.Pix[xo+1] = clampUint8(rgba[1] / sum)
			out.Pix[xo+2] = clampUint8(rgba[2] / sum)
			out.Pix[xo+3] = clampUint8(rgba[3] / sum)
		}
	}
}

func resizeRGBA64(in *image.RGBA64, out *image.RGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 8 * maxX
					default:
						xi *= 8
					}
					rgba[0] += int64(coeff) * int64(uint16(row[xi+0])<<8|uint16(row[xi+1]))
					rgba[1] += int64(coeff) * int64(uint16(row[xi+2])<<8|uint16(row[xi+3]))
					rgba[2] += int64(coeff) * int64(uint16(row[xi+4])<<8|uint16(row[xi+5]))
					rgba[3] += int64(coeff) * int64(uint16(row[xi+6])<<8|uint16(row[xi+7]))
					sum += int64(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[xo+0] = uint8(value >> 8)
			out.Pix[xo+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[xo+2] = uint8(value >> 8)
			out.Pix[xo+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[xo+4] = uint8(value >> 8)
			out.Pix[xo+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[xo+6] = uint8(value >> 8)
			out.Pix[xo+7] = uint8(value)
		}
	}
}

func resizeGray(in *image.Gray, out *image.Gray, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-newBounds.Min.X)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					gray += int32(coeff) * int32(row[xi])
					sum += int32(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x - newBounds.Min.X)
			out.Pix[offset] = clampUint8(gray / sum)
		}
	}
}

func resizeGray16(in *image.Gray16, out *image.Gray16, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 2 * maxX
					default:
						xi *= 2
					}
					gray += int64(coeff) * int64(uint16(row[xi+0])<<8|uint16(row[xi+1]))
					sum += int64(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*2
			value := clampUint16(gray / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
		}
	}
}

func resizeYCbCr(in *ycc, out *ycc, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var p [3]int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 3 * maxX
					default:
						xi *= 3
					}
					p[0] += int32(coeff) * int32(row[xi+0])
					p[1] += int32(coeff) * int32(row[xi+1])
					p[2] += int32(coeff) * int32(row[xi+2])
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*3
			out.Pix[xo+0] = clampUint8(p[0] / sum)
			out.Pix[xo+1] = clampUint8(p[1] / sum)
			out.Pix[xo+2] = clampUint8(p[2] / sum)
		}
	}
}

func nearestYCbCr(in *ycc, out *ycc, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var p [3]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = 3 * maxX
					default:
						xi *= 3
					}
					p[0] += float32(row[xi+0])
					p[1] += float32(row[xi+1])
					p[2] += float32(row[xi+2])
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*3
			out.Pix[xo+0] = floatToUint8(p[0] / sum)
			out.Pix[xo+1] = floatToUint8(p[1] / sum)
			out.Pix[xo+2] = floatToUint8(p[2] / sum)
		}
	}
}

// modified based on github.com/pokemium/hq2xgo, thanks

func interp1(a, b color.RGBA) color.RGBA {
	f := func(a, b uint8) uint8 {
		return uint8((uint(a)*3 + uint(b)) / 4)
	}

	R := f(a.R, b.R)
	G := f(a.G, b.G)
	B := f(a.B, b.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp2(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*2 + uint(b) + uint(c)) / 4)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp5(a, b color.RGBA) color.RGBA {
	f := func(a, b uint8) uint8 {
		return uint8((uint(a) + uint(b)) / 2)
	}

	R := f(a.R, b.R)
	G := f(a.G, b.G)
	B := f(a.B, b.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp6(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*5 + uint(b)*2 + uint(c)) / 8)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp7(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*6 + uint(b) + uint(c)) / 8)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp9(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*2 + uint(b)*3 + uint(c)*3) / 8)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp10(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*14 + uint(b) + uint(c)) / 16)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

const (
	TOP_LEFT = iota
	TOP
	TOP_RIGHT
	LEFT
	CENTER
	RIGHT
	BOTTOM_LEFT
	BOTTOM
	BOTTOM_RIGHT
)

var (
	contextFlag [9]uint8
)

func (p *ImageTK) EnlargeImage(src image.Image, scaleA float64) (image.Image, error) {
	srcX, srcY := src.Bounds().Dx(), src.Bounds().Dy()

	timesT := int(math.Ceil(math.Sqrt(scaleA)))

	destT, errT := p.LoadRGBAFromImage(src)
	if errT != nil {
		return nil, errT
	}

	for i := 0; i < timesT; i++ {
		destT, errT = p.HQ2x(destT)

		if errT != nil {
			return nil, errT
		}
	}

	w, h := destT.Bounds().Dx(), destT.Bounds().Dy()

	nw, nh := int(float64(srcX)*scaleA), int(float64(srcY)*scaleA)

	if (nw != w) && (nh != h) {
		newImageT := p.ResizeImage(nw, nh, destT)

		return newImageT, nil
	}

	return destT, nil

}

// HQ2x - Enlarge image by 2x with hq2x algorithm
func (p *ImageTK) HQ2x(src *image.RGBA) (*image.RGBA, error) {
	initContextFlag()
	srcX, srcY := src.Bounds().Dx(), src.Bounds().Dy()

	dest := image.NewRGBA(image.Rect(0, 0, srcX*2, srcY*2))

	columns := make(chan int, srcX)
	for x := 0; x < srcX; x++ {
		columns <- x
	}

	var wg sync.WaitGroup
	wg.Add(srcX)
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(i, src, dest, columns, &wg)
	}
	close(columns)
	wg.Wait()

	return dest, nil
}

func worker(id int, src, dest *image.RGBA, columns chan int, wg *sync.WaitGroup) {
	for column := range columns {
		hq2xColumn(src, dest, column)
		wg.Done()
	}
}

func workerx(id int, src, dest *image.RGBA, columns chan int, scaleA int, wg *sync.WaitGroup) {
	for column := range columns {
		hq2xColumnx(src, dest, column, scaleA)
		wg.Done()
	}
}

// x列目に対してhq2xアルゴリズムによる拡大処理
func hq2xColumn(src, dest *image.RGBA, x int) {
	srcY := src.Bounds().Dy()
	for y := 0; y < srcY; y++ {
		tl, tr, bl, br := hq2xPixel(src, x, y)
		tl.A, tr.A, bl.A, br.A = 0xff, 0xff, 0xff, 0xff
		dest.Set(x*2, y*2, tl)
		dest.Set(x*2+1, y*2, tr)
		dest.Set(x*2, y*2+1, bl)
		dest.Set(x*2+1, y*2+1, br)
	}
}

func hq2xColumnx(src, dest *image.RGBA, x int, scaleA int) {
	srcY := src.Bounds().Dy()
	for y := 0; y < srcY; y++ {
		tl, tr, bl, br := hq2xPixel(src, x, y)
		tl.A, tr.A, bl.A, br.A = 0xff, 0xff, 0xff, 0xff
		dest.Set(x*scaleA, y*scaleA, tl)
		dest.Set(x*scaleA+1, y*scaleA, tr)
		dest.Set(x*scaleA, y*scaleA+1, bl)
		dest.Set(x*scaleA+1, y*scaleA+1, br)
	}
}

func getPixel(src *image.RGBA, x, y int) color.RGBA {
	width, height := src.Bounds().Dx(), src.Bounds().Dy()

	if x < 0 {
		x = 0
	} else if x >= width {
		x = width - 1
	}

	if y < 0 {
		y = 0
	} else if y >= height {
		y = height - 1
	}

	return src.RGBAAt(x, y)
}

func hq2xPixel(src *image.RGBA, x, y int) (tl, tr, bl, br color.RGBA) {

	context := [9]color.RGBA{
		getPixel(src, x-1, y-1), getPixel(src, x, y-1), getPixel(src, x+1, y-1),
		getPixel(src, x-1, y), getPixel(src, x, y), getPixel(src, x+1, y),
		getPixel(src, x-1, y+1), getPixel(src, x, y+1), getPixel(src, x+1, y+1),
	}

	yuvContext := [9]color.YCbCr{}
	yuvPixel := rgbaToYCbCr(context[CENTER])
	for i := 0; i < 9; i++ {
		yuvContext[i] = rgbaToYCbCr(context[i])
	}

	var pattern uint8
	for bit := 0; bit < 9; bit++ {
		if bit != CENTER && !equalYuv(yuvContext[bit], yuvPixel) {
			pattern |= contextFlag[bit]
		}
	}

	switch pattern {
	case 0, 1, 4, 32, 128, 5, 132, 160, 33, 129, 36, 133, 164, 161, 37, 165:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 2, 34, 130, 162:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 16, 17, 48, 49:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 64, 65, 68, 69:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 8, 12, 136, 140:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 3, 35, 131, 163:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 6, 38, 134, 166:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 20, 21, 52, 53:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 144, 145, 176, 177:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM])

	case 192, 193, 196, 197:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 96, 97, 100, 101:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 40, 44, 168, 172:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 9, 13, 137, 141:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 18, 50:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 80, 81:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 72, 76:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 10, 138:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 66:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 24:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 7, 39, 135:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 148, 149, 180:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM])

	case 224, 228, 225:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 41, 169, 45:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 22, 54:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 208, 209:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 104, 108:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 11, 139:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 19, 51:
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tl = interp1(context[CENTER], context[LEFT])
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tl = interp6(context[CENTER], context[TOP], context[LEFT])
			tr = interp9(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 146, 178:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
			br = interp1(context[CENTER], context[BOTTOM])
		} else {
			tr = interp9(context[CENTER], context[TOP], context[RIGHT])
			br = interp6(context[CENTER], context[RIGHT], context[BOTTOM])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])

	case 84, 85:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP])
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			tr = interp6(context[CENTER], context[RIGHT], context[TOP])
			br = interp9(context[CENTER], context[RIGHT], context[BOTTOM])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])

	case 112, 113:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			bl = interp1(context[CENTER], context[LEFT])
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			bl = interp6(context[CENTER], context[BOTTOM], context[LEFT])
			br = interp9(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 200, 204:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
			br = interp1(context[CENTER], context[RIGHT])
		} else {
			bl = interp9(context[CENTER], context[BOTTOM], context[LEFT])
			br = interp6(context[CENTER], context[BOTTOM], context[RIGHT])
		}

	case 73, 77:
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			tl = interp1(context[CENTER], context[TOP])
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			tl = interp6(context[CENTER], context[LEFT], context[TOP])
			bl = interp9(context[CENTER], context[BOTTOM], context[LEFT])
		}
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 42, 170:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
			bl = interp1(context[CENTER], context[BOTTOM])
		} else {
			tl = interp9(context[CENTER], context[LEFT], context[TOP])
			bl = interp6(context[CENTER], context[LEFT], context[BOTTOM])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 14, 142:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
			tr = interp1(context[CENTER], context[RIGHT])
		} else {
			tl = interp9(context[CENTER], context[LEFT], context[TOP])
			tr = interp6(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 67:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 70:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 28:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 152:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 194:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 98:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 56:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 25:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 26, 31:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 82, 214:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 88, 248:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 74, 107:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 27:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 86:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 216:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 106:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 30:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 210:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 120:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])
	case 75:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 29:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 198:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 184:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 99:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 57:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 71:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 156:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 226:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 60:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 195:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 102:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 153:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 58:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 83:
		tl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 92:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 202:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[RIGHT])

	case 78:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 154:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 114:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 89:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 90:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 55, 23:
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tl = interp1(context[CENTER], context[LEFT])
			tr = context[CENTER]
		} else {
			tl = interp6(context[CENTER], context[TOP], context[LEFT])
			tr = interp9(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 182, 150:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
			br = interp1(context[CENTER], context[BOTTOM])
		} else {
			tr = interp9(context[CENTER], context[TOP], context[RIGHT])
			br = interp6(context[CENTER], context[RIGHT], context[BOTTOM])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])

	case 213, 212:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			tr = interp1(context[CENTER], context[TOP])
			br = context[CENTER]
		} else {
			tr = interp6(context[CENTER], context[RIGHT], context[TOP])
			br = interp9(context[CENTER], context[RIGHT], context[BOTTOM])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])

	case 241, 240:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			bl = interp1(context[CENTER], context[LEFT])
			br = context[CENTER]
		} else {
			bl = interp6(context[CENTER], context[BOTTOM], context[LEFT])
			br = interp9(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 236, 232:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
			br = interp1(context[CENTER], context[RIGHT])
		} else {
			bl = interp9(context[CENTER], context[BOTTOM], context[LEFT])
			br = interp6(context[CENTER], context[BOTTOM], context[RIGHT])
		}

	case 109, 105:
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			tl = interp1(context[CENTER], context[TOP])
			bl = context[CENTER]
		} else {
			tl = interp6(context[CENTER], context[LEFT], context[TOP])
			bl = interp9(context[CENTER], context[BOTTOM], context[LEFT])
		}
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 171, 43:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
			bl = interp1(context[CENTER], context[BOTTOM])
		} else {
			tl = interp9(context[CENTER], context[LEFT], context[TOP])
			bl = interp6(context[CENTER], context[LEFT], context[BOTTOM])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 143, 15:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
			tr = interp1(context[CENTER], context[RIGHT])
		} else {
			tl = interp9(context[CENTER], context[LEFT], context[TOP])
			tr = interp6(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 124:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 203:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 62:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 211:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 118:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 217:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 110:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 155:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 188:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 185:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])
	case 61:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 157:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 103:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 227:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 230:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 199:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 220:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 158:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 234:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[RIGHT])

	case 242:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 59:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 121:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 87:
		tl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 79:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 122:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 94:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 218:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 91:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 229:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 167:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 173:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 181:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM])

	case 186:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 115:
		tl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 93:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 206:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[RIGHT])

	case 205, 201:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[RIGHT])

	case 174, 46:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 179, 147:
		tl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM])

	case 117, 116:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 189:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 231:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 126:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 219:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 125:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 221:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 207:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
			tr = interp1(context[CENTER], context[RIGHT])
		} else {
			tl = interp9(context[CENTER], context[LEFT], context[TOP])
			tr = interp6(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 238:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
			br = interp1(context[CENTER], context[RIGHT])
		} else {
			bl = interp9(context[CENTER], context[BOTTOM], context[LEFT])
			br = interp6(context[CENTER], context[BOTTOM], context[RIGHT])
		}

	case 190:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
			br = interp1(context[CENTER], context[BOTTOM])
		} else {
			tr = interp9(context[CENTER], context[TOP], context[RIGHT])
			br = interp6(context[CENTER], context[RIGHT], context[BOTTOM])
		}
		bl = interp1(context[CENTER], context[BOTTOM])

	case 187:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
			bl = interp1(context[CENTER], context[BOTTOM])
		} else {
			tl = interp9(context[CENTER], context[LEFT], context[TOP])
			bl = interp6(context[CENTER], context[LEFT], context[BOTTOM])
		}
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		br = interp1(context[CENTER], context[BOTTOM])

	case 243:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			bl = interp1(context[CENTER], context[LEFT])
			br = context[CENTER]
		} else {
			bl = interp6(context[CENTER], context[BOTTOM], context[LEFT])
			br = interp9(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 119:
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tl = interp1(context[CENTER], context[LEFT])
			tr = context[CENTER]
		} else {
			tl = interp6(context[CENTER], context[TOP], context[LEFT])
			tr = interp9(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 237, 233:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp10(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[RIGHT])

	case 175, 47:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp10(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 183, 151:
		tl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp10(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM])

	case 245, 244:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp10(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 250:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 123:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 95:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 222:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 252:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp10(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 249:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp10(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 235:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp10(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[RIGHT])

	case 111:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp10(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 63:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp10(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 159:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp10(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 215:
		tl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp10(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 246:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp10(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 254:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp10(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 253:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp10(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp10(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 251:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp10(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 239:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp10(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp10(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[RIGHT])

	case 127:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp10(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 191:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp10(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp10(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 223:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp10(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 247:
		tl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp10(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp10(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 255:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp10(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp10(context[CENTER], context[TOP], context[RIGHT])
		}
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp10(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp10(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	default:
		panic(fmt.Errorf("invalid pattern: %d", pattern))
	}

	return tl, tr, bl, br
}

func equalYuv(a color.YCbCr, b color.YCbCr) bool {
	const (
		yThreshhold = 48.
		uThreshhold = 7.
		vThreshhold = 6.
	)

	aY, aU, aV := a.Y, a.Cb, a.Cr
	bY, bU, bV := b.Y, b.Cb, b.Cr

	if math.Abs(float64(aY)-float64(bY)) > yThreshhold {
		return false
	}
	if math.Abs(float64(aU)-float64(bU)) > uThreshhold {
		return false
	}
	if math.Abs(float64(aV)-float64(bV)) > vThreshhold {
		return false
	}

	return true
}

func initContextFlag() {
	curFlag := uint8(1)

	for i := 0; i < 9; i++ {
		if i == CENTER {
			continue
		}

		contextFlag[i] = curFlag
		curFlag = curFlag << 1
	}
}

func rgbaToYCbCr(c color.RGBA) color.YCbCr {
	r, g, b := c.R, c.G, c.B
	y, u, v := color.RGBToYCbCr(r, g, b)
	return color.YCbCr{
		Y:  y,
		Cb: u,
		Cr: v,
	}
}
