import requests
import json
import time
import random

# === 接口地址（根据实际部署调整） ===
LOGIN_URL = "http://localhost:8961/login"
JOB_ADD_URL = "http://localhost:8961/job/add"

# === 登录信息（替换为实际账号密码） ===
login_data = {
    "username": "test_user_001",
    "password": "Test123456"  # 必须填写正确密码
}

# === 1. 登录获取 Token ===
response = requests.post(LOGIN_URL, json=login_data)
if response.status_code != 200:
    print(f"登录失败: {response.text}")
    exit(1)

login_result = response.json()
if login_result.get("code") != 200:
    print(f"登录失败: {login_result.get('msg')}")
    exit(1)

token = login_result["data"]["token"]
headers = {
    "Authorization": token,  # 无需 Bearer 前缀
    "Content-Type": "application/json"
}

print("登录成功 ✅")

TASK_COUNT = 1000

cron_patterns = [
    "*/10 * * * * *",    # 每10秒执行一次（支持6位格式）
    "*/30 * * * * *",    # 每30秒执行
    "*/1 * * * *",       # 每1分钟执行
    "*/5 * * * *",       # 每5分钟执行
    "*/15 * * * *",      # 每15分钟执行
]

for i in range(TASK_COUNT):
    # 随机选择一个 cron 表达式
    spec = random.choice(cron_patterns)

    # 构造任务
    task = {
        "name": f"auto_task_{i}",
        "command": f"echo '执行第{i}个任务'",
        "job_type": 1,  
        "spec": spec,
        "allocation": 2,  # 自动分配节点
        "timeout": random.choice([5, 10, 20, 30]),
        "retry_times": random.choice([0, 1, 2]),
        "status": 1,
        "retry_interval": random.choice([0, 5, 10]),
        "notify_status": 2,
        "notify_type": 0,
        "note": f"自动测试任务_{i}，间隔：{spec}",
    }

    # 发请求
    resp = requests.post(JOB_ADD_URL, json=task, headers=headers)
    if resp.status_code != 200:
        print(f"❌ 任务{i} 创建失败: {resp.text}")
    else:
        result = resp.json()
        if result.get("code") != 200:
            print(f"⚠️ 任务{i} 创建失败: {result.get('msg')}")
        else:
            if i % 50 == 0:
                print(f"✅ 已成功创建 {i} 个任务...")

    # 控制请求频率，防止压测过猛
    time.sleep(0.1)

print("🎉 所有任务创建完毕！")
