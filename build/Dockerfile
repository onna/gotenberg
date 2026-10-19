# ARG instructions do not create additional layers. Instead, next layers will
# concatenate them. Also, we have to repeat ARG instructions in each build
# stage that uses them.
ARG GOLANG_VERSION=1.21

# ----------------------------------------------
# Gotenberg binary build stage
# ----------------------------------------------
FROM golang:$GOLANG_VERSION AS binary-stage

ARG GOTENBERG_VERSION=8.0.2
ENV CGO_ENABLED 0

# Define the working directory outside of $GOPATH (we're using go modules).
WORKDIR /home

# Install module dependencies.
COPY go.mod go.sum ./

RUN go mod download &&\
    go mod verify

# Copy the source code.
COPY cmd ./cmd
COPY pkg ./pkg

RUN go build -o gotenberg -ldflags "-X 'github.com/gotenberg/gotenberg/v8/cmd.Version=$GOTENBERG_VERSION'" cmd/gotenberg/main.go

# ----------------------------------------------
# Final stage
# ----------------------------------------------
FROM debian:13-slim

ARG GOTENBERG_VERSION=8.0.2
ARG GOTENBERG_USER_GID=1001
ARG GOTENBERG_USER_UID=1001
ARG NOTO_COLOR_EMOJI_VERSION=v2.038
# First trixie build with the fixes for CVE-2026-63260/63266/63267/63268/63269/63277 (DSA-6543-1).
# Keep in sync with patchedLibreOfficeVersion in pkg/modules/libreoffice/image_test.go.
ARG LIBREOFFICE_MIN_VERSION="4:25.2.3-2+deb13u8"

LABEL org.opencontainers.image.title="Gotenberg" \
      org.opencontainers.image.description="A Docker-powered stateless API for PDF files." \
      org.opencontainers.image.version="$GOTENBERG_VERSION" \
      org.opencontainers.image.authors="Julien Neuhart <neuhart.julien@gmail.com>" \
      org.opencontainers.image.documentation="https://gotenberg.dev" \
      org.opencontainers.image.source="https://github.com/gotenberg/gotenberg"

RUN \
    # Create a non-root user.
    # All processes in the Docker container will run with this dedicated user.
    groupadd --gid "$GOTENBERG_USER_GID" gotenberg &&\
    useradd --uid "$GOTENBERG_USER_UID" --gid gotenberg --shell /bin/bash --home /home/gotenberg --no-create-home gotenberg &&\
    mkdir /home/gotenberg &&\
    chown gotenberg: /home/gotenberg

RUN \
    # Install system dependencies required for the next instructions or debugging.
    # Note: tini is a helper for reaping zombie processes.
    apt-get update -qq &&\
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends ca-certificates curl gnupg tini python3 &&\
    # Cleanup.
    # Note: the Debian image does automatically a clean after each install thanks to a hook.
    # Therefore, there is no need for apt-get clean.
    # See https://stackoverflow.com/a/24417119/3248473.
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN \
    # Install fonts.
    # Credits:
    # https://github.com/arachnys/athenapdf/blob/master/cli/Dockerfile.
    # https://help.accusoft.com/PrizmDoc/v12.1/HTML/Installing_Asian_Fonts_on_Ubuntu_and_Debian.html.
    curl -o ./ttf-mscorefonts-installer_3.8.1_all.deb http://httpredir.debian.org/debian/pool/contrib/m/msttcorefonts/ttf-mscorefonts-installer_3.8.1_all.deb &&\
    apt-get update -qq &&\
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends \
    ./ttf-mscorefonts-installer_3.8.1_all.deb \
    culmus \
    fonts-beng \
    fonts-hosny-amiri \
    fonts-lklug-sinhala \
    fonts-lohit-guru \
    fonts-lohit-knda \
    fonts-samyak-gujr \
    fonts-samyak-mlym \
    fonts-samyak-taml \
    fonts-sarai \
    fonts-sil-abyssinica \
    fonts-sil-padauk \
    fonts-telu \
    fonts-thai-tlwg \
    ttf-wqy-zenhei \
    fonts-arphic-ukai \
    fonts-arphic-uming \
    fonts-ipafont-mincho \
    fonts-ipafont-gothic \
    fonts-unfonts-core \
    # LibreOffice recommends.
    fonts-crosextra-caladea \
    fonts-crosextra-carlito \
    fonts-dejavu \
    fonts-dejavu-extra \
    fonts-liberation \
    fonts-liberation2 \
    fonts-linuxlibertine \
    fonts-noto-cjk \
    fonts-noto-core \
    fonts-noto-mono \
    fonts-noto-ui-core \
    fonts-sil-gentium \
    fonts-sil-gentium-basic &&\
    rm -f ./ttf-mscorefonts-installer_3.8.1_all.deb &&\
    # Add Color and Black-and-White Noto emoji font.
    # Credits:
    # https://github.com/gotenberg/gotenberg/pull/325.
    # https://github.com/googlefonts/noto-emoji.
    curl -Ls "https://github.com/googlefonts/noto-emoji/raw/$NOTO_COLOR_EMOJI_VERSION/fonts/NotoColorEmoji.ttf" -o /usr/local/share/fonts/NotoColorEmoji.ttf &&\
    # Cleanup.
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN \
    # Install either Google Chrome stable on amd64 architecture or
    # Chromium on other architectures.
    # See https://github.com/gotenberg/gotenberg/issues/328.
    /bin/bash -c \
    'set -e &&\
    if [[ "$(dpkg --print-architecture)" == "amd64" ]]; then \
      curl -fsSL https://dl.google.com/linux/linux_signing_key.pub | gpg --dearmor -o /usr/share/keyrings/google-chrome.gpg &&\
      echo "deb [signed-by=/usr/share/keyrings/google-chrome.gpg] http://dl.google.com/linux/chrome/deb/ stable main" | tee /etc/apt/sources.list.d/google-chrome.list &&\
      apt-get update -qq &&\
      DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends google-chrome-stable &&\
      mv /usr/bin/google-chrome-stable /usr/bin/chromium; \
    else \
      apt-get update -qq &&\
      DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends chromium; \
    fi' &&\
    # Verify installation.
    chromium --version &&\
    # Cleanup.
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN \
    # Install LibreOffice & unoconverter (v0.5.0 is the first release without distutils, which Python 3.13 removed).
    apt-get update -qq &&\
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends libreoffice &&\
    dpkg --compare-versions "$(dpkg-query -W -f='${Version}' libreoffice-core)" ge "$LIBREOFFICE_MIN_VERSION" &&\
    # Untrusted documents can load JDBC drivers and run them in a JVM, so the image must not ship one.
    ! command -v java &&\
    curl -Ls https://raw.githubusercontent.com/gotenberg/unoconverter/v0.5.0/unoconv -o /usr/bin/unoconverter &&\
    chmod +x /usr/bin/unoconverter &&\
    # unoconverter will look for the Python binary, which has to be at version 3.
    ln -s /usr/bin/python3 /usr/bin/python &&\
    # Verify installations.
    libreoffice --version &&\
    unoconverter --version &&\
    # Cleanup.
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN \
    # Install QPDF (PDF engine). PDFtk is left out because it needs a JVM.
    apt-get update -qq &&\
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends qpdf &&\
    # Verify installations.
    qpdf --version &&\
    # Cleanup.
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Improve fonts subpixel hinting and smoothing.
# Credits:
# https://github.com/arachnys/athenapdf/issues/69.
# https://github.com/arachnys/athenapdf/commit/ba25a8d80a25d08d58865519c4cd8756dc9a336d.
COPY build/fonts.conf /etc/fonts/conf.d/100-gotenberg.conf

# Copy the Gotenberg binary from the binary stage.
COPY --from=binary-stage /home/gotenberg /usr/bin/

# Environment variables required by modules or else.
ENV CHROMIUM_BIN_PATH /usr/bin/chromium
ENV LIBREOFFICE_BIN_PATH /usr/lib/libreoffice/program/soffice.bin
ENV UNOCONVERTER_BIN_PATH /usr/bin/unoconverter
ENV QPDF_BIN_PATH /usr/bin/qpdf

USER gotenberg
WORKDIR /home/gotenberg

# Default API port.
EXPOSE 3000

ENTRYPOINT [ "/usr/bin/tini", "--" ]
CMD [ "gotenberg" ]
