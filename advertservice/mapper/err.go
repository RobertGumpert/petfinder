package mapper

import "errors"


var (
	ErrorNonValidData             error = errors.New("Non valid data. ")
	ErrorBadDataOperation         error = errors.New("Some problems with data operation. ")
)
