*** Settings ***
Library       REST    localhost:9090

*** Test cases ***
Send Message To Slack
  GET         /v1/reminder/send