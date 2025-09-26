package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/voioo/td/internal/task"
)

// normalUpdate handles updates in normal mode.
func (m *Model) normalUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.updateTaskCache()
		switch {
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.taskCache) {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 1 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Right):
			if m.cursor == 0 || m.cursor > len(m.taskCache) {
				break
			}
			taskToEdit := m.taskCache[m.cursor-1]
			m.mode = ModeEdit
			m.editTaskNameInput.Placeholder = taskToEdit.Name
			m.editTaskNameInput.SetValue("")
			return m, m.editTaskNameInput.Focus()
		case key.Matches(msg, m.keys.Enter):
			if m.cursor == 0 || m.cursor > len(m.taskCache) {
				break
			}
			taskToComplete := m.taskCache[m.cursor-1]

			oldState := taskToComplete.IsDone
			if completedTask := m.taskManager.CompleteTask(taskToComplete.ID); completedTask != nil {
				m.undoManager.PushUndo(task.Action{
					Type:     task.ActionTypeComplete,
					Task:     completedTask,
					OldState: oldState,
					NewState: true,
				})
				m.invalidateCache()
			}

			if len(m.taskCache) == 0 {
				m.cursor = 0
			} else if m.cursor > len(m.taskCache) {
				m.cursor = len(m.taskCache)
			}
		case key.Matches(msg, m.keys.Priority):
			if len(m.taskCache) > 0 && m.cursor > 0 && m.cursor <= len(m.taskCache) {
				taskToUpdate := m.taskCache[m.cursor-1]
				oldPriority := taskToUpdate.Priority
				newPriority := (oldPriority + 1) % 4
				if updatedTask := m.taskManager.SetTaskPriority(taskToUpdate.ID, newPriority); updatedTask != nil {
					m.undoManager.PushUndo(task.Action{
						Type:     task.ActionTypePriority,
						Task:     updatedTask,
						OldState: oldPriority,
						NewState: newPriority,
					})
					m.invalidateCache()
					m.followTask(taskToUpdate.ID)
				}
			}
		case key.Matches(msg, m.keys.Add):
			m.mode = ModeAdditional
			return m, m.newTaskNameInput.Focus()
		case key.Matches(msg, m.keys.Delete):
			if m.cursor == 0 {
				break
			}

			taskToDelete := m.taskCache[m.cursor-1]
			oldState := taskToDelete.IsDone
			if deletedTask := m.taskManager.DeleteTask(taskToDelete.ID); deletedTask != nil {
				m.undoManager.PushUndo(task.Action{
					Type:     task.ActionTypeDelete,
					Task:     deletedTask,
					OldState: oldState,
				})
				m.invalidateCache()
			}

			if len(m.taskCache) == 0 {
				m.cursor = 0
			} else if m.cursor > len(m.taskCache) {
				m.cursor = len(m.taskCache)
			}
		case key.Matches(msg, m.keys.ListType):
			if m.mode == ModeDoneTaskList {
				if len(m.taskManager.GetTasks()) == 0 {
					m.cursor = 0
				} else {
					m.cursor = 1
				}
				m.mode = ModeNormal
				return m, nil
			}
			if len(m.taskManager.GetDoneTasks()) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
			m.mode = ModeDoneTaskList
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Filter):
			m.filter = (m.filter + 1) % 5
			m.invalidateCache()
			m.updateTaskCache()
			return m, nil
		case key.Matches(msg, m.keys.Help):
			m.mode = ModeHelp
		case key.Matches(msg, m.keys.Undo):
			if m.undoManager.Undo(m.taskManager) {
				m.invalidateCache()
			}
		case key.Matches(msg, m.keys.Redo):
			if m.undoManager.Redo(m.taskManager) {
				m.invalidateCache()
			}
		case key.Matches(msg, m.keys.PriorityNone):
			if len(m.taskCache) > 0 && m.cursor > 0 && m.cursor <= len(m.taskCache) {
				taskToUpdate := m.taskCache[m.cursor-1]
				if updatedTask := m.taskManager.SetTaskPriority(taskToUpdate.ID, task.PriorityNone); updatedTask != nil {
					m.undoManager.PushUndo(task.Action{
						Type:     task.ActionTypePriority,
						Task:     updatedTask,
						OldState: taskToUpdate.Priority,
						NewState: task.PriorityNone,
					})
					m.invalidateCache()
					m.followTask(taskToUpdate.ID)
				}
			}
		case key.Matches(msg, m.keys.PriorityLow):
			if len(m.taskCache) > 0 && m.cursor > 0 && m.cursor <= len(m.taskCache) {
				taskToUpdate := m.taskCache[m.cursor-1]
				if updatedTask := m.taskManager.SetTaskPriority(taskToUpdate.ID, task.PriorityLow); updatedTask != nil {
					m.undoManager.PushUndo(task.Action{
						Type:     task.ActionTypePriority,
						Task:     updatedTask,
						OldState: taskToUpdate.Priority,
						NewState: task.PriorityLow,
					})
					m.invalidateCache()
					m.followTask(taskToUpdate.ID)
				}
			}
		case key.Matches(msg, m.keys.PriorityMedium):
			if len(m.taskCache) > 0 && m.cursor > 0 && m.cursor <= len(m.taskCache) {
				taskToUpdate := m.taskCache[m.cursor-1]
				if updatedTask := m.taskManager.SetTaskPriority(taskToUpdate.ID, task.PriorityMedium); updatedTask != nil {
					m.undoManager.PushUndo(task.Action{
						Type:     task.ActionTypePriority,
						Task:     updatedTask,
						OldState: taskToUpdate.Priority,
						NewState: task.PriorityMedium,
					})
					m.invalidateCache()
					m.followTask(taskToUpdate.ID)
				}
			}
		case key.Matches(msg, m.keys.PriorityHigh):
			if len(m.taskCache) > 0 && m.cursor > 0 && m.cursor <= len(m.taskCache) {
				taskToUpdate := m.taskCache[m.cursor-1]
				if updatedTask := m.taskManager.SetTaskPriority(taskToUpdate.ID, task.PriorityHigh); updatedTask != nil {
					m.undoManager.PushUndo(task.Action{
						Type:     task.ActionTypePriority,
						Task:     updatedTask,
						OldState: taskToUpdate.Priority,
						NewState: task.PriorityHigh,
					})
					m.invalidateCache()
					m.followTask(taskToUpdate.ID)
				}
			}
		case key.Matches(msg, m.keys.Home):
			if len(m.taskCache) > 0 {
				m.cursor = 1
			} else {
				m.cursor = 0
			}
		case key.Matches(msg, m.keys.End):
			if len(m.taskCache) > 0 {
				m.cursor = len(m.taskCache)
			} else {
				m.cursor = 0
			}
		case key.Matches(msg, m.keys.ClearCompleted):
			doneTasks := m.taskManager.GetDoneTasks()
			if len(doneTasks) > 0 {
				// Create undo action for bulk delete
				m.undoManager.PushUndo(task.Action{
					Type: task.ActionTypeDelete,
					Task: &task.Task{IsDone: true}, // Marker for bulk operation
				})
				// Clear all completed tasks
				m.taskManager = task.NewTaskManager(m.taskManager.GetTasks(), []*task.Task{}, m.taskManager.GetNextID())
			}
		}
	}

	return m, nil
}

// doneTaskListUpdate handles updates in done task list mode.
func (m *Model) doneTaskListUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		doneTasks := m.taskManager.GetDoneTasks()
		switch {
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(doneTasks) {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 1 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Delete):
			if m.cursor == 0 {
				break
			}
			m.taskManager.DeleteTask(doneTasks[m.cursor-1].ID)
			if len(doneTasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
		case key.Matches(msg, m.keys.Enter):
			if m.cursor == 0 {
				break
			}
			t := doneTasks[m.cursor-1]

			oldState := t.IsDone
			if uncompletedTask := m.taskManager.UncompleteTask(t.ID); uncompletedTask != nil {
				m.undoManager.PushUndo(task.Action{
					Type:     task.ActionTypeUncomplete,
					Task:     uncompletedTask,
					OldState: oldState,
					NewState: false,
				})
				m.invalidateCache()
			}

			if len(doneTasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Undo):
			if m.undoManager.Undo(m.taskManager) {
				m.invalidateCache()
			}
		case key.Matches(msg, m.keys.Redo):
			if m.undoManager.Redo(m.taskManager) {
				m.invalidateCache()
			}
		case key.Matches(msg, m.keys.Escape):
			m.mode = ModeNormal
			return m, nil
		}
	}

	return m, nil
}

// addingTaskUpdate handles updates in task adding mode.
func (m *Model) addingTaskUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Escape):
			m.newTaskNameInput.Reset()
			m.mode = ModeNormal
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.mode = ModeNormal
			m.newTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			taskName := SanitizeTaskName(m.newTaskNameInput.Value())
			if err := ValidateTaskName(taskName); err != nil {
				// Could show error message here, for now just ignore invalid input
				return m, nil
			}

			addedTask := m.taskManager.AddTask(taskName)
			m.undoManager.PushUndo(task.Action{
				Type: task.ActionTypeAdd,
				Task: addedTask,
			})

			m.invalidateCache()
			m.newTaskNameInput.Reset()
			m.mode = ModeNormal
			return m, nil
		case key.Matches(msg, m.keys.Undo):
			if m.undoManager.Undo(m.taskManager) {
				m.invalidateCache()
			}
			m.mode = ModeNormal
			m.newTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Redo):
			if m.undoManager.Redo(m.taskManager) {
				m.invalidateCache()
			}
			m.mode = ModeNormal
			m.newTaskNameInput.Reset()
			return m, nil
		}
	}

	m.newTaskNameInput, cmd = m.newTaskNameInput.Update(msg)
	return m, cmd
}

// editTaskUpdate handles updates in task editing mode.
func (m *Model) editTaskUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Escape):
			m.editTaskNameInput.Reset()
			m.mode = ModeNormal
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.mode = ModeNormal
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			newName := SanitizeTaskName(m.editTaskNameInput.Value())
			if err := ValidateTaskName(newName); err != nil {
				// Could show error message here, for now just ignore invalid input
				return m, nil
			}

			m.updateTaskCache()
			if m.cursor > 0 && m.cursor <= len(m.taskCache) {
				taskToEdit := m.taskCache[m.cursor-1]
				oldName := taskToEdit.Name

				if oldName != newName {
					if updatedTask := m.taskManager.UpdateTaskName(taskToEdit.ID, newName); updatedTask != nil {
						m.undoManager.PushUndo(task.Action{
							Type:     task.ActionTypeEdit,
							Task:     updatedTask,
							OldState: oldName,
							NewState: newName,
						})
					}
				}
			}

			m.mode = ModeNormal
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Undo):
			if m.undoManager.Undo(m.taskManager) {
				m.invalidateCache()
			}
			m.mode = ModeNormal
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Redo):
			if m.undoManager.Redo(m.taskManager) {
				m.invalidateCache()
			}
			m.mode = ModeNormal
			m.editTaskNameInput.Reset()
			return m, nil
		}
	}

	m.editTaskNameInput, cmd = m.editTaskNameInput.Update(msg)
	return m, cmd
}

// helpUpdate handles updates in help mode.
func (m *Model) helpUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.mode = ModeNormal
		}
	}

	return m, nil
}
