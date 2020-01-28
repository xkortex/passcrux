/*
Copyright © 2019 MICHAEL McDERMOTT
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0. If a copy of the MPL was not
distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
See also LICENSE and README.md for additional information.

Shamir's Secret Sharing copyright Hashicorp Vault, Mozilla Public License, v. 2.0.
PassCrux uses SSS wholesale without any modification.
*/

package main

import (
	"github.com/xkortex/passcrux/cmd"
)

var Version = "dev"

func main() {
	cmd.Version = Version // todo: need less hackish way of setting version

	cmd.Execute()
}
