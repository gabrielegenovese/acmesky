package prontogram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UserAuthCredentials struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type UserSignUpRequest struct {
	AccountDisplayName string `json:"display_name"`
	Password           string `json:"password"`
}

type AuthenticatedUser struct {
	UserId string `json:"userId"`
	Sid    string `json:"sid"`
}

type SendMessageRequest struct {
	Sender  AuthenticatedUser `json:"sender"`
	Content string            `json:"content"`
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

	if errSignUp != nil {
		fmt.Printf("[Prontogram] ACMESky can't sign up %s\n", errSignUp)
	} else {
		fmt.Println("[Prontogram] ACMESky Signed up")
	}
	_, errSignIn := SignInAs(credentials)

	if errSignIn != nil {
		fmt.Printf("[Prontogram] ACMESky can't sign in %s\n", errSignIn)
	} else {
		fmt.Println("[Prontogram] ACMESky Signed in")
	}

	return errSignIn
}

func SignUpAs(credentials UserAuthCredentials) error {

	var PRONTOGRAM_BASEURL string = os.Getenv("PRONTOGRAM_BASEURL")

	body := &UserSignUpRequest{
		AccountDisplayName: "ACMESky",
		Password:           credentials.Password,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	res, errReq := http.Post(PRONTOGRAM_BASEURL+"/users/"+credentials.UserId, "application/json", payloadBuf)
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

	var PRONTOGRAM_BASEURL string = os.Getenv("PRONTOGRAM_BASEURL")
	var resBody AuthenticatedUser

	body := &credentials

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	res, errReq := http.Post(PRONTOGRAM_BASEURL+"/auth/"+credentials.UserId+"/login", "application/json", payloadBuf)
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
	var PRONTOGRAM_BASEURL string = os.Getenv("PRONTOGRAM_BASEURL")
	var resBody SendMessageResponse

	authUser, errAuth := SignInAs(getACMESkyCredentials())

	if errAuth != nil {
		return resBody, errAuth
	}

	body := &SendMessageRequest{
		Sender:  authUser,
		Content: content,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	res, errSend := http.Post(PRONTOGRAM_BASEURL+"/users/"+destinationUserId+"/messages", "application/json", payloadBuf)

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
