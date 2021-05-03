package processor

type NoneProcessor struct {
}

func (lp *NoneProcessor) Optimize(media []byte, originalFormat string, targetFormat string, quality int) ([]byte, error) {
	return media, nil
}
