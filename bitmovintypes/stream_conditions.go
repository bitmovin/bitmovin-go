package bitmovintypes

type ConditionAttribute string

const (
	ConditionAttributeHeight  ConditionAttribute = "HEIGHT"
	ConditionAttributeWidth   ConditionAttribute = "WIDTH"
	ConditionAttributeFPS     ConditionAttribute = "FPS"
	ConditionAttributeBitrate ConditionAttribute = "BITRATE"
)

type ConditionType string

const (
	ConditionTypeAnd       ConditionType = "AND"
	ConditionTypeOr        ConditionType = "OR"
	ConditionTypeCondition ConditionType = "CONDITION"
)
