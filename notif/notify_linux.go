package notif

import (
	"os/exec"
)

var command = exec.Command

// Send notification for linux system
func (n *Notify) Send() error {
	notifCmdName := "notify-send"
	soundFile := "sounds/message-new-instant.oga"

	notifCmd, err := exec.LookPath(notifCmdName)

	if err != nil {
		return err
	}

	notifCommand := command(notifCmd, "-u", n.severity.String(), n.title, n.message)

	if err := notifCommand.Run(); err != nil {
		return err
	}

	// Check if `paplay` is available; otherwise, use `ogg123`
	soundCmdName, err := exec.LookPath("paplay")
	if err != nil {
		soundCmdName, err = exec.LookPath("ogg123")
		if err != nil {
			return nil
		}
	}

	soundCommand := exec.Command(soundCmdName, soundFile)
	return soundCommand.Run()
}
