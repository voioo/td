package main

import (
	"testing"

	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModel(t *testing.T) {
	t.Run("initialize model", func(t *testing.T) {
		m := initializeModel()
		model, ok := m.(model)
		if !ok {
			t.Fatal("expected model type")
		}

		if model.mode != ModeNormal {
			t.Errorf("expected initial mode to be ModeNormal, got %d", model.mode)
		}

		if model.cursor != 0 && model.cursor != 1 {
			t.Errorf("expected cursor to be 0 or 1, got %d", model.cursor)
		}
	})

	t.Run("add task", func(t *testing.T) {
		m := model{
			mode:             ModeAdditional,
			cursor:           0,
			tasks:            []*Task{},
			latestTaskID:     0,
			newTaskNameInput: input.New(),
			keys:             keys,
		}
		m.newTaskNameInput.Focus()
		m.newTaskNameInput.SetValue("Test task")

		newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		updatedModel := newModel.(model)

		if len(updatedModel.tasks) != 1 {
			t.Errorf("expected 1 task, got %d", len(updatedModel.tasks))
		}

		if updatedModel.tasks[0].Name != "Test task" {
			t.Errorf("expected task name 'Test task', got %s", updatedModel.tasks[0].Name)
		}

		if updatedModel.mode != ModeNormal {
			t.Errorf("expected mode to switch to ModeNormal, got %d", updatedModel.mode)
		}
	})

	t.Run("mark task as done", func(t *testing.T) {
		task := &Task{
			Name:   "Test task",
			ID:     1,
			IsDone: false,
		}

		m := model{
			mode:      ModeNormal,
			cursor:    1,
			tasks:     []*Task{task},
			keys:      keys,
			doneTasks: []*Task{}, // Initialize empty doneTasks slice
		}

		newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		updatedModel := newModel.(model)

		if len(updatedModel.tasks) != 0 {
			t.Errorf("expected task to be removed from tasks, got %d tasks", len(updatedModel.tasks))
		}

		if len(updatedModel.doneTasks) != 1 {
			t.Errorf("expected task to be added to doneTasks, got %d done tasks", len(updatedModel.doneTasks))
		}

		if !updatedModel.doneTasks[0].IsDone {
			t.Error("expected task to be marked as done")
		}
	})
}
