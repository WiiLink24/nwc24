package nwc24

import "errors"

var (
	InvalidContentType     = errors.New("invalid content type")
	UnsupportedContentType = errors.New("unsupported content type. Check in a newer version")
)
