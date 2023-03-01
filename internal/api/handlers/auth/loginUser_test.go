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
	authorizationRefreshKey  = "refresh_token"
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
		// TODO: 200 Login OK
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
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		// TODO: 400 no body
		{
			name:       "400 no body",
			body:       gin.H{},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 404 Login User Not Found
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

		// TODO: 500 query get user email
		{
			name: "500 query error",
			body: gin.H{
				"email":    user.Email,
				"password": plainPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Any()).Return(db.Users{}, sql.ErrConnDone).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		// TODO: 500 query create session
		{
			name: "500 query create session",
			body: gin.H{
				"email":    user.Email,
				"password": plainPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), user.Email).Return(user, nil).Times(1)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(db.Sessions{}, sql.ErrConnDone).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		// TODO: 500 Incorrect Password
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
