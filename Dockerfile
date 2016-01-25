FROM golang:1.4.2-cross
WORKDIR /go/src/github.com/feelobot/dogsights
ADD ./ /go/src/github.com/feelobot/dogsights
RUN go get && go install
CMD dogsights
