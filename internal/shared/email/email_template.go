package email

type TemplateType string

const (
	TemplateVerification  TemplateType = "verification"
	TemplateResetPassword TemplateType = "reset_password"
)

func GetSubject(t TemplateType, lang string) string {
	switch t {
	case TemplateVerification:
		return getVerificationSubject(lang)

	case TemplateResetPassword:
		return getResetPasswordSubject(lang)

	default:
		return "JamLink Notification"
	}
}
