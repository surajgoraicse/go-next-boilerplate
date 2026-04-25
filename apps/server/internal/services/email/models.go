package email

// EmailRequest represents the request payload for sending an email
type EmailRequest struct {
	ID        int                    `json:"id"`
	Subject   string                 `json:"subject"`
	Recipient string                 `json:"recipient"`
	Body      map[string]interface{} `json:"body"`
	CC        []string               `json:"cc"`
	BCC       []string               `json:"bcc"`
	Self      bool                   `json:"self"`
}

// EmailDetails represents the email details in the success response
type EmailDetails struct {
	ID        int                    `json:"id"`
	Subject   string                 `json:"subject"`
	Recipient string                 `json:"recipient"`
	Body      map[string]interface{} `json:"body"`
	CC        []string               `json:"cc"`
	BCC       []string               `json:"bcc"`
	Self      bool                   `json:"self"`
	Provider  string                 `json:"provider"`
}

// EmailSuccessResponse represents the success response from the email service
type EmailSuccessResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Details EmailDetails `json:"details"`
}

// ErrorDetail represents a single validation error detail
type ErrorDetail struct {
	Loc  []string `json:"loc"`
	Msg  string   `json:"msg"`
	Type string   `json:"type"`
}

// EmailErrorResponse represents the error response from the email service
type EmailErrorResponse struct {
	Detail []ErrorDetail `json:"detail"`
}
