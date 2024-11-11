systemctl daemon-reload
systemctl start ride-sharing-api.service
systemctl restart ride-sharing-api.service

sleep 3
systemctl status ride-sharing-api.service
