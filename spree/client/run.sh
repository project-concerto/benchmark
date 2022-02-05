#!/usr/bin/env bash

set -ex

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
SPREE_DIRECTORY="${SCRIPT_DIR}/../spree/spree"
PATCHES_DIRECTORY="${SCRIPT_DIR}/patches/checkout-v1/"

function original() {
	pushd $SPREE_DIRECTORY
	git reset --hard
	popd
}

function serializable() {
	pushd $SPREE_DIRECTORY
	git reset --hard
	git apply ${PATCHES_DIRECTORY}/serializable-remove.patch
	popd
}

function lock() {
	pushd $SPREE_DIRECTORY
	git reset --hard
	git apply ${PATCHES_DIRECTORY}/lock-remove.patch
	popd
}

for i in $(seq 1 1); do
	serializable
	THREADS=1,2,4,8,16,32,48,64,80,96,108,120 TIMEOUT=60 go run ./spree-checkout/bin/main.go -file-prefix=DBT -command=benchmark 2&> $(date +%s).log

	THREADS=1,2,4,8,16,32,48,64,80,96,108,120 TIMEOUT=60 go run ./spree-checkout/bin/main.go -file-prefix=NCDBT -command=benchmark-no-contention 2&> $(date +%s).log

	lock
	THREADS=1,2,4,8,16,32,48,64,80,96,108,120 TIMEOUT=60 go run ./spree-checkout/bin/main.go -file-prefix=AHT -command=benchmark 2&> $(date +%s).log

	THREADS=1,2,4,8,16,32,48,64,80,96,108,120 TIMEOUT=60 go run ./spree-checkout/bin/main.go -file-prefix=NCAHT -command=benchmark-no-contention &> $(date +%s).log
done

