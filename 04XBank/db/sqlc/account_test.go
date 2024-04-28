/*
* Created on 27 April 2024
* @author Sai Sumanth
 */
package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/sai7xp/xbank/utils"
	"github.com/stretchr/testify/require"
)

// creates a test account
func createTestAccount(t *testing.T) (*Account, error) {
	args := CreateAccountParams{
		Owner:    fake.Person().Name(),
		Balance:  100, // default amount while creating new account
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	if err != nil {
		return nil, err
	}

	// Assertions
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotNil(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return &account, nil
}

// Unit Tests for CreateAccount
func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

// 🧪 Unit Tests for GetAccount
func TestGetAccount(t *testing.T) {
	// create account first
	acc, err := createTestAccount(t)
	require.NoError(t, err)

	getAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAcc)

	require.Equal(t, acc.ID, getAcc.ID)
	require.Equal(t, acc.Owner, getAcc.Owner)
	require.Equal(t, acc.Balance, getAcc.Balance)
	require.Equal(t, acc.Currency, getAcc.Currency)
	require.WithinDuration(t, acc.CreatedAt, getAcc.CreatedAt, time.Second)
}

// 🧪 Unit Tests for DeleteAcount
func TestDeleteAccount(t *testing.T) {
	// create account first
	newAccToBeDeleted, err := createTestAccount(t)
	require.NoError(t, err)

	// delete account
	delErr := testQueries.DeleteAccount(context.Background(), newAccToBeDeleted.ID)
	require.NoError(t, delErr)

	// try to get account details
	tryGetAccount, err := testQueries.GetAccount(context.Background(), newAccToBeDeleted.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tryGetAccount)
}

// 🧪 Unit Tests for Update Account
func TestUpdateAccount(t *testing.T) {
	// create account first
	newAcc, err := createTestAccount(t)
	require.NoError(t, err)

	updatedArgs := UpdateAccountParams{
		ID:      newAcc.ID,
		Balance: utils.RandomMoney(),
	}
	updatedAcc, err := testQueries.UpdateAccount(context.Background(), updatedArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)

	require.Equal(t, newAcc.ID, updatedAcc.ID)
	require.Equal(t, newAcc.Owner, updatedAcc.Owner)
	require.Equal(t, updatedArgs.Balance, updatedAcc.Balance)
	require.Equal(t, newAcc.Currency, updatedAcc.Currency)
	require.WithinDuration(t, newAcc.CreatedAt, updatedAcc.CreatedAt, time.Second)
}

// 🧪 Unit Tests for List Acounts
func TestListAccounts(t *testing.T) {
	// create few test accounts
	for i := 0; i < 8; i++ {
		createTestAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
