//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package resources

import (
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dellemc-symphony/workflow-cli/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var UseEndpointFile bool

//ParseEndpointsFile parses the file for the endpoints
func ParseEndpointsFile(useEndpointFile bool) map[string]models.Endpoint {
	services := []string{"rackhd", "hostbmc", "vcenter", "scaleiogateway"}
	endpoints := make(map[string]models.Endpoint, len(services))

	if !useEndpointFile {
		return endpoints
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName("endpoint")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(dir)

	err = viper.ReadInConfig()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Warnf(`Config file "endpoint.yaml" not found.`)
		} else {
			log.Warnf("Invalid endpoint.yaml: %s", err)
		}

		log.Warnf("Will prompt user for endpoints.")
		return endpoints
	}

	for _, service := range services {

		entry := models.Endpoint{}

		endpoint := viper.GetStringSlice(service + ".endpoint")
		if len(endpoint) == 1 && ValidateEndpoint(endpoint[0]) {
			entry.EndpointURL = endpoint[0]
		} else {
			entry.EndpointURL = ""
		}

		username := viper.GetStringSlice(service + ".username")
		if len(username) == 1 {
			entry.Username = username[0]
		} else {
			entry.Username = ""
		}

		password := viper.GetStringSlice(service + ".password")
		if len(password) == 1 {
			entry.Password = password[0]
		} else {
			entry.Password = ""
		}

		endpoints[service] = entry

	}

	return endpoints
}

// ValidateEndpoint is a helper to parse endpoints entered
func ValidateEndpoint(endpoint string) bool {
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Warnf("Invalid URL: '%s' (%s)", endpoint, err)
		return false
	}

	// If not (scheme is either http or https)
	if !(u.Scheme == "http" || u.Scheme == "https") {
		log.Warnf("Endpoint must begin with http:// or https://")
		return false

	} else if u.Host == "" {
		log.Warnf("Host must be a valid IP or Hostname")
		return false

	} else if u.Port() != "" {
		portNumber, err := strconv.Atoi(u.Port())
		if err != nil {
			log.Warnf("Could not parse port number for %s (%s)", endpoint, err)
			return false
		}

		if portNumber < 0 || portNumber > 65536 {
			log.Warnf("Endpoint must include a valid Port Number (0-65535)")
			return false
		}
	}

	// Endpoint is valid
	return true
}
