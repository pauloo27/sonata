FROM fedora:39

# Setup the environment
RUN dnf --assumeyes install mingw64-gtk3 mingw64-gcc go glib2-devel
RUN bash -c "sed -i -e 's/-Wl,-luuid/-luuid/g' /usr/x86_64-w64-mingw32/sys-root/mingw/lib/pkgconfig/gdk-3.0.pc"

RUN mkdir -p /root/go/src/sonata /root/go/bin /root/go/pkg

# Setup the source code
WORKDIR /root/go/src/sonata
COPY go.mod go.sum ./
RUN go mod download

RUN PKG_CONFIG_PATH=/usr/x86_64-w64-mingw32/sys-root/mingw/lib/pkgconfig CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go get -v github.com/gotk3/gotk3/gtk

COPY common ./common
COPY gui ./gui

RUN CGO_ENABLED=1 PKG_CONFIG_PATH=/usr/x86_64-w64-mingw32/sys-root/mingw/lib/pkgconfig CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -v -o ./sonata-gui.exe -ldflags -H=windowsgui ./gui

RUN mkdir ./dist
RUN cp sonata-gui.exe dist
RUN cp -r /usr/x86_64-w64-mingw32/sys-root/mingw/bin/* ./dist
RUN cp -r /usr/x86_64-w64-mingw32/sys-root/mingw/share ./dist/
