package bitmovintypes

// DolbyDigitalPlusDownmixingPreferredMode : It indicates if downmixing mode is Dolby Surround compatible (`LT_RT`: Left total/Right total) or Dolby Pro Logic II (`PRO_LOGIC_II`).  `LO_RO` for Left only/Right only: A downmix from a multichannel to a two‐channel output that is compatible for stereo or mono reproduction.
type DolbyDigitalPlusDownmixingPreferredMode string

// List of possible DolbyDigitalPlusDownmixingPreferredMode values
const (
	DolbyDigitalPlusDownmixingPreferredMode_LO_RO        DolbyDigitalPlusDownmixingPreferredMode = "LO_RO"
	DolbyDigitalPlusDownmixingPreferredMode_LT_RT        DolbyDigitalPlusDownmixingPreferredMode = "LT_RT"
	DolbyDigitalPlusDownmixingPreferredMode_PRO_LOGIC_II DolbyDigitalPlusDownmixingPreferredMode = "PRO_LOGIC_II"
)
