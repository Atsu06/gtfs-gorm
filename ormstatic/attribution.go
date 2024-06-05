package ormstatic

import "database/sql"

type Attribution struct {
	Id               uint           `gorm:"primaryKey;auto_increment;not null"`
	AttributionId    sql.NullString `gorm:"primaryKey"`
	AgencyId         sql.NullString `gorm:"primaryKey"`
	RouteId          sql.NullString `gorm:"primaryKey"`
	TripId           sql.NullString `gorm:"primaryKey"`
	OrganizationName string         `gorm:"not null"`
	IsProducer       int16          `gorm:"default:0"`
	IsOperator       int16          `gorm:"default:0"`
	IsAuthority      int16          `gorm:"default:0"`
	AttributionUrl   sql.NullString
	AttributionEmail sql.NullString
	AttributionPhone sql.NullString
}

func (Attribution) TableName() string {
	return "attributions"
}

func (a Attribution) GetAttributionId() any {
	if a.AttributionId.Valid {
		return a.AttributionId.String
	}
	return nil
}

func (a Attribution) GetAgencyId() any {
	if a.AgencyId.Valid {
		return a.AgencyId.String
	}
	return nil
}

func (a Attribution) GetRouteId() any {
	if a.RouteId.Valid {
		return a.RouteId.String
	}
	return nil
}

func (a Attribution) GetTripId() any {
	if a.TripId.Valid {
		return a.TripId.String
	}
	return nil
}

func (a Attribution) GetOrganizationName() any {
	return a.OrganizationName
}

func (a Attribution) GetIsProducer() any {
	return a.IsProducer
}

func (a Attribution) GetIsOperator() any {
	return a.IsOperator
}

func (a Attribution) GetIsAuthority() any {
	return a.IsAuthority
}

func (a Attribution) GetAttributionUrl() any {
	if a.AttributionUrl.Valid {
		return a.AttributionUrl.String
	}
	return nil
}

func (a Attribution) GetAttributionEmail() any {
	if a.AttributionEmail.Valid {
		return a.AttributionEmail.String
	}
	return nil
}

func (a Attribution) GetAttributionPhone() any {
	if a.AttributionPhone.Valid {
		return a.AttributionPhone.String
	}
	return nil
}
