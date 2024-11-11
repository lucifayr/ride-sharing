addr="$1"
name="$2"
key="$3"
poll_id="$4"

url="http://$addr/update/logs/$name"
out=$(curl -H "SECRET_KEY: $key" -H "POLL_ID: $poll_id" "$url" 2>/tmp/poll.err || echo "STOPPED due to curl error")
echo -n "$out"

while !(echo "$out" | /bin/grep -q -E "(STOPPED|DONE)"); do
    sleep 3
    out=$(curl -H "SECRET_KEY: $key" -H "POLL_ID: $poll_id" "$url" 2>/tmp/poll.err || echo "STOPPED due to curl error")
    echo -n "$out"
done
