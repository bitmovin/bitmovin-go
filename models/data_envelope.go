package models

type DataEnvelope struct {
	RequestID string           `json:"requestId"`
	Status    string           `json:"status"`
	Data      DataEnvelopeData `json:"data"`
}

type DataEnvelopeData struct {
	Code             int                       `json:"code"`
	Message          string                    `json:"message"`
	DeveloperMessage string                    `json:"developerMessage"`
	Links            []DataEnvelopeDataLinks   `json:"links"`
	Details          []DataEnvelopeDataDetails `json:"details"`
}

type DataEnvelopeDataLinks struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type DataEnvelopeDataDetails struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Field string `json:"field"`
}
