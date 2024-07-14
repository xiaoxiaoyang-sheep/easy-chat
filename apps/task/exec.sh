kafka-topics.sh --topic msgChatTransfer --bootstrap-server 116.198.246.212:9092 --describe --exclude-internal

kafka-console-consumer.sh --bootstrap-server 116.198.246.212:9092 --topic msgChatTransfer

kafka-console-producer.sh --broker-list 116.198.246.212:9092 --topic msgChatTransfer