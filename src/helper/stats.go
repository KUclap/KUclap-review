package helper

func GetNewStatsByCreated(oldN float64, oldstat float64, newStats float64) float64 {
	return ((newStats / 5 * 100) + (oldstat * oldN)) / (oldN + 1)
}

func GetNewStatsByDeleted(oldN float64, oldstat float64, newStats float64) float64 {
	if oldN - 1 <= 0 {
		return oldstat
	} 
		return ((oldstat * oldN) - (newStats / 5 * 100) ) / (oldN - 1)
}
