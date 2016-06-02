package battery_crypto

import (
	xycrypto "guanghuan.com/xiaoyao/common/crypto"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(int64(time.Now().UnixNano()))
	xycrypto.SetKeyTable(keys[:len(keys)])
}

func Encrypt(plain_data []byte) (cipher []byte, err error) {
	return xycrypto.Encrypt(plain_data)
}
func Decrypt(cipher []byte) (plain_data []byte, err error) {
	return xycrypto.Decrypt(cipher)
}
