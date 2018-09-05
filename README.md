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
* `SLACK_CHANNEL`: the Slack channel to send messages to.
* `SLACK_ICON_EMOJI`: the Slack emoji to use as the icon.
* `GITHUB_URL`: the URL to the Github repository that Flux uses, used for Slack links.
* `WEBHOOK_URL`: if the exporter is "webhook", then the URL to use for the webhook.
* `EXPORTER_TYPE` (optional): The type of exporter to use. (Choices: slack, webhook, Default: slack)

And then apply the configuration:

```
kubectl apply -f examples/fluxcloud.yaml
```

Set the `--connect` flag on Flux to `--connect=ws://fluxcloud`.

# Exporters

There are multiple exporters that you can use with fluxcloud. If there is not a suitable
one already, feel free to contribute one by implementing the [exporter interface](https://github.com/justinbarrick/fluxcloud/blob/master/pkg/exporters/exporter.go)!

## Slack

The default exporter to use is Slack. To use the Slack exporter, set the `SLACK_URL`,
`SLACK_USERNAME`, and `SLACK_CHANNEL` environment variables to use. You can also
optionally set the `EXPORTER_TYPE` to "slack".

## Webhooks

Events can be sent to an arbitrary webhook by setting the `EXPORTER_TYPE` to "webhook" and
then setting the `WEBHOOK_URL` to the URL to send the webhook to.

Fluxcloud will send a POST request to the provided URL with [the encoded event](https://github.com/justinbarrick/fluxcloud/blob/master/pkg/msg/msg.go) as the payload.

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
