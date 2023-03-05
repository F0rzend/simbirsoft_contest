package visits

import "time"

type Entity struct {
	ID                           int64
	DateTimeOfVisitLocationPoint time.Time
	LocationPointID              int64
}
