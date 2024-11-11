cp --force deployment/ride-sharing-api.service /etc/systemd/system/
systemctl daemon-reload
systemctl start ride-sharing-api.service
systemctl restart ride-sharing-api.service
