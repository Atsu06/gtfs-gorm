package gtfsschedule

type Agency struct {
	AgencyID       string `gorm:"primary_key"`
	AgencyName     string `gorm:"not null"`
	AgencyUrl      string `gorm:"not null"`
	AgencyTimezone string `gorm:"not null"`
	AgencyLang     *string
	AgencyPhone    *string
	AgencyFareUrl  *string
	AgencyEmail    *string
}
