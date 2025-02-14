package notif

import (
	"fmt"
	"log"
	"os/exec"
)

var command = exec.Command

// Send sends a notification for macOS using terminal-notifier,
// then plays a notification sound using afplay.
func (n Notify) Send() error {
	// Send the notification using terminal-notifier.
	notifCmdName := "terminal-notifier"
	notifCmd, err := exec.LookPath(notifCmdName)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("(%s) %s", n.severity, n.title)
	notifCommand := exec.Command(notifCmd, "-title", title, "-message", n.message)
	if err := notifCommand.Run(); err != nil {
		return err
	}

	soundCmdName, err := exec.LookPath("afplay")
	if err != nil {
		return err
	}

	soundFile := "/System/Library/Sounds/Glass.aiff"
	soundCommand := command(soundCmdName, soundFile)
	if err := soundCommand.Run(); err != nil {
		log.Printf("Failed to play sound: %v", err)
	}

	return nil
}
