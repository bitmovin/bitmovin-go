package bitmovintypes

// DolbyDigitalDialogueIntelligence : Whether to use the Dolby Dialogue Intelligence feature, which identifies and analyzes dialogue segments within audio as a basis for speech gating
type DolbyDigitalDialogueIntelligence string

// List of possible DolbyDigitalDialogueIntelligence values
const (
	DolbyDigitalDialogueIntelligence_ENABLED  DolbyDigitalDialogueIntelligence = "ENABLED"
	DolbyDigitalDialogueIntelligence_DISABLED DolbyDigitalDialogueIntelligence = "DISABLED"
)
