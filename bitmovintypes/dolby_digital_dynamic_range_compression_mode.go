package bitmovintypes

// DolbyDigitalDynamicRangeCompressionMode : Dynamic range compression processing mode
type DolbyDigitalDynamicRangeCompressionMode string

// List of possible DolbyDigitalDynamicRangeCompressionMode values
const (
	DolbyDigitalDynamicRangeCompressionMode_NONE           DolbyDigitalDynamicRangeCompressionMode = "NONE"
	DolbyDigitalDynamicRangeCompressionMode_FILM_STANDARD  DolbyDigitalDynamicRangeCompressionMode = "FILM_STANDARD"
	DolbyDigitalDynamicRangeCompressionMode_FILM_LIGHT     DolbyDigitalDynamicRangeCompressionMode = "FILM_LIGHT"
	DolbyDigitalDynamicRangeCompressionMode_MUSIC_STANDARD DolbyDigitalDynamicRangeCompressionMode = "MUSIC_STANDARD"
	DolbyDigitalDynamicRangeCompressionMode_MUSIC_LIGHT    DolbyDigitalDynamicRangeCompressionMode = "MUSIC_LIGHT"
	DolbyDigitalDynamicRangeCompressionMode_SPEECH         DolbyDigitalDynamicRangeCompressionMode = "SPEECH"
)
