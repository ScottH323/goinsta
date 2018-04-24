package goinsta

import "errors"

// ErrNotFound is returned if the request responds with a 404 status code
// i.e a non existent user
var ErrNotFound = errors.New("The specified data wasn't found.")

// ErrLoggedOut is returned if the request responds with a 400 status code
var ErrLoggedOut = errors.New("The account is logged out")

// ErrRateLimit is returned when the IG rate limiter is hit
var ErrRateLimit = errors.New("Rate limit hit, chill out")
