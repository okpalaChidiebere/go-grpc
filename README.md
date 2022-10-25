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

## Environment Configurations

- You need to have go installed in your system
- You need the Protocol buffer [`protoc`] compiler to able to compile your proto file. In my case i used HomeBrew to install this
- You need the protocol buffer compiler plugins; one for proto model compiling and one for rpc compiling.
- You may also need to update your go path
- See full steps [here](https://grpc.io/docs/languages/go/quickstart/#prerequisites)
- If i wanted to set up a CI/CD process using Travis CI, so whenever i add a new proto file in a github repo to generate the protos for me on another github repo or on thesame repo i can read [this](https://docs.travis-ci.com/user/installing-dependencies/#installing-packages-on-macos) and [this](https://docs.travis-ci.com/user/installing-dependencies/#installing-projects-from-source)

## VSCode Configurations to edit proto files

- Install proto3 plugin that the ide suggested
- You may also have to set your IDE (my case was VSCode) to be able to read proto file importing another proto file without errors. In VSCode go to preferences > settings and search for 'proto3' edit the settings.json file add

```json
{
  ...
  "protoc": {
	"path": "/usr/local/bin/protoc", /*This is optional. You can get this path by running in termial `which protoc`*/
    "options": ["--proto_path=protos"] /*the path will be the path that your proto file is defined. For this project its called `protos`*/
  }
}

```

- If you want auto format of the files like tabs spaces you can install the clang format plugin. Via HomeBrew `brew install clang-format` or regular VSCode install
- Restart the VSCode! See [Video](https://www.youtube.com/watch?v=3r327rjB8qg)

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

## More to explore in gRPC

- gRPC Interceptor. See example [here](https://github.com/grpc/grpc-go/tree/master/examples/features/interceptor). A video tutorial [here](https://www.youtube.com/watch?v=kVpB-uH6X-s) where we have a server interceptor authenticate users with JWT and authorize access by roles and the client interceptor to login user and attach jwt token to the request before calling the actual grpc api.
- [https://github.com/grpc/grpc-go](https://github.com/grpc/grpc-go)

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
