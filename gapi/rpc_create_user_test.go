package gapi

import (
	"context"
	"fmt"
	"reflect"

	// "database/sql"

	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/liang3030/simple-bank/db/mock"
	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/pb"
	"github.com/liang3030/simple-bank/util"
	mockWK "github.com/liang3030/simple-bank/worker/mock"
	"github.com/stretchr/testify/require"
)

// Example: custom matcher
type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	// This is a type assertion. It checks if x is of type db.CreateUserTxParams.
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}

	//  sets the HashedPassword field in e.arg (the expected CreateUserTxParams) to the HashedPassword from arg (the actual input).
	expected.arg.HashedPassword = actualArg.HashedPassword

	return reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams)
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func equalCreateUserTxParams(arg db.CreateUserTxParams, password string) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (user db.User, password string) {
	// user, _ := createRandomUser(t)
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}
	return user, password
}

// TODO: add more test cases
func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStub     func(store *mockdb.MockIStore)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStub: func(store *mockdb.MockIStore) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
					// AfterCreate: func(user db.User) error {
					// 	taskPayload := &worker.PayloadSendVerifyEmail{
					// 		Username: user.Username,
					// 	}
					// 	opts := []asynq.Option{
					// 		asynq.MaxRetry(10),
					// 		asynq.ProcessIn(10 * time.Second),
					// 		asynq.Queue(worker.QueueCritical),
					// 	}
					// 	return server.taskDistributor.DistributeTaskSendVerifyEmail(context, taskPayload, opts...)
					// },
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), equalCreateUserTxParams(arg, password)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createUser := res.GetUser()
				require.Equal(t, user.Username, createUser.Username)
				require.Equal(t, user.FullName, createUser.FullName)
				require.Equal(t, user.Email, createUser.Email)

			},
		},
		// {
		// 	name: "InvalidUsername",
		// 	body: gin.H{
		// 		"username":  "invalid-user#1",
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStub: func(store *mockdb.MockIStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockIStore(storeCtrl)

			tc.buildStub(store)

			distributorCtrl := gomock.NewController(t)
			defer distributorCtrl.Finish()
			taskDistributor := mockWK.NewMockTaskDistributor(distributorCtrl)

			server := newTestServer(t, store, taskDistributor)
			res, err := server.CreateUser(context.Background(), tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}
