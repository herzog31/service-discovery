package main

type Settings struct {
	Hostname        string
	Notification    bool
	NotificationLog bool
	SaveLogs        bool
	SaveLogsDays    int
	HipChatToken    string
	HipChatRoom     string
}
