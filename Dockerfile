FROM openshift/origin-release:golang-1.16 AS builder

WORKDIR /tmp/set-commit-status
COPY . /tmp/set-commit-status
RUN go build -o set-commit-status -mod=readonly .

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=builder /tmp/set-commit-status/set-commit-status /usr/local/bin
