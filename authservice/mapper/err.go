package mapper

import "errors"

var (
	ErrorAuthorizationTokenExpire error = errors.New("JWT token non valid. ")
	ErrorNonValidRefreshToken     error = errors.New("Non valid refresh token. ")
	ErrorNonValidAccessToken      error = errors.New("Non valid access token. ")
	ErrorNonValidData             error = errors.New("Non valid data. ")
	ErrorNonExistUser             error = errors.New("User isn't exist. ")
	ErrorBadDataOperation         error = errors.New("Some problems with data operation. ")
	ErrorRetryingPasswordChange   error = errors.New("Retrying password change. ")
)
