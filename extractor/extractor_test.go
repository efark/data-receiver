package extractor_test

import (
	"bytes"
	"fmt"
	"github.com/efark/data-receiver/extractor"
	"github.com/gin-gonic/gin"
	"net/http"
	net_url "net/url"
	"testing"
)

func setupTest(method string, url string, body []byte, params net_url.Values, headers map[string]string) *gin.Context {
	// Parse URL
	u, err := net_url.Parse(url)
	if err != nil {
		panic(err)
	}
	// Params
	q, _ := net_url.ParseQuery(u.RawQuery)
	for k, v := range params {
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

	context, _ := gin.CreateTestContext(nil)
	context.Request = req

	return context
}

func TestHeaderExtractor_Extract(t *testing.T) {
	url := "localhost:8080/"
	body := []byte(`This is the body of the message`)
	headers := map[string]string{"x-signature": "message signature", "x-user-id": "testUserId"}

	// method string, url string, body []byte, params net_url.Values, headers map[string]string
	context := setupTest(http.MethodPost, url, body, nil, headers)

	extractorConfig := map[string]string{"signature": "x-signature", "user_id": "x-user-id"}
	ext, err := extractor.NewHeaderExtractor(extractorConfig)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	mappedValues := ext.Extract(context)

	for key, value := range extractorConfig {

		x, ok := mappedValues[key]
		if !ok {
			t.Error(fmt.Sprintf("Header not found: '%s',\n", key))
			t.FailNow()
		}
		if x != headers[value] {
			t.Error(fmt.Sprintf("Header value is different: expected '%s' - received '%s'.\n", headers[value], x))
			t.FailNow()
		}
		fmt.Printf("Header value: expected '%s' - received '%s'.\n", headers[value], x)
	}
}

func TestQueryExtractor_Extract(t *testing.T) {
	url := "localhost:8080/"
	body := []byte(`This is the body of the message`)

	params := net_url.Values{}
	params.Add("c", "testUserId")
	params.Add("s", "message signature")

	// method string, url string, body []byte, params net_url.Values, headers map[string]string
	context := setupTest(http.MethodPost, url, body, params, nil)

	extractorConfig := map[string]string{"user_id": "c", "signature": "s"}
	ext, err := extractor.NewQueryExtractor(extractorConfig)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	mappedValues := ext.Extract(context)

	for key, value := range extractorConfig {

		x, ok := mappedValues[key]
		if !ok {
			t.Error(fmt.Sprintf("Query param not found: '%s',\n", key))
			t.FailNow()
		}
		if x != params.Get(value) {
			t.Error(fmt.Sprintf("Query param value is different: expected '%s' - received '%s'.\n", params.Get(value), x))
			t.FailNow()
		}
		fmt.Printf("Query param value: expected '%s' - received '%s'.\n", params.Get(value), x)
	}
}
