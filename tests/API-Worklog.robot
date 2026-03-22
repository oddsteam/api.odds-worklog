*** Settings ***
Library       REST    localhost:9090

*** Test cases ***
# Reminder/send endpoint was removed with the Slack worker.
# ID card email is sent via POST /v1/reminder/mail/{id} (requires auth).
