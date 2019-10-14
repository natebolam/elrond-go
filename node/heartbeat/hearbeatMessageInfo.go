package heartbeat

import (
	"time"
)

// heartbeatMessageInfo retain the message info received from another node (identified by a public key)
type heartbeatMessageInfo struct {
	maxDurationPeerUnresponsive time.Duration
	maxInactiveTime             Duration
	totalUpTime                 Duration
	totalDownTime               Duration

	getTimeHandler     func() time.Time
	timeStamp          time.Time
	isActive           bool
	receivedShardID    uint32
	computedShardID    uint32
	versionNumber      string
	nodeDisplayName    string
	isValidator        bool
	lastUptimeDowntime time.Time
	genesisTime        time.Time
}

// newHeartbeatMessageInfo returns a new instance of a heartbeatMessageInfo
func newHeartbeatMessageInfo(
	maxDurationPeerUnresponsive time.Duration,
	isValidator bool,
	genesisTime time.Time,
	timer Timer,
) (*heartbeatMessageInfo, error) {

	if maxDurationPeerUnresponsive == 0 {
		return nil, ErrInvalidMaxDurationPeerUnresponsive
	}
	if timer == nil || timer.IsInterfaceNil() {
		return nil, ErrNilTimer
	}

	hbmi := &heartbeatMessageInfo{
		maxDurationPeerUnresponsive: maxDurationPeerUnresponsive,
		maxInactiveTime:             Duration{0},
		isActive:                    false,
		receivedShardID:             uint32(0),
		timeStamp:                   genesisTime,
		lastUptimeDowntime:          timer.Now(),
		totalUpTime:                 Duration{0},
		totalDownTime:               Duration{0},
		versionNumber:               "",
		nodeDisplayName:             "",
		isValidator:                 isValidator,
		genesisTime:                 genesisTime,
		getTimeHandler:              timer.Now,
	}

	return hbmi, nil
}

func (hbmi *heartbeatMessageInfo) updateFields(crtTime time.Time) {
	hbmi.updateTimes(crtTime)
}

func (hbmi *heartbeatMessageInfo) computeActive(crtTime time.Time) {
	hbmi.updateTimes(crtTime)
}

func (hbmi *heartbeatMessageInfo) updateTimes(crtTime time.Time) {
	if crtTime.Sub(hbmi.genesisTime) < 0 {
		return
	}
	hbmi.updateMaxInactiveTimeDuration(crtTime)
	hbmi.updateUpAndDownTime(crtTime)
}

// Will update the total time a node was up and down
func (hbmi *heartbeatMessageInfo) updateUpAndDownTime(crtTime time.Time) {
	if hbmi.lastUptimeDowntime.Sub(hbmi.genesisTime) < 0 {
		hbmi.lastUptimeDowntime = hbmi.genesisTime
	}

	if crtTime.Sub(hbmi.timeStamp) < 0 {
		return
	}

	durationSinceLastHeartbeat := crtTime.Sub(hbmi.timeStamp)

	lastDuration := crtTime.Sub(hbmi.lastUptimeDowntime)
	lastDuration = maxDuration(0, lastDuration)

	uptime := time.Duration(0)
	downTime := time.Duration(0)

	insideActiveWindowAfterHeartheat := durationSinceLastHeartbeat < hbmi.maxDurationPeerUnresponsive
	outSideActiveWindowAfterHeartbeat := durationSinceLastHeartbeat-lastDuration > hbmi.maxDurationPeerUnresponsive
	noHeartbeatReceived := hbmi.timeStamp == hbmi.genesisTime

	if noHeartbeatReceived || outSideActiveWindowAfterHeartbeat {
		downTime = lastDuration
	} else if insideActiveWindowAfterHeartheat {
		uptime = lastDuration
	} else {
		downTime = durationSinceLastHeartbeat - hbmi.maxDurationPeerUnresponsive
		uptime = lastDuration - downTime
	}

	hbmi.totalUpTime.Duration += uptime
	hbmi.totalDownTime.Duration += downTime

	hbmi.lastUptimeDowntime = crtTime
}

// HeartbeatReceived processes a new message arrived from a peer
func (hbmi *heartbeatMessageInfo) HeartbeatReceived(
	computedShardID uint32,
	receivedshardID uint32,
	version string,
	nodeDisplayName string,
) {
	crtTime := hbmi.getTimeHandler()
	hbmi.isActive = true
	hbmi.computedShardID = computedShardID
	hbmi.receivedShardID = receivedshardID
	hbmi.versionNumber = version
	hbmi.nodeDisplayName = nodeDisplayName

	hbmi.updateFields(crtTime)
	hbmi.updateMaxInactiveTimeDuration(crtTime)
	hbmi.timeStamp = crtTime

}

func (hbmi *heartbeatMessageInfo) updateMaxInactiveTimeDuration(currentTime time.Time) {
	crtDuration := currentTime.Sub(hbmi.timeStamp)
	crtDuration = maxDuration(0, crtDuration)

	greaterDurationThanMax := hbmi.maxInactiveTime.Duration < crtDuration
	currentTimeAfterGenesis := hbmi.genesisTime.Sub(currentTime) < 0

	if greaterDurationThanMax && currentTimeAfterGenesis {
		hbmi.maxInactiveTime.Duration = crtDuration
	}
}

func maxDuration(first, second time.Duration) time.Duration {
	if first > second {
		return first
	}

	return second
}
