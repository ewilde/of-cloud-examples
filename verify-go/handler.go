package function

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type KeyType struct {
	Id  string `json:"id"`
	PEM string `json:"pem"`
}

// Handle a serverless request
func Handle(req []byte) string {

	log.SetOutput(os.Stderr)

	resp, err := http.Get(fmt.Sprintf("http://%s:8080/certificates/callback", gateway()))
	if err != nil {
		return errorResult(err)
	}
	defer resp.Body.Close()
	pem, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errorResult(err)
	}

	keyData := &KeyType{}
	err = json.Unmarshal(pem, keyData)
	if err != nil {
		return errorResult(err)
	}

	err = verify(req, keyData.PEM)
	if err != nil {
		return errorResult(err)
	}

	log.Println("Verified OK")

	return "valid"
}

func verify(req []byte, pem string) error {
	authorizationHeader := header("Authorization")
	authParts := strings.Split(authorizationHeader, ",")[3]
	signatureEnc := strings.Split(authParts, "=\"")[1] + "="
	signature, err := base64.StdEncoding.DecodeString(signatureEnc) // signing

	signingString := getSigningString(req)
	log.Printf("Signature:\n%s\n", signatureEnc)
	log.Printf("Signing string:\n%s", signingString)
	log.Printf("Public key:\n%s\n", pem)
	publicKey, err := loadPublicKey([]byte(pem))
	if err != nil {
		return err
	}

	hashed := sha256.Sum256([]byte(signingString))
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return err
	}

	return nil
}

func getSigningString(body []byte) string {

	d := sha256.Sum256(body)
	digest := base64.StdEncoding.EncodeToString(d[:])

	res := strings.Builder{}
	res.WriteString(fmt.Sprintf("%s\n", getRequestTarget()))
	res.WriteString(fmt.Sprintf("host: %s\n", header("X-Forwarded-Host")))
	res.WriteString(fmt.Sprintf("date: %s\n", header("Date")))
	res.WriteString(fmt.Sprintf("content-type: %s\n", header("Content-Type")))
	res.WriteString(fmt.Sprintf("digest: SHA-256=%s\n", string(digest)))
	res.WriteString(fmt.Sprintf("content-length: %s", strconv.Itoa(len(body))))

	return res.String()
}

func getRequestTarget() string {
	path := header("Path")
	if path == "/" {
		path = ""
	}

	name := strings.Split(strings.Split(header("Host"), ":")[0], ".")[0]

	requestTarget := fmt.Sprintf("(request-target): %s /function/%s%s", strings.ToLower(header("Method")), name, path)
	query := header("Query")
	if query != "" {
		requestTarget = fmt.Sprintf("%s?%s", requestTarget, query)
	}

	return requestTarget
}

func errorResult(err error) string {
	log.Printf("error: %v", err)

	return err.Error()
}

func loadPublicKey(keyData []byte) (*rsa.PublicKey, error) {
	pem, _ := pem.Decode(keyData)
	if pem.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("public key is of the wrong type: %s", pem.Type)
	}

	key, err := x509.ParsePKIXPublicKey(pem.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func header(key string) string {
	key = fmt.Sprintf("Http_%s", strings.Replace(key, "-", "_", -1))
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	response := strings.Builder{}
	for _, e := range os.Environ() {
		response.WriteString(e)
		response.WriteString("\n")
	}

	return ""
}

func gateway() string {
	host, ok := os.LookupEnv("gateway_host")
	if !ok {
		return "gateway"
	}

	return host
}
