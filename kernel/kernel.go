/* Acksin STRUM - Linux Diagnostics
 * Copyright (C) 2016 Acksin <hey@acksin.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package kernel

import (
	"github.com/acksin/procfs"
)

func New() map[string]string {
	fs, err := procfs.NewFS(procfs.DefaultMountPoint)
	if err != nil {
		// No-op
	}

	sys, err2 := fs.NewSys()
	if err2 != nil {
		// No-op
	}

	return map[string]string(sys)
}
