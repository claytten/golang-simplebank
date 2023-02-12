package account_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	"github.com/claytten/golang-simplebank/internal/api/token"
	mockdb "github.com/claytten/golang-simplebank/internal/db/mock"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestListsAccountsHandler(t *testing.T) {
	user, _ := util.RandomUser(t)

	n := 5
	accounts := make([]db.Accounts, n)
	for i := 0; i < n; i++ {
		accounts[i] = util.RandomAccount(user.Username)
	}

	type Query struct {
		PageID   int
		PageSize int
	}
	tests := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		// TODO: 200 OK
		{
			name: "200 OK",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()

				arg := db.ListsAccountsParams{
					Owner:  user.Username,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().ListsAccounts(gomock.Any(), gomock.Eq(arg)).Return(accounts, nil).Times(1)
				store.EXPECT().GetTotalPageListsAccounts(gomock.Any(), gomock.Eq(user.Username)).Return(int64(n), nil).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatchAccounts(t, recorder.Body, accounts, float64(1), float64(n))
			},
		},

		// TODO: 400 InvalidPageID
		{
			name: "400 invalid page id",
			query: Query{
				PageID:   -1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 400 InvalidPageSize
		{
			name: "400 invalid page size",
			query: Query{
				PageID:   1,
				PageSize: 1000000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 401 no token
		{
			name:       "400 no token",
			query:      Query{},
			setupAuth:  func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		// TODO: 500 query error
		{
			name: "500 query error",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()

				store.EXPECT().ListsAccounts(gomock.Any(), gomock.Any()).Return([]db.Accounts{}, sql.ErrConnDone).Times(1)
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

			getAccountPath := "/api/v1/accounts/lists"
			request, err := http.NewRequest(http.MethodGet, getAccountPath, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tt.query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", tt.query.PageSize))
			request.URL.RawQuery = q.Encode()

			tt.setupAuth(t, request, server.Token)
			routes.ApplyAllPublicRoutes(server)
			server.Engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func RequireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Accounts, current_page, limit float64) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotData map[string]interface{}
	err = json.Unmarshal(data, &gotData)
	require.NoError(t, err)

	require.Equal(t, current_page, gotData["current_page"])
	require.Equal(t, limit, gotData["limit"])

	for i, v := range gotData["data"].([]interface{}) {
		item, _ := v.(map[string]interface{})
		require.Equal(t, float64(accounts[i].ID), item["id"])
		require.Equal(t, float64(accounts[i].Balance), item["balance"])
		require.Equal(t, accounts[i].Owner, item["owner"])
		require.Equal(t, accounts[i].Currency, item["currency"])
	}

}
