# gitops-commit-status

A simple Go program to set commit status. It uses [go-scm](https://github.com/jenkins-x/go-scm) to support multiple Git providers(GitHub, GitLab, etc).



## Usage

1. Run it locally

```shell
$ go build -o gitops-commit-status .

$ ./gitops-commit-status --sha <commit_sha> --url <git_repo_url> --token <token> --status <commit_status> --context <commit_status>
```

2. Tekton 

Inlcude it as part of your Tekton pipeline and set the commit status based on PipelineRun status. An example task and taskrun can be found in [tekton](tekton/).

```shell
$ kubectl create secret generic git-host-access-token --from-literal=token=<token>

$ kubectl apply -f tekton/
```
