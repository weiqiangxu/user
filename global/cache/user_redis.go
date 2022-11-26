package cache

import "fmt"

const OneHour = 3600

// GetUserInfoCacheKey user cache
func GetUserInfoCacheKey(userID string) (string, int) {
	return fmt.Sprintf("user_info_%s", userID), OneHour
}
