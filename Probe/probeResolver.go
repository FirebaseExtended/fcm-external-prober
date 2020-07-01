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
	"strconv"
	"sync"
	"time"
)

var unresolvedProbes []time.Time
var resolve = true
var probeLock sync.Mutex
var resolveLock sync.Mutex

func resolveProbes(c chan bool) {
	i := 0
	// Continue to resolve probes after no more messages are being sent
	for resolve {
		for len(unresolvedProbes) > 0 {
			res := resolveProbe(unresolvedProbes[i])
			if res {
				removeProbe(i)
				i = i % len(unresolvedProbes)
			} else {
				i = (i + 1) % len(unresolvedProbes)
			}
		}
		i = 0
		// Avoid unnecessary processing while waiting for more probes
		time.Sleep(10 * time.Second)
	}
	c <- true
}

func stopResolving() {
	resolveLock.Lock()
	resolve = false
	resolveLock.Unlock()
}

func addProbe(t time.Time) {
	probeLock.Lock()
	unresolvedProbes = append(unresolvedProbes, t)
	probeLock.Unlock()
}

func removeProbe(i int) {
	probeLock.Lock()
	// Constant time complexity, but does not maintain order
	unresolvedProbes[i] = unresolvedProbes[len(unresolvedProbes)-1]
	unresolvedProbes = unresolvedProbes[:len(unresolvedProbes)-1]
	probeLock.Unlock()
}

func resolveProbe(t time.Time) bool {
	st, err := getMessage(t.Format(timeFileFormat))
	if err != nil {

	}
	if st == "nf" {
		// File not found, so probe is still unresolved
		logProbe(t.Format(timeLogFormat), "unresolved", -1)
		return false
	} else {
		lat, err := calculateLatency(t, st)
		if err != nil {
			// Message received but data is not present/readable
			logProbe(t.Format(timeLogFormat), "error", lat)
		} else {
			logProbe(t.Format(timeLogFormat), "resolved", lat)
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
