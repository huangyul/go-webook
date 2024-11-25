-- key
local key = KEYS[1]

-- window size
local window = tonumber(ARGV[1])

-- rate
local rate = tonumber(ARGV[2])

-- now
local now = tonumber(ARGV[3])

local min = now - window

redis.call('ZREMRANGEBYSCORE', key, '-inf', min)

local count = redis.call('ZCOUNT', key, '-inf', '+inf')

if count >= rate then
    return "true"
else
    redis.call('ZADD', key, now, now)
    redis.call('PEXPIRE', key, window)
    return "false"
end