package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
)

// Encoding implements encrypt/decryt alg
type Encoding struct {
	encoding *base64.Encoding
	block    cipher.Block
	key      []byte
}

func NewEncoding(block, key string) *Encoding {
	// use cookie name as cipher key
	hash32 := sha256.New()
	hash32.Write([]byte(block))

	block32, _ := aes.NewCipher(hash32.Sum(nil)) // 32-bytes key

	// use cookie secret as crypto key
	hash16 := md5.New()
	hash16.Write([]byte(key))

	key16 := hash16.Sum(nil)

	return &Encoding{
		encoding: base64.RawURLEncoding,
		block:    block32,
		key:      key16,
	}
}

func (ec *Encoding) Encrypt(s string) string {
	data := []byte(s)
	buf := make([]byte, len(data))

	cfb := cipher.NewCFBEncrypter(ec.block, ec.key)
	cfb.XORKeyStream(buf, data)

	return ec.encoding.EncodeToString(buf)
}

func (ec *Encoding) Decrypt(s string) (string, error) {
	data, err := ec.encoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	buf := make([]byte, len(data))

	cfb := cipher.NewCFBDecrypter(ec.block, ec.key)
	cfb.XORKeyStream(buf, data)

	return string(buf), nil
}
