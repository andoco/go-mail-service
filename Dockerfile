FROM golang
 
ADD . /go/src/bitbucket.org/andoco/gomailservice
RUN go install bitbucket.org/andoco/gomailservice
ENTRYPOINT /go/bin/gomailservice
 
EXPOSE 8000
