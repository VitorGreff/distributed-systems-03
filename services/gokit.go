package services

import "context"

type Service interface {
	PublishValidationRequest(context.Context, string, uint64) error
	ConsumeValidationResponse(context.Context) error

	PublishGenerationRequest(context.Context, uint64) string
	ConsumeGenerationRequest(context.Context) error
}
