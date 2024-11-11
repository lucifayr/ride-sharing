addr="$1"
name="$2"
key="$3"
sha="$4"

url="http://$addr/update/start/$name"
curl -X POST -H "SECRET_KEY: $key" -H "GIT_SHA: $sha" "$url"
