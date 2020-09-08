package models

const (
	ORAKKI_STATE_INIT = iota
	ORAKKI_STATE_READY
	ORAKKI_STATE_PAUSED
	ORAKKI_STATE_PLAYING
	ORAKKI_STATE_EXIT
	ORAKKI_STATE_PANIC
)

const MSG_FETCH_ORAKKI_STATE = "MSG_FETCH_ORAKKI_STATE"
const MSG_HANDLE_SETUP_OFFER = "MSG_HANDLE_SETUP_OFFER"
