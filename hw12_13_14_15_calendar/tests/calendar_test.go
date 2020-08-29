package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

type calendarTest struct {
	responseStatusCode int
	responseTimeout    time.Duration
	responseBody       []byte
}

func (test *calendarTest) iSendRequestTo(httpMethod, url string) (err error) {
	var rsp *http.Response

	start := time.Now()

	switch httpMethod {
	case http.MethodGet, http.MethodDelete:
		var req *http.Request

		client := &http.Client{}

		req, err = http.NewRequest(httpMethod, url, nil)
		if err != nil {
			return
		}

		rsp, err = client.Do(req)
		defer rsp.Body.Close()
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}

	test.responseTimeout = time.Now().Sub(start)
	test.responseStatusCode = rsp.StatusCode
	test.responseBody, err = ioutil.ReadAll(rsp.Body)

	return
}

func (test *calendarTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}

	return nil
}

func (test *calendarTest) iSendRequestToWithData(
	httpMethod, url, contentType string,
	data *messages.PickleStepArgument_PickleDocString,
) (err error) {
	var rsp *http.Response

	switch httpMethod {
	case http.MethodPost, http.MethodPut:
		var req *http.Request

		client := &http.Client{}

		bodyData := preparePostData(data)
		req, err = http.NewRequest(httpMethod, url, strings.NewReader(bodyData))
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", contentType)

		rsp, err = client.Do(req)
		defer rsp.Body.Close()
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}

	test.responseStatusCode = rsp.StatusCode
	test.responseBody, err = ioutil.ReadAll(rsp.Body)

	return
}

func (test *calendarTest) iReceiveResponseWithData(data *messages.PickleStepArgument_PickleDocString) error {
	replacer := strings.NewReplacer("\n", "", "\t", "", ": ", ":")

	cleanData := replacer.Replace(data.Content)
	cleanResponseBody := replacer.Replace(string(test.responseBody))

	if cleanResponseBody != cleanData {
		return fmt.Errorf(`response with data "%s" was not found in "%s"`, cleanData, test.responseBody)
	}

	return nil
}

func CalendarFeatureContext(s *godog.Suite) {
	test := new(calendarTest)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)

	s.Step(`^I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^I receive response with data:$`, test.iReceiveResponseWithData)
}
