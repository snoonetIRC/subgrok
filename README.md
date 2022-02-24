# subgrok

`subgrok` is an IRC bot which alerts IRC channels when new posts are made to
subreddits.

## Bot usage

The bot has three commands:

### `!subscribe <subreddit name here>`

Subscribes the current channel to receive messages when a new post is made on
the provided subreddit.

Permitted to channel half-operators and above only.

```
<@Mike>    !subscribe /r/metal
< subgrok> Subscribed #metal to metal
```

### `!subscriptions`

Lists the subreddits the current channel is subscribed to.

Permitted to channel half-operators and above only.

```
<@Mike> !subscriptions
< subgrok> #metal subscribes to: metal
```

### `!unsubscribe <subreddit name here>`

Unsubscribes the current channel from messages about new posts on the provided
subreddit.

Permitted to channel half-operators and above only.

```
<@Mike>    !unsubscribe /r/metal
< subgrok> Unsubscribed #metal from metal
```

## Hosting your own

Check the [releases](https://github.com/snoonetIRC/subgrok/releases) page for
the application's most recent releases.

Linux x86_64 binaries are automatically attached to every build.

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