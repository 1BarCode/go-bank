package db

import (
	"context"
	"testing"

	"github.com/1BarCode/go-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAcct, toAcct Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAcct.ID,
		ToAccountID: toAcct.ID,
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAcct := createRandomAccount(t)
	toAcct := createRandomAccount(t)
	createRandomTransfer(t, fromAcct, toAcct)
}

func TestGetTransfer(t *testing.T) {
	fromAcct := createRandomAccount(t)
	toAcct := createRandomAccount(t)

	newTransfer := createRandomTransfer(t, fromAcct, toAcct)
	fetchedTransfer, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTransfer)

	require.Equal(t, newTransfer.ID, fetchedTransfer.ID)
	require.Equal(t, newTransfer.FromAccountID, fetchedTransfer.FromAccountID)
	require.Equal(t, newTransfer.ToAccountID, fetchedTransfer.ToAccountID)
	require.Equal(t, newTransfer.Amount, fetchedTransfer.Amount)
	require.Equal(t, newTransfer.CreatedAt, fetchedTransfer.CreatedAt)
}

func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID: account1.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}