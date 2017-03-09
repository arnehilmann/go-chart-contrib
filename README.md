go-chart-contrib
================

contribute some missing(?) features to the
[github.com/wcharczuk/go-chart](https://github.com/wcharczuk/go-chart) library

# Installation

```bash
> go get -u github.com/arnehilmann/go-chart-contrib
```

# Features

## ContinuousRange with nice Ticks on YAxis

specify the desired linespacing, and `ContinuousRangeWithTicksLinespacing(float64)` will
create a `Range` object with nice intervalls.
These intervalls are considered as _nice_: .5, 1, 2, 2.5, 5, 10

see [some sample output with different linespacings](testdata/go-chart-contrib-test.all.gif)

## &hellip;
