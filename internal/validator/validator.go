package validator

import (
	"context"
	"time"
	"userrestapigorm/internal/repository"

	"github.com/go-playground/validator/v10"
)

/*
* New initializes the validator and registers custom validation functions.
* It sets up the validator to be used for validating request payloads and other data structures.
 */
func New() *validator.Validate {
	validate := validator.New()
	// Register custom validation functions
	validate.RegisterValidation("is_unique_email", IsUserEmailExists)
	validate.RegisterValidation("is_birthdate_in_past", IsBirthdateInPast)
	return validate
}

/*
* IsUserEmailExists is a custom validation function that checks if the provided email is unique in the database.
* It queries the database to determine if the email already exists.
@param fl validator.FieldLevel - The field level information provided by the validator.
@return bool - Returns true if the email is unique (does not exist in the database), otherwise returns false.
*/
func IsUserEmailExists(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	ctx := context.Background()
	exists, err := repository.UserEmailExists(ctx, email)
	if err != nil {
		return false
	}

	return !exists
}

func IsBirthdateInPast(fl validator.FieldLevel) bool {
	birthdate, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	return birthdate.Before(today)
}
