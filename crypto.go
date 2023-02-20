package main


type Crypto interface {
	calculateHash() []byte
}
