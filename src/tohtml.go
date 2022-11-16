package main

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

func Markdown2HTML(markdown []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", "https://api.github.com/markdown", bytes.NewBuffer(markdown))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "token "+GHAPIToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("api error: Status code: " + strconv.Itoa(resp.StatusCode))
	}

	return io.ReadAll(resp.Body)
}
