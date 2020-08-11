package main

import (
	"net/url"
	"strings"
	"time"

	"github.com/cucumber/messages-go/v10"
)

func preparePostData(data *messages.PickleStepArgument_PickleDocString) string {
	replacer := strings.NewReplacer("\t", "")
	cleanData := replacer.Replace(data.Content)

	params := url.Values{}

	rows := strings.Split(cleanData, "\n")
	for _, row := range rows {
		rowData := strings.Split(row, "=")

		value := replaceTimeNow(rowData[1])

		params.Add(rowData[0], value)
	}

	return params.Encode()
}

func replaceTimeNow(value string) string {
	replacer := strings.NewReplacer("@timeNow", time.Now().Format("2006-01-02T15:04:00-07:00"))
	return replacer.Replace(value)
}
