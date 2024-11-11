cp --force deployment/ride-sharing-web-app.service /etc/systemd/system/
systemctl daemon-reload
systemctl start ride-sharing-web-app.service
systemctl restart ride-sharing-web-app.service
