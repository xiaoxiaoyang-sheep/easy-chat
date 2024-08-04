kafka-topics.sh --topic msgChatTransfer --bootstrap-server 120.26.209.19:9092 --describe --exclude-internal

kafka-console-consumer.sh --bootstrap-server 120.26.209.19:9092 --topic msgChatTransfer

kafka-console-producer.sh --broker-list 120.26.209.19:9092 --topic msgChatTransfer