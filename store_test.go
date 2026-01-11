package main

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStoreLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	store := NewStore(path)
	tasks, err := store.Load()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty tasks, got %d", len(tasks))
	}
}

func TestStoreSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	store := NewStore(path)
	today := time.Date(2025, 1, 11, 0, 0, 0, 0, time.Local)
	tasks := []Task{
		NewTask("Task 1", today),
		NewTask("Task 2", today),
	}

	err := store.Save(tasks)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(loaded))
	}
	if loaded[0].Title != "Task 1" {
		t.Errorf("expected title 'Task 1', got %q", loaded[0].Title)
	}
}
