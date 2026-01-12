import requests
import json

# 检查登录接口返回的数据结构
BASE_URL = "http://localhost:8080"

# 登录数据
login_data = {
    "username": "testuser_sig",
    "password": "123456"
}

print("登录测试用户...")
response = requests.post(f"{BASE_URL}/api/v1/user/login", json=login_data)
print(f"状态码: {response.status_code}")
print(f"响应内容: {json.dumps(response.json(), indent=2, ensure_ascii=False)}")