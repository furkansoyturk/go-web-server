package main

import (
	"testing"
)

func TestAsterisksWhenReqBodyContainsBlacklistKeyword(t *testing.T) {
	req := Request{
		Body: "This is a kerfuffle opinion I need to share with the world",
	}
	actualResponse := censorRequestBody(req)

	actualResult := actualResponse.Body
	expectedResult := "This is a **** opinion I need to share with the world"
	if actualResult != expectedResult {
		t.Fatalf("Expected result: %v. Actual result: %v", expectedResult, actualResult)
	}
}
