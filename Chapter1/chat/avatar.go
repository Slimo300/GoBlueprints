package main

import (
	"errors"
	"io/ioutil"
	"log"
	"path"
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
	GetAvatarURL(ChatUser) (string, error)
}

//AuthAvatar is a type for obtaining avatar info
type AuthAvatar struct{}

// UseAuthAvatar is an object of AuthAvatar
var UseAuthAvatar AuthAvatar

// GetAvatarURL returns avatar_url in form of a string or error if something goes off
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

// GravatarAvatar type provides functionality for using gravatar service
type GravatarAvatar struct{}

// UseGravatarAvatar is an object of GravatarAvatar
var UseGravatarAvatar GravatarAvatar

// GetAvatarURL method implementation for GravatarAvatar type
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar type provides funcionality for storing uploaded by users pictures
type FileSystemAvatar struct{}

// UseFileSystemAvatar is an object of FileSystemAvatar
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL method implementation for FileSystemAvatar type
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {

	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		log.Println("Err when ReadDir")
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
			return "/avatars/" + file.Name(), nil
		}
	}

	log.Println("Escaping Function")
	return "", ErrNoAvatarURL
}

//TryAvatars type provides funcionality for running multiple Avatar
// implementations in given order
type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
