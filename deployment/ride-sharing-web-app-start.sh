systemctl daemon-reload
systemctl start ride-sharing-web-app.service
systemctl restart ride-sharing-web-app.service

sleep 3
systemctl status ride-sharing-web-app.service
