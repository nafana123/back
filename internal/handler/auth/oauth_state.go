package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

func generateOAuthState(secret string) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	nonceEncoded := base64.RawURLEncoding.EncodeToString(nonce)
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	payload := nonceEncoded + "." + ts

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return payload + "." + signature, nil
}

func validateOAuthState(state, secret string) error {
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		return errors.New("invalid state format")
	}

	payload := parts[0] + "." + parts[1]

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedSig := mac.Sum(nil)

	providedSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return err
	}

	if !hmac.Equal(providedSig, expectedSig) {
		return errors.New("invalid state signature")
	}

	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return err
	}

	if time.Since(time.Unix(ts, 0)) > 5*time.Minute {
		return errors.New("state expired")
	}

	return nil
}
