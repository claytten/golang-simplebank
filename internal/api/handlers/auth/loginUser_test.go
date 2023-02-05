package auth_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	mockdb "github.com/claytten/golang-simplebank/internal/db/mock"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	authorizationHeaderKey   = "authorization"
	authorizationTypeBearer  = "bearer"
	authorizationPayloadKey  = "authorization_payload"
	authorizationPassword    = "password"
	authorizationOldPassword = "oldPassword"
	authorizationUsername    = "username"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func TestPostLoginUserRoute(t *testing.T) {
	user, plainPassword := util.RandomUser(t)
	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		// TODO: Login OK
		{
			name: "200 OK",
			body: gin.H{
				"email":    user.Email,
				"password": plainPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		// TODO: Login User Not Found
		{
			name: "404 Not Found",
			body: gin.H{
				"email":    "notfound@email.com",
				"password": plainPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserUsingEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Users{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		// TODO: Incorrect Password
		{
			name: "500 Incorrect Pass",
			body: gin.H{
				"email":    user.Email,
				"password": "incorrectPass",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			tt.buildStubs(store)

			server := api.NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to json
			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			loginPath := "/api/v1/auth/login"
			request, err := http.NewRequest(http.MethodPost, loginPath, bytes.NewReader(data))
			require.NoError(t, err)

			routes.ApplyAllPublicRoutes(server)
			server.Engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
