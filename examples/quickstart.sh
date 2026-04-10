#!/usr/bin/env sh
set -eu

ikuai-cli auth set-url https://192.168.1.1
ikuai-cli auth set-token '<your-api-token>'

ikuai-cli auth status
ikuai-cli monitor system
ikuai-cli network wan
ikuai-cli system get
