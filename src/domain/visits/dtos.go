package visits

import (
	"net/http"
	"time"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Response struct {
	ID                           int64     `json:"id"`
	DateTimeOfVisitLocationPoint time.Time `json:"dateTimeOfVisitLocationPoint"`
	LocationPointID              int64     `json:"locationPointId"`
}

func (*Response) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type ResponsesList []Response

func (ResponsesList) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type GetAnimalLocationsQuery struct {
	StartDateTime *time.Time `json:"startDateTime"`
	EndDateTime   *time.Time `json:"endDateTime"`
	From          int        `json:"from" validate:"gte=0"`
	Size          int        `json:"size" validate:"gt=0"`
}

func NewGetAnimalLocationsQuery(r *http.Request) (*GetAnimalLocationsQuery, error) {
	tv, err := common.TranslatedValidatorFromRequest(r)
	if err != nil {
		return nil, err
	}

	values := r.URL.Query()

	var invalid []common.InvalidRequestParameter

	startDateTime, invalidParameter := common.GetDatetimeFromQuery(
		values,
		"startDateTime",
	)
	if invalidParameter != nil {
		invalid = append(invalid, *invalidParameter)
	}

	endDateTime, invalidParameter := common.GetDatetimeFromQuery(
		values,
		"endDateTime",
	)
	if invalidParameter != nil {
		invalid = append(invalid, *invalidParameter)
	}
	from, invalidParameter := common.GetIntFromQuery(values, "from", 0)
	if invalidParameter != nil {
		invalid = append(invalid, *invalidParameter)
	}

	size, invalidParameter := common.GetIntFromQuery(values, "from", 10)
	if invalidParameter != nil {
		invalid = append(invalid, *invalidParameter)
	}

	if len(invalid) != 0 {
		return nil, common.NewValidationError(invalid...)
	}

	params := &GetAnimalLocationsQuery{
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		From:          from,
		Size:          size,
	}

	if err := tv.ValidateStruct(params); err != nil {
		return nil, err
	}

	return params, nil
}
