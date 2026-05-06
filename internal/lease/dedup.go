package lease

// ShouldAlert determines whether an alert should be sent for the given lease
// entry. It returns false if the entry is nil, has already been alerted, or
// has already been marked as expired — preventing duplicate notifications.
func ShouldAlert(entry *Entry) bool {
	if entry == nil {
		return false
	}
	if entry.Alerted {
		return false
	}
	if entry.Expired {
		return false
	}
	return true
}
