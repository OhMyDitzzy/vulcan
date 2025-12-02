#!/bin/bash

# Start three interconnected nodes
echo "Starting Vulcan nodes..."

# Node 1
./bin/vulcan --api-port=8080 --port=6000 --db-path=./data/node1 &
sleep 2

# Node 2 (connects to Node 1)
./bin/vulcan --api-port=8081 --port=6001 --db-path=./data/node2 --peers=localhost:6000 &
sleep 2

# Node 3 (connects to Node 1 and 2)
./bin/vulcan --api-port=8082 --port=6002 --db-path=./data/node3 --peers=localhost:6000,localhost:6001 &

echo "Nodes started!"
echo "Node 1: API=8080, P2P=6000"
echo "Node 2: API=8081, P2P=6001"
echo "Node 3: API=8082, P2P=6002"