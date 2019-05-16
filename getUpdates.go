package rano

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
	UpdateId int64                           `json:"update_id"`
	Message  *TLGUpdateResponseResultMessage `json:"message"`
}

type TLGUpdateResponseResultMessage struct {
	MessageId   int64                               `json:"message_id"`
	UnixDate    int64                               `json:"date"`
	Text        string                              `json:"text"`
	From        *TLGUpdateResponseResultFrom        `json:"from"`
	ChatChannel *TLGUpdateResponseResultChatChannel `json:"chat"`
}

type TLGUpdateResponseResultChatChannel struct {
	Id   int64  `json:"id"`
	Name string `json:"title"`
}

type TLGUpdateResponseResultFrom struct {
	Id   int64  `json:"id"`
	Name string `json:"first_name"`
}

func (rano *Rano) getUpdates() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
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
	// fmt.Println("body", string(body))
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
			author := &Author{
				Id:   result.Message.From.Id,
				Name: result.Message.From.Name,
			}
			group := &Group{
				Id:   result.Message.ChatChannel.Id,
				Name: result.Message.ChatChannel.Name,
			}
			message := &Message{
				UpdateId: result.UpdateId,
				From:     author,
				Group:    group,
				Text:     result.Message.Text,
				Date:     time.Unix(result.Message.UnixDate, 0),
			}
			list = append(list, message)
		}
	}
	return list
}

type Author struct {
	Id   int64
	Name string
}

type Group struct {
	Id   int64
	Name string
}

type Message struct {
	UpdateId int64
	From     *Author
	Group    *Group
	Text     string
	Date     time.Time
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