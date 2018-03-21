package chart

type ChartDataOutput struct {
	Title string
	Items []*ChartDataOutputItem
}

type ChartDataOutputItem struct {
	Name string
	Data [][2]float64
}

type ChartDataOutputItemsSorter [][2]float64

func (c ChartDataOutputItemsSorter) Len() int {
	return len(c)
}
func (c ChartDataOutputItemsSorter) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ChartDataOutputItemsSorter) Less(i, j int) bool {
	return c[i][0] < c[j][0]
}
