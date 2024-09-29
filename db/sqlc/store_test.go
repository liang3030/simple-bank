package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("before ", account1.Balance, "account2: ", account2.Balance)

	// run a concurrent transfer transaction
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// TODO: txName for test
		// txName := fmt.Sprintf("tx %d", i+1),
		// TODO: check concurrent in golang
		go func() {
			// TODO: what is txKey here?
			// ctx := context.WithValue(context.Background(), txKey, txName)
			// result, err := store.TransferTx(ctx), TransferTxParams{
			result, err := store.TransferTx(context.Background(), TransferTxParams{

				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		transfer := result.Transfer
		require.NotEmpty(t, result)
		// require.Equal(t, account1.ID, result.FromAccount.ID)
		// require.Equal(t, account2.ID, result.ToAccount.ID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.Equal(t, transfer.Amount, amount)

		_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// TODO: check account balances

		// check accounts' balance
		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		fmt.Println("during0, from: ", fromAccount.Balance, "to: ", toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

		fmt.Println("during, from: ", fromAccount.Balance, "to: ", toAccount.Balance)
	}

	// check the final updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

	fmt.Println("after, from: ", updateAccount1.Balance, "account2: ", updateAccount2.Balance)

}
