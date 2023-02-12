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
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type EqUpdateUserParamsMatcher struct {
	arg db.UpdateUserParams
}

func (e *EqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	if arg.FullName.String != e.arg.FullName.String && arg.Email.String != e.arg.Email.String {
		return false
	}

	argUpdatedAt := arg.UpdatedAt.Round(time.Second)
	expected := e.arg.UpdatedAt.Round(time.Second)
	if ok := argUpdatedAt.Equal(expected); !ok {
		return false
	}

	e.arg.FullName.String = arg.FullName.String
	e.arg.Email.String = arg.Email.String
	e.arg.UpdatedAt = arg.UpdatedAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e *EqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqUpdateUserParams(arg db.UpdateUserParams) gomock.Matcher {
	return &EqUpdateUserParamsMatcher{arg}
}

func TestUpdateUserProfileHandler(t *testing.T) {
	now := time.Now()
	user, plainPassword := util.RandomUser(t)

	newUser := user
	newUser.FullName = util.RandomOwner()
	newUser.Email = util.RandomEmail()
	newUser.UpdatedAt = now

	newUser1 := user
	newUser1.Email = util.RandomEmail()
	newUser1.UpdatedAt = now

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
				"full_name": newUser.FullName,
				"email":     newUser.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, plainPassword)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Return(user, nil).AnyTimes()

				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newUser.Username,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newUser.Email,
						Valid:  true,
					},
					UpdatedAt: now,
				}
				store.EXPECT().UpdateUser(gomock.Any(), EqUpdateUserParams(arg)).
					Times(1).Return(newUser, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatchProfile(t, recorder.Body, newUser)
			},
		},

		// TODO: 200 OK just one params change
		{
			name: "200 OK one params",
			body: gin.H{
				"full_name": newUser1.FullName,
				"email":     newUser1.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, plainPassword)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Times(1).Return(user, nil)
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)

				arg := db.UpdateUserParams{
					Username: newUser1.Username,
					FullName: sql.NullString{
						String: newUser1.Username,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newUser1.Email,
						Valid:  true,
					},
					UpdatedAt: now,
				}
				store.EXPECT().UpdateUser(gomock.Any(), EqUpdateUserParams(arg)).
					Times(1).Return(newUser1, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatchProfile(t, recorder.Body, newUser1)
			},
		},

		// TODO: 401 Token not found
		{
			name:       "401 token not found",
			body:       gin.H{},
			setupAuth:  func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		// TODO: 401 no header params (exclude token)
		{
			name: "401 no header",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 401 user params not match with user token
		{
			name: "401 user params not match",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, "notfounded")
				request.Header.Set(authorizationOldPassword, plainPassword)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Times(1).Return(db.Users{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		// TODO: 401 oldpassword doesn't match
		{
			name: "401 mismatch oldpassword",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, "notFounded")
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		// TODO: 400 no body profile (fullname and email)
		{
			name: "400 badrequest body",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, plainPassword)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 400 if one of header doesn't fill it
		{
			name: "400 none one of header",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: 500 query error
		{
			name: "500 query error",
			body: gin.H{
				"full_name": newUser.FullName,
				"email":     newUser.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
				request.Header.Set(authorizationUsername, user.Username)
				request.Header.Set(authorizationOldPassword, plainPassword)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserUsingEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Return(user, nil).AnyTimes()
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(1).Return(newUser, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			//Marshal body to json
			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			getUserPath := "/api/v1/auth/profile"
			request, err := http.NewRequest(http.MethodPut, getUserPath, bytes.NewReader(data))
			require.NoError(t, err)

			tt.setupAuth(t, request, server.Token)
			routes.ApplyAllPublicRoutes(server)
			server.Engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func RequireBodyMatchProfile(t *testing.T, body *bytes.Buffer, user db.Users) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.Users
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
}
