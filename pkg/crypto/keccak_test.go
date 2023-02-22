package crypto

import (
	"testing"
)

func TestNewKeccak_padding(t *testing.T) {
	message := []byte("The quick brown fox jumps over the lazy dog")

	k := NewKeccak(message)

	k.padding()

	t.Logf("paddedMessage len: %+v", len(k.paddedMessage))

}