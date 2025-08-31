package notification

type MailMessage struct {
	Subject     string
	Body        string
	HTMLBody    string
	From        string
	ReplyTo     string
	CC          []string
	BCC         []string
	Attachments []Attachment
}

type SMSMessage struct {
	Body string
	From string
}

type Attachment struct {
	Name        string
	Content     []byte
	ContentType string
}

func NewMailMessage() *MailMessage {
	return &MailMessage{
		CC:          make([]string, 0),
		BCC:         make([]string, 0),
		Attachments: make([]Attachment, 0),
	}
}

func (m *MailMessage) SetSubject(subject string) *MailMessage {
	m.Subject = subject
	return m
}

func (m *MailMessage) SetBody(body string) *MailMessage {
	m.Body = body
	return m
}

func (m *MailMessage) SetHTMLBody(htmlBody string) *MailMessage {
	m.HTMLBody = htmlBody
	return m
}

func (m *MailMessage) SetFrom(from string) *MailMessage {
	m.From = from
	return m
}

func (m *MailMessage) SetReplyTo(replyTo string) *MailMessage {
	m.ReplyTo = replyTo
	return m
}

func (m *MailMessage) AddCC(cc string) *MailMessage {
	m.CC = append(m.CC, cc)
	return m
}

func (m *MailMessage) AddBCC(bcc string) *MailMessage {
	m.BCC = append(m.BCC, bcc)
	return m
}

func (m *MailMessage) AddAttachment(name string, content []byte, contentType string) *MailMessage {
	m.Attachments = append(m.Attachments, Attachment{
		Name:        name,
		Content:     content,
		ContentType: contentType,
	})
	return m
}

func NewSMSMessage() *SMSMessage {
	return &SMSMessage{}
}

func (s *SMSMessage) SetBody(body string) *SMSMessage {
	s.Body = body
	return s
}

func (s *SMSMessage) SetFrom(from string) *SMSMessage {
	s.From = from
	return s
}