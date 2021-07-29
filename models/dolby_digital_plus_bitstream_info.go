package models

import "github.com/streamco/bitmovin-go/bitmovintypes"

// DolbyDigitalPlusBitstreamInfo model
type DolbyDigitalPlusBitstreamInfo struct {
	// This parameter indicates to a decoder whether the two‐channel encoded bitstream contains a Dolby Surround (Lt/Rt) program that requires Dolby Pro Logic decoding.  When downmixing to stereo from a multichannel input, set this value according to the type of downmix performed (Lt/Rt: `ENABLED`, Lo/Ro: `DISABLED`). When transcoding a stereo Dolby Digital, Dolby Digital Plus, or Dolby E input, the value must be passed through from the input bitstream to the output bitstream. When transcoding a third-party stereo input to stereo Dolby Digital set the value to `NOT_INDICATED`.
	SurroundMode bitmovintypes.DolbyDigitalPlusSurroundMode `json:"surroundMode,omitempty"`
	// This is used to identify the encoded audio as material encoded in Dolby Digital Surround EX. This parameter is used only if the encoded audio has two surround channels.  An amplifier or receiver with Dolby Digital Surround EX decoding can use this parameter as a flag to switch the decoding on or off automatically. The behavior is similar to that of the `surroundMode` parameter.
	SurroundExMode bitmovintypes.DolbyDigitalPlusSurroundMode `json:"surroundExMode,omitempty"`
}
