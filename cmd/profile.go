package cmd

import (
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/rs/zerolog/log"
)

var (
	tf             *os.File
	memProfileFile *os.File
)

func StartProfile(traceFile string) {
	if len(traceFile) > 0 {
		f, err := os.Create(traceFile)
		if err != nil {
			panic(err)
		}
		if err := trace.Start(f); err != nil {
			panic(err)
		}
		tf = f
		log.Debug().Msg("trace enabled")
	}

	if len(FlagCPUProfile) > 0 {
		f, err := os.Create(FlagCPUProfile)
		if err != nil {
			panic(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		log.Debug().Msg("cpuprofile enabled")
	}

	if len(FlagMemProfile) > 0 {
		f, err := os.Create(FlagMemProfile)
		if err != nil {
			panic(err)
		}
		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}
		memProfileFile = f
		log.Debug().Msg("memprofile enabled")
	}
}

func CloseProfile(traceFile string) {
	if len(FlagCPUProfile) > 0 {
		pprof.StopCPUProfile()
		log.Debug().Msg("cpu profile closed")
	}

	if len(FlagMemProfile) > 0 {
		if err := memProfileFile.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close mem profile file")
		}

		log.Debug().Msg("mem profile closed")
	}

	if len(traceFile) > 0 {
		trace.Stop()
		if err := tf.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close trace file")
		}

		log.Debug().Msg("trace closed")
	}
}
