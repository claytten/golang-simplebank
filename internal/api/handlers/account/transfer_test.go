package account_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	"github.com/stretchr/testify/require"
)

func TestPostCreateTransferAccountHandler(t *testing.T) {
	amount := int64(10)
	user1, oldPassword1 := util.RandomUser(t)
	user2, _ := util.RandomUser(t)
	user3, _ := util.RandomUser(t)

	acc1 := util.RandomAccount(user1.Username)
	acc2 := util.RandomAccount(user2.Username)
	acc3 := util.RandomAccount(user3.Username)

	acc1.Currency = "USD"
	acc2.Currency = "USD"
	acc3.Currency = "CAD"

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
				"amount":   amount,
				"currency": util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Email, time.Minute)
				request.Header.Set(authorizationUsername, user1.Username)
				request.Header.Set(authorizationOldPassword, oldPassword1)
				request.Header.Set("from_account_id", strconv.Itoa(int(acc1.ID)))
				request.Header.Set("to_account_id", strconv.Itoa(int(acc2.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user1.Email)).Return(user1, nil).AnyTimes()

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).Return(acc1, nil).Times(1)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).Return(acc2, nil).Times(1)

				arg := db.TransferTxParams{
					FromAccountID: acc1.ID,
					ToAccountID:   acc2.ID,
					Amount:        amount,
				}

				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		// TODO: 400 one of account id missing
		{
			name: "400 id missing",
			body: gin.H{
				"amount":   amount,
				"currency": util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Email, time.Minute)
				request.Header.Set(authorizationUsername, user1.Username)
				request.Header.Set(authorizationOldPassword, oldPassword1)
				request.Header.Set("from_account_id", strconv.Itoa(int(acc1.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user1.Email)).Return(user1, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 400 one of body json missing
		{
			name: "400 body missing",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Email, time.Minute)
				request.Header.Set(authorizationUsername, user1.Username)
				request.Header.Set(authorizationOldPassword, oldPassword1)
				request.Header.Set("from_account_id", strconv.Itoa(int(acc1.ID)))
				request.Header.Set("to_account_id", strconv.Itoa(int(acc2.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user1.Email)).Return(user1, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

		// TODO: 401 token and user mismatch
		{
			name: "401 mismatch user and token",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Email, time.Minute)
				request.Header.Set(authorizationUsername, user2.Username)
				request.Header.Set(authorizationOldPassword, oldPassword1)
				request.Header.Set("from_account_id", strconv.Itoa(int(acc1.ID)))
				request.Header.Set("to_account_id", strconv.Itoa(int(acc2.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user1.Email)).Return(user1, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		// TODO: 404 not found to account
		{
			name: "404 not found to account",
			body: gin.H{
				"amount":   amount,
				"currency": util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Email, time.Minute)
				request.Header.Set(authorizationUsername, user1.Username)
				request.Header.Set(authorizationOldPassword, oldPassword1)
				request.Header.Set("from_account_id", strconv.Itoa(int(acc1.ID)))
				request.Header.Set("to_account_id", strconv.Itoa(int(acc2.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user1.Email)).Return(user1, nil).AnyTimes()
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).Return(acc1, nil).Times(1)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).Return(db.Accounts{}, sql.ErrNoRows).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		// TODO: 400 currency mismatch
		{
			name: "400 currency mismatch",
			body: gin.H{
				"amount":   amount,
				"currency": "XYZ",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Email, time.Minute)
				request.Header.Set(authorizationUsername, user1.Username)
				request.Header.Set(authorizationOldPassword, oldPassword1)
				request.Header.Set("from_account_id", strconv.Itoa(int(acc1.ID)))
				request.Header.Set("to_account_id", strconv.Itoa(int(acc2.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user1.Email)).Return(user1, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 500 error query from account
		{
			name: "500 error query from account",
			body: gin.H{
				"amount":   amount,
				"currency": util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Email, time.Minute)
				request.Header.Set(authorizationUsername, user1.Username)
				request.Header.Set(authorizationOldPassword, oldPassword1)
				request.Header.Set("from_account_id", "999999")
				request.Header.Set("to_account_id", strconv.Itoa(int(acc1.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user1.Email)).Return(user1, nil).AnyTimes()
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(db.Accounts{}, sql.ErrConnDone).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		// TODO: 500 query transfer error
		{
			name: "500 query transfers",
			body: gin.H{
				"amount":   amount,
				"currency": util.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Email, time.Minute)
				request.Header.Set(authorizationUsername, user1.Username)
				request.Header.Set(authorizationOldPassword, oldPassword1)
				request.Header.Set("from_account_id", strconv.Itoa(int(acc1.ID)))
				request.Header.Set("to_account_id", strconv.Itoa(int(acc2.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user1.Email)).Return(user1, nil).AnyTimes()

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).Return(acc1, nil).Times(1)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).Return(acc2, nil).Times(1)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Return(db.TransferTxResult{}, sql.ErrTxDone).Times(1)
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

			// body data to json
			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			getAccountPath := "/api/v1/accounts/transfer"
			request, err := http.NewRequest(http.MethodPost, getAccountPath, bytes.NewReader(data))
			require.NoError(t, err)

			tt.setupAuth(t, request, server.Token)
			routes.ApplyAllPublicRoutes(server)
			server.Engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
