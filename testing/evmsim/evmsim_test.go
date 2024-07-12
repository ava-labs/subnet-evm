package evmsim

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestAddr(t *testing.T) {
	keys := []string{
		"0x56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027",
		"0x750839e9dbbd2a0910efe40f50b2f3b2f2f59f5580bb4b83bd8c1201cf9a010a",
	}
	b := NewWithHexKeys(t, keys)

	var got []common.Address
	for i := range keys {
		got = append(got, b.Addr(i))
	}

	want := []common.Address{
		// The above private keys were used in Hardhat, resulting in these
		// addresses (i.e. we have equivalence with ethers.js).
		common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"),
		common.HexToAddress("0x0Fa8EA536Be85F32724D57A37758761B86416123"),
	}
	assert.Equal(t, want, got, "NewWithHexKeys().Addr([...])")
}

func FuzzNewWithNumKeys(f *testing.F) {
	f.Add(uint8(3), uint8(5))
	f.Add(uint8(10), uint8(10))

	f.Fuzz(func(t *testing.T, n1, n2 uint8) {
		min := n1
		if n2 < n1 {
			min = n2
		}

		b1 := NewWithNumKeys(t, uint(n1))
		b2 := NewWithNumKeys(t, uint(n2))

		t.Run("same prefix of addresses", func(t *testing.T) {
			for i := 0; i < int(min); i++ {
				assert.Equalf(t, b1.Addr(i), b2.Addr(i), "Addr(%d) on %Ts with n=%d and n=%d deterministic keys", i, b1, n1, n2)
			}
		})
	})
}
func TestLockDeterministicKeys(t *testing.T) {
	// The actual value of the deterministic keys isn't important, as long as we
	// don't change them. This address is calculated through repeated hashing,
	// so we know that those before it must be correct too (assuming the same
	// algorithm is used for all keys).
	const n = 100
	b := NewWithNumKeys(t, n)
	assert.Equalf(t, b.Addr(n-1), common.HexToAddress("0x529B45f562B9399D67435868A7474175E3B99Be3"), "NewWithNumKeys().Addr(%d)", n-1)
}
