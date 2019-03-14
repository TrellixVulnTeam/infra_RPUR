// Copyright 2018 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dutinfo

import (
	"go.chromium.org/luci/common/errors"

	"infra/cmd/skylab_swarming_worker/internal/swmbot"
	"infra/libs/skylab/inventory"
)

// LoadDUTName returns the Swarming bot's DUT's name.
func LoadDUTName(b *swmbot.Info) (string, error) {
	ddir, err := inventory.ReadSymlink(b.Inventory.DataDir)
	if err != nil {
		return "", errors.Annotate(err, "load DUT name").Err()
	}
	lab, err := inventory.LoadLab(ddir)
	if err != nil {
		return "", errors.Annotate(err, "load DUT name").Err()
	}
	for _, d := range lab.GetDuts() {
		c := d.GetCommon()
		if c.GetId() == b.DUTID {
			return c.GetHostname(), nil
		}
	}
	return "", errors.Reason("load DUT name: no DUT").Err()
}