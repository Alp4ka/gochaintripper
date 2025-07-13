package logging

import (
	"bytes"
	"context"
	"github.com/Alp4ka/gochaintripper"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Interceptor implements gochaintripper.Interceptor.
type Interceptor struct {
	cfg config
}

func New(options ...Option) (*Interceptor, error) {
	cfg := defaultConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Interceptor{cfg: cfg}, nil
}

func (i *Interceptor) RoundTripWithTransport(rt http.RoundTripper, req *http.Request) (resp *http.Response, err error) {
	ctx := req.Context()
	logger := i.cfg.Logger.With(
		slog.String("request_id", i.cfg.RequestIDFn()),
		slog.String("http_method", req.Method),
		slog.String("url", i.maskURL(req)),
	)

	if i.cfg.LogRequest {
		err = i.logRequest(ctx, logger, req)
		if err != nil {
			logger.LogAttrs(ctx, i.cfg.DefaultLevel, "Failed to read request body!", slog.Any("err", err))
			return nil, err
		}
	}

	startTime := i.cfg.TimeFn()
	defer func() {
		logger = logger.With(slog.String("duration", time.Since(startTime).String()))

		if err != nil {
			logger.LogAttrs(ctx, i.cfg.ErrorLevel, "HTTP Request failed!", slog.Any("err", err))
			return
		}

		if i.cfg.LogResponse {
			err = i.logResponse(ctx, logger, resp)
			if err != nil {
				logger.LogAttrs(ctx, i.cfg.ErrorLevel, "Failed to read response body!", slog.Any("err", err))
			}
		}
	}()

	return rt.RoundTrip(req)
}

func (i *Interceptor) logRequest(ctx context.Context, logger *slog.Logger, req *http.Request) error {
	maskedHeaders := i.maskSensitiveHeaders(req.Header)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	logger.LogAttrs(
		ctx,
		i.cfg.DefaultLevel,
		"Sending request!",
		slog.String("req_body", i.maskSensitiveJSONFields(string(body))),
		slog.Any("req_headers", maskedHeaders),
	)

	return nil
}

func (i *Interceptor) logResponse(ctx context.Context, logger *slog.Logger, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	logger.LogAttrs(
		ctx,
		i.defineLogLevelByStatusCode(resp.StatusCode),
		"Response received!",
		slog.String("resp_body", i.maskSensitiveJSONFields(string(body))),
		slog.Int("status_code", resp.StatusCode),
	)

	return nil
}

func (i *Interceptor) maskSensitiveJSONFields(body string) string {
	if len(i.cfg.SensitiveJSONFields) == 0 || len(body) == 0 || i.cfg.DebugMode || !i.cfg.JSONGetter.IsJSON(body) {
		return body
	}

	for _, field := range i.cfg.SensitiveJSONFields {
		if i.cfg.JSONGetter.Exists(body, field) {
			body, _ = i.cfg.JSONSetter.SetValue(body, field, i.cfg.ReplaceString)
		}
	}

	return body
}

func (i *Interceptor) maskSensitiveHeaders(headers http.Header) http.Header {
	maskedHeaders := headers.Clone()

	if i.cfg.DebugMode {
		return maskedHeaders
	}

	for _, header := range i.cfg.SensitiveHeaders {
		if maskedHeaders.Get(header) != "" {
			maskedHeaders.Set(header, i.cfg.ReplaceString)
		}
	}

	return maskedHeaders
}

func (i *Interceptor) maskURL(req *http.Request) string {
	urlVal := req.URL.Path
	if i.cfg.LogFullURL {
		urlVal = req.URL.String()
	}

	return urlVal
}

func (i *Interceptor) defineLogLevelByStatusCode(statusCode int) slog.Level {
	logLevel := i.cfg.DefaultLevel
	if statusCode/100 != 2 {
		logLevel = i.cfg.ErrorLevel
	}

	return logLevel
}

var _ gochaintripper.Interceptor = (*Interceptor)(nil)
