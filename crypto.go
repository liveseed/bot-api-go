package bot

import (
	"crypto/ed25519"
	"crypto/sha512"

	"filippo.io/edwards25519"
)

func PrivateKeyToCurve25519(curve25519Private *[32]byte, privateKey ed25519.PrivateKey) {
	h := sha512.New()
	h.Write(privateKey.Seed())
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(curve25519Private[:], digest)
}

func PublicKeyToCurve25519(publicKey ed25519.PublicKey) ([]byte, error) {
	p, err := (&edwards25519.Point{}).SetBytes(publicKey[:])
	if err != nil {
		return nil, err
	}
	return p.BytesMontgomery(), nil
}
