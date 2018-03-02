package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type KubernetesClusterConfigurationRequest struct {
	ParallelEncodings  *int64 `json:"parallelEncodings,omitempty"`
	WorkersPerEncoding *int64 `json:"workersPerEncoding,omitempty"`
}

type KubernetesClusterConfigurationDetail struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	ID                 string `json:"id"`
	ParallelEncodings  int64  `json:"parallelEncodings"`
	WorkersPerEncoding int64  `json:"workersPerEncoding"`
}

type KubernetesClusterConfigurationResponseData struct {
	Result KubernetesClusterConfigurationDetail `json:"result"`
}

type KubernetesClusterConfigurationResponse struct {
	RequestID *string                                    `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus               `json:"status,omitempty"`
	Data      KubernetesClusterConfigurationResponseData `json:"data,omitempty"`
}
