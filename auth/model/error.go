package model

import "errors"

var WrongPasswordError = errors.New("Wrong password")
var InvalidJWTError = errors.New("Invalid JWT Token")
