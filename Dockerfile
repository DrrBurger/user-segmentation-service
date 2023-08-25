FROM ubuntu:latest
LABEL authors="dr.burger"

ENTRYPOINT ["top", "-b"]