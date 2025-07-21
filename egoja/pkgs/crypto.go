package pkgs

import (
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/deatil/go-cryptobin/cryptobin/rsa"
	"github.com/deatil/go-cryptobin/cryptobin/sm2"
	"goframe-ex/egoja/gojaapi"
	"goframe-ex/encrypt/cryptobin"
)

func init() {
	gojaapi.RegisterImport("cryptobin", map[string]any{
		"FromString":       crypto.FromString,
		"FromBase64String": crypto.FromBase64String,
		"FromBytes":        crypto.FromBytes,
		"FromHexString":    crypto.FromHexString,
		//"Cryptobin":       crypto.Cryptobin{},
	})

	gojaapi.RegisterImport("cryptobin/sm2", map[string]any{
		"NewSM2": sm2.NewSM2,
		//"SM2":   (*sm2.SM2)(nil),
	})

	gojaapi.RegisterImport("cryptobin/rsa", map[string]any{
		"New":                    rsa.New,
		"RsaBigDataEncrypt":      cryptobin.RsaBigDataEncrypt,
		"RsaBigDataDecrypt":      cryptobin.RsaBigDataDecrypt,
		"RsaBigDataEncryptByPri": cryptobin.RsaBigDataEncryptByPri,
		"RsaBigDataDecryptByPub": cryptobin.RsaBigDataDecryptByPub,
		//"RSA":                   (*rsa.RSA)(nil),
	})
}
