FROM docker.io/library/golang:1.22.7 AS build
COPY . /build
RUN make -C /build


FROM scratch
COPY --from=build /build/dist /
CMD [ "/bin/ozon-test-task" ]
