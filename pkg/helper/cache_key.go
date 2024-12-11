package helper

import "fmt"

// CacheKeyLockTrxCreditUser is a key of redis cache
func CacheKeyLockTrxCreditUser(userID string) string {
	return fmt.Sprintf("movie-fest-skilltest:trx-credit-user:%s", userID)
}
