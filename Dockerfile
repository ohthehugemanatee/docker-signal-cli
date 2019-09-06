FROM golang:latest
WORKDIR /go/src/github.com/ohthehugemanatee/signal-cli-container/
COPY . .
RUN set -e \
&& go get ./... \
&& go build -o signal-nixplay-bridge .


FROM openjdk:14-alpine
ENV VERSION "0.6.2"
RUN set -e \
  && apk add ca-certificates tzdata wget tar gzip \
  && rm -rf /var/cache/apk/* \
  && wget https://github.com/AsamK/signal-cli/releases/download/v"${VERSION}"/signal-cli-"${VERSION}".tar.gz -O /tmp/signal-cli.tar.gz \
  && tar -xzf /tmp/signal-cli.tar.gz -C /opt \
  && ln -sf /opt/signal-cli-"${VERSION}"/bin/signal-cli /usr/local/bin/
COPY --from=0 /go/src/github.com/ohthehugemanatee/signal-cli-container/signal-nixplay-bridge /usr/local/bin
CMD /usr/local/bin/signal-nixplay-bridge -p $MYPHONE -g $GROUPID -e $DESTMAIL -user $SMTPUSER -pass $SMTPPASS -s $SMTPSERVER -f $SMTPFROM