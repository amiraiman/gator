# Gator

CLI Tools to Track Your RSS Feed

## Requirements

- Postgres
- Go

### Installation

1. Install the `gator` CLI:

   ```sh
   go install github.com/amiraiman/gator@latest
   ```

2. The `gator` binary is placed in your Go bin directory (`$GOBIN`, or `$GOPATH/bin` — default `~/go/bin`). Make sure it's on your `PATH`:

   ```sh
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

   Add that line to your shell profile (`~/.bashrc`, `~/.zshrc`) to make it permanent.

3. Verify:

   ```sh
   gator users
   ```

### Configuration

- Create a config file at `$HOME/.gatorconfig.json`, and place the following content inside of it
```json
{"db_url":"postgres://{POSTGRES_USERNAME}:{POSTGRES_PASSWORD}@{POSTGRES_HOST}:{POSTGRES_PORT}/{POSTGRES_DB}?sslmode=disable","current_user_name":""}
```
- Replace the placeholder `{*}` to its appropriate values as described below:

| VAR NAME | Description |
|----------|-------------|
| POSTGRES_USERNAME | The name of the postgres user that will be used to connect to the database |
| POSTGRES_PASSWORD | The password for the postgres user |
| POSTGRES_HOST | The hostname for the postgres database server |
| POSTGRES_PORT | The port that postgres is listening on |
| POSTGRES_DB | The name of the database that gator will used to store its data |

## Commands

Run any command with `gator <command> [args]`. Commands marked 🔒 requires a
logged-in user (see [`login`](#authentication--users)).

### Authentication & Users

| Command | Usage | Description |
|---------|-------|-------------|
| `login` | `gator login <username>` | Log in as an existing user. The username must already be registered. |
| `register` | `gator register <username>` | Create a new user and log in as them. Fails if the username is taken. |
| `users` | `gator users` | List all users, marking the one currently logged in. |

### Feeds

| Command | Usage | Description |
|---------|-------|-------------|
| `addfeed` 🔒 | `gator addfeed <name> <url>` | Register a new feed and auto-follow it. Fails if the feed already exists. |
| `feeds` | `gator feeds` | List all registered feeds with their URL and creator. |
| `follow` 🔒 | `gator follow <url>` | Follow an existing feed. Fails if the feed isn't registered or is already followed. |
| `unfollow` 🔒 | `gator unfollow <url>` | Stop following a feed. Fails if the feed doesn't exist or isn't followed. |
| `following` 🔒 | `gator following` | List the feeds the current user follows. |

### Posts

| Command | Usage | Description |
|---------|-------|-------------|
| `agg` | `gator agg <interval>` | Continuously poll followed feeds every `<interval>`. Accepts a Go duration like `1s`, `5m`, `1h`. |
| `browse` 🔒 | `gator browse [limit]` | Show recent posts from followed feeds. `limit` is optional, defaults to `2`. |

### Database

| Command | Usage | Description |
|---------|-------|-------------|
| `reset` | `gator reset` | ⚠️ **Destructive, irreversible.** Deletes all users, feeds, and posts. |
