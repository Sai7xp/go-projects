/*
* Created on 27 April 2024
* @author Sai Sumanth
 */

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestTransfer creates new Transfer in 'transfers' Table
func createTestTransfer(t *testing.T, fromAcc *Account, toAcc *Account) Transfer {
	args := CreateTransferParams{FromAccountID: fromAcc.ID,
		ToAccountID: toAcc.ID,
		Amount:      int64(fake.Currency().Number()),
	}
	newTransfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, newTransfer)
	require.Equal(t, args.FromAccountID, newTransfer.FromAccountID)
	require.Equal(t, args.ToAccountID, newTransfer.ToAccountID)
	require.Equal(t, args.Amount, newTransfer.Amount)

	assert.NotZero(t, newTransfer.ID, "New Transfer ID can't be zero")

	return newTransfer
}

// 🧪 Unit Tests for CreateTransfer
func TestCreateTransfer(t *testing.T) {
	// create two accounts
	acc1, err := createTestAccount(t)
	require.NoError(t, err)

	acc2, err := createTestAccount(t)
	require.NoError(t, err)

	newTransfer := createTestTransfer(t, acc1, acc2)
	require.NotEmpty(t, newTransfer)

}

// 🧪 Unit Tests for GetTransfer
func TestGetTransfer(t *testing.T) {
	// create two accounts
	acc1, err := createTestAccount(t)
	require.NoError(t, err)

	acc2, err := createTestAccount(t)
	require.NoError(t, err)

	// create new Transfer
	newTransfer := createTestTransfer(t, acc1, acc2)
	require.NotEmpty(t, newTransfer)

	// try to get the created transfer
	getTransfer, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getTransfer)

	assert.Equal(t, getTransfer.FromAccountID, newTransfer.FromAccountID, "From Account IDs should be equal")
	assert.Equal(t, getTransfer.ToAccountID, newTransfer.ToAccountID, "To Account IDs should be equal")
	assert.Equal(t, getTransfer.Amount, newTransfer.Amount, "New Transfer Amount & Get Transfer Amount should match")
}

// 🧪 Unit Tests for ListTransfers
func TestListTransfers(t *testing.T) {
	// create two accounts
	acc1, _ := createTestAccount(t)
	acc2, _ := createTestAccount(t)

	// create few transfers between two accounts
	for i := 0; i < 5; i++ {
		createTestTransfer(t, acc1, acc2)
		createTestTransfer(t, acc2, acc1)
	}

	args := ListTransfersParams{FromAccountID: acc1.ID,
		ToAccountID: acc2.ID,
		Limit:       5,
	}
	args2 := ListTransfersParams{FromAccountID: acc2.ID,
		ToAccountID: acc1.ID,
		Limit:       5,
	}
	t1, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)

	t2, err := testQueries.ListTransfers(context.Background(), args2)
	require.NoError(t, err)

	transfers := append(t1, t2...)
	require.Len(t, transfers, 10)

	for _, eachTransfer := range transfers {
		assert.NotEmpty(t, eachTransfer, "Transfer can't be empty")
		require.True(t, eachTransfer.FromAccountID == acc1.ID || eachTransfer.ToAccountID == acc1.ID)
	}
}
