package utils

func Authorization(role int) bool {
	if role == 2 {
		return false
	}
	return true
}
