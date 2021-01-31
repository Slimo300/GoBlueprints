package main

import (
	"errors"
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
	userid, ok := c.userData["userid"]
	if !ok {
		return "", ErrNoAvatarURL
	}

	useridStr, ok := userid.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}

	return "//www.gravatar.com/avatar/" + useridStr, nil
}

// FileSystemAvatar type provides funcionality for storing uploaded by users pictures
type FileSystemAvatar struct{}

// UseFileSystemAvatar is an object of FileSystemAvatar
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL method implementation for FileSystemAvatar type
func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	userid, ok := c.userData["userid"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	useridStr, ok := userid.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	return "/avatars/" + useridStr + ".jpg", nil
}
