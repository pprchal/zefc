//go:build linux
// +build linux

package utils

import (
	"fmt"
	structs "zefc/structs"

	"github.com/godbus/dbus/v5"
)

const (
	dbusNotifyInterface = "org.freedesktop.Notifications"
	dbusNotifyPath      = "/org/freedesktop/Notifications"
)

func PrintGui(showGUI bool, msg string) {
	fmt.Print(msg)
	if showGUI {
		sendNotification("zefc", msg, 0)
	}
}

func ShowGUI(zipFile string, errors int, zipCount int, pattern structs.Profile) {
	var msg string
	var urgency uint32 = 1 // Normal urgency

	if errors == 0 {
		msg = fmt.Sprintf("All %d files are ok", zipCount)
		urgency = 0 // Low urgency for success
	} else {
		msg = fmt.Sprintf("Some files are missing or differ - %d!", errors)
		urgency = 2 // Critical urgency for errors
	}

	sendNotification(zipFile, msg, urgency)
}

func sendNotification(summary, body string, urgency uint32) {
	conn, err := dbus.SessionBus()
	if err != nil {
		// Silently fail if D-Bus is not available
		return
	}
	defer conn.Close()

	obj := conn.Object(dbusNotifyInterface, dbusNotifyPath)
	call := obj.Call(
		dbusNotifyInterface+".Notify",
		0,
		"zefc",     // app_name
		uint32(0),  // replaces_id
		"",         // app_icon
		summary,    // summary
		body,       // body
		[]string{}, // actions
		map[string]dbus.Variant{ // hints
			"urgency": dbus.MakeVariant(urgency),
		},
		int32(5000), // timeout in ms
	)

	if call.Err != nil {
		// Silently ignore notification errors
		return
	}
}
