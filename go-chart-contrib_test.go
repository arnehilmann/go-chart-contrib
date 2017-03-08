package chartcontrib_test

import "testing"

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func TestScale(t *testing.T) {
	out, err := exec.Command("rrdtool", "fetch", "testdata/temperature.rrd", "AVERAGE",
		"-s", "02/22/2017", "-e", "02/24/2017").Output()
	PanicIf(err)

	var timelines []Timeline
	for _, line := range strings.Split(string(out), "\n") {
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
	for _, line := range timelines {
		log.Println(line.name)
		for i := range line.epochs {
			log.Println(line.epochs[i], line.values[i])
		}
		break
	}

	log.Println("assembling graph")

	graph := chart.Chart{
		Width:  600,
		Height: 400,
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
			Range: chartcontrib.ContinuousRangeWithTicks(5),
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
	err = graph.Series[0].Validate()
	PanicIf(err)
	err = graph.Series[1].Validate()
	PanicIf(err)

	log.Println("tmp dir", os.TempDir())
	prefix := "go-chart-contrib-test-"
	old_charts, err := filepath.Glob(os.TempDir() + prefix + "*")
	PanicIf(err)
	for _, old_chart := range old_charts {
		log.Println("old", old_chart)
		err := os.Remove(old_chart)
		PanicIf(err)
	}

	f, err := ioutil.TempFile("", prefix)
	PanicIf(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
	log.Println("graph can be found in", f.Name())
}
