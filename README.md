# Intro

Generate customizable scaffold command on your wish.

# Usage

By default, this project would create a comamnd named `scaffold`, which will create a dummy greeting project scaffold.

One would like to customize and generate real scaffold on his/her wish. So one would download this project first. Then remove *template.go*, and modify both *template* and *main.go*.

**NOTE**: All files in *template* directory are treated as template, any valid template directives will be `Exec()` by `template` package.

Afterwards, run `go generate && go install` from root directory. If everything goes fine, you would get your own `scaffold` command.
