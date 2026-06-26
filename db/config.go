package db

type DB string

const Ride DB = "ride"

func (receiver DB) String() string {
	return string(receiver)
}
