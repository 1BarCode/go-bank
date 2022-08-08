package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/1BarCode/go-bank/db/mock"
	db "github.com/1BarCode/go-bank/db/sqlc"
	"github.com/1BarCode/go-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestCreateAccountAPI(t *testing.T) {
	user := util.RandomOwner()
	account := randomAccount(user)

	testCases := []struct{
		name			string
		body			gin.H
		buildStubs		func(store *mockdb.MockStore)
		checkResponse 	func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"owner": account.Owner,
				"currency": account.Currency,
			} ,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				arg := db.CreateAccountParams{
					Owner: account.Owner,
					Currency: account.Currency,
					Balance: 0,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner": account.Owner,
				"currency": account.Currency,
			} ,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
		
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone) // should match with return signature of GetAccount method
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"owner": account.Owner,
				"currency": "invalid",
			} ,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)  // should match with return signature of GetAccount method
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
		
			// this is when the mock API call gets made
			server.router.ServeHTTP(recorder, request)	
			// check response
			tc.checkResponse(t, recorder)
		})
	}	
}

func TestGetAccountAPI(t *testing.T) {
	user := util.RandomOwner()
	account := randomAccount(user)

	// create a list of test case scenarios
	testCases := []struct{
		name string
		accountID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil) // should match with return signature of GetAccount method
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()). 
					// since accountID 0 is invalid, this 'get' method should not be called at all
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// TODO: add more cases
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
		
			server.router.ServeHTTP(recorder, request)	
			// check response
			tc.checkResponse(t, recorder)
		})
	}	

}

func TestListAccountsAPI(t *testing.T) {
	user := util.RandomOwner() 
	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user)
	}

	type listAccountsQuery struct {
		pageID 		int
		pageSize	int
	}

	testCases := []struct {
		name			string
		query			listAccountsQuery
		buildStubs 		func(store *mockdb.MockStore)
		checkResponse 	func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: listAccountsQuery{
				pageID: 1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				arg := db.ListAccountsParams{
					Limit: int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name: "InternalError",
			query: listAccountsQuery{
				pageID: 1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: listAccountsQuery{
				pageID: -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: listAccountsQuery{
				pageID: 1,
				pageSize: 100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response and body
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()
		
			server.router.ServeHTTP(recorder, request)	
			// check response
			tc.checkResponse(t, recorder)
		})
	}	
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID: util.RandomInt(1, 1000),
		Owner: owner,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}