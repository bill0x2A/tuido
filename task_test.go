package main

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	title := "Buy groceries"
	dueDate := time.Date(2025, 1, 11, 0, 0, 0, 0, time.Local)

	task := NewTask(title, dueDate)

	if task.Title != title {
		t.Errorf("expected title %q, got %q", title, task.Title)
	}
	if task.Done {
		t.Error("new task should not be done")
	}
	if task.ID == "" {
		t.Error("task should have an ID")
	}
	if !task.DueDate.Equal(dueDate) {
		t.Errorf("expected due date %v, got %v", dueDate, task.DueDate)
	}
}

func TestTaskIsRolledOver(t *testing.T) {
	yesterday := time.Date(2025, 1, 10, 0, 0, 0, 0, time.Local)
	today := time.Date(2025, 1, 11, 0, 0, 0, 0, time.Local)

	task := NewTask("Test task", yesterday)

	if task.IsRolledOver() {
		t.Error("task created today should not be rolled over")
	}

	task.DueDate = today
	if !task.IsRolledOver() {
		t.Error("task with different due date than created date should be rolled over")
	}
}
