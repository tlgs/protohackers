# protohackers

Solutions to the [Protohackers server programming challenge](https://protohackers.com/).

## Setup

`mkdir -p bin && go build -o bin ./...`

## Deployment

The server is deployed/hosted on AWS: see the
[CloudFormation template](internal/deploy/cfn.yaml).
Additionally, a useful collection of tiny Bash scripts
is defined in `internal/deploy`.

I typically use [direnv](https://github.com/direnv/direnv) to set
project-relevant environment variables and add helper utilities to PATH.
My `.envrc` looks something like:

```bash
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=

export EC2_KEYNAME=
export CHECKER_ADDR=

PATH_add internal/deploy
```
