package bitmovintypes

type H264AQMode string

const (
	H264AQModeDisabled               H264AQMode = "DISABLED"
	H264AQModeVariance               H264AQMode = "VARIANCE"
	H264AQModeAutoVariance           H264AQMode = "AUTO_VARIANCE"
	H264AQModeAutoVarianceDarkScenes H264AQMode = "AUTO_VARIANCE_DARK_SCENES"
)
