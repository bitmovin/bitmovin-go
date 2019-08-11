package bitmovintypes

type StreamMode string

const (
	StreamModeStandard                                  StreamMode = "STANDARD"
	StreamModePerTitleTemplate                          StreamMode = "PER_TITLE_TEMPLATE"
	StreamModePerTitleResult                            StreamMode = "PER_TITLE_RESULT"
	StreamModePerTitleTemplateFixedResolution           StreamMode = "PER_TITLE_TEMPLATE_FIXED_RESOLUTION"
	StreamModePerTitleTemplateFixedResolutionAndBitrate StreamMode = "PER_TITLE_TEMPLATE_FIXED_RESOLUTION_AND_BITRATE"
)
