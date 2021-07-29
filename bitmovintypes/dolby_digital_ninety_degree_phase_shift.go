package bitmovintypes

// DolbyDigitalNinetyDegreePhaseShift : A 90° phase shift can be applied to the surround channels during encoding. This is useful for generating multichannel bitstreams which, when downmixed, can create a true Dolby Surround compatible output (Left/Right)
type DolbyDigitalNinetyDegreePhaseShift string

// List of possible DolbyDigitalNinetyDegreePhaseShift values
const (
	DolbyDigitalNinetyDegreePhaseShift_ENABLED  DolbyDigitalNinetyDegreePhaseShift = "ENABLED"
	DolbyDigitalNinetyDegreePhaseShift_DISABLED DolbyDigitalNinetyDegreePhaseShift = "DISABLED"
)
