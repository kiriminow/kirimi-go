# Kirimi Go SDK

Go SDK untuk Kirimi Console API - Platform WhatsApp messaging dengan fitur OTP dan media support.

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
    // Membuat client baru
    client := kirimi.NewClient()
    
    // Health check
    health, err := client.HealthCheck()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("API Status: %+v\n", health)
    
    // Kirim pesan
    msgReq := kirimi.SendMessageRequest{
        UserCode: "USER123",
        DeviceID: "DEVICE456",
        Receiver: "628123456789",
        Message:  "Hello from Kirimi Go SDK!",
        Secret:   "your-secret-key",
    }
    
    resp, err := client.SendMessage(msgReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Message sent: %+v\n", resp)
}
```

## Fitur

- ✅ **Send Message** - Kirim pesan WhatsApp dengan atau tanpa media
- ✅ **Generate OTP** - Generate dan kirim kode OTP 6 digit
- ✅ **Validate OTP** - Validasi kode OTP yang dikirim
- ✅ **Health Check** - Cek status API
- ✅ **Error Handling** - Penanganan error yang komprehensif
- ✅ **Package Type Helpers** - Helper functions untuk cek tipe paket

## API Reference

### Client

#### Membuat Client

```go
// Client dengan base URL default
client := kirimi.NewClient()

// Client dengan custom base URL
client := kirimi.NewClientWithBaseURL("https://custom-api.kirimi.id")

// Set timeout
client.SetTimeout(60 * time.Second)
```

### Send Message

Kirim pesan WhatsApp dengan atau tanpa media.

```go
// Text message
msgReq := kirimi.SendMessageRequest{
    UserCode: "USER123",
    DeviceID: "DEVICE456", 
    Receiver: "628123456789",
    Message:  "Hello World!",
    Secret:   "your-secret-key",
}

// Message dengan media
msgReq := kirimi.SendMessageRequest{
    UserCode: "USER123",
    DeviceID: "DEVICE456",
    Receiver: "628123456789", 
    Message:  "Check this image!",
    Secret:   "your-secret-key",
    MediaURL: "https://example.com/image.jpg",
}

resp, err := client.SendMessage(msgReq)
```

**Batasan:**
- Maksimal 1200 karakter per pesan
- Media URL harus valid
- Fitur media hanya tersedia untuk paket Lite, Basic, dan Pro

### Generate OTP

Generate dan kirim kode OTP 6 digit ke nomor WhatsApp.

```go
otpReq := kirimi.GenerateOTPRequest{
    UserCode: "USER123",
    DeviceID: "DEVICE456",
    Phone:    "628123456789",
    Secret:   "your-secret-key",
}

resp, err := client.GenerateOTP(otpReq)
```

**Requirements:**
- Hanya tersedia untuk paket Basic dan Pro
- Device harus dalam status 'connected'
- Device belum expired
- Kuota tersedia (jika tidak unlimited)

### Validate OTP

Validasi kode OTP yang telah dikirim.

```go
validateReq := kirimi.ValidateOTPRequest{
    UserCode: "USER123",
    DeviceID: "DEVICE456",
    Phone:    "628123456789",
    OTP:      "123456",
    Secret:   "your-secret-key",
}

resp, err := client.ValidateOTP(validateReq)
```

**Notes:**
- OTP berlaku selama 5 menit
- OTP hanya bisa digunakan sekali (one-time use)
- OTP akan otomatis dihapus setelah divalidasi atau expired

### Health Check

Cek status API.

```go
health, err := client.HealthCheck()
```

## Package Types

Kirimi mendukung berbagai tipe paket dengan fitur yang berbeda:

| Package | ID | Fitur |
|---------|----|---------|
| Free | 1 | Text only + watermark |
| Lite | 2, 6, 9 | Text + Media |
| Basic | 3, 7, 10 | Text + Media + OTP |
| Pro | 4, 8, 11 | Text + Media + OTP |

### Helper Functions

```go
// Cek apakah paket mendukung OTP
isOTPSupported := kirimi.IsBasicOrProPackage(packageID)

// Cek apakah paket mendukung media
isMediaSupported := kirimi.IsMediaSupportedPackage(packageID)

// Cek apakah paket adalah free
isFree := kirimi.IsFreePackage(packageID)
```

## Error Handling

SDK menggunakan custom error type `APIError` untuk error dari API:

```go
resp, err := client.SendMessage(msgReq)
if err != nil {
    if apiErr, ok := err.(*kirimi.APIError); ok {
        fmt.Printf("API Error %d: %s\n", apiErr.StatusCode, apiErr.Message)
        fmt.Printf("Response: %+v\n", apiErr.Response)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

### Common Error Codes

- **400**: Parameter tidak lengkap, secret tidak valid, device expired, dll
- **403**: Fitur tidak tersedia untuk tipe paket
- **404**: User atau device tidak ditemukan
- **500**: Server error

## Constants

```go
// Package Types
kirimi.PackageFree     // 1
kirimi.PackageLite1    // 2
kirimi.PackageBasic1   // 3
kirimi.PackagePro1     // 4
kirimi.PackageLite2    // 6
kirimi.PackageBasic2   // 7
kirimi.PackagePro2     // 8
kirimi.PackageLite3    // 9
kirimi.PackageBasic3   // 10
kirimi.PackagePro3     // 11

// Other Constants
kirimi.DefaultBaseURL     // "https://api.kirimi.id"
kirimi.APIVersion         // "v1"
kirimi.MaxMessageLength   // 1200
```

## Contoh Lengkap

Lihat file `example/main.go` untuk contoh penggunaan lengkap semua fitur.

```bash
go run example/main.go
```

## Requirements

- Go 1.21 atau lebih baru
- Akun Kirimi Console dengan secret key yang valid
- Device yang sudah terdaftar dan dalam status 'connected'

## License

MIT License

## Contributing

Kontribusi sangat diterima! Silakan buat issue atau pull request.

## Support

Untuk dukungan teknis, silakan hubungi tim Kirimi atau buat issue di repository ini.# kirimi-go
