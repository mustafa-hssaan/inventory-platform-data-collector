package analytics

type TrafficReportResponse struct {
	Warnings  []Warning             `json:"warnings,omitempty"`
	Header    HeaderParameters      `json:"header"`
	Records   []TrafficReportRecord `json:"records"`
	StartDate string                `json:"startDate"`
	EndDate   string                `json:"EndDate"`
}

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type TrafficReportRecord struct {
	Key     string   `json:"key"`
	Metrics []Metric `json:"metrics"`
}

type Warning struct {
	ErrorId    int                 `json:"errorId"`
	Domain     string              `json:"domain"`
	Category   string              `json:"category"`
	Message    string              `json:"message"`
	Parameters []WarningParameters `json:"parameters"`
}
type WarningParameters struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type HeaderParameters struct {
	DimensionKey []KeysParameters `json:"dimensionKeys"`
	Metrics      []KeysParameters `json:"metrics"`
}
type KeysParameters struct {
	Key           string `json:"key"`
	LocalizedName string `json:"localizedName"`
	DataType      string `json:"dataType"`
}
