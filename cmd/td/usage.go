package main

const (
	usage = `

--Normal Mode--

↓ - move cursor one line down
↑ - move cursor one line up
a - add a new task(move to additional mode)
d - remove a task
→ - edit the task name(mode to edit mode)
? - help(switch to help mode)
enter - mark as done
t - switch to done tasks list mode
esc, ctrl+c - save tasks and close this app

--Done Tasks List Mode--

↓ - move cursor one line down
↑ - move cursor one line up
d - remove a task
t - switch to normal mode
enter - mark as done
ctrl+c - save tasks and close this app

--Edit Mode--

← - go back
esc - switch to normal mode
enter - submit

--Help Mode--
esc - switch to normal mode
`
)
