# create-github-app-token

## How to use

```bash
$ make build
$ export ORG=xxx
$ export APP_ID=xxx
$ export PRIVATE_KEY=$(cat /path/to/private-key.pem)
$ export GITHUB_TOKEN=$(./create-github-app-token)
$ gh api /orgs/$ORG/repos
```
