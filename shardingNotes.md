# Sharding in distributed systems 

## Static Sharding
divide data into shards -> diff servers

### why ? 
* all data in single server -> not good idea 

### How ?
* H(key) % n_shards =>  shard | key 
* static n_shards when change => resharding hard 
* old keys wont work with new sharding func

### Resharding
* theoretically O(n(keys)/2 * Reshard())
* downtime of system needed(bad I guess lol)
* each bucket has to split keys into other shards

### Resharding *2 
*  n_ shards  is pow(2)
*  1/2 of each shard is split into new shard
*  doubling gets bad with scale (expecially +1 key)

### Other
* range based sharding (eg : [A-P] Sh1 ; [Q-Z] Sh2 )
* Consistent Hashing (kinda ???  state stuff all painful)
* use Geography (shard location and end user location)

