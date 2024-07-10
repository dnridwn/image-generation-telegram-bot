package main

import (
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://image.pollinations.ai/prompt"

func SendImageGenerationRequest(p string) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, p))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
