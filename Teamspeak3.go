package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/multiplay/go-ts3"
	"gopkg.in/ini.v1"
)

type serverQueryConfig struct {
	ServerAddress string
	Username      string
	Password      string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	queryConfig, err := readConfig()
	if err != nil {
		return err
	}

	// Print Check_MK section header
	fmt.Println("<<<Teamspeak3>>>")
	fmt.Println("ConfigError: No")

	// Establish connection to Teamspeak3 server query
	c, err := ts3.NewClient(queryConfig.ServerAddress)

	// Determine if we can actually reach the server's query console
	// In case of error we want to exit with zero-code because the check application itself ran correctly
	if err != nil {
		fmt.Println("QueryPortReachable: No")
		return err
	}
	fmt.Println("QueryPortReachable: Yes")

	// Make sure the query connection will be closed when application terminates
	defer c.Close()

	// Try to authenticate with Teamspeak3 server query
	if err := c.Login(queryConfig.Username, queryConfig.Password); err != nil {
		fmt.Println("AuthSuccess: No")
		return err
	}
	fmt.Println("AuthSuccess: Yes")

	// Try to get server's current version
	v, err := c.Version()
	if err != nil {
		fmt.Println("Version: None")
		fmt.Println("Platform: None")
		fmt.Println("Build: None")
		return err
	}
	fmt.Println("Version:", v.Version)
	fmt.Println("Platform:", v.Platform)
	fmt.Println("Build:", v.Build)

	// Iterate through list of virtual servers
	l, err := c.Server.List()
	if err != nil {
		return err
	}
	for _, server := range l {
		if err := c.Use(server.ID); err != nil {
			continue
		}

		var serverAutoStart string = "no"
		var trafficIngressBytesTotal uint64 = 0
		var trafficEgressBytesTotal uint64 = 0

		// Convert boolean value to string
		if server.AutoStart {
			serverAutoStart = "yes"
		}

		// When the server is stopped this method exits with an error
		if connInfo, err := c.Server.ServerConnectionInfo(); err == nil {
			trafficIngressBytesTotal = connInfo.BytesReceivedTotal
			trafficEgressBytesTotal = connInfo.BytesSentTotal
		}

		// Scheme: "VirtualServer: ($PORT $STATUS $ONLINE_CLIENTS $MAX_CLIENTS $CURRENT_CHANNELS $AUTO_START $BANDWIDTH_INGRESS_TOTAL $BANDWIDTH_EGRESS_TOTAL )"
		fmt.Printf(
			"VirtualServer: (%d %s %d %d %d %s %d %d)\n",
			server.Port,
			server.Status,
			server.ClientsOnline,
			server.MaxClients,
			server.ChannelsOnline,
			serverAutoStart,
			trafficIngressBytesTotal,
			trafficEgressBytesTotal,
		)
	}
	return nil
}

func readConfig() (*serverQueryConfig, error) {
	// Load the user configuration and process the required sections and keys
	// Exit application with non-zero code when configuration was not read successful
	// Also display a Check_MK Agent-compatible error code to process on monitoring server
	cpath := filepath.Join(os.Getenv("MK_CONFDIR"), "teamspeak3.cfg")
	cfg, err := ini.Load(cpath)
	if err != nil {
		fmt.Println("<<<Teamspeak3>>>")
		fmt.Println("ConfigError: Yes, 1")
		return nil, fmt.Errorf("could not load config: %s", cpath)
	}
	sect, err := cfg.GetSection("serverquery")
	if err != nil {
		fmt.Println("<<<Teamspeak3>>>")
		fmt.Println("ConfigError: Yes, 2")
		return nil, errors.New("config section missing: serverquery")
	}

	conf_address, err := sect.GetKey("address")
	if err != nil {
		fmt.Println("<<<Teamspeak3>>>")
		fmt.Println("ConfigError: Yes, 3")
		return nil, errors.New("config value missing: address")
	}

	conf_user, err := sect.GetKey("user")
	if err != nil {
		fmt.Println("<<<Teamspeak3>>>")
		fmt.Println("ConfigError: Yes, 4")
		return nil, errors.New("config value missing: user")
	}

	conf_password, err := sect.GetKey("password")
	if err != nil {
		fmt.Println("<<<Teamspeak3>>>")
		fmt.Println("ConfigError: Yes, 5")
		return nil, errors.New("config value missing: password")
	}

	// Build ServerQueryConfig struct and return pointer
	return &serverQueryConfig{
		ServerAddress: conf_address.String(),
		Username:      conf_user.String(),
		Password:      conf_password.String(),
	}, nil
}
