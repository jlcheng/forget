package db

type BatchBuffer struct {
	// notifications
	arrivals chan(struct{})
	finish   chan(struct{})

	buffer []struct{}
	begin func([]struct{})
}