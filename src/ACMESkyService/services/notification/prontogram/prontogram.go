package prontogram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UserAuthCredentials struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type UserSignUpRequest struct {
	AccountDisplayName string              `json:"display_name"`
	Credentials        UserAuthCredentials `json:"credentials"`
}

type AuthenticatedUser struct {
	UserId string `json:"user_id"`
	Sid    string `json:"sid"`
}

type SendMessageRequest struct {
	Sender     AuthenticatedUser `json:"sender"`
	ReceiverId string            `json:"receiver_user_id"`
	Content    string            `json:"content"`
}

type SendMessageResponse struct {
	MessageId string `json:"id"`
}

func getACMESkyCredentials() UserAuthCredentials {
	return UserAuthCredentials{
		UserId:   os.Getenv("PRONTOGRAM_USERNAME"),
		Password: os.Getenv("PRONTOGRAM_PASSWORD"),
	}
}

func Init() error {
	credentials := getACMESkyCredentials()
	errSignUp := SignUpAs(credentials)

	if errSignUp == nil {
		fmt.Printf("[Prontogram] ACMESky Signed up")
	} else {
		fmt.Printf("[Prontogram] ACMESky can't sign up")
	}
	_, errSignIn := SignInAs(credentials)

	if errSignUp == nil {
		fmt.Printf("[Prontogram] ACMESky Signed in")
	} else {
		fmt.Printf("[Prontogram] ACMESky can't sign in")
	}

	return errSignIn
}

func SignUpAs(credentials UserAuthCredentials) error {

	var PRONTOGRAM_SERVICE_ADDRESS string = os.Getenv("PRONTOGRAM_ADDRESS")

	body := &UserSignUpRequest{
		AccountDisplayName: "ACMESky",
		Credentials:        credentials,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	res, errReq := http.Post("http://"+PRONTOGRAM_SERVICE_ADDRESS+"/api/auth/signup", "application/json", payloadBuf)
	if errReq != nil {
		return errReq
	}

	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return fmt.Errorf("HTTP_ERROR:" + res.Status)
	}

	return nil
}

func SignInAs(credentials UserAuthCredentials) (AuthenticatedUser, error) {

	var PRONTOGRAM_SERVICE_ADDRESS string = os.Getenv("PRONTOGRAM_ADDRESS")
	var resBody AuthenticatedUser

	body := &credentials

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	res, errReq := http.Post("http://"+PRONTOGRAM_SERVICE_ADDRESS+"/api/auth/login", "application/json", payloadBuf)
	if errReq != nil {
		return resBody, errReq
	}

	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return resBody, fmt.Errorf("HTTP_ERROR:" + res.Status)
	}

	decodeErr := json.NewDecoder(res.Body).Decode(&resBody)

	if decodeErr != nil {
		return resBody, fmt.Errorf("PARSE_ERROR:" + decodeErr.Error())
	}

	return resBody, nil
}

func SendMessage(content string, destinationUserId string) (SendMessageResponse, error) {
	var PRONTOGRAM_SERVICE_ADDRESS string = os.Getenv("PRONTOGRAM_ADDRESS")
	var resBody SendMessageResponse

	authUser, errAuth := SignInAs(getACMESkyCredentials())

	if errAuth != nil {
		return resBody, errAuth
	}

	body := &SendMessageRequest{
		Sender:     authUser,
		ReceiverId: destinationUserId,
		Content:    content,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	res, errSend := http.Post("http://"+PRONTOGRAM_SERVICE_ADDRESS+"/api/messages", "application/json", payloadBuf)

	if errSend != nil {
		return resBody, errSend
	}

	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return resBody, fmt.Errorf("HTTP_ERROR:" + res.Status)
	}

	decodeErr := json.NewDecoder(res.Body).Decode(&resBody)

	if decodeErr != nil {
		return resBody, fmt.Errorf("PARSE_ERROR:" + decodeErr.Error())
	}

	if errSend != nil {
		return resBody, errSend
	}

	return resBody, nil
}
