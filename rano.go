package rano

import (
	"fmt"
	"net/http"
)

type Rano struct {
	isDisable bool

	token      string
	chatIdList []string
	baseUrl    string
	client     *http.Client
}

func NewRano(token string, chatIdList []string) *Rano {
	if token == "" {
		fmt.Println("Rano is disable. No token")
		return &Rano{
			isDisable: true,
		}
	}
	return &Rano{
		isDisable:  false,
		token:      token,
		chatIdList: chatIdList,
		baseUrl:    fmt.Sprintf("https://api.telegram.org/bot%s/", token),
		client:     &http.Client{},
	}
}

func (rano *Rano) Send(text string) {
	if rano.isDisable {
		return
	}
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
