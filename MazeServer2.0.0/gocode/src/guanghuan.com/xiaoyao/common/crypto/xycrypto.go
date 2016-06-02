// xiaoyao cryption package
// 	2014.4.21
/*
	Package xycrypto: a wrapper of cryption methods
*/
package xycrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	//	"fmt"
	//	xylog "guanghuan.com/xiaoyao/common/log"
	"io"
)

var (
	ErrInvalidKeyTable = errors.New("Invalid Key table")      // 错误的秘钥表
	ErrInvalidKey      = errors.New("Invalid Key")            // 错误的秘钥
	ErrInvalidIv       = errors.New("Invalid Initial Vector") // 错误的初始化向量
)

// 一个加密块的大小，
// 16字节
func BlockSize() int {
	return 16
}

// 秘钥的长度，
// 我们采用256位AES加密模式，所以秘钥长度是32字节
func KeySize() int {
	return 32
}

// AES 加密/解密过程 (采用ofb模式，所以加密解密的过程是相同的)
// 	参考：http://en.wikipedia.org/wiki/Block_cipher_modes_of_operation#Cipher_feedback_.28CFB.29
func AES_ofb_crypt(in_buf []byte, key []byte, iv []byte) (out_buf []byte, err error) {

	block, err := aes.NewCipher(key)
	out_buf = make([]byte, len(in_buf))

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(out_buf, in_buf)

	return
}

// 加密过程
// 	约定初始化向量保存在密文首部
func AESEncrypt(src_buf []byte, key []byte, key_idx byte) (out_buf []byte, err error) {
	if len(key) != KeySize() {
		err = ErrInvalidKey
		return
	}
	out_buf = make([]byte, aes.BlockSize+len(src_buf))

	// 初始化向量约定保存在密文首部
	iv := out_buf[:aes.BlockSize]
	// 随机
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	iv[0] = byte(key_idx)
	enc := make([]byte, len(src_buf))

	enc, err = AES_ofb_crypt(src_buf, key, iv)

	copy(out_buf[aes.BlockSize:], enc)

	return
}

// 解密过程
func AESDecrypt(enc_buf []byte, key []byte) (dec_buf []byte, err error) {
	//	xylog.Debug("buf: %x", enc_buf)
	if len(key) != KeySize() {
		err = ErrInvalidKey
		return
	}
	dec_buf = make([]byte, len(enc_buf)-aes.BlockSize)
	// 初始化向量约定保存在密文首部
	iv := enc_buf[:aes.BlockSize]
	src_buf := enc_buf[aes.BlockSize:]

	dec_buf, err = AES_ofb_crypt(src_buf, key, iv)

	return
}

/*
func aes_ofb_demo(key16 []byte, plaintext []byte) {
	key := key16 //[]byte("example key 1234")
	//	plaintext := []byte("some plaintext")

	fmt.Printf("key: %x\n", string(key))
	fmt.Printf("in : %x\n", plaintext)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := []byte("1111111111111111")
	//	iv := ciphertext[:aes.BlockSize]
	//	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//		panic(err)
	//	}
	fmt.Printf("iv : %x\n", iv)

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	fmt.Printf("enc: %x\n", ciphertext[aes.BlockSize:])
	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	// OFB mode is the same for both encryption and decryption, so we can
	// also decrypt that ciphertext with NewOFB.

	plaintext2 := make([]byte, len(plaintext))
	stream = cipher.NewOFB(block, iv)
	stream.XORKeyStream(plaintext2, ciphertext[aes.BlockSize:])

	fmt.Printf("dec: %x\n", plaintext2)

	fmt.Printf("value : [%s]\n", plaintext2)
}
*/
