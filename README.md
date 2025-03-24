# passpo








### **Step 1: Create a private key for the CA**
```sh
openssl genrsa -out certs/ca.key 4096
```

### **Step 2: Create a self-signed CA certificate**
```sh
openssl req -x509 -new -nodes -key certs/ca.key -sha256 -days 365 -out certs/ca.crt -subj "/C=US/ST=California/L=San Francisco/O=MyCA/OU=IT/CN=MyCA"
```

This will generate `ca.crt`, which is the root certificate used to sign server and client certificates.

---

## 2️⃣ Generate Server Certificate and Key
### **Step 3: Create a private key for the server**
```sh
openssl genrsa -out certs/server.key 2048
```

### **Step 4: Create a Certificate Signing Request (CSR) for the server**
Replace `yourserver.com` with your actual server's hostname or IP.
```sh
openssl req -new -key certs/server.key -out certs/server.csr -subj "/C=US/ST=California/L=San Francisco/O=MyCompany/OU=IT/CN=localhost"
```

### **Step 5: Sign the server certificate with the CA**
```sh
openssl x509 -req -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt -days 365 -sha256
```

This will generate `server.crt`, which the gRPC server will use for secure communication.

---

## 3️⃣ Generate Client Certificate and Key
### **Step 6: Create a private key for the client**
```sh
openssl genrsa -out certs/client.key 2048
```

### **Step 7: Create a Certificate Signing Request (CSR) for the client**
```sh
openssl req -new -key certs/client.key -out certs/client.csr -subj "/C=US/ST=California/L=San Francisco/O=MyCompany/OU=IT/CN=client"
```

### **Step 8: Sign the client certificate with the CA**
```sh
openssl x509 -req -in certs/client.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/client.crt -days 365 -sha256
```

This will generate `client.crt`, which the gRPC client will use to authenticate itself.

---

## 4️⃣ Verify the Certificates
To ensure everything is correctly signed, run:

### **Verify CA certificate**
```sh
openssl x509 -in certs/ca.crt -noout -text
```

### **Verify Server Certificate**
```sh
openssl verify -CAfile certs/ca.crt certs/server.crt
```

### **Verify Client Certificate**
```sh
openssl verify -CAfile certs/ca.crt certs/client.crt
```

If all checks pass, the certificates are correctly signed and ready for use.

---

## 5️⃣ Using the Certificates in gRPC
- **Server:** Uses `server.crt` and `server.key`.
- **Client:** Uses `client.crt` and `client.key`.
- **Both Server and Client** trust `ca.crt`.