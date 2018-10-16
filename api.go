package rano

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func (rano *Rano) sendRequest(method string, params map[string]string) (responseData map[string]interface{}, err error) {
	return rano.sendRequestWithFile(method, params, nil)
}

func (rano *Rano) SendRequestWithFile(method string, params map[string]string, fileParams map[string]string) (responseData map[string]interface{}, err error) {
	return rano.sendRequestWithFile(method, params, fileParams)
}

func (rano *Rano) sendRequestWithFile(method string, params map[string]string, fileParams map[string]string) (responseData map[string]interface{}, err error) {
	urlStr := fmt.Sprintf("%s%s", rano.baseUrl, method)
	var bodyBuffer bytes.Buffer
	w := multipart.NewWriter(&bodyBuffer)
	for key, value := range params {
		err = w.WriteField(key, value)
		if err != nil {
			return nil, err
		}
	}
	// log.Log("finish params")
	for key, filePath := range fileParams {
		fileObj, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer fileObj.Close()

		fileForm, err := w.CreateFormFile(key, "certificate")
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fileForm, fileObj)
		if err != nil {
			return nil, err
		}
	}
	w.Close()
	// log.Log("close things")

	request, err := http.NewRequest("POST", urlStr, &bodyBuffer)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", w.FormDataContentType())
	// log.Log("set header")

	response, err := rano.client.Do(request)
	if err != nil {
		return nil, err
	}
	// log.Log("done sent")
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	bodyString := string(body)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &data)
	if err != nil {
		return nil, err
	}
	// log.Log("done sent 2")

	ok := data["ok"].(bool)
	if !ok {
		err = errors.New(data["description"].(string))
		return nil, err
	}
	// log.Log("done sent 3")
	// fmt.Println("response", bodyString)
	return data, nil
}
