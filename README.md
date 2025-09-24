# asset-discovery
Due to time constraint couldn't complete all given test cases. 
But below are commands to test each stage.

To start first node ->  go run main.go --port=8081
can add multiple nodes with peers   -> go run main.go --port=8082 --peers=localhost:8081
                                    -> go run main.go --port=8083 --peers=localhost:8081

To view peers in each node          -> http://localhost:8081/peers (Can change port for other nodes)
To join node with API               -> http://localhost:8081/join 
Payload                             -> {
                                            "id": "localhost:8085",
                                            "address": "http://localhost:8085",
                                            "last_seen": "2025-09-24T10:54:41.339733127+05:30",
                                            "alive": true
                                        }
Check if node is active (Heartbeat) -> http://localhost:8084/heartbeat

<--------->
Initially started with basic API with starting service on port and adding peers in map.
Later added logic to start each node with different port as peers logic already completed in first step.
Added health check logic as we already had map with all peers. created simple healthcheck API which was called from parent node.



# counter
counter is continously synchronised with all peers in each node.
To add counter in node              ->  http://localhost:8083/counter/increment
To view counter value               ->  http://localhost:8083/counter/value
To verify if counter syncronised    ->  http://localhost:8081/counter/value

Due to time constraint couldn't complete test suit.

<-------------------------------------->
Started with basic API to increment, decrement counter in local variable.
Created sync API which will update provided counter for given node.
Added boradcast logic which was called everytime we manupulate counter.
