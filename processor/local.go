package processor

// LocalProcessor process images directly in traefik itself, unsupported with interpreter limitations.
type LocalProcessor struct{}

// Optimize optimize image with given params.
func (lp *LocalProcessor) Optimize(media []byte, of string, tf string, q, w int) ([]byte, string, error) {
	// newImage, err := bimg.NewImage(media).Convert(bimg.WEBP)
	// if err != nil {
	// 	return nil, err
	// }
	return media, tf, nil
}
