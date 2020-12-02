package mapper

import "time"

type TypeAdvert uint64

const (
	TypeLost                    TypeAdvert    = 1
	TypeFound                   TypeAdvert    = 2
	MaxLatitude                 float64       = 90.0
	MinLatitude                 float64       = -90.0
	MaxLongitude                float64       = 180.0
	MinLongitude                float64       = -180.0
	OneKilometerScale           float64       = 0.00899
	CompareAdvertTime           time.Duration = 24 * time.Hour
	LifetimeOfFoundAnimalAdvert time.Duration = 24 * 3 * time.Hour
	LifetimeOfLostAnimalAdvert  time.Duration = 24 * 14 * time.Hour
)
