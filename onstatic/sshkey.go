package onstatic

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/go-git/go-billy/v5"
	"github.com/morikuni/failure"
	"github.com/sters/onstatic/conf"
	"golang.org/x/crypto/ssh"
)

// key wrapper ed25519 key.
type key struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// Save private and public keys to directory.
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

	if err := k.savePrivateKey(fs); err != nil {
		return failure.Wrap(err)
	}
	if err := k.savePublicKey(fs); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (k *key) savePrivateKey(fs billy.Filesystem) (err error) {
	f, err := fs.Create(conf.Variables.SSHKeyFilename)
	if err != nil {
		return failure.Wrap(err)
	}
	defer func() {
		err = failure.Wrap(f.Close())
	}()

	b, err := x509.MarshalPKCS8PrivateKey(k.privateKey)
	if err != nil {
		return failure.Wrap(err)
	}

	return pem.Encode(f, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	})
}

func (k *key) savePublicKey(fs billy.Filesystem) (err error) {
	f, err := fs.Create(conf.Variables.SSHPubKeyFilename)
	if err != nil {
		return failure.Wrap(err)
	}
	defer func() {
		err = failure.Wrap(f.Close())
	}()

	pk, err := ssh.NewPublicKey(k.publicKey)
	if err != nil {
		return failure.Wrap(err)
	}

	_, err = f.Write(ssh.MarshalAuthorizedKey(pk))
	return failure.Wrap(err)
}

// GenerateKey returns new Key instance
// TODO: pass phrease.
func generateKey(size int, privateKeyName string, publicKeyName string) (*key, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	testMessage := []byte("hello world!")
	signed := ed25519.Sign(privateKey, testMessage)
	if res := ed25519.Verify(publicKey, testMessage, signed); !res {
		return nil, failure.Unexpected("failed to verify generated key")
	}

	return &key{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}
