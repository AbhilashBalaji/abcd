set -e

trap 'killall abcd' SIGINT

cd $(dirname $0)

killall abcd || true
sleep 0.1

go install -v

abcd -db-location=./sh0.db -http-addr=127.0.0.1:8080 -config-file=sharding.toml -shard=sh0 &
abcd -db-location=./sh1.db -http-addr=127.0.0.1:8081 -config-file=sharding.toml -shard=sh1 &
abcd -db-location=./sh2.db -http-addr=127.0.0.1:8082 -config-file=sharding.toml -shard=sh2 &

wait


