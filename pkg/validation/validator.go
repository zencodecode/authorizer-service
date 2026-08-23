package validation

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func InitValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("IsDateFormat", IsDate)
	v.RegisterValidation("IsDateTimeZFormat", IsDateTimeZ)
	v.RegisterValidation("GtDateAddDay", GtDate)
	v.RegisterValidation("IsJson", IsJson)
	v.RegisterValidation("IsConstantCase", IsConstantCase)

	return v
}

func (cv *CustomValidator) Validate(c *gin.Context, i interface{}) error {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	c.Set("structType", t)
	return cv.Validator.Struct(i)
}

func IsJson(fl validator.FieldLevel) bool {
	s := fl.Field().Interface()
	json, err := json.Marshal(s)
	if err != nil {
		return false
	}
	return string(json) != "{}"
}

func IsDate(fl validator.FieldLevel) bool {
	if _, err := time.Parse("2006-01-02", fl.Field().String()); err != nil {
		return false
	}

	return true
}

func IsDateTimeZ(fl validator.FieldLevel) bool {
	if _, err := time.Parse("2006-01-02 15:04:05 -07:00", fl.Field().String()); err != nil {
		return false
	}

	return true
}

func IsConstantCase(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}

	re := regexp.MustCompile(`^[A-Z0-9]+(_[A-Z0-9]+)*$`)
	return re.MatchString(fl.Field().String())
}

func GtDate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := strings.Split(fl.Param(), `:`)
	paramField := param[0]
	paramValue := param[1]

	if paramField == `` {
		return true
	}

	var paramFieldValue reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
	} else {
		paramFieldValue = fl.Parent().FieldByName(paramField)
	}

	date1, err := time.Parse("2006-01-02", value)
	if err != nil {
		return false
	}

	date2, err := time.Parse("2006-01-02", paramFieldValue.String())
	if err != nil {
		return false
	}

	addDate, err := strconv.Atoi(paramValue)
	if err != nil {
		addDate = 0
	}
	date2 = date2.AddDate(0, 0, addDate)

	return date1.After(date2)
}

func ParseFields(field string) (string, string) {
	fields := strings.Split(field, `=`)
	if len(fields) == 1 {
		return fields[0], ``
	}

	return fields[0], fields[1]
}
