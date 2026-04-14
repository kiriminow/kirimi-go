# Kirimi Go SDK

Go SDK untuk Kirimi Console API - Platform WhatsApp messaging dengan fitur OTP, broadcast, WABA, dan banyak lagi.

## Instalasi

```bash
go get github.com/kiriminow/kirimi-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    kirimi "github.com/yolk/kirimi-go"
)

func main() {
    client := kirimi.NewClient()

    health, err := client.HealthCheck()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("API Status: %+v\n", health)
}
```

## API Reference

### Client

```go
client := kirimi.NewClient()
client := kirimi.NewClientWithBaseURL("https://custom-api.kirimi.id")
client.SetTimeout(60 * time.Second)
```

### Send Message

```go
resp, err := client.SendMessage(kirimi.SendMessageRequest{
    UserCode: "USER123",
    DeviceID: "DEVICE456",
    Phone:    "628123456789",
    Message:  "Hello!",
    Secret:   "your-secret",
    MediaURL: "https://example.com/image.jpg", // optional
})
```

### Send Message Fast (tanpa efek mengetik)

```go
resp, err := client.SendMessageFast(kirimi.SendMessageRequest{
    UserCode: "USER123",
    DeviceID: "DEVICE456",
    Phone:    "628123456789",
    Message:  "Hello cepat!",
    Secret:   "your-secret",
})
```

### Send Message File (multipart, max 50MB)

```go
f, _ := os.Open("document.pdf")
defer f.Close()

resp, err := client.SendMessageFile(kirimi.SendMessageFileRequest{
    UserCode: "USER123",
    DeviceID: "DEVICE456",
    Phone:    "628123456789",
    Secret:   "your-secret",
    File:     f,
    Filename: "document.pdf",
    Message:  "Ini dokumennya", // optional
})
```

### Send WABA Message

```go
resp, err := client.SendWabaMessage(kirimi.SendWabaMessageRequest{
    UserCode: "USER123",
    Secret:   "your-secret",
    DeviceID: "WABA_DEVICE",
    Phone:    "628123456789",
    Message:  "Hello via WABA!",
})
```

### List Devices

```go
resp, err := client.ListDevices("USER123", "your-secret")
```

### Device Status

```go
resp, err := client.DeviceStatus("USER123", "your-secret", "DEVICE456")
```

### Device Status Enhanced

```go
resp, err := client.DeviceStatusEnhanced("USER123", "your-secret", "DEVICE456")
```

### User Info

```go
resp, err := client.UserInfo("USER123", "your-secret")
```

### Save Contact

```go
resp, err := client.SaveContact(kirimi.SaveContactRequest{
    UserCode: "USER123",
    Secret:   "your-secret",
    Phone:    "628123456789",
    Name:     "Budi",    // optional
    Email:    "b@x.com", // optional
})
```

### Generate OTP (v1)

```go
resp, err := client.GenerateOTP(kirimi.GenerateOTPRequest{
    UserCode:         "USER123",
    DeviceID:         "DEVICE456",
    Phone:            "628123456789",
    Secret:           "your-secret",
    OtpLength:        6,            // optional
    OtpType:          "numeric",    // optional: numeric/alphabetic/alphanumeric
    CustomOtpMessage: "Kode OTP: ", // optional
})
```

### Validate OTP (v1)

```go
resp, err := client.ValidateOTP(kirimi.ValidateOTPRequest{
    UserCode: "USER123",
    DeviceID: "DEVICE456",
    Phone:    "628123456789",
    OTP:      "123456",
    Secret:   "your-secret",
})
```

### Send OTP v2

```go
resp, err := client.SendOtpV2(kirimi.SendOtpV2Request{
    UserCode:      "USER123",
    Secret:        "your-secret",
    Phone:         "628123456789",
    DeviceID:      "DEVICE456",
    Method:        "device",    // optional: device/waba
    AppName:       "MyApp",     // optional
    TemplateCode:  "tmpl_001",  // optional (waba)
    CustomMessage: "OTP kamu:", // optional (device)
})
```

### Verify OTP v2

```go
resp, err := client.VerifyOtpV2(kirimi.VerifyOtpV2Request{
    UserCode: "USER123",
    Secret:   "your-secret",
    Phone:    "628123456789",
    OtpCode:  "123456",
})
```

### Broadcast Message

```go
resp, err := client.BroadcastMessage(kirimi.BroadcastMessageRequest{
    UserCode: "USER123",
    Secret:   "your-secret",
    DeviceID: "DEVICE456",
    Phones:   "628111,628222,628333", // comma-separated
    Message:  "Promo hari ini!",
    Delay:    3, // optional: delay antar pesan (detik)
})
```

### List Deposits

```go
// status: "" / "paid" / "unpaid" / "expired"
resp, err := client.ListDeposits("USER123", "your-secret", "paid")
```

### List Packages

```go
resp, err := client.ListPackages("USER123", "your-secret")
```

### Health Check

```go
health, err := client.HealthCheck()
```

## Error Handling

```go
resp, err := client.SendMessage(req)
if err != nil {
    if apiErr, ok := err.(*kirimi.APIError); ok {
        fmt.Printf("API Error %d: %s\n", apiErr.StatusCode, apiErr.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Package Types

| Package | ID | Fitur |
|---------|----|-------|
| Free | 1 | Text only + watermark |
| Lite | 2, 6, 9 | Text + Media |
| Basic | 3, 7, 10 | Text + Media + OTP |
| Pro | 4, 8, 11 | Text + Media + OTP |

```go
kirimi.IsBasicOrProPackage(packageID)
kirimi.IsMediaSupportedPackage(packageID)
kirimi.IsFreePackage(packageID)
```

## Requirements

- Go 1.21 atau lebih baru
- Akun Kirimi Console dengan secret key valid
- Device yang sudah terdaftar dan connected

## License

MIT License
