# ghi
An interactive GitHub CLI.

Think of it like the official gh CLI, with no oversight and a small fraction of the functionality.

## Usage

ghi is a Conda application with two subcommands, `view` and `explore`. 

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

## Author
James Taylor
