---
title: Usage
nav_order: 2
layout: default
permalink: usage/
description: "Usage"
---

# Usage

Burnout Barometer can be conveniently used within Slack. Use it for individual
life-logging or for measuring your team's health.
{: .fs-6 .fw-300 }


## Logging from Slack 


The workhorse operation for the barometer is a
[slash-command](https://api.slack.com/interactivity/slash-commands)&mdash; it's
quick, easy, and memorable. Assuming that you created a `/barometer` command,
make your first log by typing:

```
/barometer 5 "my first barometer log"
```

Let's break this down into components:
- **The burnout barometer slash command** (`/barometer`): this could be anything,
    depending on how you set-up your Slack application, but we highly-recommend
    using `/barometer` since that we'll be using throughout the documentation.
- **Your "mood-level," an integer between 1-5** (`5`): we use this to assign a
    quantitative health-measurement throughout all event logs. Check the next
    section on our recommended approach in assigning numbers to moods.
- **The log message** (`"my first barometer log"`): a short note on the current
    log event. You can talk about what happened, or what you're currently
    feeling at the moment. Message length can vary, but we recommend keeping it
    short and sweet.


### On mood-levels

Mood-levels are a great way to (1) **quantitatively assess** what you've felt
throughout the day and (2) **label your thoughts** and feelings. It is an
integer between 1 to 5, and the measurement is self-reported. You can interpret the numbers in any way to tailor-fit yourself or your
organization, but we recommend this as a guide:

| Level | Description                                                                                                                                                                |
|-------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 5     | Things are going well and you can truly say you're happy. Just had dinner with your closest friends, on a trip you've planned months ago, a pat-on-the-back at work, etc. |
| 4     | You're looking forward to something: going home in the weekends, a friday night with friends, etc. Things are getting optimistic and you feel much good. |
| 3     | It's neither good nor bad. We recommend doing frequent spot-checks throughout the day to practice a certain level of  mindfulness.                                                             |
| 2     | Something happened that ticked you off. You're nervous about what's going to happen later throughout the day, an annoying comment that you didn't like, etc.                                                         |
| 1     | Bad things have happened and you need an avenue to vent out. If you've been feeling a lot of 1s through the days, we recommend reaching-out to someone.                    |

- <span class="label label-yellow">Coming Soon</span> Use emojis in
    the slash command and map it to the integer mood levels. 
