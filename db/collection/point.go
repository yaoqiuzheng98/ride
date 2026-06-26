package collection

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"ride/db"
)

type Point struct {
	Id        int64   `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

// Upsert inserts the point or updates the matching row keyed by (latitude, longitude).
func (receiver Point) Upsert(ctx context.Context) (sql.Result, error) {
	_, err := db.GetClient().ExecContext(ctx,
		"INSERT INTO point (latitude, longitude, name) VALUES (?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE name = VALUES(name)",
		receiver.Latitude, receiver.Longitude, receiver.Name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (receiver Point) getName() string {
	return "point"
}

func GetPoints() *Points {
	return _points
}

var _points = &Points{}

type Points struct {
	rwMutex sync.RWMutex
	points  []Point
}

func (receiver *Points) getName() string {
	return "point"
}

func (receiver *Points) Load() error {
	rows, err := db.GetClient().Query("SELECT id, latitude, longitude, name FROM point")
	if err != nil {
		return err
	}
	defer rows.Close()

	var points []Point
	for rows.Next() {
		var p Point
		if err := rows.Scan(&p.Id, &p.Latitude, &p.Longitude, &p.Name); err != nil {
			return err
		}
		points = append(points, p)
	}
	if err = rows.Err(); err != nil {
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
	_, err := db.GetClient().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			" id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"+
			" latitude DOUBLE NOT NULL,"+
			" longitude DOUBLE NOT NULL,"+
			" name VARCHAR(255) NOT NULL DEFAULT '',"+
			" PRIMARY KEY (id),"+
			" UNIQUE KEY uk_lat_lng (latitude, longitude)"+
			" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		receiver.getName()))
	return err
}

func init() {
	if err := _points.ensureTable(); err != nil {
		panic(err)
	}
}
