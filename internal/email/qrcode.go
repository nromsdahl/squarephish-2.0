package email

import (
	"bytes"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a QR code for the given data string and returns
// it as a PNG image in bytes.
// Parameters:
//   - data:  The string data to encode (e.g., a URL).
//   - size:  The desired width and height of the image in pixels.
//
// It returns qr code byte sequence.
func GenerateQRCode(data string, size int) ([]byte, error) {
	var qr []byte

	qr, err := qrcode.Encode(data, qrcode.Medium, size)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

// GenerateQRCodeASCII generates a QR code for the given data string and returns
// it as a string of ASCII characters formatted for HTML.
// This is a Golang implementation of https://github.com/nromsdahl/squarephish2/pull/2.
// Parameters:
//   - data:  The string data to encode (e.g., a URL).
//
// It returns qr code ascii string.
func GenerateQRCodeASCII(data string) (string, error) {
	var qr *qrcode.QRCode
	var ascii_qrcode string

	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return ascii_qrcode, err
	}

	bits := qr.Bitmap()
	var buf bytes.Buffer

	inverseColor := false

	// Based on: https://github.com/skip2/go-qrcode/blob/master/qrcode.go#L577
	buf.WriteString(`<pre style="line-height:1">`)
	for y := 0; y < len(bits)-1; y += 2 {
		for x := range bits[y] {
			if bits[y][x] == bits[y+1][x] {
				if bits[y][x] != inverseColor {
					buf.WriteString("█")
				} else {
					buf.WriteString(" ")
				}
			} else {
				if bits[y][x] != inverseColor {
					buf.WriteString("▀")
				} else {
					buf.WriteString("▄")
				}
			}
		}
		buf.WriteString("\n")
	}
	buf.WriteString("</pre>")

	ascii_qrcode = buf.String()

	return ascii_qrcode, nil
}
