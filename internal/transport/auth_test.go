package transport

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BalamutDiana/crud_movie_manager/internal/domain"
	repo "github.com/BalamutDiana/crud_movie_manager/internal/repository"
	mock_service "github.com/BalamutDiana/crud_movie_manager/internal/transport/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, ctx context.Context, user domain.SignUpInput)

	testTable := []struct {
		name               string
		inputBody          string
		inputUser          domain.SignUpInput
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","email":"test@gmail.com","password":"qwerty"}`,
			inputUser: domain.SignUpInput{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockUser, ctx context.Context, user domain.SignUpInput) {
				s.EXPECT().SignUp(gomock.Any(), user).Return(nil)
			},
			expectedStatusCode: 200,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init deps
			c := gomock.NewController(t)
			//c, ctx := gomock.WithContext(context.TODO(), t)
			defer c.Finish()

			auth := mock_service.NewMockUser(c)
			testCase.mockBehavior(auth, context.Background(), testCase.inputUser)

			movieRepo := &repo.Movies{}
			h := NewHandler(movieRepo, auth)

			// Test Server
			r := mux.NewRouter()
			a := r.PathPrefix("/auth").Subrouter()
			{
				a.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
			}

			//Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)

		})
	}

}
