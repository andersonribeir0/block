package internal

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"time"
)

type Block struct {
	hash      []byte
	data      []byte
	prevHash  []byte
	signature []byte
	pubKey    []byte
	nonce     uint64
	timestamp time.Time
}

func NewBlock(data []byte, kg *KeyGen) (*Block, error) {
	block := &Block{data: data, timestamp: time.Now()}

	if err := kg.Sign(block.data); err != nil {
		return nil, err
	}

	block.signature = append(kg.r.Bytes(), kg.s.Bytes()...)
	block.pubKey = elliptic.Marshal(elliptic.P256(), kg.key.PublicKey.X, kg.key.PublicKey.Y)

	return block, nil
}

func (b *Block) CalculateHash() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(b.prevHash)

	err := binary.Write(buf, binary.LittleEndian, byte(b.timestamp.UnixNano()))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, byte(b.nonce))
	if err != nil {
		return nil, err
	}

	buf.Write(b.data)

	hashedBlock := sha256.Sum256(buf.Bytes())

	return hashedBlock[:], nil
}

func (b *Block) MineBlock(difficulty int) error {
	var intHash big.Int
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	for b.nonce = 0; ; b.nonce++ {
		hash, err := b.CalculateHash()
		if err != nil {
			return err
		}
		intHash.SetBytes(hash)

		if intHash.Cmp(target) == -1 {
			b.hash = hash
			break
		}
	}

	return nil
}

func (b *Block) IsValidHash(difficulty int) bool {
	var intHash big.Int

	intHash.SetBytes(b.hash)

	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	return intHash.Cmp(target) == -1
}

func (b *Block) IsValid(difficulty int, kg *KeyGen) bool {
	if !b.IsValidHash(difficulty) {
		return false
	}

	r := big.NewInt(0).SetBytes(b.signature[:len(b.signature)/2])
	s := big.NewInt(0).SetBytes(b.signature[len(b.signature)/2:])

	x, y := elliptic.Unmarshal(elliptic.P256(), b.pubKey)
	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	return ecdsa.Verify(&publicKey, b.data, r, s)
}
