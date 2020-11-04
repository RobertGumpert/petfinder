package mapper

import "errors"

type TypeAdvert uint64

const (
	TypeLost          TypeAdvert = 1
	TypeFound         TypeAdvert = 2
	MaxLatitude       float64    = 90.0
	MinLatitude       float64    = -90.0
	MaxLongitude      float64    = 180.0
	MinLongitude      float64    = -180.0
	OneKilometerScale float64    = 0.00899
)

var (
	ErrorNonValidData             error = errors.New("Non valid data. ")
	ErrorBadDataOperation         error = errors.New("Some problems with data operation. ")

)
