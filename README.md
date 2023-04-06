# google-cloud-pubsub-operator [![go](https://github.com/quipper/google-cloud-pubsub-operator/actions/workflows/go.yaml/badge.svg)](https://github.com/quipper/google-cloud-pubsub-operator/actions/workflows/go.yaml)

This is a Kubernetes operator for Google Cloud Pub/Sub Topic and Subscription.

## Getting Started

To create a Topic:

```yaml
apiVersion: googlecloudpubsuboperator.quipper.github.io/v1
kind: Topic
metadata:
  name: your-topic
spec:
  projectID: your-project
  topicID: your-topic
```

To create a Subscription:

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
