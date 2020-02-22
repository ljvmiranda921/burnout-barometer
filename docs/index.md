---
title: Home
nav_order: 1
layout: default
description: "Burnout Barometer Homepage"
permalink: /
---

# Mindfulness throughout the day  
{: .fs-9 }

Burnout Barometer is a simple Slack tool to log, track, and assess you or your
team's stress and work life.
{: .fs-6 .fw-300 }


[Download]({{ site.baseurl }}/download){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 } [View it on GitHub](https://github.com/ljvmiranda921/burnout-barometer){: .btn .fs-5 .mb-4 .mb-md-0 }

---


<a class="github-button" href="https://github.com/ljvmiranda921/burnout-barometer/subscription" data-icon="octicon-eye" data-size="large" data-show-count="true" aria-label="Watch ljvmiranda921/burnout-barometer on GitHub">Watch</a>
<a class="github-button" href="https://github.com/ljvmiranda921/burnout-barometer" data-icon="octicon-star" data-size="large" data-show-count="true" aria-label="Star ljvmiranda921/burnout-barometer on GitHub">Star</a>
<a class="github-button" href="https://github.com/ljvmiranda921/burnout-barometer/fork" data-icon="octicon-repo-forked" data-size="large" data-show-count="true" aria-label="Fork ljvmiranda921/burnout-barometer on GitHub">Fork</a>


Burnout Barometer functions as a Slack application where you log your current mood:

![]({{ site.baseurl }}/assets/demo.gif)

All logs are then stored in a database like PostgreSQL or BigQuery that you set
up. With that said, you own your own data and it's not shared to other
entities or organizations.

## Personal

The Barometer is first and foremost a personal tool. I made it to meet specific
needs that works for me. However, I'm releasing it in public with the hope that
other people will find this useful as much as I did.  These tools may not work
for you, but I'll try my best so that you can also set-up the Barometer easily
as any other open-source code.

## Open-source

I'm open-sourcing the Barometer to encourage myself to write and deliver
quality software. I admit that my Golang and schema design skills aren't
topnotch, so if you find any code smells, bugs, or incorrect practices, please
feel free to create an Issue or make a Pull Request in the [Github
page](https://github.com/ljvmiranda921/burnout-barometer). I will also
appreciate if you can suggest features that you would like the Barometer to
have, it would be nice to build this project together!

## Private

As you will see later on, the Barometer requires a PostgreSQL database (or
BigQuery) to house collected data. That dataset is yours. You set-up your own
database instance and manage your own machine. I do all of these inside [Google
Cloud Run](https://cloud.google.com/run), but you can do it in any platform of
your choice (or perhaps even in a Raspberry Pi!). 

