package main

type UserEvent struct {
	After  *UserEventPayload `json:"after"`
	Before *UserEventPayload `json:"before"`
	Op     string            `json:"op"`
}

type UserEventPayload struct {
	UserID   string `json:"_id"`
	Password string `json:"password"`
}
