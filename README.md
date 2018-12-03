# Fingerprint (README out of date)

A customer authentication/authorization microsevice in Go. It allows for user authorization both synchronously and asynchronously. 

To validate synchronously, simply give the token to the validate session endpoint, which returns information about the session.

To validate asynchronously, decrypt the bearer token, which is also a JWT token. If you want to validate if its been revolked, call the validate session endpoint. Because you already have the session information from the token decryption, the validation can either be skipped or allowed to fail. 

The bearer token is a JWT token, it can be decrypted to provide information about the session (see below).

Passwords are hashed via bcrypt.

The backing database is postgres. 

Has the concept of expiring scopes, allowing one session to have multiple groupings of scopes that expire at different times.  
Useful for making customers re-login to perform sensisive actions after a period of time. 

Has the concept of a guest customer.  
This works by generating a customer with an addendum to their email address that fingerprint splits off for you automatically.  
This customer is generated a password that cannot be recovered, and is never exposed to any client. 
Sessions can be requested for guest users to grant them access to things with out having to register.  

## Setup

Secret for hashing.   
Secret for token decoding.      

## Token Format
```javascript
{
    "version": 1,
    "session": {
        "customer_id": 1,
        "session_id": 1,
        "expiration": 1538523728,
        "is_guest": false,
        "scope_groupings": [
            {
                "scopes": ["read", "comment"],
                "expiration": 1538523728
            },
            {
                "scopes": ["write"],
                "expiration": 1538523720
            }
        ]
    }
}
```

Version specifies the format of the token.   
Scope groupings are collections of scopes with each set of scopes experation date.
Dates are a unix timestamp.  

## GRPC/API Endpoints (Not for external use)

### User Exists  
    Request: email
    Respone: status 

### Create User
    Request: email, password, scopes 
    Response: token 

### Create Guest User
    Appends .guest.random_id, scrambles a password it does not tell you
    Request: email, scopes
    Response: token 

### Create Session
    Request: email, password, scopes  
    Response: token  

### Update Password
    Request: reset token, new password 
    Response: status 

### Validate Session
    Request: token
    Response: status, scopes 

### Create Password Reset Token
    Request: email
    Response: reset token

### Create Session Revoke 
    Request: session_id OR customer_id
    Response: status 

## Tables
### Customers
| Field | Type |
|---| --- |
| uuid  |
| email  |
| reset_token  |
| first_name  |
| last_name   |
| is_guest   |
| updated_at |
| created_at   |

* Has many Sessions
* Has many PasswordResets 

## Sessions
| Field | Type |
|---| --- |
| customer_id |
| uuid  |
| experation  |
| updated_at |
| created_at   |

* Has many ScopeGroupings
* Belongs to a Customer

## ScopeGroupings
| Field | Type |
|---| --- |
| session_id |
| uuid  |
| scopes  | [String] | 
| updated_at |
| created_at   |

* Belongs to a Session

## SessionRevokes
| Field | Type |
|---| --- |
| session_id |
| uuid  |
| updated_at |
| created_at   |

* Belongs to a Session

## PasswordResets
| Field | Type |
|---| --- |
| customer_id |
| uuid  |
| reset_hash |
| updated_at |
| created_at   |

* Belongs to a Customer
