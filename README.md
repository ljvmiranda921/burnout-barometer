# burnout-barometer [![Build Status](https://dev.azure.com/ljvmiranda/ljvmiranda/_apis/build/status/ljvmiranda921.burnout-barometer?branchName=master)](https://dev.azure.com/ljvmiranda/ljvmiranda/_build/latest?definitionId=6&branchName=master) [![GoDoc](https://godoc.org/github.com/ljvmiranda921/burnout-barometer?status.svg)](https://godoc.org/github.com/ljvmiranda921/burnout-barometer)


A Slack tool to log, track, and asses you or your team's stress and work-life.

## Setup

All executables can be downloaded in the Releases tab. You can then setup the
server by running:

```sh
barometer init
```

This will trigger a set of prompts to configure the server. You can then start
it by typing:

```sh
barometer serve
```

Burnout Barometer is also packaged as a Docker image. You can start the server
by running the following command:

```sh
docker run ljvmiranda921/burnout-barometer:latest 
```

You can then deploy this however you want. We highly-recommend using
Functions-as-a-Service (FaaS) platforms like Google Cloud Functions, Cloud Run,
AWS Lambda, OpenFaas, etc. 

## Deploying to Google Cloud Functions

This application also provides a [Cloud
Function](https://cloud.google.com/functions/) that stores results into
[BigQuery](https://cloud.google.com/bigquery/). To use it, first clone the
repository and go the the `function` directory:

```sh
git clone git@github.com:ljvmiranda921/burnout-barometer.git
cd function
```

And then configure `config.json`: 

```json
{
    "PROJECT_ID": "<my-project>", 
    "BQ_TABLE": "<my-project>.<my-dataset>.<my-table>", 
    "SLACK_TOKEN": "<base64-encoded-token>",
    "AREA": "<locale-for-datetime>"
}
```


Once complete, deploy the Cloud Function!

```sh
gcloud functions deploy BurnoutBarometerFn --runtime go111 --trigger-http
```

This command will provide you with the endpoint you'll use for Slack's
slash-command. For more information, please refer to [this
link](https://cloud.google.com/functions/docs/tutorials/slack).

