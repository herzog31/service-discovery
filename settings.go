package main

type Settings struct {
	Hostname         string
	Notification     bool
	NotificationLog  bool
	SaveLogs         bool
	SaveLogsDays     int64
	SaveLogsInterval int64
	HipChatToken     string
	HipChatRoom      string
}
