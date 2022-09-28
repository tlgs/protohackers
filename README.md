# protohackers

Solutions to the [Protohackers server programming challenge](https://protohackers.com/).

## Setup

`mkdir -p bin && go build -o bin ./...`

## Deployment

The server is deployed/hosted on AWS: see the
[CloudFormation template](internal/deploy/cfn.yaml).

A useful collection of tiny Bash scripts is defined in `internal/deploy/`.
These can be automatically loaded into the current shell by using
[direnv](https://github.com/direnv/direnv) and a `.envrc` file like:

```bash
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=

export EC2_KEYNAME=
export CHECKER_ADDR=

PATH_add internal/deploy
```
