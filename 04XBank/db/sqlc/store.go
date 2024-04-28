/*
* Created on 28 April 2024
* @author Sai Sumanth
 */
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries // composition
	db       *sql.DB
}

// Creates a New Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// 🗃️ Execute Transaction
// execTx executes a function (fn callback function) within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// create new queries
	q := New(tx)
	// execute the provided fn (Queries)
	err = fn(q)
	if err != nil {
		// check for any rollback errors
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// Input params required for transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// result of transfer Transaction
type TransferTxResult struct {
	TransferDetails Transfer `json:"transfer"`
	FromAccount     Account  `json:"from_account"`
	ToAccount       Account  `json:"to_account"`
	FromEntry       Entry    `json:"from_entry"`
	ToEntry         Entry    `json:"to_entry"`
}

// TransferTx - Money Transfer Transaction
//
// Performs a money transfer from one account to another account.
// Creates a ✅ transfer record, ✅ add entries for each account, ✅ and update accounts'
// balance within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		// 1️⃣ Add new Transfer entry in 'transfers' Table
		var err error
		result.TransferDetails, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2️⃣ Add two new entries for each user (from and to) inside 'entries' Table
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // Amount should be negative since money is being transferred from this account
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}


		// 3️⃣ Update Account Balance for both users
		// TODO: Update Balance


		return nil
	})

	return result, err
}
