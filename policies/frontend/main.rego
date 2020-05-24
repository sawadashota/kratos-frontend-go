package frontend

import input

default allow = false

allow {
    input.method == "GET"
    input.path = ["my", "salary"]
}

allow {
    some user
    input.method == "GET"
    input.path = ["salary", user]
    input.user == user
}
