package net

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

type JSONRequest map[string]interface{}

func RequestBodyMatcher(expectedBody string) RequestMatcher {
	return func(t *testing.T, request *http.Request) {
		bodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			assert.Fail(t,"Error reading request body: %s",err)
		}

		actualBody := &JSONRequest{}
		json.Unmarshal(bodyBytes,actualBody)
		expectedBody := &JSONRequest{}
		json.Unmarshal(bodyBytes,expectedBody)

		assert.Equal(t,actualBody,expectedBody)

		assert.Equal(t,request.Header.Get("content-type"), "application/json", "Content Type was not application/json.")
	}
}

func RequestBodyMatcherWithContentType(expectedBody, expectedContentType string) RequestMatcher {
	return func(t *testing.T, request *http.Request) {
		bodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			assert.Fail(t,"Error reading request body: %s",err)
		}

		actualBody := string(bodyBytes)
		assert.Equal(t,RemoveWhiteSpaceFromBody(actualBody),RemoveWhiteSpaceFromBody(expectedBody), "Body did not match.")

		actualContentType := request.Header.Get("content-type")
		assert.Equal(t,actualContentType,expectedContentType, "Content Type did not match.")
	}
}

func RemoveWhiteSpaceFromBody(body string) string {
	body = strings.Replace(body, " ", "", -1)
	body = strings.Replace(body, "\n", "", -1)
	body = strings.Replace(body, "\r", "", -1)
	body = strings.Replace(body, "\t", "", -1)
	return body
}
