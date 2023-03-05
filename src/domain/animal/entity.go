package animal

import "time"

type Entity struct {
	ID                 int64
	Types              []int64
	Weight             float32
	Length             float32
	Height             float32
	Gender             string
	LifeStatus         string
	ChippingDateTime   *time.Time
	ChipperID          int
	ChippingLocationID int
	VisitedLocations   []int
	DeathDateTime      *time.Time
}
