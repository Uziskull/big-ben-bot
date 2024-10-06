# Big Ben

Discord bot that joins any active voice channel, and lets you know what time is it. Good to annoy your friends (and yourself).
Inspired by some random video online.

## Usage

The bot is not publicly available, so you'll have to host it yourself.

1. Create a Discord Bot (scopes: `bot`; permissions: `connect`, `speak` (maybe `view_channels` too, unsure)).
2. Get the Docker image, stick it somewhere hostable, and feed it the necessary environment variables.
3. Add it to a server.

## Env Variables

| Variable  | Description                                         | Default Value                                            |
|-----------|-----------------------------------------------------|----------------------------------------------------------|
| BOT_TOKEN | **Mandatory.** Connects the app to the Discord bot. | None; the service exits with an error if it's not given. |
| TZ        | Lets the bot know which timezone it should follow.  | `UTC`                                                    |

## Notes

* The `go.mod` file has little information, as the Docker image does the `go mod tidy`.
  This is due to me not wanting to install Go on my workstation; it's bad design and I should fix it eventually, one day, maybe (not really).
* Sounds are converted to .DCA within the Dockerfile build, and the names are hardcoded.
  If you want to change them, consider changing the names within `Dockerfile`, `main.go` and the files themselves.
