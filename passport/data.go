package passport

import (
	"sorgulat-api/passport/models"
	"sorgulat-api/timezones/utils"
)

var Countries = utils.LoadData[models.PassportCountries]("countries", "passport/")
