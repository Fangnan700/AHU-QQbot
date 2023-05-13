package ahu

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"math/big"
)

func AhuPasswdEncode(et string, nt string, password string) string {

	eBytes, err := hex.DecodeString(et)
	if err != nil {

	}
	e := new(big.Int).SetBytes(eBytes)

	nBytes, err := hex.DecodeString(nt)
	if err != nil {

	}
	n := new(big.Int).SetBytes(nBytes)

	rsaPublicKey := &rsa.PublicKey{N: n, E: int(e.Int64())}

	passwordBytes := []byte(password)
	passwordEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, passwordBytes)
	if err != nil {

	}
	passwordHex := hex.EncodeToString(passwordEncrypted)

	return passwordHex
}
