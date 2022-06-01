package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func APIHandler(url, str string) []byte {

	payload := strings.NewReader(str)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return body
}
