package models

import (
	"time"

	"rearatrox/event-booking-api/pkg/db"
)

type Event struct {
	ID          int64
	Name        string    `binding: "required"`
	Description string    `binding: "required"`
	Location    string    `binding: "required"`
	DateTime    time.Time `binding: "required"`
	CreatorID   int64
}

// events slice removed; models now use database storage.

func (e *Event) SaveEvent() error {
	query := `INSERT INTO events (name, description, location, datetime, creator_id) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	if err := db.DB.QueryRow(db.Ctx, query, e.Name, e.Description, e.Location, e.DateTime, e.CreatorID).Scan(&e.ID); err != nil {
		return err
	}
	return nil
}

func (e *Event) UpdateEvent() error {
	query := `UPDATE events SET name=$1, description=$2, location=$3, datetime=$4, creator_id=$5 WHERE id=$6`
	_, err := db.DB.Exec(db.Ctx, query, e.Name, e.Description, e.Location, e.DateTime, e.CreatorID, e.ID)
	return err
}

func (e *Event) DeleteEvent() error {
	query := `DELETE FROM events WHERE id=$1`
	_, err := db.DB.Exec(db.Ctx, query, e.ID)
	return err
}

func GetEvents() ([]Event, error) {
	query := `SELECT id, name, description, location, datetime, creator_id FROM events`
	rows, err := db.DB.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.CreatorID); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	var event Event
	query := `SELECT id, name, description, location, datetime, creator_id FROM events WHERE id=$1`
	row := db.DB.QueryRow(db.Ctx, query, id)
	if err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.CreatorID); err != nil {
		return nil, err
	}
	return &event, nil
}

func (e Event) Register(userId int64) error {
	query := `INSERT INTO event_registrations(event_id, user_id) VALUES ($1, $2)`
	_, err := db.DB.Exec(db.Ctx, query, e.ID, userId)
	return err
}

func (e Event) DeleteRegistration(userId int64) error {
	query := `DELETE FROM event_registrations WHERE event_id=$1 AND user_id=$2`
	_, err := db.DB.Exec(db.Ctx, query, e.ID, userId)
	return err
}
