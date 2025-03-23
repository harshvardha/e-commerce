package utility

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var TWILIO_ACCOUNT_SID string = os.Getenv("TWILIO_ACCOUNT_SID")

// func GenerateOtp() string {
// 	rand.Seed(time.Now().UnixNano())
// 	return fmt.Sprintf("%06d", rand.Intn(1000000)) // 6-digit otp
// }

// creating a map to store verification sid as key and phonenumber as value
var verificationTokenToPhonenumber map[string]string

func SendOtpViaPhonenumber(to string, client *twilio.RestClient, VERIFY_SERVICE_SID string) (string, error) {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(VERIFY_SERVICE_SID, params)

	if err != nil {
		return "", err
	}

	// adding verification token and phone number to map
	if verificationTokenToPhonenumber == nil {
		verificationTokenToPhonenumber = make(map[string]string)
	}
	verificationTokenToPhonenumber[*resp.Sid] = to
	return *resp.Sid, nil
}

func VerifyOtp(client *twilio.RestClient, VERIFY_SERVICE_SID string, verificationToken string, code string) (bool, string, error) {
	to := verificationTokenToPhonenumber[verificationToken]
	if len(to) == 0 {
		return false, "", fmt.Errorf("invalid verification token")
	}
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(VERIFY_SERVICE_SID, params)

	if err != nil {
		return false, "", err
	} else if *resp.Status == "approved" {
		delete(verificationTokenToPhonenumber, verificationToken)
		return true, to, nil
	} else {
		return false, "", fmt.Errorf("incorrect otp")
	}
}

// func SendOtpViaEmail(to string) error {
// 	fromEmail := "harshvardhansingh458@gmail.com"
// 	password := "nsau diiq azkz vifv"
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := 587
// 	otp := generateOtp()

// 	mail := gomail.NewMessage()
// 	mail.SetHeader("From", fromEmail)
// 	mail.SetHeader("To", to)
// 	mail.SetHeader("Subject", "one time otp")
// 	mail.SetBody("text/plain", "OTP: "+otp)

// 	dialer := gomail.NewDialer(smtpHost, smtpPort, fromEmail, password)
// 	if err := dialer.DialAndSend(mail); err != nil {
// 		return fmt.Errorf("failed to send email: %v", err.Error())
// 	}

// 	return nil
// }
