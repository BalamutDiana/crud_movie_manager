package transport

import (
	"net/http/httptest"
	"testing"

	//mock_service "github.com/BalamutDiana/crud_movie_manager/internal/transport/mocks"

	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_getTokenFromRequest(t *testing.T) {
	testTable := []struct {
		name             string
		headerName       string
		headerValue      string
		token            string
		expectedResponse string
		expectedError    error
	}{
		{
			name:             "OK",
			headerName:       "Authorization",
			headerValue:      "Bearer token",
			token:            "token",
			expectedResponse: "token",
			expectedError:    nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			tokenFromReq, errFromReq := getTokenFromRequest(req)

			assert.Equal(t, tokenFromReq, testCase.expectedResponse)
			assert.Equal(t, errFromReq, testCase.expectedError)
		})
	}
}
