local key = KEYS[1]
local input_code = ARGV[2]
local cnt_key = key.. "_cnt"

local code = redis.call("get", key)
local cnt = tonumber(redis.call("get", cnt_key))

if cnt == nil or cnt <= 0 then
    return -1
end

if input_code == code then
    redis.call("set", cnt_key, 0)
    return 0
else
    redis.call("decr", cnt_key)
    return -2
end

