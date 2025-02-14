package notif

import (
	"fmt"
	"os/exec"
)

var command = exec.Command

// Send sends notification for MacOs, if an error
// occurs it is returned and the notification is not sent.
func (n Notify) Send() error {
	notifCmdName := "terminal-notifier"

	notifCmd, err := exec.LookPath(notifCmdName)

	if err != nil {
		return err
	}

	title := fmt.Sprintf("(%s) %s", n.severity, n.title)
	notifCommand := command(notifCmd, "-title", title, "-message", n.message)

	return notifCommand.Run()
}
