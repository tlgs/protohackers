# protohackers

## Motivation

- Mess around a bit with Go and network programming
- Get some more experience managing resources using AWS

## Setup

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
