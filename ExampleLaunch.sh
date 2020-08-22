set -e

trap 'killall abcd' SIGINT

cd $(dirname $0)

killall abcd || true
sleep 0.1

go install -v

abcd -db-location=./sh0.db -http-addr=127.0.0.2:8080 -config-file=sharding.toml -shard=sh0 &
abcd -db-location=./sh0-replica.db -http-addr=127.0.0.2:8080 -config-file=sharding.toml -shard=sh0 -replica &
abcd -db-location=./sh1.db -http-addr=127.0.0.3:8080 -config-file=sharding.toml -shard=sh1 &
abcd -db-location=./sh1-replica.db -http-addr=127.0.0.2:8080 -config-file=sharding.toml -shard=sh1 -replica &
abcd -db-location=./sh2.db -http-addr=127.0.0.4:8080 -config-file=sharding.toml -shard=sh2 &
abcd -db-location=./sh2-replica.db -http-addr=127.0.0.2:8080 -config-file=sharding.toml -shard=sh2 -replica &
abcd -db-location=./sh3.db -http-addr=127.0.0.5:8080 -config-file=sharding.toml -shard=sh3 &
abcd -db-location=./sh3-replica.db -http-addr=127.0.0.2:8080 -config-file=sharding.toml -shard=sh3 -replica &



wait


