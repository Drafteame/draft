package utils_test

import (
	"testing"

	"github.com/Drafteame/draft/utils"
)

func Test_ToSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		data string
		exp  string
	}{
		{"Simple word", "Users", "users"},
		{"Word with underscore", "Users_", "users_"},
		{"Words with underscore", "Users_Users", "users_users"},
		{"Words wit dot", "Users.Users", "users_users"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ToSnakeCase(tt.data); got != tt.exp {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.exp)
			}
		})
	}
}
