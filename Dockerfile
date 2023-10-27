FROM gcr.io/distroless/base-debian12:debug-nonroot
COPY --chown=nonroot:nonroot repo-util /usr/sbin/repo-util
ENTRYPOINT ["sh", "-c" , "repo-util summary"]