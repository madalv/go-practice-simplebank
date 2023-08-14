package db

import (
	"context"
	"testing"
	"time"

	db "simplebank/db/sqlc"
	"simplebank/util"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, from db.Account, to db.Account) db.Transfer {
	arg := db.CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        util.RandomAmount(),
	}

	tran, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, tran.ID)
	require.NotZero(t, tran.CreatedAt)
	require.Equal(t, tran.Amount, arg.Amount)
	require.Equal(t, tran.FromAccountID, arg.FromAccountID)
	require.Equal(t, tran.ToAccountID, arg.ToAccountID)

	return tran
}

func TestCreateTransfer(t *testing.T) {
	from := createRandomAccount(t)
	to := createRandomAccount(t)
	createRandomTransfer(t, from, to)
}

func TestGetTransfer(t *testing.T) {
	from := createRandomAccount(t)
	to := createRandomAccount(t)
	created := createRandomTransfer(t, from, to)

	tran, err := testQueries.GetTransfer(context.Background(), created.ID)

	require.NoError(t, err)
	require.NotEmpty(t, created)
	require.Equal(t, tran.ID, created.ID)
	require.Equal(t, tran.FromAccountID, created.FromAccountID)
	require.Equal(t, tran.ToAccountID, created.ToAccountID)
	require.WithinDuration(t, created.CreatedAt, tran.CreatedAt, time.Second)
}

func testTransferArray(t *testing.T, tran []db.Transfer, len int) {
	require.NotEmpty(t, tran)
	require.Len(t, tran, len)

	for _, tr := range tran {
		require.NotEmpty(t, tr)
	}
}

func TestListTransfers(t *testing.T) {
	from := createRandomAccount(t)
	to := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, from, to)
		createRandomTransfer(t, to, from)
	}

	arg1 := db.ListTransfersParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Limit:         5,
		Offset:        0,
	}

	arg2 := db.ListTransfersParams{
		FromAccountID: to.ID,
		ToAccountID:   from.ID,
		Limit:         5,
		Offset:        0,
	}

	trans1, err := testQueries.ListTransfers(context.Background(), arg1)
	require.NoError(t, err)
	testTransferArray(t, trans1, 5)

	trans2, err := testQueries.ListTransfers(context.Background(), arg2)
	require.NoError(t, err)
	testTransferArray(t, trans2, 5)
}
