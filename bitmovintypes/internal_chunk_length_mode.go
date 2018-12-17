package bitmovintypes

type InternalChunkLengthMode string

const (
	InternalChunkLengthModeSpeedOptimized   InternalChunkLengthMode = "SPEED_OPTIMIZED"
	InternalChunkLengthModeQualityOptimized InternalChunkLengthMode = "QUALITY_OPTIMIZED"
	InternalChunkLengthModeAdaptive         InternalChunkLengthMode = "ADAPTIVE"
	InternalChunkLengthModeCustom           InternalChunkLengthMode = "CUSTOM"
)
