/*
Example test package for data handler.
*/
package webserver_test

import (
	"bytes"
	"fmt"
	"github.com/efark/data-receiver/authenticator"
	"github.com/efark/data-receiver/extractor"
	"github.com/efark/data-receiver/webserver"
	"github.com/efark/data-receiver/writer"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	net_url "net/url"
	"testing"
)

var (
	mw *writer.MemoryWriter
)

func TestDataHandler(t *testing.T) {
	var err error
	mw, err = writer.NewMemoryWriter()
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	destroy := setupTest(t, mw)
	defer destroy()

	method := http.MethodPost
	url := "localhost:8080"
	body := []byte(`test message`)
	urlParams := []gin.Param{{"service", "test"}}
	queryParams := net_url.Values{}
	headers := map[string]string{"x-user-id": "test_id", "x-signature": "GXjQXzGexUuSH444qEyMI-b9Lif_Uq39gElhs_7PMVY="}

	c, record := createGinContext(method, url, body, urlParams, queryParams, headers)

	webserver.DataHandler(c)

	if record.Result().StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("Status code: %v\n", record.Result().StatusCode))
		t.FailNow()
	}

	// Test output
	messages := mw.GetMessages()
	if len(messages) != 1 {
		t.Error("Incorrect number of messages.")
		t.FailNow()
	}

	message := messages[0]
	if message != "test message" {
		t.Error("message content error.")
		t.FailNow()
	}
}

//key []byte, hasher func() hash.Hash, encrypter func([]byte) string
func setupTest(t *testing.T, w writer.Writer) func() {
	t.Log("Setting up test service.")

	extConfig := map[string]string{"signature": "x-signature"}
	ext, err := extractor.NewHeaderExtractor(extConfig)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	authConfig := map[string]string{"Key": "magicKey", "Hasher": "sha256", "Encrypter": "base64.URL"}
	auth, err := authenticator.NewSigner(authConfig)
	w = mw

	webserver.SetService("test", ext, auth, w)
	return func() {
		t.Log("Closing service.")
		webserver.CloseWriters()
	}
}

func createGinContext(method string, url string, body []byte, urlParams []gin.Param, queryParams net_url.Values, headers map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	// Parse URL
	u, err := net_url.Parse(url)
	if err != nil {
		panic(err)
	}
	// Params
	q, _ := net_url.ParseQuery(u.RawQuery)
	for k, v := range queryParams {
		q.Add(k, v[0])
	}
	u.RawQuery = q.Encode()
	url = u.String()

	// New request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		panic("Cannot create an instance of request object")
	}

	// Add headers to request
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resw := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resw)
	context.Request = req

	context.Params = urlParams

	return context, resw
}
