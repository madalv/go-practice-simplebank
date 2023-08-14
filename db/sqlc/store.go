package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute queries & transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Store provides all functions to execute SQL queries & transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within db transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("TX err: %v, RB err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs money transfer from an acc to anoter.
// It creates a transfer record, entries for accounts, and updates accounts' balances.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var res TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create transfer record
		res.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))

		if err != nil {
			return err
		}

		// create entries
		res.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		res.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// update accounts' balances
		// important! always update accounts in the same order to avoid deadlocking concurrent transactions
		if arg.FromAccountID < arg.ToAccountID {
			res.FromAccount, res.ToAccount, err = addAmountToBalance(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			res.ToAccount, res.FromAccount, err = addAmountToBalance(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})

	return res, err
}

func addAmountToBalance(
	ctx context.Context,
	q *Queries,
	acc1ID int64,
	amount1 int64,
	acc2ID int64,
	amount2 int64,
) (acc1 Account, acc2 Account, err error) {
	acc1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     acc1ID,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	acc2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     acc2ID,
		Amount: amount2,
	})

	return
}
