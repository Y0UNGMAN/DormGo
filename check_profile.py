import requests
import json

# 检查获取用户资料接口返回的数据结构
BASE_URL = "http://localhost:8080"

# 先登录获取 token
login_data = {
    "username": "testuser_sig",
    "password": "123456"
}

print("登录测试用户...")
response = requests.post(f"{BASE_URL}/api/v1/user/login", json=login_data)
login_result = response.json()
print(f"登录状态码: {response.status_code}")

if response.status_code == 200 and login_result.get('token'):
    token = login_result['token']
    headers = {"Authorization": f"Bearer {token}"}
    
    # 获取用户资料
    print("\n获取用户资料...")
    response = requests.get(f"{BASE_URL}/api/v1/user/profile", headers=headers)
    profile_result = response.json()
    print(f"用户资料状态码: {response.status_code}")
    print(f"用户资料响应: {json.dumps(profile_result, indent=2, ensure_ascii=False)}")
else:
    print("登录失败")