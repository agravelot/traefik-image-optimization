package processor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/agravelot/image_optimizer/config"
)

type ImaginaryProcessor struct {
	Url string
}

func isValidUrl(s string) error {

	if s == "" {
		return fmt.Errorf("url cannot be empty")
	}

	u, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unvalid imaginary scheme")
	}

	return nil
}

func NewImaginary(conf config.Config) (*ImaginaryProcessor, error) {

	err := isValidUrl(conf.Imaginary.Url)
	if err != nil {
		return nil, err
	}

	return &ImaginaryProcessor{
		Url: conf.Imaginary.Url,
	}, nil
}

func (ip *ImaginaryProcessor) Optimize(media []byte, origialFormat string, targetFormat string, quality int) ([]byte, error) {

	url := fmt.Sprintf("%s/convert?type=webp&field=file", ip.Url)
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
