// Copyright 2020 Lester James V. Miranda. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

/*
Burnout Barometer is a simple Slack tool to log, track, and assess you or your
team's stress and work life.

Installation

You can get the barometer executable from the Releases page:
https://github.com/ljvmiranda921/burnout-barometer/releases. Binaries are
available for Windows, Linux, and MacOS operating systems with amd64
architecture.

Alternatively, you can build your own binary from source by cloning the
repository and compiling it. The following steps work for Go 1.11 and above:

    $ git clone git@github.com:ljvmiranda921/burnout-barometer.git
    $ cd burnout-barometer
    $ go get
    $ go build .

Burnout Barometer is a server-side application, and can be deployed by various
means. You can also check out our Docker image located in the Azure Container
Registry. For deployment options, head over this link:
https://ljvmiranda921.github.io/burnout-barometer/installation/#deployment-options

Usage

You can find usage instructions for Burnout Barometer in our external
documentation: https://ljvmiranda921.github.io/burnout-barometer/usage/

Contributing

Contributions are welcome, and they are greatly appreciated! Every little bit
helps, and credit will always be given.  Simply fork the Github repository and
make a Pull Request. We're open to any kind of contribution, but we'd
definitely appreciate (1) implementation of new features (2) writing
documentation and (3) testing.

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
