package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Id         uuid.UUID
	Type       string
	Payload    json.RawMessage
	OccurredAt time.Time
}
