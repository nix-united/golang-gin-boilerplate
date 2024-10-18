package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestEncrypt(t *testing.T) {
	tests := []struct {
		name   string
		fields struct {
			cost int
		}
		args struct {
			pass string
		}
		wantErr bool
	}{
		{
			"test successful encrypting a password",
			struct{ cost int }{cost: bcrypt.DefaultCost},
			struct{ pass string }{pass: "test pass"},
			false,
		},
		{
			"test returning an error",
			// Below we set a cost that is greater than max available cost value.
			// So, it causes the InvalidCostError, which is returned
			// by bcrypt.GenerateFromPassword function.
			struct{ cost int }{cost: 100},
			struct{ pass string }{pass: "test pass"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBcryptEncoder(tt.fields.cost).Encrypt(tt.args.pass)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			}
		})
	}
}
