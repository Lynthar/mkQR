package encoder

// Encoder interface for encoding data into QR-compatible strings
type Encoder interface {
	Encode() string
}

// ContentType represents the type of QR code content
type ContentType string

const (
	TypeText    ContentType = "text"
	TypeURL     ContentType = "url"
	TypeWiFi    ContentType = "wifi"
	TypeVCard   ContentType = "vcard"
	TypeEmail   ContentType = "email"
	TypePhone   ContentType = "phone"
	TypeSMS     ContentType = "sms"
	TypeOTP     ContentType = "otp"
	TypeGeo     ContentType = "geo"
	TypeEvent   ContentType = "event"
	TypeProxy   ContentType = "proxy"
	TypeUnknown ContentType = "unknown"
)
