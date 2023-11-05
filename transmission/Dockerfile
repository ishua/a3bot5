FROM python:3.11

RUN mkdir /srv/data

ADD transmission/app /srv
COPY transmission/app/ /srv/app/
ADD transmission/main.py /srv 

ADD transmission/requirements.txt /srv

WORKDIR /srv
RUN pip install -r requirements.txt

CMD ["python", "/srv/main.py"]
