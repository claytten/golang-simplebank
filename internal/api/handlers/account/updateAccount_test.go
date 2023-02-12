package account_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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

type EqUpdateAccountParamsMatcher struct {
	arg       db.UpdateAccountParams
	updatedAt time.Time
}

func (e *EqUpdateAccountParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateAccountParams)
	if !ok {
		return false
	}

	argUpdatedAt := arg.UpdatedAt.Round(time.Second)
	expected := e.updatedAt.Round(time.Second)
	if ok := argUpdatedAt.Equal(expected); !ok {
		return false
	}

	e.arg.UpdatedAt = arg.UpdatedAt
	return reflect.DeepEqual(e.arg, arg)
}

func (e *EqUpdateAccountParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.updatedAt)
}

func EqUpdateAccountParams(arg db.UpdateAccountParams, updatedAt time.Time) gomock.Matcher {
	return &EqUpdateAccountParamsMatcher{arg, updatedAt}
}

func TestPutUpdateAccountHandler(t *testing.T) {
	user, oldPassword := util.RandomUser(t)
	account := util.RandomAccount(user.Username)
	now := time.Now()
	addNewBalance := int64(200)

	newAccount := account
	newAccount.Balance = addNewBalance
	newAccount.UpdatedAt = now

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
				"balance": addNewBalance,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, oldPassword)
				request.Header.Set("id", strconv.Itoa(int(account.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()

				arg := db.UpdateAccountParams{
					ID:      account.ID,
					Balance: addNewBalance,
				}
				store.EXPECT().UpdateAccount(gomock.Any(), EqUpdateAccountParams(arg, now)).Return(newAccount, nil).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

		// TODO: 400 missing id header
		{
			name: "400 no id header",
			body: gin.H{
				"balance": addNewBalance,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, oldPassword)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 400 missing body balance
		{
			name: "400 missing body balance",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, oldPassword)
				request.Header.Set("id", strconv.Itoa(int(account.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 500 run query
		{
			name: "500 error query",
			body: gin.H{
				"balance": addNewBalance,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, oldPassword)
				request.Header.Set("id", strconv.Itoa(int(account.ID)))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Return(db.Accounts{}, sql.ErrConnDone).Times(1)
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

			getAccountPath := "/api/v1/accounts/update"
			request, err := http.NewRequest(http.MethodPut, getAccountPath, bytes.NewReader(data))
			require.NoError(t, err)

			tt.setupAuth(t, request, server.Token)
			routes.ApplyAllPublicRoutes(server)
			server.Engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
