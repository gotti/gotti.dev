FROM node as builder

WORKDIR /frontend

COPY ./contents /frontend
COPY ./frontend /frontend

RUN npm install
RUN npm run build

RUN rm -rf /frontend/node_modules

FROM nginx

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /frontend/ /srv/http/
COPY --from=builder /frontend/out/ /srv/http/
