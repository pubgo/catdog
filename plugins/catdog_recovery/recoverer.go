package catdog_recovery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pubgo/xlog"

	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/plugins/catdog_tracing/tracing"
)

func (p *Plugin) handlerWrap() error {
	return catdog_app.WrapHandler(func(handlerFunc server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			t := time.Now()
			defer func() {
				fieldsMap := make(map[string]interface{})
				fieldsMap["service"] = fmt.Sprintf("%v.%v", req.Service(), req.Endpoint())
				fieldsMap["cost"] = int(time.Since(t).Seconds() * 1000) // ms
				fieldsMap["unique_id"] = tracing.GetRequestIDFromContext(ctx)
				fieldsMap["receive_time"] = t.Format(time.RFC3339Nano)

				var params string
				if b := req.Body(); b != nil {
					if v, ok := b.(string); ok {
						params = v
					} else if v, ok := b.([]byte); ok {
						params = string(v)
					} else {
						body, err := json.Marshal(b)
						if err != nil {
							//p.log.ErrorF("handler error, body marshal, ", err)
						}
						params = string(body)
					}
				}
				fieldsMap["params"] = params

				msg, err := json.Marshal(fieldsMap)
				if err != nil {
					xlog.Errorf("handler error, msg marshal, %s", err)
				}

				var fields []xlog.Field
				for k, v := range fieldsMap {
					fields = append(fields, xlog.Any(k, v))
				}

				xlog.Info(string(msg))

				if err := recover(); err != nil {
					//p.log.ErrorF("handler error: %v", err)
				}
			}()

			return handlerFunc(ctx, req, rsp)
		}
	})
}

func (p *Plugin) clientWrap() error {
	return catdog_app.WrapCall(func(callFunc client.CallFunc) client.CallFunc {
		return func(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) error {
			defer func() {
				if err := recover(); err != nil {
					//p.log.ErrorF("Call service error: %v", err)
				}
			}()

			return callFunc(ctx, addr, req, rsp, opts)
		}
	})
}
