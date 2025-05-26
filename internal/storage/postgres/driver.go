package postgres

import (
	"fmt"

	"github.com/kartikey1188/build-in-progress_01/internal/models"
)

func (p *Postgres) EndDelivery(requestID int64) error {
	// Fetching the pickup request
	var request models.PickupRequest
	if err := p.GormDB.First(&request, "request_id = ?", requestID).Error; err != nil {
		return fmt.Errorf("pickup request not found: %w", err)
	}

	// Updating the status to completed
	request.Status = "Completed"
	if err := p.GormDB.Save(&request).Error; err != nil {
		return fmt.Errorf("failed to update pickup request status: %w", err)
	}

	return nil
}
