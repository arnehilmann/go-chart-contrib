package chartcontrib_test

import "testing"

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/andybons/gogif"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"

	"github.com/arnehilmann/go-chart-contrib"

	. "github.com/arnehilmann/goutils"
)

type Timeline struct {
	name   string
	epochs []time.Time
	values []float64
}

func NewTimeline(name string) Timeline {
	return Timeline{name, []time.Time{}, []float64{}}
}

func createChart(timelines []Timeline, linespacing float64) (chart.Chart, error) {
	c := chart.Chart{
		Width:  800,
		Height: 600,
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				return chart.TimeValueFormatterWithFormat(v, "2006-01-02T15:04")
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show:        true,
				StrokeColor: drawing.Color{255, 0, 0, 255},
			},
			ValueFormatter: func(v interface{}) string {
				return chart.FloatValueFormatterWithFormat(v, "%0.1f")
			},
		},
		YAxisSecondary: chart.YAxis{
			Style: chart.Style{
				Show:        true,
				StrokeColor: drawing.Color{0, 255, 0, 255},
			},
			ValueFormatter: func(v interface{}) string {
				return chart.FloatValueFormatterWithFormat(v, "%0.1f")
			},
			Range: chartcontrib.ContinuousRangeWithTicksLinespacing(linespacing),
			Name:  fmt.Sprintf("linespacing: %.1f", linespacing),
			NameStyle: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    timelines[0].name,
				XValues: timelines[0].epochs,
				YValues: timelines[0].values,
				Style: chart.Style{
					Show:        true,
					StrokeWidth: 2.0,
					StrokeColor: drawing.Color{255, 0, 0, 255},
				},
			},
			chart.TimeSeries{
				Name:    timelines[0].name,
				XValues: timelines[0].epochs,
				YValues: timelines[0].values,
				YAxis:   chart.YAxisSecondary,
				Style: chart.Style{
					Show:        true,
					StrokeWidth: 2.0,
					StrokeColor: drawing.Color{0, 255, 0, 255},
				},
			},
		},
	}
	if err := c.Series[0].Validate(); err != nil {
		return chart.Chart{}, err
	}
	if err := c.Series[1].Validate(); err != nil {
		return chart.Chart{}, err
	}
	return c, nil
}

func TestScale(t *testing.T) {
	file, err := os.Open("testdata/temperature.dump")
	PanicIf(err)
	defer file.Close()

	var timelines []Timeline
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(timelines) == 0 {
			for _, name := range fields {
				timelines = append(timelines, NewTimeline(name))
			}
			continue
		}
		epoch, err := strconv.Atoi(strings.TrimSuffix(fields[0], ":"))
		PanicIf(err)
		for index, field := range fields[1:] {
			value, err := strconv.ParseFloat(field, 32)
			PanicIf(err)
			if math.IsNaN(value) {
				continue
			}
			timelines[index].epochs = append(timelines[index].epochs, time.Unix(int64(epoch), 0))
			timelines[index].values = append(timelines[index].values, float64(value))
		}
	}
	err = scanner.Err()
	PanicIf(err)
	for _, line := range timelines {
		log.Println(line.name)
		for i := range line.epochs {
			log.Println(line.epochs[i], line.values[i])
		}
		break
	}

	var filenames []string
	for i := 1; i <= 10; i++ {
		c, err := createChart(timelines, float64(i))
		PanicIf(err)
		log.Println("assembling chart with linespacing", i)
		filename := filepath.Join(os.TempDir(), fmt.Sprintf("go-chart-contrib-test.%d.png", i))
		filenames = append(filenames, filename)
		f, err := os.Create(filepath.Join(os.TempDir(), fmt.Sprintf("go-chart-contrib-test.%d.png", i)))
		PanicIf(err)
		defer f.Close()
		c.Render(chart.PNG, f)
		log.Println("chart can be found in", f.Name())
	}

	outGif := &gif.GIF{}
	for i := len(filenames) - 1; i >= 0; i-- {
		filename := filenames[i]
		f, err := os.Open(filename)
		PanicIf(err)
		defer f.Close()
		simage, err := png.Decode(f)
		PanicIf(err)

		bounds := simage.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 64}
		quantizer.Quantize(palettedImage, bounds, simage, image.ZP)

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 1)
	}

	// save to out.gif
	filename := filepath.Join(os.TempDir(), "go-chart-contrib-test.all.png")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	PanicIf(err)
	defer f.Close()
	gif.EncodeAll(f, outGif)
	log.Println("all charts can be found in", f.Name())
}
