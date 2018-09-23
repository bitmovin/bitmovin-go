package bitmovintypes

type MotionEstimationMethod string

const (
	MotionEstimationMethodDIA     = `DIA`
	MotionEstimationMethodHEX     = `HEX`
	MotionEstimationMethodUMH     = `UMH`
	MotionEstimationMethodESA     = `ESA`
	MotionEstimationMethodTESA    = `TESA`
	MotionEstimationMethodDefault = MotionEstimationMethodUMH
)
