package parse

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/gregdhill/go-openrpc/types"
	"github.com/test-go/testify/require"
)

func TestParse(t *testing.T) {
	data, err := ioutil.ReadFile("openrpc.json")
	require.NoError(t, err)
	spec := types.NewOpenRPCSpec1()
	err = json.Unmarshal(data, spec)
	require.NoError(t, err)
	GetTypes(spec, spec.Objects)

	result := spec.Objects.Get("GetBlockByHashResult")
	require.NotNil(t, result)
	require.Len(t, result.GetKeys(), 1)

	result = spec.Objects.Get("EthBlockNumberResult")
	require.NotNil(t, result)
	require.Len(t, result.GetKeys(), 1)

	result = spec.Objects.Get("BlockNumber")
	require.NotNil(t, result)
	require.Len(t, result.GetKeys(), 1)

	result = spec.Objects.Get("EthAccountsResult")
	require.NotNil(t, result)
	bt := result.Get("Addresses")
	require.Equal(t, "Addresses", bt.Name)
	require.Equal(t, "[]string", bt.Type)

	result = spec.Objects.Get("EthGetTransactionReceiptResult")
	require.NotNil(t, result)
	bt = result.Get("Receipt")
	require.Equal(t, "Receipt", bt.Name)
	require.Equal(t, "Receipt", bt.Type)

	result = spec.Objects.Get("EthSyncingResult")
	require.NotNil(t, result)
	bt = result.Get("Syncing")
	require.Equal(t, "Syncing", bt.Name)
	require.Equal(t, "SyncStatus", bt.Type)

	result = spec.Objects.Get("EthGetTransactionByHashResult")
	require.NotNil(t, result)
	bt = result.Get("Transaction")
	require.Equal(t, "Transaction", bt.Name)

	result = spec.Objects.Get("Transaction")
	require.NotNil(t, result)
	bt = result.Get("BlockHash")
	require.Equal(t, "BlockHash", bt.Name)
}
