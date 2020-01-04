package rano

import (
	"fmt"
	"net/http"
	"time"
)

type Rano struct {
	isDisable bool

	token      string
	chatIdList []string
	baseUrl    string
	client     *http.Client

	lastUpdateId            int64
	alreadyReceivingMessage bool
	MessageChan             chan *Message
}

type MessageHandler func(message string)

func NewRano(token string, chatIdList []string) *Rano {
	if token == "" {
		fmt.Println("Rano is disable. No token")
		return &Rano{
			isDisable: true,
		}
	}
	rano := &Rano{
		isDisable:  false,
		token:      token,
		chatIdList: chatIdList,
		baseUrl:    fmt.Sprintf("https://api.telegram.org/bot%s/", token),
		client: &http.Client{
			Timeout: 5 * 60 * time.Second,
		},
	}
	return rano
}

func (rano *Rano) SendTo(chatId int64, text string) error {
	if rano.isDisable {
		return nil
	}
	chatIdStr := fmt.Sprintf("%d", chatId)
	_, err := rano.sendRequest(
		"sendMessage",
		map[string]string{
			"chat_id": chatIdStr,
			"text":    text,
		})
	// fmt.Println("RanoSend: ", chatId, text, response, err)
	if err != nil {
		return err
	}
	return nil
}

func (rano *Rano) SendPhoto(chatId int64, photo []byte) error {
	if rano.isDisable {
		return nil
	}
	chatIdStr := fmt.Sprintf("%d", chatId)
	_, err := rano.sendRequestWithBinaryFile(
		"sendPhoto",
		map[string]string{
			"chat_id": chatIdStr,
		},
		map[string][]byte{
			"photo": photo,
		})
	// fmt.Println("RanoSend: ", chatId, text, response, err)
	if err != nil {
		return err
	}
	return nil
}

func (rano *Rano) Send(text string) error {
	if rano.isDisable {
		return nil
	}
	for _, chatId := range rano.chatIdList {
		_, err := rano.sendRequest(
			"sendMessage",
			map[string]string{
				"chat_id": chatId,
				"text":    text,
			})
		// fmt.Println("RanoSend: ", chatId, text, response, err)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rano *Rano) StartReceivingMessage() {
	if rano.alreadyReceivingMessage {
		return
	}
	rano.alreadyReceivingMessage = true
	rano.MessageChan = make(chan *Message, 1)
	go rano.getUpdates()
}
