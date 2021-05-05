// Package processor Handle image processing, including resizing, converting and stripping unwanted metadata.
package processor

import (
	"fmt"

	"github.com/agravelot/imageopti/config"
)

// Processor Define processor interface.
type Processor interface {
	Optimize(media []byte, originalFormat string, targetFormat string, quality, width int) ([]byte, string, error)
}

// New Processor factory from dynamic configurations.
func New(conf config.Config) (Processor, error) {
	if conf.Processor == "imaginary" {
		p, err := NewImaginary(conf)
		if err != nil {
			return nil, err
		}

		return p, nil
	}

	if conf.Processor == "local" {
		return &LocalProcessor{}, nil
	}

	if conf.Processor == "none" {
		return &NoneProcessor{}, nil
	}

	return nil, fmt.Errorf("unable to resolver given optimizer %s", conf.Processor)
}
