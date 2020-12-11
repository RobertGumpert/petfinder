package mapper

import "errors"

var (
	ErrorNonValidAccessToken       error = errors.New("Non valid access token. ")
	ErrorNonValidData              error = errors.New("Non valid data. ")
	ErrorMessageDoesntBelongDialog error = errors.New("Message doesn't belong dialog ")
	ErrorNonExistDialog            error = errors.New("Dialog isn't exist. ")
	ErrorNonExistMessage           error = errors.New("Message isn't exist. ")
	ErrorBadDataOperation          error = errors.New("Some problems with data operation. ")
)
