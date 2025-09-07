# 서울시 주택공고 모니터링 봇

서울시 LH 및 SH 주택공고를 자동으로 모니터링하고 텔레그램으로 알림을 전송하는 Go 애플리케이션입니다.

## 기능

- 🏠 **서울시 주택공고 전용**: LH, SH 주택공고 사이트 특화 스크래핑
- 📊 **요약 알림**: 최근 7일간 모든 공고를 요약하여 전송 
- 🔗 **다중 사이트 지원**: LH와 SH 사이트 동시 모니터링
- 📨 **텔레그램 알림**: 공고 제목, 게시일, URL 포함한 상세 정보 전송
- ⏰ **GitHub Actions 자동화**: 매일 정해진 시간에 자동 실행
- 🚀 **상태 비저장**: 매 실행마다 최신 정보만 확인하는 단순한 구조

## GitHub Actions 설정 (권장)

### 1. 프로젝트를 GitHub에 업로드

### 2. GitHub Secrets 설정
Repository Settings > Secrets and variables > Actions에서 다음 secrets를 추가하세요:

- `TELEGRAM_BOT_TOKEN`: 텔레그램 봇 토큰
- `TELEGRAM_CHAT_ID`: 텔레그램 채팅 ID  
- `TARGET_URL`: 모니터링할 주택공고 URL (콤마로 구분하여 여러 개 설정 가능)

#### TARGET_URL 설정 예시
```
https://example.com/board1,https://example.com/board2
```

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
TARGET_URL=https://example.com/board1,https://example.com/board2
```

### 2. 로컬 실행
```bash
go mod download
go run main.go
```

## 알림 방식

### 공고가 있는 경우
- 최근 7일간의 모든 주택공고를 요약하여 전송
- 각 공고마다 제목, 게시일, URL, 출처(LH/SH) 정보 포함
- 최대 10개까지 표시하고 초과 시 개수 안내

### 공고가 없는 경우  
- "최근 7일간 새로운 공고가 없습니다" 알림 전송
- 확인 완료 시간 포함

## 지원 사이트

- **LH 서울지역본부**: 공공분양/임대 공고
- **SH 서울주택도시공사**: 공공분양/임대 공고

각 사이트의 테이블 구조에 맞춰 최적화된 스크래핑 로직 적용

## 프로젝트 구조

```
go-newPost/
├── main.go                    # 메인 애플리케이션 - 다중 URL 처리 및 통합 알림
├── config/                    # 설정 관리
│   └── config.go             # 환경변수 로드 및 다중 URL 파싱
├── scraper/                   # 웹 스크래핑
│   └── scraper.go            # LH/SH 사이트 특화 스크래핑 로직
├── telegram/                  # 텔레그램 봇
│   └── bot.go                # 요약 알림 및 에러 메시지 전송
├── .env.example              # 환경변수 예시 (실제 서울 주택공고 URL 포함)
├── .github/workflows/        # GitHub Actions 워크플로우
│   └── monitor.yml          # 자동화 스케줄 설정
└── README.md
```

## 주의사항

- 서울시 주택공고 사이트 전용으로 제작되었습니다
- 사이트의 이용약관을 준수하며 사용해주세요
- 봇 토큰과 채팅 ID는 안전하게 관리하세요
- 하루 한 번 실행으로 적절한 요청 간격을 유지합니다