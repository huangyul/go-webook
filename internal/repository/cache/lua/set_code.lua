-- set success: 0
-- set too frequently: -1
-- not expiration time: -2

local key = KEYS[1]
local code = ARGV[1]
local cnt_key = key .. "_cnt"

local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, code)
    redis.call("set", cnt_key, 1)
    redis.call("expire", key, 600)
    redis.call("expire", cnt_key, 600)
    return 0
else
    return -1
end