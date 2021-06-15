package bot

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
	httpUri    string
	blazeUri   string
)

func Request(ctx context.Context, method, path string, body []byte, accessToken string) ([]byte, error) {
	return RequestWithId(ctx, method, path, body, accessToken, UuidNewV4().String())
}
func RequestWithId(ctx context.Context, method, path string, body []byte, accessToken, requestID string) ([]byte, error) {
	req, err := http.NewRequest(method, httpUri+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Request-Id", requestID)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, ServerError(ctx, nil)
	}
	return ioutil.ReadAll(resp.Body)
}

func init() {
	httpClient = &http.Client{Timeout: 10 * time.Second}
	httpUri = "https://api.mixin.one"
	blazeUri = "blaze.mixin.one"
}

func SetBaseUri(base string) {
	httpUri = base
}

func SetBlazeUri(blaze string) {
	blazeUri = blaze
}
