package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	//"fmt"
)

// GenerateKeyPairs generates an RSA key pair with the given key size and returns the PEM encoded public and private keys.
func RSA_GenerateKeyPairs(bits int) (*string, *string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	privateKeyStr := string(privateKeyPEM)
	publicKeyStr := string(publicKeyPEM)
	return &publicKeyStr, &privateKeyStr, nil
}

// Encrypt encrypts the given message with the provided PEM encoded public key.
func RSA_Encrypt(pubKeyPEM *string, message *string) (*[]byte, error) {
	block, _ := pem.Decode([]byte(*pubKeyPEM))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("invalid public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(*message))
	if err != nil {
		return nil, err
	}

	return &encryptedMessage, nil
}

// Decrypt decrypts the given encrypted message with the provided PEM encoded private key.
func RSA_Decrypt(privKeyPEM *string, encryptedMessage *[]byte) (*string, error) {
	block, _ := pem.Decode([]byte(*privKeyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, *encryptedMessage)
	if err != nil {
		return nil, err
	}

	decryptedMessageStr := string(decryptedMessage)
	return &decryptedMessageStr, nil
}
