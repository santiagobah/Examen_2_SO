FROM kalilinux/kali-rolling:latest


ENV DEBIAN_FRONTEND=noninteractive


RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    unzip \
    xz-utils \
    fontconfig \
    libfontconfig1 \
    libfreetype6 \
    libgl1 \
    libglu1-mesa \
    libasound2 \
    libpulse0 \
    libxcursor1 \
    libxinerama1 \
    libxrandr2 \
    libxi6 \
    libxxf86vm1 \
    libglx0 \
    libegl1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && fc-cache -fv


RUN wget --no-verbose https://github.com/godotengine/godot/releases/download/4.4-stable/Godot_v4.4-stable_linux.arm64.zip -O /tmp/godot.zip \
    && unzip /tmp/godot.zip -d /tmp \
    && mv /tmp/Godot_v4.4-stable_linux.arm64 /usr/local/bin/godot \
    && chmod +x /usr/local/bin/godot \
    && rm /tmp/godot.zip


RUN mkdir -p /usr/share/fonts/truetype \
    && ln -s /etc/fonts /usr/share/fontconfig

WORKDIR /project

ENTRYPOINT ["godot"]
CMD ["--headless"]
