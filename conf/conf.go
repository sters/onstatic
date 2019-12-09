package conf

import (
	"os"
	"strconv"
	"strings"
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
}

// Init configuration variales
func Init() {
	const prefix = "ONSTATIC_"
	Variables = c{
		CGIMode:                           getenvBool(prefix+"CGI_MODE", false),
		HTTPPort:                          getenv(prefix+"HTTP_PORT", "18888"),
		Salt:                              getenv(prefix+"SALT", "saltsaltsalt"),
		SSHKeySize:                        getenvInt(prefix+"SSH_KEY_SIZE", 4096),
		SSHKeyFilename:                    getenv(prefix+"SSH_KEY_FILENAME", "id_rsa"),
		SSHPubKeyFilename:                 getenv(prefix+"SSH_PUB_KEY_FILENAME", "id_rsa.pub"),
		RepositoriesDirectory:             getenv(prefix+"REPOSITORIES_DIRECTORY", "repositories/"),
		KeyDirectoryRelatedFromRepository: getenv(prefix+"KEY_DIRECTORY_RELATED_FROM_REPOSITORY", "."),
		HTTPHeaderKey:                     getenv(prefix+"HTTP_HEADER_KEY", "onstaticonstaticonstatic"),
	}
}

func getenv(k string, d string) string {
	s := strings.TrimSpace(os.Getenv(k))
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
