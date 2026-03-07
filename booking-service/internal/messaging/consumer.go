package messaging

import (
	"quicket/booking-service/internal/snapshot"
)

type EventConsumer struct {
	snapshotService snapshot.EventSnapshotService
}

func NewEventConsumer(snapshotService snapshot.EventSnapshotService) *EventConsumer {
	return &EventConsumer{
		snapshotService: snapshotService,
	}
}

/* EXAMPLE
func (c *EventConsumer) HandleEventUpdated(body []byte) error {
	var event snapshot.EventUpdatedMessage
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}

	return c.snapshotService.UpdateEventSnapshot(context.Background(), event)
}
*/
