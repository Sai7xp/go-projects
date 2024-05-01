/*
* Created on 28 April 2024
* @author Sai Sumanth
 */

package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// first create 2 new accounts
	acc1, _ := createTestAccount(t)
	acc2, _ := createTestAccount(t)
	fmt.Println(">> Acc Balance Before:", acc1.Balance, acc2.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errChan := make(chan error)
	resultsChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		// new goroutine
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
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

		// check acounts' balances
		fmt.Println(">> Acc Balance After Each Transaction:", fromAccount.Balance, toAccount.Balance)
		diff1 := acc1.Balance - fromAccount.Balance // amount that goes from acc1
		diff2 := toAccount.Balance - acc2.Balance   // amount that gets added to acc2
		require.Equal(t, diff1, diff2)              // both amounts should be equal
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
	}

	// check for final updated balances
	updatedAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> Acc Balance After:", updatedAcc1.Balance, updatedAcc2.Balance)
	// we are performing n transactions so we need to multiply amount with n
	require.Equal(t, acc1.Balance-int64(n)*amount, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updatedAcc2.Balance)
}

// let's check for db deadlock
// check run_queries.sql file to understand this deadlock condition and
// how to avoid id just by changing the order of query execution
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	// first create 2 new accounts
	acc1, _ := createTestAccount(t)
	acc2, _ := createTestAccount(t)
	fmt.Println(">> Acc Balance Before:", acc1.Balance, acc2.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errChan := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := acc1.ID
		toAccountId := acc2.ID

		// transfer money FROM both the accounts
		if i%2 == 1 {
			fromAccountId = acc2.ID
			toAccountId = acc1.ID
		}

		// new goroutine
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errChan <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)
	}

	// check for final updated balances
	updatedAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> Acc Balance After Transferring the same amount from both accounts(should be equal):", updatedAcc1.Balance, updatedAcc2.Balance)
	require.Equal(t, acc1.Balance, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
