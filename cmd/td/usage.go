package main

const (
	usage = `
Task Manager Controls
====================

Normal Mode
----------
Navigation:
  ↑/k        Move cursor up
  ↓/j        Move cursor down
  ←/h        Move left
  →/l        Move right/edit task
  Enter      Mark task as done/undone

Task Actions:
  a          Add new task
  d          Delete task
  t          Toggle between active/completed tasks

Other:
  ?          Show/hide this help
  Esc/Ctrl+C Save and exit


Done Tasks Mode
--------------
Navigation:
  ↑/k        Move cursor up
  ↓/j        Move cursor down
  Enter      Mark task as undone

Task Actions:
  d          Delete task
  t          Switch to normal mode
  Esc/Ctrl+C Save and exit


Edit Mode
---------
  Enter      Save changes
  Esc        Cancel and return to normal mode
  Ctrl+C     Save and exit


Add Task Mode
------------
  Enter      Save new task
  Esc        Cancel and return to normal mode
  Ctrl+C     Save and exit


Press Esc to exit this help screen
`
)
