package intercept

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[30;47m"
	yellow  = "\033[30;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
	proxy   = "\033[97;40m"
)

type LogFormatterParams struct {
	Request *http.Request
	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]any
}

func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// Logger 重写了GIN框架的Logger，主要是为了自定义日志格式
func Logger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	var param = LogFormatterParams{
		Request: c.Request,
		isTerm:  true,
		Keys:    c.Keys,
	}

	param.TimeStamp = time.Now()
	param.Latency = param.TimeStamp.Sub(start)
	param.ClientIP = c.ClientIP()
	param.Method = c.Request.Method
	param.StatusCode = c.Writer.Status()
	if raw != "" {
		path = path + "?" + raw
	}

	param.Path = path

	var statusColor, methodColor, resetColor string
	statusColor = param.StatusCodeColor()
	methodColor = param.MethodColor()
	resetColor = param.ResetColor()

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	var log = fmt.Sprintf("[WPServer] %v |%s %3d %s| %3v | %s %s => %s => %s %s |%s %-7s %s %#v",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		proxy, param.ClientIP,
		c.Request.Host,
		c.Request.URL.Host, reset,
		methodColor, param.Method, resetColor,
		param.Path,
	)
	fmt.Println(log)
	c.Next()
}
