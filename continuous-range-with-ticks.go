package chartcontrib

import (
	"log"
	"math"

	"github.com/wcharczuk/go-chart"
)

type MyRange struct {
	*chart.ContinuousRange
	count int
}

func (r MyRange) GetTicks(re chart.Renderer, defaults chart.Style, vf chart.ValueFormatter) []chart.Tick {
	log.Println("MyRange.GetTicks called")
	log.Println(r.GetDelta())
	log.Println(r.GetDomain())
	log.Println(r.GetMin())
	log.Println(r.GetMax())

	steplength := r.GetDelta() / float64(r.count)
	log.Println(steplength)

	factor := math.Pow10(int(math.Log10(steplength)))
	log.Println(factor)

	normlength := steplength / factor
	log.Println(normlength)

	var minindex int
	mindiff := math.Inf(1)
	for _, steptry := range []int{1, 2, 5, 10} {
		diff := math.Abs(normlength - float64(steptry))
		log.Println(diff)
		if diff < mindiff {
			mindiff = diff
			minindex = steptry
		}
	}
	log.Println(minindex)

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

func ContinuousRangeWithTicks(count int) MyRange {
	return MyRange{ContinuousRange: &chart.ContinuousRange{}, count: count}
}
