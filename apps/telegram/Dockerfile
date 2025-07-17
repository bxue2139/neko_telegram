ARG BASE_IMAGE=ghcr.io/m1k1o/neko/base:latest
FROM $BASE_IMAGE

#
# install telegram desktop
RUN set -eux; apt-get update; \
    apt-get install -y --no-install-recommends \
        openbox \
        wget \
        curl \
        ca-certificates \
        xz-utils \
        libgtk-3-0 \
        libxss1 \
        libgconf-2-4 \
        libxtst6 \
        libxrandr2 \
        libasound2 \
        libpangocairo-1.0-0 \
        libatk1.0-0 \
        libcairo-gobject2 \
        libgtk-3-0 \
        libgdk-pixbuf2.0-0; \
    #
    # Try to install from package manager first (more reliable)
    if apt-get install -y --no-install-recommends telegram-desktop; then \
        echo "Telegram installed from package manager"; \
    else \
        echo "Package manager installation failed, trying direct download"; \
        ARCH=$(dpkg --print-architecture); \
        if [ "${ARCH}" = "amd64" ]; then \
            # Use a more reliable download method with retries
            for i in 1 2 3; do \
                if wget --timeout=30 --tries=3 -O /tmp/telegram.tar.xz "https://telegram.org/dl/desktop/linux"; then \
                    break; \
                elif [ $i -eq 3 ]; then \
                    echo "Failed to download Telegram after 3 attempts"; \
                    exit 1; \
                fi; \
                sleep 5; \
            done; \
            mkdir -p /opt/telegram; \
            tar -xf /tmp/telegram.tar.xz -C /opt/telegram --strip-components=1; \
            ln -s /opt/telegram/Telegram /usr/bin/telegram-desktop; \
            rm -f /tmp/telegram.tar.xz; \
        else \
            echo "Unsupported architecture for direct download: ${ARCH}"; \
            exit 1; \
        fi; \
    fi; \
    #
    # create telegram data directory
    mkdir -p /home/neko/.local/share/TelegramDesktop; \
    chown -R neko:neko /home/neko/.local; \
    #
    # clean up
    apt-get clean -y; \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/*

#
# copy configuration files
COPY supervisord.conf /etc/neko/supervisord/telegram.conf
COPY openbox.xml /etc/neko/openbox.xml
