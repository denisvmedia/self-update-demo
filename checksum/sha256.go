package checksum

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

// ValidateSha256Checksum calculates sha256 of the given file and compares it to the expected one.
// Note, on Windows you can calculate sha256 of the existing file using: powershell "(get-filehash app.exe).Hash"
func ValidateSha256Checksum(f *os.File, expected string) (bool, error) {
	if _, err := f.Seek(0, 0); err != nil {
		return false, err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}

	ck := fmt.Sprintf("%x", h.Sum(nil))

	return strings.Compare(ck, strings.ToLower(expected)) == 0, nil
}
