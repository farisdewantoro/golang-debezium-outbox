curl -s http://localhost:8083/connectors/outbox-connector/config | jq
curl -s http://localhost:8083/connectors/outbox-connector/status | jq
curl -X DELETE http://localhost:8083/connectors/outbox-connector
curl -X POST -H "Content-Type: application/json" --data @outbox.json http://localhost:8083/connectors
curl -X POST http://localhost:8083/connectors/outbox-connector/restart

kafka-console-consumer --bootstrap-server localhost:9092 --topic dbserver1.public.outbox_events --from-beginning

docker exec -it go-outbox-debezium-kafka bash 
kafka-topics --create --topic dbz.public.outbox_events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 || true &&
kafka-topics --create --topic eventdriven-examples.dlq --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 || true &&
kafka-topics --create --topic notif.user.registration --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 || true &&