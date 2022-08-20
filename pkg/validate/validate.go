package validate

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

// emailRegex is the regular expression used to determine if a string is an email.
// https://github.com/go-playground/validator/blob/v10.10.0/regexes.go#L73
var emailRegex *regexp.Regexp

func Init() {

	// Instantiate a validator.
	validate = validator.New()

	// Create a translator for english so the error messages are
	// more human-readable than technical.
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the english error messages for use.
	en_translations.RegisterDefaultTranslations(validate, translator)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = validate.RegisterTranslation("mdate", translator, func(ut ut.Translator) error {
		return ut.Add("mdate", "{0} must be valid date format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mdate", fe.Field())
		return t
	})

	_ = validate.RegisterValidation("mdate", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		layout := "2006-01-02 15:04:05"
		_, err := time.Parse(layout, str)
		return err == nil
	})

	// emailRegexString is the regular expression string used to compile into a regexp.
	// https://github.com/go-playground/validator/blob/v10.10.0/regexes.go#L18
	const emailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	emailRegex = regexp.MustCompile(emailRegexString)
}

// =======================================================================================

// SliceStruct menerima payload []struct dan mengembalikan validasi error yang
// didapatkan dalam bentuk string dan error aslinya. jika valid akan mengembalikan nil
func SliceStruct(input interface{}) (string, error) {

	if input == nil || reflect.TypeOf(input).Kind() != reflect.Slice {
		errorMsg := "developer_vault. input cant be nil or other than slice struct"
		return errorMsg, errors.New(errorMsg)
	}

	err := validate.Var(input, "omitempty,dive")
	if err != nil {

		// not accepted error by validator
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err.Error(), err
		}

		// iterate error message
		var sb strings.Builder
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			sb.WriteString(e.Translate(translator) + ", ")
		}
		return strings.TrimRight(sb.String(), ", "), errs
	}

	return "", nil
}

// Struct menerima payload struct dan mengembalikan validasi error yang
// didapatkan dalam bentuk string dan error aslinya. jika valid akan mengembalikan nil
func Struct(input interface{}) (string, error) {

	if input == nil {
		errorMsg := "developer_vault. input cant be nil"
		return errorMsg, errors.New(errorMsg)
	}

	err := validate.Struct(input)
	if err != nil {
		// not accepted error by validator
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err.Error(), err
		}

		// iterate error message
		var sb strings.Builder
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			sb.WriteString(e.Translate(translator) + ", ")
		}
		return strings.TrimRight(sb.String(), ", "), errs
	}

	return "", nil
}

// Var menerima payload variable dan mengembalikan validasi error yang
// didapatkan dalam bentuk string dan error aslinya. jika valid akan mengembalikan nil
func Var(input interface{}, tag string) (string, error) {

	err := validate.Var(input, tag)
	if err != nil {
		// not accepted error by validator
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err.Error(), err
		}

		// iterate error message
		var sb strings.Builder
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			sb.WriteString(e.Translate(translator) + ", ")
		}
		return strings.TrimRight(sb.String(), ", "), errs
	}

	return "", nil
}

func CheckEmail(email string) bool {
	return emailRegex.MatchString(email)
}
