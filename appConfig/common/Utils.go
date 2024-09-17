package common

import (

	"regexp"
	"time"

	"github.com/google/uuid"
)

func FormatDate(d time.Time) string {
	return d.Format("2006-01-02 15:04:05.999999 -0700 MST")
}

func GenerateGUID() string {
	guid := uuid.New().String()
	return guid
}

func IsValidEmail(email string) bool {
    // // Regular expression for basic email validation
    emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

    // // Compile the regex
    regex := regexp.MustCompile(emailRegex)

    // // Check if the email matches the regex pattern
    return regex.MatchString(email)

}