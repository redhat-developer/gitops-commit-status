FROM openshift/origin-release:golang-1.16 AS builder

WORKDIR /tmp/gitops-commit-status
COPY . /tmp/gitops-commit-status
RUN go build -o gitops-commit-status -mod=readonly .

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=builder /tmp/gitops-commit-status/gitops-commit-status /usr/local/bin
