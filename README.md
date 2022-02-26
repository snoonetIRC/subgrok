# subgrok

`subgrok` is an IRC bot which alerts IRC channels when new posts are made to
subreddits.

## Bot usage

The bot should be invited to your channel (`/INVITE subgrok #channelname`).

The bot has three commands:

### `!subscribe <subreddit name here>`

Subscribes the current channel to receive messages when a new post is made on
the provided subreddit.

**Please note that subreddit names are case-sensitive.**

Permitted to channel half-operators and above only.

```
<@Mike>    !subscribe /r/Metal
< subgrok> Subscribed #metal to Metal
```

### `!subscriptions`

Lists the subreddits the current channel is subscribed to.

Permitted to channel half-operators and above only.

```
<@Mike> !subscriptions
< subgrok> #metal subscribes to: Metal
```

### `!unsubscribe <subreddit name here>`

Unsubscribes the current channel from messages about new posts on the provided
subreddit.

**Please note that subreddit names are case-sensitive.**

Permitted to channel half-operators and above only.

```
<@Mike>    !unsubscribe /r/Metal
< subgrok> Unsubscribed #metal from Metal
```

## Hosting your own

Check the [releases](https://github.com/snoonetIRC/subgrok/releases) page for
the application's most recent releases.

Linux x86_64 binaries are automatically attached to every build. This guide
assumes you would like to host the bot on a 64 bit Linux server.

### Configuration reference

```yaml
---
# "irc" is configuration used for the IRC connection
irc:
  # admin_channels are joined by default, even if they do not have any subscriptions.
  admin_channels:
    - '##my-subgrok-admin-channel'
  debug: false            # Display verbose IRC debug-level information
  ident: subgrok          # The "username" the bot will connect with
  modes: +B               # umodes that'll be set against the bot (at least +B recommended)
  nickname: subgrok       # The nickname the bot will use
  port: 6697              # IRC port
  real_name: subgrok      # Realname
  server: irc.snoonet.org # IRC server
  use_tls: true           # Boolean, whether the bot should connect using SSL

  nickserv_account: my-nickserv-account   # The nickserv username the bot will identify with
  nickserv_password: my-nickserv-password # The nickserv password the bot will identify with

# "reddit" is configuration for our use of the reddit API
reddit:
  poll_wait_time: 600 # Time to wait (seconds) between reddit API calls

# "database" is configuration for the database
database:
  filepath: '~/.config/snoonet/subgrok/file.db' # Path to use as a file database

# "application" is configuration for the bot itself
application:
  channel_maximum_subscriptions: 20 # The maximum number of subreddits any channel may watch
```

### Running the bot

By default, on Linux systems, the bot will look for a configuration file in
`~/.config/snoonet/subgrok/config.yaml`. It won't launch if one isn't found.

After downloading the most recent binary, you can run it by issuing the following
commands:

```
% chmod +x subgrok
% ./subgrok
```

## Development

### Running the test suite

```
% go test -v ./...
```

### Building

If you would like to build the application yourself, you will need Golang 1.17+.

```
% go build -o subgrok ./cmd/subgrok/main.go
```
