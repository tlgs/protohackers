# protohackers

Solutions to the [Protohackers server programming challenge](https://protohackers.com/).

## Motivation

- Mess around a bit with Go and network programming
- Get some more experience managing resources using AWS

## Setup

`go build -o build/ ./...`

### AWS CLI recipes

Create stack:

```bash
aws cloudformation create-stack \
  --stack-name protohackers \
  --template-body file://cfn.yaml \
  --parameters ParameterKey=KeyName,ParameterValue="$EC2_KEY"
```

Delete stack:

```bash
aws cloudformation delete-stack --stack-name protohackers
```

Get the ID of the created EC2 instance:

```bash
aws cloudformation describe-stacks \
  | jq -r '
      .Stacks[]
      | select(.StackName == "protohackers").Outputs[]
      | select(.OutputKey == "InstanceId").OutputValue
    '
```

Start instance:

```bash
aws ec2 start-instances --instance-ids $(
  aws cloudformation describe-stacks \
    | jq -r '
        .Stacks[]
        | select(.StackName == "protohackers").Outputs[]
        | select(.OutputKey == "InstanceId").OutputValue
      '
)
```

Stop instance:

```bash
aws ec2 stop-instances --instance-ids $(
  aws cloudformation describe-stacks \
    | jq -r '
        .Stacks[]
        | select(.StackName == "protohackers").Outputs[]
        | select(.OutputKey == "InstanceId").OutputValue
      '
)
```

Get the Public DNS name from the created EC2 instance:

```bash
aws ec2 describe-instances \
  | jq -r '
      .Reservations[].Instances[]
      | select(.Tags[] | .Key == "aws:cloudformation:stack-name" and .Value == "protohackers")
      | .NetworkInterfaces[].Association.PublicDnsName
    '
```
