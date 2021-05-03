package processor

type LocalProcessor struct {
}

func (lp *LocalProcessor) Optimize(media []byte, originalFormat string, targetFormat string, quality int) ([]byte, error) {

	// newImage, err := bimg.NewImage(media).Convert(bimg.WEBP)
	// if err != nil {
	// 	return nil, err
	// }

	// if bimg.NewImage(newImage).Type() == "webp" {
	// 	fmt.Fprintln(os.Stderr, "The image was converted into webp")
	// }

	return media, nil
}
