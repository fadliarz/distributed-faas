package application

type CheckpointEvent struct {
	CheckpointID string `json:"_id"`
	OutputURL    string `json:"output_url"`
}
