# Chat Stream Completion

```golang
_, err = bob.ChatCompletionStream(context.Background(), func(self *agents.Agent, content string, err error) error {
    fmt.Print(content)
    return nil
})

if err != nil {
    panic(err)
}
```