FROM ubuntu

RUN apt-get install mkisofs
RUN cp /usr/bin/mkisofs /usr/local/bin/mkisofs
ADD ./mkclouddrive /bin/mkclouddrive
RUN chmod +x /bin/mkclouddrive

ENTRYPOINT ["/bin/mkclouddrive"]
CMD ["-help"]
