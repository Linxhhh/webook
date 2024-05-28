local key = KEYS[1]
local cntKey = key..":cnt"
local inputCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    return -1
end

local code = redis.call("get", key)
if code == inputCode then
    -- 验证通过，把 cntKey 置为无效
    redis.call("set", cntKey, -1)
    return 0
else
    -- 验证失败，把 cntKey 减一
    redis.call("decr", cntKey)
    return 1
end

-- 返回 0，表示成功响应
-- 返回 1，表示错误响应
-- 返回-1，表示验证频繁