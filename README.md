<div align="center">
  <h1>Kinctx</h1>
</div>

<div align="center">
  <img src=".github/kinctx.png" alt="kinctx banner" width="40%" height="40%" />
</div>

`kinctx` is a terminal-first command library for the snippets you keep rewriting, re-Googling, or forgetting.

It lets you store commands in a local SQLite database, search them with `fzf`, fill in placeholder values when needed, and execute the final command directly from your terminal.

## Usage First

### Command overview

```bash
kin init
kin add
kin list
kin search
kin rm
kin help
```

### `kin init`

Initializes the local SQLite database used to store your commands.

Run this once before adding, listing, searching, or deleting snippets:

```bash
kin init
```

### `kin add`

Adds a command to your kin store.

You can add commands in two ways:

#### 1. Pipe a command in

```bash
echo 'docker logs -f ${CONTAINER:=api}' | kin add
```

```bash
echo 'sudo systemctl restart NetworkManager' | kin add
```

#### 2. Open your editor

```bash
kin add
```

When no stdin is piped in, `kinctx` opens your `$EDITOR`. If `$EDITOR` is not set, it falls back to `nano`.

This is the best option for:

- multi-line commands
- long commands with many flags
- commands with placeholders

Example multi-line snippet:

```bash
kubectl rollout restart deployment/${DEPLOYMENT:=api-server} \
  -n ${NAMESPACE:=default} && \
kubectl rollout status deployment/${DEPLOYMENT:=api-server} -n ${NAMESPACE:=default}
```

### `kin list`

Lists all saved commands in a styled table, along with usage count and creation date.

```bash
kin list
```

### `kin search`

Opens an interactive `fzf` picker so you can find a saved command, fill placeholders, and run it.

```bash
kin search
```

Workflow:

1. `fzf` shows your saved commands.
2. You select one.
3. `kinctx` detects placeholders like `${BRANCH:=main}`.
4. It prompts you for values.
5. It substitutes the values and executes the final command.


`kinctx` is smart enough to detect if your command doesn't have any placeholder synatx and runs it directly. 

### `kin rm`

Deletes one or more saved commands using a multi-select `fzf` screen.

```bash
kin rm
```

Use `Tab` to mark multiple entries, then press `Enter` to delete them.

## How To Use Stored Commands

### Static commands

A static command has no placeholders:

```bash
echo 'ss -tuln' | kin add
kin search
```

Selecting it executes `ss -tuln` immediately.

### Dynamic commands

A dynamic command includes placeholders:

```bash
echo 'git checkout ${BRANCH:=main}' | kin add
kin search
```

When selected, `kinctx` prompts for `BRANCH`. If you accept the default, it runs:

```bash
git checkout main
```

## Placeholder Syntax

`kinctx` uses a placeholder format inspired by shell parameter expansion:

```bash
${NAME}
${NAME:=default}
```

Supported placeholder keys must look like:

```text
[a-zA-Z_][a-zA-Z0-9_]*
```

### What it means

#### `${NAME}`

Prompt the user for a value with no default.

Example:

```bash
echo 'ssh ${HOST}' | kin add
```

When you run it through `kin search`, you will be asked for `HOST`.

#### `${NAME:=default}`

Prompt the user for a value, but prefill a default.

Example:

```bash
echo 'git checkout ${BRANCH:=main}' | kin add
```

If you press enter without changing it, `main` is used.

In a normal shell, parameter expansion is handled by the shell itself. In `kinctx`, the syntax is used as a templating format before execution.

That means:

- `kinctx` detects placeholders in the saved text
- it asks you for values interactively
- it replaces the placeholders with the values you entered
- it then executes the finished command

Important difference:

`kinctx` is not implementing full shell expansion. It supports the placeholder pattern above, not the full shell parameter expansion language.

### What is not supported

These are not currently part of `kinctx` placeholder parsing:

```bash
${VAR:-default}
${VAR:+alt}
${VAR?message}
${#VAR}
${VAR%suffix}
${VAR/pattern/repl}
```

If you want shell logic like that, wrap the command in a shell explicitly:

```bash
echo 'bash -lc '\''echo ${VAR:-fallback}'\''' | kin add
```

## Install And Run

### Requirements

- Go
- `fzf`

### Build locally

```bash
go build -o kin .
```

Then run:

```bash
./kin init
```

## Troubleshooting

### `no such table: snippets` or `table doesn't exist`

Run:

```bash
kin init
```

This is the most important setup step. `kinctx init` creates the SQLite table used by the app. If you skip it, commands like `add`, `list`, `search`, or `rm` can fail because the database file may exist but the `snippets` table has not been created yet.

### `fzf` is missing

`search` and `rm` depend on `fzf`.

Check:

```bash
fzf --version
```

If it is not installed, install `fzf` and try again.


### I cannot add an empty command

That is expected. `kinctx` rejects blank or whitespace-only input.

### My placeholder syntax is rejected

Check for:

- an unclosed `${`
- an empty key like `${}`
- malformed text like `${:=}`

Valid examples:

```bash
${HOST}
${BRANCH:=main}
${NAMESPACE:=default}
```

### Duplicate command error

The store keeps commands unique. If you try to add the exact same command twice, `kinctx` rejects it.

## Storage

`kinctx` stores its SQLite database under your user config directory:

```text
$XDG_CONFIG_HOME/kinctx/sqlite-database.db
```

On systems without `XDG_CONFIG_HOME`, this resolves through Go's user config directory behavior for your platform.

## License

MIT. See [LICENSE](/home/karthik/Projects/kinctx/LICENSE).
