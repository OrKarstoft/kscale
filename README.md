# kscale

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/orkarstoft/kscale?sort=semver&style=flat-square) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/orkarstoft/kscale?style=flat-square)

kscale is an application designed to annotate Kubernetes namespaces based on Pub/Sub messages. It serves as an extension to kube-downscaler, providing additional functionality for managing forced-uptime periods within namespaces.

## Overview

Kubernetes clusters often utilize kube-downscaler to scale down namespaces and conserve resources. However, this scaling-down process can render services within a namespace temporarily inaccessible. kscale addresses this limitation by enabling the annotation of forced-uptime periods, instructing kube-downscaler to scale up the affected namespaces during specific time frames.

## Features

- Annotations based on Pub/Sub messages: kscale listens for Pub/Sub messages and expects specific attributes (`action`, `namespace`, `cluster`, and `duration`) to be present in each message.
- Forced-uptime annotations: When a relevant Pub/Sub message is received, kscale adds an annotation to the targeted namespace, specifying a forced-uptime period using the `downscaler/force-uptime` key-value pair.
- Integration with kube-downscaler: The annotations created by kscale trigger kube-downscaler to scale up the namespace during the annotated forced-uptime period, ensuring the availability of services.

## Getting Started

### Prerequisites

- Kubernetes cluster
- Google Cloud Pub/Sub
- Service Account with Pub/Sub permissions and Workload Identity enabled
- Helm

### Installation
To be added.

### Usage
1. Create an interface (e.g., a website) that allows users to trigger actions that require forced-uptime periods.
2. Configure the interface to send Pub/Sub messages with the required attributes (action, namespace, cluster, and duration) to the kscale application.
3. kscale will receive the Pub/Sub messages, process them, and add annotations to the specified namespaces accordingly.
4. kube-downscaler will respond to the annotations created by kscale, scaling up the namespaces during the forced-uptime periods.

## Contributing
Contributions are welcome! If you have any ideas, suggestions, or bug reports, please open an issue or submit a pull request.