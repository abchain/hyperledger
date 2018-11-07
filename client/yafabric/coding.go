package client

import (
	"encoding/base64"
)

func toArgument(b []byte) string {
	return base64.RawStdEncoding.EncodeToString(b)
}

func FromArgument(arg string) []byte {
	b, err := base64.RawStdEncoding.DecodeString(arg)
	if err != nil {
		return nil
	}

	return b
}
