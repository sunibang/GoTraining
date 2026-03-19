package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/romangurevitch/go-training/pkg/api/apierror"
)

// DoRequest executes the request with the given context, checks the response status,
// and returns the body bytes. Returns *apierror.APIError on non-expected status codes.
func DoRequest(ctx context.Context, client *http.Client, r *http.Request, expectedResponses ...int) ([]byte, error) {
	resp, err := client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err = checkResponse(resp, expectedResponses...); err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func checkResponse(resp *http.Response, expectedResponses ...int) error {
	for _, expected := range expectedResponses {
		if expected == resp.StatusCode {
			return nil
		}
	}
	return getAPIError(resp)
}

func getAPIError(resp *http.Response) error {
	if resp.StatusCode == http.StatusUnauthorized {
		return apierror.APIError{Message: resp.Status}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiErr *apierror.APIError
	if err = json.Unmarshal(body, &apiErr); err != nil {
		return apierror.ErrInternalServerError
	}
	return apiErr
}

// GetURL builds a full URL by joining baseURL with path p and appending a formatted query string.
func GetURL(baseURL, p, queryFormat string, args ...any) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, p)
	if queryFormat != "" {
		u.RawQuery = fmt.Sprintf(queryFormat, args...)
	}
	return u.String(), nil
}

// HeaderApplicationJSON returns the Content-Type header key and value for JSON.
func HeaderApplicationJSON() (key, value string) {
	return "Content-Type", "application/json"
}
