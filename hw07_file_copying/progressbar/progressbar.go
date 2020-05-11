package progressbar

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/PavelVershinin/home_work/hw07_file_copying/progressbar/drawers"
)

type Drawer interface {
	io.Closer
	Draw(percent float64, text string) error
}

type ProgressBar struct {
	drawer Drawer
	min    float64
	max    float64
	val    float64
}

func New() *ProgressBar {
	p := &ProgressBar{}

	return p
}

func (p *ProgressBar) Min(n float64) *ProgressBar {
	p.min = n
	return p
}

func (p *ProgressBar) Max(n float64) *ProgressBar {
	p.max = n
	return p
}

func (p *ProgressBar) Val(n float64) *ProgressBar {
	p.val = n
	return p
}

func (p *ProgressBar) Add(n float64) *ProgressBar {
	p.val += n
	return p
}

func (p *ProgressBar) Drawer(drawer Drawer) *ProgressBar {
	p.drawer = drawer
	return p
}

//nolint:gomnd
func (p *ProgressBar) Percent() float64 {
	max := p.max - p.min
	val := p.val - p.min
	percent := val / (max / 100)

	return math.Min(100, percent)
}

func (p *ProgressBar) Left() float64 {
	return p.max - p.val
}

func (p *ProgressBar) Draw(format string) error {
	percent := p.Percent()
	format = strings.ReplaceAll(format, ":percent", fmt.Sprintf("%.2f", percent))
	format = strings.ReplaceAll(format, ":left", fmt.Sprintf("%.2f", p.Left()))
	format = strings.ReplaceAll(format, ":val", fmt.Sprintf("%.2f", p.val))
	format = strings.ReplaceAll(format, ":min", fmt.Sprintf("%.2f", p.min))
	format = strings.ReplaceAll(format, ":max", fmt.Sprintf("%.2f", p.max))

	if p.drawer == nil {
		p.drawer = &drawers.Stdout{}
	}

	return p.drawer.Draw(percent, format)
}

func (p *ProgressBar) Close() error {
	if p.drawer != nil {
		return p.drawer.Close()
	}
	return nil
}
