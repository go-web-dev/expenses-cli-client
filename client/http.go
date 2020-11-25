package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// HTTPClient represents the HTTP client which communicates with reminders backend API
type HTTPClient struct {
	client     *http.Client
	BackendURI string
}

// NewHTTPClient creates a new instance of HTTPClient
func NewHTTPClient(uri string) HTTPClient {
	return HTTPClient{
		BackendURI: uri,
		client:     &http.Client{},
	}
}

type expenseRequestBody struct {
	Title    string  `json:"title"`
	Currency string  `json:"currency"`
	Price    float64 `json:"price"`
}

type authReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create calls the create API endpoint
func (c HTTPClient) Create(title, currency string, price float64) error {
	body := expenseRequestBody{
		Title:    title,
		Currency: currency,
		Price:    price,
	}
	req, err := c.newReqWithToken(http.MethodPost, "/expenses", body)
	if err != nil {
		return err
	}
	_, err = c.apiCall(req, http.StatusCreated)
	return err
}

// Update calls the update API endpoint
func (c HTTPClient) Update(id, title, currency string, price float64) error {
	body := expenseRequestBody{
		Title:    title,
		Currency: currency,
		Price:    price,
	}
	req, err := c.newReqWithToken(http.MethodPatch, "/expenses/"+id, body)
	if err != nil {
		return err
	}
	_, err = c.apiCall(req, http.StatusNoContent)
	return err
}

// Delete calls the delete API endpoint
func (c HTTPClient) Delete(id string) error {
	req, err := c.newReqWithToken(http.MethodDelete, "/expenses/"+id, nil)
	if err != nil {
		return err
	}
	_, err = c.apiCall(req, http.StatusNoContent)
	return err
}

// GetAll calls the get-all API endpoint
func (c HTTPClient) GetAll(page, pageSize string) ([]byte, error) {
	params := url.Values{}
	params.Add("page", page)
	params.Add("page_size", pageSize)
	req, err := c.newReqWithToken(http.MethodPatch, "/expenses?"+params.Encode(), nil)
	if err != nil {
		return []byte{}, err
	}
	return c.apiCall(req, http.StatusNoContent)
}

// GetByIDs calls the get-by-ids API endpoint
func (c HTTPClient) GetByIDs(ids ...string) ([]byte, error) {
	req, err := c.newReqWithToken(http.MethodPatch, "/expenses/"+strings.Join(ids, ","), nil)
	if err != nil {
		return []byte{}, err
	}
	return c.apiCall(req, http.StatusNoContent)
}

// Login calls the login API endpoint
func (c HTTPClient) Login(email, password string) ([]byte, error) {
	body := authReqBody{
		Email:    email,
		Password: password,
	}
	req, err := c.newReq(http.MethodPost, "/login", body)
	if err != nil {
		return []byte{}, err
	}
	return c.apiCall(req, http.StatusOK)
}

// Signup calls the signup API endpoint
func (c HTTPClient) Signup(email, password string) ([]byte, error) {
	body := authReqBody{
		Email:    email,
		Password: password,
	}
	req, err := c.newReq(http.MethodPost, "/signup", body)
	if err != nil {
		return []byte{}, err
	}
	return c.apiCall(req, http.StatusOK)
}

// Logout calls the logout API endpoint
func (c HTTPClient) Logout() error {
	req, err := c.newReq(http.MethodPost, "/logout", nil)
	if err != nil {
		return err
	}
	_, err = c.apiCall(req, http.StatusOK)
	return err
}

// apiCall makes a new backend api call
func (c HTTPClient) apiCall(req *http.Request, resCode int) ([]byte, error) {
	res, err := c.client.Do(req)
	if err != nil {
		e := errors.Wrap(err, "could not make http call")
		return []byte{}, e
	}

	resBody, err := c.readResBody(res.Body)
	if err != nil {
		return []byte{}, err
	}

	if res.StatusCode != resCode {
		if len(resBody) > 0 {
			fmt.Printf("got this response body:\n%s\n", resBody)
		}
		return []byte{}, fmt.Errorf(
			"expected response code: %d, got: %d",
			resCode,
			res.StatusCode,
		)
	}

	return []byte(resBody), err
}

// readBody reads response body
func (c HTTPClient) readResBody(b io.Reader) (string, error) {
	bs, err := ioutil.ReadAll(b)
	if err != nil || len(bs) == 0 {
		return "", errors.Wrap(err, "could not read response body")
	}

	var buff bytes.Buffer
	if err := json.Indent(&buff, bs, "", "\t"); err != nil {
		return "", errors.Wrap(err, "could not indent json")
	}

	return buff.String(), nil
}

func (c HTTPClient) newReq(method, path string, body interface{}) (*http.Request, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		e := errors.Wrap(err, "could not marshal request body")
		return nil, e
	}
	req, err := http.NewRequest(method, c.BackendURI+path, bytes.NewReader(bs))
	if err != nil {
		return nil, errors.Wrap(err, "could not create http request")
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c HTTPClient) newReqWithToken(method, path string, body interface{}) (*http.Request, error) {
	req, err := c.newReq(method, path, body)
	if err != nil {
		return nil, err
	}
	credentials, err := readCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "could not read credentials")
	}
	req.Header.Add("Bearer", credentials.AccessToken)
	return req, nil
}
