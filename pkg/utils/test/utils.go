package test_utils

import (
	"bytes"

	"github.com/topfreegames/fluxcloud/pkg/utils"
	fluxevent "github.com/weaveworks/flux/event"
)

func NewFluxSyncEvent() fluxevent.Event {
	event, _ := utils.ParseFluxEvent(bytes.NewBufferString(`{
    "id": 0,
    "serviceIDs": [
        "default:deployment/test"
    ],
    "type": "sync",
    "startedAt": "2018-07-07T02:45:35.247153012Z",
    "endedAt": "2018-07-07T02:45:35.247153012Z",
    "logLevel": "info",
    "metadata": {
        "commits": [
            {
                "revision": "810c2e6f22ac5ab7c831fe0dd697fe32997b098f",
                "message": "change test image"
            }
        ],
        "includes": {
            "other": true
        }
    }
}`))

	return event
}

func NewFluxSyncErrorEvent() fluxevent.Event {
	event, _ := utils.ParseFluxEvent(bytes.NewBufferString(`{
  "id": 0,
  "serviceIDs": [
    "default:persistentvolumeclaim/test"
  ],
  "type": "sync",
  "startedAt": "2018-09-05T01:44:17.427541601Z",
  "endedAt": "2018-09-05T01:44:17.427541601Z",
  "logLevel": "info",
  "metadata": {
    "commits": [
      {
        "revision": "4997efcd4ac6255604d0d44eeb7085c5b0eb9d48",
        "message": "create invalid resource"
      }
    ],
    "includes": {
      "other": true
    },
    "errors": [
      {
        "ID": "default:persistentvolumeclaim/test",
        "Path": "manifests/test.yaml",
        "Error": "running kubectl: The PersistentVolumeClaim \"test\" is invalid: spec: Forbidden: field is immutable after creation"
      },
      {
        "ID": "default:persistentvolumeclaim/lol",
        "Path": "manifests/lol.yaml",
        "Error": "running kubectl: The PersistentVolumeClaim \"lol\" is invalid: spec: Forbidden: field is immutable after creation"
      }
    ]
  }
}`))

	return event
}

func NewFluxCommitEvent() fluxevent.Event {
	event, _ := utils.ParseFluxEvent(bytes.NewBufferString(`{
    "id": 0,
    "serviceIDs": [
        "default:deployment/test"
    ],
    "type": "commit",
    "startedAt": "2018-07-07T03:02:16.042202426Z",
    "endedAt": "2018-07-07T03:02:16.042202426Z",
    "logLevel": "info",
    "metadata": {
        "revision": "d644e1a05db6881abf0cdb78299917b95f442036",
        "spec": {
            "type": "policy",
            "cause": {
                "Message": "",
                "User": "Justin Barrick \u003cjustin.m.barrick@gmail.com\u003e"
            },
            "spec": {
                "default:deployment/test": {
                    "add": {
                        "automated": "true"
                    },
                    "remove": {}
                }
            }
        },
        "result": {
            "default:deployment/test": {
                "Status": "success",
                "PerContainer":null
            }
        }
    }
}`))

	return event
}

func NewFluxAutoReleaseEvent() fluxevent.Event {
	event, _ := utils.ParseFluxEvent(bytes.NewBufferString(`{
    "id": 0,
    "serviceIDs": [
        "default:deployment/test"
    ],
    "type": "autorelease",
    "startedAt": "2018-07-07T03:29:28.419542197Z",
    "endedAt": "2018-07-07T03:29:29.403503538Z",
    "logLevel": "info",
    "metadata": {
        "Revision": "4d030af4f8e4af14ae35154483b1355bdfeefb73",
        "result": {
            "default:deployment/test": {
                "Status": "success",
                "PerContainer": [
                    {
                        "Container": "test2",
                        "Current": "justinbarrick/nginx:test1",
                        "Target": "justinbarrick/nginx:test3"
                    }
                ]
            }
        },
        "spec": {
            "Changes": [
                {
                    "ServiceID": "default:deployment/test",
                    "Container": {
                        "Name": "test2",
                        "Image": "justinbarrick/nginx:test1"
                    },
                    "ImageID": "justinbarrick/nginx:test3"
                }
            ]
        }
    }
}`))

	return event
}

func NewFluxUpdatePolicyEvent() fluxevent.Event {
	event, _ := utils.ParseFluxEvent(bytes.NewBufferString(`{
    "id": 0,
    "serviceIDs": [
        "default:deployment/test"
    ],
    "type": "sync",
    "startedAt": "2018-07-07T03:02:24.611208878Z",
    "endedAt": "2018-07-07T03:02:24.611208878Z",
    "logLevel": "info",
    "metadata": {
        "commits": [
            {
                "revision": "d644e1a05db6881abf0cdb78299917b95f442036",
                "message": "Automated: default:deployment/test"
            }
        ],
        "includes": {
            "update_policy": true
        }
    }
}`))
	return event
}
