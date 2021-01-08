package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ruleFunctions = map[string]func(fieldValue, ruleValue string) (bool, error){
	"len": func(fieldValue, ruleValue string) (bool, error) {
		mustLen, err := strconv.ParseFloat(ruleValue, 64)
		if err != nil {
			return false, ErrorRuleValueIsNotNumber
		}
		return len([]rune(fieldValue)) == int(mustLen), nil
	},
	"min": func(fieldValue, ruleValue string) (bool, error) {
		val, err := strconv.ParseFloat(fieldValue, 64)
		if err != nil {
			return false, ErrorRuleValueIsNotNumber
		}
		mustMin, err := strconv.ParseFloat(ruleValue, 64)
		if err != nil {
			return false, ErrorRuleValueIsNotNumber
		}
		return val >= mustMin, nil
	},
	"max": func(fieldValue, ruleValue string) (bool, error) {
		val, err := strconv.ParseFloat(fieldValue, 64)
		if err != nil {
			return false, ErrorRuleValueIsNotNumber
		}
		mustMax, err := strconv.ParseFloat(ruleValue, 64)
		if err != nil {
			return false, ErrorRuleValueIsNotNumber
		}
		return val <= mustMax, nil
	},
	"regexp": func(fieldValue, ruleValue string) (bool, error) {
		reg, err := regexp.Compile(ruleValue)
		if err != nil {
			return false, fmt.Errorf("rule value is not valid regexp: %w", err)
		}
		return reg.MatchString(fieldValue), nil
	},
	"in": func(fieldValue, ruleValue string) (bool, error) { //nolint:unparam
		for _, s := range strings.Split(ruleValue, ",") {
			if s == fieldValue {
				return true, nil
			}
		}
		return false, nil
	},
}
