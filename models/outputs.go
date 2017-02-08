package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type AzureOutput struct {
	ID          *string                `json:"id"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	CustomData  map[string]interface{} `json:"customData"`
	AccountName *string                `json:"accountName"`
	AccountKey  *string                `json:"accountKey"`
	Container   *string                `json:"container"`
}

type FTPOutput struct {
	ID          *string                `json:"id"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	CustomData  map[string]interface{} `json:"customData"`
	Host        *string                `json:"host"`
	UserName    *string                `json:"username"`
	Password    *string                `json:"password"`
	Passive     *bool                  `json:"passive"`
}

type GCSOutput struct {
	ID          *string                         `json:"id"`
	Name        *string                         `json:"name"`
	Description *string                         `json:"description"`
	CustomData  map[string]interface{}          `json:"customData"`
	AccessKey   *bool                           `json:"accessKey"`
	SecretKey   *bool                           `json:"secretKey"`
	BucketName  *bool                           `json:"bucketName"`
	CloudRegion bitmovintypes.GoogleCloudRegion `json:"cloudRegion"`
}

type S3Output struct {
	ID          *string                      `json:"id"`
	Name        *string                      `json:"name"`
	Description *string                      `json:"description"`
	CustomData  map[string]interface{}       `json:"customData"`
	AccessKey   *string                      `json:"accessKey"`
	SecretKey   *string                      `json:"secretKey"`
	BucketName  *string                      `json:"bucketName"`
	CloudRegion bitmovintypes.AWSCloudRegion `json:"cloudRegion"`
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
	ID          *string                `json:"id"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	CustomData  map[string]interface{} `json:"customData"`
	Host        *string                `json:"host"`
	UserName    *string                `json:"username"`
	Password    *string                `json:"password"`
	Passive     *bool                  `json:"passive"`
}
