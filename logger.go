package elasticsearch

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/illidaris/core"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	logger *zap.Logger // log core without context
)

type Logger struct {
	requestBodyEnabled  bool
	responseBodyEnabled bool
}

func NewLogger(hasReq, hasRes bool) estransport.Logger {
	logger = zap.L().WithOptions(zap.AddCallerSkip(6))
	return &Logger{
		requestBodyEnabled:  hasReq,
		responseBodyEnabled: hasRes,
	}
}

func (l *Logger) LogRoundTrip(req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	ctx := req.Context()
	fields := FieldsFromCtx(ctx)
	fields = append(fields, zap.Int64(core.Duration.String(), dur.Milliseconds()))
	url := req.URL.String()
	// query
	query := req.URL.RawQuery
	// method
	method := req.Method
	// request
	var request string
	if l.RequestBodyEnabled() && req != nil && req.Body != nil && req.Body != http.NoBody {
		var buf bytes.Buffer
		if req.GetBody != nil {
			b, _ := req.GetBody()
			buf.ReadFrom(b)
		} else {
			buf.ReadFrom(req.Body)
		}
		request = buf.String()
	}
	// code
	code := res.StatusCode
	// Response
	var response string
	if l.ResponseBodyEnabled() && res != nil && res.Body != nil && res.Body != http.NoBody {
		defer res.Body.Close()
		var buf bytes.Buffer
		buf.ReadFrom(res.Body)
		response = buf.String()
	}
	resultMsg := fmt.Sprintf("[%s-%d-%dms]%s,%s,request:%s,response:%s,err:%s", method, code, dur.Milliseconds(), url, query, request, response, err)
	// Error
	if err != nil {
		logger.Error(resultMsg, fields...)
	}
	logger.Info(resultMsg, fields...)
	return err
}

// RequestBodyEnabled makes the client pass a copy of request body to the logger.
func (l *Logger) RequestBodyEnabled() bool {
	return l.requestBodyEnabled
}

// ResponseBodyEnabled makes the client pass a copy of response body to the logger.
func (l *Logger) ResponseBodyEnabled() bool {
	return l.responseBodyEnabled
}
