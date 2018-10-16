package rano

import (
	"fmt"
	"net/http"
)

type Rano struct {
	token      string
	chatIdList []string
	baseUrl    string
	client     *http.Client
}

func NewRano(token string, chatIdList []string) *Rano {
	return &Rano{
		token:      token,
		chatIdList: chatIdList,
		baseUrl:    fmt.Sprintf("https://api.telegram.org/bot%s/", token),
		client:     &http.Client{},
	}
}

func (rano *Rano) Send(text string) {
	for _, chatId := range rano.chatIdList {
		response, err := rano.sendRequest(
			"sendMessage",
			map[string]string{
				"chat_id": chatId,
				"text":    text,
			})
		fmt.Println("RanoSend: ", response, err)
	}
}
