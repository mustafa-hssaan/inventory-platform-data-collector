package analytics

type TrafficReportParam struct {
	Dimension string `url:"dimension"`
	Filter    string `url:"filter"`
	Metric    string `url:"metric"`
	Sort      string `url:"sort"`
}
