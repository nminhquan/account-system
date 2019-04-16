package transaction

import (
	"strconv"
	"strings"
)

func assignPeers() []string {
	maxId = maxId + 1
	thisPeer := maxId % int64(len(peerList))
	thisPeerStr := strconv.FormatInt(thisPeer+1, 10)
	peers := strings.Split(peerList[thisPeerStr], ",")
	return peers
}
