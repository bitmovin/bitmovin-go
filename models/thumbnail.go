package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type Thumbnail struct {
	Name         *string                     `json:"name,omitempty"`
	Description  *string                     `json:"description,omitempty"`
	Height       int                         `json:"height,omitempty"`
	PositionUnit *bitmovintypes.PositionUnit `json:"unit,omitempty"`
	Positions    []float64                   `json:"positions"`
	Pattern      *string                     `json:"pattern"`
	Outputs      []Output                    `json:"outputs,omitempty"`
}

func NewThumbnail(height int, positions []float64, outputs []Output) *Thumbnail {
	return &Thumbnail{
		Height:    height,
		Positions: positions,
		Outputs:   outputs,
	}
}

func (t *Thumbnail) Builder() *ThumbnailBuilder {
	return &ThumbnailBuilder{
		Thumbnail: t,
	}
}
