package nwc24

import (
	"fmt"
)

type Multipart struct {
	filename    string
	fileData    []byte
	contentType ContentType
	content     string
	charset     Charset
}

func NewMultipart() *Multipart {
	return &Multipart{}
}

func (m *Multipart) AddFile(filename string, data []byte, ct ContentType) {
	m.filename = filename
	m.fileData = data
	m.contentType = ct
}

func (m *Multipart) SetText(text string, charset Charset) {
	m.content = text
	m.charset = charset
	m.contentType = PlainText
}

func (m *Multipart) SetContentType(ct ContentType) {
	m.contentType = ct
}

func (m *Multipart) toString() (string, error) {
	switch m.contentType {
	case Binary, WiiMessageBoard, Jpeg:
		encoded := Base64Encode(m.fileData)
		return fmt.Sprint(
			CRLF,
			"Content-Type: ", m.contentType, ";", CRLF,
			" name=", m.filename, CRLF,
			"Content-Transfer-Encoding: base64", CRLF,
			"Content-Disposition: attachment;", CRLF,
			" filename=", m.filename, CRLF, CRLF,
			encoded,
			CRLF,
			CRLF,
		), nil
	case PlainText:
		encoded := Base64Encode([]byte(m.content))

		return fmt.Sprint(
			CRLF,
			"Content-Type: ", m.contentType, "; ", "charset=", m.charset, CRLF,
			"Content-Transfer-Encoding: base64", CRLF, CRLF,
			encoded,
			CRLF,
			CRLF,
			CRLF,
			CRLF,
			CRLF,
		), nil
	default:
		return "", UnsupportedContentType
	}
}
