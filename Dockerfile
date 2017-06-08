FROM scratch

ADD ./ovpndns /

ENTRYPOINT [ "/ovpndns" ]
