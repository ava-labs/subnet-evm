#!/usr/bin/env bash

set -e

function setStatus() {
    cat <<EOF >local_status.sh
export CHAIN_ID=$(echo "$OUTPUT" | awk -F'|' '/node1/{print $4}' | awk -F'/' '{print $6}')
export LOGS_PATH="$(echo "$OUTPUT" | awk -F': ' '/Node log path: /{print $2}')"
EOF

    cat <<EOF >~/.hubblenet.json
{
    "chain_id": "$(echo "$OUTPUT" | awk -F'|' '/node1/{print $4}' | awk -F'/' '{print $6}')"
}
EOF
}

function showLogs() {
    if ! command -v multitail &>/dev/null; then
        echo "multitail could not be found; please install using 'brew install multitail'"
        exit
    fi

    source local_status.sh
    if [ -z "$1" ]; then
        # Define colors and nodes

        colors=("magenta" "green" "white" "yellow" "cyan")
        nodes=("1" "2" "3" "4" "5")

        # Use for loop to iterate through nodes
        for index in ${!nodes[*]}
        do
            node=${nodes[$index]}
            color=${colors[$index]}
            logs_path=$(echo $LOGS_PATH | sed -e "s/<i>/$node/g")
            # Add multitail command for each node
            cmd_part+=" -ci $color --label \"[node$node]\" -I ${logs_path}/$CHAIN_ID.log"
        done

        # Execute multitail with the generated command parts
        eval "multitail -D $cmd_part"

    else
        if [ -z "$2" ]; then
            # from the beginning
            tail -f -n +1 "${LOGS_PATH/<i>/$1}/$CHAIN_ID.log"
        else
            grep --color=auto -i "$2" "${LOGS_PATH/<i>/$1}/$CHAIN_ID.log"
        fi
    fi

}
