package service

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const jwtKey = "!ChangeThisMercureHubJWTSecretKey!"

var jwtToken string

func CreateJwt() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mercure": map[string]any{
			"publish":   []string{"*"},
			"subscribe": []string{"*"},
		},
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		fmt.Println("Error creating the token")
	}
	jwtToken = tokenString
	logrus.Info("JWT Token: ", tokenString)
	return tokenString
}

func GetJwt() string {
	return jwtToken
}

func PushToMercureHub(topic, data string) error {
	form := url.Values{}
	form.Add("topic", topic)
	form.Add("data", data)

	req, err := http.NewRequest(
		"POST",
		"http://localhost:3000/.well-known/mercure",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+GetJwt())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to push data to Mercure Hub, status code: %d", resp.StatusCode)
	}

	return nil
}
