package domain

// SchedulerTrigger represents a scheduler that can be triggered to update its schedule.
// This interface is typically implemented by infrastructure components that manage
// scheduled tasks, such as closing polls at their scheduled times.
type SchedulerTrigger interface {
	TriggerUpdate()
}
