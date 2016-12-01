package syslog

import (
	sl "log/syslog"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	lrs "github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/casualjim/go-app/logging"
	"github.com/spf13/viper"
)

var (
	prioMap     map[string]sl.Priority
	prioMapLock *sync.Mutex

	// DefaultSeverity for syslog, value is syslog.LOG_INFO
	DefaultSeverity sl.Priority
	// DefaultFacility for syslog, value is syslog.LOG_LOCAL0
	DefaultFacility sl.Priority
)

func mapSev(sev string) sl.Priority {
	return mapPrio(sev, DefaultSeverity)
}

func mapFac(fac string) sl.Priority {
	return mapPrio(fac, DefaultFacility)
}

func mapPrio(prio string, defVal sl.Priority) sl.Priority {
	prioMapLock.Lock()
	if r, ok := prioMap[strings.ToLower(prio)]; ok {
		prioMapLock.Unlock()
		return r
	}

	prioMapLock.Unlock()
	return defVal
}

func init() {
	DefaultSeverity = sl.LOG_INFO
	DefaultFacility = sl.LOG_LOCAL0
	prioMapLock = new(sync.Mutex)
	prioMap = map[string]sl.Priority{
		"emerg":     sl.LOG_EMERG,
		"emergency": sl.LOG_EMERG,
		"alert":     sl.LOG_ALERT,
		"crit":      sl.LOG_CRIT,
		"critical":  sl.LOG_CRIT,
		"err":       sl.LOG_ERR,
		"error":     sl.LOG_ERR,
		"warn":      sl.LOG_WARNING,
		"warning":   sl.LOG_WARNING,
		"notice":    sl.LOG_NOTICE,
		"info":      sl.LOG_INFO,
		"debug":     sl.LOG_DEBUG,
		"kern":      sl.LOG_KERN,
		"kernel":    sl.LOG_KERN,
		"user":      sl.LOG_USER,
		"mail":      sl.LOG_MAIL,
		"daemon":    sl.LOG_DAEMON,
		"auth":      sl.LOG_AUTH,
		"syslog":    sl.LOG_SYSLOG,
		"lpr":       sl.LOG_LPR,
		"news":      sl.LOG_NEWS,
		"uucp":      sl.LOG_UUCP,
		"cron":      sl.LOG_CRON,
		"authpriv":  sl.LOG_AUTHPRIV,
		"ftp":       sl.LOG_FTP,
		"local0":    sl.LOG_LOCAL0,
		"local1":    sl.LOG_LOCAL1,
		"local2":    sl.LOG_LOCAL2,
		"local3":    sl.LOG_LOCAL3,
		"local4":    sl.LOG_LOCAL4,
		"local5":    sl.LOG_LOCAL5,
		"local6":    sl.LOG_LOCAL6,
		"local7":    sl.LOG_LOCAL7,
	}

	logging.RegisterHook("syslog", func(v *viper.Viper) logrus.Hook {
		nw := v.GetString("network")
		raddr := v.GetString("address")
		sev := v.GetString("severity")
		fac := v.GetString("facility")
		tag := v.GetString("tag")

		slg, err := lrs.NewSyslogHook(nw, raddr, mapSev(sev)|mapFac(fac), tag)
		if err != nil {
			panic(err)
		}
		return slg
	})
}
