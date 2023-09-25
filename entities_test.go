package bookclub_test

import (
	"testing"

	"github.com/gosom/bookclub"

	"github.com/stretchr/testify/require"
)

func TestNewEmail(t *testing.T) {
	cases := []struct {
		value string
		valid bool
	}{
		{"", false},
		{"a", false},
		{"a@", false},
		{"a@b", true}, // notice that this a valid email address
		{"a@b.c", true},
		{"a@b.c.d", true},
	}

	for _, c := range cases {
		_, err := bookclub.NewEmail(c.value)
		if c.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			require.Equal(t, bookclub.ErrInvalidEmail, err)
		}
	}
}

func TestNewPassword(t *testing.T) {
	cases := []struct {
		value string
		valid bool
	}{
		{"", false},
		{"a", false},
		{"12345678", false},
		{"1234567a", false},
		{"1234567aA", false},
		{"1234567aA1", false},
		{"1234567aA1!", true},
		{"1234567aA1!11111111111111113329494943933", false},
	}

	for _, c := range cases {
		passwd, err := bookclub.NewPassword(c.value)
		if c.valid {
			require.NoError(t, err)
			err = passwd.Compare(c.value, passwd)
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			require.Equal(t, bookclub.ErrInvalidPassword, err)
		}
	}
}

func Test_UserPassword(t *testing.T) {
	passwd, err := bookclub.NewPassword("1234567aA1!")
	require.NoError(t, err)

	user := bookclub.User{}

	user.SetPassword(passwd)

	err = user.ComparePassword("1234567aA1!")
	require.NoError(t, err)

	err = user.ComparePassword("1234567aA1!1")
	require.Error(t, err)
	require.Equal(t, bookclub.ErrPasswordMismatch, err)
}
