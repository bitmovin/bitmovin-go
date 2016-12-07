package main

import (
	"fmt"

	"github.com/bitmovin/bitmovin-go/bitmovintypes"
)

func ErrorHandler(responseStatus bitmovintypes.ResponseStatus, err error) {
	if err != nil {
		fmt.Println("go error")
		fmt.Println(err)
	} else if responseStatus == "ERROR" {
		fmt.Println("api error")
	}
}

func StringToPtr(s string) *string {
	return &s
}

func IntToPtr(i int64) *int64 {
	return &i
}

func BoolToPtr(b bool) *bool {
	return &b
}

func FloatToPtr(f float64) *float64 {
	return &f
}
