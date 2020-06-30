Fluxcloud is a tool to receive events from the [Weave flux](https://github.com/weaveworks/flux).

![build status](https://ci.codesink.net/api/badges/justinbarrick/fluxcloud/status.svg)

Weave Flux is a useful tool for managing the state of your Kubernetes cluster.

Fluxcloud is a valid upstream for Weave, allowing you to send Flux events to Slack or a
webhook without using Weave Cloud.

# Setup

Please see the [Weave Flux setup documentation](https://github.com/weaveworks/flux/blob/master/site/standalone/installing.md) for setting up Flux.

To use Fluxcloud, you can deploy fluxcloud as either a sidecar to Flux or a seperate deployment.

To deploy as a sidecar, see `examples/flux-deployment-sidecar.yaml`.
To deploy independently, see `examples/fluxcloud.yaml`.

Set the following environment variables in your chosen deployment:

* `SLACK_URL`: the Slack [webhook URL](https://api.slack.com/incoming-webhooks) to use.
* `SLACK_USERNAME`: the Slack username to use when sending messages.
* `SLACK_TOKEN` (optional): legacy Slack API token to use.
* `SLACK_CHANNEL`: the Slack channel to send messages to.
* `SLACK_ICON_EMOJI`: the Slack emoji to use as the icon.
* `SLACK_ICON_EMOJI`: the Slack emoji to use as the icon.
* `SLACK_ICON_EMOJI`: the Slack emoji to use as the icon.
* `DATADOG_API_KEY`: the Datadog API key used to push events.
* `DATADOG_API_KEY`: the Datadog APP key used to push events.
* `DATADOG_ADITIONAL_TAGS`: Datadog aditional tags to be added to the generated event.
* `MSTEAMS_URL`: the Microsoft Teams [webhook URL](https://docs.microsoft.com/en-us/outlook/actionable-messages/send-via-connectors#sending-actionable-messages-via-office-365-connectors) to use
* `GITHUB_URL`: the URL to the Github repository that Flux uses, used for Slack links.
* `WEBHOOK_URL`: if the exporter is "webhook", then the URL to use for the webhook.
* `EXPORTER_TYPE` (optional): The types of exporter to use in comma delimited form. (Ex: `slack,webhook`) (Choices: slack, msteams, datadog, webhook, Default: slack)
* `JAEGER_ENDPOINT` (optional): endpoint to report Jaeger traces to.

And then apply the configuration:

```
kubectl apply -f examples/fluxcloud.yaml
```

Set the `--connect` flag on Flux to `--connect=ws://fluxcloud`.

# Exporters

There are multiple exporters that you can use with fluxcloud. If there is not a suitable
one already, feel free to contribute one by implementing the [exporter interface](https://github.com/topfreegames/fluxcloud/blob/master/pkg/exporters/exporter.go)!

## Slack

The default exporter to use is Slack. To use the Slack exporter, set the `SLACK_URL`,
`SLACK_USERNAME`, and `SLACK_CHANNEL` environment variables to
use. You can also optionally set the `EXPORTER_TYPE` to "slack".

### Sending notifications to multiple channels

If sending notifications to only one channel is unsufficient for your use case you can
configure fluxcloud to send them to multiple channels based upon the namespace(s) from
the created and/or updated resources. This is done by setting a comma separated
`<channel>=<namespace>` string as the `SLACK_CHANNEL` environment variable.

If you for example want to send notifications of all events to `#k8s-events` but only
events from namespace `team-b` to `#teamb` you would set the following string:
`SLACK_CHANNEL=#k8s-events=*,#team-b=team-b`.

## Microsoft Teams

Set the environment variable `MSTEAMS_URL` to the URL generated on activation of an
Incoming Webhook in a Microsoft Teams channel.

## Datadog

Events can be sent to Datadog by adding "datadog" to to `EXPORTER_TYPE` and then setting
the `DATADOG_API_KEY` and the `DATADOG_APP_KEY`. More information about generating those
keys can be found in [Datadog documentation](https://docs.datadoghq.com/account_management/api-app-keys/).

You can also add additional tags to the event by setting `DATADOG_ADDITIONAL_TAGS`.

## Webhooks

Events can be sent to an arbitrary webhook by setting the `EXPORTER_TYPE` to "webhook" and
then setting the `WEBHOOK_URL` to the URL to send the webhook to.

Fluxcloud will send a POST request to the provided URL with [the encoded event](https://github.com/topfreegames/fluxcloud/blob/master/pkg/msg/msg.go) as the payload.

# Formatting commit links

By default, commit links are formatted for Github. It is possible to format them
for another VCS system, such as Bitbucket, by overriding the commit template.

The commit template is a go template that supports two variables:

* `VCSLink`: which is the GITHUB_URL configuration option.
* `Commit`: which is the commit id.

The default is:

```
{{ .VCSLink }}/commit/{{ .Commit }}
```

For example, to override to work for Bitbucket, set the `COMMIT_TEMPLATE` environment
variable to:

```
{{ .VCSLink }}/commits/{{ .Commit }}
```

# Versioning

Fluxcloud follows semver for versioning, but also publishes development images tagged
with `$BRANCH-$COMMIT`.

To track release images:

```
fluxctl policy -c kube-system:deployment/fluxcloud --tag-all='v0*'
```

To track the latest pre-release images:

```
fluxctl policy -c kube-system:deployment/fluxcloud --tag-all='master-*'
```

And then you can automate it:

```
fluxctl automate -c kube-system:deployment/fluxcloud
```

# Build

To build fluxcloud, you can either use go:

```
go build -o fluxcloud ./cmd/
```

Or, to run a full CI build, download [hone](https://github.com/topfreegames/hone):

```
hone
```
