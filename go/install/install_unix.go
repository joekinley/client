// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

// +build linux freebsd openbsd

package install

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/keybase/client/go/libkb"
)

// Similar to the Brew install on OSX, the Unix install happens in two steps.
// First, the system package manager installs all the binaries as root. Second,
// an autostart file needs to be written to the user's home dir, so that
// Keybase launches when that user logs in. The second step is done the first
// time the user starts Keybase.
//
// ".desktop" files and the ~/.config/autostart directory are part of the
// freedesktop.org set of standards, which the popular desktop environments
// like Gnome and KDE all support. See
// http://standards.freedesktop.org/desktop-entry-spec/latest/.

// TODO: See if upgrading to Electron 1.8.x removes the need to start
// with XDG_CURRENT_DESKTOP=Unity; see
// https://github.com/electron/electron/issues/10887 .
const autostartFileText = `# This file is generated the first time Keybase starts, along with a sentinel
# file at ~/.config/keybase/autostart_created. As long as the sentinel exists,
# this file won't be automatically recreated, so you can edit it or delete it
# as you like. Manually changing autostart settings in the GUI can stomp on
# your edits here though.

[Desktop Entry]
Name=Keybase
Comment=Keybase Filesystem Service and GUI
Type=Application
Exec=env KEYBASE_AUTOSTART=1 run_keybase
`

const sentinelFileText = `This file is created the first time Keybase starts, along with
~/.config/autostart/keybase_autostart.desktop. As long as this
file exists, the autostart file won't be automatically recreated.
`

func autostartDir(context Context) string {
	// strip off the "keybase" folder on the end of the config dir
	return path.Join(context.GetConfigDir(), "..", "autostart")
}

func autostartFilePath(context Context) string {
	return path.Join(autostartDir(context), "keybase_autostart.desktop")
}

func sentinelFilePath(context Context) string {
	return path.Join(context.GetConfigDir(), "autostart_created")
}

// AutoInstall installs auto start on unix
func AutoInstall(context Context, _ string, _ bool, timeout time.Duration, log Log) ( /* newProc */ bool, error) {
	_, err := os.Stat(sentinelFilePath(context))
	if err == nil {
		// The sentinel exists. Don't recreate the autostart file.
		return false, nil
	} else if !os.IsNotExist(err) {
		// The error is something unexpected. Return it.
		return false, err
	}

	// The sentinel doesn't exist. Create the autostart file, and then create
	// the sentinel. This might stomp on old user edits one time, but we need
	// to do that to add in the KEYBASE_AUTOSTART variable.
	err = os.MkdirAll(autostartDir(context), 0755)
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(autostartFilePath(context), []byte(autostartFileText), 0644)
	if err != nil {
		return false, err
	}
	err = os.MkdirAll(context.GetConfigDir(), 0755)
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(sentinelFilePath(context), []byte(sentinelFileText), 0644)
	if err != nil {
		return false, err
	}

	return false, nil
}

// CheckIfValidLocation is not used on unix
func CheckIfValidLocation() error {
	return nil
}

// KBFSBinPath returns the path to the KBFS executable
func KBFSBinPath(runMode libkb.RunMode, binPath string) (string, error) {
	return kbfsBinPathDefault(runMode, binPath)
}

// kbfsBinName returns the name for the KBFS executable
func kbfsBinName() string {
	return "kbfsfuse"
}

func updaterBinName() (string, error) {
	return "", fmt.Errorf("Updater isn't supported on unix")
}

// RunApp starts the app
func RunApp(context Context, log Log) error {
	// TODO: Start app, see run_keybase: /opt/keybase/Keybase
	return nil
}

func InstallLogPath() (string, error) {
	return "", nil
}

// WatchdogLogPath doesn't exist on linux as an independent log file
func WatchdogLogPath(string) (string, error) {
	return "", nil
}

func SystemLogPath() string {
	return ""
}
