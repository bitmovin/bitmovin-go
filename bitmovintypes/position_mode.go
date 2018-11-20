package bitmovintypes

type PositionMode string

const (
	Keyframe PositionMode = "KEYFRAME"
	Time     PositionMode = "TIME"
	Segment  PositionMode = "SEGMENT"
)
