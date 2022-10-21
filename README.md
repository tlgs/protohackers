# protohackers

Solutions to the [Protohackers](https://protohackers.com/) server programming challenges:

0. [Smoke Test](cmd/smoke-test/main.go) -
   Echo service from [RFC 862](https://www.rfc-editor.org/rfc/rfc862.html) (TCP).
1. [Prime Time](cmd/prime-time/main.go) -
   Primality testing service using a JSON-based response-request protocol (TCP).
2. [Means to an End](cmd/means-to-an-end/main.go) -
   Asset price tracking service using a binary message format (TCP).
3. [Budget Chat](cmd/budget-chat/main.go) -
   Simple chat room (TCP).
4. [Unusual Database Program](cmd/unusual-database-program/main.go) -
   Key-value store (UDP).
5. [Mob in the Middle](cmd/mob-in-the-middle/main.go) -
   Man-in-the-middle attack on _Budget Chat_ (TCP).

## Setup

```console
$ mkdir -p bin
$ go build -o bin ./...
```

## Deployment

The server is deployed/hosted on AWS: see the
[CloudFormation template](internal/deploy/cfn.yaml).
Additionally, a useful collection of tiny Bash scripts
is defined in `internal/deploy`.

I typically use [direnv](https://github.com/direnv/direnv) to automatically set
project-relevant environment variables and add helper utilities to PATH.
My `.envrc` looks something like:

```bash
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=

export EC2_KEYNAME=
export CHECKER_ADDR=

PATH_add internal/deploy
PATH_add bin
```
