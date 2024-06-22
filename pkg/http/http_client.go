package external

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpClientContract interface {
	GetRequest(reqUrl string, queryParams map[string]string, headers map[string]string) (string, error)
	PostRequest(reqUrl string, jsonBody string, headers map[string]string) (string, error)
}

type HttpClient struct {
}

func NewHttpClient(isTesting bool) HttpClientContract {
	if isTesting { // Todo: Get from somewhere else
		return &MockHttpClient{}
	} else {
		return &HttpClient{}
	}
}

func (c *HttpClient) GetRequest(reqUrl string, queryParams map[string]string, headers map[string]string) (string, error) {
	u, err := url.Parse(reqUrl)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%d", resp.StatusCode)
	}

	return string(body), nil
}

func (c *HttpClient) PostRequest(reqUrl string, jsonBody string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return "", err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("ERROR %d %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

var mockSequences map[string]mockSequence = make(map[string]mockSequence)

type MockHttpClient struct {
}

func NewMockHttpClient() *MockHttpClient {
	return &MockHttpClient{}
}

type mockSequence struct {
	Response string
	Error    error
}

func FakeHttp(reqUrl string, response string, respErr error) {
	mockSequences[reqUrl] = mockSequence{Response: response, Error: respErr}
}

func (c *MockHttpClient) getResponse(reqUrl string) (string, error) {
	seq, ok := mockSequences[reqUrl]

	if ok {
		return seq.Response, seq.Error
	}

	return "", fmt.Errorf("%s can not be mocked, expectations are not specified.", reqUrl)
}

func (c *MockHttpClient) GetRequest(reqUrl string, queryParams map[string]string, headers map[string]string) (string, error) {
	return c.getResponse(reqUrl)
}

func (c *MockHttpClient) PostRequest(reqUrl string, jsonBody string, headers map[string]string) (string, error) {
	return c.getResponse(reqUrl)
}
