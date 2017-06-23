package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type AzureOutput struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	AccountName *string                `json:"accountName,omitempty"`
	AccountKey  *string                `json:"accountKey,omitempty"`
	Container   *string                `json:"container,omitempty"`
}

type FTPOutput struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Host        *string                `json:"host,omitempty"`
	UserName    *string                `json:"username,omitempty"`
	Password    *string                `json:"password,omitempty"`
	Passive     *bool                  `json:"passive,omitempty"`
}

type GCSOutput struct {
	ID          *string                         `json:"id,omitempty"`
	Name        *string                         `json:"name,omitempty"`
	Description *string                         `json:"description,omitempty"`
	CustomData  map[string]interface{}          `json:"customData,omitempty"`
	AccessKey   *string                         `json:"accessKey,omitempty"`
	SecretKey   *string                         `json:"secretKey,omitempty"`
	BucketName  *string                         `json:"bucketName,omitempty"`
	CloudRegion bitmovintypes.GoogleCloudRegion `json:"cloudRegion,omitempty"`
}

type GCSOutputItem struct {
	ID          *string                      `json:"id,omitempty"`
	Name        *string                      `json:"name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	BucketName  *string                      `json:"bucketName,omitempty"`
	CloudRegion bitmovintypes.AWSCloudRegion `json:"cloudRegion,omitempty"`
	CreatedAt   *string                      `json:"createdAt,omitempty"`
	UpdatedAt   *string                      `json:"updatedAt,omitempty"`
}

type GCSOutputData struct {
	//Success fields
	Result   GCSOutputItem `json:"result,omitempty"`
	Messages []Message     `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type GCSOutputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      GCSOutputData                `json:"data,omitempty"`
}

type GCSOutputListResult struct {
	TotalCount *int64          `json:"totalCount,omitempty"`
	Previous   *string         `json:"previous,omitempty"`
	Next       *string         `json:"next,omitempty"`
	Items      []GCSOutputItem `json:"items,omitempty"`
}

type GCSOutputListData struct {
	Result GCSOutputListResult `json:"result,omitempty"`
}

type GCSOutputListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      GCSOutputListData            `json:"data,omitempty"`
}

type S3Output struct {
	ID          *string                      `json:"id,omitempty"`
	Name        *string                      `json:"name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	CustomData  map[string]interface{}       `json:"customData,omitempty"`
	AccessKey   *string                      `json:"accessKey,omitempty"`
	SecretKey   *string                      `json:"secretKey,omitempty"`
	BucketName  *string                      `json:"bucketName,omitempty"`
	CloudRegion bitmovintypes.AWSCloudRegion `json:"cloudRegion,omitempty"`
}

type S3OutputItem struct {
	ID          *string                      `json:"id,omitempty"`
	Name        *string                      `json:"name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	BucketName  *string                      `json:"bucketName,omitempty"`
	CloudRegion bitmovintypes.AWSCloudRegion `json:"cloudRegion,omitempty"`
	CreatedAt   *string                      `json:"createdAt,omitempty"`
	UpdatedAt   *string                      `json:"updatedAt,omitempty"`
}

type S3OutputData struct {
	//Success fields
	Result   S3OutputItem `json:"result,omitempty"`
	Messages []Message    `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type S3OutputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      S3OutputData                 `json:"data,omitempty"`
}

type S3OutputListResult struct {
	TotalCount *int64         `json:"totalCount,omitempty"`
	Previous   *string        `json:"previous,omitempty"`
	Next       *string        `json:"next,omitempty"`
	Items      []S3OutputItem `json:"items,omitempty"`
}

type S3OutputListData struct {
	Result S3OutputListResult `json:"result,omitempty"`
}

type S3OutputListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      S3OutputListData             `json:"data,omitempty"`
}

type SFTPOutput struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	Host        *string                `json:"host,omitempty"`
	UserName    *string                `json:"username,omitempty"`
	Password    *string                `json:"password,omitempty"`
	Passive     *bool                  `json:"passive,omitempty"`
}
