// internal/commands/result/result.go
package result

import "fmt"

type CommandResult struct {
	Success  bool
	Error    error
	Data     interface{}
	Warnings []string
	Messages []string
}

func NewSuccess(data interface{}, messages ...string) CommandResult {
	return CommandResult{
		Success:  true,
		Data:     data,
		Messages: messages,
	}
}

func NewError(err error, warnings ...string) CommandResult {
	return CommandResult{
		Success:  false,
		Error:    err,
		Warnings: warnings,
	}
}

func (r CommandResult) WithWarnings(warnings ...string) CommandResult {
	r.Warnings = append(r.Warnings, warnings...)
	return r
}

func (r CommandResult) WithMessages(messages ...string) CommandResult {
	r.Messages = append(r.Messages, messages...)
	return r
}

// Helper methods for common result patterns
func NotFound(itemType, identifier string) CommandResult {
	return NewError(fmt.Errorf("%s not found: %s", itemType, identifier))
}

func ValidationFailed(err error, warnings ...string) CommandResult {
	return NewError(fmt.Errorf("validation failed: %v", err), warnings...)
}

func StorageError(err error) CommandResult {
	return NewError(fmt.Errorf("storage operation failed: %v", err))
}
