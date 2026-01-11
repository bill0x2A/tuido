package main

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	DueDate   time.Time `json:"due_date"`
}

func NewTask(title string, dueDate time.Time) Task {
	normalizedDueDate := normalizeDate(dueDate)
	return Task{
		ID:        uuid.New().String(),
		Title:     title,
		Done:      false,
		CreatedAt: normalizedDueDate,
		DueDate:   normalizedDueDate,
	}
}

func (t Task) IsRolledOver() bool {
	return !sameDay(t.CreatedAt, t.DueDate)
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
