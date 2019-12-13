package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/spikeekips/avl"
	"github.com/spikeekips/avl/cmd"
)

var (
	sigc            chan os.Signal
	exitHooks       []func()
	exitCode        int = 0
	log             zerolog.Logger
	logOutput       io.Writer
	whitespaceSplit *regexp.Regexp = regexp.MustCompile(`\s+`)
)

func init() {
	zerolog.TimestampFieldName = "t"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.MessageFieldName = "m"
	zerolog.DisableSampling(true)

	exitHooks = append(exitHooks, func() {
		cmd.CloseProfile(cmd.FlagTrace)
	})
}

var rootCmd = &cobra.Command{
	Use:   "avl-print [<key> ...]",
	Short: "avl-print tree",
	PersistentPreRun: func(c *cobra.Command, args []string) {
		{ // logging
			if cmd.FlagLogOut == "null" {
				logOutput = nil
			} else if len(cmd.FlagLogOut) < 1 {
				if cmd.FlagLogFormat.String() == "terminal" {
					f := zerolog.NewConsoleWriter()
					f.NoColor = false
					logOutput = f
				} else {
					logOutput = os.Stderr
				}
			} else if f, err := os.OpenFile(
				cmd.FlagLogOut, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600,
			); err != nil {
				c.Println("Error:", err.Error())
				os.Exit(1)
			} else {
				logOutput = f
			}

			logContext := zerolog.
				New(os.Stderr).
				With().
				Timestamp()

			if cmd.FlagLogLevel.Lvl == zerolog.DebugLevel {
				logContext = logContext.
					Caller().
					Stack()
			}

			log = logContext.Logger().Level(cmd.FlagLogLevel.Lvl).Output(logOutput)
		}

		log.Debug().
			RawJSON("flags", cmd.PrintFlagsJSON(c)).
			Msg("parsed flags")

		cmd.StartProfile(cmd.FlagTrace)

		sigc = make(chan os.Signal, 1)
		signal.Notify(sigc,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		go func() {
			s := <-sigc

			for _, h := range exitHooks {
				h()
			}

			log.Info().
				Str("sig", s.String()).
				Int("exit", exitCode).
				Msg("stopped by force")

			os.Exit(exitCode)
		}()
	},
	PersistentPostRun: func(c *cobra.Command, args []string) {
		for _, h := range exitHooks {
			h()
		}

		log.Info().
			Int("exit", exitCode).
			Msg("stopped")
		os.Exit(exitCode)
	},
	Run: func(c *cobra.Command, args []string) {
		input := args
		{ // from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {

				var s string
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					s += scanner.Text() + " "
				}
				input = append(input, whitespaceSplit.Split(s, -1)...)
			}
		}
		log.Debug().Strs("input", input).Send()

		// check keys
		var maxLength int
		var keys []int64
		for _, i := range input {
			if strings.TrimSpace(i) == "" {
				continue
			}

			k, err := strconv.ParseInt(i, 10, 64)
			if err != nil {
				c.Printf("Error: node key must be integer: %q\n", i)
				os.Exit(1)
			} else if k < 0 {
				c.Printf("Error: node key must be equal or greater than 0: %q\n", i)
				os.Exit(1)
			}

			keys = append(keys, k)

			l := len(strconv.FormatInt(k, 10))
			if maxLength < l {
				maxLength = l
			}
		}
		log.Debug().Int("max-length", maxLength).Ints64("keys", keys).Send()

		var b bytes.Buffer

		started := time.Now()

		tr := avl.NewTree(avl.NewMapNodePool())
		_ = tr.SetLogger(log)

		keyFormat := fmt.Sprintf("%%0%dd", maxLength)

		for _, k := range keys {
			_, err := tr.Add(
				avl.NewBaseNode(
					[]byte(fmt.Sprintf(keyFormat, k)),
				),
			)
			if err != nil {
				c.Println("Error: failed to add node;", err.Error())
				os.Exit(1)
			}
		}
		log.Debug().Dur("elapsed", time.Since(started)).Msg("finished")

		if !cmd.FlagQuiet {
			avl.PrintDotGraph(tr, &b)

			fmt.Fprint(os.Stdout, b.String())
		}
	},
}

func main() {
	rootCmd.PersistentFlags().Var(&cmd.FlagLogLevel, "log-level", "log level: {debug error warn info crit}")
	rootCmd.PersistentFlags().StringVar(&cmd.FlagLogOut, "log", cmd.FlagLogOut, "log output directory")
	rootCmd.PersistentFlags().Var(&cmd.FlagLogFormat, "log-format", "log format: {json terminal}")
	rootCmd.PersistentFlags().StringVar(
		&cmd.FlagCPUProfile, "cpuprofile", cmd.FlagCPUProfile, "write cpu profile to file",
	)
	rootCmd.PersistentFlags().StringVar(
		&cmd.FlagMemProfile, "memprofile", cmd.FlagMemProfile, "write memory profile to file",
	)
	rootCmd.PersistentFlags().StringVar(&cmd.FlagTrace, "trace", cmd.FlagTrace, "write trace to file")
	rootCmd.PersistentFlags().BoolVar(&cmd.FlagQuiet, "quiet", cmd.FlagQuiet, "no output")

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println("Error:", err.Error())
		os.Exit(1)
	}
}
