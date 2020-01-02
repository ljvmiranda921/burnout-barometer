# burnout-barometer [![GoDoc](https://godoc.org/github.com/ljvmiranda921/burnout-barometer?status.svg)](https://godoc.org/github.com/ljvmiranda921/burnout-barometer)

A Slack tool to log, track, and asses you or your team's stress and work-life.

## Usage

To use within Slack, just type:

```
/barometer 0 "what a stressful day"
/barometer 5 "today's really great!"
/barometer [integer between 0 to 5] [optional string] 
```

## Installation

This application works as a [Cloud
Function](https://cloud.google.com/functions/) that stores results into
[BigQuery](https://cloud.google.com/bigquery/) in the Google Cloud Platform.
Ideally in the future, we'd update this to accommodate more open platforms like
OpenFaas, PostgresSQL, and more.

**Lastly**, it assumes that you know how to create a Slack Application. Create
a Slack App with a corresponding slash-command (we recommend using
`/barometer`), then take note of the API Token provided for you.  You can find
more information [here](https://api.slack.com/start).

### Deploying the Cloud Function
First, clone the repository:

```sh
git clone git@github.com:ljvmiranda921/burnout-barometer.git
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
gcloud functions deploy BurnoutBarometer --runtime go111 --trigger-http
```

This command will provide you with the endpoint you'll use for Slack's
slash-command. For more information, please refer to [this
link](https://cloud.google.com/functions/docs/tutorials/slack).

