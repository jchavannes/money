package chart

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

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

type ChartDataPointSorter []*ChartDataPoint

func (c ChartDataPointSorter) Len() int {
	return len(c)
}
func (c ChartDataPointSorter) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ChartDataPointSorter) Less(i, j int) bool {
	return c[i].Timestamp.Unix() < c[j].Timestamp.Unix()
}
