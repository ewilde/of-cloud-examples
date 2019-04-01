package function

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	type args struct {
		r *CSR
	}
	tests := []struct {
		name     string
		args     args
		evaluate func(privateKey string, publicKey string, signError error) error
	}{
		{
			name: "valid RSA based certificate",
			args: args{
				r: &CSR{
					Host:       "example.com",
					ValidFor:   time.Hour * 24 * 365 * 2,
					RSAKeySize: 2048,
				},
			},
			evaluate: func(privateKey string, publicKey string, signError error) error {
				if !strings.HasPrefix(privateKey, "-----BEGIN RSA PRIVATE KEY-----") {
					return fmt.Errorf("private key not expected. %s", privateKey)
				}

				if !strings.HasPrefix(publicKey, "-----BEGIN CERTIFICATE-----") {
					return fmt.Errorf("public key not expected. %s", publicKey)
				}

				return nil
			},
		},
		{
			name: "valid ECDSA based certificate",
			args: args{
				r: &CSR{
					Host:       "example.com",
					ValidFor:   time.Hour * 24 * 365 * 2,
					ECDSACurve: "P256",
					RSAKeySize: 256,
				},
			},
			evaluate: func(privateKey string, publicKey string, signError error) error {
				if !strings.HasPrefix(privateKey, "-----BEGIN EC PRIVATE KEY-----") {
					return fmt.Errorf("private key not expected. %s", privateKey)
				}

				if !strings.HasPrefix(publicKey, "-----BEGIN CERTIFICATE-----") {
					return fmt.Errorf("public key not expected. %s", publicKey)
				}

				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrivateKey, gotPublicCertificate, gotErr := Sign(tt.args.r)
			if err := tt.evaluate(gotPrivateKey, gotPublicCertificate, gotErr); err != nil {
				t.Errorf("Sign() %v", err)
			}
		})
	}
}
