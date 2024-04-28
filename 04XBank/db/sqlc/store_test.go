/*
* Created on 28 April 2024
* @author Sai Sumanth
 */

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// first create 2 new accounts
	acc1, _ := createTestAccount(t)
	acc2, _ := createTestAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errChan := make(chan error)
	resultsChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// new goroutine
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errChan <- err
			resultsChan <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)

		res := <-resultsChan
		require.NotEmpty(t, res)

		// check transfer
		transfer := res.TransferDetails
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// Now try querying the Transfer Table
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check for entries
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)
		
		// TODO: check acounts' balances
	}
}
