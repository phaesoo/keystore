package memdb

import "github.com/gomodule/redigo/redis"

// HSETX only updates a field in a hash if it already exists.
var HSETX = redis.NewScript(1, `
if redis.call('hexists', KEYS[1], ARGV[1]) == 1 then
	redis.call('hset', KEYS[1], ARGV[1], ARGV[2])
	return 1
else
	return 0
end`)

// HDELZEROINCRBY increments a field in a hash by an amount, and deletes the key if it is zero.
var HDELZEROINCRBY = redis.NewScript(1, `
local res = tonumber(redis.call('HINCRBY', KEYS[1], ARGV[1], ARGV[2]))
if res == 0 then
	redis.call('HDEL', KEYS[1], ARGV[1])
end
return res
`)

// PopAndPushSortedSet requires 2 keys and 1 argument.
// KEYS[1]: source set
// KEYS[2]: dest set
// ARGV[1]: number of members to pop
// It pops {argv[1]} members from sources set and put it into dest set.
var PopAndPushSortedSet = redis.NewScript(2, `
local popped = redis.call('ZPOPMIN', KEYS[1], ARGV[1])
local res = {}
local i = 1
while i <= table.getn(popped) do
	local elem, score = popped[i], popped[i+1]
	table.insert(res, elem)
	redis.call('ZADD', KEYS[2], score, elem)
	i = i + 2
end
return res
`)

// PopAndPushSortedSetWithMember requires 2 keys and 2 arguments.
// KEYS[1]: source set
// KEYS[2]: dest set
// ARGV[1]: member
// ARGV[2]: new score for member
// It pops target member from source set and put into dest set with new score.Vy
var PopAndPushSortedSetWithMember = redis.NewScript(2, `
local numPopped = redis.call('ZREM', KEYS[1], ARGV[1])
if numPopped == 0 then
	return false
end
redis.call('ZADD', KEYS[2], ARGV[2], ARGV[1])
return true
`)

// MoveSetMembersWithinScore requires 2 keys and 3 arguments.
// KEYS[1]: source set
// KEYS[2]: dest set
// ARGV[1]: min score
// ARGV[2]: max score
// ARGV[3]: score incremental for backoff
// It pops members with score range(min <= x <= max) from source set and put it into dest set with
// each score incremented by {argv[3]}.
var MoveSetMembersWithinScore = redis.NewScript(2, `
local moved = redis.call('ZRANGEBYSCORE', KEYS[1], ARGV[1], ARGV[2], 'WITHSCORES')
redis.call('ZREMRANGEBYSCORE', KEYS[1], ARGV[1], ARGV[2])
local numMoved = table.getn(moved)/2
local i = 1
while i <= table.getn(moved) do
	local elem, score = moved[i], moved[i+1]
	local newScore = score + ARGV[3]
	redis.call('ZADD', KEYS[2], newScore, elem)
	i = i + 2
end
return numMoved
`)
