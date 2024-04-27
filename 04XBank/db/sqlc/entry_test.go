/*
* Created on 27 April 2024
* @author Sai Sumanth
 */
package db

import (
	"context"
	"testing"
	"time"

	"github.com/sai7xp/xbank/utils"
	"github.com/stretchr/testify/require"
)

// createTestEntry creates new Entry for the given account
func createTestEntry(t *testing.T, account *Account) Entry {
	createEntryArgs := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}
	newEntry, err := testQueries.CreateEntry(context.Background(), createEntryArgs)
	require.NoError(t, err)
	require.NotEmpty(t, newEntry)

	require.Equal(t, newEntry.AccountID, account.ID)
	require.Equal(t, newEntry.AccountID, createEntryArgs.AccountID)
	require.Equal(t, newEntry.Amount, createEntryArgs.Amount)

	require.NotZero(t, newEntry.ID)
	require.NotZero(t, newEntry)

	return newEntry
}

// 🧪 Unit Tests for CreateEntry
func TestCreateEntry(t *testing.T) {
	/// create new account first
	acc, err := createTestAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	createTestEntry(t, acc)
}

// 🧪 Unit Tests for GetEntry
func TestGetEntry(t *testing.T) {
	// create an account
	acc, err := createTestAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	// then create entry
	entry1 := createTestEntry(t, acc)

	getEntry, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getEntry)

	require.Equal(t, entry1.AccountID, getEntry.AccountID)
	require.Equal(t, entry1.Amount, getEntry.Amount)
	require.Equal(t, entry1.CreatedAt, getEntry.CreatedAt)
	require.WithinDuration(t, entry1.CreatedAt, getEntry.CreatedAt, time.Second)
}

// 🧪 Unit Tests for ListEntries
func TestListEntries(t *testing.T) {
	// create an account
	acc, err := createTestAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	// create 10 entries for [acc]
	for i := 0; i < 10; i++ {
		createTestEntry(t, acc)
	}

	listEntriesParams := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), listEntriesParams)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
