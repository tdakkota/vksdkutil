package zapvk

import (
	"fmt"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	sdkutil "github.com/tdakkota/vksdkutil/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is middleware which logs VK API request info.
func Log(l *zap.Logger, lvl zapcore.Level, skipResponse bool) func(handler sdkutil.Handler) sdkutil.Handler {
	return func(handler sdkutil.Handler) sdkutil.Handler {
		return func(method string, params ...api.Params) (api.Response, error) {
			start := time.Now()
			r, err := handler(method, params...)

			if e := l.Check(lvl, "VK request"); e != nil {
				var lastField zap.Field
				if skipResponse {
					lastField = zap.Skip()
				} else {
					lastField = zap.Inline(zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
						if len(r.Response) > 0 {
							if err := encoder.AddReflected("response", r.Response); err != nil {
								return fmt.Errorf("encode response: %w", err)
							}
						}

						if r.Error.Code != 0 {
							if err := encoder.AddReflected("error", r.Error); err != nil {
								return fmt.Errorf("encode error: %w", err)
							}
						}

						if len(r.ExecuteErrors) > 0 {
							if err := encoder.AddReflected("execute_errors", r.ExecuteErrors); err != nil {
								return fmt.Errorf("encode execute_errors: %w", err)
							}
						}
						return nil
					}))
				}

				e.Write(
					zap.String("method", method),
					zap.Object("params", zapParams(params)),
					zap.Duration("took", time.Since(start)),
					zap.Error(err),
					lastField,
				)
			}

			return r, err
		}
	}
}

type zapParams []api.Params

func (z zapParams) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	for _, params := range z {
		for k, v := range params {
			if k == "access_token" {
				continue
			}

			if err := encoder.AddReflected(k, v); err != nil {
				return fmt.Errorf("encode %q: %w", k, err)
			}
		}
	}
	return nil
}
