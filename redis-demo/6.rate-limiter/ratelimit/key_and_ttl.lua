-- 如果键不存在，返回 nil
if redis.call("EXISTS", KEYS[1]) == 0 then
    return {nil, nil}
end

-- 获取键的值和剩余时间
local value = redis.call("GET", KEYS[1])
local ttl = redis.call("TTL", KEYS[1])
return {value, ttl}
