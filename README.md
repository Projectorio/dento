# dento
## Getting Started
In order to use `dento` you need to have a `ChatGPT`-[account](https://chat.openai.com/auth/login).

Get Api-Key and OrgId for ChatGpt and Set them as environment-variables

```bash
export OPEN_AI_API_KEY=<Your-Api-Key>
export OPEN_AI_API_ORG=<Your-Org-Id>
```

## Commands
### Explain
`dento explain` will try to explain the given keyword in the context of kubernetes using `gpt-3.5-turbo` model.


## Mockery
To create a test mocking the communication with ChatGPT using Mockery, you can follow these steps:

1. Install the Mockery library by running `go get github.com/vektra/mockery/v2/...` in your terminal.

2. Create a new file called `mock_client.go` in the same directory as your `explainCmd` file.

3. In `mock_client.go`, create an interface called `MockClient` with the same method signatures as the `goopenai.Client` interface:

```
type MockClient interface {
    CreateCompletions(ctx context.Context, req goopenai.CreateCompletionsRequest) (*goopenai.CreateCompletionsResponse, error)
}
```

4. Use Mockery to generate a mock implementation of `MockClient`:

```
mockery --name=MockClient
```

This will create a `mock_client.go` file in a `mocks` directory. This file will contain a mock implementation of the `MockClient` interface.

5. In your test file, import the `mocks` package:

```
import "your_project_path/mocks"
```

6. Create a new test case for your `explainCmd` function.

7. Create an instance of your mock client:

```
mockClient := &mocks.MockClient{}
```

8. Set up the expectations for the `CreateCompletions` method:

```
expectedRequest := goopenai.CreateCompletionsRequest{
    Model:       "gpt-3.5-turbo",
    Messages:    []goopenai.Message{},
    Temperature: 0.2,
}
mockClient.On("CreateCompletions", context.Background(), expectedRequest).Return(
    &goopenai.CreateCompletionsResponse{
        Choices: []goopenai.CompletionChoice{
            {
                Message: goopenai.Message{
                    Content: "your_json_response",
                },
            },
        },
    },
    nil,
)
```

9. Inject the mock client into your `explainCmd` function:

```
client = mockClient
```

10. Call your `explainCmd` function with your mock client and assert the output.