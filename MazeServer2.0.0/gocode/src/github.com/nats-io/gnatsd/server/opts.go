// Copyright 2012-2013 Apcera Inc. All rights reserved.

package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/gnatsd/conf"
)

// Options block for gnatsd server.
type Options struct {
	Host               string        `json:"addr"`
	Port               int           `json:"port"`
	Trace              bool          `json:"-"`
	Debug              bool          `json:"-"`
	NoLog              bool          `json:"-"`
	NoSigs             bool          `json:"-"`
	Logtime            bool          `json:"-"`
	MaxConn            int           `json:"max_connections"`
	Username           string        `json:"user,omitempty"`
	Password           string        `json:"-"`
	Authorization      string        `json:"-"`
	PingInterval       time.Duration `json:"ping_interval"`
	MaxPingsOut        int           `json:"ping_max"`
	HTTPPort           int           `json:"http_port"`
	SslTimeout         float64       `json:"ssl_timeout"`
	AuthTimeout        float64       `json:"auth_timeout"`
	MaxControlLine     int           `json:"max_control_line"`
	MaxPayload         int           `json:"max_payload"`
	MaxPending         int           `json:"max_pending_size"`
	ClusterHost        string        `json:"addr"`
	ClusterPort        int           `json:"port"`
	ClusterUsername    string        `json:"-"`
	ClusterPassword    string        `json:"-"`
	ClusterAuthTimeout float64       `json:"auth_timeout"`
	ProfPort           int           `json:"-"`
	PidFile            string        `json:"-"`
	LogFile            string        `json:"-"`
	Syslog             bool          `json:"-"`
	RemoteSyslog       string        `json:"-"`
	Routes             []*url.URL    `json:"-"`
	RoutesStr          string        `json:"-"`
}

type authorization struct {
	user    string
	pass    string
	timeout float64
}

// ProcessConfigFile processes a configuration file.
// FIXME(dlc): Hacky
func ProcessConfigFile(configFile string) (*Options, error) {
	opts := &Options{}

	if configFile == "" {
		return opts, nil
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}

	m, err := conf.Parse(string(data))
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		switch strings.ToLower(k) {
		case "port":
			opts.Port = int(v.(int64))
		case "host", "net":
			opts.Host = v.(string)
		case "debug":
			opts.Debug = v.(bool)
		case "trace":
			opts.Trace = v.(bool)
		case "logtime":
			opts.Logtime = v.(bool)
		case "authorization":
			am := v.(map[string]interface{})
			auth := parseAuthorization(am)
			opts.Username = auth.user
			opts.Password = auth.pass
			opts.AuthTimeout = auth.timeout
		case "http_port", "monitor_port":
			opts.HTTPPort = int(v.(int64))
		case "cluster":
			cm := v.(map[string]interface{})
			if err := parseCluster(cm, opts); err != nil {
				return nil, err
			}
		case "logfile", "log_file":
			opts.LogFile = v.(string)
		case "syslog":
			opts.Syslog = v.(bool)
		case "remote_syslog":
			opts.RemoteSyslog = v.(string)
		case "pidfile", "pid_file":
			opts.PidFile = v.(string)
		case "prof_port":
			opts.ProfPort = int(v.(int64))
		case "max_control_line":
			opts.MaxControlLine = int(v.(int64))
		case "max_payload":
			opts.MaxPayload = int(v.(int64))
		case "max_pending_size", "max_pending":
			opts.MaxPending = int(v.(int64))
		case "max_connections", "max_conn":
			opts.MaxConn = int(v.(int64))
		}
	}
	return opts, nil
}

// parseCluster will parse the cluster config.
func parseCluster(cm map[string]interface{}, opts *Options) error {
	for mk, mv := range cm {
		switch strings.ToLower(mk) {
		case "port":
			opts.ClusterPort = int(mv.(int64))
		case "host", "net":
			opts.ClusterHost = mv.(string)
		case "authorization":
			am := mv.(map[string]interface{})
			auth := parseAuthorization(am)
			opts.ClusterUsername = auth.user
			opts.ClusterPassword = auth.pass
			opts.ClusterAuthTimeout = auth.timeout
		case "routes":
			ra := mv.([]interface{})
			opts.Routes = make([]*url.URL, 0, len(ra))
			for _, r := range ra {
				routeURL := r.(string)
				url, err := url.Parse(routeURL)
				if err != nil {
					return fmt.Errorf("error parsing route url [%q]", routeURL)
				}
				opts.Routes = append(opts.Routes, url)
			}
		}
	}
	return nil
}

// Helper function to parse Authorization configs.
func parseAuthorization(am map[string]interface{}) authorization {
	auth := authorization{}
	for mk, mv := range am {
		switch strings.ToLower(mk) {
		case "user", "username":
			auth.user = mv.(string)
		case "pass", "password":
			auth.pass = mv.(string)
		case "timeout":
			at := float64(1)
			switch mv.(type) {
			case int64:
				at = float64(mv.(int64))
			case float64:
				at = mv.(float64)
			}
			auth.timeout = at
		}
	}
	return auth
}

// MergeOptions will merge two options giving preference to the flagOpts
// if the item is present.
func MergeOptions(fileOpts, flagOpts *Options) *Options {
	if fileOpts == nil {
		return flagOpts
	}
	if flagOpts == nil {
		return fileOpts
	}
	// Merge the two, flagOpts override
	opts := *fileOpts

	if flagOpts.Port != 0 {
		opts.Port = flagOpts.Port
	}
	if flagOpts.Host != "" {
		opts.Host = flagOpts.Host
	}
	if flagOpts.Username != "" {
		opts.Username = flagOpts.Username
	}
	if flagOpts.Password != "" {
		opts.Password = flagOpts.Password
	}
	if flagOpts.Authorization != "" {
		opts.Authorization = flagOpts.Authorization
	}
	if flagOpts.HTTPPort != 0 {
		opts.HTTPPort = flagOpts.HTTPPort
	}
	if flagOpts.Debug {
		opts.Debug = true
	}
	if flagOpts.Trace {
		opts.Trace = true
	}
	if flagOpts.Logtime {
		opts.Logtime = true
	}
	if flagOpts.LogFile != "" {
		opts.LogFile = flagOpts.LogFile
	}
	if flagOpts.PidFile != "" {
		opts.PidFile = flagOpts.PidFile
	}
	if flagOpts.ProfPort != 0 {
		opts.ProfPort = flagOpts.ProfPort
	}
	if flagOpts.RoutesStr != "" {
		mergeRoutes(&opts, flagOpts)
	}
	return &opts
}

func RoutesFromStr(routesStr string) []*url.URL {
	routes := strings.Split(routesStr, ",")
	if len(routes) == 0 {
		return nil
	}
	routeUrls := []*url.URL{}
	for _, r := range routes {
		r = strings.TrimSpace(r)
		u, _ := url.Parse(r)
		routeUrls = append(routeUrls, u)
	}
	return routeUrls
}

// This will merge the flag routes and override anything that was present.
func mergeRoutes(opts, flagOpts *Options) {
	routeUrls := RoutesFromStr(flagOpts.RoutesStr)
	if routeUrls == nil {
		return
	}
	opts.Routes = routeUrls
	opts.RoutesStr = flagOpts.RoutesStr
}

func RemoveSelfReference(clusterPort int, routes []*url.URL) ([]*url.URL, error) {
	var cleanRoutes []*url.URL
	cport := strconv.Itoa(clusterPort)

	selfIPs := getInterfaceIPs()
	for _, r := range routes {
		host, port, err := net.SplitHostPort(r.Host)
		if err != nil {
			return nil, err
		}

		if cport == port && isIpInList(selfIPs, getUrlIp(host)) {
			Noticef("Self referencing IP found: ", r)
			continue
		}
		cleanRoutes = append(cleanRoutes, r)
	}

	return cleanRoutes, nil
}

func isIpInList(list1 []net.IP, list2 []net.IP) bool {
	for _, ip1 := range list1 {
		for _, ip2 := range list2 {
			if ip1.Equal(ip2) {
				return true
			}
		}
	}
	return false
}

func getUrlIp(ipStr string) []net.IP {
	ipList := []net.IP{}

	ip := net.ParseIP(ipStr)
	if ip != nil {
		ipList = append(ipList, ip)
		return ipList
	}

	hostAddr, err := net.LookupHost(ipStr)
	if err != nil {
		Errorf("Error looking up host with route hostname: %v", err)
		return ipList
	}
	for _, addr := range hostAddr {
		ip = net.ParseIP(addr)
		if ip != nil {
			ipList = append(ipList, ip)
		}
	}
	return ipList
}

func getInterfaceIPs() []net.IP {
	var localIPs []net.IP

	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		Errorf("Error getting self referencing address: %v", err)
		return localIPs
	}

	for i := 0; i < len(interfaceAddr); i++ {
		interfaceIP, _, _ := net.ParseCIDR(interfaceAddr[i].String())
		if net.ParseIP(interfaceIP.String()) != nil {
			localIPs = append(localIPs, interfaceIP)
		} else {
			Errorf("Error parsing self referencing address: %v", err)
		}
	}
	return localIPs
}

func processOptions(opts *Options) {
	// Setup non-standard Go defaults
	if opts.Host == "" {
		opts.Host = DEFAULT_HOST
	}
	if opts.Port == 0 {
		opts.Port = DEFAULT_PORT
	} else if opts.Port == RANDOM_PORT {
		// Choose randomly inside of net.Listen
		opts.Port = 0
	}
	if opts.MaxConn == 0 {
		opts.MaxConn = DEFAULT_MAX_CONNECTIONS
	}
	if opts.PingInterval == 0 {
		opts.PingInterval = DEFAULT_PING_INTERVAL
	}
	if opts.MaxPingsOut == 0 {
		opts.MaxPingsOut = DEFAULT_PING_MAX_OUT
	}
	if opts.SslTimeout == 0 {
		opts.SslTimeout = float64(SSL_TIMEOUT) / float64(time.Second)
	}
	if opts.AuthTimeout == 0 {
		opts.AuthTimeout = float64(AUTH_TIMEOUT) / float64(time.Second)
	}
	if opts.ClusterAuthTimeout == 0 {
		opts.ClusterAuthTimeout = float64(AUTH_TIMEOUT) / float64(time.Second)
	}
	if opts.MaxControlLine == 0 {
		opts.MaxControlLine = MAX_CONTROL_LINE_SIZE
	}
	if opts.MaxPayload == 0 {
		opts.MaxPayload = MAX_PAYLOAD_SIZE
	}
	if opts.MaxPending == 0 {
		opts.MaxPending = MAX_PENDING_SIZE
	}
}
