package processor

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// TODO imaginary url = http://imaginary:9000
type ImaginaryProcessor struct {
}

func (ip *ImaginaryProcessor) Optimize(media []byte, origialFormat string, targetFormat string, quality int) ([]byte, error) {

	url := "http://imaginary:9000/convert?type=webp&field=file"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	fileWritrer, err := writer.CreateFormFile("file", "tmp.jpg")
	if err != nil {
		return nil, err
	}

	_, err = fileWritrer.Write(media)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil
	}

	return body, nil
}
