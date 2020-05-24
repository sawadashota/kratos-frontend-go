package frontend

import input

default allow = false

allow {
    input.method == "GET"
    input.path = ["my", "salary"]
    user_owns_token
}

allow {
    some user
    input.method == "GET"
    input.path = ["salary", user]
    input.user == user
}

# Ensure that the token was issued to the user supplying it.
user_owns_token { input.user == token.payload.sub }

# Helper to get the token payload.
token = {"payload": payload} {
  response := http.send({"url": "http://oathkeeper:4456/.well-known/jwks.json", "method": "GET"})
  io.jwt.verify_rs256(input.token, response.raw_body)
  [header, payload, signature] := io.jwt.decode(input.token)
}

