package main

import (
	"flag"
	"fmt"
	"os"
)

var FILENAME = "supervisord.conf"
var FILESECONDARY = "supervisord_secondary.conf"

//公共信息
func echoCommon(fd *os.File) {
	fd.WriteString(`; Sample supervisor config file.
;
; For more information on the config file, please see:
; Note: shell expansion ("~" or "$HOME") is not supported.  Environment
; variables can be expanded using this syntax: "%(ENV_HOME)s".

[unix_http_server]
file=/tmp/supervisor.sock   ; (the path to the socket file)
;chmod=0700                 ; socket file mode (default 0700)
;password=123               ; (default is no password (open server))

[inet_http_server]         ; inet (TCP) server disabled by default
port=*:9011        ; (ip_address:port specifier, *:port for all iface)
username=xiaoyao              ; (default is no username (open server))
password=xiaoyao               ; (default is no password (open server))

[supervisord]
logfile=/logs/supervisord.log ; (main log file;default $CWD/supervisord.log)
logfile_maxbytes=50MB        ; (max main logfile bytes b4 rotation;default 50MB)
logfile_backups=10           ; (num of main logfile rotation backups;default 10)
loglevel=info                ; (log level;default info; others: debug,warn,trace)
pidfile=/tmp/supervisord.pid ; (supervisord pidfile;default supervisord.pid)
nodaemon=false               ; (start in foreground if true;default false)
minfds=1024                  ; (min. avail startup file descriptors;default 1024)
minprocs=200                 ; (min. avail process descriptors;default 200)
environment=GODEBUG="gctrace=1"
;umask=022                   ; (process file creation umask;default 022)
;user=chrism                 ; (default is current user, required if root)
;identifier=supervisor       ; (supervisord identifier, default is 'supervisor')
;directory=/tmp              ; (default is not to cd during start)
;nocleanup=true              ; (don't clean up tempfiles at start;default false)
;childlogdir=/tmp            ; ('AUTO' child log dir, default $TEMP)
;strip_ansi=false            ; (strip ansi escape codes in logs; def. false)
; the below section must remain in the config file for RPC
; (supervisorctl/web interface) to work, additional interfaces may be
; added by defining them in separate rpcinterface: sections
[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisorctl]
serverurl=unix:///tmp/supervisor.sock ; use a unix:// URL  for a unix socket
serverurl=http://127.0.0.1:9011 ; use an http:// url to specify an inet socket
;username=chris              ; should be same as http_username if set
;password=123                ; should be same as http_password if set
;prompt=mysupervisor         ; cmd line prompt (default "supervisor")
;history_file=~/.sc_history  ; use readline history if available

; The below sample program section shows all possible program subsection values,
; create one or more 'real' program: sections to be able to control them under
; supervisor.

[program:gnatsd]
command=/server/bin/gnatsd -c /server/bin/gnatsd.conf
autorestart=true
autostart=true
stdout_logfile=/logs/gnats/stdout_gnatsd.log
stderr_logfile=/logs/gnats/stderr_gnatsd.log

[program:main.battery_gateway_server]
command=/server/bin/latest/battery_gateway_server -config=/server/bin/battery_gateway_server.ini -nodeid=0
autorestart=true
autostart=false
stdout_logfile=/logs/gateway/main.stdout_gateway.log
stderr_logfile=/logs/gateway/main.stderr_gateway.log

[program:main.battery_file_server]
command=/server/bin/latest/battery_file_server -httpdoc=./httpdoc -logpath=/logs/file -stdout=0
autorestart=true
autostart=false
stdout_logfile=/logs/file/main.stdout_file.log
stderr_logfile=/logs/file/main.stderr_file.log

[program:main.battery_apn_server]
command=/server/bin/latest/battery_apn_server -production -dburl=mongodb://localhost:20143/briosdb -db=briosdb -stamina=10 -pushcd=24h -logpath=/logs/apn -cert=/server/certs/aps_production.pem -key=/server/certs/key.unencrypted.pem -test=false
autorestart=true
autostart=false
stdout_logfile=/logs/apn/main.stdout_apn.log
stderr_logfile=/logs/apn/main.stderr_apn.log

`)
}

func echoAppAndTransaction(programPrefix string, fd *os.File, appBeginNode, appCount, transactionBeginNode, transactionCount int) {

	fmt.Printf("appBeginNode : %d, appCount : %d, transactionBeginNode : %d, transactionCount : %d",
		appBeginNode,
		appCount,
		transactionBeginNode,
		transactionCount)

	for i := appBeginNode; i < appCount+appBeginNode; i++ {
		str := fmt.Sprintf("[program:%s.battery_app_server_%02d]\n"+
			"command=/server/bin/latest/battery_app_server -config=/server/bin/battery_app_server.ini -nodeid=%d\n"+
			"autorestart=true\n"+
			"autostart=false\n"+
			"stdout_logfile=/logs/app/main.stdout_app_%02d.log\n"+
			"stderr_logfile=/logs/app/main.stderr_app_%02d.log\n\n",
			programPrefix, i, i, i, i)
		fd.WriteString(str)
	}

	for i := transactionBeginNode; i < transactionCount+transactionBeginNode; i++ {
		str := fmt.Sprintf("[program:%s.battery_transaction_server_%02d]\n"+
			"command=/server/bin/latest/battery_transaction_server -config=/server/bin/battery_transaction_server.ini -nodeid=%d\n"+
			"autorestart=true\n"+
			"autostart=false\n"+
			"stdout_logfile=/logs/transaction/main.stdout_transaction_%02d.log\n"+
			"stderr_logfile=/logs/transaction/main.stderr_transaction_%02d.log\n\n",
			programPrefix, i, i, i, i)
		fd.WriteString(str)
	}
}

func echoGroup(prefix string, fd *os.File, appBeginNode, appCount, transactionBeginNode, transactionCount int) {
	str := fmt.Sprintf("[group:%s]\nprograms=", prefix)

	num := 0
	for i := appBeginNode; i < appCount+appBeginNode; i++ {
		if num == 0 {
			str += fmt.Sprintf("%s.battery_app_server_%02d", prefix, i)
			num++
		} else {
			str += fmt.Sprintf(",%s.battery_app_server_%02d", prefix, i)
		}

	}

	for i := transactionBeginNode; i < transactionCount+transactionBeginNode; i++ {
		if num == 0 {
			str += fmt.Sprintf("%s.battery_transaction_server_%02d", prefix, i)
			num++
		} else {
			str += fmt.Sprintf(",%s.battery_transaction_server_%02d", prefix, i)
		}
	}
	str += "\n\n"
	fd.WriteString(str)
}

func main() {
	var (
		appCount                  int
		appBeginNode              int
		secondaryAppCount         int
		transactionCount          int
		transactionBeginNode      int
		secondaryTransactionCount int
	)

	flag.IntVar(&transactionCount, "tc", 4, "transaction node count")
	flag.IntVar(&transactionBeginNode, "tbegin", 0, "transaction begin node")
	flag.IntVar(&secondaryTransactionCount, "stc", 4, "secondary transaction node count")
	flag.IntVar(&appCount, "ac", 4, "app node count")
	flag.IntVar(&appBeginNode, "abegin", 0, "app begin node")
	flag.IntVar(&secondaryAppCount, "sac", 4, "secondary app node count")
	flag.Parse()
	fd, err := os.Create(FILENAME)
	if err != nil {
		fmt.Printf("create file %s failed : %v", FILENAME, err)
		return
	}

	fmt.Printf("transactionCount : %d, secondaryTransactionCount : %d, appCount : %d, secondaryAppCount : %d",
		transactionCount,
		secondaryTransactionCount,
		appCount,
		secondaryAppCount)

	echoCommon(fd)
	echoAppAndTransaction("primary", fd, appBeginNode, appCount, transactionBeginNode, transactionCount)
	//echoAppAndTransaction("secondary", fd, secondaryAppCount, secondaryAppCount, secondaryTransactionCount, secondaryTransactionCount)
	echoGroup("primary", fd, appBeginNode, appCount, transactionBeginNode, transactionCount)
	//echoGroup("secondary", fd, appCount, appCount, transactionCount, transactionCount)
	fd.Close()

	//fd, err = os.Create(FILESECONDARY)
	//if err != nil {
	//	fmt.Printf("create file %s failed : %v", FILESECONDARY, err)
	//	return
	//}

	//echoCommon(fd)
	//echoAppAndTransaction("primary", fd, appCount, appCount, transactionCount, transactionCount)
	//echoAppAndTransaction("secondary", fd, 0, secondaryAppCount, 0, secondaryTransactionCount)
	//echoGroup("primary", fd, appCount, appCount, transactionCount, transactionCount)
	//echoGroup("secondary", fd, 0, appCount, 0, transactionCount)
	//fd.Close()

}
