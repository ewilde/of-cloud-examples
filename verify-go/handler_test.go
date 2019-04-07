package function

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strconv"
	"testing"
)

func TestVerify(t *testing.T) {
	//	t.Skip("Need better test data")
	body := "hello\n"

	publicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArMIYRx8eE0cdoZ5paukK
ipGJkte+kdYjpdXbg5MOrJgfnAjVVvtb2KbiNvr8tCWj4rV71MNOKHGZwXIT+Nus
zaZgJzi3YmaVRVZnc2oVGC/SMz5HUWA8mbj38XoO1Ok8y19Ggeh62D5o1N4QwMR7
e0hW+F8Z/NpahTRfLVg2VTI/LAi9SLP9zf0fysTFKzFfxQMrFiWH8oK8O0dye77k
p8EqYG3hMGt5awRzDcN/TSz19rTwQXrQf+KABn4EuOGi2e0AM2yGQmaMmWbr3QyR
LUUwFO1pfBMAawScP7SyclGp1urJPGrIEBIE7MoNuYddB0B449pH+apsrT2V6nah
MQIDAQAB
-----END PUBLIC KEY-----
`
	os.Setenv("Http_Method", "POST")
	os.Setenv("Http_Path", "/")
	os.Setenv("Http_Query", "")
	os.Setenv("Http_X_Forwarded_Host", "gateway:8080")
	os.Setenv("Http_Host", "verify-go:8080")
	os.Setenv("Http_ContentLength", strconv.Itoa(len(body)))
	os.Setenv("Http_Date", "Fri, 15 Mar 2019 19:43:40 GMT")
	os.Setenv("Http_Content_Type", "application/x-www-form-urlencoded")
	os.Setenv("Http_Digest", "WJG1tSLV3whtD/CxEPvZ0hu0/HFjrzTQgoai6Eb2vgM=")
	os.Setenv("Http_Authorization", `Signature keyId="callback",algorithm="rsa-sha256",headers="(request-target) host date content-type digest content-length",signature="qp0bUPrwtOKGaYFF35r1/eM8mgj6yPs48/CTX4NG3o6Y9NcJ7GpVTF8P9352Dp9y27UO81IF/ad4qUJG+VrkfUfsxRzyLyNsyMDOMDaV5S0XJTonIaf03WaEo4pJG9X5EisHp8+JYf5t+QGLIHZSq10L/HGKFvvEpfTVFu5iNeTl66dVJyStLmNnB2a227i288upRvzsV/HIEvYfCpRWpLrYNzie4tsIaIHlM+H7lBeY8YRzMvemC91NyxLDpEW0zGf/mqfk8xQBTreFudLvl7Z0bIWFSHslXhYEWLuvfkfLYMXlaW00zBuxMlr8kBNbV6zYSByHtNHG5HR3dTFjuw=="`)

	err := verify([]byte(body), publicKey)

	if err != nil {
		t.Error(err)
	}

}

const (
	testSpecPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`
)

func TestBasicRsaVerify(t *testing.T) {
	signingString := `(request-target): post /foo?param=value&pet=dog
host: example.com
date: Sun, 05 Jan 2014 21:31:40 GMT
content-type: application/json
digest: SHA-256=X48E9qOokqqrvdts8nOJRJN3OWDUoyWxBf7kbu9DBPE=
content-length: 18`

	publicKey, err := loadPublicKey([]byte(testSpecPublicKeyPEM))
	if err != nil {
		t.Fatal(err)
	}

	signature, err := base64.StdEncoding.DecodeString(`vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE=`)
	if err != nil {
		t.Fatal(err)
	}

	hashed := sha256.Sum256([]byte(signingString))
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRequestTargetKubernetes(t *testing.T) {

	os.Setenv("Http_Host", "verify-go.openfaas-fn.svc.cluster.local:8080")
	os.Setenv("Http_Method", "post")
	target := getRequestTarget()
	expected := "(request-target): post /function/verify-go"
	if target != expected {
		t.Errorf("want %s, got %s", expected, target)
	}
}

func TestGetRequestTargetSwarm(t *testing.T) {

	os.Setenv("Http_Host", "verify-go:8080")
	os.Setenv("Http_Method", "post")
	target := getRequestTarget()
	expected := "(request-target): post /function/verify-go"
	if target != expected {
		t.Errorf("want %s, got %s", expected, target)
	}
}

func getDigest(body string) string {
	d := sha256.Sum256([]byte(body))

	return base64.StdEncoding.EncodeToString(d[:])
}
