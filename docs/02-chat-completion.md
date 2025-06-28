# Chat Completion

```golang
response, err := bob.ChatCompletion(context.Background())

if err != nil {
    panic(err)
}

println("Response from Bob:", response)
```