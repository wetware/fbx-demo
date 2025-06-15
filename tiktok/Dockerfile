FROM python:3.9

ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1

# Everything we need to run TikTok-API and pycapnp
RUN apt-get update && apt-get install -y \
    gcc \
    cmake \
    capnproto \
    libcapnp-dev \
    && rm -rf /var/lib/apt/lists/*

COPY ./requirements.txt /requirements.txt

RUN pip install --no-cache-dir -r /requirements.txt

COPY . /tiktok

RUN mkdir -p /tiktok/api/cap

WORKDIR /

CMD ["python", "-c", "print('I should be replaced in docker-compose.yaml')"]
