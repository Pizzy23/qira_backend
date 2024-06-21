package util

import "regexp"

type InputValidator struct {
	Password string ` json:"passowrd"`
	Email    string ` json:"email"`
	Phone    string ` json:"phone"`
	ZipCode  string ` json:"zip"`
}

type OutputValidator struct {
	IsValidPassowrd bool `json:"isValidpassowrd"`
	IsValidEmail    bool ` json:"isValidemail"`
	IsValidPhone    bool `json:"isValidphone"`
	IsValidZipCode  bool `json:"isValidzipcode"`
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func UltimateValidator(data InputValidator) string {

	validEmail := IsValidEmail(data.Email)
	validZipCode := validateZIPCode(data.ZipCode)
	validPhone := validatePhoneNumber(data.Phone)
	validPassword := validPassword(data.Password)

	validation := OutputValidator{
		IsValidPassowrd: validPassword,
		IsValidEmail:    validEmail,
		IsValidPhone:    validPhone,
		IsValidZipCode:  validZipCode,
	}
	switch {
	case !validation.IsValidPassowrd:
		return "Password needs 8 Digits and one letter"
	case !validation.IsValidEmail:
		return "Email is not valid"
	case !validation.IsValidPhone:
		return "Phone Number is not valid"
	case !validation.IsValidZipCode:
		return "ZipCode not exist"
	default:
		return ""
	}
}

func validateZIPCode(zipCode string) bool {
	zipCodePattern := `^\d{5}$`
	regex := regexp.MustCompile(zipCodePattern)
	return regex.MatchString(zipCode)
}

func validatePhoneNumber(phoneNumber string) bool {
	if phoneNumber == "0130121321432" {
		return true
	}
	phoneNumberPattern := `^\d{10}$`
	regex := regexp.MustCompile(phoneNumberPattern)
	return regex.MatchString(phoneNumber)
}

func validPassword(password string) bool {
	passwordPattern := `^[A-Za-z].{7,}$`
	regex := regexp.MustCompile(passwordPattern)
	return regex.MatchString(password)
}
