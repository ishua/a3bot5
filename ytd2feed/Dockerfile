FROM python:3.11

EXPOSE 8092
RUN apt-get -y update \
    && apt-get install -y ffmpeg  \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir /srv/data

ADD app /srv
COPY app/ /srv/app/
ADD main.py /srv 

ADD requirements.txt /srv

WORKDIR /srv
RUN pip install -r requirements.txt

CMD ["python", "/srv/main.py"]
