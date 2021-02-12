package test_case

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
)

type TypeOfRequest string

const (
	POST TypeOfRequest = "post"
	GET  TypeOfRequest = "get"
)

//TestCase denotes a single test case and contains all information required
// to run this test
type TestCase struct {
	// URL of the API
	URL string `json:"url"`

	// TypeOfRequest what type of request this is
	// Possible values are GET, POST, get, post
	TypeOfRequest TypeOfRequest `json:"type_of_request"`

	// ExpectedStatusCode 200, 400, 404, 401, 403
	ExpectedStatusCode int16 `json:"expected_status_code"`

	//RequestParams contains query params for this request
	// only used if this is a GET request, ignored otherwise
	RequestParams map[string]string `json:"request_params"`

	//RequestBody contains post body if this is a POST request
	// ignored otherwise
	RequestBody map[string]interface{} `json:"request_body"`

	// Headers to pass along with the request
	Headers map[string]string `json:"headers"`

	//ExpectedResponse contains expected response keys with their values
	// as the type of data they will contain in the actual response
	ExpectedResponse map[string]interface{} `json:"expected_response"`

	// Debug enable debug mode
	Debug bool `json:"debug"`

	httpRunner *resty.Client
}

func (testCase TestCase) Run(httpRunner *resty.Client) error {
	if httpRunner != nil {
		testCase.httpRunner = httpRunner
	} else {
		testCase.httpRunner = resty.New()
	}

	switch testCase.TypeOfRequest {
	case POST:
		{
			return testCase.runPostRequest()
		}
	case GET:
		{
			return testCase.runGetRequest()
		}
	default:
		return errors.New("unsupported type of request")
	}
}

// runGetRequest runs the specified get request
// if the request is not a get request, it aborts and returns an error
func (testCase TestCase) runGetRequest() error {
	if testCase.TypeOfRequest != GET {
		return errors.New("not a GET request")
	}

	response, err := testCase.httpRunner.
		R().
		SetQueryParams(testCase.RequestParams).
		SetHeaders(testCase.Headers).
		Get(testCase.URL)

	return validateTestCase(testCase, response, err)
}

func (testCase TestCase) runPostRequest() error {
	if testCase.TypeOfRequest != POST {
		return errors.New("not a POST request")
	}

	response, err := testCase.httpRunner.
		R().
		SetBody(testCase.RequestBody).
		SetHeaders(testCase.Headers).
		Post(testCase.URL)

	return validateTestCase(testCase, response, err)
}

func validateTestCase(testCase TestCase, response *resty.Response, err error) error {
	if testCase.Debug {
		log.Printf("\n\n url = %+v, response = %+v, err = %+v \n\n", testCase.URL, response, err)
	}

	if testCase.ExpectedStatusCode < 400 && (err != nil || response.IsError()) {
		return err
	}

	parsedResponse := gjson.Parse(response.String())

	if testCase.Debug {
		log.Printf("\n\n parsedResponse = %+v\n\n", parsedResponse)
	}

	return matchResponseWithExpectedTypes(testCase.ExpectedResponse, parsedResponse)
}

func matchResponseWithExpectedTypes(expectedResponse map[string]interface{}, parsedResponse gjson.Result) error {
	for key, value := range expectedResponse {
		data := parsedResponse.Get(key)
		switch value.(type) {
		case string:
			{

				if data.Type.String() != value.(string) {
					return errors.New(fmt.Sprintf("\n\ntype of %+v (%v) does not match. here's what was returned = %v\n\n", key, value.(string), data))
				}
			}

		case map[string]interface{}:
			{
				return matchResponseWithExpectedTypes(value.(map[string]interface{}), data)
			}
		}
	}
	return nil
}
