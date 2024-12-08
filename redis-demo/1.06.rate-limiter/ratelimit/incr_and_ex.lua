local count = redis.call('incr', KEYS[1]);
-- 如果是第一次调用，设置过期时间
if count == 1 then
    redis.call('expire', KEYS[1], ARGV[1]);
end
-- 返回最新计数
return count;