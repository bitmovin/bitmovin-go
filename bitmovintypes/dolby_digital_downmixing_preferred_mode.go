package bitmovintypes

// DolbyDigitalDownmixingPreferredMode : It indicates if downmixing mode is Dolby Surround compatible (`LT_RT`: Left total/Right total) or Dolby Pro Logic II (`PRO_LOGIC_II`).  `LO_RO` for Left only/Right only: A downmix from a multichannel to a two‐channel output that is compatible for stereo or mono reproduction.
type DolbyDigitalDownmixingPreferredMode string

// List of possible DolbyDigitalDownmixingPreferredMode values
const (
	DolbyDigitalDownmixingPreferredMode_LO_RO        DolbyDigitalDownmixingPreferredMode = "LO_RO"
	DolbyDigitalDownmixingPreferredMode_LT_RT        DolbyDigitalDownmixingPreferredMode = "LT_RT"
	DolbyDigitalDownmixingPreferredMode_PRO_LOGIC_II DolbyDigitalDownmixingPreferredMode = "PRO_LOGIC_II"
)
