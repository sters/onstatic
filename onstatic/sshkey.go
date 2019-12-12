package onstatic

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/morikuni/failure"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-billy.v4"
)

// key wrapper rsa key
type key struct {
	*rsa.PrivateKey
	privateKeyName string
	publicKeyName  string
}

// Save private and public keys to directory
func (k *key) save(fs billy.Filesystem, dir string) error {
	if s, err := fs.Stat(dir); err != nil {
		return failure.Wrap(err, failure.Message("failed to stat"))
	} else if !s.IsDir() {
		return fmt.Errorf("failed to load directory")
	}

	fs, err := fs.Chroot(dir)
	if err != nil {
		return failure.Wrap(err)
	}

	if err := k.savePrivateKey(fs, k.privateKeyName); err != nil {
		return failure.Wrap(err)
	}
	if err := k.savePublicKey(fs, k.publicKeyName); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (k *key) savePrivateKey(fs billy.Filesystem, filename string) error {
	f, err := fs.Create(filename)
	defer f.Close()
	if err != nil {
		return failure.Wrap(err)
	}

	b, err := x509.MarshalPKCS8PrivateKey(k.PrivateKey)
	if err != nil {
		return failure.Wrap(err)
	}

	return pem.Encode(f, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	})
}

func (k *key) savePublicKey(fs billy.Filesystem, filename string) error {
	f, err := fs.Create(filename)
	defer f.Close()
	if err != nil {
		return failure.Wrap(err)
	}

	pk, err := ssh.NewPublicKey(&k.PublicKey)
	if err != nil {
		return failure.Wrap(err)
	}

	_, err = f.Write(ssh.MarshalAuthorizedKey(pk))
	return failure.Wrap(err)
}

// GenerateKey returns new Key instance
// TODO: pass phrease
func generateKey(size int, privateKeyName string, publicKeyName string) (*key, error) {
	k, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	if err = k.Validate(); err != nil {
		return nil, failure.Wrap(err)
	}

	return &key{
		PrivateKey:     k,
		privateKeyName: privateKeyName,
		publicKeyName:  publicKeyName,
	}, nil
}
