package auth_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	tokens "github.com/claytten/golang-simplebank/internal/api/token"
	mockdb "github.com/claytten/golang-simplebank/internal/db/mock"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func addRefreshTokenHeader(
	t *testing.T,
	request *http.Request,
	tokenMaker tokens.Maker,
	email string,
	duration time.Duration,
) (string, *tokens.Payload) {
	refreshToken, refreshTokenPayload, err := tokenMaker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, refreshTokenPayload)

	return refreshToken, refreshTokenPayload
}

func TestPostRenewTokenUserHandler(t *testing.T) {
	user, _ := util.RandomUser(t)

	durationMinute := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(durationMinute)

	session := db.Sessions{
		Email:     user.Email,
		IsBlocked: false,
		ExpiresAt: expiredAt,
		CreatedAt: issuedAt,
	}

	tests := []struct {
		name          string
		setupToken    func(t *testing.T, request *http.Request, tokenMaker tokens.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder, accessTokenDuration util.Config)
	}{
		// TODO: 200 OK
		{
			name: "200 OK",
			setupToken: func(t *testing.T, request *http.Request, tokenMaker tokens.Maker) {
				refreshToken, refreshTokenPayload := addRefreshTokenHeader(t, request, tokenMaker, user.Email, durationMinute)
				request.Header.Set(authorizationRefreshKey, refreshToken)
				session.ID = refreshTokenPayload.ID
				session.RefreshToken = refreshToken
				session.UserAgent = request.UserAgent()
				session.ClientIp = request.RemoteAddr

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(session.ID)).Return(session, nil).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, accessTokenDuration util.Config) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatchRenewToken(t, recorder.Body, accessTokenDuration.AccessTokenDuration, issuedAt)
			},
		},

		// TODO: 400 no header
		{
			name:       "400 no header",
			setupToken: func(t *testing.T, request *http.Request, tokenMaker tokens.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder, accessTokenDuration util.Config) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 401 not verify token
		{
			name: "401 no verify token",
			setupToken: func(t *testing.T, request *http.Request, tokenMaker tokens.Maker) {
				request.Header.Set(authorizationRefreshKey, "refreshToken")
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder, accessTokenDuration util.Config) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		// TODO: 404 not found data session
		{
			name: "404 not found data",
			setupToken: func(t *testing.T, request *http.Request, tokenMaker tokens.Maker) {
				refreshToken, refreshTokenPayload := addRefreshTokenHeader(t, request, tokenMaker, user.Email, durationMinute)
				request.Header.Set(authorizationRefreshKey, refreshToken)
				session.ID = refreshTokenPayload.ID
				session.RefreshToken = refreshToken
				session.UserAgent = request.UserAgent()
				session.ClientIp = request.RemoteAddr
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Any()).Return(db.Sessions{}, sql.ErrNoRows).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, accessTokenDuration util.Config) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		// TODO: 500 error query get data session
		{
			name: "500 error query get data session",
			setupToken: func(t *testing.T, request *http.Request, tokenMaker tokens.Maker) {
				refreshToken, refreshTokenPayload := addRefreshTokenHeader(t, request, tokenMaker, user.Email, durationMinute)
				request.Header.Set(authorizationRefreshKey, refreshToken)
				session.ID = refreshTokenPayload.ID
				session.RefreshToken = refreshToken
				session.UserAgent = request.UserAgent()
				session.ClientIp = request.RemoteAddr
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Any()).Return(db.Sessions{}, sql.ErrConnDone).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, accessTokenDuration util.Config) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		// TODO: 401 isBlocked?
		{
			name: "401 isBlocked",
			setupToken: func(t *testing.T, request *http.Request, tokenMaker tokens.Maker) {
				refreshToken, refreshTokenPayload := addRefreshTokenHeader(t, request, tokenMaker, user.Email, durationMinute)
				request.Header.Set(authorizationRefreshKey, refreshToken)
				session.ID = refreshTokenPayload.ID
				session.RefreshToken = refreshToken
				session.UserAgent = request.UserAgent()
				session.ClientIp = request.RemoteAddr
				session.IsBlocked = true
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(session.ID)).Return(session, nil).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, accessTokenDuration util.Config) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)

			server := api.NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			renewTokenPath := "/api/v1/auth/renew-token"
			request, err := http.NewRequest(http.MethodPost, renewTokenPath, nil)
			require.NoError(t, err)

			tt.setupToken(t, request, server.Token)
			tt.buildStubs(store)
			routes.ApplyAllPublicRoutes(server)
			server.Engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder, server.Config)
		})
	}
}

func RequireBodyMatchRenewToken(t *testing.T, body *bytes.Buffer, accessTokenDuration time.Duration, issuedAt time.Time) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotData map[string]interface{}
	err = json.Unmarshal(data, &gotData)
	require.NoError(t, err)
	require.NotEmpty(t, gotData)

	layout := "2006-01-02T15:04:05.999999999-07:00"
	accessExpiresAt := gotData["access_token_expires_at"].(string)
	timeParse, err := time.Parse(layout, accessExpiresAt)
	require.NoError(t, err)
	require.WithinDuration(t, issuedAt.Add(accessTokenDuration), timeParse, time.Second)
	require.NotEmpty(t, gotData["access_token"].(string))
}
