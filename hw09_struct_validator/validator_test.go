package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	Coords struct {
		Lat float64 `validate:"min:44|max:45"`
		Lon float64 `validate:"min:38|max:39"`
	}
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int            `validate:"min:18|max:50"`
		Email  string         `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole       `validate:"in:admin,stuff"`
		Phones []string       `validate:"len:11"`
		Data   map[string]int `validate:"in:50,60,70"`
		Coords Coords         `validate:"nested"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	positiveTests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:    "d8f4590320e1343a915lb69410650a8f359d",
				Name:  "Positive",
				Age:   50,
				Email: "test@test.ru",
				Role:  "admin",
				Phones: []string{
					"89999999999",
				},
				Data: map[string]int{
					"test_1": 50,
					"test_2": 60,
					"test_3": 70,
				},
				Coords: Coords{
					Lat: 44.5,
					Lon: 38.89,
				},
			},
			nil,
		},
		{
			App{
				Version: "12.05",
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "",
			},
			nil,
		},
	}

	for i, tt := range positiveTests {
		t.Run(fmt.Sprintf("positive case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}

	negativeTests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:    "d8f4590320e1343a915b69410650a8f359d",
				Name:  "Negative",
				Age:   51,
				Email: "@test.ru",
				Role:  "guest",
				Phones: []string{
					"+79999999999",
				},
				Data: map[string]int{
					"test_1": 51,
					"test_2": 59,
					"test_3": 69,
				},
				Coords: Coords{
					Lat: 45.5,
					Lon: 37.89,
				},
			},
			ValidationErrors{
				{
					Field: "ID",
					Err:   ErrorFieldValueIsInvalid,
				},
				{
					Field: "Age",
					Err:   ErrorFieldValueIsInvalid,
				},
				{
					Field: "Email",
					Err:   ErrorFieldValueIsInvalid,
				},
				{
					Field: "Role",
					Err:   ErrorFieldValueIsInvalid,
				},
				{
					Field: "Phones",
					Err:   ErrorFieldValueIsInvalid,
				},
				{
					Field: "Data",
					Err:   ErrorFieldValueIsInvalid,
				},
			},
		},
		{
			App{
				Version: "12.051",
			},
			ValidationErrors{
				{
					Field: "App",
					Err:   ErrorFieldValueIsInvalid,
				},
			},
		},
		{
			Response{
				Code: 301,
				Body: "",
			},
			ValidationErrors{
				{
					Field: "Code",
					Err:   ErrorFieldValueIsInvalid,
				},
			},
		},
	}

	for i, tt := range negativeTests {
		t.Run(fmt.Sprintf("negative case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Error(t, err)
			for _, e := range err.(ValidationErrors) {
				for _, expected := range tt.expectedErr.(ValidationErrors) {
					if e.Field == expected.Field && !errors.Is(e.Err, expected.Err) {
						require.Failf(t, "Wrong error", "Expected: %s, Actual: %s", expected.Err, e.Err)
					}
				}
			}
		})
	}

	t.Run("no struct", func(t *testing.T) {
		err := Validate(5)
		require.Error(t, err)
		if !errors.Is(err, ErrorInterfaceIsNotStructure) {
			require.Failf(t, "Wrong error", "Expected: %s, Actual: %s", ErrorInterfaceIsNotStructure, err)
		}
	})

}
