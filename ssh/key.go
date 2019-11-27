package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// Key wrapper rsa key
type Key struct {
	*rsa.PrivateKey
	privateKeyName string
	publicKeyName  string
}

// Save private and public keys to directory
func (k *Key) Save(dir string) error {
	dir = filepath.Clean(dir)
	if s, err := os.Stat(dir); err != nil {
		return err
	} else if !s.IsDir() {
		return fmt.Errorf("failed to load directory: %s", dir)
	}

	if err := k.savePrivateKey(filepath.Join(dir, k.privateKeyName)); err != nil {
		return err
	}
	if err := k.savePublicKey(filepath.Join(dir, k.publicKeyName)); err != nil {
		return err
	}

	return nil
}

func (k *Key) savePrivateKey(filename string) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	return pem.Encode(f, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(k.PrivateKey),
	})
}

func (k *Key) savePublicKey(filename string) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	pk, err := ssh.NewPublicKey(k.PublicKey)
	if err != nil {
		return err
	}

	_, err = f.Write(ssh.MarshalAuthorizedKey(pk))
	return err
}

// GenerateKey returns new Key instance
// TODO: pass phrease
func GenerateKey(size int, privateKeyName string, publicKeyName string) (*Key, error) {
	k, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}

	if err = k.Validate(); err != nil {
		return nil, err
	}

	return &Key{
		PrivateKey:     k,
		privateKeyName: privateKeyName,
		publicKeyName:  publicKeyName,
	}, nil
}
