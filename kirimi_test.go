package kirimi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Test helper function to create a test server
func createTestServer(statusCode int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}))
}

// createRequestBodyServer captures the decoded JSON request body and returns a fixed response
func createRequestBodyServer(t *testing.T, statusCode int, response interface{}, capturePath *string, captureBody *map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*capturePath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(captureBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}))
}

func okResponse() Response {
	return Response{Success: true, Data: json.RawMessage(`{}`), Message: "OK"}
}

func assertStringField(t *testing.T, body map[string]interface{}, key, want string) {
	t.Helper()
	got, ok := body[key]
	if !ok {
		t.Errorf("expected key %q in request body, got none", key)
		return
	}
	if got != want {
		t.Errorf("expected %s=%q, got %v", key, want, got)
	}
}

func assertAbsentField(t *testing.T, body map[string]interface{}, key string) {
	t.Helper()
	if _, ok := body[key]; ok {
		t.Errorf("expected key %q to be absent from request body", key)
	}
}

func assertAPIError(t *testing.T, err error, wantStatus int, wantMessage string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != wantStatus {
		t.Errorf("expected status code %d, got %d", wantStatus, apiErr.StatusCode)
	}
	if wantMessage != "" && apiErr.Message != wantMessage {
		t.Errorf("expected message %q, got %q", wantMessage, apiErr.Message)
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client.BaseURL != DefaultBaseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", DefaultBaseURL, client.BaseURL)
	}
	if client.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", client.HTTPClient.Timeout)
	}
}

func TestNewClientWithBaseURL(t *testing.T) {
	customURL := "https://custom.api.com"
	client := NewClientWithBaseURL(customURL)
	if client.BaseURL != customURL {
		t.Errorf("Expected BaseURL to be %s, got %s", customURL, client.BaseURL)
	}
}

func TestSetTimeout(t *testing.T) {
	client := NewClient()
	newTimeout := 60 * time.Second
	client.SetTimeout(newTimeout)
	if client.HTTPClient.Timeout != newTimeout {
		t.Errorf("Expected timeout to be %v, got %v", newTimeout, client.HTTPClient.Timeout)
	}
}

func TestHealthCheck(t *testing.T) {
	// Test successful health check
	successResponse := Response{
		Success: true,
		Data:    json.RawMessage(`{}`),
		Message: "Kirimi API v1",
	}

	server := createTestServer(200, successResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	resp, err := client.HealthCheck()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Message != "Kirimi API v1" {
		t.Errorf("Expected message 'Kirimi API v1', got %s", resp.Message)
	}
}

func TestGenerateOTP(t *testing.T) {
	// Test successful OTP generation
	successResponse := Response{
		Success: true,
		Data:    json.RawMessage(`{"phone":"628123456789","message":"OTP berhasil dikirim","expires_in":"5 menit"}`),
		Message: "OTP berhasil digenerate dan dikirim",
	}

	server := createTestServer(200, successResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	req := GenerateOTPRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Phone:    "628123456789",
		Secret:   "test-secret",
	}

	resp, err := client.GenerateOTP(req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Phone != "628123456789" {
		t.Errorf("Expected phone '628123456789', got %s", resp.Phone)
	}
	if resp.ExpiresIn != "5 menit" {
		t.Errorf("Expected expires_in '5 menit', got %s", resp.ExpiresIn)
	}
}

func TestValidateOTP(t *testing.T) {
	// Test successful OTP validation
	successResponse := Response{
		Success: true,
		Data:    json.RawMessage(`{"phone":"628123456789","verified":true,"verified_at":"2024-01-15T10:30:00.000Z"}`),
		Message: "OTP berhasil divalidasi",
	}

	server := createTestServer(200, successResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	req := ValidateOTPRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Phone:    "628123456789",
		OTP:      "123456",
		Secret:   "test-secret",
	}

	resp, err := client.ValidateOTP(req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Phone != "628123456789" {
		t.Errorf("Expected phone '628123456789', got %s", resp.Phone)
	}
	if !resp.Verified {
		t.Errorf("Expected verified to be true, got %v", resp.Verified)
	}
}

func TestSendMessage(t *testing.T) {
	// Test successful message sending
	successResponse := Response{
		Success: true,
		Data:    json.RawMessage(`{"message_length":25,"media_url":"https://example.com/image.jpg","has_media":true}`),
		Message: "Berhasil mengirim pesan dengan media",
	}

	server := createTestServer(200, successResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	req := SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  "Hello from test!",
		Secret:   "test-secret",
		MediaURL: "https://example.com/image.jpg",
	}

	resp, err := client.SendMessage(req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.MessageLength != 25 {
		t.Errorf("Expected message_length 25, got %d", resp.MessageLength)
	}
	if !resp.HasMedia {
		t.Errorf("Expected has_media to be true, got %v", resp.HasMedia)
	}
}

func TestSendMessageTooLong(t *testing.T) {
	client := NewClient()
	req := SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  string(make([]byte, MaxMessageLength+1)), // Message too long
		Secret:   "test-secret",
	}

	_, err := client.SendMessage(req)

	if err == nil {
		t.Error("Expected error for message too long, got nil")
	}
}

func TestAPIError(t *testing.T) {
	// Test API error response
	errorResponse := Response{
		Success: false,
		Data:    json.RawMessage(`{}`),
		Message: "Parameter tidak boleh kosong",
	}

	server := createTestServer(400, errorResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	req := GenerateOTPRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Phone:    "628123456789",
		Secret:   "test-secret",
	}

	_, err := client.GenerateOTP(req)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("Expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Parameter tidak boleh kosong" {
		t.Errorf("Expected message 'Parameter tidak boleh kosong', got %s", apiErr.Message)
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test IsBasicOrProPackage
	if !IsBasicOrProPackage(PackageBasic1) {
		t.Error("Expected PackageBasic1 to support OTP")
	}
	if !IsBasicOrProPackage(PackagePro1) {
		t.Error("Expected PackagePro1 to support OTP")
	}
	if IsBasicOrProPackage(PackageFree) {
		t.Error("Expected PackageFree to not support OTP")
	}
	if IsBasicOrProPackage(PackageLite1) {
		t.Error("Expected PackageLite1 to not support OTP")
	}

	// Test IsMediaSupportedPackage
	if IsMediaSupportedPackage(PackageFree) {
		t.Error("Expected PackageFree to not support media")
	}
	if !IsMediaSupportedPackage(PackageLite1) {
		t.Error("Expected PackageLite1 to support media")
	}
	if !IsMediaSupportedPackage(PackageBasic1) {
		t.Error("Expected PackageBasic1 to support media")
	}
	if !IsMediaSupportedPackage(PackagePro1) {
		t.Error("Expected PackagePro1 to support media")
	}

	// Test IsFreePackage
	if !IsFreePackage(PackageFree) {
		t.Error("Expected PackageFree to be free package")
	}
	if IsFreePackage(PackageLite1) {
		t.Error("Expected PackageLite1 to not be free package")
	}
	if IsFreePackage(PackageBasic1) {
		t.Error("Expected PackageBasic1 to not be free package")
	}
	if IsFreePackage(PackagePro1) {
		t.Error("Expected PackagePro1 to not be free package")
	}
}

func TestConstants(t *testing.T) {
	if DefaultBaseURL != "https://api.kirimi.id" {
		t.Errorf("Expected DefaultBaseURL to be 'https://api.kirimi.id', got %s", DefaultBaseURL)
	}
	if ContentType != "application/json" {
		t.Errorf("Expected ContentType to be 'application/json', got %s", ContentType)
	}
	if MaxMessageLength != 1200 {
		t.Errorf("Expected MaxMessageLength to be 1200, got %d", MaxMessageLength)
	}
}

// --- Wire format: canonical field names ---

func TestSendMessageSendsReceiver(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.SendMessage(SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  "hi",
		Secret:   "test-secret",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/send-message" {
		t.Errorf("expected path /v1/send-message, got %s", path)
	}
	assertStringField(t, body, "receiver", "628987654321")
	assertAbsentField(t, body, "phone")
}

func TestSendMessageFastSendsReceiver(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.SendMessageFast(SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  "hi",
		Secret:   "test-secret",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/send-message-fast" {
		t.Errorf("expected path /v1/send-message-fast, got %s", path)
	}
	assertStringField(t, body, "receiver", "628987654321")
	assertAbsentField(t, body, "phone")
}

func TestSendMessageFileSendsReceiver(t *testing.T) {
	var contentType string
	var receiver string
	var phone string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Errorf("expected multipart form, got %v", err)
		}
		receiver = r.FormValue("receiver")
		phone = r.FormValue("phone")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(okResponse())
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.SendMessageFile(SendMessageFileRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Secret:   "test-secret",
		File:     strings.NewReader("file-bytes"),
		Filename: "doc.pdf",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		t.Errorf("expected multipart Content-Type, got %s", contentType)
	}
	if receiver != "628987654321" {
		t.Errorf("expected receiver=628987654321, got %q", receiver)
	}
	if phone != "" {
		t.Errorf("expected phone to be absent, got %q", phone)
	}
}

func TestBroadcastMessageSendsNumbersArray(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.BroadcastMessage(BroadcastMessageRequest{
		UserCode: "USER123",
		Secret:   "test-secret",
		DeviceID: "DEVICE456",
		Label:    "promo",
		Numbers:  []string{"628111", "628222"},
		Message:  "hello all",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/broadcast-message" {
		t.Errorf("expected path /v1/broadcast-message, got %s", path)
	}
	assertAbsentField(t, body, "phones")
	assertStringField(t, body, "label", "promo")

	numbers, ok := body["numbers"].([]interface{})
	if !ok {
		t.Fatalf("expected numbers to be a JSON array, got %T", body["numbers"])
	}
	if len(numbers) != 2 || numbers[0] != "628111" || numbers[1] != "628222" {
		t.Errorf("unexpected numbers array: %v", numbers)
	}
}

func TestSaveContactSendsNamaNomor(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.SaveContact(SaveContactRequest{
		UserCode: "USER123",
		Secret:   "test-secret",
		Nama:     "Budi",
		Nomor:    "6281234567890",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/save-contact" {
		t.Errorf("expected path /v1/save-contact, got %s", path)
	}
	assertStringField(t, body, "nama", "Budi")
	assertStringField(t, body, "nomor", "6281234567890")
	assertAbsentField(t, body, "name")
	assertAbsentField(t, body, "phone")
}

func TestSaveContactsBulk(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.SaveContactsBulk("USER123", "test-secret", []BulkContact{
		{Nama: "Budi", Nomor: "628111"},
		{Nama: "Siti", Nomor: "628222"},
	}, "DEVICE456")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/save-contacts-bulk" {
		t.Errorf("expected path /v1/save-contacts-bulk, got %s", path)
	}
	assertStringField(t, body, "device_id", "DEVICE456")

	contacts, ok := body["contacts"].([]interface{})
	if !ok {
		t.Fatalf("expected contacts to be a JSON array, got %T", body["contacts"])
	}
	if len(contacts) != 2 {
		t.Fatalf("expected 2 contacts, got %d", len(contacts))
	}
	first, ok := contacts[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected contact to be an object, got %T", contacts[0])
	}
	assertStringField(t, first, "nama", "Budi")
	assertStringField(t, first, "nomor", "628111")
}

func TestSendWabaMessage(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.SendWabaMessage(SendWabaMessageRequest{
		UserCode:     "USER123",
		Secret:       "test-secret",
		WabaID:       "WABA789",
		To:           "6281234567890",
		TemplateName: "order_update",
		Variables:    []string{"Budi", "12345"},
		Header:       &WabaTemplateHeader{Type: "image", Link: "https://example.com/h.jpg"},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/waba/send-message" {
		t.Errorf("expected path /v1/waba/send-message, got %s", path)
	}
	assertStringField(t, body, "waba_id", "WABA789")
	assertStringField(t, body, "to", "6281234567890")
	assertStringField(t, body, "template_name", "order_update")
	assertAbsentField(t, body, "device_id")
	assertAbsentField(t, body, "phone")

	vars, ok := body["variables"].([]interface{})
	if !ok || len(vars) != 2 {
		t.Fatalf("expected variables array of 2, got %v", body["variables"])
	}
	header, ok := body["header"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected header object, got %T", body["header"])
	}
	assertStringField(t, header, "type", "image")
	assertStringField(t, header, "link", "https://example.com/h.jpg")
}

func TestSendOtpV2Methods(t *testing.T) {
	cases := []struct {
		name string
		req  SendOtpV2Request
		want map[string]string
	}{
		{
			name: "whatsapp",
			req: SendOtpV2Request{
				UserCode: "USER123",
				Secret:   "test-secret",
				Phone:    "6281234567890",
				Method:   OtpMethodWhatsApp,
				AppName:  "Kirimi.id",
			},
			want: map[string]string{"method": "whatsapp", "phone": "6281234567890", "app_name": "Kirimi.id"},
		},
		{
			name: "device",
			req: SendOtpV2Request{
				UserCode:      "USER123",
				Secret:        "test-secret",
				Phone:         "6281234567890",
				Method:        OtpMethodDevice,
				DeviceID:      "DEVICE456",
				CustomMessage: "Kode OTP Anda: {{otp}}",
			},
			want: map[string]string{"method": "device", "device_id": "DEVICE456", "custom_message": "Kode OTP Anda: {{otp}}"},
		},
		{
			name: "waba_user",
			req: SendOtpV2Request{
				UserCode:     "USER123",
				Secret:       "test-secret",
				Phone:        "6281234567890",
				Method:       OtpMethodWabaUser,
				WabaID:       "WABA789",
				TemplateName: "otp_login",
			},
			want: map[string]string{"method": "waba_user", "waba_id": "WABA789", "template_name": "otp_login"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var path string
			var body map[string]interface{}
			server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
			defer server.Close()

			client := NewClientWithBaseURL(server.URL)
			if _, err := client.SendOtpV2(tc.req); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if path != "/v2/otp/send" {
				t.Errorf("expected path /v2/otp/send, got %s", path)
			}
			for k, v := range tc.want {
				assertStringField(t, body, k, v)
			}
		})
	}
}

// --- WABA ---

func TestWabaReply(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.WabaReply("USER123", "test-secret", "WABA789", "6281234567890", WabaReplyMessage{
		Type: "text",
		Text: "halo",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/waba/messages/reply" {
		t.Errorf("expected path /v1/waba/messages/reply, got %s", path)
	}
	assertStringField(t, body, "waba_id", "WABA789")
	assertStringField(t, body, "to", "6281234567890")

	msg, ok := body["message"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected message object, got %T", body["message"])
	}
	assertStringField(t, msg, "type", "text")
	assertStringField(t, msg, "text", "halo")
	assertAbsentField(t, msg, "media_url")
	assertAbsentField(t, msg, "interactive")
}

func TestWabaReplyMedia(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.WabaReply("USER123", "test-secret", "WABA789", "6281234567890", WabaReplyMessage{
		Type:     "document",
		MediaURL: "https://example.com/doc.pdf",
		Caption:  "invoice",
		Filename: "invoice.pdf",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	msg, ok := body["message"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected message object, got %T", body["message"])
	}
	assertStringField(t, msg, "type", "document")
	assertStringField(t, msg, "media_url", "https://example.com/doc.pdf")
	assertStringField(t, msg, "caption", "invoice")
	assertStringField(t, msg, "filename", "invoice.pdf")
	assertAbsentField(t, msg, "text")
}

func TestWabaConversations(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.WabaConversations("USER123", "test-secret", 50, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/waba/conversations" {
		t.Errorf("expected path /v1/waba/conversations, got %s", path)
	}
	if body["limit"] != float64(50) {
		t.Errorf("expected limit 50, got %v", body["limit"])
	}
	if body["page"] != float64(2) {
		t.Errorf("expected page 2, got %v", body["page"])
	}
}

func TestWabaTemplatesSync(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.WabaTemplatesSync("USER123", "test-secret", "WABA789")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/waba/templates/sync" {
		t.Errorf("expected path /v1/waba/templates/sync, got %s", path)
	}
	assertStringField(t, body, "waba_id", "WABA789")
	assertStringField(t, body, "user_code", "USER123")
	assertStringField(t, body, "secret", "test-secret")
}

func TestWabaSendOtp(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.WabaSendOtp("USER123", "test-secret", "WABA789", "6281234567890", "otp_login")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/waba/send-otp" {
		t.Errorf("expected path /v1/waba/send-otp, got %s", path)
	}
	assertStringField(t, body, "waba_id", "WABA789")
	assertStringField(t, body, "to", "6281234567890")
	assertStringField(t, body, "template_name", "otp_login")
}

func TestWabaVerifyOtp(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.WabaVerifyOtp("USER123", "test-secret", "WABA789", "6281234567890", "123456")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/waba/verify-otp" {
		t.Errorf("expected path /v1/waba/verify-otp, got %s", path)
	}
	assertStringField(t, body, "waba_id", "WABA789")
	assertStringField(t, body, "to", "6281234567890")
	assertStringField(t, body, "otp_code", "123456")
}

// --- Devices ---

func TestCreateDevice(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.CreateDevice("USER123", "test-secret", 4, "VOUCHER10")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/create-device" {
		t.Errorf("expected path /v1/create-device, got %s", path)
	}
	if body["package_id"] != float64(4) {
		t.Errorf("expected package_id 4, got %v", body["package_id"])
	}
	assertStringField(t, body, "voucher_code", "VOUCHER10")
	assertStringField(t, body, "user_code", "USER123")
	assertStringField(t, body, "secret", "test-secret")
}

func TestConnectDevice(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.ConnectDevice("USER123", "test-secret", "DEVICE456")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/connect-device" {
		t.Errorf("expected path /v1/connect-device, got %s", path)
	}
	assertStringField(t, body, "device_id", "DEVICE456")
}

func TestRenewDevice(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.RenewDevice("USER123", "test-secret", "DEVICE456", "pkg-pro", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/renew-device" {
		t.Errorf("expected path /v1/renew-device, got %s", path)
	}
	assertStringField(t, body, "device_id", "DEVICE456")
	assertStringField(t, body, "package_id", "pkg-pro")
	assertAbsentField(t, body, "voucher_code")
}

// --- OTP Reverse ---

func TestOtpReverseCreate(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.OtpReverseCreate("USER123", "test-secret", OtpReverseCreateRequest{
		Phone:         "6281234567890",
		DeviceID:      "DEVICE456",
		CallbackURL:   "https://example.com/cb",
		CustomMessage: "Kirim {{token}} dari {{phone}}",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v2/otp-reverse/create" {
		t.Errorf("expected path /v2/otp-reverse/create, got %s", path)
	}
	assertStringField(t, body, "phone", "6281234567890")
	assertStringField(t, body, "device_id", "DEVICE456")
	assertStringField(t, body, "callback_url", "https://example.com/cb")
	assertStringField(t, body, "custom_message", "Kirim {{token}} dari {{phone}}")
	assertStringField(t, body, "user_code", "USER123")
	assertStringField(t, body, "secret", "test-secret")
}

func TestOtpReverseStatus(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.OtpReverseStatus("USER123", "test-secret", "TOKEN123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v2/otp-reverse/status" {
		t.Errorf("expected path /v2/otp-reverse/status, got %s", path)
	}
	assertStringField(t, body, "token", "TOKEN123")
}

// --- Deposits ---

func TestCreateDeposit(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.CreateDeposit("USER123", "test-secret", 50000)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/create-deposit" {
		t.Errorf("expected path /v1/create-deposit, got %s", path)
	}
	if body["nominal"] != float64(50000) {
		t.Errorf("expected nominal 50000, got %v", body["nominal"])
	}
}

func TestDepositStatus(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.DepositStatus("USER123", "test-secret", "REF123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/deposit-status" {
		t.Errorf("expected path /v1/deposit-status, got %s", path)
	}
	assertStringField(t, body, "ref", "REF123")
}

func TestCancelDeposit(t *testing.T) {
	var path string
	var body map[string]interface{}
	server := createRequestBodyServer(t, 200, okResponse(), &path, &body)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.CancelDeposit("USER123", "test-secret", "REF123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/v1/cancel-deposit" {
		t.Errorf("expected path /v1/cancel-deposit, got %s", path)
	}
	assertStringField(t, body, "ref", "REF123")
}

// --- Errors on new endpoints ---

func TestNewEndpointsReturnAPIError(t *testing.T) {
	errorResponse := Response{
		Success: false,
		Data:    json.RawMessage(`{}`),
		Message: "Nominal minimal 100",
	}

	cases := []struct {
		name string
		call func(c *Client) error
	}{
		{"create-device", func(c *Client) error { _, e := c.CreateDevice("U", "S", 1, ""); return e }},
		{"connect-device", func(c *Client) error { _, e := c.ConnectDevice("U", "S", "D"); return e }},
		{"renew-device", func(c *Client) error { _, e := c.RenewDevice("U", "S", "D", 1, ""); return e }},
		{"save-contacts-bulk", func(c *Client) error {
			_, e := c.SaveContactsBulk("U", "S", []BulkContact{{Nama: "a", Nomor: "b"}}, "")
			return e
		}},
		{"waba-reply", func(c *Client) error {
			_, e := c.WabaReply("U", "S", "W", "T", WabaReplyMessage{Type: "text", Text: "x"})
			return e
		}},
		{"waba-conversations", func(c *Client) error { _, e := c.WabaConversations("U", "S", 50, 1); return e }},
		{"waba-templates-sync", func(c *Client) error { _, e := c.WabaTemplatesSync("U", "S", "W"); return e }},
		{"waba-send-otp", func(c *Client) error { _, e := c.WabaSendOtp("U", "S", "W", "T", "tpl"); return e }},
		{"waba-verify-otp", func(c *Client) error { _, e := c.WabaVerifyOtp("U", "S", "W", "T", "1234"); return e }},
		{"otp-reverse-create", func(c *Client) error {
			_, e := c.OtpReverseCreate("U", "S", OtpReverseCreateRequest{Phone: "p", DeviceID: "d"})
			return e
		}},
		{"otp-reverse-status", func(c *Client) error { _, e := c.OtpReverseStatus("U", "S", "TOKEN"); return e }},
		{"create-deposit", func(c *Client) error { _, e := c.CreateDeposit("U", "S", 50); return e }},
		{"deposit-status", func(c *Client) error { _, e := c.DepositStatus("U", "S", "REF"); return e }},
		{"cancel-deposit", func(c *Client) error { _, e := c.CancelDeposit("U", "S", "REF"); return e }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := createTestServer(400, errorResponse)
			defer server.Close()

			client := NewClientWithBaseURL(server.URL)
			assertAPIError(t, tc.call(client), 400, "Nominal minimal 100")
		})
	}
}

func TestUnauthorizedIsAPIError(t *testing.T) {
	errorResponse := Response{
		Success: false,
		Data:    json.RawMessage(`null`),
		Message: "Secret salah",
	}

	server := createTestServer(401, errorResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.SendMessage(SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  "hi",
		Secret:   "wrong",
	})
	assertAPIError(t, err, 401, "Secret salah")
}
