local key = KEYS[1]

local cntKey = key .. ":cnt"

local code = ARGV[1]

local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- expire time not exist
    return -2
elseif ttl == -2 or ttl < 540 then
    --- set new code
    redis.call("SET", key, code)
    redis.call("EXPIRE", key, 600)
    redis.call("SET", cntKey, 3)
    redis.call("EXPIRE", cntKey, 600)
    return 0
else
    -- send too frequent
    return -1
end

