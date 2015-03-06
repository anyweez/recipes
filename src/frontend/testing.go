package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"lib/config"
	log "logging"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
)

const (
	ConfigPath = "/home/luke/git/recipes/recipes.test.conf"
)

type ErrorResponse struct {
	ErrorText string
	HttpCode  int
	ErrorCode int
}

type TestClient struct {
	Client   http.Client
	Requests []*http.Request
}

func NewTestClient() (TestClient, error) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		return TestClient{}, err
	}

	return TestClient{
		Client: http.Client{
			Jar: jar,
		},
	}, nil
}

/**
 * Queue up the http.Request to be run when Run() is executed.
 */
func (t *TestClient) QueueRequest(req *http.Request, err error) error {
	t.Requests = append(t.Requests, req)

	return err
}

/**
 * Run all of the requests that have been queued.
 */
func (t *TestClient) Run() (ErrorResponse, int, error) {
	er := ErrorResponse{}
	status := 0
	var err error

	for _, req := range t.Requests {
		er, status, err = t.DoTest(req)
	}

	t.Requests = make([]*http.Request, 0)

	return er, status, err
}

func (t *TestClient) DoTest(req *http.Request) (ErrorResponse, int, error) {
	resp, err := t.Client.Do(req)
	defer resp.Body.Close()

	obj := ErrorResponse{}

	if err != nil {
		return ErrorResponse{}, 0, err
	}

	body, rerr := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &obj)

	if rerr != nil {
		return ErrorResponse{}, 0, rerr
	}

	return obj, resp.StatusCode, nil
}

/**
 * Run a test by sending an HTTP request to `path` using the specified `method.` This function
 * also checks that the response code that comes back is equivalent to `expected` and will return
 * an error if not.
 */
func NewHttpTest(name string, method string, path string, body []byte, expected int) (*FrontendServer, []string, error) {
	// Initialize the test configuration.
	le := log.New(name, nil)
	conf, cerr := config.New(ConfigPath)

	if cerr != nil {
		return &FrontendServer{}, nil, cerr
	}

	fes, ferr := NewFrontendServer(conf, le)

	if ferr != nil {
		return &FrontendServer{}, nil, cerr
	}

	// Once we've got the frontend server initialized, the work is the same as if the FES had already
	// existed.
	return NextHttpTest(name, method, path, body, expected, &fes, nil)
}

func NextHttpTest(name string, method string, path string, body []byte, expected int, fes *FrontendServer, cookies []string) (*FrontendServer, []string, error) {
	// Make the request and record the response.
	req, err := http.NewRequest(method, path, bytes.NewReader(body))

	if err != nil {
		return fes, nil, err
	}

	if len(cookies) > 0 {
		req.Header.Set("Set-Cookie", cookies[0])
	}

	fmt.Println("cookies!")
	fmt.Println(req.Header.Get("Set-Cookie"))

	w := httptest.NewRecorder()
	fes.Router.ServeHTTP(w, req)

	// Check whether the response code matches the expected value. If not, we
	// should return an error.
	if w.Code != expected {
		return fes, w.Header()["Set-Cookie"], errors.New(fmt.Sprintf("Returned HTTP code %d instead of expected %d", w.Code, expected))
	}

	return fes, w.Header()["Set-Cookie"], nil
}
