package models

type Query struct {
	Dimension  string    `json:"dimension"`
	Start      string    `json:"start"`
	End        string    `json:"end"`
	Interval   *string   `json:"interval,omitempty"`
	GroupBy    []string  `json:"groupBy,omitempty"`
	Limit      *int32    `json:"limit,omitempty"`
	Offset     *int32    `json:"offset,omitempty"`
	LicenseKey *string   `json:"licenseKey,omitempty"`
	Filters    []Filter  `json:"filters,omitempty"`
	OrderBy    []OrderBy `json:"orderBy,omitempty"`
}

type Filters struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type OrderBy struct {
	Name  string `json:"name"`
	Order string `json:"order"`
}

type QueryResponse struct {
	Data      Data_  `json:"data"`
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
}

type Data_ struct {
	Result   Results    `json:"result"`
	Messages []Messages `json:"messages"`
}

type Results struct {
	Rows         []interface{}  `json:"rows"`
	RowCount     int32          `json:"rowCount"`
	CollumnLabel []collumnLabel `json:"columnLabels"`
}

type collumnLabel struct {
	Key   string `json:"key"`
	label string `json:"label"`
}

type Messages struct {
	Date  *string       `json:"date,omitempty"`
	Type  string        `json:"type"`
	Text  string        `json:"text"`
	Links []interface{} `json:"links,omitempty"`
	Field *string       `json:"field,omitempty"`
}
