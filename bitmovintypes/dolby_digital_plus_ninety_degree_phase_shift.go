package bitmovintypes

// DolbyDigitalPlusNinetyDegreePhaseShift : A 90° phase shift can be applied to the surround channels during encoding. This is useful for generating multichannel bitstreams which, when downmixed, can create a true Dolby Surround compatible output (Left/Right)
type DolbyDigitalPlusNinetyDegreePhaseShift string

// List of possible DolbyDigitalPlusNinetyDegreePhaseShift values
const (
	DolbyDigitalPlusNinetyDegreePhaseShift_ENABLED  DolbyDigitalPlusNinetyDegreePhaseShift = "ENABLED"
	DolbyDigitalPlusNinetyDegreePhaseShift_DISABLED DolbyDigitalPlusNinetyDegreePhaseShift = "DISABLED"
)
