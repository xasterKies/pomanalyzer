package notif

import (
	"fmt"
	"os/exec"
)

var command = exec.Command

// Send notification for Windows
func (n *Notify) Send() error {
	notifCmdName := "powershell.exe"

	notifCmd, err := exec.LookPath(notifCmdName)
	if err != nil {
		return err
	}

	// Powershell notification script - It is loosely based on the BaloonTip script
	// developed by Boe Prox
	psscript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms
	$notify = New-Object System.Windows.Forms.NotifyIcon
	$notify.Icon = [System.Drawing.SystemIcons]::Information
	$notify.BalloonTipIcon = %q
	$notify.BalloonTipTitle = %q
	$notify.BalloonTipText = %q
	$notify.Visible = $True
	$notify.ShowBalloonTip(10000)`,
		n.severity, n.title, n.message,
	)

	// Slice of strings with required powershell arguements to run silently
	args := []string{
		"-NoProfile",
		"-NonInteractive",
	}

	//Append the script to the arguements slice to pass it to the function that creates the command
	args = append(args, psscript)

	notifCommand := command(notifCmd, args...)
	return notifCommand.Run()
}
