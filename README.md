# TeamSpeak to Telegram
Updates a pinned Telegram message with online TeamSpeak users and optionally prepends the user count to the chat title.
Supports both TeamSpeak3 and TeamSpeak6.

![pinned-message-screenshot](.github/screenshots/pinned-message.png)

## Config
Copy `config.example.yaml` to `config.yaml`.

### TeamSpeak 3 (over Telnet)
- Connect to your server with `Server Admin` permissions
- Go to `Tools` → `ServerQuery Login`
- Copy username, password, and server ID (usually 1) to the config

### TeamSpeak 6 (over HTTP)
- Enable HTTP query in your server config
- If you use a query allowlist, add your IP
- Get your API key:
  - Option 1: Check server logs on first startup
  - Option 2: Enable SSH query, ssh into your server, then run `use 1` and `apikeyadd scope=read lifetime=0`
- Copy the API key to the config

### Telegram
- Create a bot via [@BotFather](https://t.me/BotFather)
- Add the bot to your group
- Give it following admin permissions:
  - Pin messages (required)
  - Change group info (needed for the optional user count title feature)
  - Delete messages (needed for the optional user count title feature)
- Get the chat ID
  you could use [@username_to_id_bot](https://t.me/username_to_id_bot) (not official, use with caution)
- Copy the bot token and chat ID to config

## Usage
### Run the container
Images are available here on GitHub.
- setup your `config.yaml`
- `docker run --volume ./config.yaml:/config.yaml --env CONFIG_PATH=/config.yaml ghcr.io/kevincali/teamspeak-to-telegram:latest`

### Build and run the binary
- clone the repository
- `make build` and run the binary or `make run`

