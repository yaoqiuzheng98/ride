package table

import (
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ride/db"
)

type Point struct {
	gorm.Model
	Latitude  float64 `gorm:"uniqueIndex:uk_lat_lng" json:"latitude"`
	Longitude float64 `gorm:"uniqueIndex:uk_lat_lng" json:"longitude"`
	Name      string  `gorm:"size:255;not null;default:''" json:"name"`
}

func (Point) TableName() string {
	return "point"
}

// Upsert inserts the point or updates the matching row keyed by (latitude, longitude).
func (receiver Point) Upsert() error {
	return db.GetClient().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "latitude"}, {Name: "longitude"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&receiver).Error
}

func GetPoints() *Points {
	return _points
}

var _points = &Points{}

type Points struct {
	rwMutex sync.RWMutex
	points  []Point
}

func (receiver *Points) Load() error {
	var points []Point
	if err := db.GetClient().Find(&points).Error; err != nil {
		return err
	}
	receiver.rwMutex.Lock()
	defer receiver.rwMutex.Unlock()
	receiver.points = points
	return nil
}

func (receiver *Points) GetRefreshGap() time.Duration {
	return time.Minute
}

func (receiver *Points) GetPoints() []Point {
	receiver.rwMutex.RLock()
	defer receiver.rwMutex.RUnlock()
	if receiver.points == nil {
		return []Point{}
	}
	return receiver.points
}

// ensureTable makes sure the point table exists. Called once at startup.
func (receiver *Points) ensureTable() error {
	return db.GetClient().AutoMigrate(&Point{})
}

func init() {
	if err := _points.ensureTable(); err != nil {
		panic(err)
	}
}
