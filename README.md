# Go NewPost Monitor

특정 홈페이지 게시판의 새 글을 모니터링하고 텔레그램으로 알림을 전송하는 Go 애플리케이션입니다.

## 기능

- 🔍 웹사이트 게시판 자동 모니터링
- 📨 새 글 발견 시 텔레그램 알림
- ⏰ GitHub Actions를 통한 주기적 실행
- 💾 중복 알림 방지를 위한 데이터 저장

## GitHub Actions 설정 (권장)

### 1. 프로젝트를 GitHub에 업로드

### 2. GitHub Secrets 설정
Repository Settings > Secrets and variables > Actions에서 다음 secrets를 추가하세요:

- `TELEGRAM_BOT_TOKEN`: 텔레그램 봇 토큰
- `TELEGRAM_CHAT_ID`: 텔레그램 채팅 ID  
- `TARGET_URL`: 모니터링할 게시판 URL

### 3. 텔레그램 봇 설정

1. [@BotFather](https://t.me/botfather)에게 `/newbot` 명령을 보내 새 봇 생성
2. 봇 토큰을 받아 GitHub Secrets에 `TELEGRAM_BOT_TOKEN`으로 설정
3. 봇과 대화를 시작하고 채팅 ID 확인:
   - [@userinfobot](https://t.me/userinfobot)에게 메시지를 보내 채팅 ID 확인
   - 또는 봇에게 메시지를 보낸 후 `https://api.telegram.org/bot<TOKEN>/getUpdates`로 확인
   - 받은 채팅 ID를 GitHub Secrets에 `TELEGRAM_CHAT_ID`로 설정

### 4. 워크플로우 실행
- **자동 실행**: 매일 UTC 자정 (한국시간 오전 9시)에 자동 실행
- **수동 실행**: Actions 탭에서 "Board Monitor" 워크플로우를 수동으로 실행 가능

## 로컬 실행 (개발/테스트용)

### 1. 환경변수 설정
```bash
cp .env.example .env
```

`.env` 파일을 편집하여 필요한 값들을 설정하세요:

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
TELEGRAM_CHAT_ID=your_telegram_chat_id_here
TARGET_URL=https://example.com/board
DATA_FILE=data.json
```

### 2. 로컬 실행
```bash
go mod download
go run main.go
```

## 설정 옵션

- `TARGET_URL`: 모니터링할 게시판 URL
- `DATA_FILE`: 데이터 저장 파일 경로 (기본값: data.json)

## 프로젝트 구조

```
go-newPost/
├── main.go              # 메인 애플리케이션
├── config/              # 설정 관리
│   └── config.go
├── scraper/             # 웹 스크래핑
│   └── scraper.go
├── telegram/            # 텔레그램 봇
│   └── bot.go
├── storage/             # 데이터 저장
│   └── storage.go
├── .env.example         # 환경변수 예시
└── README.md
```

## 주의사항

- 모니터링 대상 사이트의 이용약관을 확인하세요
- 너무 자주 요청하지 않도록 적절한 간격을 설정하세요
- 봇 토큰과 채팅 ID는 안전하게 관리하세요