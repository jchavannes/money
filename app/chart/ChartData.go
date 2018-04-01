package chart

import (
	"sort"
)

type ChartData struct {
	Title      string
	ChartItems []*ChartItem
}

func (c *ChartData) GetChartDataOutput() *ChartDataOutput {
	chartDataOutput := &ChartDataOutput{
		Title: c.Title,
		Items: []*ChartDataOutputItem{},
	}
	for _, chartItem := range c.ChartItems {
		chartDataOutputItem := &ChartDataOutputItem{
			Name: chartItem.Name,
			Data: [][2]float64{},
		}
		for _, chartDataPoint := range chartItem.ChartDataPoints {
			outputDataPoint := [2]float64{
				float64(chartDataPoint.Timestamp.Unix() * 1000),
				float64(chartDataPoint.Amount),
			}
			chartDataOutputItem.Data = append(chartDataOutputItem.Data, outputDataPoint)
		}
		sort.Sort(ChartDataOutputItemsSorter(chartDataOutputItem.Data))
		chartDataOutput.Items = append(chartDataOutput.Items, chartDataOutputItem)
	}
	return chartDataOutput
}
