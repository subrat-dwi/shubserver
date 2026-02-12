🧠 Core principle (memorize this)

The server never sees secrets. Ever.
It only stores opaque blobs.

Everything else follows from this.

🧩 Actors in the system
1. Client (CLI / Web / Telegram bot)

Knows master password

Derives encryption keys

Encrypts / decrypts secrets

Holds JWT for auth

2. Server (your Go backend)

Authenticates users

Stores encrypted data

Never decrypts anything

3. Database

Stores ciphertext + nonce

Useless without client keys

🔐 Cryptographic building blocks
Purpose	Algorithm
Password → key	Argon2id
Encryption	AES-256-GCM
Login hashing	bcrypt
Auth	JWT
IDs	UUID
🟢 PHASE 1: User registration (first time)
Step 1: User enters

Email

Master password

⚠️ Master password is NOT the login password
(You can combine them later, but keep them separate now)

Step 2: Client derives keys (locally)

Client runs:

master_key = Argon2id(
  password = master_password,
  salt = user_email (or random salt),
  memory = high,
  iterations = moderate
)


This produces 256-bit key.

🧠 Important:

This key is deterministic

Same password + salt → same key every time

Server never sees it

Step 3: Client sends only login data

Client sends to server:

{
  "email": "...",
  "password": "login-password"
}


Server:

Hashes password with bcrypt

Stores email + hash

🚫 Master password never leaves device.

🟢 PHASE 2: Login
Step 1: Client sends login credentials

Server:

Verifies bcrypt hash

Issues JWT

{
  "token": "eyJhbGciOi..."
}

Step 2: Client re-derives master key

Client:

Asks user for master password

Runs same Argon2id

Gets same encryption key

🧠 No storage needed. Pure math.

🟢 PHASE 3: Storing a password

User wants to store:

Service: Gmail
Username: user@gmail.com
Password: Gmail123!

Step 1: Client prepares plaintext
{
  "service": "gmail",
  "username": "user@gmail.com",
  "password": "Gmail123!"
}

Step 2: Client encrypts locally

Client:

Generates random nonce (12 bytes)

Encrypts using AES-GCM

ciphertext, tag = AES-256-GCM(
  key = master_key,
  nonce = random_nonce,
  plaintext = json_data
)

Step 3: Client sends encrypted blob
{
  "ciphertext": "...",
  "nonce": "...",
  "version": 1
}


Server:

Stores it

Cannot read it

Cannot modify it (GCM detects tampering)

🟢 PHASE 4: Fetching passwords
Step 1: Client requests data
GET /passwords
Authorization: Bearer <JWT>


Server:

Returns encrypted blobs

Step 2: Client decrypts

For each entry:

plaintext = AES-256-GCM-DECRYPT(
  key = master_key,
  nonce = stored_nonce,
  ciphertext = stored_ciphertext
)


Client displays password.

🟢 PHASE 5: Updating a password

Same as create:

New nonce

New ciphertext

Old one replaced

🟢 PHASE 6: Deleting a password

Server deletes row.
No crypto involved.

🟢 PHASE 7: Changing master password (hard part)
Step 1: Client decrypts everything

Using old key

Step 2: Client derives new key

From new master password

Step 3: Client re-encrypts everything

New nonces

New ciphertexts

Step 4: Uploads updated blobs

Server still never sees plaintext.

🟢 PHASE 8: Multiple clients (CLI, Telegram, Web)

Each client:

Requests encrypted data

Asks user for master password

Derives key

Decrypts locally

No syncing secrets needed.

🔐 Why this design is strong
Threat	Result
DB leaked	Data unreadable
Server hacked	Still unreadable
Insider attack	No access
MITM	TLS + GCM
Replay attack	JWT expiry
🚨 What happens if user forgets master password?

Data permanently encrypted

No recovery

This is correct behavior

🧠 Key mental model (burn this in)

Passwords are encrypted data, not records.
The server is a dumb storage engine.

🧱 Where code lives (your Go project)
internal/
  auth/        → login, JWT
  users/       → user records
  passwords/
    model.go
    repo.go
    service.go
    handler.go
crypto/        → AES + Argon2 helpers

❌ What NOT to do

Never store master password

Never derive keys on server

Never reuse nonce

Never decrypt on backend