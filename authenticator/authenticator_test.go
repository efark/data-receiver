package authenticator_test

import (
	"github.com/efark/data-receiver/authenticator"
	"testing"
)

func TestSigner_Authenticate(t *testing.T) {
	params := map[string]string{"Key": `magickey`, "Hasher": "sha1", "Encrypter": "base64.RawURL"}
	s, err := authenticator.NewSigner(params)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	m := []byte(`Example message`)
	signature := `eZIp7BDQLn3PuZrDPWSlW3x6dgo`

	err = s.Authenticate(m, signature)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

}

func TestSigner_Invalid(t *testing.T) {
	params := map[string]string{"Key": `magickey`, "Hasher": "sha256", "Encrypter": "base64.URL"}
	s, err := authenticator.NewSigner(params)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	m := []byte(`Example message`)
	signature := `Wrong signature`

	err = s.Authenticate(m, signature)
	if err == nil {
		t.Error(err.Error())
		t.FailNow()
	}

}

func TestEmptyAuthenticator_Authenticate(t *testing.T) {
	s, err := authenticator.NewEmptyAuthenticator()
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	m := []byte(`Example message`)
	signature := `Wrong signature`

	err = s.Authenticate(m, signature)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}
