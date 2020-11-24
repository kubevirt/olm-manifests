FROM golang:1.15.2 AS builder

WORKDIR /go/src/github.com/kubevirt/hyperconverged-cluster-operator/
COPY . .
RUN make build-webhook

FROM registry.access.redhat.com/ubi8/ubi-minimal
ENV WEBHOOK=/usr/local/bin/hyperconverged-cluster-webhook \
    APP=WEBHOOK \
    USER_UID=1001 \
    USER_NAME=hyperconverged-cluster-webhook

COPY --from=builder /go/src/github.com/kubevirt/hyperconverged-cluster-operator/build/bin/ /usr/local/bin

# ensure $HOME exists and is accessible by group 0 (we don't know what the runtime UID will be)
RUN mkdir -p ${HOME} && \
    chown ${USER_UID}:0 ${HOME} && \
    chmod ug+rwx ${HOME} && \
    # runtime user will need to be able to self-insert in /etc/passwd
    chmod g+rw /etc/passwd

COPY --from=builder /go/src/github.com/kubevirt/hyperconverged-cluster-operator/_out/hyperconverged-cluster-webhook $APP
ENTRYPOINT ["/usr/local/bin/entrypoint"]
USER ${USER_UID}

ARG git_url=https://github.com/kubevirt/hyperconverged-cluster-operator.git
ARG git_sha=NONE

LABEL multi.GIT_URL=${git_url} \
      multi.GIT_SHA=${git_sha}
