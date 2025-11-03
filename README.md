# Ticketly Backend

JWT와 Redis를 활용한 안전한 인증 시스템을 갖춘 티켓 관리 백엔드 서비스

## 기능

- ✅ 회원가입 (비밀번호 bcrypt 해싱)
- ✅ 로그인 (JWT Access Token + Refresh Token)
- ✅ 토큰 갱신
- ✅ 로그아웃 (토큰 블랙리스트)
- ✅ 인증 미들웨어
- ✅ Redis를 통한 토큰 관리

## 기술 스택

- **언어**: Go 1.25.2
- **프레임워크**: Fiber v2
- **데이터베이스**: MySQL (Ent ORM)
- **캐시**: Redis
- **인증**: JWT (golang-jwt/jwt/v5)
- **비밀번호 해싱**: bcrypt

## 프로젝트 구조

```
ticketly-backend/
├── cmd/
│   └── gocore/
│       └── main.go                 # 애플리케이션 진입점
├── config/
│   └── config.go                   # 환경 변수 설정
├── internal/
│   ├── db/
│   │   ├── db.go                   # MySQL 연결
│   │   └── redis.go                # Redis 연결
│   ├── domain/
│   │   ├── errors.go               # 도메인 에러 정의
│   │   └── user.go                 # User 도메인 모델
│   ├── handler/
│   │   └── authHandler.go          # 인증 HTTP 핸들러
│   ├── middleware/
│   │   ├── authMiddleware.go       # JWT 인증 미들웨어
│   │   └── middleware.go
│   ├── repository/
│   │   ├── mysql/
│   │   │   └── userRepository.go   # User MySQL 저장소
│   │   └── redis/
│   │       └── tokenRepository.go  # Token Redis 저장소
│   ├── usecase/
│   │   ├── authUsecase.go          # 인증 비즈니스 로직
│   │   └── userUsecase.go          # 사용자 비즈니스 로직
│   └── util/
│       └── jwt.go                  # JWT 유틸리티
└── lib/
    └── ent/                        # Ent ORM 생성 파일
```

## 시작하기

### 사전 요구사항

- Go 1.25.2 이상
- MySQL 8.0 이상
- Redis 6.0 이상

### 설치

1. 저장소 클론
```bash
git clone <repository-url>
cd ticketly-backend
```

2. 의존성 설치
```bash
go mod download
```

3. 환경 변수 설정
```bash
cp .env.example .env.dev
```

`.env.dev` 파일을 열어 다음 값들을 설정하세요:
- 데이터베이스 연결 정보
- Redis 연결 정보
- JWT 시크릿 키 (프로덕션에서는 반드시 변경)

4. MySQL 데이터베이스 생성
```sql
CREATE DATABASE ticketly CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. Redis 실행
```bash
redis-server
```

6. 애플리케이션 실행
```bash
go run cmd/gocore/main.go
```

서버는 기본적으로 `:3000` 포트에서 실행됩니다.

## API 엔드포인트

### 공개 엔드포인트 (인증 불필요)

#### 회원가입
```http
POST /auth/register
Content-Type: application/json

{
  "first_name": "홍",
  "last_name": "길동",
  "nick_name": "hong123",
  "birthday": "1990-01-01",
  "email": "hong@example.com",
  "password": "SecurePassword123!",
  "phone_number": "010-1234-5678"
}
```

**응답**:
```json
{
  "message": "user created successfully",
  "user": {
    "id": "uuid",
    "first_name": "홍",
    "last_name": "길동",
    "nick_name": "hong123",
    "email": "hong@example.com",
    ...
  }
}
```

#### 로그인
```http
POST /auth/login
Content-Type: application/json

{
  "email": "hong@example.com",
  "password": "SecurePassword123!"
}
```

**응답**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "hong@example.com",
    ...
  }
}
```

#### 토큰 갱신
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**응답**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 보호된 엔드포인트 (인증 필요)

모든 요청에 Authorization 헤더 필요:
```
Authorization: Bearer <access_token>
```

#### 내 정보 조회
```http
GET /api/me
Authorization: Bearer <access_token>
```

**응답**:
```json
{
  "user_id": "uuid",
  "email": "hong@example.com"
}
```

#### 로그아웃
```http
POST /api/logout
Authorization: Bearer <access_token>
```

**응답**:
```json
{
  "message": "logged out successfully"
}
```

## 보안 기능

### JWT 토큰 관리

1. **Access Token**: 15분 만료
   - API 요청 인증에 사용
   - 짧은 수명으로 보안 강화

2. **Refresh Token**: 7일 만료
   - Access Token 갱신에 사용
   - Redis에 저장되어 관리
   - 로그아웃 시 삭제됨

### Redis를 통한 토큰 보안

1. **Refresh Token 저장**
   - 사용자별로 Redis에 저장
   - Key: `refresh_token:{user_id}`
   - TTL: 7일

2. **Token Blacklist**
   - 로그아웃 시 Access Token을 블랙리스트에 추가
   - Key: `blacklist:{token}`
   - TTL: 토큰 만료 시간과 동일

3. **토큰 검증 프로세스**
   - JWT 서명 검증
   - 만료 시간 확인
   - 블랙리스트 확인
   - Refresh Token의 경우 Redis에 저장된 값과 비교

### 비밀번호 보안

- bcrypt 알고리즘 사용 (DefaultCost)
- 솔트 자동 생성
- 평문 비밀번호는 저장하지 않음

## 환경 변수

| 변수명 | 설명 | 기본값 |
|--------|------|--------|
| DB_HOST | MySQL 호스트 | localhost |
| DB_PORT | MySQL 포트 | 3306 |
| DB_USER | MySQL 사용자명 | root |
| DB_PASSWORD | MySQL 비밀번호 | - |
| DB_NAME | 데이터베이스 이름 | ticketly |
| REDIS_ADDR | Redis 주소 | localhost:6379 |
| REDIS_PASSWORD | Redis 비밀번호 | - |
| JWT_ACCESS_SECRET | Access Token 시크릿 | (기본값 있음, 변경 필요) |
| JWT_REFRESH_SECRET | Refresh Token 시크릿 | (기본값 있음, 변경 필요) |
| PORT | 서버 포트 | 3000 |

## 개발

### 코드 빌드
```bash
go build -o bin/server cmd/gocore/main.go
```

### 실행
```bash
./bin/server
```

### 테스트
```bash
go test ./...
```

## 라이선스

MIT License

## 기여

이슈와 풀 리퀘스트는 언제나 환영합니다!
