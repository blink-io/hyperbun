package verbose

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/uptrace/bun"
)

var bufPool = &sync.Pool{
	New: func() any { return &bytes.Buffer{} },
}

type QueryHookOptions struct {
	// Show time taken by the query.
	ShowTimeTaken bool

	// If true, logs are shown as plaintext (no color).
	NoColor bool

	// Verbose query interpolation, which shows the query before and after
	// interpolating query arguments. The logged query is interpolated by
	// default, InterpolateVerbose only controls whether the query before
	// interpolation is shown. To disable query interpolation entirely, look at
	// HideArgs.
	InterpolateVerbose bool

	// Explicitly hides arguments when logging the query (only the query
	// placeholders will be shown).
	HideArgs bool
}

type hook struct {
	logger *log.Logger
	opts   QueryHookOptions
}

func (h *hook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *hook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	h.logQuery(ctx, event)
}

// New returns new instance
func New(w io.Writer, prefix string, flag int, opts QueryHookOptions) bun.QueryHook {
	return &hook{
		logger: log.New(w, prefix, flag),
		opts:   opts,
	}
}

func (h *hook) logQuery(ctx context.Context, event *bun.QueryEvent) {
	var reset, red, green, blue, purple string
	envNoColor, _ := strconv.ParseBool(os.Getenv("NO_COLOR"))
	if !h.opts.NoColor && !envNoColor {
		reset = colorReset
		red = colorRed
		green = colorGreen
		blue = colorBlue
		purple = colorPurple
	}
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	if event.Err == nil {
		buf.WriteString(green + "[OK]" + reset)
	} else {
		buf.WriteString(red + "[FAIL]" + reset)
	}
	buf.WriteString("[" + event.Operation() + "]")
	if h.opts.HideArgs {
		buf.WriteString(" " + event.Query + ";")
	} else if !h.opts.InterpolateVerbose {
		if event.Err != nil {
			buf.WriteString(" " + event.Query + ";")
			if len(event.QueryArgs) > 0 {
				buf.WriteString(" [")
			}
			for i := 0; i < len(event.QueryArgs); i++ {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(fmt.Sprintf("%#v", event.QueryArgs[i]))
			}
			if len(event.QueryArgs) > 0 {
				buf.WriteString("]")
			}
		} else {
			query, err := Sprintf(event.DB.Dialect(), event.Query, event.QueryArgs)
			if err != nil {
				query += " " + err.Error()
			}
			buf.WriteString(" " + query + ";")
		}
	}
	if event.Err != nil {
		errStr := event.Err.Error()
		if i := strings.IndexByte(errStr, '\n'); i < 0 {
			buf.WriteString(blue + " err" + reset + "={" + event.Err.Error() + "}")
		}
	}
	if h.opts.ShowTimeTaken {
		timeTakenInMs := time.Until(event.StartTime).Milliseconds()
		buf.WriteString(blue + " timeTaken" + reset + "=" + strconv.FormatInt(timeTakenInMs, 10))
	}
	//if queryStats.RowCount.Valid {
	//	buf.WriteString(blue + " rowCount" + reset + "=" + strconv.FormatInt(queryStats.RowCount.Int64, 10))
	//}
	if rowAffected, err := event.Result.RowsAffected(); err == nil {
		buf.WriteString(blue + " rowsAffected" + reset + "=" + strconv.FormatInt(rowAffected, 10))
	}
	if lastInsertId, err := event.Result.LastInsertId(); err == nil {
		buf.WriteString(blue + " lastInsertId" + reset + "=" + strconv.FormatInt(lastInsertId, 10))
	}
	//if h.opts.ShowCaller {
	//	buf.WriteString(blue + " caller" + reset + "=" + queryStats.CallerFile + ":" + strconv.Itoa(queryStats.CallerLine) + ":" + filepath.Base(queryStats.CallerFunction))
	//}
	if !h.opts.HideArgs && h.opts.InterpolateVerbose {
		buf.WriteString("\n" + purple + "----[ Executing query ]----" + reset)
		buf.WriteString("\n" + event.Query + "; " + fmt.Sprintf("%#v", event.QueryArgs))
		buf.WriteString("\n" + purple + "----[ with bind values ]----" + reset)
		query, err := Sprintf(event.DB.Dialect(), event.Query, event.QueryArgs)
		query += ";"
		if err != nil {
			query += " " + err.Error()
		}
		buf.WriteString("\n" + query)
	}
	//if h.opts.ShowResults > 0 && event.Err == nil {
	//	buf.WriteString("\n" + purple + "----[ Fetched result ]----" + reset)
	//	buf.WriteString(queryStats.Results)
	//	if queryStats.RowCount.Int64 > int64(l.config.ShowResults) {
	//		buf.WriteString("\n...\n(Fetched " + strconv.FormatInt(queryStats.RowCount.Int64, 10) + " rows)")
	//	}
	//}
	if buf.Len() > 0 {
		h.logger.Println(buf.String())
	}
}
