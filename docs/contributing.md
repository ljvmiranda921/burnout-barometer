---
title: Contributing
nav_order: 4
layout: default
description: "Contributing"
---

# Contributing
{: .no_toc}

Contributions are welcome, and they are greatly appreciated! Every
little bit helps, and credit will always be given.
{: .fs-6 .fw-300 }
--- 
## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

## Types of Contributions

### Report Bugs

Report bugs at this [link](https://github.com/ljvmiranda921/burnout-barometer/issues)

If you are reporting a bug, please include:

* Your operating system name and version.
* Any details about your local setup that might be helpful in troubleshooting.
* Detailed steps to reproduce the bug.

### Fix Bugs

Look through the GitHub issues for bugs. Anything tagged with "bug"
and "help wanted" is open to whoever wants to implement it.


### Implement Features

Look through the GitHub issues for features. Anything tagged with "enhancement"
and "help wanted" is open to whoever wants to implement it. Those that are
tagged with "first-timers-only" is suitable for those getting started in open-source software.

### Write Documentation

Burnout Barometer could always use more documentation, whether as part of the
official Geomancer docs, in docstrings, or even on the web in blog posts,
articles, and such.

### Submit Feedback

The best way to send feedback is to file an issue at this [link](https://github.com/ljvmiranda921/burnout-barometer/issues)


If you are proposing a feature:

* Explain in detail how it would work.
* Keep the scope as narrow as possible, to make it easier to implement.
* Remember that this is a volunteer-driven project, and that contributions
  are welcome :)

## Get Started!

Ready to contribute? Here's how to set up `burnout-barometer` for local development.

1. Fork the `burnout-barometer` repo in Github.
2. Clone your fork locally

    ```bash
    git clone git@github.com:your_name_here/burnout-barometer.git
    ```
3. Once cloned, you can then download all dependencies

    ```bash
    export GO111MODULE=on
    go get -v
    ```
4. Create a branch for local development

    ```bash
    git checkout -b name-of-your-bugfix-or-feature
    ```

5. When you're done making changes, run `gofmt` and `go test`:

    ```bash
    gofmt path/to/file.go
    go test -v ./... -cover
    ```
6. Commit your changes and push your branch to Github:

    ```bash
    git add .
    git commit -m "Your detailed description of your changes."
    git push origin name-of-your-bugfix-or-feature
    ```
7. Submit a Pull Request through the Github website.


## Pull Request Guidelines

Before you submit a pull request, check that it meets these guidelines:

1. The pull request should include tests.
2. If the pull request adds functionality, the docs should be updated. Put
   your new functionality into a function with a docstring, and add the
   feature to the list in README.rst.
3. The pull request should work for Go 1.11, and above. Check [Azure
   Pipelines](https://dev.azure.com/ljvmiranda/ljvmiranda/_build/latest?definitionId=6&branchName=master)
   and make sure that the tests pass for all supported operating systems.
