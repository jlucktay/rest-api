package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func getResponse(get url.URL) (io.Reader, error) {
	req, errReq := http.NewRequest("GET", get.String(), nil)
	if errReq != nil {
		return nil, errReq
	}

	res, errDo := http.DefaultClient.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("URL '%s': status code error: [%d] %s", get.String(), res.StatusCode, res.Status)
	}

	buf := new(bytes.Buffer)
	if _, errRead := buf.ReadFrom(res.Body); errRead != nil {
		return nil, errRead
	}

	return buf, nil
}
