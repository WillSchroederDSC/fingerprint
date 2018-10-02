# Fingerprint

A customer authentication/authorization microsevice in Go. 

Allows for sync and async (via decrypting token) user authentication

Uses a PAST token to store information about a session. 
https://github.com/o1egl/paseto

Uses Devise/Rails/has_secure_password style method for password hashing, uses bcrypt
Based on https://github.com/consyse/go-devise-encryptor

Has the concept of expiring scopes

Has the concept of generating a guest user 

## Config

Secret for hashing
Secret for token decoding 

## Jobs

Cleanup job for session revolks 

## Token Format
`{
    version: 1
    customer_id: 1,
    session_id: 1,
    scope_groupings: [
        {
            scopes: ["read", "comment"],
            experation: timestamp(format?)
        },
        {
            scopes: ["write"],
            experation: timestamp 
        }
    ]
}`

Version specifies the format of the token. 
Customer id is the obfuscated id
Session id is the obfuscated id
Scopes is an array of strings 

## GRPC/API Endpoints (Not for external use)

User Exists
Request: email
Respone: status 

Create User
Request: email, password, scopes 
Response: token 

Create Guest User
Appends .guest.random_id, scrambles a password it does not tell you
Request: email, scopes
Response: token 

Create Session
Request: email, password, scopes 
Response: token 

Update Password
Request: reset token, new password 
Response: status 

Validate Session
Request: token
Response: status, scopes 

Create Reset Token
Request: email
Response: reset token

Create Session Revoke 
Request: session id
Response: status 

Is Session Revoked (To support async validation)
Request: session_id
Response: status

## Tables
Customers
    obfuscated id
    email
    reset token
    first name
    last name
    is_guest
    created_at
    (has many sessions)
    (has man password resets)

Sessions
    customer_id
    obfuscated id
    experation 
    created_at
    (has many Scopes)
    (belongs to a customer)

ScopeGroupings
    session_id
    obfuscated id
    scopes [string array]
    experation 
    (belongs to a Session)

SessionRevokes
    revoked session_id
    obfuscated id
    created_at
    remove_at (Default, 90 days)

PasswordResets
    customer_id
    obfuscated id
    reset_hash
