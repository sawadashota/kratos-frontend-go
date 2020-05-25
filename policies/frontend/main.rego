package frontend

import input

default allow = false

allow {
    verify_token
}

# Ensure that the token was issued to the user supplying it.
# user_owns_token { input.user == token.payload.sub }

verify_token {
  response := http.send({"url": "http://oathkeeper:4456/.well-known/jwks.json", "method": "GET"})
  io.jwt.verify_rs256(input.token, response.raw_body)
}

# Helper to get the token payload.
token = {"payload": payload} {
  [header, payload, signature] := io.jwt.decode(input.token)
}

