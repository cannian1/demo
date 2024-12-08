local resource = redis.call('SRANDMEMBER', KEYS[1])
if resource then
    if redis.call('SMOVE', KEYS[1], KEYS[2], resource) == 1 then
        return resource
    end
end
return nil