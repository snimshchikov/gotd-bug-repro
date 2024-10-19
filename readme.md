# Reproduction of https://github.com/gotd/td/issues/1030

Required environment variables:

- `SESSION_FILE`: path to where session file will be stored
- `BOT_TOKEN`: bot token from BotFather (can be new bot)
- `APP_ID`: app_id of Telegram app
- `APP_HASH`: app_hash of Telegram app

## How to reproduce

1. Run `go run .` in this directory
2. Wait
3. After ~2000 requests, connection will be closed
4. Read logs from debug.log file
