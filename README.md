# GRPC

The purpose of this project is to learn gRPC protocols like Unary gRPC, Server Streaming, Client Streaming, and Bidirectional (combination of client and server streaming). I also learned how to follow the Hexagonal Architecture to arrange our gRPC server files.

## Unary RPC

This is a mode of communication where the client makes a request to the server, then that server might make another rpc request to another service and then response with a result. We do see this in action with the `CreateTodo` and `ReadTodos` RPC for the `TodoService`

## Server Streaming RPC

This is a mode of communication where the client makes a request to the server but the client is expecting a lot of data (like a stream of data) to come back from the server. We do see this in action with the `ReadTodosStream` RPC for the `TodoService`

## Client Streaming RPC

This is a situation where the client is constantly sending information. Eg Asynchronously uploading a huge files. We dp see this in action with the `AddPhoto` RPC for the `TodoService`. The `AddPhoto`method request expects a stream of file bytes or todoID and will return a response of the fileSize uploaded and a bool where the upload was successful or not.

## More info about information about proto that

```proto
//the go_package should be the name of the package where the generated will be.
//eg if you want this to be on a different repo project like "github.com/okpalaChidiebere/generated-protos/todos/v1
option go_package = "github.com/okpalaChidiebere/go-grpc;go_grpc";
// used for java application
option java_outer_classname = "TodoImageProto"; //good name convention is ${fileName}Proto
option java_package = "com.examplecompanyname.path.v1.generated";
```

## Takeaways on Streaming

- Sometimes you make want to receive streams on the background thread in go and using stream context you can listen on the main thread for whenever the background task is done. See this [video](https://www.youtube.com/watch?v=l_74x_qQZB8) as an example. You can read this [article](https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go) about go contexts as well. Read more about concurrency in go [here](https://go.dev/blog/pipelines)
- You may have to write code like this

```go

	// Stream generates values with DoSomething and sends them to out
 // until DoSomething returns an error or ctx.Done from gRPC stream is closed.
 func Stream(ctx context.Context, out chan<- Value) error {
 	for {
 		v, err := DoSomething(ctx)
 		if err != nil {
 			return err
 		}
 		select {
 		case <-ctx.Done():
 			return ctx.Err()
 		case out <- v:
 		}
 	}
 }
```

## Go Read File articles

- [Read FIle](https://zetcode.com/golang/readfile/)
- [Go file Upload](https://tutorialedge.net/golang/go-file-upload-tutorial/)
- [Read file by chunks](https://gist.github.com/rodkranz/90c82583987a15e3d0f2c4678f2835c7)
- [convert file content into an array of bytes](https://socketloop.com/tutorials/convert-file-content-into-array-of-bytes-in-go)

## Useful links

- Another very good gRPC golang tutorial [here](https://earthly.dev/blog/golang-grpc-example/)
- UUID links in go [here](https://yourbasic.org/golang/generate-uuid-guid/) and [here](https://stackoverflow.com/questions/67729822/how-to-generate-a-deterministic-set-of-uuids-in-golang)
- About Go Flags [here](https://www.developer.com/languages/flag-package-go-golang/)
- [Go logs](https://stackoverflow.com/questions/70521948/why-dont-i-see-fmt-logs-in-my-terminal-when-running-go-app-locally)
