package chartcontrib

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"log"
	"math"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type MyRange struct {
	*chart.ContinuousRange
	count       int
	linespacing float64
}

func (r MyRange) GetTicks(re chart.Renderer, defaults chart.Style, vf chart.ValueFormatter) []chart.Tick {
	log.Println("MyRange.GetTicks called")
	log.Println(r.GetDelta())
	log.Println(r.GetDomain())
	log.Println(r.GetMin())
	log.Println(r.GetMax())

	count := r.count
	if count == 0 {
		log.Println("count is zero")
		log.Println("domain", r.GetDomain())
		log.Println("font", defaults.GetFont().Name(truetype.NameIDFontFullName))
		log.Println("size", defaults.GetFontSize())
		log.Println("bounds", defaults.GetFont().Bounds(fixed.Int26_6(defaults.GetFontSize())))
		font_height := float64(defaults.GetFont().Bounds(fixed.Int26_6(defaults.GetFontSize())).Max.Y)
		log.Println("font height", font_height)

		log.Println(float64(r.GetDomain()) / float64(font_height))
		log.Println(float64(r.GetDomain()) / float64(font_height) / r.linespacing)
		count = int(float64(r.GetDomain()) / float64(font_height) / r.linespacing)
		log.Println("count", count)

		extents := drawing.Extents(defaults.GetFont(), defaults.GetFontSize())
		log.Println("ascent", extents.Ascent)
		log.Println("descent", extents.Descent)
		log.Println("height", extents.Height)

		count = int(float64(r.GetDomain()) / (float64(extents.Height) * r.linespacing))
		log.Println("count", count)
	}

	steplength := r.GetDelta() / float64(count)
	log.Println(steplength)

	factor := math.Pow10(int(math.Log10(steplength)))
	log.Println(factor)

	normlength := steplength / factor
	log.Println(normlength)

	var minindex float64
	mindiff := math.Inf(1)
	for _, steptry := range []float64{.1, .2, .5, 1, 2, 2.5, 5, 10} {
		diff := math.Abs(normlength - steptry)
		log.Println("diff", steptry, diff)
		if diff < mindiff {
			mindiff = diff
			minindex = steptry
		}
	}
	log.Println("min", minindex)

	newsteplength := float64(minindex) * factor
	log.Println(newsteplength)

	min := r.GetMin()
	max := r.GetMax()

	ticks := []chart.Tick{}
	for actual := chart.Math.RoundUp(min, newsteplength); actual <= max; actual += newsteplength {
		value := float64(actual)
		ticks = append(ticks, chart.Tick{Value: value, Label: vf(value)})
	}

	return ticks
}

func ContinuousRangeWithTicksLinespacing(linespacing float64) MyRange {
	return MyRange{ContinuousRange: &chart.ContinuousRange{}, linespacing: linespacing}
}

func ContinuousRangeWithTicksCount(count int) MyRange {
	return MyRange{ContinuousRange: &chart.ContinuousRange{}, count: count}
}
