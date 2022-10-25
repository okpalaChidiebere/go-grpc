This directory contains x509 certificates and associated private keys used in implementing server-side TLS where the ONLy the server need to provide its TLS certificate to the client; Our server don't care about which client is calling its API. Full [Tutorial](https://www.youtube.com/watch?v=jmqLJMFS_y) and from [techSchool](https://dev.to/techschoolguru/how-to-secure-grpc-connection-with-ssl-tls-in-go-4ph)

If we want need our server to make sure that the right client is calling its API, then we will implement Mutual TLS. See [here](https://github.com/grpc/grpc-go/tree/master/examples/features/encryption/mTLS)

# With Micro Services

- You will use [Let'sEncrypt](https://letsencrypt.org/) to help with Automatic Certificate renewal when certificate expires. We will learn about that in our BackEnd Course

## Flag explantion

- `-x509` This ensures that this key is in a self-signed certificate format
- `-newKey rsa:4096` tells openssl to create a 4096 bit private key
- `-days 365` specify the number of days that this cert is valid for
- `-keyout ca_key.pem` specifies that we want to write the private key to a file called ca_key.pem
- `-out ca_cert.pem` specify the file we want to write the certificate to
- `-subj` option helps us to set the identity information used in completing signing the certificate non-interactively. This step is usually interactive in the termial but we can skip it and just feed the values in
- `-nodes` tells ssl to skip asking us to a enter passphrase to use to encrypt the private key before writing it to the pem file. A private key will still be generated but it will not be encrypted. This is fine to skip because we are using this certificate for dev mode :) This pass phrase is good to use just because if your private key is stolen by an attacker, the attacker will need to know to pass phrase in ordrr to decrypt your private key
- `-in ca-cert.perm` we are giving the x509 our original certificate
- `-noout` tells x509 not to print the original encoded value
- `-text` we want the text format.
- `-req` tells ssl that we will be passing in a certifcate request
- CAcreateserial ensures that the CA signing this self-signed certificate uses a unique serial number to sign the certificate. You will see this number when you print out the server_cert.pem file in human readable format like in the next comand
- `-extfile server-ext.cnf` is used to list all the domains, emails or IP address that we will want this signed certificate request to support. See [here](https://man.openbsd.org/x509v3.cnf.5#Subject_alternative_name)
- See about more flags [https://man.openbsd.org/openssl.1](https://man.openbsd.org/openssl.1)
