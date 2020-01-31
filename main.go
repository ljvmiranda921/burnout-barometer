// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

/*
Burnout Barometer is a simple Slack tool to log, track, and assess you or your
team's stress and work life.

Installation

You can get the binaries from our Github releases:
https://github.com/ljvmiranda921/burnout-barometer/releases

Or, you can compile this from source by cloning the repository and building it:

    $ git clone git@github.com:ljvmiranda921/burnout-barometer.git
    $ cd burnout-barometer
    $ go get
    $ go build .


Usage

Usage instructions can be found in the documentation:
https://ljvmiranda921.github.io/burnout-barometer/usage.html

Contributing

Simply fork the Github repository and make a Pull Request. We're open to any
kind of contribution, but we'd definitely appreciate (1) implementation of new
features (2) writing documentation and (3) testing.

License

MIT License (c) 2020, Lester James V. Miranda
*/
package main

import (
	"fmt"
	"os"

	"github.com/ljvmiranda921/burnout-barometer/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
