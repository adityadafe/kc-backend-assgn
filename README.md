# kc-backend-assgn

#### Description:

    This is backend service for Retail Pulse to process images.

    This service has two enpoints:

    1) /api/submit :
        This endpoint accepts the job with payload and return job id instantly
    2) /api/status?jobid=123 :
        This endpoint returns the state of the job like "completed" | "ongoing" | "failed"

#### Assumptions you took if any:

    1) Server will keep running:
        I haven't used any persistent storage and instead opted for in-memory storage since the scope is small. This saves a lot of time because storage operations are straightforward.

    2) Perimeter is not needed:
        Since the status endpoint returns the state of the job but not the perimeter, after calculating the perimeter I'm putting it in a channel. As mentioned in the scope, there's no need to return it.    
    
    3) Enough Memory in server:
        I assume the system has enough memory since every image will be processed in a separate goroutine. I've currently capped max concurrent operations to 100. This number can also be configured in the process module.    

#### Installing (setup) and testing instructions:

    1) make build:
        This command builds and runs the server.

    2) make installswag:
        Installs Swag locally so it can be used to create Swagger documentation.

    3) make swagger:
        This generates a directory called docs/. You can copy the YAML and JSON files and paste them into an online Swagger editor to view the docs.
        Why not spin up the UI locally?
        → That's a good question, but that seemed overkill. I wanted to keep important stuff lightweight.

    4)make test:
        There's very small test coverage, but if you want to check, this command will perform unit tests to verify the perimeter calculation functionality.

    5) make dockerbuild {your-docker-hub-username}:
        Make sure to have BuildKit installed. This will build a multi-platform Docker image and push it to Docker Hub.
        Alternatively, you can run this to avoid build steps:
        docker run --net=host -it kakashifr/kc-backend.

#### Brief description of the work environment used to run this project (Computer/operating system, text editor/IDE, libraries, etc).

    1) Operating System: Arch Linux btw 
    2) Text Editor: Nvim 
    3) libaries: google/uuid, swaggo/swag and most of the heavy lifting is done by std lib.

#### If given more time, what improvements will you do?

    1) Use a message queue:
        For a service that receives >1000 images per job (with multiple concurrent jobs), a message queue is essential. This would also allow extending the architecture with distributed consumers for performance gains.

    2) Use a backend framework:
        While the standard library works for small services, a framework like Gin or Chi would provide battle-tested tools (e.g., middlewares) as the project scales.

    3) Improve test coverage:
        Write more unit tests to cover additional code and functionality. I consider tests very productive long-term and prefer writing them thoroughly.

    4) Set up observability:
        "What gets measured gets improved." Setting up observability tools would provide actionable data through dashboards.

    5) Profiling of the binary:
        I'm learning to profile Go binaries using tools like pprof. Profiling helps understand memory usage and performance bottlenecks – interesting stuff!    

#### Regards:

     I tried my best to follow best practices that I have picked by looking at open source software. This was fun assignment, really looking forward to talk more. 

     Thank you,
     @@author adityadafe