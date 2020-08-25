/*
Package authenticator implements authentication for the http requests based on HMAC.
Gin-Gonic has integrated Basic auth in case you need it.
Other authentication methods can be implemented with Authenticator interface.
*/
package authenticator

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
)

// These maps work as translators to get the corresponding functions from the values in the configuration.
var hashFuncs = map[string]func() hash.Hash{"sha1": sha1.New, "sha256": sha256.New}
var encryptFuncs = map[string]func([]byte) string{"hex": hex.EncodeToString,
	"base64.URL":    base64.URLEncoding.EncodeToString,
	"base64.RawURL": base64.RawURLEncoding.EncodeToString}

// Authenticator is the main interface of the package, it has only one method to implement.
type Authenticator interface {
	Authenticate(message []byte, signature string) error
}

// CreateAuthenticator is the function that initializes an authenticator of the appropriate kind based on the configuration received.
func CreateAuthenticator(class string, params map[string]string) (Authenticator, error) {
	var auth Authenticator
	var err error
	switch class {
	case "Signer":
		auth, err = NewSigner(params)
	default:
		auth, err = NewEmptyAuthenticator()
	}
	return auth, err
}

/*
EmptyAuthenticator implements the interface but does not have any logic, it accepts every input as valid.
*/
type EmptyAuthenticator struct {
}

// NewEmptyAuthenticator returns an EmptyAuthenticator struct.
func NewEmptyAuthenticator() (EmptyAuthenticator, error) {
	return EmptyAuthenticator{}, nil
}

// Authenticate always returns nil.
func (e EmptyAuthenticator) Authenticate(_ []byte, _ string) error {
	return nil
}

/*
Signer type stores key, hasher and encrypter to generate a signature based on the received message.
*/
type Signer struct {
	key       []byte
	hasher    func() hash.Hash
	encrypter func([]byte) string
}

// NewSigner creates a Signer struct with the received parameters.
func NewSigner(params map[string]string) (Signer, error) {
	var s Signer
	// Validate received parameters.
	key, ok := params["Key"]
	if !ok {
		return s, errors.New("Key not received for authenticator.")
	}
	hasherP, ok := params["Hasher"]
	if !ok {
		return s, errors.New("Hasher not received for authenticator.")
	}
	encrypterP, ok := params["Encrypter"]
	if !ok {
		return s, errors.New("Encrypter not received for authenticator.")
	}

	// Look for the functions in the functions map.
	hashF, ok := hashFuncs[hasherP]
	if !ok {
		return s, errors.New("Hashing function not found in hashFuncs.")
	}
	encryptF, ok := encryptFuncs[encrypterP]
	if !ok {
		return s, errors.New("Hashing function not found in hashFuncs.")
	}

	//fmt.Printf("Key: %q, Hasher: %q, Encrypter: %q \n", key, hasherP, encrypterP)
	return Signer{[]byte(key), hashF, encryptF}, nil
}

// Authenticate authenticates a message using the received signature and the parameters of the Signer.
func (s Signer) Authenticate(message []byte, signature string) error {
	newSignature := s.encrypter(getHMAC(message, s.key, s.hasher))
	if signature == newSignature {
		return nil
	}
	return fmt.Errorf("Signatures don't match. Received %q - Generated %q", signature, newSignature)
}

// Aux function to calculate the HMAC using the message, the hashing function and the key.
func getHMAC(message, key []byte, hasher func() hash.Hash) []byte {
	mac := hmac.New(hasher, key)
	mac.Write(message)

	return mac.Sum(nil)
}
