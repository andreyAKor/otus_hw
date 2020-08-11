package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/streadway/amqp"
)

var (
	amqpDSN      = os.Getenv("TESTS_AMQP_DSN")
	queueName    = os.Getenv("TESTS_RMQ_QUEUE_NAME")
	bindingKey   = os.Getenv("TESTS_RMQ_BINDING_KEY")
	exchangeName = os.Getenv("TESTS_RMQ_EXCHANGE_NAME")
)

func init() {
	if amqpDSN == "" {
		amqpDSN = "amqp://guest:guest@localhost:5672/"
	}

	if queueName == "" {
		queueName = "senders"
	}
	if exchangeName == "" {
		exchangeName = "stats"
	}
}

type senderTest struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	messages      [][]byte
	messagesMutex *sync.RWMutex
	stopSignal    chan struct{}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (test *senderTest) startConsuming(*messages.Pickle) {
	test.messages = make([][]byte, 0)
	test.messagesMutex = new(sync.RWMutex)
	test.stopSignal = make(chan struct{})

	var err error

	test.conn, err = amqp.Dial(amqpDSN)
	panicOnErr(err)

	test.ch, err = test.conn.Channel()
	panicOnErr(err)

	// Consume
	_, err = test.ch.QueueDeclare(queueName, true, false, false, false, nil)
	panicOnErr(err)

	err = test.ch.QueueBind(queueName, bindingKey, exchangeName, false, nil)
	panicOnErr(err)

	events, err := test.ch.Consume(queueName, "", false, false, false, false, nil)
	panicOnErr(err)

	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case event := <-events:
				test.messagesMutex.Lock()
				test.messages = append(test.messages, event.Body)
				test.messagesMutex.Unlock()
			}
		}
	}(test.stopSignal)
}

func (test *senderTest) stopConsuming(*messages.Pickle, error) {
	test.stopSignal <- struct{}{}

	panicOnErr(test.ch.Close())
	panicOnErr(test.conn.Close())
	test.messages = nil
}

func (test *senderTest) iSendRequestToWithData(
	httpMethod, url, contentType string,
	data *messages.PickleStepArgument_PickleDocString,
) (err error) {
	var rsp *http.Response

	switch httpMethod {
	case http.MethodPost:
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

	return
}

func (test *senderTest) iReceiveEventWithData(data *messages.PickleStepArgument_PickleDocString) error {
	fmt.Println("wait 7 sec for producing sender stats...")
	time.Sleep(7 * time.Second)

	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	replacer := strings.NewReplacer("\n", "", "\t", "", ": ", ":")

	cleanData := replacer.Replace(data.Content)
	cleanData = replaceTimeNow(cleanData)

	for _, msg := range test.messages {
		cleanMsg := replacer.Replace(string(msg))
		cleanMsg = replaceTimeNow(cleanMsg)

		if cleanMsg == cleanData {
			return nil
		}
	}

	return fmt.Errorf(`event with data "%s" was not found in "%s"`, cleanData, test.messages)
}

func SenderFeatureContext(s *godog.Suite) {
	test := new(senderTest)

	s.BeforeScenario(test.startConsuming)

	s.Step(`^I send "([^"]*)" request to "([^"]*)" with "([^"]*)" create new event, with data:$`, test.iSendRequestToWithData)
	s.Step(`^I receive event with data:$`, test.iReceiveEventWithData)

	s.AfterScenario(test.stopConsuming)
}
