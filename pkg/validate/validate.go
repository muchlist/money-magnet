package validate

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	// Real mengembalian validator instance asli
	Real() *validator.Validate

	// SliceStruct menerima payload []struct dan mengembalikan validasi error yang
	// didapatkan dalam bentuk map string, dan error aslinya. inputan valid apabila err nil
	SliceStruct(input interface{}) (map[string]string, error)

	// SliceStruct menerima payload []struct dan mengembalikan validasi error yang
	// didapatkan dalam bentuk map string dan error aslinya. inputan valid apabila err nil
	Struct(input interface{}) (map[string]string, error)

	// Var menerima payload variable dan mengembalikan validasi error yang
	// didapatkan dalam bentuk map string dan error aslinya. inputan valid apabila err nil
	// eg.
	// var i int
	// validate.Var(i, "gt=1,lt=10")
	Var(input interface{}, tag string) (map[string]string, error)
}

type mValidator struct {
	instance   *validator.Validate
	translator ut.Translator
}

type envelop map[string]string

func (e envelop) makeError() error {
	var sb strings.Builder
	for _, message := range e {
		sb.WriteString(message + ", ")
	}
	return errors.New(strings.TrimRight(sb.String(), ", "))
}

// =======================================================================================
// register validator baru pada bagian code New()  <<<<<<

// New menginisiasi validator dan mengembalikan custom validator
// dengan format error map
func New(regs ...Register) *mValidator {

	// init instance of 'validate' with sane defaults
	// init default translator
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Panic("translator not found")
	}

	// register default translation
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	// register tag e.Field() use json tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// init register validator
	// init register do this to you
	//  _ = validate.RegisterTranslation("custom_date", trans, func(ut ut.Translator) error {
	// 	return ut.Add("custom_date", "{0} must be valid date format", true)
	// }, func(ut ut.Translator, fe validator.FieldError) string {
	// 	t, _ := ut.T("custom_date", fe.Field())
	// 	return t
	// })

	// _ = validate.RegisterValidation("custom_date", func(fl validator.FieldLevel) bool {
	// 	str := fl.Field().String()
	// 	layout := "2006-01-02 15:04:05"
	// 	_, err := time.Parse(layout, str)
	// 	return err == nil
	// })

	for _, v := range regs {
		_ = validate.RegisterTranslation(v.Key, trans, func(ut ut.Translator) error {
			return ut.Add(v.Key, v.Translate, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(v.Key, fe.Field())
			return t
		})
		_ = validate.RegisterValidation(v.Key, v.ValidFunc)
	}

	return &mValidator{
		instance:   validate,
		translator: trans,
	}
}

type Register struct {
	// Key example : "custom_date"
	Key string
	// Translate example: "{0} must be valid date format"
	Translate string
	// ValidFunc like you add validator to golang validator/v10
	ValidFunc func(fl validator.FieldLevel) bool
}

// =======================================================================================

// Real mengembalian validator instance asli
func (m *mValidator) Real() *validator.Validate {
	return m.instance
}

func (m *mValidator) Engine() interface{} {
	return m.instance
}

// SliceStruct menerima payload []struct dan mengembalikan validasi error yang
// didapatkan dalam bentuk map string, dan error aslinya jika valid akan mengembalikan nil
func (m *mValidator) SliceStruct(input interface{}) (map[string]string, error) {

	errMap := envelop{}

	isInputSlice := reflect.TypeOf(input).Kind() == reflect.Slice
	if input == nil || !isInputSlice {
		errorMsg := "500: input cant be nil or other than slice struct"
		errMap["message"] = errorMsg
		return errMap, errors.New(errorMsg)
	}

	err := m.instance.Var(input, "omitempty,dive")
	if err != nil {

		// not accepted error by validator
		if _, ok := err.(*validator.InvalidValidationError); ok {
			errMap["message"] = err.Error()
			return errMap, err
		}

		// iterate error message
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			errMap[e.Field()] = e.Translate(m.translator)
		}
		return errMap, errMap.makeError()
	}

	return nil, nil
}

// SliceStruct menerima payload []struct dan mengembalikan validasi error yang
// didapatkan dalam bentuk map string dan error aslinya jika valid akan mengembalikan nil
func (m *mValidator) Struct(input interface{}) (map[string]string, error) {

	errMap := envelop{}

	if input == nil {
		errorMsg := "500: input cant be nil"
		errMap["message"] = errorMsg
		return errMap, errors.New(errorMsg)
	}

	err := m.instance.Struct(input)
	if err != nil {
		// not accepted error by validator
		if _, ok := err.(*validator.InvalidValidationError); ok {
			errMap["message"] = err.Error()
			return errMap, err
		}

		// iterate error message
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			errMap[e.Field()] = e.Translate(m.translator)
		}
		return errMap, errMap.makeError()
	}

	return nil, nil
}

// Var menerima payload variable dan mengembalikan validasi error yang
// didapatkan dalam bentuk map string dan error aslinya jika valid akan mengembalikan nil
func (m *mValidator) Var(input interface{}, tag string) (map[string]string, error) {

	errMap := envelop{}

	err := m.instance.Var(input, tag)
	if err != nil {
		// not accepted error by validator
		if _, ok := err.(*validator.InvalidValidationError); ok {
			errMap["message"] = err.Error()
			return errMap, err
		}

		// iterate error message
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			errMap[e.Field()] = e.Translate(m.translator)
		}
		return errMap, errMap.makeError()
	}

	return nil, nil
}
