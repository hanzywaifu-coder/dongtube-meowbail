package meowbail

import "github.com/hanzywaifu-coder/dongtube-meowbail/core"

// Expose core types to package meowbail
type LIDResolver = core.LIDResolver
type RetrySpiralingTracker = core.RetrySpiralingTracker

func NewLIDResolver() *LIDResolver {
	return core.NewLIDResolver()
}

func NewRetrySpiralingTracker(maxRetries int) *RetrySpiralingTracker {
	return core.NewRetrySpiralingTracker(maxRetries)
}
