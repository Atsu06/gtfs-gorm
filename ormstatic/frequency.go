package ormstatic

import "gorm.io/datatypes"

type Frequency struct {
	TripId      string         `gorm:"primaryKey;index;not null"`
	StartTime   datatypes.Time `gorm:"index;not null"`
	EndTime     datatypes.Time `gorm:"index;not null"`
	HeadwaySecs int            `gorm:"not null"`
	ExactTimes  int16          `gorm:"default:0"`
}

func (Frequency) TableName() string {
	return "frequencies"
}

func (f Frequency) GetTripId() any {
	return f.TripId
}

func (f Frequency) GetStartTime() any {
	return f.StartTime
}

func (f Frequency) GetEndTime() any {
	return f.EndTime
}

func (f Frequency) GetHeadwaySecs() any {
	return f.HeadwaySecs
}

func (f Frequency) GetExactTimes() any {
	return f.ExactTimes
}
