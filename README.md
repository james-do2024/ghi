# ghi
An interactive GitHub CLI.

Think of it like the official gh CLI, with no oversight and a small fraction of the functionality.

## Usage

ghi is a Cobra application with two subcommands, `view` and `explore`. 

The former is a non-interactive command that displays contents of a repository or file and immediately
exits. 

`explore` is interactive and provides the user a prompt to explore the repository. It also attempts 
syntax highlighting based on the syntax of the file, and paging for file output.


The repository name and path format is designed to be as forgiving as possible.

Given a valid repository, both of the following formats are accepted:

- `ghi explore org/repo dir`
- `ghi explore org/repo/dir`


When using the `explore` subcommand, the commands are as follows:

```
[0-9] : Select entry by number (can also be 10 or higher)
^  : The caret character navigates to the root of the repository
.. : Navigates one directory back, if possible
q  : Quits the interactive session
```

Use the `-h` argument for online documentation.

## Environment Variables

The following must be set for ghi to work:

`GH_AUTH_TOKEN`: ghi is tested and working with "fine-grained" GitHub auth tokens

Additionally, the user can supply any color scheme supported by the [alecthomas/chroma](https://github.com/alecthomas/chroma) package
by setting `GH_COLOR_THEME`. The author recommends "monokai".

In the event that your base API URL is other than api.github.com (as it is for some GitHub Enterprise subscribers), an alternate may be specified 
in `GH_BASE_URL`.

Lastly, if you wish to identify with a different user agent than the official GitHub Go API package provides, you may do so by setting
`GH_USER_AGENT` to anything but an empty string. The GitHub API requires a user agent.

## Author
James Taylor
