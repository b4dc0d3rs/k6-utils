package k6utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

func (k6utils *K6Utils) GenerateKeyPair() (string, string) {
	curve := elliptic.P256()
	privateKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	publicKey := elliptic.Marshal(curve, privateKey.X, privateKey.Y)
	return base64.StdEncoding.EncodeToString(privateKey.D.Bytes()), base64.StdEncoding.EncodeToString(publicKey)
}

func (K6utils *K6Utils) SignData(data string, privKeyStr string, pubKeyStr string) (string, error) {
	dataBytes := []byte(data)
	signedDataByte, err := SignDataLocal(dataBytes, privKeyStr, pubKeyStr)
	if err != nil {
		return "", fmt.Errorf("Error signing data")
	}

	signedData := base64.StdEncoding.EncodeToString(signedDataByte)
	return signedData, nil
}

// SignData signs the data using the ECDSA private key and SHA256 hashing algorithm
func SignDataLocal(data []byte, privKeyStr string, pubKeyStr string) ([]byte, error) {
	curve := elliptic.P256()

	privKeyBytes, err := base64.StdEncoding.DecodeString(privKeyStr)
	if err != nil {
		return nil, fmt.Errorf("Error decoding private key: %v", err)
	}

	privKey := new(ecdsa.PrivateKey)
	privKey.Curve = curve
	privKey.D = new(big.Int).SetBytes(privKeyBytes)

	// Might want to delete later
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, fmt.Errorf("Error decoding public key: %v", err)
	}

	pubKeyX, pubKeyY := elliptic.Unmarshal(curve, pubKeyBytes)
	if pubKeyX == nil {
		return nil, fmt.Errorf("Error unmarshalling public key")
	}

	privKey.PublicKey = ecdsa.PublicKey{Curve: curve, X: pubKeyX, Y: pubKeyY}

	// Hash the data
	h := sha256.New()
	h.Write(data)
	hash := h.Sum(nil)

	// Sign the hash using the private key
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		return nil, err
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	// Encode the signature to base64
	return signature, nil
}
