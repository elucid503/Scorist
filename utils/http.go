package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetAndDecode(url string, target interface{}) error {

	resp, err := http.Get(url)

	if err != nil {

		return err

	}

	defer resp.Body.Close()

	return DecodeJSON(resp.Body, target)

}

func DecodeJSON(body io.ReadCloser, target interface{}) error {

	decoder := json.NewDecoder(body)

	return decoder.Decode(target)

}
