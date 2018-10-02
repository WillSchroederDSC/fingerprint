# Fingerprint

A customer authentication/authorization microsevice in Go. It allows for user authentication both synchronously (validation via endpoints) and asynchronously (validation via token decryption). 

The bearer token will also be a PAST style token containing information about the session.  
https://github.com/o1egl/paseto

Uses Devise/Rails/has_secure_password style method for password hashing, using bcrypt.  
Based on https://github.com/consyse/go-devise-encryptor

Has the concept of expiring scopes, allowing one session to have multiple groupings of scopes that expire at different times.  
Useful for making customers re-login to perform sensisive actions after a period of time. 

Has the concept of a guest customer.  
This works by generating a customer with an addendum to their email address that fingerprint splits off for you automatically.  
This customer is genearted a password that cannot be recovered.  
Sessions can be requested for guest users to grant them access to things with out having to register.  

## Setup

Secret for hashing.   
Secret for token decoding.      

## Jobs

Cleanup job for session revolks.

## Token Format
```json
{
    version: 1
    customer_id: 1,
    session_id: 1,
    experation: 1538523728
    scope_groupings: [
        {
            scopes: ["read", "comment"],
            experation: 1538523728
        },
        {
            scopes: ["write"],
            experation: 1538523720
        }
    ]
}
```

Version specifies the format of the token.   
Scope groupings are collections of scopes with each set of scopes experation date.
Dates are a unix timestamp.  

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
    (has many password resets)

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
