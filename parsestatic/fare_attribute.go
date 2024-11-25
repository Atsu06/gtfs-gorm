package parsestatic

import (
	"fmt"
	"github.com/ITNS-LAB/gtfs-gorm/gtfsjp"
	"github.com/ITNS-LAB/gtfs-gorm/internal/dataframe"
	"github.com/ITNS-LAB/gtfs-gorm/internal/util"
)

func ParseFareAttributes(path string) ([]gtfsjp.FareAttribute, error) {
	var fareAttributes []gtfsjp.FareAttribute
	df, err := dataframe.OpenCsv(path)
	if err != nil {
		return fareAttributes, err
	}

	for df.HasNext() {
		_, err := df.Next()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		fareAttributes = append(fareAttributes, gtfsjp.FareAttribute{
			FareId:           util.IsBlank(df.GetElement("fare_id")),
			Price:            util.ParseFloat64(df.GetElement("price")),
			CurrencyType:     util.IsBlank(df.GetElement("currency_type")),
			PaymentMethod:    util.ParseEnum(df.GetElement("payment_method")),
			Transfers:        util.ParseEnum(df.GetElement("transfers")),
			AgencyId:         util.IsBlank(df.GetElement("agency_id")),
			TransferDuration: util.ParseInt(df.GetElement("transfer_duration")),
		})
	}
	return fareAttributes, nil
}
