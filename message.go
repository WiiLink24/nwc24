package nwc24

import (
	"fmt"
	"net/mail"
	"time"
)

const (
	CRLF = "\r\n"
	LF   = "\n"
)

type Message struct {
	to          *mail.Address
	from        *mail.Address
	subject     string
	contentType ContentType
	multipart   []*Multipart
	boundary    string
	tags        map[string]string
	content     string
	charset     Charset
}

func NewMessage(from, to *mail.Address) *Message {
	return &Message{to: to, from: from, tags: make(map[string]string)}
}

func (m *Message) SetSubject(subject string) {
	m.subject = subject
}

func (m *Message) AddMultipart(m_ ...*Multipart) {
	m.multipart = append(m.multipart, m_...)
}

func (m *Message) SetBoundary(b string) {
	m.boundary = b
}

func (m *Message) SetText(text string, charset Charset) {
	m.content = text
	m.charset = charset
}

func (m *Message) SetTag(key, value string) {
	m.tags[key] = value
}

func (m *Message) ToString() (string, error) {
	if m.multipart != nil && (m.contentType != MultipartAlternative && m.contentType != MultipartMixed && m.contentType != MultipartRelated) {
		return "", InvalidContentType
	}

	data := fmt.Sprint(
		"Content-Type: text/plain", CRLF, CRLF,
		"Date: ", time.Now().Format("02 Jan 2006 15:04:05 -0700"), CRLF,
		"From: ", m.from.Address, CRLF,
		"To: ", m.to.Address, CRLF,
		// TODO: Change this to dynamic
		"Message-ID: <776DCLBHYHD.2QBO4Y3I2Y04S@JavaMail.w9999999900000000@rc24.xyz>\r\n",
		"Subject: ", m.subject, CRLF,
		"MIME-Version: 1.0", CRLF,
	)

	for k, v := range m.tags {
		data += fmt.Sprintf("%s: %s%s", k, v, CRLF)
	}

	if m.contentType == MultipartMixed {
		data += fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s%s", m.boundary, CRLF)
		data += "Content-Transfer-Encoding: base64" + CRLF

		// Form the multiple parts now
		for _, multipart := range m.multipart {
			data += "--" + m.boundary

			content, err := multipart.toString()
			if err != nil {
				return "", err
			}

			data += content
		}

		// Finalize the message
		data += "--" + m.boundary + "--" + CRLF + CRLF
	}

	return data, nil
}

func CreateMessageToSend(boundary string, messages ...*Message) (string, error) {
	base := fmt.Sprint(
		boundary,
		CRLF,
		"Content-Type: text/plain",
		CRLF,
		CRLF,
		"This part is ignored.",
		CRLF,
		CRLF,
		CRLF,
		CRLF,
	)

	for _, message := range messages {
		str, err := message.ToString()
		if err != nil {
			return "", err
		}

		base += boundary + CRLF
		base += str
	}

	base += boundary + "--" + CRLF

	return base, nil
}
