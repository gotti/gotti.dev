FROM node as builder

WORKDIR /frontend

COPY ./frontend /frontend

RUN npm install
RUN npm run build
RUN ls

FROM nginx

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /frontend/out/ /srv/http/
COPY ./contents/ /srv/http/
