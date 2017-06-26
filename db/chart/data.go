package chart

import (
	"time"
	"github.com/jchavannes/jgo/jerr"
	"sort"
)

type ChartData struct {
	Title      string
	ChartItems []*ChartItem
}

type ChartDataOutputItemsSorter [][2]float64

func (c ChartDataOutputItemsSorter) Len() int           { return len(c) }
func (c ChartDataOutputItemsSorter) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ChartDataOutputItemsSorter) Less(i, j int) bool { return c[i][0] < c[j][0] }

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

type ChartItem struct {
	Name            string
	ChartDataPoints []*ChartDataPoint
}

func (i *ChartItem) GetDataPoint(timestamp time.Time) (*ChartDataPoint, error) {
	for _, chartDataPoint := range i.ChartDataPoints {
		if chartDataPoint.Timestamp == timestamp {
			return chartDataPoint, nil
		}
	}
	return nil, jerr.New("Unable to find data point")
}

type ChartDataPoint struct {
	Timestamp time.Time
	Amount    float32
}
