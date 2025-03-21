local key = KEYS[1]

local window = tonumber(ARGV[1])

local threshold = tonumber(ARGV[2])

local now = tonumber(ARGV[3])

local min = now - window

redis.call('zremrangebyscore', key, "-inf", min)

local count = redis.call('zcount', key, '-inf', '+inf')

if count >= threshold then
    return "true"
else
    redis.call('zadd', key, now, now)
    redis.call('pexpire', key, window)
    return "false"
end