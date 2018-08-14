Fluxcloud is a tool to receive events from the [Weave flux](https://github.com/weaveworks/flux).

![build status](https://ci.codesink.net/api/badges/justinbarrick/fluxcloud/status.svg)
[![image version](https://images.microbadger.com/badges/version/justinbarrick/fluxcloud.svg)](https://microbadger.com/images/justinbarrick/fluxcloud)
[![image size](https://images.microbadger.com/badges/image/justinbarrick/fluxcloud.svg)](https://microbadger.com/images/justinbarrick/fluxcloud "Get your own image badge on microbadger.com")

Weave Flux is a useful tool for managing the state of your Kubernetes cluster.

Fluxcloud is a valid upstream for Weave, allowing you to send Flux events to Slack
without using Weave Cloud.

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

And then apply the configuration:

```
kubectl apply -f examples/fluxcloud.yaml
```

Set the `--connect` flag on Flux to `--connect=ws://fluxcloud`.

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
