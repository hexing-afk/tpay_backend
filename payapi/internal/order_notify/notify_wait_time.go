package order_notify

const (
	WaitTimesOne      = 1
	WaitTimesTwo      = 2
	WaitTimesThree    = 3
	WaitTimesFour     = 4
	WaitTimesFive     = 5
	WaitTimesSex      = 6
	WaitTimesSeven    = 7
	WaitTimesEight    = 8
	WaitTimesNine     = 9
	WaitTimesTen      = 10
	WaitTimesEleven   = 11
	WaitTimesTwelve   = 12
	WaitTimesThirteen = 13
	WaitTimesFourteen = 14
	WaitTimesFifteen  = 15
)

var WaitTimes = make(map[int64]int64)

func init() {
	WaitTimes[WaitTimesOne] = 15         // 15s
	WaitTimes[WaitTimesTwo] = 15         // 15s
	WaitTimes[WaitTimesThree] = 30       // 30s
	WaitTimes[WaitTimesFour] = 60        // 60s
	WaitTimes[WaitTimesFive] = 180       // 3m
	WaitTimes[WaitTimesSex] = 600        // 10m
	WaitTimes[WaitTimesSeven] = 600      // 10m
	WaitTimes[WaitTimesEight] = 1800     // 30m
	WaitTimes[WaitTimesNine] = 3600      // 1h
	WaitTimes[WaitTimesTen] = 7200       // 2h
	WaitTimes[WaitTimesEleven] = 10800   // 3h
	WaitTimes[WaitTimesTwelve] = 10800   // 3h
	WaitTimes[WaitTimesThirteen] = 10800 // 3h
	WaitTimes[WaitTimesFourteen] = 21600 // 6h
	WaitTimes[WaitTimesFifteen] = 21600  // 6h
}

func GetNotifyWaitTimeById(id int64) int64 {
	if v, ok := WaitTimes[id]; ok {
		return v
	}

	return -1
}
