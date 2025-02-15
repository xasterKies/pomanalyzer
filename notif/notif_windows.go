package notif

import (
	"fmt"
	"os/exec"
)

var command = exec.Command

// Send sends a Windows notification and plays a sound.
func (n *Notify) Send() error {
	notifCmdName := "powershell.exe"

	// Look up the path for powershell.exe
	notifCmd, err := exec.LookPath(notifCmdName)
	if err != nil {
		return err
	}

	// Build the PowerShell script for the balloon notification.
	// Using single quotes to enclose strings inside the script.
	psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms;
$notify = New-Object System.Windows.Forms.NotifyIcon;
$notify.Icon = [System.Drawing.SystemIcons]::Information;
$notify.BalloonTipIcon = '%s';
$notify.BalloonTipTitle = '%s';
$notify.BalloonTipText = '%s';
$notify.Visible = $True;
$notify.ShowBalloonTip(10000);`, n.severity, n.title, n.message)

	// Create the command to send the notification.
	notifArgs := []string{
		"-NoProfile",
		"-NonInteractive",
		"-Command", psScript,
	}
	notifCommand := command(notifCmd, notifArgs...)
	if err := notifCommand.Run(); err != nil {
		return err
	}

	soundScript := `(New-Object Media.SoundPlayer 'C:\Windows\Media\notify.wav').PlaySync();`
	soundArgs := []string{
		"-NoProfile",
		"-NonInteractive",
		"-Command", soundScript,
	}
	soundCommand := command(notifCmd, soundArgs...)
	return soundCommand.Run()
}
