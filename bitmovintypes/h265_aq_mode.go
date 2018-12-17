package bitmovintypes

type H265AQMode string

const (
	H265AQModeDisabled               H265AQMode = "DISABLED"
	H265AQModeVariance               H265AQMode = "VARIANCE"
	H265AQModeAutoVariance           H265AQMode = "AUTO_VARIANCE"
	H265AQModeAutoVarianceDarkScenes H265AQMode = "AUTO_VARIANCE_DARK_SCENES"
)
