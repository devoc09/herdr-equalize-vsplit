# Equalize VSplit

A [Herdr](https://herdr.dev) plugin that splits the current pane to the right
and equalizes column widths in the active tab.

## Requirements

- [Herdr](https://herdr.dev) v0.7.0 or later (macOS)
- [Go](https://go.dev) 1.26.2 or later (the plugin is built from source at install
  time; no prebuilt binaries are shipped)

## Installation

```console
herdr plugin install devoc09/herdr-equalize-vsplit
```

Herdr clones the repo, runs the `[[build]]` step (which compiles the Go source),
and registers the plugin with the ID `devoc09.equalize-vsplit`.

## Usage

### List the available action

```console
herdr plugin action list --plugin devoc09.equalize-vsplit
```

### Invoke the action

```console
herdr plugin action invoke split --plugin devoc09.equalize-vsplit
```

The first invocation splits the current pane to the right and sets both columns
to equal width. Invoke it again on one of the new right-side panes to create
three equal columns. Each invocation finds every sibling split in the `right`
direction and recalculates the equal ratio for each column.

### Comparison with normal vsplit

Starting from two equal-width columns and splitting the right pane, the
difference between a normal vsplit and Equalize VSplit is:

```
Before            50%         |        50%
                  A          |         B
───────────────────────────────────────────────
Normal vsplit     50%         |  25%  |  25%
                  A          |   B   | New
───────────────────────────────────────────────
Equalize VSplit   33.3%  |  33.3%  |  33.3%
                  A      |    B    |  New
```

A normal vsplit divides only the selected pane. Equalize VSplit divides the
selected pane and then redistributes all horizontal columns to equal widths.

### Bind a key (optional)

Add a `[[keys.command]]` entry to your Herdr config
(`~/.config/herdr/config.toml`):

```toml
[[keys.command]]
key = "prefix+v"
type = "plugin_action"
command = "devoc09.equalize-vsplit.split"
description = "Split and equalize columns"
```

Reload the running config without restarting Herdr:

```console
herdr server reload-config
```

### Check the action log

```console
herdr plugin log list --plugin devoc09.equalize-vsplit
```

## Uninstall

```console
herdr plugin uninstall devoc09.equalize-vsplit
```

Or by its install source:

```console
herdr plugin uninstall devoc09/herdr-equalize-vsplit
```

## Development

### Local link

```console
herdr plugin link /path/to/herdr-equalize-vsplit
```

### Run tests

```console
go test ./...
```

### Build

```console
go build -o bin/herdr-equalize-vsplit .
```

## How it works

The plugin communicates with the running Herdr server through its
[socket API](https://herdr.dev/docs/socket-api/). It calls `layout.export` to
read the current tab's layout tree, walks every `right`-directional split, and
computes equal ratios based on the number of columns spanned by each split. It
then calls `layout.set_split_ratio` for each split, ordered from root to leaf.

## Limitations

- Only `right`-direction splits are equalized; `down` splits are left untouched.
- The plugin targets the pane identified by the `HERDR_PANE_ID` environment
  variable.
- Supported on macOS only (as declared in the manifest). Extending to Linux
  requires no code changes beyond updating `platforms` in `herdr-plugin.toml`.
