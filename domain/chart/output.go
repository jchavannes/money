package chart

type ChartDataOutput struct {
	Title string
	Items []*ChartDataOutputItem
}

type ChartDataOutputItem struct {
	Name string
	Data [][2]float64
}
