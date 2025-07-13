package logging

import (
	"errors"
	"log/slog"
	"time"
)

const _secureMask = "[SECURE]"

func _noopRequestIDFn() string {
	return ""
}

func defaultRequestIDFn() string {
	return _noopRequestIDFn()
}

func defaultTimeFn() time.Time {
	return time.Now()
}

func defaultConfig() config {
	return config{
		DebugMode:           false,
		LogRequest:          true,
		LogResponse:         true,
		LogFullURL:          true,
		SensitiveHeaders:    nil,
		SensitiveJSONFields: nil,
		ReplaceString:       _secureMask,
		DefaultLevel:        slog.LevelDebug,
		ErrorLevel:          slog.LevelWarn,
		Logger:              slog.Default(),
		RequestIDFn:         defaultRequestIDFn,
		TimeFn:              defaultTimeFn,
		JSONSetter:          nil,
		JSONGetter:          nil,
	}
}

type config struct {
	DebugMode           bool
	LogRequest          bool       // Flag indicating whether the body of the request is needed to log.
	LogResponse         bool       // Flag indicating whether the body of the response is needed to log.
	LogFullURL          bool       // Flag indicating whether the full URL of the request is needed to log.
	SensitiveHeaders    []string   // Headers that need to be hidden (for example, "authorization").
	SensitiveJSONFields []string   // Fields in json (you can indicate the ways, for example, "User.password")
	ReplaceString       string     // Line for replacement (by default "[SECURE]")
	DefaultLevel        slog.Level // Default log level. If nothing is specified, the default level is debug.
	ErrorLevel          slog.Level // Error log level. If nothing is specified, the default level is warn.
	Logger              *slog.Logger
	RequestIDFn         func() string
	TimeFn              func() time.Time
	JSONSetter          JSONSetter
	JSONGetter          JSONGetter
}

func validateConfig(c config) error {
	if c.JSONGetter == nil {
		return errors.New("parameter 'JSONGetter' is not set")
	}

	if c.JSONSetter == nil {
		return errors.New("parameter 'JSONSetter' is not set")
	}

	if c.Logger == nil {
		return errors.New("parameter 'Logger' is not set")
	}

	if c.RequestIDFn == nil {
		return errors.New("parameter 'RequestIDFn' is not set")
	}

	if c.TimeFn == nil {
		return errors.New("parameter 'TimeFn' is not set")
	}

	return nil
}

type Option func(*config)

func WithDebugMode() Option {
	return func(c *config) {
		c.DebugMode = true
	}
}

func WithOmitRequest() Option {
	return func(c *config) {
		c.LogRequest = false
	}
}

func WithOmitResponse() Option {
	return func(c *config) {
		c.LogResponse = false
	}
}

func WithOmitFullURL() Option {
	return func(c *config) {
		c.LogFullURL = false
	}
}

func WithSensitiveHeaders(headers []string) Option {
	return func(c *config) {
		c.SensitiveHeaders = headers
	}
}

func WithSensitiveJSONFieldsHeaders(fields []string) Option {
	return func(c *config) {
		c.SensitiveJSONFields = fields
	}
}

func WithDefaultLogLevel(level slog.Level) Option {
	return func(c *config) {
		c.DefaultLevel = level
	}
}

func WithErrorLogLevel(level slog.Level) Option {
	return func(c *config) {
		c.ErrorLevel = level
	}
}
