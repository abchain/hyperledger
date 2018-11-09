package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/ripemd160"
)

func SHA256RIPEMD160(data []byte) ([]byte, error) {

	hash1pass := sha256.Sum256(data)
	if len(hash1pass) != sha256.Size {
		return nil, errors.New("Wrong sha256 hashing")
	}

	rmd160h := ripemd160.New()
	if nn, err := rmd160h.Write(hash1pass[:]); nn != len(hash1pass) || err != nil {
		return nil, errors.New("Wrong ripemd write")
	}

	hash2pass := rmd160h.Sum([]byte{})
	if len(hash2pass) != ripemd160.Size {
		return nil, errors.New("Wrong ripemd160 hashing")
	}

	return hash2pass[:], nil
}

func DoubleSHA256(data []byte) ([]byte, error) {

	hash1pass := sha256.Sum256(data)
	if len(hash1pass) != sha256.Size {
		return nil, errors.New("Wrong sha256 hashing 1pass")
	}

	hash2pass := sha256.Sum256(hash1pass[:])
	if len(hash2pass) != sha256.Size {
		return nil, errors.New("Wrong sha256 hashing 2pass")
	}

	return hash2pass[:], nil
}

func HMACSHA256(data []byte, seed []byte) ([]byte, error) {
	hasher := hmac.New(sha256.New, seed)

	_, err := hasher.Write(data)
	if err != nil {
		return nil, errors.New("Wrong hmac write")
	}

	return hasher.Sum(nil), nil
}
