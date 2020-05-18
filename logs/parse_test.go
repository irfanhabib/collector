package logs_test

import (
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/pganalyze/collector/logs"
	"github.com/pganalyze/collector/output/pganalyze_collector"
	"github.com/pganalyze/collector/state"
)

type parseTestpair struct {
	prefixIn  string
	lineIn    string
	lineOut   state.LogLine
	lineOutOk bool
}

var parseTests = []parseTestpair{
	// rsyslog format
	{
		"",
		"Feb  1 21:48:31 ip-172-31-14-41 postgres[9076]: [3-1] LOG:  database system is ready to accept connections",
		state.LogLine{
			OccurredAt: time.Date(time.Now().Year(), time.February, 1, 21, 48, 31, 0, time.UTC),
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 9076,
			Content:    "database system is ready to accept connections",
		},
		true,
	},
	{
		"",
		"Feb  1 21:48:31 ip-172-31-14-41 postgres[9076]: [3-2] #011 something",
		state.LogLine{
			OccurredAt: time.Date(time.Now().Year(), time.February, 1, 21, 48, 31, 0, time.UTC),
			LogLevel:   pganalyze_collector.LogLineInformation_UNKNOWN,
			BackendPid: 9076,
			Content:    "\t something",
		},
		false,
	},
	{
		"",
		"Feb  1 21:48:31 ip-172-31-14-41 postgres[123]: [8-1] [user=postgres,db=postgres,app=[unknown]] LOG: connection received: host=[local]",
		state.LogLine{
			OccurredAt: time.Date(time.Now().Year(), time.February, 1, 21, 48, 31, 0, time.UTC),
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 123,
			Username:   "postgres",
			Database:   "postgres",
			Content:    "connection received: host=[local]",
		},
		true,
	},
	// RDS format
	{
		"",
		"2018-08-22 16:00:04 UTC:ec2-1-1-1-1.compute-1.amazonaws.com(48808):myuser@mydb:[18762]:LOG:  duration: 3668.685 ms  execute <unnamed>: SELECT 1",
		state.LogLine{
			OccurredAt: time.Date(2018, time.August, 22, 16, 0, 4, 0, time.UTC),
			Username:   "myuser",
			Database:   "mydb",
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 18762,
			Content:    "duration: 3668.685 ms  execute <unnamed>: SELECT 1",
		},
		true,
	},
	{
		"",
		"2018-08-22 16:00:03 UTC:127.0.0.1(36404):myuser@mydb:[21495]:LOG:  duration: 1630.946 ms  execute 3: SELECT 1",
		state.LogLine{
			OccurredAt: time.Date(2018, time.August, 22, 16, 0, 3, 0, time.UTC),
			Username:   "myuser",
			Database:   "mydb",
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 21495,
			Content:    "duration: 1630.946 ms  execute 3: SELECT 1",
		},
		true,
	},
	{
		"",
		"2018-08-22 16:00:03 UTC:[local]:myuser@mydb:[21495]:LOG:  duration: 1630.946 ms  execute 3: SELECT 1",
		state.LogLine{
			OccurredAt: time.Date(2018, time.August, 22, 16, 0, 3, 0, time.UTC),
			Username:   "myuser",
			Database:   "mydb",
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 21495,
			Content:    "duration: 1630.946 ms  execute 3: SELECT 1",
		},
		true,
	},
	// Custom 3 format
	{
		"",
		"2018-09-27 06:57:01.030 UTC [20194] [user=[unknown],db=[unknown],app=[unknown]] LOG:  connection received: host=[local]",
		state.LogLine{
			OccurredAt: time.Date(2018, time.September, 27, 6, 57, 1, 30*1000*1000, time.UTC),
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 20194,
			Content:    "connection received: host=[local]",
		},
		true,
	},
	{
		"",
		"2018-09-27 06:57:02.779 UTC [20194] [user=postgres,db=postgres,app=psql] ERROR:  canceling statement due to user request",
		state.LogLine{
			OccurredAt:  time.Date(2018, time.September, 27, 6, 57, 2, 779*1000*1000, time.UTC),
			Username:    "postgres",
			Database:    "postgres",
			Application: "psql",
			LogLevel:    pganalyze_collector.LogLineInformation_ERROR,
			BackendPid:  20194,
			Content:     "canceling statement due to user request",
		},
		true,
	},
	// Custom 4 format
	{
		"",
		"2018-09-27 06:57:01.030 UTC [20194] [user=[unknown],db=[unknown],app=[unknown],host=[local]] LOG:  connection received: host=[local]",
		state.LogLine{
			OccurredAt: time.Date(2018, time.September, 27, 6, 57, 1, 30*1000*1000, time.UTC),
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 20194,
			Content:    "connection received: host=[local]",
		},
		true,
	},
	{
		"",
		"2018-09-27 06:57:02.779 UTC [20194] [user=postgres,db=postgres,app=psql,host=127.0.0.1] ERROR:  canceling statement due to user request",
		state.LogLine{
			OccurredAt:  time.Date(2018, time.September, 27, 6, 57, 2, 779*1000*1000, time.UTC),
			Username:    "postgres",
			Database:    "postgres",
			Application: "psql",
			LogLevel:    pganalyze_collector.LogLineInformation_ERROR,
			BackendPid:  20194,
			Content:     "canceling statement due to user request",
		},
		true,
	},
	// Custom 5 format
	{
		"",
		"2018-09-28 07:37:59 UTC [331]: [1-1] user=[unknown],db=[unknown] - PG-00000 LOG:  connection received: host=127.0.0.1 port=49738",
		state.LogLine{
			OccurredAt:    time.Date(2018, time.September, 28, 7, 37, 59, 0, time.UTC),
			LogLevel:      pganalyze_collector.LogLineInformation_LOG,
			BackendPid:    331,
			LogLineNumber: 1,
			Content:       "connection received: host=127.0.0.1 port=49738",
		},
		true,
	},
	{
		"",
		"2018-09-28 07:39:48 UTC [347]: [3-1] user=postgres,db=postgres - PG-57014 ERROR:  canceling statement due to user request",
		state.LogLine{
			OccurredAt:    time.Date(2018, time.September, 28, 7, 39, 48, 0, time.UTC),
			Username:      "postgres",
			Database:      "postgres",
			LogLevel:      pganalyze_collector.LogLineInformation_ERROR,
			BackendPid:    347,
			LogLineNumber: 3,
			Content:       "canceling statement due to user request",
		},
		true,
	},
	// Custom 6 format
	{
		"",
		"2018-10-16 01:25:58 UTC [93897]: [4-1] user=,db=,app=,client= LOG:  database system is ready to accept connections",
		state.LogLine{
			OccurredAt:    time.Date(2018, time.October, 16, 1, 25, 58, 0, time.UTC),
			LogLevel:      pganalyze_collector.LogLineInformation_LOG,
			BackendPid:    93897,
			LogLineNumber: 4,
			Content:       "database system is ready to accept connections",
		},
		true,
	},
	{
		"",
		"2018-10-16 01:26:09 UTC [93907]: [1-1] user=[unknown],db=[unknown],app=[unknown],client=::1 LOG:  connection received: host=::1 port=61349",
		state.LogLine{
			OccurredAt:    time.Date(2018, time.October, 16, 1, 26, 9, 0, time.UTC),
			LogLevel:      pganalyze_collector.LogLineInformation_LOG,
			BackendPid:    93907,
			LogLineNumber: 1,
			Content:       "connection received: host=::1 port=61349",
		},
		true,
	},
	{
		"",
		"2018-10-16 01:26:33 UTC [93911]: [3-1] user=postgres,db=postgres,app=psql,client=::1 ERROR:  canceling statement due to user request",
		state.LogLine{
			OccurredAt:    time.Date(2018, time.October, 16, 1, 26, 33, 0, time.UTC),
			Username:      "postgres",
			Database:      "postgres",
			LogLevel:      pganalyze_collector.LogLineInformation_ERROR,
			BackendPid:    93911,
			LogLineNumber: 3,
			Content:       "canceling statement due to user request",
		},
		true,
	},
	// Custom 7 format
	{
		"",
		"2019-01-01 01:59:42 UTC [1]: [4-1] [trx_id=0] user=,db= LOG:  database system is ready to accept connections",
		state.LogLine{
			OccurredAt:    time.Date(2019, time.January, 1, 1, 59, 42, 0, time.UTC),
			LogLevel:      pganalyze_collector.LogLineInformation_LOG,
			BackendPid:    1,
			LogLineNumber: 4,
			Content:       "database system is ready to accept connections",
		},
		true,
	},
	{
		"",
		"2019-01-01 02:00:28 UTC [35]: [1-1] [trx_id=0] user=[unknown],db=[unknown] LOG:  connection received: host=::1 port=38842",
		state.LogLine{
			OccurredAt:    time.Date(2019, time.January, 1, 2, 0, 28, 0, time.UTC),
			LogLevel:      pganalyze_collector.LogLineInformation_LOG,
			BackendPid:    35,
			LogLineNumber: 1,
			Content:       "connection received: host=::1 port=38842",
		},
		true,
	},
	{
		"",
		"2019-01-01 02:00:28 UTC [34]: [3-1] [trx_id=120950] user=postgres,db=postgres ERROR:  canceling statement due to user request",
		state.LogLine{
			OccurredAt:    time.Date(2019, time.January, 1, 2, 0, 28, 0, time.UTC),
			Username:      "postgres",
			Database:      "postgres",
			LogLevel:      pganalyze_collector.LogLineInformation_ERROR,
			BackendPid:    34,
			LogLineNumber: 3,
			Content:       "canceling statement due to user request",
		},
		true,
	},
	// Custom 8 format
	{
		"",
		"[1127]: [8-1] db=postgres,user=pganalyze LOG:  duration: 2001.842 ms  statement: SELECT pg_sleep(2);",
		state.LogLine{
			OccurredAt:    time.Time{},
			Username:      "pganalyze",
			Database:      "postgres",
			LogLevel:      pganalyze_collector.LogLineInformation_LOG,
			BackendPid:    1127,
			LogLineNumber: 8,
			Content:       "duration: 2001.842 ms  statement: SELECT pg_sleep(2);",
		},
		true,
	},
	// Simple format
	{
		"",
		"2018-05-04 03:06:18.360 UTC [3184] LOG:  pganalyze-collector-identify: server1",
		state.LogLine{
			OccurredAt: time.Date(2018, time.May, 4, 3, 6, 18, 360*1000*1000, time.UTC),
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 3184,
			Content:    "pganalyze-collector-identify: server1",
		},
		true,
	},
	{
		"",
		"2018-05-04 03:06:18.360 +0100 [3184] LOG:  pganalyze-collector-identify: server1",
		state.LogLine{
			OccurredAt: time.Date(2018, time.May, 4, 3, 6, 18, 360*1000*1000, time.FixedZone("+0100", 3600)),
			LogLevel:   pganalyze_collector.LogLineInformation_LOG,
			BackendPid: 3184,
			Content:    "pganalyze-collector-identify: server1",
		},
		true,
	},
}

func TestParseLogLineWithPrefix(t *testing.T) {
	for _, pair := range parseTests {
		l, lOk := logs.ParseLogLineWithPrefix(pair.prefixIn, pair.lineIn)

		cfg := pretty.CompareConfig
		cfg.SkipZeroFields = true

		if pair.lineOutOk != lOk {
			t.Errorf("For \"%v\": expected parsing ok? to be %v, but was %v\n", pair.lineIn, pair.lineOutOk, lOk)
		}

		if diff := cfg.Compare(pair.lineOut, l); diff != "" {
			t.Errorf("For \"%v\": log line diff: (-got +want)\n%s", pair.lineIn, diff)
		}
	}
}
