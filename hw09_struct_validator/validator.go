package hw09_struct_validator //nolint:golint,stylecheck
import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrorInterfaceIsNotStructure = errors.New("this interface is not a structure")
	ErrorFieldValueIsInvalid     = errors.New("this field value is invalid")
	ErrorRuleValueIsNotNumber    = errors.New("this rule value is not a number")
	ErrorRuleUnknownName         = errors.New("unknown validator rule")
	ErrorRuleInvalid             = errors.New("this rule is invalid")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	for i, e := range v {
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(") Field: ")
		b.WriteString(e.Field)
		b.WriteString(", Error: ")
		b.WriteString(e.Err.Error())
		if i < len(v)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func Validate(v interface{}) error {
	var vErr ValidationErrors

	vOf := reflect.ValueOf(v)

	if vOf.Kind() != reflect.Struct {
		return ErrorInterfaceIsNotStructure
	}

	for i := 0; i < vOf.NumField(); i++ {
		field := vOf.Type().Field(i)
		tag := field.Tag.Get("validate")
		if len(tag) == 0 {
			continue
		}
		if tag == "nested" {
			if err := Validate(vOf.Field(i).Interface()); err != nil {
				if validationErrors, ok := err.(ValidationErrors); ok { //nolint:errorlint
					vErr = append(vErr, validationErrors...)
				} else {
					vErr = append(vErr, ValidationError{
						Field: field.Name,
						Err:   err,
					})
				}
			}
			continue
		}
		var values []interface{}
		switch vOf.Field(i).Kind() {
		case reflect.Slice, reflect.Array:
			for j := 0; j < vOf.Field(i).Len(); j++ {
				values = append(values, vOf.Field(i).Index(j).Interface())
			}
		case reflect.Map:
			for _, key := range vOf.Field(i).MapKeys() {
				values = append(values, vOf.Field(i).MapIndex(key).Interface())
			}
		default:
			values = append(values, vOf.Field(i).Interface())
		}
		for _, value := range values {
			vErr = append(vErr, validate(value, field.Name, tag)...)
		}
	}
	if len(vErr) == 0 {
		return nil
	}
	return vErr
}

func validate(value interface{}, fieldName, rules string) ValidationErrors {
	var vErr ValidationErrors
	for _, s := range strings.Split(rules, "|") {
		a := strings.SplitN(s, ":", 2)
		if len(a) != 2 {
			vErr = append(vErr, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("%w, %s", ErrorRuleInvalid, s),
			})
			continue
		}
		ruleName := a[0]
		ruleValue := a[1]
		ruleFunction, ok := ruleFunctions[ruleName]
		if !ok {
			vErr = append(vErr, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("%w, %s", ErrorRuleUnknownName, ruleName),
			})
			continue
		}
		ok, err := ruleFunction(toString(value), ruleValue)
		if err != nil {
			vErr = append(vErr, ValidationError{
				Field: fieldName,
				Err:   err,
			})
			continue
		}
		if !ok {
			vErr = append(vErr, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("%w, %s, value:%s", ErrorFieldValueIsInvalid, s, toString(value)),
			})
		}
	}
	if len(vErr) == 0 {
		return nil
	}
	return vErr
}

func toString(i interface{}) string {
	format := "%s"
	switch i.(type) {
	case float32, float64:
		format = "%f"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		format = "%d"
	}
	return fmt.Sprintf(format, i)
}
