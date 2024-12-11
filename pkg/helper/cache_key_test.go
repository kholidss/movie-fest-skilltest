package helper

import "testing"

func TestCacheKeyLockTrxCreditUser(t *testing.T) {
	tests := []struct {
		userID   string
		expected string
	}{
		{"user123", "movie-fest-skilltest:trx-credit-user:user123"},
		{"abc456", "movie-fest-skilltest:trx-credit-user:abc456"},
		{"testUser", "movie-fest-skilltest:trx-credit-user:testUser"},
		{"", "movie-fest-skilltest:trx-credit-user:"},
	}

	for _, tt := range tests {
		t.Run("CacheKeyLockTrxCreditUser", func(t *testing.T) {
			result := CacheKeyLockTrxCreditUser(tt.userID)
			if result != tt.expected {
				t.Errorf("expected %s, got %s for userID %v", tt.expected, result, tt.userID)
			}
		})
	}
}
