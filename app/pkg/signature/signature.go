package signature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
)

type Signature struct {
	R []byte `json:"r"`
	S []byte `json:"s"`
}

func VerifySignature(signature Signature, data []byte, publicKey *ecdsa.PublicKey) (bool, error) {
	// Hash data menggunakan SHA-256
	hash := sha256.Sum256(data)

	// Konversi tanda tangan ke *big.Int
	r := new(big.Int).SetBytes(signature.R)
	s := new(big.Int).SetBytes(signature.S)

	// Verifikasi tanda tangan menggunakan kunci publik
	isValid := ecdsa.Verify(publicKey, hash[:], r, s)
	return isValid, nil
}

func SignData(privateKey *ecdsa.PrivateKey, message string) (r, s []byte, err error) {
	// Hash pesan menggunakan SHA-256
	hash := sha256.Sum256([]byte(message))

	// Tanda tangani hash menggunakan kunci privat
	rInt, sInt, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, nil, err
	}

	// Konversi tanda tangan menjadi byte array
	r = rInt.Bytes()
	s = sInt.Bytes()
	return r, s, nil
}

func SerializePublicKey(pubKey *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pubKey)
}

func DeserializePublicKey(data []byte) (*ecdsa.PublicKey, error) {
	pubKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}
	return ecdsaPubKey, nil
}

func LoadOrCreateKeyPair(filename string) (*ecdsa.PrivateKey, error) {
	// Check if the file exists
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		// If not, create a new private key
		fmt.Println("No key found. Generating a new key pair...")
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		// Save the private key to a file
		err = savePrivateKey(filename, privateKey)
		if err != nil {
			return nil, err
		}
		return privateKey, nil
	}

	// Load the existing private key
	privateKey, err := loadPrivateKey(filename)
	if err != nil {
		return nil, err
	}
	fmt.Println("Loaded existing key pair.")
	return privateKey, nil
}

// savePrivateKey saves a private key to a PEM file
func savePrivateKey(filename string, privateKey *ecdsa.PrivateKey) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}

	// Encode as PEM
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.Encode(file, pemBlock)
}

// loadPrivateKey loads a private key from a PEM file
func loadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Decode PEM
	block, _ := pem.Decode(file)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// Parse the EC private key
	return x509.ParseECPrivateKey(block.Bytes)
}
