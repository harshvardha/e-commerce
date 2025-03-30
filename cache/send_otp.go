package cache

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

// creating a map to store verification sid as key and phonenumber as value for login
type userLoginCacheDetails struct {
	phonenumber string
	resendAt    time.Time
	expiresAt   time.Time
}
type loginCache struct {
	verificationTokenToPhonenumber map[string]userLoginCacheDetails
	lock                           sync.Mutex
}

// creating a map to store verification sid as key and a struct containing phone_number, email and password as value for registration
type userRegistrationCacheDetails struct {
	email       string
	phonenumber string
	password    string
	resendAt    time.Time
	expiresAt   time.Time
}
type registrationCache struct {
	verificationTokenToUserRegistrationCacheDetails map[string]userRegistrationCacheDetails
	lock                                            sync.Mutex
}

// struct to hold the registration details of user which will be used to register user if registration otp is correct
type userRegistrationDetails struct {
	Email       string
	Phonenumber string
	Password    string
}

// variables to hold caches
var userLoginCache *loginCache
var userRegistrationCache *registrationCache

// storeOTP function for registration cache
func (registerCache *registrationCache) setRegistrationOTPCache(otpVerificationToken string, userRegistrationDetails userRegistrationCacheDetails) {
	registerCache.lock.Lock()
	defer registerCache.lock.Unlock()

	registerCache.verificationTokenToUserRegistrationCacheDetails[otpVerificationToken] = userRegistrationDetails
}

// getOTP function for registration cache
func (registerCache *registrationCache) getRegistrationOTPCache(otpVerificationToken string) (userRegistrationCacheDetails, bool) {
	registerCache.lock.Lock()
	defer registerCache.lock.Unlock()

	value, ok := registerCache.verificationTokenToUserRegistrationCacheDetails[otpVerificationToken]
	if !ok || time.Now().After(value.expiresAt) {
		delete(registerCache.verificationTokenToUserRegistrationCacheDetails, otpVerificationToken)
		return value, false
	}

	return value, true
}

// deleteOTP function for registration cache
func (registerCache *registrationCache) deleteRegistrationOTPFromCache(otpVerificationToken string) (bool, error) {
	registerCache.lock.Lock()
	defer registerCache.lock.Unlock()

	if _, ok := registerCache.verificationTokenToUserRegistrationCacheDetails[otpVerificationToken]; ok {
		delete(registerCache.verificationTokenToUserRegistrationCacheDetails, otpVerificationToken)
		return true, nil
	}

	return false, errors.New("otp verification token not found")
}

// storeOTP function for login cache
func (loginCache *loginCache) setLoginOTPCache(otpVerificationToken string, userLoginDetails userLoginCacheDetails) {
	loginCache.lock.Lock()
	defer loginCache.lock.Unlock()

	loginCache.verificationTokenToPhonenumber[otpVerificationToken] = userLoginDetails
}

// getOTP function for login cache
func (loginCache *loginCache) getLoginOTPCache(otpVerificationToken string) (userLoginCacheDetails, bool) {
	loginCache.lock.Lock()
	defer loginCache.lock.Unlock()

	value, ok := loginCache.verificationTokenToPhonenumber[otpVerificationToken]
	if !ok || time.Now().After(value.expiresAt) {
		delete(loginCache.verificationTokenToPhonenumber, otpVerificationToken)
		return value, false
	}

	return value, true
}

// deleteOTP function for login cache
func (loginCache *loginCache) deleteLoginOTPFromCache(otpVerificationToken string) (bool, error) {
	loginCache.lock.Lock()
	defer loginCache.lock.Unlock()

	if _, ok := loginCache.verificationTokenToPhonenumber[otpVerificationToken]; ok {
		delete(loginCache.verificationTokenToPhonenumber, otpVerificationToken)
		return true, nil
	}

	return false, errors.New("otp verification token not found")
}

/*
@params userDetails: to hold the multiple values from either registration request or login request
and if userDetails has 3 values (phonenumber, email, password) then it is a registration request
and if userDetails has 1 value (phonenumber) then it is a login request
*/
func SendOtpViaPhonenumber(client *twilio.RestClient, VERIFY_SERVICE_SID string, userDetails ...string) (string, error) {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(userDetails[0])
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(VERIFY_SERVICE_SID, params)

	if err != nil {
		return "", err
	}

	if len(userDetails) == 3 {
		if userRegistrationCache == nil {
			userRegistrationCache = &registrationCache{
				verificationTokenToUserRegistrationCacheDetails: make(map[string]userRegistrationCacheDetails),
			}
		}
		userRegistrationCache.setRegistrationOTPCache(*resp.Sid, userRegistrationCacheDetails{
			phonenumber: userDetails[0],
			email:       userDetails[1],
			password:    userDetails[2],
			resendAt:    time.Now().Add(time.Minute),
			expiresAt:   time.Now().Add(time.Minute * 5),
		})
	} else if len(userDetails) == 1 {
		if userLoginCache == nil {
			userLoginCache = &loginCache{
				verificationTokenToPhonenumber: make(map[string]userLoginCacheDetails),
			}
		}
		userLoginCache.setLoginOTPCache(*resp.Sid, userLoginCacheDetails{
			phonenumber: userDetails[0],
			resendAt:    time.Now().Add(time.Minute),
			expiresAt:   time.Now().Add(time.Minute * 5),
		})
	} else {
		return "", fmt.Errorf("invalid user details")
	}
	return *resp.Sid, nil
}

func verifyOTP(client *twilio.RestClient, VERIFY_SERVICE_SID string, phonenumber string, code string) (bool, error) {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(phonenumber)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(VERIFY_SERVICE_SID, params)

	if err != nil {
		return false, err
	} else if *resp.Status == "approved" {
		return true, nil
	} else {
		return false, errors.New("incorrect otp")
	}
}

func VerifyRegistrationOTP(client *twilio.RestClient, VERIFY_SERVICE_SID string, verificationToken string, code string) (bool, userRegistrationDetails, error) {
	userDetails, isPresent := userRegistrationCache.getRegistrationOTPCache(verificationToken)
	if !isPresent {
		return false, userRegistrationDetails{
			Email:       "",
			Phonenumber: "",
			Password:    "",
		}, errors.New("otp expired")
	}

	isOTPCorrect, err := verifyOTP(client, VERIFY_SERVICE_SID, userDetails.phonenumber, code)
	if !isOTPCorrect {
		return false, userRegistrationDetails{
			Email:       "",
			Phonenumber: "",
			Password:    "",
		}, err
	}

	isDeleted, err := userRegistrationCache.deleteRegistrationOTPFromCache(verificationToken)
	if !isDeleted {
		log.Println(err)
	}
	return true, userRegistrationDetails{
		Email:       userDetails.email,
		Phonenumber: userDetails.phonenumber,
		Password:    userDetails.password,
	}, nil
}

func VerifyLoginOTP(client *twilio.RestClient, VERIFY_SERVICE_SID string, verificationToken string, code string) (bool, error) {
	userDetails, isPresent := userLoginCache.getLoginOTPCache(verificationToken)
	if !isPresent {
		return false, errors.New("otp expired")
	}

	isOTPCorrect, err := verifyOTP(client, VERIFY_SERVICE_SID, userDetails.phonenumber, code)
	if !isOTPCorrect {
		return false, err
	}

	isDeleted, err := userLoginCache.deleteLoginOTPFromCache(verificationToken)
	if !isDeleted {
		log.Println(err)
	}
	return true, nil
}

func ResendOTP(client *twilio.RestClient, VERIFY_SERVICE_SID string, verificationToken string, requestType int) (string, error) {
	if requestType == 3 {
		// requestType == 3 means registration otp
		userDetails, isPresent := userRegistrationCache.getRegistrationOTPCache(verificationToken)
		if !isPresent {
			if userDetails.phonenumber == "" && userDetails.email == "" && userDetails.password == "" {
				return "", errors.New("please provide your details again")
			} else {
				return SendOtpViaPhonenumber(client, VERIFY_SERVICE_SID, userDetails.phonenumber, userDetails.email, userDetails.password)
			}
		} else {
			// checking if resend is allowed for this otp request or not
			if time.Now().After(userDetails.resendAt) {
				return SendOtpViaPhonenumber(client, VERIFY_SERVICE_SID, userDetails.email, userDetails.phonenumber, userDetails.password)
			}

			return "", fmt.Errorf("resend allowed after %.2f seconds", time.Until(userDetails.resendAt).Seconds())
		}
	} else if requestType == 1 {
		// requestType == 1 means login otp
		userDetails, isPresent := userLoginCache.getLoginOTPCache(verificationToken)
		if !isPresent {
			if userDetails.phonenumber == "" {
				return "", errors.New("please provide your details again")
			} else {
				return SendOtpViaPhonenumber(client, VERIFY_SERVICE_SID, userDetails.phonenumber)
			}
		} else {
			if time.Now().After(userDetails.resendAt) {
				return SendOtpViaPhonenumber(client, VERIFY_SERVICE_SID, userDetails.phonenumber)
			}
			return "", fmt.Errorf("resend allowed after %.2f seconds", time.Until(userDetails.resendAt).Seconds())
		}
	}

	return "", errors.New("invalid request")
}
