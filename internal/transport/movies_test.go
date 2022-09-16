package transport

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
)

func TestHandler_getIdFromRequest(t *testing.T) {
	testTable := []struct {
		name          string
		id            string
		expectedId    int64
		expectedError error
	}{
		{
			name:          "OK",
			id:            "1",
			expectedId:    1,
			expectedError: nil,
		},
		{
			name:          "Parse ID error",
			id:            "oops",
			expectedId:    0,
			expectedError: errors.New("id parsing error"),
		},
		{
			name:          "Parse ID error",
			id:            "0",
			expectedId:    0,
			expectedError: errors.New("id can't be 0"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := httptest.NewRequest("GET", "/movies/1", nil)

			vars := map[string]string{
				"id": testCase.id,
			}
			r = mux.SetURLVars(r, vars)

			idfromReq, errFromReq := getIdFromRequest(r)

			assert.Equal(t, idfromReq, testCase.expectedId)
			assert.Equal(t, errFromReq, testCase.expectedError)
		})
	}
}
