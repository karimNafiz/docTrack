package utils

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const BytesInBytes int = 1
const KiloBytesInBytes int = BytesInBytes * 1024
const MegaBytesInBytes int = KiloBytesInBytes * 1024
const GigaBytesInBytes int64 = int64(MegaBytesInBytes) * 1024

var jwtSuperSecretKey = []byte("supersecretkey")

func VerifyPassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, err
}

func GenerateTokens(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSuperSecretKey)

}

func VerifyTokens(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// right now just return the secret key
		return jwtSuperSecretKey, nil
	})

	// check the validity of the token or check for errors
	if err != nil || !token.Valid {
		return nil, errors.New("unauthorized: invalid token")
	}

	// we need to extract the claims
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, errors.New("unauthorized: invalid claims")
	}

	return claims, nil

}

// parseByteSize parses strings like "4B", "5KB", "12 MB", "3GB"
// and returns the corresponding byte count as an int.
func ParseByteSize(sizeStr string) (int, error) {
	s := strings.TrimSpace(sizeStr)
	sUp := strings.ToUpper(s)

	const (
		B  = 1
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)

	var (
		unit       string
		multiplier int64
	)

	switch {
	case strings.HasSuffix(sUp, "GB"):
		unit = "GB"
		multiplier = GB
	case strings.HasSuffix(sUp, "MB"):
		unit = "MB"
		multiplier = MB
		unit = "KB"
		multiplier = KB
	case strings.HasSuffix(sUp, "B"):
		unit = "B"
		multiplier = B
	default:
		return 0, fmt.Errorf("unknown unit in %q", sizeStr)
	}

	numPart := strings.TrimSpace(sUp[:len(sUp)-len(unit)])
	if numPart == "" {
		return 0, fmt.Errorf("missing number before %q", unit)
	}
	n, err := strconv.ParseInt(numPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q: %w", numPart, err)
	}

	total := n * multiplier
	if total > int64(^uint(0)>>1) {
		return 0, fmt.Errorf("size %d%s too large", n, unit)
	}
	return int(total), nil
}

// ReadBytes reads from the io.Reader in chunks of size buffSize (e.g. "4KB"),
// appending each to outBuff and returning the full byte slice.
func ReadBytes(stream io.Reader, buffSize string, callBack func([]byte) error) error {
	// 1) Parse the buffer size once
	bufSize, err := ParseByteSize(buffSize)
	if err != nil {
		return err
	}

	// 2) Allocate a single reusable buffer of that exact size
	readBuff := make([]byte, bufSize)

	// 3) Loop until EOF, appending each read to outBuff
	for {
		n, err := stream.Read(readBuff)
		if n > 0 {
			err := callBack(readBuff)
			if err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
	}
	return nil
}

func WriteBytes(stream io.Writer, buffSize string, data []byte) error {
	// 1) Parse the buffer size (e.g. “4KB” → 4096)
	bufSize, err := ParseByteSize(buffSize)
	if err != nil {
		return err
	}

	// 2) Compute how many writes we need (ceiling of dataLen/bufSize)
	dataLen := len(data)
	noIOCalls := (dataLen + (bufSize - 1)) / bufSize

	// 3) Loop with a classic for-loop
	for i := 0; i < noIOCalls; i++ {
		offset := i * bufSize
		limit := offset + bufSize
		if limit > dataLen {
			limit = dataLen
		}
		// Write it—and wrap any error with context
		if _, err := stream.Write(data[offset:limit]); err != nil {
			return fmt.Errorf("write chunk %d failed: %w", i, err)
		}
	}

	return nil
}

func Copy(reader io.Reader, writer io.Writer) (int64, error) { // in the future, have another parameter for buffer size

	// 1. Create buffer
	buffer := make([]byte, 32*(1<<20)) // or default if empty

	// 2. Track total bytes written
	var totalWritten int64 = 0

	// 3. Loop to read and write
	for {
		n, readErr := reader.Read(buffer)

		if n > 0 {
			writeN, writeErr := writer.Write(buffer[:n])
			totalWritten += int64(writeN)

			if writeN != n {
				// Important: handle partial writes
				// If writeN != n, then we have to write the remaining bytes
				// (very rare but CAN happen with some writers like network streams)
				// Optional at first: panic or error out
				return totalWritten, io.ErrShortWrite
			}

			if writeErr != nil {
				return totalWritten, fmt.Errorf("write error: %w", writeErr)
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				break // Clean EOF
			}
			return totalWritten, fmt.Errorf("read error: %w", readErr)
		}
	}

	return totalWritten, nil
}
