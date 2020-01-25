---
title: Installation
nav_order: 1
layout: default
description: "Set-up and Installation"
---


# Installation
{: .no_toc}


Burnout Barometer is easy-to-configure and deployable as a serverless application.
{: .fs-6 .fw-300 }

Assuming that you already know how to create a [Slack
App](https://api.slack.com/start), this page will walk you through setting-up
the barometer and various deployment options at your disposal.

---


## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

## Initial Setup 

1. **Download the executable**. Ensure that you have downloaded the `barometer`
   executable. Follow the [download instructions]({{ site.baseurl }}/downloads)
   for more info.
2. **Initialize configuration**. Run `barometer init`. It will start a series
   of prompts that will walk you through in configuring the barometer. The
   following values need to be set:


    v1.0.0-alpha
    {: .label .label }

    | Option         | Description                                                                                                                                                                                                                                                          |
    |----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
    | GCP Project ID | The Google Cloud Project ID (GCP) for easy-access of GCP resources. This will be deprecated in the first major release.                                                                                                                                              |
    | Table          | The database connection URL to store Barometer logs. For Bigquery, use the `bq` protocol like so: `bq://my-gcp-project.my-dataset.my-table`                                                                                                                          |
    | Slack Token    | The Slack Token generated whenever you create an App. This is used to verify that the incoming request came from the authorized account. See this [page](https://slack.com/intl/en-ph/help/articles/215770388-Create-and-regenerate-API-tokens) for more information |
    | Area           | The IANA compliant area for correcting the timezone. For example, `Asia/Manila`. This will be deprecated in the first major release.                                                                                                                                 |

3. **Check if config file has been generated**. After running the `init`
   command, you should see a `config.json` file that contains all necessary
   configurations. We will use that later on when deploying or starting the
   server.

## Deployment Options

Burnout Barometer lives inside a server, you can serve the application by
various means. Lastly, you can also take advantage of our Docker image located
in the Azure Container Registry.

### Deploy via Google Cloud Functions

To deploy via [Google Cloud Functions](https://cloud.google.com/functions/),
clone the Github repository to access the functions file:

```bash
git clone git@github.com:ljvmiranda921/burnout-barometer.git
```

Copy over the `config.json` that you've generated into the `function`
directory:

```bash
cp path/of/config.json function
```

Lastly, head to the `function/` directory and execute the Cloud Function deploy
command:

```bash
cd function/
gcloud functions deploy BurnoutBarometerFn --runtime go111 --triger-http
```

Once successful, Cloud Functions will provide you with a URL that you can now
add in your Slack Application's Slash command.

### Deploy via Google Cloud Run
