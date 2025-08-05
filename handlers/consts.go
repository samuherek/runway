package handlers

type HxHeaderCustom string

const (
	HxNotifications HxHeaderCustom = "X-Notifications"
	HxErrors        HxHeaderCustom = "X-Errors"
)

type HxHeader string

const (
	HxTrigger HxHeader = "HX-Trigger"
	HxPushUrl HxHeader = "HX-Push-Url"
)
