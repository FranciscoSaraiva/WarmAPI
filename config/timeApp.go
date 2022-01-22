package config

import "time"

/*
*	Public Functions
 */

// GetDateForRequest function that returns the date to search for in
// the Greenhouse requests.
// It reads the date from a time.Now() and subtracts the date from the number
// received in parameters.
// It then slices the time from the string and returns only the date
// in the format yyyy-mm-dd.
func GetDateForRequest(daysBefore int) string {
	timestamp := time.Now()
	pastTime := timestamp.AddDate(0, 0, -daysBefore)
	dateCandidate, err := time.Parse("2006-01-02", pastTime.String()[0:10])
	CheckError(err)
	return dateCandidate.String()[0:10] // 0:10 equates to the format yyyy-mm-dd
}
