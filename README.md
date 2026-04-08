# TeamSpeak to Telegram
Updates a pinned Telegram message with online TeamSpeak users and optionally prepends the user count to the chat title for groups.
Supports both TeamSpeak6 and TeamSpeak3.

![pinned-message-screenshot](.github/screenshots/pinned-message.png)

## Usage
Images are available on GitHub Container Registry. Pinning to a specific commit is recommended as there might be breaking changes.

```yaml
services:
  teamspeak-to-telegram:
    image: ghcr.io/kevincali/teamspeak-to-telegram:<SHORT_COMMIT_HASH> # or latest
    restart: unless-stopped
    volumes:
      - ./state.json:/home/nonroot/state.json
    environment:
      TS6_ADDRESS: http://your-ts6-host:10080
      TS6_API_KEY: your-api-key-here
      TELEGRAM_BOT_TOKEN: your-bot-token
      TELEGRAM_CHAT_ID: "-1001234567890"
```

Create the state file with necessary permissions before the first start:
```sh
touch state.json && chmod 666 state.json
```

## Config
Configuration is done via environment variables. Use either TS3 or TS6, not both.

| Variable                      | Description                                            | Required           | Default      |
|-------------------------------|--------------------------------------------------------|--------------------|--------------|
| `TS6_ADDRESS`                 | TS6 HTTP query address (e.g. `http://localhost:10080`) | Yes (if using TS6) | -            |
| `TS6_API_KEY`                 | HTTP query API key                                     | Yes (if using TS6) | -            |
| `TS3_ADDRESS`                 | TS3 server query address (e.g. `localhost:10011`)      | Yes (if using TS3) | -            |
| `TS3_USERNAME`                | Server query username                                  | Yes (if using TS3) | -            |
| `TS3_PASSWORD`                | Server query password                                  | Yes (if using TS3) | -            |
| `TS3_VIRTUAL_SERVER_ID`       | Virtual server ID                                      | Yes (if using TS3) | -            |
| `TS_FAVORITE_USERS`           | Comma-separated list of usernames to filter            | No                 | (all users)  |
| `TS_POLLING_INTERVAL_SECONDS` | Polling interval in seconds                            | No                 | `60`         |
| `TELEGRAM_BOT_TOKEN`          | Telegram bot token                                     | Yes                | -            |
| `TELEGRAM_CHAT_ID`            | Telegram chat ID                                       | Yes                | -            |
| `TELEGRAM_SEPARATOR`          | Separator between usernames                            | No                 | ` \| `       |
| `TELEGRAM_ZERO_USERS`         | Text when no users are online                          | No                 | `¯\_(ツ)_/¯` |
| `TELEGRAM_UPDATE_TITLE`       | Prepend online user count to chat title (only groups)  | No                 | `false`      |

The message ID is auto-saved to `state.json` and reused on restart.

### Getting your credentials

**TeamSpeak 6**
- Get your API key from the server logs on first startup, or via SSH query: `use 1` then `apikeyadd scope=read lifetime=0`
- If you use a query allowlist, add your IP

**TeamSpeak 3**
- Connect with Server Admin permissions
- Go to `Tools` -> `ServerQuery Login` to get the username/password

**Telegram**
- Create a bot via [@BotFather](https://t.me/BotFather) and add it to your group
- Grant it admin permissions: *Pin messages* (required), *Change group info* and *Delete messages* (only needed for `TELEGRAM_UPDATE_TITLE`)
- Get the chat ID with [@username_to_id_bot](https://t.me/username_to_id_bot) (not official, use with caution)
