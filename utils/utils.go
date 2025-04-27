package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
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

// the function takes a io.Reader
// the buffer size (buffSize) everytime it makes read call it will at max read data buffSize bytes of data
// the outSize is the amount of data you want to read
// the reason im taking a pointer to a slice is becuasae after appending if the
// underlying data is shifted the pointer to the start of the contiguous memory won't be the same
// but the byte slice passed to this function frm outside this function will still point to the older start
func ReadBytes(stream io.Reader, buffSize string, outBuff []byte) ([]byte, error) {

	readByteSize, err := parseByteSize(buffSize)
	if err != nil {
		return nil, err
	}
	var readBuff []byte
	switch v := readByteSize.(type) {
	case int:
		// v is already int
		readBuff = make([]byte, v)

	case int64:
		// v is already int64
		// but make([]byte, int(v)) if you really must convert
		readBuff = make([]byte, v)
	}
	for {
		n, err := stream.Read(readBuff)
		if n > 0 {
			outBuff = append(outBuff, readBuff[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("read error: %v ", err)
			return nil, err
		}

	}
	return outBuff, nil
}

// parseByteSize parses strings like "4B", "5kb", "12 MB", "3GB"
// and returns either an int (if it fits) or an int64.
func parseByteSize(sizeStr string) (interface{}, error) {
	s := strings.TrimSpace(sizeStr)
	sUp := strings.ToUpper(s)

	// Map units to their byte-size constants (mixed int / int64).
	unitMap := map[string]interface{}{
		"B":  BytesInBytes,
		"KB": KiloBytesInBytes,
		"MB": MegaBytesInBytes,
		"GB": GigaBytesInBytes,
	}

	// Try longer suffixes first so "KB" matches before "B".
	// todo make a list of all the keys
	// and not hard code this
	for _, u := range []string{"GB", "MB", "KB", "B"} {
		if strings.HasSuffix(sUp, u) {
			numPart := strings.TrimSpace(sUp[:len(sUp)-len(u)])
			if numPart == "" {
				return nil, fmt.Errorf("missing number before %q", u)
			}
			n, err := strconv.ParseInt(numPart, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number %q: %w", numPart, err)
			}

			// Decide whether to use int or int64.
			maxInt := int64(int(^uint(0) >> 1))
			switch v := unitMap[u].(type) {
			case int:
				total := n * int64(v)
				if total <= maxInt {
					return int(total), nil
				}
				return total, nil // too big for int, return int64
			case int64:
				total := v * n
				if total <= maxInt {
					return int(total), nil
				}
				return total, nil
			default:
				return nil, fmt.Errorf("unsupported unit type %T", v)
			}
		}
	}

	return nil, fmt.Errorf("unknown unit in %q", sizeStr)
}
