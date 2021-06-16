package checksum

// This file is based on https://github.com/sassoftware/relic/blob/master/lib/authenticode/checksum.go
// Unfortunately, unused because go does NOT calculate and write PE checksums.

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash"
	"io"
	"os"
)

const dosHeaderSize = 64

func readDosHeader(r io.Reader, d io.Writer) (int64, error) {
	dosheader, err := readAndHash(r, d, dosHeaderSize)
	if err != nil {
		return 0, err
	} else if dosheader[0] != 'M' || dosheader[1] != 'Z' {
		return 0, errors.New("not a PE file")
	}
	return int64(binary.LittleEndian.Uint32(dosheader[0x3c:])), nil
}

// read bytes from a stream and return a byte slice, also feeding a hash
func readAndHash(r io.Reader, d io.Writer, n int) ([]byte, error) {
	if n == 0 {
		return nil, nil
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	if d != nil {
		_, _ = d.Write(buf)
	}
	return buf, nil
}

// ValidatePEChecksum reads and a checksum and compares it with the calculated one
func ValidatePEChecksum(f *os.File) (bool, error) {
	if _, err := f.Seek(0, 0); err != nil {
		return false, err
	}
	peStart, err := readDosHeader(f, nil)
	if err != nil {
		return false, err
	}
	ck := NewPEChecksum(int(peStart))
	if _, err := io.Copy(ck, f); err != nil {
		return false, err
	}

	peChecksum := make([]byte, 4)

	if _, err := f.Seek(0, 0); err != nil {
		return false, err
	}
	if _, err = f.ReadAt(peChecksum, peStart+88); err != nil {
		return false, err
	}
	newChecksum := ck.Sum(nil)

	result := bytes.Compare(newChecksum, peChecksum)

	return result == 0, nil
}

type peChecksum struct {
	cksumPos  int
	sum, size uint32
	odd       bool
}

// NewPEChecksum Hasher that calculates the undocumented, non-CRC checksum used in PE images.
// peStart is the offset found at 0x3c in the DOS header.
func NewPEChecksum(peStart int) hash.Hash {
	var cksumPos int
	if peStart <= 0 {
		cksumPos = -1
	} else {
		cksumPos = peStart + 88
	}
	return &peChecksum{cksumPos: cksumPos}
}

func (peChecksum) Size() int {
	return 4
}

func (peChecksum) BlockSize() int {
	return 2
}

func (h *peChecksum) Reset() {
	h.cksumPos = -1
	h.sum = 0
	h.size = 0
}

func (h *peChecksum) Write(d []byte) (int, error) {
	// tolerate odd-sized files by adding a final zero byte, but odd writes anywhere but the end are an error
	n := len(d)
	if h.odd {
		return 0, errors.New("odd write")
	} else if n%2 != 0 {
		h.odd = true
		d2 := make([]byte, n+1)
		copy(d2, d)
		d = d2
	}
	ckpos := -1
	if h.cksumPos > n {
		h.cksumPos -= n
	} else if h.cksumPos >= 0 {
		ckpos = h.cksumPos
		h.cksumPos = -1
	}
	sum := h.sum
	for i := 0; i < n; i += 2 {
		val := uint32(d[i+1])<<8 | uint32(d[i])
		if i == ckpos || i == ckpos+2 {
			val = 0
		}
		sum += val
		sum = 0xffff & (sum + (sum >> 16))
	}
	h.sum = sum
	h.size += uint32(n)
	return n, nil
}

func (h *peChecksum) Sum(buf []byte) []byte {
	sum := h.sum
	sum = 0xffff & (sum + (sum >> 16))
	sum += h.size
	d := make([]byte, 4)
	binary.LittleEndian.PutUint32(d, sum)
	return append(buf, d...)
}
