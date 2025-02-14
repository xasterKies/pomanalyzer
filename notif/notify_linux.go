package notif

import "os/exec"

var command = exec.Command

// Send notification for linux system
func (n *Notify) Send() error {
	notifCmdName := "notify-send"

	notifCmd, err := exec.LookPath(notifCmdName)

	if err != nil {
		return err
	}

	notifCommand := command(notifCmd, "-u", n.severity.String(), n.title, n.message)

	return notifCommand.Run()
}
