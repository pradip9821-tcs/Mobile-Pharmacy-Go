package utils

func Authorization(role int) bool {
	if role == 1 {
		return false
	}
	return true
}
