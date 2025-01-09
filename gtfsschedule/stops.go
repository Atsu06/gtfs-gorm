package gtfsschedule

type Stop struct {
	StopId             string `gorm:"primary_key"`
	StopCode           *string
	StopName           string `gorm:"not null"`
	TtsStopName        *string
	StopDesc           *string
	StopLat            float64 `gorm:"not null"`
	StopLon            float64 `gorm:"not null"`
	ZoneId             *string
	StopUrl            *string
	LocationType       *string
	ParentStation      int `gorm:"not null"`
	StopTimezone       *string
	WheelchairBoarding *string
	LevelId            *string
	PlatformCode       *string
}
