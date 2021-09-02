package chain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress(t *testing.T) {
	addr, err := CreateAddress()
	require.NoError(t, err)

	t.Log(addr.Address)
	t.Log(addr.PrivateKey)

}
