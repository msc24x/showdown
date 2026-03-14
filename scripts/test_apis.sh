#!/usr/bin/env bash

set -e

########################################
# Configuration
########################################

BASE_URL="http://localhost:7070"
STATUS_URL="$BASE_URL/status"
TOKEN="your-access-token"

########################################
# Colors
########################################

GREEN="\033[0;32m"
RED="\033[0;31m"
BLUE="\033[0;34m"
NC="\033[0m"

########################################
# Helpers
########################################

print_title() {
    echo -e "\n${BLUE}========== $1 ==========${NC}"
}

success() {
    echo -e "${GREEN}✔ $1${NC}"
}

fail() {
    echo -e "${RED}✖ $1${NC}"
}

pretty() {
    if command -v jq >/dev/null 2>&1; then
        jq . 2>/dev/null || cat
    else
        cat
    fi
}

########################################
# Status Endpoint Test
########################################

test_status() {

    print_title "Testing /status"

    response=$(curl -s -w "\n%{http_code}" \
        "$STATUS_URL" \
        -H "Access-Token: $TOKEN")

    body=$(echo "$response" | head -n -1)
    code=$(echo "$response" | tail -n 1)

    if [[ "$code" == "200" ]]; then
        success "/status returned HTTP 200"
    else
        fail "/status returned HTTP $code"
    fi

    echo "$body" | pretty
}

########################################
# Judge Helper
########################################

run_judge() {

    language="$1"
    code="$2"
    input="$3"
    expected="$4"

    print_title "Testing language: $language"

    payload=$(jq -n \
        --arg lang "$language" \
        --arg code "$code" \
        --arg input "$input" \
        --arg output "$expected" \
'{
  judge_params: {
    donotjudge: false,
    webhook: ""
  },
  exe: {
    language: $lang,
    code: $code,
    input: $input,
    output: $output
  }
}')

    response=$(curl -s -w "\n%{http_code}" \
        "$BASE_URL/judge" \
        -X POST \
        -H "Access-Token: $TOKEN" \
        -H "Content-Type: application/json" \
        -d "$payload")

    body=$(echo "$response" | head -n -1)
    code=$(echo "$response" | tail -n 1)

    if [[ "$code" == "200" ]]; then
        success "$language execution succeeded"
    else
        fail "$language execution failed (HTTP $code)"
    fi

    echo "$body" | pretty
}

########################################
# Language Tests
########################################

test_cpp() {
run_judge "cpp" \
'#include <iostream>
int main() {
    std::cout << "Hello CPP";
    return 0;
}' \
"" \
"Hello CPP"
}

test_c() {
run_judge "c" \
'#include <stdio.h>
int main() {
    printf("Hello C");
    return 0;
}' \
"" \
"Hello C"
}

test_python() {
run_judge "py" \
'print("Hello Python")' \
"" \
"Hello Python"
}

test_go() {
run_judge "go" \
'package main
import "fmt"
func main() {
    fmt.Println("Hello Go")
}' \
"" \
"Hello Go"
}

test_js() {
run_judge "js" \
'console.log("Hello JS")' \
"" \
"Hello JS"
}

test_ts() {
run_judge "ts" \
'console.log("Hello TS")' \
"" \
"Hello TS"
}

########################################
# Run All Tests
########################################

main() {

    print_title "Showdown API Test Suite"

    test_status

    test_cpp
    test_c
    test_python
    test_go
    test_js
    test_ts

    print_title "All tests completed"
}

main