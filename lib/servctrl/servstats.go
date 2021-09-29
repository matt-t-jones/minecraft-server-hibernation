package servctrl

import (
	"fmt"
	"sync"
	"time"

	"msh/lib/errco"
	"msh/lib/logger"
)

// Stats contains the info relative to server
var Stats *serverStats

type serverStats struct {
	M              *sync.Mutex
	Status         int     // represent the status of the minecraft server
	PlayerCount    int     // tracks players connected to the server
	StopMSRequests int32   // tracks active StopMSRequest() instances. (int32 for atomic operations)
	LoadProgress   string  // tracks loading percentage of starting server
	BytesToClients float64 // tracks bytes/s server->clients
	BytesToServer  float64 // tracks bytes/s clients->server
}

func init() {
	Stats = &serverStats{
		M:              &sync.Mutex{},
		Status:         errco.SERVER_STATUS_OFFLINE,
		PlayerCount:    0,
		StopMSRequests: 0,
		LoadProgress:   "0%",
		BytesToClients: 0,
		BytesToServer:  0,
	}
}

// PrintDataUsage prints each second bytes/s to clients and to server.
// (must be launched after ServTerm.IsActive has been set to true)
// [goroutine]
func printDataUsage() {
	for ServTerm.IsActive {
		if Stats.BytesToClients != 0 || Stats.BytesToServer != 0 {
			logger.Logln(fmt.Sprintf("data/s: %8.3f KB/s to clients | %8.3f KB/s to server", Stats.BytesToClients/1024, Stats.BytesToServer/1024))

			Stats.M.Lock()
			Stats.BytesToClients = 0
			Stats.BytesToServer = 0
			Stats.M.Unlock()
		}

		time.Sleep(time.Second)
	}
}
