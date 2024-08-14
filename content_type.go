package nwc24

import "encoding/binary"

// ContentType represents the different content types a message can have.
type ContentType string

const (
	Binary               ContentType = "application/octet-stream"
	Jpeg                 ContentType = "image/jpeg"
	WiiPicture           ContentType = "image/x-wii-picture"
	WiiMessageBoard      ContentType = "application/x-wii-msgboard"
	MultipartAlternative ContentType = "multipart/alternative"
	MultipartMixed       ContentType = "multipart/mixed"
	MultipartRelated     ContentType = "multipart/related"
	PlainText            ContentType = "text/plain"
)

// Charset represents the different types of charsets a string can be
type Charset string

const (
	UTF16BE Charset = "utf-16be"
)

func (m *Message) SetContentType(ct ContentType) {
	m.contentType = ct
}

func UTF16ToString(uint16s []uint16) string {
	byteArray := make([]byte, len(uint16s)*2)
	for i, v := range uint16s {
		binary.BigEndian.PutUint16(byteArray[i*2:], v)
	}

	return string(byteArray)
}
