package logic

import "testing"

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		login   string
		wantErr bool
	}{
		{"user1", false},
		{"us", true},
		{"user!@#", true},
		{"valid_user-1", false},
		{"", true},
		{"user withspace", true},
		{"user.name", true},
		{"user-name", false},
		{"юзер", true},
		{"veryveryveryveryveryveryveryverylonguser", true},
	}
	for _, tt := range tests {
		err := validateLogin(tt.login)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateLogin(%q) error = %v; wantErr = %v", tt.login, err, tt.wantErr)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		wantErr bool
	}{
		{"password123", false},
		{"short", true},
		{"validPASS!@#", false},
		{"русскийПароль123№@", true},
		{"with space", true},
		{"", true},
		{"12345678", false},
		{"abc!@#$%^&*()_+=-", false},
		{"abc/\\", true},
		{"!@#$%^&*()_+=-", false},
		{"abcdefgh", false},
		{"1234567890", false},
		{"a2345678", false},
		{string(make([]byte, 65)), true},
	}
	for _, tt := range tests {
		err := validatePassword(tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("validatePassword(%q) error = %v; wantErr = %v", tt.password, err, tt.wantErr)
		}
	}
}
