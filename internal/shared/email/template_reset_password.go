package email

func getResetPasswordSubject(lang string) string {
	switch lang {
	case "fr-FR":
		return "Réinitialiser votre mot de passe JamLink"
	default:
		return "Reset your JamLink password"
	}
}
