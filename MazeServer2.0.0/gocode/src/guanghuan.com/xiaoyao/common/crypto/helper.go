package xycrypto

import (
	//	"log"
	xylog "guanghuan.com/xiaoyao/common/log"
	"math/rand"
	"time"
)

func Init() {
	rand.Seed(time.Now().UnixNano())
}

type CryptoHelper struct {
	key_table []byte
}

var (
	DefCrypto CryptoHelper
)

func (helper *CryptoHelper) SetKeyTable(key_table []byte) {
	helper.key_table = make([]byte, len(key_table))
	helper.key_table = key_table
}

func (helper *CryptoHelper) Key(idx int, size int) []byte {
	return helper.key_table[idx : idx+size]
}
func (helper *CryptoHelper) KeySize() int {
	return KeySize()
}
func (helper *CryptoHelper) Encrypt(plain_data []byte) (cipher []byte, err error) {
	idx := rand.Intn(256)
	cipher, err = AESEncrypt(plain_data, helper.Key(idx, KeySize()), byte(idx))
	return
}

func (helper *CryptoHelper) Decrypt(cipher []byte) (plain_data []byte, err error) {
	idx := int(cipher[0])
	plain_data, err = AESDecrypt(cipher, helper.Key(idx, KeySize()))
	return
}

func SetKeyTable(key_table []byte) {
	DefCrypto.SetKeyTable(key_table)
}
func Encrypt(plain_data []byte) (cipher []byte, err error) {
	if len(DefCrypto.key_table) < DefCrypto.KeySize() {
		xylog.ErrorNoId("Key Table is invalid, size: %d", len(DefCrypto.key_table))
		err = ErrInvalidKeyTable
		return
	}
	return DefCrypto.Encrypt(plain_data)
}
func Decrypt(cipher []byte) (plain_data []byte, err error) {
	if len(DefCrypto.key_table) < DefCrypto.KeySize() {
		xylog.ErrorNoId("Key Table is invalid, size: %d", len(DefCrypto.key_table))
		err = ErrInvalidKeyTable
		return
	}
	return DefCrypto.Decrypt(cipher)
}
