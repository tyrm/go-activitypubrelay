package httpsign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
)

func Sign(pk *rsa.PrivateKey, kid string, req *http.Request) error {
	requestTarget := req.URL.Path
	host := req.Header.Get("Host")
	date := req.Header.Get("Date")
	digest := req.Header.Get("Digest")
	signedString := fmt.Sprintf("(request-target): post %s\nhost: %s\ndate: %s\ndigest: %s", requestTarget, host, date, digest)

	fmt.Println(signedString)

	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(signedString))
	if err != nil {
		return err
	}
	msgHashSum := msgHash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA256, msgHashSum)
	if err != nil {
		return err
	}

	headerTemplate := "keyId=\"%s\",headers=\"(request-target) host date digest\",signature=\"%s\""
	header := fmt.Sprintf(headerTemplate, kid, base64.StdEncoding.EncodeToString(signature))

	fmt.Println()
	fmt.Println(header)

	req.Header.Set("Signature", header)

	return nil
}
