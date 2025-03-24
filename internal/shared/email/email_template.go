package email

type TemplateType string

const (
	TemplateVerification TemplateType = "verification"
)

func GetSubject(t TemplateType, lang string) string {
	switch t {
	case TemplateVerification:
		return getVerificationSubject(lang)
	default:
		return "JamLink Notification"
	}
}
