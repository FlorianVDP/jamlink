package email

func getVerificationSubject(lang string) string {
	switch lang {
	case "fr-FR":
		return "Vérifie ton compte JamLink"
	default:
		return "Verify your JamLink account"
	}
}
