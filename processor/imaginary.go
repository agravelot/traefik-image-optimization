// Package processor Add support for imaginary imge processing service.
package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/agravelot/imageopti/config"
)

const httpTimeout = 5 * time.Second

type pipelineOperationParams struct {
	Font      string  `json:"font,omitempty"`
	Height    int     `json:"height,omitempty"`
	Opacity   float64 `json:"opacity,omitempty"`
	Rotate    int     `json:"rotate,omitempty"`
	Text      string  `json:"text,omitempty"`
	Textwidth int     `json:"textwidth,omitempty"`
	Type      string  `json:"type,omitempty"`
	Width     int     `json:"width,omitempty"`
	StripMeta bool    `json:"stripmeta,omitempty"`
}

type pipelineOperation struct {
	Operation string                  `json:"operation"`
	Params    pipelineOperationParams `json:"params"`
}

// ImaginaryProcessor define imaginary processor settings.
type ImaginaryProcessor struct {
	URL    string
	client http.Client
}

func isValidURL(s string) error {
	if s == "" {
		return fmt.Errorf("url cannot be empty")
	}

	u, err := url.ParseRequestURI(s)
	if err != nil {
		return fmt.Errorf("unable to parse imaginary url: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unvalid imaginary scheme")
	}

	return nil
}

// NewImaginary instantiate a new imaginary instance with given config.
func NewImaginary(conf config.Config) (*ImaginaryProcessor, error) {
	err := isValidURL(conf.Imaginary.URL)
	if err != nil {
		return nil, err
	}

	return &ImaginaryProcessor{
		client: http.Client{
			Timeout: httpTimeout,
		},
		URL: conf.Imaginary.URL,
	}, nil
}

// Optimize method to process image with imaginary with given parameters.
func (ip *ImaginaryProcessor) Optimize(media []byte, of string, tf string, q, width int) ([]byte, string, error) {
	ope := []pipelineOperation{
		{Operation: "convert", Params: pipelineOperationParams{Type: "webp", StripMeta: true}},
	}

	if width > 0 {
		ope = append(ope, pipelineOperation{Operation: "resize", Params: pipelineOperationParams{Width: width}})
	}

	opString, err := json.Marshal(ope)
	if err != nil {
		return nil, "", fmt.Errorf("unable generate imaginary operations: %w", err)
	}

	u := fmt.Sprintf("%s/pipeline?operations=%s", ip.URL, url.QueryEscape(string(opString)))
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	fileWriter, err := writer.CreateFormFile("file", "tmp.jpg")
	if err != nil {
		return nil, "", fmt.Errorf("unable to create file to imaginary file writer: %w", err)
	}

	_, err = fileWriter.Write(media)
	if err != nil {
		return nil, "", fmt.Errorf("unable to write file to imaginary file writer: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("unable to close imaginary file writer: %w", err)
	}

	req, err := http.NewRequest(method, u, payload)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create imaginary request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := ip.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("unable to send imaginary request: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", fmt.Errorf("unable to read imaginary response body: %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return nil, "", fmt.Errorf("unable to close imaginary body response: %w", err)
	}

	return body, "image/webp", nil
}
