# Gotex

Gotex is a small terminal text editor written in Go. It is built around a deliberately simple editing loop: open a file, edit it directly in the terminal, and write the result back when you exit.

## Requirements

- Go 1.26.5 or newer
- A terminal that supports ANSI escape sequences

## Usage

Gotex edits an existing file supplied as its first positional argument:

```sh
go run . path/to/file.txt
```

You can also build and run the binary:

```sh
go build -o gotex .
./gotex path/to/file.txt
```

Running Gotex without a file prints a short usage message. The file must already exist; Gotex currently opens files for reading before entering the editor.

## Controls

| Key | Action |
| --- | --- |
| Arrow keys | Move the cursor |
| `Enter` | Split the current line |
| `Backspace` | Delete the character before the cursor or remove an empty line |
| `Tab` | Insert a tab |
| `Ctrl+Z` | Page up |
| `Ctrl+X` | Page down |
| `Ctrl+C` | Exit and save |

Printable characters are inserted at the cursor. Gotex restores the terminal mode when it exits, including when the editor returns normally.

## Saving

There is no separate save command or confirmation prompt. Press `Ctrl+C` to leave the editor; the in-memory contents are then written back to the file that was opened. Because the file is rewritten on exit, keep a backup when experimenting with important files.

## Current limitations

- Input is read one byte at a time, so non-ASCII input is not fully supported.
- Lines are displayed without soft wrapping; long lines can extend beyond the terminal width.
- Tabs are displayed using a fixed width of eight columns.
- The editor does not currently offer a new-file workflow, save-as support, undo/redo, search, or a save prompt.
- File line endings are normalized through the scanner/writer path, and the current writer does not add newline characters when saving.

## Development

The Makefile provides the main checks:

```sh
make fmt       # Format Go files
make vet       # Format and run go vet
make staticcheck
make build     # Format, vet, staticcheck, and build
```

The project uses [`golang.org/x/term`](https://pkg.go.dev/golang.org/x/term) for raw terminal mode and terminal dimensions.
