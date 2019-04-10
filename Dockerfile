FROM nicholasjackson/gocv-alpine:4.0.0-buildstage as build-stage

ENV GO111MODULE on

COPY . $GOPATH/src/github.com/emojify-app/face-detection
RUN cd $GOPATH/src/github.com/emojify-app/face-detection && go get ./... && go build -o facedetection .

FROM nicholasjackson/gocv-alpine:runtime

RUN mkdir /app
COPY --from=build-stage $GOPATH/src/github.com/emojify-app/face-detection/facedetection /app/facedetection
COPY ./cascades /app/cascades

EXPOSE 9090

ENV PKG_CONFIG_PATH /usr/local/lib64/pkgconfig
ENV LD_LIBRARY_PATH /usr/local/lib64
ENV CGO_CPPFLAGS -I/usr/local/include
ENV CGO_CXXFLAGS "--std=c++1z"
ENV CGO_LDFLAGS "-L/usr/local/lib -lopencv_core -lopencv_face -lopencv_videoio -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_xfeatures2d -lopencv_plot -lopencv_tracking"
ENV CASCADE_FOLDER /app/cascades

ENTRYPOINT /app/facedetection
