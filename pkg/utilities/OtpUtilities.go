package utilities

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
)

type Digits int
const (
	OtpInterval = 30 // Interval in seconds
	SixDigits   Digits = 6  // Number of digits in the OTP code
	EightDigits Digits = 8
)
func (d Digits) Length() int {
	return int(d)
}

type otpUtilities struct {
	key []byte
}

func NewOtpUtilities(key []byte) OTPUtilities {
	return &otpUtilities{
		key: key,
	}
}

type OTPUtilities interface {
	TOTPToken(counter int64, digits Digits) (int, error)
}


func (o *otpUtilities) TOTPToken(counter int64, digits Digits) (int, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))
	h := hmac.New(sha1.New, o.key)
	_, err := h.Write(buf)
	if err != nil {
		return 0, err
	}
	hash := h.Sum(nil)
	offset := hash[len(hash)-1] & 0xf
	code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff
	l := digits.Length()
	return int(code % uint32(pow10(l))), nil
}

func pow10(n int) int {
	if n == 0 {
		return 1
	}
	return 10 * pow10(n-1)
}