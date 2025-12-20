import requests
import json

# 最终测试个性签名功能
BASE_URL = "http://localhost:8080"

# 登录现有用户
login_data = {
    "username": "testuser_sig",
    "password": "123456"
}

print("1. 登录测试用户...")
response = requests.post(f"{BASE_URL}/api/v1/user/login", json=login_data)
login_result = response.json()
print(f"登录响应: {response.status_code}")

if response.status_code == 200 and login_result.get('token'):
    token = login_result['token']
    headers = {"Authorization": f"Bearer {token}"}
    
    # 2. 获取用户资料，查看当前个性签名
    print("\n2. 获取用户资料...")
    response = requests.get(f"{BASE_URL}/api/v1/user/profile", headers=headers)
    profile_result = response.json()
    print(f"用户资料响应: {response.status_code}")
    print(f"当前个性签名: '{profile_result.get('data', {}).get('intro', '未设置')}'")
    
    # 3. 更新个性签名
    new_signature = "这是我最新的个性签名！"
    update_data = {
        "nickname": profile_result['data']['nickname'],
        "avatar": profile_result['data']['avatar'],
        "dorm_id": profile_result['data']['dorm_id'],
        "bio": new_signature
    }
    
    print(f"\n3. 更新个性签名为: '{new_signature}'")
    response = requests.put(f"{BASE_URL}/api/v1/user/profile", headers=headers, json=update_data)
    update_result = response.json()
    print(f"更新响应: {response.status_code} - {update_result.get('msg')}")
    
    # 4. 再次获取用户资料，确认个性签名已更新
    print("\n4. 再次获取用户资料，确认更新...")
    response = requests.get(f"{BASE_URL}/api/v1/user/profile", headers=headers)
    profile_result = response.json()
    current_signature = profile_result.get('data', {}).get('intro', '未设置')
    print(f"更新后的个性签名: '{current_signature}'")
    
    # 5. 获取用户统计数据，确认个性签名也在这里返回
    print("\n5. 获取用户统计数据...")
    response = requests.get(f"{BASE_URL}/api/v1/user/stats", headers=headers)
    stats_result = response.json()
    stats_signature = stats_result.get('data', {}).get('bio', '未设置')
    print(f"统计数据中的个性签名: '{stats_signature}'")
    
    # 6. 验证两个接口返回的签名一致
    if current_signature == new_signature and stats_signature == new_signature:
        print("\n✅ 测试成功！个性签名功能正常工作。")
    else:
        print("\n❌ 测试失败！接口返回的签名不一致。")
else:
    print("登录失败，无法继续测试")