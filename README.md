# TeamSpeak to Telegram
`teamspeak-to-telegram` updates a pinned Telegram message with online TeamSpeak users.

## Configure
- copy `config.example.yaml` to `config.yaml`
- open `config.yaml` in your favorite editor

### TeamSpeak
- connect to your TeamSpeak server with Server Admin permissions
- click on `Tools` and then on `ServerQuery Login`
- enter a username
- copy the generated password to the config
- add your virtual server id (usually `1`) to the config

### Telegram
- create a Telegram bot by contacting [@BotFather](https://t.me/BotFather)
- copy the bot token to the config
- contact the bot in a DM or add it to a group
- get the chat ID
    - you could use [@username_to_id_bot](https://t.me/username_to_id_bot) (not an official Telegram bot, use with caution!)
- copy the chat ID to the config

## Usage
### Run the container
Images are available on [Docker Hub](https://hub.docker.com/r/kevincali/teamspeak-to-telegram).
- `docker pull kevincali/teamspeak-to-telegram:latest`
- `docker run --env CONFIG_PATH=config.yaml kevincali/teamspeak-to-telegram:latest`

### Build and run the binary
- clone the repository
- `make build`
- `make run`

