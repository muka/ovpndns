FROM scratch

ADD ./build/ovpndns-amd64 /ovpndns

ENTRYPOINT [ "/ovpndns" ]
