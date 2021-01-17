package conf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/morikuni/failure"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Variables of configure
var Variables c

type c struct {
	CGIMode                           bool
	HTTPPort                          string
	Salt                              string
	SSHKeySize                        int
	SSHKeyFilename                    string
	SSHPubKeyFilename                 string
	RepositoriesDirectory             string
	KeyDirectoryRelatedFromRepository string
	HTTPHeaderKey                     string
	logger                            *zap.Logger
}

const prefix = "ONSTATIC_"

// Init configuration variales
func Init() {
	Variables = c{
		CGIMode:                           getenvBool("CGI_MODE", false),
		HTTPPort:                          getenv("HTTP_PORT", "18888"),
		Salt:                              getenv("SALT", "saltsaltsalt"),
		SSHKeySize:                        getenvInt("SSH_KEY_SIZE", 4096),
		SSHKeyFilename:                    getenv("SSH_KEY_FILENAME", "id_rsa"),
		SSHPubKeyFilename:                 getenv("SSH_PUB_KEY_FILENAME", "id_rsa.pub"),
		RepositoriesDirectory:             getenv("REPOSITORIES_DIRECTORY", "repositories/"),
		KeyDirectoryRelatedFromRepository: getenv("KEY_DIRECTORY_RELATED_FROM_REPOSITORY", "."),
		HTTPHeaderKey:                     getenv("HTTP_HEADER_KEY", "onstaticonstaticonstatic"),
		logger: logger(
			zapcore.InfoLevel,
			getenv("STDLOG_OUTPUT_PATH", "stdout"), // "/var/log/onstatic/stdout.log"),
			getenv("ERRLOG_OUTPUT_PATH", "stderr"), // "/var/log/onstatic/stderr.log"),
		),
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func getenv(k string, d string) string {
	s := strings.TrimSpace(os.Getenv(prefix + k))
	if s == "" {
		return d
	}
	return s
}

func getenvInt(k string, d int) int {
	n, e := strconv.Atoi(getenv(k, ""))
	if e != nil {
		return d
	}
	return n
}

func getenvBool(k string, d bool) bool {
	n, e := strconv.ParseBool((getenv(k, "")))
	if e != nil {
		return d
	}
	return n
}

func logger(logLevel zapcore.Level, logOutputPath string, logErrorPath string) *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(logLevel)
	zapConfig.DisableStacktrace = true
	zapConfig.OutputPaths = []string{logOutputPath}
	zapConfig.ErrorOutputPaths = []string{logErrorPath}

	l, err := zapConfig.Build()
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "no such file or directory") {
			if err := createLogFileIfEmpty(logOutputPath); err != nil {
				log.Fatal(err)
			}
			if err := createLogFileIfEmpty(logErrorPath); err != nil {
				log.Fatal(err)
			}

			return logger(logLevel, logOutputPath, logErrorPath)
		}

		log.Fatalf("failed to initialize logger: %+v", err)
	}

	zap.ReplaceGlobals(l)

	return l
}

func createLogFileIfEmpty(p string) error {
	if _, err := os.Stat(p); err == nil {
		return nil
	}

	dir := filepath.Dir(p)
	if _, err := os.Stat(p); err != nil {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return failure.Wrap(err)
		}
	}

	f, err := os.OpenFile(p, os.O_CREATE, os.ModePerm)
	f.Close()
	if err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func getCurrentGoFilePath() string {
	_, file, _, _ := runtime.Caller(1)
	p, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	return p
}
