package crypto


const (
	rateInBytes  = 136 // number of bytes in a Keccak-256 block
	hashSize     = 32  // output size in bytes
	paddingValue = 0x06
	stateSizeBytes = 25 * 8 // number of bytes in the Keccak-256 state
	stateSizeWords = stateSizeBytes / 8
)


func Keccak256(data []byte) []byte {
	var state [stateSizeWords]uint64
	block := make([]byte, rateInBytes)

	// Absorb the input data
	for len(data) >= len(block) {
		copy(block, data[:len(block)])
		data = data[len(block):]
		absorbBlock(state, block)
	}

	// Add padding and absorb the last block
	copy(block, data)
	block[len(data)] = paddingValue
	block[rateInBytes-1] |= 0x80
	absorbBlock(state, block)

	// Squeeze the output
	hash := make([]byte, hashSize)
	squeeze(state, hash)
	return hash
}

func absorbBlock(state [25]uint64, block []byte) {
	// XOR the block into the state
	for i := 0; i < rateInBytes/8; i++ {
		state[i] ^= bytesToUint64(block[i*8 : (i+1)*8])
	}

	// Permute the state
	keccakF1600(state)
}

func squeeze(state [25]uint64, out []byte) {
	// Squeeze full blocks
	for len(out) >= rateInBytes {
		out = uint64ToBytes(state[:rateInBytes/8])
		out = out[rateInBytes:]
		keccakF1600(state)
	}

	// Squeeze the last partial block
	if len(out) > 0 {
		block := make([]byte, rateInBytes)
		copy(block, out)
		block[len(out)] = paddingValue
		block[rateInBytes-1] |= 0x80
		absorbBlock(state, block)
		uint64ToBytes(out, state[:rateInBytes/8])
	}
}

func keccakF1600(state [25]uint64) {
	var C [5]uint64
	var D [5]uint64
	var roundConstants = [24]uint64{
		0x0000000000000001, 0x0000000000008082, 0x800000000000808a, 0x8000000080008000,
		0x000000000000808b, 0x0000000080000001, 0x8000000080008081, 0x8000000000008009,
		0x000000000000008a, 0x0000000000000088, 0x0000000080008009, 0x000000008000000a,
		0x000000008000808b, 0x800000000000008b, 0x8000000000008089, 0x8000000000008003,
		0x8000000000008002, 0x8000000000000080, 0x000000000000800a, 0x800000008000000a,
		0x8000000080008081, 0x8000000000008080, 0x0000000080000001, 0x8000000080008008,
	}


	for i := 0; i < 24; i++ {
		// Theta step
		for x := 0; x < 5; x++ {
			C[x] = state[x] ^ state[x+5] ^ state[x+10] ^ state[x+15] ^ state[x+20]
		}
		for x := 0; x < 5; x++ {
			D[x] = C[(x+4)%5] ^ rotateLeft(C[(x+1)%5], 1)
		}
		for x := 0; x < 5; x++ {
			for y := 0; y < 25; y += 5 {
				state[x+y] ^= D[x]
			}
		}

		// Rho and Pi steps
		var temp uint64
		x, y := 1, 0
		for t := 0; t < 24; t++ {
			temp = state[x+y*5]
			state[x+y*5] = rotateLeft(D[x], (t+1)*(t+2)/2)
			D[x] = temp
			x, y = y, (2*x+3*y)%5
		}

		// Chi step
		for y := 0; y < 25; y += 5 {
			for x := 0; x < 5; x++ {
				C[x] = state[x+y]
			}
			for x := 0; x < 5; x++ {
				state[x+y] = state[x+y] ^ (^C[(x+1)%5] & C[(x+2)%5])
			}
		}

		// Iota step
		state[0] ^= roundConstants[i]
	}
}

func rotateLeft(x uint64, n int) uint64 {
	return (x << n) | (x >> (64 - n))
}


func bytesToUint64(b []byte) (uint64) {
	if len(b) != 8 {
		return 0
	}

	var u uint64
	for i := 0; i < 8; i++ {
		u |= uint64(b[i]) << ((7 - i) * 8)
	}

	return u
}

func uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(u >> ((7 - i) * 8))
	}
	return b
}
