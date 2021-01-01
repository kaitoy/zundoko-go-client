package util

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// UUID represents UUID
type UUID [16]byte

func (uuid UUID) String() string {
	buf := make([]byte, 36)
	hex.Encode(buf[:8], uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])
	return string(buf)
}

// NewUUID generates a new UUID that complies with UUIDv4
func NewUUID() UUID {
	var uuid UUID
	io.ReadFull(rand.Reader, uuid[:])
	uuid[6] = (uuid[6] & 0x0F) | 0x40
	uuid[8] = (uuid[8] & 0x3F) | 0x80
	return uuid
}
