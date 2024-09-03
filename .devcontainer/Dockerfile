FROM golang:1.18-bullseye

# Install packages as root
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        sudo \
        nano \ 
        git \
        postgresql-client \
        # Needed for opening host's browser with "xdg-open https://url.com"
        xdg-utils \
        # Fix setting locale failed warning
        locales-all && \
    rm -rf /var/lib/apt/lists/* && \
    # Latest supported versions for go 1.18
    go install github.com/go-delve/delve/cmd/dlv@v1.20.2 && \
    go install golang.org/x/tools/gopls@v0.11.0 && \
    # Add write permissions for $GOPATH/pkg
    chmod -R a+w ${GOPATH}/pkg && \
    # Install command line yaml editor
    wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq && \
    chmod +x /usr/bin/yq && \
    # Install kubectl
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2

# With docker and sudo groups the host containers can be accessed from within the devcontainer
RUN groupadd docker && \
    useradd --groups sudo,docker -ms /bin/bash go-user && \
    echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
USER go-user

# Krew will be installed on this path
ENV PATH="$PATH:/home/go-user/.krew/bin" \
    # Fixes error downloading libraries from private repository
    GOPRIVATE="devstack.vwgroup.com/bitbucket"

# Install krew and oidc-login plugin for kubectl
RUN set -x; cd "$(mktemp -d)" && \
    OS="$(uname | tr '[:upper:]' '[:lower:]')" && \
    ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" && \
    KREW="krew-${OS}_${ARCH}" && \
    curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" && \
    tar zxvf "${KREW}.tar.gz" && \
    ./"${KREW}" install krew && \
    kubectl krew install oidc-login && \
    rm "${KREW}.tar.gz" "${KREW}" && \
    # Install brew and mockery
    NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" && \
    (echo; echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"') >> /home/go-user/.bashrc && \
    eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)" && \
    brew install mockery

# By default the ssl backend config is copied from Windows. Schannel is the one working for Windows, but gnutls is the one working in Linux.
RUN git config --global http.sslBackend gnutls