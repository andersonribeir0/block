package internal

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"math/big"
)

type KeyGenOptsFunc func(*KeyGenOpts)

type KeyGenOpts struct {
	key  *ecdsa.PrivateKey
	s, r *big.Int
}

func defaultKeyGenOpts() KeyGenOpts {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return KeyGenOpts{
		key: privKey,
	}
}

func WithPrivateKey(privKey []byte) KeyGenOptsFunc {
	return func(opts *KeyGenOpts) {
		privateKey, err := x509.ParseECPrivateKey(privKey)
		if err != nil {
			panic(err)
		}
		opts.key = privateKey
	}
}

type KeyGen struct {
	KeyGenOpts
}

func NewKeyGen(opts ...KeyGenOptsFunc) *KeyGen {
	o := defaultKeyGenOpts()

	for _, fn := range opts {
		fn(&o)
	}

	return &KeyGen{
		KeyGenOpts: o,
	}
}

func (kg *KeyGen) Sign(payload []byte) error {
	r, s, err := ecdsa.Sign(rand.Reader, kg.key, payload)
	if err != nil {
		return err
	}

	kg.r = r
	kg.s = s

	return nil
}

func (kg *KeyGen) Verify(payload []byte) bool {
	return ecdsa.Verify(&kg.key.PublicKey, payload, kg.r, kg.s)
}
