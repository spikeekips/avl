package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

var (
	FlagLogLevel      LogLevel  = LogLevel{Lvl: zerolog.ErrorLevel}
	FlagLogFormat     LogFormat = LogFormat{f: "terminal"}
	FlagLogOut        string
	FlagCPUProfile    string
	FlagMemProfile    string
	FlagTrace         string
	FlagExitAfter     time.Duration
	FlagNumberOfNodes uint = 3
	FlagQuiet         bool
	FlagQueries       []string
	FlagJSONPretty    bool
)

type LogLevel struct {
	Lvl zerolog.Level
}

func (f LogLevel) String() string {
	return f.Lvl.String()
}

func (f *LogLevel) Set(v string) error {
	lvl, err := zerolog.ParseLevel(v)
	if err != nil {
		return err
	}

	f.Lvl = lvl

	return nil
}

func (f LogLevel) Type() string {
	return "log-level"
}

type LogFormat struct {
	f string
}

func (f LogFormat) String() string {
	return f.f
}

func (f *LogFormat) Set(v string) error {
	s := strings.ToLower(v)
	switch s {
	case "json":
	case "terminal":
	default:
		return xerrors.Errorf("invalid log format: %q", v)
	}

	f.f = s

	return nil
}

func (f LogFormat) Type() string {
	return "log-format"
}

func escapeFlagValue(v interface{}, q string) string {
	if len(q) < 1 {
		return fmt.Sprintf("%v", v)
	}

	return q + strings.Replace(fmt.Sprintf("%v", v), "'", "\\"+q, -1) + q
}

func PrintFlagsJSON(cmd *cobra.Command) json.RawMessage {
	out := map[string]interface{}{}

	cmd.Flags().VisitAll(func(pf *pflag.Flag) {
		if pf.Name == "help" {
			return
		}

		out[fmt.Sprintf("--%s", pf.Name)] = map[string]interface{}{
			"default": escapeFlagValue(pf.DefValue, ""),
			"value":   escapeFlagValue(pf.Value, ""),
		}
	})

	b, _ := json.Marshal(out)

	return b
}
