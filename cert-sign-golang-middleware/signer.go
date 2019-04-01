package function

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"strings"
	"time"
)

type CSR struct {
	// Host, a comma-separated hostnames and IPs to generate a certificate for
	Host string

	// ValidFrom, (optional) the creation date formatted as Jan 1 15:04:05 2011. Defaults to the current date/time.
	ValidFrom string

	// ValidFro
	ValidFor time.Duration

	// IsCertificateAuthority, whether this cert should be its own Certificate Authority
	IsCertificateAuthority bool

	// Size of RSA key to generate. Ignored if --ecdsa-curve is set
	RSAKeySize int

	// ECDSA, (optional) the ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521
	ECDSACurve string
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) (*pem.Block, error) {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal ECDSA private key: %v", err)

		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, nil
	}
}

func Sign(r *CSR) (privateKey string, publicCertificate string, err error) {
	flag.Parse()

	if len(r.Host) == 0 {
		return "", "", fmt.Errorf("missing required --host parameter")
	}

	var priv interface{}
	switch r.ECDSACurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, r.RSAKeySize)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return "", "", fmt.Errorf("uinrecognized elliptic curve: %q", r.ECDSACurve)

	}
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}

	var notBefore time.Time
	if len(r.ValidFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", r.ValidFrom)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse creation date: %s\n", err)
		}
	}

	notAfter := notBefore.Add(r.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(r.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if r.IsCertificateAuthority {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	certOut := strings.Builder{}
	if err := pem.Encode(&certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return "", "", fmt.Errorf("failed to write data to cert.pem: %s", err)
	}

	keyOut := strings.Builder{}

	privKeyPEM, err := pemBlockForKey(priv)
	if err != nil {
		return "", "", fmt.Errorf("error creating private key PEM from key. %v", err)
	}

	if err := pem.Encode(&keyOut, privKeyPEM); err != nil {
		log.Fatalf("failed to write data to key.pem: %s", err)
	}

	return keyOut.String(), certOut.String(), nil
}
