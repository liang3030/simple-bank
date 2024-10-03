package api

import (
	"bytes"
	"fmt"
	"reflect"

	// "database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/liang3030/simple-bank/db/mock"
	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/util"
	"github.com/stretchr/testify/require"
)

// Example: custom matcher
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// This is a type assertion. It checks if x is of type db.CreateUserParams.
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	//  sets the HashedPassword field in e.arg (the expected CreateUserParams) to the HashedPassword from arg (the actual input).
	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func equalCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

// TODO: add more test cases
func TestCreateUser(t *testing.T) {
	user, password := createRandomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStub     func(store *mockdb.MockIStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStub: func(store *mockdb.MockIStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), equalCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "invalid-user#1",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStub: func(store *mockdb.MockIStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockIStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"

			request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func createRandomUser(t *testing.T) (user db.User, password string) {
	// user, _ := createRandomUser(t)
	password = util.RandomString(6)
	_, err := util.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
		// HashedPassword: hashedPassword,
	}
	return user, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	fmt.Printf("data: %s\n", data)

	var getUser db.User
	err = json.Unmarshal(data, &getUser)
	require.NoError(t, err)
	require.Equal(t, user, getUser)
}
