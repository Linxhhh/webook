
local key = KEYS[1]
local cnt = key..":cnt"
local val = ARGV[1]
local ttl = tonumber(redis.call("ttl", key))

if ttl == -2 then
    -- key 不存在，设置 key 和 cnt
    redis.call("set", key, val)
    redis.call("expire", key, 180)
    redis.call("set", cnt, 2)
    redis.call("expire", cnt, 180)
    return 0
elseif ttl == -1 then
    -- key 存在，但无过期时间，系统错误
    return -1
else
    -- key 存在
    return 1
end

-- 返回 0，表示成功响应
-- 返回 1，表示错误响应
-- 返回-1，表示系统异常