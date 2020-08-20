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

	"github.com/topxeq/tk"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var versionG = "0.9a"

type ImageTK struct {
	Version string
}

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

func (p *ImageTK) SaveImageAs(imageA image.Image, filePathA string, formatA string) error {
	fileT, errT := os.Create(filePathA)
	if errT != nil {
		return errT
	}
	defer fileT.Close()

	formatA = tk.ToLower(formatA)

	switch formatA {
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
	return p.ResizeImage(newWidth, newHeight, img, interp)
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

func (p *ImageTK) ResizeImage(width, height uint, img image.Image, interp InterpolationFunction) image.Image {
	scaleX, scaleY := calcFactors(width, height, float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	if width == 0 {
		width = uint(0.7 + float64(img.Bounds().Dx())/scaleX)
	}
	if height == 0 {
		height = uint(0.7 + float64(img.Bounds().Dy())/scaleY)
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
