/* Acksin STRUM - Linux Diagnostics
 * Copyright (C) 2016 Acksin <hey@acksin.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package stats

import (
	"github.com/acksin/procfs"
	"github.com/acksin/strum/io"
	"github.com/acksin/strum/kernel"
	"github.com/acksin/strum/memory"
	"github.com/acksin/strum/network"

	"github.com/acksin/strum/cloud"
)

// Stats contains both the system and process statistics.
type Stats struct {
	// System specific information
	System System
	// Cloud specific information
	Cloud *cloud.Cloud
	// Processes are the process information of the system
	Processes []Process
}

func (n *Stats) processes() procfs.Procs {
	fs, err := procfs.NewFS(procfs.DefaultMountPoint)
	if err != nil {
		return nil
	}

	procs, err := fs.AllProcs()
	if err != nil {
		return nil
	}

	return procs
}

func (n *Stats) containsPid(pids []int, proc procfs.Proc) bool {
	for _, pid := range pids {
		if proc.PID == pid {
			return true
		}
	}

	return false
}

// New returns stats of the machine with pids filtering for
// processes. If pids are empty then it returns all process stats.
func New(pids []int) (s *Stats) {
	s = &Stats{}

	s.System.Memory = memory.New()
	s.System.Network = network.New()

	s.System.Kernel = kernel.New()

	s.Cloud = cloud.New()

	for _, proc := range s.processes() {
		exe, err := proc.Executable()
		if err != nil || exe == "" {
			status, _ := proc.NewStatus()
			exe = status.Name
		}

		p := Process{
			Exe:    exe,
			PID:    proc.PID,
			Memory: memory.NewProcess(proc),
			IO:     io.NewProcess(proc),
		}

		if len(pids) == 0 || s.containsPid(pids, proc) {
			s.Processes = append(s.Processes, p)
		}
	}

	return s
}
