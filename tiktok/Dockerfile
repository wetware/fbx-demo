FROM python:3.9

# Everything we need to run TikTok-API and pycapnp
RUN apt-get update && apt-get install -y \
    gcc \
    cmake \
    capnproto \
    libcapnp-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /

COPY ./requirements.txt /requirements.txt

RUN pip install --no-cache-dir -r /requirements.txt

COPY . /tiktok

CMD ["python", "-c", "print('I should be replaced in docker-compose.yaml')"]
