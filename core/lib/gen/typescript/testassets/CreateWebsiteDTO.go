package testassets

type CreateWebsiteDTO struct {
	Title                   string
	BaseURL                 string
	ScheduleIntervalMinutes int64
	IsActive                int
	IDs                     []int64
}
