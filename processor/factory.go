package processor

import "fmt"

type Processor interface {
	Optimize(media []byte, origialFormat string, targetFormat string, quality int) ([]byte, error)
}

func New(driver string) (Processor, error) {
	if driver == "local" {
		return &LocalProcessor{}, nil
	}

	if driver == "imaginary" {
		return &ImaginaryProcessor{}, nil
	}

	return nil, fmt.Errorf("unable to resolver given optimizer %s", driver)
}
