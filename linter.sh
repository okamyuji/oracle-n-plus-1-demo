#!/usr/bin/env bash

set -euo pipefail

# linter.sh
# - gofmt チェック（必要なら --fix で整形）
# - go vet
# - golangci-lint（ローカル/コンテナ/セルフインストールの順で実行）
# 参考: https://github.com/golangci/golangci-lint

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT_DIR"

HAVE_DOCKER=0
if command -v docker >/dev/null 2>&1; then
    HAVE_DOCKER=1
fi

MODE_FIX=0
if [[ "${1:-}" == "--fix" ]]; then
    MODE_FIX=1
fi

# OSごとの cgo 設定（macOS は Homebrew lame を考慮）
export CGO_ENABLED=${CGO_ENABLED:-1}
if [[ "$(uname)" == "Darwin" ]]; then
    if command -v brew >/dev/null 2>&1 && brew --prefix lame >/dev/null 2>&1; then
        export CGO_CFLAGS="-I$(brew --prefix lame)/include"
        export CGO_LDFLAGS="-L$(brew --prefix lame)/lib"
    fi
fi

# 1) gofmt チェック/整形
if [[ ${MODE_FIX} -eq 1 ]]; then
    echo "[fmt] applying gofmt -s -w ..."
    gofmt -s -w .
else
    echo "[fmt] checking gofmt -s -l ..."
    fmt_out=$(gofmt -s -l . | grep -v '^vendor/' || true)
    if [[ -n "${fmt_out}" ]]; then
        echo "gofmt failed for the following files:" >&2
        echo "${fmt_out}" >&2
        echo "Run: ./linter.sh --fix" >&2
        exit 1
    fi
fi

# 2) go vet
echo "[vet] running go vet ..."
GOFLAGS=${GOFLAGS:-} go vet ./...

# 3) golangci-lint
GCL_VERSION="v1.59.1"
GCL_ARGS=("run" "--timeout=5m")
if [[ ${MODE_FIX} -eq 1 ]]; then
    GCL_ARGS+=("--fix")
fi

run_golangci() {
    echo "[lint] running golangci-lint ${GCL_VERSION} ..."
    golangci-lint "${GCL_ARGS[@]}"
}

run_golangci_local_bin() {
    echo "[lint] running local ./bin/golangci-lint ${GCL_VERSION} ..."
    "${ROOT_DIR}/bin/golangci-lint" "${GCL_ARGS[@]}"
}

run_golangci_docker() {
    echo "[lint] running docker golangci-lint ${GCL_VERSION} ..."
    # macOS の場合は brew の lame をコンテナにマウントしてヘッダ/ライブラリを解決
    local extra_mounts=()
    if [[ "$(uname)" == "Darwin" ]] && command -v brew >/dev/null 2>&1; then
        local BREW_PREFIX
        BREW_PREFIX=$(brew --prefix lame 2>/dev/null || true)
        if [[ -n "${BREW_PREFIX}" ]]; then
            extra_mounts+=("-v" "${BREW_PREFIX}/include:${BREW_PREFIX}/include")
            extra_mounts+=("-v" "${BREW_PREFIX}/lib:${BREW_PREFIX}/lib")
        fi
    fi
    docker run --rm \
        -e CGO_ENABLED -e CGO_CFLAGS -e CGO_LDFLAGS \
        -v "${ROOT_DIR}:/app" -w /app \
        "${extra_mounts[@]}" \
        golangci/golangci-lint:${GCL_VERSION} \
        golangci-lint "${GCL_ARGS[@]}"
}

ensure_golangci_local_bin() {
    mkdir -p "${ROOT_DIR}/bin"
    if [[ ! -x "${ROOT_DIR}/bin/golangci-lint" ]]; then
        echo "[lint] installing golangci-lint ${GCL_VERSION} to ./bin ..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
          | sh -s -- -b "${ROOT_DIR}/bin" ${GCL_VERSION}
    fi
}

# 実行戦略: ローカルバイナリ > セルフインストール > docker（cgo依存解決のため）
if command -v golangci-lint >/dev/null 2>&1; then
    run_golangci
elif [[ -x "${ROOT_DIR}/bin/golangci-lint" ]]; then
    run_golangci_local_bin
elif [[ ${HAVE_DOCKER} -eq 1 ]]; then
    run_golangci_docker || { echo "[lint] docker 実行に失敗。ローカルへ切替えます"; ensure_golangci_local_bin; run_golangci_local_bin; }
else
    ensure_golangci_local_bin
    run_golangci_local_bin
fi

echo "All linters passed."