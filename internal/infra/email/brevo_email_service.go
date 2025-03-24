package emailinfra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"jamlink-backend/internal/shared/email"
)

type BrevoEmailService struct {
	apiKey      string
	senderName  string
	senderEmail string
	baseURL     string
	templateDir string
}

func NewBrevoEmailService() *BrevoEmailService {
	return &BrevoEmailService{
		apiKey:      os.Getenv("BREVO_API_KEY"),
		senderName:  os.Getenv("BREVOS_SENDER_NAME"),
		senderEmail: os.Getenv("BREVO_SENDER_EMAIL"),
		baseURL:     os.Getenv("FRONTEND_VERIFY_URL"),
		templateDir: "internal/shared/email/templates",
	}
}

func (s *BrevoEmailService) Send(to string, templateType email.TemplateType, lang string, data map[string]string) error {
	htmlContent, err := s.renderTemplate(templateType, data, lang)
	if err != nil {
		return err
	}

	subject := email.GetSubject(templateType, lang)

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  s.senderName,
			"email": s.senderEmail,
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject":     subject,
		"htmlContent": htmlContent,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("brevo error: %s", resp.Status)
	}

	return nil
}

func (s *BrevoEmailService) renderTemplate(templateType email.TemplateType, data any, lang string) (string, error) {
	path := filepath.Join(s.templateDir, lang, string(templateType)+".html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
