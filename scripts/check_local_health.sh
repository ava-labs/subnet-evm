#!/usr/bin/env bash
set -e

source local_status.sh

apis=(
  "http://127.0.0.1:9650/ext/bc/$CHAIN_ID/rpc"
  "http://127.0.0.1:9652/ext/bc/$CHAIN_ID/rpc"
  "http://127.0.0.1:9654/ext/bc/$CHAIN_ID/rpc"
  "http://127.0.0.1:9656/ext/bc/$CHAIN_ID/rpc"
  "http://127.0.0.1:9658/ext/bc/$CHAIN_ID/rpc"
)

# Flag to track if any API took longer than 1 second to respond
error_flag=false

# Loop through each API endpoint
for api in "${apis[@]}"; do
  # Call the API endpoint with a timeout of 1 second
  if ! curl --connect-timeout 1 --max-time 1 -s "$api" > /dev/null; then
    echo "API $api did not respond within 1 second."
    error_flag=true
  fi
done

# Check if any API took longer to respond
if [ "$error_flag" = true ]; then
  echo "Error: One or more APIs did not respond within 1 second."
else
  echo "OK: All APIs responded within 1 second."
fi
