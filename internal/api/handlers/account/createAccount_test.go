package account_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	"github.com/claytten/golang-simplebank/internal/api/token"
	mockdb "github.com/claytten/golang-simplebank/internal/db/mock"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	email string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

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

func TestPostCreateAccountHandler(t *testing.T) {
	user, password := util.RandomUser(t)
	account := db.Accounts{
		Owner:    user.Username,
		Balance:  0,
		Currency: util.RandomCurrency(),
	}

	tests := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		// TODO: 200 OK
		{
			name: "200 OK",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, password)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()

				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  account.Balance,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).Return(account, nil).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatchUser(t, recorder.Body, account)
			},
		},

		// TODO: 401 no token
		{
			name:       "401 no token",
			body:       gin.H{},
			setupAuth:  func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		// TODO: 400 something missing header
		{
			name: "400 missing some header",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Any()).Return(db.Users{}, nil).Times(0)
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(db.Accounts{}, nil).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 400 something missing body
		{
			name: "400 missing some body",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, password)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(db.Accounts{}, nil).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 400 invalid currency (USD, CAD, EUR)
		{
			name: "400 invalid currency",
			body: gin.H{
				"currency": "SGD",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, password)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(db.Accounts{}, nil).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 500 error while run query
		{
			name: "500 query",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, password)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(db.Accounts{}, sql.ErrConnDone).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		// TODO: 403 unique violation
		{
			name: "403 unique violation",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, password)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()

				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(db.Accounts{}, &pq.Error{Code: "23505"}).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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

			// body data to json
			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			getUserPath := "/api/v1/accounts/create"
			request, err := http.NewRequest(http.MethodPost, getUserPath, bytes.NewReader(data))
			require.NoError(t, err)

			tt.setupAuth(t, request, server.Token)
			routes.ApplyAllPublicRoutes(server)
			server.Engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func RequireBodyMatchUser(t *testing.T, body *bytes.Buffer, account db.Accounts) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Accounts
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account.Owner, gotAccount.Owner)
	require.Equal(t, account.Currency, gotAccount.Currency)
	require.Equal(t, account.Balance, gotAccount.Balance)
	require.NotNil(t, gotAccount.CreatedAt)
	require.NotNil(t, gotAccount.UpdatedAt)

	//checking balance account are zero while first creating
	//intial to int64()
	require.Equal(t, int64(0), gotAccount.Balance)
}
