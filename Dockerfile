# pull the base go image 1.20.4
# create a directory known as app
# make app the working directory
# copy the go mod & go sum files i the working directory
# download the dependencies and verify them all
# make a copy of what is in the parent folder into app directory
# build your code from the main folder/file in the parent folder to a executable file brokerApp
# run the executable file
#Create a default base image -alpine to execute the build image

FROM golang:latest as builder

RUN mkdir /app
WORKDIR /app

COPY go.mod /app
COPY go.sum /app

RUN go mod download && go mod verify
COPY . /app

RUN CGO_ENABLED=0 go build -o travasAdmin ./cmd/web
RUN chmod +x /app/travasAdmin

CMD [ "/app/travasAdmin" ]
