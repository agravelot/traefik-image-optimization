package processor

// NoneProcessor dummy processor, using null pattern.
type NoneProcessor struct{}

// Optimize return same data from media.
func (lp *NoneProcessor) Optimize(media []byte, of string, tf string, q, w int) ([]byte, string, error) {
	return media, of, nil
}
