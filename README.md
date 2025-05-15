# WSO2 Asgardeo - Project Thunder ⚡
### The Lighting Fast Identity Management Suite

<a alt="Asgardeo logo" href="https://asgardeo.io" target="_blank" rel="noreferrer" style="background: #fff; display: inline-block"><img src="asgardeo-logo.png" width="400"></a>

**Asgardeo - Project Thunder** is a modern, identity management service by WSO2. It empowers you to design tailored login, registration, and recovery flows using a flexible identity flow designer.

Thunder secures users, applications, services, and AI agents by managing their identities and offering a complete suite of supporting capabilities.

Designed for extensibility, scalability, and seamless containerized deployment, Thunder integrates naturally with microservices and DevOps environments—serving as the core identity layer for your cloud platform.

---

## 🚀 Features (WIP)

- ✅ **Standards-Based**
  - OAuth 2.1, OpenID Connect (OIDC)
  - SCIM 2.0
- 🛠️ **Visual Identity Flow Designer**
- 👤 **User & Identity Management**
- 🔗 **Social Login**
- 🔐 **Multi-Factor Authentication (MFA)**
- 🌐 **RESTful APIs**

---

## ⚡ Quickstart

### ✅ Prerequisites

- Go 1.23+
- cURL
- Node.js 14+
- PNPM 10+
- React 19+

---

### 🛠 Step 1: Build and Run the Product

```bash
make run
```

---

### 🔑 Step 2: Try Out the Product

#### 1️⃣ Try Out Client Credentials Flow

```bash
curl -k -X POST https://localhost:8090/oauth2/token \
  -H 'Authorization: Basic Y2xpZW50MTIzOnNlY3JldDEyMw==' \
  -d 'grant_type=client_credentials'
```

- **Client ID:** `client123`
- **Client Secret:** `secret123`

#### 2️⃣ Try Out Authorization Code Flow

- Open the following URL in your browser:

  ```bash
  https://localhost:8090/oauth2/authorize?response_type=code&client_id=client123&redirect_uri=https://localhost:3000&scope=openid&state=state_1
  ```

- Enter the following credentials:

  - **Username:** `thor`
  - **Password:** `thor123`

    **Note:** The credentials can be configured in the `repository/conf/deployment.yaml` file under the `user_store` section.

- After successful authentication, you will be redirected to the redirect URI with the authorization code and state.

  ```bash

  https://localhost:3000/?code=<code>&state=state_1
  ```

- Copy the authorization code and exchange it for an access token using the following cURL command:

  ```bash
  curl -k --location 'https://localhost:8090/oauth2/token' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header 'Authorization: Basic Y2xpZW50MTIzOnNlY3JldDEyMw==' \
  --data-urlencode 'grant_type=authorization_code' \
  --data-urlencode 'redirect_uri=https://localhost:3000' \
  --data-urlencode 'code=<code>'
  ```

  - **Client ID:** `client123`
  - **Client Secret:** `secret123`

---

## 🧪 Running Integration Tests

Building the product with `make all` will run the integration tests by default. However if you want to run the tests manually, follow the steps below.

### 1️⃣ Build the Product

```bash
make clean build
```

### 2️⃣ Run the Tests

```bash
make test
```

---

## Running Development Environment

```bash
make run
```

## License

Licenses this source under the Apache License, Version 2.0 ([LICENSE](LICENSE)), You may not use this file except in compliance with the License.

---------------------------------------------------------------------------
(c) Copyright 2025 WSO2 LLC.
