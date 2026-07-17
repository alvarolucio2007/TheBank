package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/alvarolucio2007/TheBank/db/mock"
	db "github.com/alvarolucio2007/TheBank/db/sqlc"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransferAPI(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()
	testCases := []struct {
		name          string
		account1ID    int64
		account2ID    int64
		amount        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			account1ID: account1.ID,
			account2ID: account2.ID,
			amount:     account1.Balance,
			buildStubs: func(store *mockdb.MockStore) {
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
			server, err := NewServer(store)
			require.NoError(t, err)
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

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
