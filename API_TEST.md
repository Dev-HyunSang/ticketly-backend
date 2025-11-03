# API 테스트 명령어

## 1. 회원가입
```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "홍",
    "last_name": "길동",
    "nick_name": "hong123",
    "birthday": "1990-01-01",
    "email": "hong@example.com",
    "password": "SecurePassword123!",
    "phone_number": "010-1234-5678"
  }'
```

## 2. 로그인
```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "hong@example.com",
    "password": "SecurePassword123!"
  }'
```

응답 예시:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "hong@example.com"
  }
}
```

## 3. 내 정보 조회 (인증 필요)
```bash
# ACCESS_TOKEN을 로그인 응답에서 받은 토큰으로 교체하세요
export ACCESS_TOKEN="your_access_token_here"

curl -X GET http://localhost:3000/api/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 4. 토큰 갱신
```bash
# REFRESH_TOKEN을 로그인 응답에서 받은 토큰으로 교체하세요
export REFRESH_TOKEN="your_refresh_token_here"

curl -X POST http://localhost:3000/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }"
```

## 5. 로그아웃 (인증 필요)
```bash
curl -X POST http://localhost:3000/api/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 6. 헬스 체크
```bash
curl http://localhost:3000/health
```

## Postman / Insomnia 사용

위 curl 명령어를 Postman이나 Insomnia에서 import하여 사용할 수 있습니다.

### Authorization Header 설정
Protected endpoints를 테스트할 때:
- Header Name: `Authorization`
- Header Value: `Bearer <your_access_token>`

## 테스트 시나리오

1. 회원가입으로 새 사용자 생성
2. 로그인하여 access_token과 refresh_token 획득
3. access_token을 사용하여 /api/me 호출
4. 15분 후 access_token 만료 시 refresh_token으로 새 access_token 발급
5. 로그아웃하여 토큰 무효화
6. 로그아웃 후에는 같은 access_token으로 API 호출 불가 (블랙리스트에 등록됨)
