package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: unable to get an avatar url")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

//AuthAvatar is a type for obtaining avatar info
type AuthAvatar struct{}

// UseAuthAvatar is an object of AuthAvatar
var UseAuthAvatar AuthAvatar

// GetAvatarURL returns avatar_url in form of a string or error if something goes off
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	return urlStr, nil
}

// GravatarAvatar type provides functionality for using gravatar service
type GravatarAvatar struct{}

// UseGravatarAvatar is an object of GravatarAvatar
var UseGravatarAvatar GravatarAvatar

// GetAvatarURL method implementation for GravatarAvatar type
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	email, ok := c.userData["email"]
	if !ok {
		return "", ErrNoAvatarURL
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}

	m := md5.New() //TODO: Hashing everytime, Memoization?
	io.WriteString(m, strings.ToLower(emailStr))
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
}
