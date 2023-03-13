#!/usr/bin/env bash

set -e

function setStatus() {
    cat <<EOF > local_status.sh
export CHAIN_ID=$(echo "$OUTPUT" | awk -F'|' '/node1/{print $4}' | awk -F'/' '{print $6}')
export LOGS_PATH="$(echo "$OUTPUT" | awk -F': ' '/Node log path: /{print $2}')"
EOF
}

function showLogs() {
    if ! command -v multitail &> /dev/null
    then
        echo "multitail could not be found; please install using 'brew install multitail'"
        exit
    fi

    source local_status.sh

    multitail -D -ci magenta --label "[node1]" $(echo $LOGS_PATH | sed -e 's/<i>/1/g')/$CHAIN_ID.log \
        -ci green --label "[node2]" -I  $(echo $LOGS_PATH | sed -e 's/<i>/2/g')/$CHAIN_ID.log  \
        -ci white --label "[node3]" -I  $(echo $LOGS_PATH | sed -e 's/<i>/3/g')/$CHAIN_ID.log \
        -ci yellow --label "[node4]" -I  $(echo $LOGS_PATH | sed -e 's/<i>/4/g')/$CHAIN_ID.log \
        -ci cyan --label "[node5]" -I  $(echo $LOGS_PATH | sed -e 's/<i>/5/g')/$CHAIN_ID.log

}
