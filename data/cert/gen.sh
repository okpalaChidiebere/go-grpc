# for our app running in a dev env locally, we will be using a self-signed certificate
# Watch this video https://www.youtube.com/watch?v=T4Df5_cojAs to see how HTTPS works using SSL. TSL is outdated and old :(

# Generate CA's Private Key and Self-Signed Public Key certificate
openssl req -x509 \
  -newkey rsa:4096 \
  -nodes \
  -days 365 \
  -keyout ca_key.pem \
  -out ca_cert.pem \
  -subj /C=CA/ST=SK/L=Saskatoon/O=TestOrganizationName/CN=test-app-name-server_ca/ \
  -sha256

#uncomment if you want to see text format of the certificate
# echo "CA self-signed certificate"
# openssl x509 -in ca_cert.pem -noout -text
#According to https://man.openbsd.org/openssl.1#x509 x509 is used to print the cetificate information

# Generate web server's Private Key and Certificate Signing Request (CSR)
# Notice that we remove the -x509 flag because we dont want this key to be a self-signed certificate
openssl req \
  -newkey rsa:4096 \
  -nodes \
  -keyout server_key.pem \
  -out server_csr.pem \
  -subj /C=CA/ST=SK/L=Saskatoon/O=TestOrganizationName/CN=test-app-name-client_ca/ \
  -sha256

# Use the CA's private key to sign the web server's CSR and get back the signed certificate
openssl x509 -req \
  -CAkey ca_key.pem \
  -in server_csr.pem \
  -days 60 \
  -CA ca_cert.pem \
  -CAcreateserial \
  -out server_cert.pem \
  -extfile server-ext.cnf \
  -sha256

#uncomment if you want to see text format of the certificate
# echo  "Server's signed certificate"
# openssl x509 -in server_cert.pem -noout -text

# To verify that the server certificate is valid
openssl verify -CAfile ca_cert.pem server_cert.pem

# remove all Certificate Signing Request because we already have signed certificate for the web server and/or client (if we are using implementing Mutual TLS which we are not for this project)
rm *_csr.pem