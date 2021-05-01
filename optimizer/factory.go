package optimizer

import "fmt"

type Optimizer interface {
	Optimize(media []byte, origialFormat string, targetFormat string, quality int) ([]byte, error)
}

func New(driver string) (Optimizer, error) {
	if driver == "local" {
		return &LocalOptimizer{}, nil
	}

	if driver == "imaginary" {
		return &ImaginaryOptimizer{}, nil
	}

	return nil, fmt.Errorf("unable to resolver given optimizer %s", driver)
}
