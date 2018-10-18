package bitmovintypes

type EncodingMode string

const (
	EncodingModeStandard   EncodingMode = "STANDARD"
	EncodingModeSinglePass EncodingMode = "SINGLE_PASS"
	EncodingModeTwoPass    EncodingMode = "TWO_PASS"
	EncodingModeThreePass  EncodingMode = "THREE_PASS"
)
