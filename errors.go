package workersync

import "errors"

var (
	ErrNoWork                = errors.New("there is no work to do")
	ErrFailToUpdateStatus    = errors.New("failed to update status")
	ErrFailToAcquireWork     = errors.New("failed to acquire work")
	ErrFailToFinalizeAcquire = errors.New("failed to acquire work")

	ErrFailToInitWorkScheduler = errors.New("failed to initialize work scheduler")
)
