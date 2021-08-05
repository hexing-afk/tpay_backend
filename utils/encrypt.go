package utils

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/wiwii/pkcs7"
)

// md5计算
func Md5(src string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

// md5+Salted 计算
func Md5WithSalt(src, salt string) string {
	return Md5(src + salt)
}

// sha256散列
func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func SignPKCS7(priv *rsa.PrivateKey, cert, hashed []byte, hash crypto.Hash) ([]byte, error) {
	//Initialize a SignedData struct with content to be signed
	signedData, err := pkcs7.NewSignedData(hashed, hash)
	if err != nil {
		fmt.Printf("Cannot initialize signed data: %s", err)
	}

	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		log.Printf("err=%v\n", "private key error!")
		return nil, errors.New("PrivateKeyError")
	}
	cer, _ := x509.ParseCertificate(block.Bytes)

	// Add the signing cert and private key
	if err := signedData.AddSigner(cer, priv, pkcs7.SignerInfoConfig{}, hash); err != nil {
		fmt.Printf("Cannot add signer: %s", err)
	}

	// Call Detach() is you want to remove content from the signature
	// and generate an S/MIME detached signature
	signedData.Detach()

	// Finish() to obtain the signature bytes
	detachedSignature, err := signedData.Finish()
	//fmt.Printf("detachedSignature=[%v]", base64.StdEncoding.EncodeToString(detachedSignature))
	if err != nil {
		fmt.Printf("Cannot finish signing data: %s", err)
	}

	return detachedSignature, nil
}

func loopEncryptPKCS1v15(rand io.Reader, pub *rsa.PublicKey, msg []byte) (ret []byte, err error) {
	k := (pub.N.BitLen()+7)/8 - 11
	i := 0
	retArr := [][]byte{}
	for i = 0; i+k < len(msg); i = i + k {
		retI, err := rsa.EncryptPKCS1v15(rand, pub, msg[i:i+k])
		if err != nil {
			log.Printf("err=[%v]", err)
			return []byte{}, err
		}
		//log.Printf("retI=[%v]", retI)
		retArr = append(retArr, retI)
	}
	retI, err := rsa.EncryptPKCS1v15(rand, pub, msg[i:])
	if err != nil {
		log.Printf("err=[%v]", err)
		return []byte{}, err
	}
	//log.Printf("retI=[%v]", retI)
	retArr = append(retArr, retI)
	ret = bytes.Join(retArr, []byte{})
	return ret, nil
}

func loopDecryptPKCS1v15(rand io.Reader, priv *rsa.PrivateKey, msg []byte) (ret []byte, err error) {
	k := priv.N.BitLen() / 8
	i := 0
	retArr := [][]byte{}
	for i = 0; i+k < len(msg); i = i + k {
		retI, err := rsa.DecryptPKCS1v15(rand, priv, msg[i:i+k])
		if err != nil {
			log.Printf("err=[%v]", err)
			return []byte{}, err
		}
		//log.Printf("retI=[%v]", retI)
		retArr = append(retArr, retI)
	}
	retI, err := rsa.DecryptPKCS1v15(rand, priv, msg[i:])
	if err != nil {
		log.Printf("err=[%v]", err)
		return []byte{}, err
	}
	//log.Printf("retI=[%v]", retI)
	retArr = append(retArr, retI)
	ret = bytes.Join(retArr, []byte{})
	return ret, nil
}
