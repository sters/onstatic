package conf

import (
	"os"
	"strconv"
	"strings"
)

// Variables of configure
var Variables c

type c struct {
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
	Variables = c{
		HTTPPort:                          getenv("ONSTATIC_HTTP_PORT", "18888"),
		Salt:                              getenv("ONSTATIC_SALT", "saltsaltsalt"),
		SSHKeySize:                        getenvInt("ONSTATIC_SSH_KEY_SIZE", 4096),
		SSHKeyFilename:                    getenv("ONSTATIC_SSH_KEY_FILENAME", "id_rsa"),
		SSHPubKeyFilename:                 getenv("ONSTATIC_SSH_PUB_KEY_FILENAME", "id_rsa.pub"),
		RepositoriesDirectory:             getenv("ONSTATIC_REPOSITORIES_DIRECTORY", "repositories/"),
		KeyDirectoryRelatedFromRepository: getenv("ONSTATIC_KEY_DIRECTORY_RELATED_FROM_REPOSITORY", "."),
		HTTPHeaderKey:                     getenv("ONSTATIC_HTTP_HEADER_KEY", "onstaticonstaticonstatic"),
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
	n, _ := strconv.Atoi(getenv(k, ""))
	if n == 0 {
		return d
	}
	return n
}
