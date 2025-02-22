# Copyright 2021 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# policycoreutils-python libseccomp - containerd.io dependencies
SYSTEM_PACKAGES="curl wget virt-what bash-completion lvm2 parted sudo yum-utils nfs-utils tar xz device-mapper-persistent-data net-tools libseccomp checkpolicy"
KUBERNETES_DEPENDENCIES="conntrack-tools ebtables ethtool iproute iptables socat util-linux"
if bb-is-centos-version? 7; then
  SYSTEM_PACKAGES="${SYSTEM_PACKAGES} policycoreutils-python"
fi
if bb-is-centos-version? 8; then
  SYSTEM_PACKAGES="${SYSTEM_PACKAGES} policycoreutils-python-utils libcgroup"
fi
if bb-is-centos-version? 9; then
  SYSTEM_PACKAGES="${SYSTEM_PACKAGES} policycoreutils-python-utils"
fi
# yum-plugin-versionlock is needed for bb-yum-install
if yum --version | grep -q dnf; then
  bb-yum-install python3-dnf-plugin-versionlock
else
  bb-yum-install yum-plugin-versionlock
fi

bb-yum-install ${SYSTEM_PACKAGES} ${KUBERNETES_DEPENDENCIES}

bb-rp-install "jq:{{ .images.registrypackages.jq16 }}" "curl:{{ .images.registrypackages.d8Curl801 }}"

if bb-is-centos-version? 8; then
  bb-rp-install "inotify-tools:{{ .images.registrypackages.inotifyToolsCentos831419 }}"
fi
if bb-is-centos-version? 9; then
  bb-rp-install "inotify-tools:{{ .images.registrypackages.inotifyToolsCentos9322101 }}"
fi
bb-yum-remove yum-cron
