#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source_file="${script_dir}/.env"

if [[ ! -f "${source_file}" ]]; then
  echo "Создайте deploy/env/.env из deploy/env/.env.example"
  exit 1
fi

if ! command -v envsubst >/dev/null 2>&1; then
  echo "Для генерации нужен envsubst"
  exit 1
fi

set -a
source "${source_file}"
set +a

IFS=',' read -ra services <<< "${SERVICES}"
for service in "${services[@]}"; do
  template="${script_dir}/${service}.env.template"
  output="${script_dir}/../compose/${service}/.env"
  mkdir -p "$(dirname "${output}")"
  envsubst < "${template}" > "${output}"
  echo "Создан ${output}"
done
