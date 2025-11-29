package centrifugo

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateCentrifugoToken(userID string, userName string, centrifugoHMACSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		// Standard Claims
		"sub": userID,                                  // Subject (User ID)
		"exp": time.Now().Add(time.Minute * 30).Unix(), // Expiration time (e.g., 30 mins)
		"iat": time.Now().Unix(),                       // Issued at

		// Custom Centrifugo "info" claim
		// This data is visible to other users in presence/join messages
		"info": map[string]string{
			"name": userName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign using the Centrifugo-specific secret
	return token.SignedString(centrifugoHMACSecret)
}
