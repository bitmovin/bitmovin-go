package models

import "github.com/streamco/bitmovin-go/bitmovintypes"

// DolbyDigitalPlusLoudnessControl model
type DolbyDigitalPlusLoudnessControl struct {
	// Dialogue Normalization value to be set on the bitstream metadata. Required if the mode is `PASSTHROUGH`, or if the mode is `CORRECTION` and regulationType is `MANUAL`. For all other combinations dialnorm must not be set.
	Dialnorm *int32 `json:"dialnorm,omitempty"`
	// This may only be set if the mode is `PASSTHROUGH`, or if the mode is `CORRECTION` and regulationType is `MANUAL`. For all other regulationType values dialogueIntelligence must not be set.
	DialogueIntelligence bitmovintypes.DolbyDigitalPlusDialogueIntelligence `json:"dialogueIntelligence,omitempty"`
	Mode                 bitmovintypes.DolbyDigitalPlusLoudnessControlMode  `json:"mode,omitempty"`
	// The peak value in dB to use for loudness correction. This may only be set if the mode is `PASSTHROUGH`, or if the mode is `CORRECTION` and regulationType is `MANUAL`. For all other regulationType values peakLimit must not be set.
	PeakLimit *float64 `json:"peakLimit,omitempty"`
	// This is only allowed if the mode is CORRECTION. <table> <tr><th colspan=4 align=\"left\"> Predefined values for each regulation type: </th></tr> <tr><td> Regulation Type </td><td> EBU R128 </td><td> ATSC A/85 Fixed </td><td> ATSC A/85 Agile</td></tr> <tr><td> Limit Mode </td><td> `True Peak` </td><td> `True Peak` </td><td> `True Peak` </td></tr> <tr><td> Correction Mode </td><td> `PCM Normalization` </td><td> `PCM Normalization` </td><td> `Metadata Update` </td></tr> <tr><td> Peak Limit </td><td> `–3 dBTP` </td><td> `–2 dBTP` </td><td> `N/A` </td></tr> <tr><td> Dialogue Intelligence </td><td> `Off` </td><td> `On` </td><td> `On` </td></tr> <tr><td> Meter Mode </td><td> `ITU-R BS.1770-3` </td><td> `ITU-R BS.1770-3` </td><td> `ITU-R BS.1770-3` </td></tr> <tr><td> Speech Threshold </td><td> `20` </td><td> `20` </td><td> `20` </td></tr> <tr><td> Dialogue Normalization </td><td> `-23 dB` </td><td> `-24 dB` </td><td> `Set to measured loudness` </td></tr> </table>
	RegulationType bitmovintypes.DolbyDigitalPlusLoudnessControlRegulationType `json:"regulationType,omitempty"`
}
