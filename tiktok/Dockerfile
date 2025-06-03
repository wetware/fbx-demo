FROM python:3.9

# Everything we need to run TikTok-API and pycapnp
RUN apt-get update && apt-get install -y \
    gcc \
    cmake \
    capnproto \
    libcapnp-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /

RUN pip install --no-cache-dir \
    TikTokApi

RUN pip install --no-cache-dir \
    pycapnp

CMD ["python", "-c", "print('Hi :) I should be replaced in docker-compose.yaml')"]
