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

###Create Session
    Request: email, password, scopes 
    Response: token 

### Update Password
    Request: reset token, new password 
    Response: status 

### Validate Session
    Request: token
    Response: status, scopes 

### Create Reset Token
    Request: email
    Response: reset token

### Create Session Revoke 
    Request: session id
    Response: status 

## Tables
### Customers
| Field | Type |
|---| --- |
| obfuscated_id  |
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
| obfuscated_id  |
| experation  |
| updated_at |
| created_at   |

* Has many ScopeGroupings
* Belongs to a Customer

## ScopeGroupings
| Field | Type |
|---| --- |
| session_id |
| obfuscated_id  |
| scopes  | [String] | 
| updated_at |
| created_at   |

* Belongs to a Session

## SessionRevokes
| Field | Type |
|---| --- |
| session_id |
| obfuscated_id  |
| updated_at |
| created_at   |

* Belongs to a Session

## PasswordResets
| Field | Type |
|---| --- |
| customer_id |
| obfuscated_id  |
| reset_hash |
| updated_at |
| created_at   |

* Belongs to a Customer
