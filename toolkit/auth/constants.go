package auth

const (
	// Header / Context keys
	ContextUserIDKey   = "userId"
	ContextTenantIDKey = "tenantId"
	ContextClientIDKey = "clientId"

	// Redis key pattern
	RedisTokenKeyPattern = "token:%s:%s" // token:{clientId}:{userId}

	// Redis hash fields
	RedisFieldToken         = "token"
	RedisFieldActiveTimeout = "activeTimeout"
	RedisFieldCurrentTime   = "currentTime"
	RedisFieldExpireTime    = "expireTime"
)
