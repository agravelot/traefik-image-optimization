package processor

type NoneProcessor struct {
}

func (lp *NoneProcessor) Optimize(media []byte, originalFormat string, targetFormat string, quality, width int) ([]byte, string, error) {
	return media, originalFormat, nil
}
