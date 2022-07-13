FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12 AS builder
COPY . .
RUN go build -o gitops-commit-status -mod=readonly .

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=builder /opt/app-root/src/gitops-commit-status /usr/local/bin
