#!/usr/bin/env sh
set -eu

ikuai-cli monitor system --format json
ikuai-cli network dns get
ikuai-cli network dhcp list --page 1 --page-size 20
ikuai-cli users accounts list --format json
ikuai-cli version --format json
