webserver
receives messages sent by Step7Sender
for each message received Step7Server calls the grpc-server implemented in step6, which saves the messages in a file
clients of type Step7Client can register themselves such that a websocket connection is created
after registration the Step7Client client will receive the above mentioned messages via a websocket connection
