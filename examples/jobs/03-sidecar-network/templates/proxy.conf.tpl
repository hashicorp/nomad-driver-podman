server {
    listen 8080;
    server_name localhost;

    location / {
        # The app shares this network namespace and listens on localhost:5678.
        proxy_pass http://127.0.0.1:5678;
        proxy_set_header Host $host;
    }
}
