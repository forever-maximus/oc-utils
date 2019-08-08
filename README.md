# Openshift CLI tool

This is a CLI tool created to add extra functionality for managing Openshift.

It is written in golang and uses the Openshift rest API.

## Commands

There are three available commands:
* scaleup - this will scale up all pods in a given namespace.
* scaledown - this will scale down all pods in a given namespace.
* restartpods - this will restart all pods older than a given threshold in a namespace. Threshold is measured in days.

## Flags

There is one flag:
* --prod - this can be added to point the tool to Openshift prod environment. The default is non-prod.

## Examples

Usage: ./executable [flags] command [args ...]

### Scale up
```
./oc-utils scaleup integration-dev
```

### Scale down
```
./oc-utils scaledown integration-dev
```

### Restart Pods
```
./oc-utils restartpods integration-dev 7
```