local key = KEYS[1]

local biz = ARGV[1]

local cnt = tonumber(ARGV[2])

local exists = redis.call("EXISTS", key)

if exists == 1 then
    redis.call("HINCRBY", key, biz, cnt)
    return 1
else
    return 0
end

