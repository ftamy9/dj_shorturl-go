package auth

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"time"
	"crypto/sha256"
)

type tokenPayload struct {
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

func createSign(message []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func CreateToken(userID, authSecret string) (string, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	payload := tokenPayload{
		UserID:    userID,
		Timestamp: now,
	}

	data, err := json.Marshal(map[string]string{
		"user_id":   userID,
		"timestamp": now,
	})
	if err != nil {
		return "", err
	}

	payload.Sign = createSign(data, authSecret)

	tokenData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(tokenData), nil
}

func VerifyToken(tokenBase64, authSecret string, ttlMinutes int) (bool, string, error) {
	tokenBytes, err := base64.StdEncoding.DecodeString(tokenBase64)
	if err != nil {
		return false, "", err
	}

	var payload tokenPayload
	if err := json.Unmarshal(tokenBytes, &payload); err != nil {
		return false, "", err
	}

	timestamp, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		return false, "", err
	}

	if time.Since(timestamp) > time.Duration(ttlMinutes)*time.Minute {
		return false, "timeout", nil
	}

	data, err := json.Marshal(map[string]string{
		"user_id":   payload.UserID,
		"timestamp": payload.Timestamp,
	})
	if err != nil {
		return false, "", err
	}

	expectedSign := createSign(data, authSecret)
	if !hmac.Equal([]byte(payload.Sign), []byte(expectedSign)) {
		return false, "sign verify failed", nil
	}

	return true, payload.UserID, nil
}
