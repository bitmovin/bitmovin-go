package main

import (
	"encoding/json"
	"fmt"

	"github.com/streamco/bitmovin-go/bitmovin"
	"github.com/streamco/bitmovin-go/bitmovintypes"
	"github.com/streamco/bitmovin-go/models"
	"github.com/streamco/bitmovin-go/services"
)

func main() {

	// Creating Bitmovin object
	bitmovin := bitmovin.NewBitmovin("YOUR_API_KEY", "https://api.bitmovin.com/v1/", 5)

	awsInfrastructureR := &models.CreateAWSInfrastructureRequest{
		Name:          stringToPtr("Test AWS Infrastructure"),
		Description:   stringToPtr("Test AWS Infrastructure"),
		AccessKey:     "YOUR_ACCESS_KEY",
		SecretKey:     "YOUR_SECRET_KEY",
		AccountNumber: "YOUR_ACCOUNT_NUMBER",
	}

	awsInfrastructureS := services.NewAWSInfrastructureService(bitmovin)
	awsInfrastructure, err := awsInfrastructureS.Create(awsInfrastructureR)

	if err != nil {
		panic(err)
	}

	awsInfrastructureRegionSettingsR := &models.CreateAWSInfrastructureRegionSettingsRequest{
		SecurityGroupId: "YOUR_SECURITY_GROUP_ID",
		SubnetId:        "YOUR_SUBNET_ID",
	}

	awsInfrastructureRegionSettingsS := services.NewAWSInfrastructureRegionSettingsService(bitmovin)
	_, err = awsInfrastructureRegionSettingsS.Create(awsInfrastructure.ID, bitmovintypes.AWSCloudRegionUSEast1, awsInfrastructureRegionSettingsR)

	if err != nil {
		panic(err)
	}

	awsInfrastructureRegionSettingsResp, err := awsInfrastructureRegionSettingsS.Retrieve(awsInfrastructure.ID, bitmovintypes.AWSCloudRegionUSEast1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("AWSInfrastructureRegionSettings Response: %s\n", getAsJsonString(*awsInfrastructureRegionSettingsResp))
}

func stringToPtr(s string) *string {
	return &s
}

func intToPtr(i int64) *int64 {
	return &i
}

func boolToPtr(b bool) *bool {
	return &b
}

func floatToPtr(f float64) *float64 {
	return &f
}

func getAsJsonString(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(j)
}
