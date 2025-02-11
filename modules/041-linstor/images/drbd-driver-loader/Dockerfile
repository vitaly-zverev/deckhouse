ARG BASE_UBUNTU

FROM $BASE_UBUNTU as builder
ARG DRBD_GITREPO=https://github.com/LINBIT/drbd
ARG DRBD_VERSION=9.2.3

RUN apt-get update \
 && apt-get install -y make git \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# Copy patches
COPY ./patches /patches

# Using source code from GitHub repository
RUN git clone ${DRBD_GITREPO} /drbd \
 && cd /drbd \
 && git reset --hard drbd-${DRBD_VERSION} \
 && git apply /patches/*.patch \
 && git config --global user.email "builder@deckhouse.io" \
 && git config --global user.name "builder" \
 && git add docker/entry.sh \
 && git commit -m "changes in entry.sh for SELinux support" \
 && make tarball \
 && mv ./drbd-*.tar.gz /drbd.tar.gz

# # Using source code provided by LINBIT
# RUN DRBD_PKG=https://pkg.linbit.com//downloads/drbd/9/drbd-${DRBD_VERSION}.tar.gz \
#  && wget -O /drbd.tar.gz ${DRBD_PKG}

RUN tar xvf /drbd.tar.gz --strip-components=2 --wildcards 'drbd-*/docker/entry.sh' \
 && chmod +x /entry.sh

FROM $BASE_UBUNTU

# Debian 11 requires gcc-10
# RHEL 9 requires libelf-dev
RUN apt-get update \
 && apt-get install -y kmod gnupg wget make gcc gcc-10 patch curl libelf-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /entry.sh /drbd.tar.gz /
ENV SSL_CERT_FILE=/etc/spaas/certs/ca.crt

ENV LB_HOW compile
ENTRYPOINT [ "/entry.sh" ]
