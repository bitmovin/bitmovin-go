package bitmovintypes

type PresetConfiguration string

const (
	PresetConfigurationLiveHighQuality = "LIVE_HIGH_QUALITY"
	PresetConfigurationLiveLowLatency  = "LIVE_LOW_LATENCY"

	PresetConfigurationVodHighQuality    = "VOD_HIGH_QUALITY"
	PresetConfigurationVodStandard       = "VOD_STANDARD"
	PresetConfigurationVodSpeed          = "VOD_SPEED"
	PresetConfigurationVodHighSpeed      = "VOD_HIGH_SPEED"
	PresetConfigurationVodVeryHighSpeed  = "VOD_VERYHIGH_SPEED"
	PresetConfigurationVodExtraHighSpeed = "VOD_EXTRAHIGH_SPEED"
	PresetConfigurationVodSuperHighSpeed = "VOD_SUPERHIGH_SPEED"
	PresetConfigurationVodUltraHighSpeed = "VOD_ULTRAHIGH_SPEED"
)
