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

package main

import (
	"container/list"
	"log"
	"strconv"
	"sync"
	"time"
)

var unresolvedProbes list.List
var resolve = true
var probeLock sync.Mutex
var resolveLock sync.Mutex

func resolveProbes(wg *sync.WaitGroup) {
	// Continue to resolve probes after no more messages are being sent
	for resolve || unresolvedProbes.Len() > 0 {
		// Process probes in batches
		for unresolvedProbes.Len() > 0 {
			p := removeProbe()
			res := resolveProbe(p)
			if !res {
				addProbe(p)
				if unresolvedProbes.Len() == 1 {
					// Wait before reattempting to resolve this probe
					time.Sleep(time.Duration(probeInterval) * time.Second)
				}
			}
		}
		// Avoid unnecessary polling while waiting for more probes
		time.Sleep(time.Duration(probeInterval) * time.Second)
	}
	wg.Done()
}

func stopResolving() {
	resolveLock.Lock()
	resolve = false
	resolveLock.Unlock()
}

func addProbe(t time.Time) {
	probeLock.Lock()
	unresolvedProbes.PushBack(t)
	probeLock.Unlock()
}

func removeProbe() time.Time {
	probeLock.Lock()
	rem := unresolvedProbes.Remove(unresolvedProbes.Front())
	ret, ok := rem.(time.Time)
	if !ok {
		log.Fatal("removeProbe: item in unresolvedProbes no of type time.Time")
	}
	probeLock.Unlock()
	return ret
}

func resolveProbe(t time.Time) bool {
	st, err := getMessage(t)
	if err != nil {
		pLog.logProbe(t.Format(timeLogFormat), "error", -1)
		return true
	}
	if st == "nf" {
		// Time out probe if it has been unresolved for too long
		if clock.Now().After(t.Add(time.Duration(probeTimeout) * time.Second)) {
			pLog.logProbe(t.Format(timeLogFormat), "timeout", -1)
			return true
		}
		// File not found, so probe is still unresolved
		pLog.logProbe(t.Format(timeLogFormat), "unresolved", -1)
		return false
	} else {
		lat, err := calculateLatency(t, st)
		if err != nil {
			// Message received but data is not present/readable
			pLog.logProbe(t.Format(timeLogFormat), "error", lat)
		} else {
			pLog.logProbe(t.Format(timeLogFormat), "resolved", lat)
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
	return int(int64(t2) - t1), nil
}
