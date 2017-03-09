package chartcontrib

import (
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
	count := r.count
	if count == 0 {
		extents := drawing.Extents(defaults.GetFont(), defaults.GetFontSize())
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

// ContinuousRangeWithTicksLinespacing renders "nice" ticks on a YAxis, depending on the linespacing parameter.
// the actual linespacing depends on the min/max values and the height of the axis.
func ContinuousRangeWithTicksLinespacing(linespacing float64) MyRange {
	return MyRange{ContinuousRange: &chart.ContinuousRange{}, linespacing: linespacing}
}

// ContinuousRangeWithTicksCount renders "nice" ticks on a YAxis, depending on the count parameter.
// the actual tick count depends on the min/max values of the axis.
func ContinuousRangeWithTicksCount(count int) MyRange {
	return MyRange{ContinuousRange: &chart.ContinuousRange{}, count: count}
}
