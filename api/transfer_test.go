package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/alvarolucio2007/TheBank/db/mock"
	db "github.com/alvarolucio2007/TheBank/db/sqlc"
	"github.com/alvarolucio2007/TheBank/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransferAPI(t *testing.T) {
	user1, _ := randomUser(t)
	account1 := randomAccount(user1.Username)
	user2, _ := randomUser(t)
	account2 := randomAccount(user2.Username)
	testCases := []struct {
		name          string
		account1ID    int64
		account2ID    int64
		amount        int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			account1ID: account1.ID,
			account2ID: account2.ID,
			amount:     account1.Balance,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        account1.Balance,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TransferTxResult{
						Transfer:    db.Transfer{ID: 1, FromAccountID: account1.ID, ToAccountID: account2.ID, Amount: account1.Balance},
						FromAccount: account1,
						ToAccount:   account2,
						FromEntry:   db.Entry{ID: 1, AccountID: account1.ID, Amount: -account1.Balance},
						ToEntry:     db.Entry{ID: 2, AccountID: account2.ID, Amount: account1.Balance},
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := "/transfers"
			body := db.CreateTransferParams{
				FromAccountID: tc.account1ID,
				ToAccountID:   tc.account2ID,
				Amount:        tc.amount,
			}
			bodyJSON, err := json.Marshal(&body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyJSON))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
