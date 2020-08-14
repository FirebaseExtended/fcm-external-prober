/*
 *  Copyright 2020 Google LLC
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package probe

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const maxUnresolved int = 2000

var (
	resolve     bool
	resolveLock sync.Mutex
	closeLock   sync.Mutex
	closed      bool
	// Use buffered channel so that resolving blocks on having no probes to resolve
	unresolved chan *sentProbe
	latencyOffset int
)

type sentProbe struct {
	sendTime time.Time
	probe    *probe
}

func newSentProbe(tim time.Time, p *probe) *sentProbe {
	return &sentProbe{tim, p}
}

func initResolver() error {
	unresolved = make(chan *sentProbe, maxUnresolved)
	resolve = true
	closed = false
	var err error
	latencyOffset, err = findTimeOffset()
	if err != nil {
		return err
	}
	return nil
}

func resolveProbes(wg *sync.WaitGroup) {
	// Continue to resolve probes after no more messages are being sent
	for resolve {
		sp := removeProbe()
		// use a nil value sent on the channel to indicate that there are no more probes
		closeLock.Lock()
		// If the channel is closed, repeatedly attempt to resolve each probe. Otherwise, add probe back to queue
		if closed {
			for !resolveProbe(sp) {
			}
		} else {
			if !resolveProbe(sp) {
				addProbe(sp)
			}
		}
		closeLock.Unlock()
	}
	wg.Done()
}

func stopResolving() {
	resolveLock.Lock()
	resolve = false
	resolveLock.Unlock()
}

func addProbe(sp *sentProbe) {
	unresolved <- sp
}

func closeUnresolved() {
	closeLock.Lock()
	closed = true
	unresolved <- nil
	close(unresolved)
	closeLock.Unlock()
}

func removeProbe() *sentProbe {
	return <-unresolved
}

func resolveProbe(sp *sentProbe) bool {
	if sp == nil {
		stopResolving()
		return true
	}
	st, err := getMessage(fmt.Sprintf("%d%s", sp.probe.config.GetType(), sp.sendTime.Format(timeFileFormat)))
	if err != nil {
		logger.LogProbe(sp, "error", -1, deviceToken)
		return true
	}
	if st == "nf" {
		// Time out probe if it has been unresolved for too long
		if clock.Now().After(sp.sendTime.Add(time.Duration(sp.probe.config.GetReceiveTimeout()) * time.Second)) {
			logger.LogProbe(sp, "timeout", -1, deviceToken)
			return true
		}
		// File not found, so probe is still unresolved
		logger.LogProbe(sp, "unresolved", -1, deviceToken)
		return false
	} else {
		lat, err := calculateLatency(sp.sendTime, st)
		if err != nil {
			// Message received but data is not present/readable
			logger.LogProbe(sp, "error", lat, deviceToken)
		} else {
			logger.LogProbe(sp, "resolved", lat, deviceToken)
		}
		return true
	}
}

func calculateLatency(st time.Time, rt string) (int, error) {
	t1 := st.UnixNano() / 1000000
	t2, err := strconv.Atoi(rt)
	if err != nil {
		return -1, err
	}
	return int(int64(t2) - t1) + latencyOffset, nil
}
