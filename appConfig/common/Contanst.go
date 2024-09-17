package common

import "time"

// should be lowercase for error messages
const (
	Success                string = "success"
	ErrUnauthorized        string = "unauthorized"
	ErrNotFound            string = "not found"
	ErrInternalServerError string = "internal server error"
)

const (
	CodeSuccess             int32 = 200
	CodeNotFound            int32 = 404
	CodeInternalServerError int32 = 500
)

const ACCESSTOKEN_DURATION = 30 * time.Minute
const REFRESH_TOKEN_DURATION = 24 * time.Hour
