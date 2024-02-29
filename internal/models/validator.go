package models

import (
	"net/url"
	"reflect"
	"sync"

	"github.com/ai-readone/go-url-shortner/logger"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var _ binding.StructValidator = &DefaultValidator{}

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

// this validator parses the url reieved,
// and marks it as invalid if the there is no host(example: perizer.com) in the url
var urlValidation validator.Func = func(fl validator.FieldLevel) bool {
	u, err := url.Parse(fl.Field().String())
	if err != nil {
		logger.Error(err)
		return false
	}

	if u.Host == "" {
		logger.Error("invalid URL")
		return false
	}

	return true
}

var shortUrlLenght validator.Func = func(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	if len(url) != 7 {
		logger.Error("invalid url length")
		return false
	}

	return true
}

func (v *DefaultValidator) init() {
	v.once.Do(func() {
		v.validate = validator.New()
		// sets tag-name which will be use in the struct to assign contraints
		// example: `binding:"required"`
		v.validate.SetTagName("binding")

		// adds custom url validation
		v.validate.RegisterValidation("url", urlValidation, false)
		v.validate.RegisterValidation("shortUrlLenght", shortUrlLenght, false)
	})
}

// gets the kind of an object
// checks the kind of an elemnt of a pointer,
// if the object is a pointer,
func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Pointer {
		valueType = value.Elem().Kind()
	}

	return valueType
}

// implement the structValidator interface to validate struct objects
func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.init()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.init()
	return v.validate
}
