package imagetk

import (
	"bytes"
	"fmt"
	// "github.com/topxeq/tk"
	"gonum.org/v1/plot"
	"image"
	// "image/color"
	"gonum.org/v1/plot/vg"
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
