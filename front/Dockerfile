FROM node:lts-alpine

COPY . /app 

WORKDIR /app

RUN npm install

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
