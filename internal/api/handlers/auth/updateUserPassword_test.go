package auth_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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

type eqUpdatePassUserParamsMatcher struct {
	arg      db.UpdateUserParams
	password string
}

func (e *eqUpdatePassUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	err := util.ComparePassword(arg.HashedPassword.String, e.password)
	if err != nil || arg.PasswordChangedAt.Time == e.arg.PasswordChangedAt.Time || arg.UpdatedAt == e.arg.UpdatedAt {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	e.arg.PasswordChangedAt = arg.PasswordChangedAt
	e.arg.UpdatedAt = arg.UpdatedAt
	return reflect.DeepEqual(e.arg, arg)
}

func (e *eqUpdatePassUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqUpdatePassUserParams(arg db.UpdateUserParams, password string) gomock.Matcher {
	return &eqUpdatePassUserParamsMatcher{arg, password}
}

func TestUpdateUserPasswordHandler(t *testing.T) {
	now := time.Now()
	newPassword := "ajibaru"
	hashedNewPassword, err := util.HashingPassword(newPassword)
	require.NoError(t, err)

	oldUser, oldPassword := util.RandomUser(t)

	newUser := oldUser
	newUser.HashedPassword = hashedNewPassword
	newUser.PasswordChangedAt = now
	newUser.PasswordChangedAt = now

	tests := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		//TODO: 200 OK
		{
			name: "200 OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, oldUser.Email, time.Minute)
				request.Header.Set(authorizationUsername, oldUser.Username)
				request.Header.Set(authorizationOldPassword, oldPassword) //user plain password (oldpassword)
				request.Header.Set(authorizationPassword, newPassword)    //user new password
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(oldUser.Email)).Return(oldUser, nil).AnyTimes()
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(oldUser.Username)).Return(oldUser, nil).AnyTimes()

				args := db.UpdateUserParams{
					Username: oldUser.Username,
					PasswordChangedAt: sql.NullTime{
						Time:  now,
						Valid: true,
					},
					UpdatedAt: now,
				}
				store.EXPECT().UpdateUser(gomock.Any(), EqUpdatePassUserParams(args, newPassword)).Return(newUser, nil).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatchProfilePassword(t, recorder.Body, newUser)
			},
		},

		//TODO: 401 no token header
		{
			name:       "401 no token header",
			setupAuth:  func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		//TODO: 401 old password doesn't match
		{
			name: "401 old password mismatch",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, oldUser.Email, time.Minute)
				request.Header.Set(authorizationUsername, oldUser.Username)
				request.Header.Set(authorizationOldPassword, "notfoundoldPassword") //user plain password (oldpassword)
				request.Header.Set(authorizationPassword, newPassword)              //user new password
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(oldUser.Email)).Return(db.Users{}, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		//TODO: 400 header middleware one or multiple params are missing
		{
			name: "400 missing header",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, oldUser.Email, time.Minute)
				request.Header.Set(authorizationUsername, oldUser.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		//TODO: 400 header handler one or multiple params are missing
		{
			name: "400 header handler mismatch",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, oldUser.Email, time.Minute)
				request.Header.Set(authorizationUsername, oldUser.Username)
				request.Header.Set(authorizationOldPassword, oldPassword) //user plain password (oldpassword)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(oldUser.Email)).Return(oldUser, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		//TODO: 401 mismatch header username with payload token
		{
			name: "401 mismatch header username",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, oldUser.Email, time.Minute)
				request.Header.Set(authorizationUsername, "NewUserComing")
				request.Header.Set(authorizationOldPassword, oldPassword) //user plain password (oldpassword)
				request.Header.Set(authorizationPassword, newPassword)    //user new password
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(oldUser.Email)).Return(oldUser, nil).AnyTimes()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			getUserPath := "/api/v1/auth/profile/password"
			request, err := http.NewRequest(http.MethodPut, getUserPath, nil)
			require.NoError(t, err)

			tt.setupAuth(t, request, server.Token)
			routes.ApplyAllPublicRoutes(server)
			server.Engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func RequireBodyMatchProfilePassword(t *testing.T, body *bytes.Buffer, newUser db.Users) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.Users
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.WithinDuration(t, newUser.PasswordChangedAt, gotUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, newUser.UpdatedAt, gotUser.UpdatedAt, time.Second)
}
