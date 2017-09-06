# go-scripting

Handy toolset when using Go as a shell script replacement

# CAVEAT EMPTOR

The API is still very unstable.

# DESCRIPTION

Using Go to automate administrative tasks (e.g. deployment) is actually quite
useful. Except, doing something as trivial as `execute a command, and then
grep for a particular pattern` doesn't feel as easy as writing a real shell
script.

We can not just make everything magical like shell, but we can make certain
operations a bit easier. That is what these libraries are for.

# EXECUTING A SIMPLE COMMAND

If you don't care to capture stdout or stderr, and all you care if the command
runs successfully, you can use the shorthand `cmd.Exec`

```go
  cmd.Exec("make", "install")
```

# ADVANCED USAGE

You can create a `cmd.Command` instance by using `cmd.New`

```go
  c := cmd.New("ls", "-l")
```

You can all `cmd.Do` to execute this command. Either pass nil, or the
`context.Context` that you would like to use

```go
  res, err := c.Do(nil)
```

The first return value is an instance of `cmd.Result`. It will contain
information on the result of executing the command, such as captured
output.

Output is not captured by default. You must explicitly state to do so.

```go
  c.CaptureStdout(true)
  c.CaptureStderr(true)
  res, _ := c.Do(nil)
  fmt.Println(res.OutputString())
```

You can run certain filters on your output, such as `Grep`

```go
  c.Grep(`regular expression pattern`)
```

If you specify this before callin `Do()`, your result output will be filtered
accordingly.

Finally, this was getting really long. To avoid having to type everything in
separate method calls, you can chain them all together

```go
  res, err := cmd.New("ls", "-l").
    CaptureStdout(true).
    CaptureStderr(true).
    Grep(`regular expression pattern`).
    Do(nil)
```

For other options, please consult the godoc.
