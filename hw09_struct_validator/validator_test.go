package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte `validate:"len:5"`
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	testsSuccess := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "LNL-b5un6PiLrgv",
				Name:   "John Smith",
				Age:    30,
				Email:  "johnsmith@yahoo.com",
				Role:   "admin",
				Phones: []string{"+347880011"},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "6.2.1",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "Success message",
			},
			expectedErr: nil,
		},

		{
			in: Token{
				Header:    []byte("Content-Length"),
				Payload:   []byte("12345"),
				Signature: []byte("SHA-256"),
			},
			expectedErr: nil,
		},
	}

	testsFail := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "12345678",
				Name:   "John",
				Age:    30,
				Email:  "john.@example.com",
				Role:   "admin",
				Phones: []string{"1234567890"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("validation error: field 'Email'"),
				},
			},
		},
		{
			in: Token{
				Header:    []byte("Content-Size"),
				Payload:   []byte("123456"),
				Signature: []byte("md5"),
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Payload",
					Err:   fmt.Errorf("validation error: field 'Payload'"),
				},
			},
		},
		{
			in: User{
				ID:     "98283q3487987329847",
				Name:   "John Deer",
				Age:    15,
				Email:  "john.doe@example.com",
				Role:   "admin",
				Phones: []string{"986745321"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("validation error: field 'Age'"),
				},
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("validation error: field 'Email'"),
				},
			},
		},
	}

	for i, tt := range testsSuccess {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, err, tt.expectedErr)
		})
	}

	for i, tt := range testsFail {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if !errors.As(err, &ValidationErrors{}) {
				t.Errorf("unexpected error: got %v, expected %v", err, tt.expectedErr)
			} else {
				require.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
