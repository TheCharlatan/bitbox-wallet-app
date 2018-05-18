package addresses_test

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/shiftdevices/godbb/backend/coins/btc/addresses"
	"github.com/shiftdevices/godbb/backend/coins/btc/addresses/test"
	"github.com/shiftdevices/godbb/backend/signing"
	"github.com/stretchr/testify/require"
)

var scriptTypes = []signing.ScriptType{
	signing.ScriptTypeP2PKH,
	signing.ScriptTypeP2WPKHP2SH,
	signing.ScriptTypeP2WPKH,
}

func TestSigScriptWitnessSize(t *testing.T) {
	// A signature can be 70 or 71 bytes (excluding sighash op).
	// We take one that has 71 bytes, as the size function returns the maximum possible size.
	sigBytes, err := hex.DecodeString(
		`3045022100a97dc23e47bb79dbff73e33be4a4e476d6ef67c8c23a9ee4a9ee21f4dd80f0f202201c5d4be437308539e1193d9118fae03bae1942e9ce27c86803bb5f18aa044a46`)
	require.NoError(t, err)
	sig, err := btcec.ParseDERSignature(sigBytes, btcec.S256())
	require.NoError(t, err)

	for _, scriptType := range scriptTypes {
		t.Run(string(scriptType), func(t *testing.T) {
			address := test.GetAddress(scriptType)
			sigScriptSize, hasWitness := addresses.SigScriptWitnessSize(address.Configuration)
			sigScript, witness := address.SignatureScript([]*btcec.Signature{sig})
			require.Equal(t, len(sigScript), sigScriptSize)
			require.Equal(t, witness != nil, hasWitness)
		})
	}
}
