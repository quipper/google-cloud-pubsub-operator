# google-cloud-pubsub-operator [![go](https://github.com/quipper/google-cloud-pubsub-operator/actions/workflows/go.yaml/badge.svg)](https://github.com/quipper/google-cloud-pubsub-operator/actions/workflows/go.yaml)

The Google Cloud Pub/Sub Operator is a Kubernetes operator designed to efficiently manage Google Cloud Pub/Sub Topics and Subscriptions.

Google Cloud Pub/Sub is a real-time messaging service provided by Google Cloud Platform that efficiently delivers large volumes of messages from producers to consumers.
This operator allows you to manage Pub/Sub resources within your Kubernetes cluster, representing them as Kubernetes resources.

## Installation

To install this operator,

```sh
kubectl apply -f https://github.com/quipper/google-cloud-pubsub-operator/releases/download/v0.12.0/google-cloud-pubsub-operator.yaml
```

## Usage

### Creating a Topic

```yaml
apiVersion: googlecloudpubsuboperator.quipper.github.io/v1
kind: Topic
metadata:
  name: your-topic
spec:
  projectID: your-project
  topicID: your-topic
```

### Creating a Subscription

```yaml
apiVersion: googlecloudpubsuboperator.quipper.github.io/v1
kind: Subscription
metadata:
  name: your-subscription
spec:
  subscriptionProjectID: your-project
  subscriptionID: your-subscription

  topicProjectID: your-project
  topicID: your-topic
```

## Contributing

Contributions to this project are welcome!
For instructions on how to contribute, please see [CONTRIBUTING.md](CONTRIBUTING.md).
