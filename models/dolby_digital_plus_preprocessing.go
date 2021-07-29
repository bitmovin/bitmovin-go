package models

import "github.com/streamco/bitmovin-go/bitmovintypes"

// DolbyDigitalPlusPreprocessing model
type DolbyDigitalPlusPreprocessing struct {
	// It indicates a gain change to be applied in the Dolby Digital decoder in order to implement dynamic range compression.  The values typically indicate gain reductions (cut) during loud passages and gain increases (boost) during quiet passages based on desired compression characteristics.
	DynamicRangeCompression *DolbyDigitalPlusDynamicRangeCompression `json:"dynamicRangeCompression,omitempty"`
	// It applies a 120 Hz low-pass filter to the low-frequency effects (LFE) channel.  This is only allowed if the `channelLayout` contains a LFE channel.
	LfeLowPassFilter       bitmovintypes.DolbyDigitalPlusLfeLowPassFilter       `json:"lfeLowPassFilter,omitempty"`
	NinetyDegreePhaseShift bitmovintypes.DolbyDigitalPlusNinetyDegreePhaseShift `json:"ninetyDegreePhaseShift,omitempty"`
	ThreeDbAttenuation     bitmovintypes.DolbyDigitalPlusThreeDbAttenuation     `json:"threeDbAttenuation,omitempty"`
}
