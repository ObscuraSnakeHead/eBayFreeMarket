package webapp

import (
	"fmt"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type StaffStats struct {
	NumberOfResolvedTickets  int
	NumberOfResolvedDisputes int
	NumberOfApprovedItems    int
	NumberOfApprovedVendors  int
}

/*
	Cache
*/

func CacheGetStaffStats(userUuid string) StaffStats {
	queryStats := func() StaffStats {
		return StaffStats{
			NumberOfResolvedDisputes: CountNumberOfDisputesResolved(userUuid),
			NumberOfResolvedTickets:  CountNumberOfTicketsResolved(userUuid),
			NumberOfApprovedItems:    CountNumberOfApprovedItems(userUuid),
			NumberOfApprovedVendors:  CountNumberOfApprovedVendors(userUuid),
		}
	}

	key := fmt.Sprintf("staff-stats-%s", userUuid)
	cStats, _ := util.Cache15m.Get(key)
	if cStats == nil {
		stats := queryStats()
		util.Cache15m.Set(key, stats)
		return stats
	}

	return cStats.(StaffStats)
}
