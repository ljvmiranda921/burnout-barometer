---
title: Download
nav_order: 3
layout: default
description: "Download"
---

# Download

Burnout Barometer is built on Go and packaged into a single binary.
{: .fs-6 .fw-300 }

You can get the `barometer` executable from the
[Releases](https://github.com/ljvmiranda921/burnout-barometer/releases) page. 
For the latest version, follow the instructions below.

1. First, specify your operating system in the `OS` environment variable. Choose
between `linux`, `windows`, or `darwin`:

    ```bash
    export OS=<my-operating-system>  # [linux|windows|darwin]
    ```

2. Then download the executable:


    ```bash
    curl -s https://api.github.com/repos/ljvmiranda921/burnout-barometer/releases/latest \
    | grep "barometer-amd64-${OS}" \
    | cut -d '"' -f 4 \
    | wget -i - 
    ```

3. Set-up permissions so that it can be executed. Let's also rename the executable into `barometer`:

    ```bash
    chmod +x barometer-amd64-${OS}
    mv barometer-amd64-${OS} barometer
    ```

Ensure that you have downloaded barometer correctly, run `barometer --version`
