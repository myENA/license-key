package lk

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/sha3"
)

// ErrInvalidSignature indicates failed signature validation
var ErrInvalidSignature = errors.New("Signature validation failed.")

// ErrInvalidKey indicates an invalid key string
var ErrInvalidKey = errors.New("Invalid key format")

// secret should be changed using SetSecret(s string) before generating keys
var secret = "d8a4228a21cd4394c4b52c4a555e6e9dd0c9a6935753957da5d4429d54ccc995"

const (
	dataSize       = 16
	sumSize        = 8
	chunkSize      = 4
	chunkCount     = (dataSize + sumSize) / chunkSize
	chunkSeparator = `-`
)

// Key represents the internal key structure
type Key struct {
	data [dataSize]byte
	sum  [sumSize]byte
}

// String implements the stringer interface on the key object
func (k *Key) String() string {
	var s []string // output string

	// nil protection
	if k == nil {
		return ""
	}

	// encode data bytes
	for i := 0; i < dataSize; i += chunkSize {
		s = append(s, hex.EncodeToString(k.data[i:i+chunkSize]))
	}

	// encode sum bytes
	for i := 0; i < sumSize; i += chunkSize {
		s = append(s, hex.EncodeToString(k.sum[i:i+chunkSize]))
	}

	// join and return
	return strings.Join(s, chunkSeparator)
}

// New generates a signed key
func New() (*Key, error) {
	var k Key            // key object
	var h sha3.ShakeHash // hash holder
	var err error        // general error handler

	// read random bytes
	if _, err = rand.Read(k.data[:]); err != nil {
		return nil, err
	}

	// init hash
	h = sha3.NewShake128()

	// write salt and handle errors
	if _, err = h.Write([]byte(secret)); err != nil {
		return nil, err
	}

	// write data and handle errors
	if _, err = h.Write(k.data[:]); err != nil {
		return nil, err
	}

	// write sum
	if _, err = h.Read(k.sum[:]); err != nil {
		return nil, err
	}

	// return key - no error
	return &k, nil
}

// Parse parses a string key to a key object and validates resulting object
func Parse(s string) (*Key, error) {
	var k Key           // key object
	var buf []byte      // temp byte buffer
	var chunks []string // string chunks
	var err error       // general error handler

	// validate basic form
	chunks = strings.Split(s, chunkSeparator)

	// check length
	if len(chunks) != chunkCount {
		return nil, ErrInvalidKey
	}

	// decode data and check error
	if buf, err = hex.DecodeString(strings.Join(chunks[0:dataSize/chunkSize], "")); err != nil {
		return nil, ErrInvalidKey
	}

	// set data bytes
	copy(k.data[:], buf)

	// decode checksum and check error
	if buf, err = hex.DecodeString(strings.Join(chunks[dataSize/chunkSize:], "")); err != nil {
		return nil, ErrInvalidKey
	}

	// set checksum
	copy(k.sum[:], buf)

	// validate key
	if k.Validate() != true {
		return nil, ErrInvalidSignature
	}

	// return key object - no error
	return &k, nil
}

// Validate validates a key object
func (k *Key) Validate() bool {
	var h sha3.ShakeHash // hash holder
	var sum [8]byte      // checksum
	var err error        // error holder

	// nil protection
	if k == nil {
		return false
	}

	// init hash
	h = sha3.NewShake128()

	// write salt and handle errors
	if _, err = h.Write([]byte(secret)); err != nil {
		return false
	}

	// write data and handle errors
	if _, err = h.Write(k.data[:]); err != nil {
		return false
	}

	// read sum into buffer
	if _, err = h.Read(sum[:]); err != nil {
		return false
	}

	// return checksum equality
	return bytes.Equal(sum[:], k.sum[:])
}

// SetSecret sets the secret salt
func SetSecret(s string) {
	secret = s
}

// Secret returns the secret salt
func Secret() string {
	return secret
}
