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
		HTTPPort:                          getenv("STATICMAN_HTTP_PORT", "18888"),
		Salt:                              getenv("STATICMAN_SALT", "saltsaltsalt"),
		SSHKeySize:                        getenvInt("STATICMAN_SSH_KEY_SIZE", 4096),
		SSHKeyFilename:                    getenv("STATICMAN_SSH_KEY_FILENAME", "id_rsa"),
		SSHPubKeyFilename:                 getenv("STATICMAN_SSH_PUB_KEY_FILENAME", "id_rsa.pub"),
		RepositoriesDirectory:             getenv("STATICMAN_REPOSITORIES_DIRECTORY", "repositories/"),
		KeyDirectoryRelatedFromRepository: getenv("STATICMAN_KEY_DIRECTORY_RELATED_FROM_REPOSITORY", "."),
		HTTPHeaderKey:                     getenv("STATICMAN_HTTP_HEADER_KEY", "staticmanstaticmanstaticman"),
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
