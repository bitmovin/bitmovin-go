package models

import (
	"fmt"
)

type BitmovinError struct {
	DataEnvelope DataEnvelope
}

func (e BitmovinError) Error() string {
	// return e.Data.Message
	return fmt.Sprintf("%s %d (ReqId#%s): %s", e.DataEnvelope.Status, e.DataEnvelope.Data.Code, e.DataEnvelope.RequestID, e.DataEnvelope.Data.Message)
}
