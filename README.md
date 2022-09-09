# protohackers

## Motivation

- Learn a bit of network programming
- Get some more experience managing resources using AWS

## Setup

Validate template:

Create stack:

```
aws cloudformation create-stack \
  --stack-name protohackers \
  --template-body file://cfn.yaml \
  --parameters ParameterKey=KeyName,ParameterValue=$USER-default
```

Delete stack:

```
aws cloudformation delete-stack --stack-name protohackers
```
