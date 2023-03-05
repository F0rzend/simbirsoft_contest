package animal

import (
	"net/http"
	"time"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Response struct {
	ID                 int64      `json:"id"`
	Types              []int64    `json:"animalTypes"`
	Weight             float32    `json:"weight"`
	Length             float32    `json:"length"`
	Height             float32    `json:"height"`
	Gender             string     `json:"gender"`
	LifeStatus         string     `json:"lifeStatus"`
	ChippingDateTime   *time.Time `json:"chippingDateTime"`
	ChipperID          int        `json:"chipperId"`
	ChippingLocationID int        `json:"chippingLocationId"`
	VisitedLocations   []int      `json:"visitedLocations"`
	DeathDateTime      *time.Time `json:"deathDateTime"`
}

func (r *Response) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type SearchParams struct {
	StartDateTime      *time.Time `json:"startDateTime"`
	EndDateTime        *time.Time `json:"endDateTime"`
	ChipperID          int        `json:"chipperId"`
	ChippingLocationID int64      `json:"chippingLocationId"`
	LifeStatus         string     `json:"lifeStatus" validate:"live-status"`
	Gender             string     `json:"gender" validate:"gender"`
	From               int        `json:"from" validate:"gte=0"`
	Size               int        `json:"size" validate:"gt=0"`
}

func NewSearchParams(r *http.Request) (*SearchParams, error) {
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

	chipperID, invalidParameter := common.GetIntFromQuery(values, "chipperId", 0)
	if invalidParameter != nil {
		invalid = append(invalid, *invalidParameter)
	}

	chippingLocationID, invalidParameter := common.GetInt64FromQuery(values, "chippingLocationId", 0)
	if invalidParameter != nil {
		invalid = append(invalid, *invalidParameter)
	}

	liveStatus := values.Get("liveStatus")
	gender := values.Get("gender")

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

	params := &SearchParams{
		StartDateTime:      startDateTime,
		EndDateTime:        endDateTime,
		ChipperID:          chipperID,
		ChippingLocationID: chippingLocationID,
		LifeStatus:         liveStatus,
		Gender:             gender,
		From:               from,
		Size:               size,
	}

	if err := tv.ValidateStruct(params); err != nil {
		return nil, err
	}

	return params, nil
}

type ResponseList []*Response

func (ResponseList) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
