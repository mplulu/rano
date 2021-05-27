package rano

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mplulu/utils"
)

type TLGUpdateRequest struct {
	Offset         int64    `json:"offset"`
	AllowedUpdates []string `json:"allowed_updates"`
	Timeout        int      `json:"timeout"`
}

type TLGUpdateResponse struct {
	Ok     bool                       `json:"ok"`
	Result []*TLGUpdateResponseResult `json:"result"`
}

type TLGUpdateResponseResult struct {
	UpdateId      int64                           `json:"update_id"`
	Message       *TLGUpdateResponseResultMessage `json:"message"`
	EditedMessage *TLGUpdateResponseResultMessage `json:"edited_message"`
}

type TLGUpdateResponseResultMessage struct {
	MessageId   int64                               `json:"message_id"`
	UnixDate    int64                               `json:"date"`
	Text        string                              `json:"text"`
	From        *TLGUpdateResponseResultUser        `json:"from"`
	ChatChannel *TLGUpdateResponseResultChatChannel `json:"chat"`
	ReplyTo     *TLGUpdateResponseResultReplyTo     `json:"reply_to_message"`
	Entities    []*TLGUpdateResponseResultEntity    `json:"entities"`
}

type TLGUpdateResponseResultChatChannel struct {
	Id   int64  `json:"id"`
	Name string `json:"title"`
}

type TLGUpdateResponseResultUser struct {
	Id   int64  `json:"id"`
	Name string `json:"first_name"`
}

type TLGUpdateResponseResultReplyTo struct {
	UnixDate    int64                               `json:"date"`
	Text        string                              `json:"text"`
	ChatChannel *TLGUpdateResponseResultChatChannel `json:"chat"`
	From        *TLGUpdateResponseResultUser        `json:"from"`
}

type TLGUpdateResponseResultEntity struct {
	Type string                       `json:"type"`
	User *TLGUpdateResponseResultUser `json:"user"`
}

func (rano *Rano) getUpdates() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("", time.Now(), r, GetStack())
			// delay a bit to prevent internet or something offline, and this run repeately very fast
			utils.DelayInDuration(10 * time.Second)
			go rano.getUpdates()
		}
	}()
	urlStr := fmt.Sprintf("%sgetUpdates", rano.baseUrl)

	requestParams := &TLGUpdateRequest{
		Offset:  rano.lastUpdateId + 1,
		Timeout: 5 * 60,
	}
	requestParamsRaw, err := json.Marshal(requestParams)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(requestParamsRaw))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := rano.client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	response.Body.Close()
	fmt.Println("body", string(body))
	var responseData *TLGUpdateResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		panic(err)
	}

	messageList := responseData.toMessageList()
	for _, message := range messageList {
		rano.lastUpdateId = message.UpdateId
		rano.MessageChan <- message
	}
	go rano.getUpdates()
}

func (r *TLGUpdateResponse) toMessageList() []*Message {
	list := []*Message{}
	if r.Ok {
		for _, result := range r.Result {
			message := result.Message
			if result.EditedMessage != nil {
				message = result.EditedMessage
			}
			author := &User{
				Id:   message.From.Id,
				Name: message.From.Name,
			}
			group := &Group{
				Id:   message.ChatChannel.Id,
				Name: message.ChatChannel.Name,
			}
			var replyTo *Message
			if message.ReplyTo != nil {
				replyToAuthor := &User{
					Id:   message.ReplyTo.From.Id,
					Name: message.ReplyTo.From.Name,
				}
				replyToGroup := &Group{
					Id:   message.ReplyTo.ChatChannel.Id,
					Name: message.ReplyTo.ChatChannel.Name,
				}
				replyTo = &Message{
					From:  replyToAuthor,
					Group: replyToGroup,
					Text:  message.ReplyTo.Text,
					Date:  time.Unix(message.ReplyTo.UnixDate, 0),
				}
			}

			entities := []*Entity{}
			for _, entityResp := range message.Entities {
				entityObjc := &Entity{
					Type: entityResp.Type,
				}
				if entityResp.User != nil {
					entityObjc.User = &User{
						Id:   entityResp.User.Id,
						Name: entityResp.User.Name,
					}
				}
			}
			messageObjc := &Message{
				UpdateId: result.UpdateId,
				From:     author,
				Group:    group,
				Text:     message.Text,
				Date:     time.Unix(message.UnixDate, 0),
				ReplyTo:  replyTo,
				Entities: entities,
			}
			list = append(list, messageObjc)
		}
	}
	return list
}

type User struct {
	Id   int64
	Name string
}

type Group struct {
	Id   int64
	Name string
}

type Message struct {
	UpdateId int64
	From     *User
	Group    *Group
	Text     string
	Date     time.Time
	ReplyTo  *Message
	Entities []*Entity
}

type Entity struct {
	Type string
	User *User
}

// {
//   "ok": true,
//   "result": [
//     {
//       "update_id": 473131607,
//       "message": {
//         "message_id": 21246,
//         "from": {
//           "id": 684139036,
//           "is_bot": false,
//           "first_name": "PhDung"
//         },
//         "chat": {
//           "id": -311058787,
//           "title": "CS & Army Training",
//           "type": "group",
//           "all_members_are_administrators": true
//         },
//         "date": 1557744238,
//         "reply_to_message": {
//           "message_id": 21233,
//           "from": {
//             "id": 677314948,
//             "is_bot": true,
//             "first_name": "MissTeen",
//             "username": "p2playsvno_bot"
//           },
//           "chat": {
//             "id": -311058787,
//             "title": "CS & Army Training",
//             "type": "group",
//             "all_members_are_administrators": true
//           },
//           "date": 1557736044,
//           "text": "DWTrans Mismatch FROM btu3 agent TO btu30284 member | withdraw | 16740 | 13/05/2019 15:27:19 | game backend credit mismatch"
//         },
//         "text": "Cái này báo lỗi tại vì deposit vs withdraw nhiều lần . Tiền vẫn hiện đúng"
//       }
//     }
//   ]
// }
