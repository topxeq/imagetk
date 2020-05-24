package imagetk

import (
	"bytes"
	"fmt"
	"github.com/topxeq/tk"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
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

func (p *ImageTK) NewNRGBAFromHex(strA string) color.NRGBA {
	r, g, b, a := tk.ParseHexColor(strA)

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (p *ImageTK) NewRGBAFromHex(strA string) color.RGBA {
	r, g, b, a := tk.ParseHexColor(strA)

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
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
