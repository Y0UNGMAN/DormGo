import requests
import json

# 测试个性签名功能
BASE_URL = "http://localhost:8080"

# 1. 先注册一个测试用户
signup_data = {
    "username": "testuser_sig",
    "password": "123456",
    "re_password": "123456",
    "student_id": "20230001",
    "dorm_id": 1
}

print("1. 注册测试用户...")
response = requests.post(f"{BASE_URL}/api/v1/user/signup", json=signup_data)
print(f"注册响应: {response.status_code} - {response.text}")

# 2. 登录获取 token
login_data = {
    "username": "testuser_sig",
    "password": "123456"
}

print("\n2. 登录获取 token...")
response = requests.post(f"{BASE_URL}/api/v1/user/login", json=login_data)
login_result = response.json()
print(f"登录响应: {response.status_code}")
print(f"Token: {login_result.get('token')}")

if response.status_code == 200 and login_result.get('token'):
    token = login_result['token']
    headers = {"Authorization": f"Bearer {token}"}
    
    # 3. 获取用户资料，查看当前个性签名
    print("\n3. 获取用户资料...")
    response = requests.get(f"{BASE_URL}/api/v1/user/profile", headers=headers)
    profile_result = response.json()
    print(f"用户资料响应: {response.status_code}")
    print(f"当前个性签名: {profile_result.get('data', {}).get('intro', '未设置')}")
    
    # 4. 更新个性签名
    update_data = {
        "nickname": profile_result['data']['nickname'],
        "avatar": profile_result['data']['avatar'],
        "dorm_id": profile_result['data']['dorm_id'],
        "bio": "这是我的个性签名测试！"
    }
    
    print("\n4. 更新个性签名...")
    response = requests.put(f"{BASE_URL}/api/v1/user/profile", headers=headers, json=update_data)
    update_result = response.json()
    print(f"更新响应: {response.status_code} - {update_result}")
    
    # 5. 再次获取用户资料，确认个性签名已更新
    print("\n5. 再次获取用户资料，确认更新...")
    response = requests.get(f"{BASE_URL}/api/v1/user/profile", headers=headers)
    profile_result = response.json()
    print(f"更新后的个性签名: {profile_result.get('data', {}).get('intro', '未设置')}")
    
    # 6. 获取用户统计数据，确认个性签名也在这里返回
    print("\n6. 获取用户统计数据...")
    response = requests.get(f"{BASE_URL}/api/v1/user/stats", headers=headers)
    stats_result = response.json()
    print(f"统计数据中的个性签名: {stats_result.get('data', {}).get('bio', '未设置')}")
else:
    print("登录失败，无法继续测试")