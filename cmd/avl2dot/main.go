package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
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

var (
	flagSkipDot bool
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
	Use:   "avl2dot [<key> ...]",
	Short: "print dot graph from avl tree",
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

			log.Debug().
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

		log.Debug().
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
		if len(input) < 1 {
			_ = c.Help()
			os.Exit(1)
		}

		var b bytes.Buffer

		started := time.Now()

		tg := avl.NewTreeGenerator()
		_ = tg.SetLogger(log)

		for _, k := range input {
			if len(k) < 1 {
				log.Debug().Msg("skip empty key")
				continue
			}

			log.Debug().Str("key", k).Msg("trying to add node")
			_, err := tg.Add(cmd.NewMutableNode([]byte(k)))
			if err != nil {
				c.Println("Error: failed to add node;", err.Error())
				os.Exit(1)
			}
		}

		elapsed := time.Since(started)
		log.Info().
			Int("keys", len(input)).
			Str("elapsed", elapsed.String()).
			Dur("elapsed_ms", elapsed).
			Msg("tree generated")

		if !flagSkipDot {
			tr, err := tg.Tree()
			if err != nil {
				c.Println("Error: failed to make tree;", err.Error())
				os.Exit(1)
			}

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
	rootCmd.PersistentFlags().BoolVar(&flagSkipDot, "skip-dot", flagSkipDot, "don't print dot")

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println("Error:", err.Error())
		os.Exit(1)
	}
}
