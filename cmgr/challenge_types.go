package cmgr

func (m *Manager) getDockerfile(challengeType string) []byte {
	if challengeType == "custom" {
		return nil
	}

	return m.challengeDockerfiles[challengeType]
}

func (m *Manager) initDockerfiles() {
	m.challengeDockerfiles = make(map[string][]byte)
	m.challengeDockerfiles["hacksport"] = []byte(hacksportDockerfile)
	m.challengeDockerfiles["flask"] = []byte(flaskDockerfile)
	m.challengeDockerfiles["php"] = []byte(phpDockerfile)
	m.challengeDockerfiles["solver"] = []byte(solverDockerfile)
}

const hacksportDockerfile = `
FROM aci/hacksport

COPY . /app

ARG FORMAT
ARG FLAG
ARG SEED
`

const flaskDockerfile = `
FROM ubuntu:20.04

RUN apt-get update
RUN apt-get -y install python3-pip build-essential
RUN pip3 install flask
RUN groupadd -r flask && useradd -r -d /app -g flask flask

ENV FLASK_RUN_HOST=0.0.0.0
ENV FLASK_RUN_PORT=8000

# End of shared layers for all flask challenges

COPY Dockerfile packages.txt* .
RUN if [ -f packages.txt ]; then xargs -a packages.txt sudo apt install -y; fi

COPY Dockerfile requirements.txt* .
RUN if [ -f requirements.txt ]; then pip3 install -r requirements.txt; fi

COPY --chown=flask:flask . /app

# End of share layers for all builds of the same flask challenge

ARG FLAG
ARG SEED

RUN install -d -m 0700 /challenge && \
    echo "{\"flag\":\"$FLAG\"}" > /challenge/metadata.json

USER flask:flask

RUN sed -i -e "s|{{flag}}|$FLAG|g"                                           \
           -e "s|{{secret_key}}|$(echo $FLAG | sha256sum | cut -d' ' -f1)|g" \
           -e "s|{{seed}}|$SEED|g"                                           \
        /app/app.py

WORKDIR /app
CMD flask run

EXPOSE 8000
# PUBLISH 8000 AS http
`

const phpDockerfile = `
FROM ubuntu:20.04

RUN apt-get update
RUN apt-get -y install php build-essential
RUN groupadd -r php && useradd -r -d /app -g php php

# End of shared layers for all php challenges

COPY Dockerfile packages.txt* .
RUN if [ -f packages.txt ]; then xargs -a packages.txt sudo apt install -y; fi

COPY --chown=php:php . /app

# End of share layers for all builds of the same php challenge

ARG FLAG
ARG SEED

RUN install -d -m 0700 /challenge && \
    echo "{\"flag\":\"$FLAG\"}" > /challenge/metadata.json

USER php:php

RUN find /app -name \( -name *.php -o -name *.txt -o -name *.html \) \
              -exec sed -i -e "s|{{flag}}|$FLAG|g"                   \
                           -e "s|{{seed}}|$SEED|g"                   \
                        {} \;

WORKDIR /app
CMD php -S 0.0.0.0:8000 -t /app

EXPOSE 8000
# PUBLISH 8000 AS http
`

const solverDockerfile = `
FROM ubuntu:20.04

RUN apt-get update
RUN apt-get -y install python3-pip python3-dev git libssl-dev libffi-dev build-essential
RUN pip3 install pwntools

COPY Dockerfile packages.txt* .
RUN if [ -f packages.txt ]; then xargs -a packages.txt sudo apt install -y; fi

COPY Dockerfile requirements.txt* .
RUN if [ -f requirements.txt ]; then pip3 install -r requirements.txt; fi

COPY . /solve
WORKDIR /solve

CMD python3 solve.py
`
