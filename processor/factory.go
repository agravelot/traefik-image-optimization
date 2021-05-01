package processor

import (
	"fmt"

	"github.com/agravelot/image_optimizer/config"
)

type Processor interface {
	Optimize(media []byte, origialFormat string, targetFormat string, quality int) ([]byte, error)
}

func New(conf config.Config) (Processor, error) {
	if conf.Processor == "local" {
		return &LocalProcessor{}, nil
	}

	if conf.Processor == "imaginary" {
		return NewImaginary(conf), nil
	}

	return nil, fmt.Errorf("unable to resolver given optimizer %s", conf.Processor)
}
