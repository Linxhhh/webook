local key = KEYS[1]
local cntKey = ARGV[1]
local delta = tonumber(ARGV[2]) -- 表示 +1 或 -1
local exists = redis.call("exists", key)

if exists == 1 then
    -- key 存在，更新阅读计数
    redis.call("hincrby", key, cntKey, delta)
    return 1
else
    -- key 不存在
    return 0
end