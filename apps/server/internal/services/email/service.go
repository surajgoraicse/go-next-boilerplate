package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendEmail(templateID int, to, subject string, body map[string]interface{}, cc, bcc []string, self bool) (*EmailSuccessResponse, error)
}

// EmailServiceImpl implements the EmailService interface
type EmailServiceImpl struct {
	config *config.Config
}

// NewEmailService creates a new email service instance
func NewEmailService(cfg *config.Config) EmailSender {
	return &EmailServiceImpl{
		config: cfg,
	}
}

// SendEmail sends an email using the TeamShiksha email service
func (e *EmailServiceImpl) SendEmail(templateID int, to, subject string, body map[string]interface{}, cc, bcc []string, self bool) (*EmailSuccessResponse, error) {
	emailServiceBaseURL := e.config.EmailServiceBaseURL
	emailServiceToken := e.config.EmailServiceToken

	fullURL, err := url.JoinPath(emailServiceBaseURL, "/email")
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	reqBody := EmailDetails{
		ID:        templateID,
		Subject:   subject,
		Recipient: to,
		Body:      body,
		CC:        cc,
		BCC:       bcc,
		Self:      self,
		Provider:  e.config.EmailProvider,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", emailServiceToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp EmailErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d and failed to parse error response: %s", resp.StatusCode, string(respBody))
		}

		var errorMsg string
		if len(errorResp.Detail) > 0 {
			errorMsg = fmt.Sprintf("validation error: %s", errorResp.Detail[0].Msg)
		} else {
			errorMsg = fmt.Sprintf("request failed with status %d", resp.StatusCode)
		}

		return nil, fmt.Errorf("%s", errorMsg)
	}

	var successResp EmailSuccessResponse
	if err := json.Unmarshal(respBody, &successResp); err != nil {
		return nil, fmt.Errorf("failed to parse success response: %w", err)
	}

	return &successResp, nil
}
