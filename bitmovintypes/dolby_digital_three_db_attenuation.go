package bitmovintypes

// DolbyDigitalThreeDbAttenuation : –3 dB attenuation can be used to reduce the levels of the surround channels to compensate between the calibration of film dubbing stages and consumer replay environments. The surround channels in film studios are set 3 dB lower than the front channels (unlike consumer applications of 5.1), leading to the level on tape being 3 dB higher. Apply the 3 dB attenuation when using a master mixed in a film room.
type DolbyDigitalThreeDbAttenuation string

// List of possible DolbyDigitalThreeDbAttenuation values
const (
	DolbyDigitalThreeDbAttenuation_ENABLED  DolbyDigitalThreeDbAttenuation = "ENABLED"
	DolbyDigitalThreeDbAttenuation_DISABLED DolbyDigitalThreeDbAttenuation = "DISABLED"
)
