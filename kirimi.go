package kirimi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Constants
const (
	DefaultBaseURL = "https://api.kirimi.id"
	APIVersion     = "v1"
	ContentType    = "application/json"
	MaxMessageLength = 1200
)

// Package Types
const (
	PackageFree  = 1
	PackageLite1 = 2
	PackageBasic1 = 3
	PackagePro1  = 4
	PackageLite2 = 6
	PackageBasic2 = 7
	PackagePro2  = 8
	PackageLite3 = 9
	PackageBasic3 = 10
	PackagePro3  = 11
)

// Client represents the Kirimi API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Kirimi API client
func NewClient() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithBaseURL creates a new client with custom base URL
func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}

// Response represents the standard API response format
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// GenerateOTPRequest represents the request for generating OTP
type GenerateOTPRequest struct {
	UserCode string `json:"user_code"`
	DeviceID string `json:"device_id"`
	Phone    string `json:"phone"`
	Secret   string `json:"secret"`
}

// GenerateOTPResponse represents the response data for OTP generation
type GenerateOTPResponse struct {
	Phone     string `json:"phone"`
	Message   string `json:"message"`
	ExpiresIn string `json:"expires_in"`
}

// ValidateOTPRequest represents the request for validating OTP
type ValidateOTPRequest struct {
	UserCode string `json:"user_code"`
	DeviceID string `json:"device_id"`
	Phone    string `json:"phone"`
	OTP      string `json:"otp"`
	Secret   string `json:"secret"`
}

// ValidateOTPResponse represents the response data for OTP validation
type ValidateOTPResponse struct {
	Phone      string    `json:"phone"`
	Verified   bool      `json:"verified"`
	VerifiedAt time.Time `json:"verified_at"`
}

// SendMessageRequest represents the request for sending message
type SendMessageRequest struct {
	UserCode string `json:"user_code"`
	DeviceID string `json:"device_id"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
	Secret   string `json:"secret"`
	MediaURL string `json:"media_url,omitempty"`
}

// SendMessageResponse represents the response data for sending message
type SendMessageResponse struct {
	MessageLength int    `json:"message_length"`
	MediaURL      string `json:"media_url,omitempty"`
	HasMedia      bool   `json:"has_media"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// APIError represents an API error
type APIError struct {
	StatusCode int
	Message    string
	Response   *Response
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

// makeRequest makes an HTTP request to the API
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*Response, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, APIVersion, endpoint)
	
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", ContentType)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &apiResp, &APIError{
			StatusCode: resp.StatusCode,
			Message:    apiResp.Message,
			Response:   &apiResp,
		}
	}

	return &apiResp, nil
}

// makeGetRequest makes a GET request (for health check)
func (c *Client) makeGetRequest(endpoint string) (*Response, error) {
	url := c.BaseURL
	if endpoint != "" {
		url = fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &apiResp, nil
}

// GenerateOTP generates an OTP and sends it to the specified phone number
func (c *Client) GenerateOTP(req GenerateOTPRequest) (*GenerateOTPResponse, error) {
	resp, err := c.makeRequest("POST", "generate-otp", req)
	if err != nil {
		return nil, err
	}

	var otpResp GenerateOTPResponse
	if resp.Data != nil {
		dataBytes, err := json.Marshal(resp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}
		if err := json.Unmarshal(dataBytes, &otpResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal OTP response: %w", err)
		}
	}

	return &otpResp, nil
}

// ValidateOTP validates the provided OTP
func (c *Client) ValidateOTP(req ValidateOTPRequest) (*ValidateOTPResponse, error) {
	resp, err := c.makeRequest("POST", "validate-otp", req)
	if err != nil {
		return nil, err
	}

	var validateResp ValidateOTPResponse
	if resp.Data != nil {
		dataBytes, err := json.Marshal(resp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}
		if err := json.Unmarshal(dataBytes, &validateResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal validate response: %w", err)
		}
	}

	return &validateResp, nil
}

// SendMessage sends a WhatsApp message with optional media
func (c *Client) SendMessage(req SendMessageRequest) (*SendMessageResponse, error) {
	// Validate message length
	if len(req.Message) > MaxMessageLength {
		return nil, fmt.Errorf("message length exceeds maximum of %d characters", MaxMessageLength)
	}

	resp, err := c.makeRequest("POST", "send-message", req)
	if err != nil {
		return nil, err
	}

	var msgResp SendMessageResponse
	if resp.Data != nil {
		dataBytes, err := json.Marshal(resp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}
		if err := json.Unmarshal(dataBytes, &msgResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message response: %w", err)
		}
	}

	return &msgResp, nil
}

// HealthCheck checks the API health status
func (c *Client) HealthCheck() (*HealthCheckResponse, error) {
	resp, err := c.makeGetRequest("")
	if err != nil {
		return nil, err
	}

	var healthResp HealthCheckResponse
	if resp.Data != nil {
		dataBytes, err := json.Marshal(resp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}
		if err := json.Unmarshal(dataBytes, &healthResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal health response: %w", err)
		}
	}

	// Handle different response formats
	if healthResp.Status == "" && healthResp.Message == "" {
		healthResp.Message = resp.Message
	}

	return &healthResp, nil
}

// Helper functions

// IsBasicOrProPackage checks if the package supports OTP features
func IsBasicOrProPackage(packageID int) bool {
	return packageID == PackageBasic1 || packageID == PackagePro1 ||
		packageID == PackageBasic2 || packageID == PackagePro2 ||
		packageID == PackageBasic3 || packageID == PackagePro3
}

// IsMediaSupportedPackage checks if the package supports media features
func IsMediaSupportedPackage(packageID int) bool {
	return packageID != PackageFree
}

// IsFreePackage checks if the package is free (text only with watermark)
func IsFreePackage(packageID int) bool {
	return packageID == PackageFree
}