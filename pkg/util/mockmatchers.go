package util

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/golang/mock/gomock"
)

type httpRequestEqMatcher struct {
	expected *http.Request
}

// HTTPRequestEq returns a matcher that compares http.Request instants.
// Only Method, URL, Body, and Header are taken into account.
func HTTPRequestEq(expected *http.Request) gomock.Matcher {
	return httpRequestEqMatcher{expected}
}

func (m httpRequestEqMatcher) Matches(actual interface{}) bool {
	if actualRequest, ok := actual.(*http.Request); ok {
		if m.expected.Body == nil && actualRequest.Body != nil {
			return false
		}
		if m.expected.Body != nil {
			if actualRequest.Body == nil {
				return false
			}

			expectedBody := new(strings.Builder)
			if _, err := io.Copy(expectedBody, m.expected.Body); err != nil {
				panic(err)
			}

			actualBody := new(strings.Builder)
			if _, err := io.Copy(actualBody, actualRequest.Body); err != nil {
				panic(err)
			}

			if actualBody.String() != expectedBody.String() {
				return false
			}
		}

		return actualRequest.Method == m.expected.Method &&
			actualRequest.URL.String() == m.expected.URL.String() &&
			reflect.DeepEqual(actualRequest.Header, m.expected.Header)
	}

	return false
}

func (m httpRequestEqMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}
