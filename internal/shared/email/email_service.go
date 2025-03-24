package email

type EmailService interface {
	Send(to string, template TemplateType, lang string, data map[string]string) error
}
