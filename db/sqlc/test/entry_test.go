package db

import (
	"context"
	"testing"
	"time"

	db "simplebank/db/sqlc"
	"simplebank/util"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, acc db.Account) db.Entry {
	arg := db.CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomAmount(),
	}

	ent, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, ent.ID)
	require.NotZero(t, ent.CreatedAt)
	require.Equal(t, ent.Amount, arg.Amount)
	require.Equal(t, ent.AccountID, arg.AccountID)

	return ent
}

func TestCreateEntry(t *testing.T) {
	acc := createRandomAccount(t)
	createRandomEntry(t, acc)
}

func TestGetEntry(t *testing.T) {
	acc := createRandomAccount(t)
	created := createRandomEntry(t, acc)

	ent, err := testQueries.GetEntry(context.Background(), created.ID)

	require.NoError(t, err)
	require.NotEmpty(t, created)
	require.Equal(t, ent.ID, created.ID)
	require.Equal(t, ent.AccountID, created.AccountID)
	require.Equal(t, ent.Amount, created.Amount)
	require.WithinDuration(t, created.CreatedAt, ent.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	acc := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc)
	}

	arg := db.ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    5,
	}

	ents, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, ents)
	require.Len(t, ents, 5)

	for _, ent := range ents {
		require.NotEmpty(t, ent)
	}
}
