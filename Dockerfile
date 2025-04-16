#Use the official Go image
FROM golang:latest

#Set the working directory inside the container
WORKDIR /app

#Copy go.mod and go.sum first (for caching dependencies)
#COPY go.mod go.sum ./

#Copy go file into the directory
COPY go.* ./

#Download dependencies
RUN go mod download

#Copy the rest of the code
COPY . .

RUN go build -o main main.go

#Command to run the application
CMD ["go", "run", "main.go"]


#Command to buld docker image
#docker build -t my-go-test .

#Command to run docker container
#docker run -d -p 8080:8080 my-go-test