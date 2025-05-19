package domain

import "errors"

var ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
var ErrConfirmationTokenNotFound = errors.New("confirmation token not found")
var ErrUnsubscribeTokenNotFound = errors.New("unsubscribe token not found")
var ErrSubscriptionNotFound = errors.New("subscription not found")
var ErrCityNotFound = errors.New("city not found")
var ErrUnableToSubscribe = errors.New("failed to subscribe, try again later")
var ErrBadRequest = errors.New("bad request")
var ErrInternalServerError = errors.New("internal server error")
