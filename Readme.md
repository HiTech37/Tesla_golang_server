- Build executive file.

    go build .

- Update the permission of executive file (tesla_server)
  
    sudo chmod -R 775 main
    sudo chown ec2-user:ec2-user tesla_server

sudo mv tesla_server /usr/local/bin/

sudo nano /etc/systemd/system/tesla_server.service

        [Unit]
        Description=My Go Application
        After=network.target

        [Service]
        ExecStart=/usr/local/bin/tesla_server
        Restart=on-failure

        [Install]
        WantedBy=multi-user.target

sudo systemctl daemon-reload

sudo systemctl start tesla_server

sudo systemctl enable tesla_server

sudo systemctl status tesla_server
sudo systemctl stop tesla_server




sudo firewall-cmd --zone=public --add-port=9443/tcp --permanent
sudo firewall-cmd --reload

-To proxy, can use nginx and it is set on fleet telemetry deployment

sudo ss -tulnp

export TESLA_AUTH_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InFEc3NoM2FTV0cyT05YTTdLMzFWV0VVRW5BNCJ9.eyJpc3MiOiJodHRwczovL2F1dGgudGVzbGEuY29tL29hdXRoMi92My9udHMiLCJhenAiOiI2OWU1NTgxNC0xNjc5LTQ2ZDMtYTNiNi1hYzcxM2Y3N2YyODciLCJzdWIiOiJhMTZmMDc2MS1kYzQ1LTQwYjktOGY4ZC0xZWFhNTUwZWY3OGMiLCJhdWQiOlsiaHR0cHM6Ly9mbGVldC1hcGkucHJkLm5hLnZuLmNsb3VkLnRlc2xhLmNvbSIsImh0dHBzOi8vZmxlZXQtYXBpLnByZC5ldS52bi5jbG91ZC50ZXNsYS5jb20iLCJodHRwczovL2F1dGgudGVzbGEuY29tL29hdXRoMi92My91c2VyaW5mbyJdLCJzY3AiOlsib3BlbmlkIiwidmVoaWNsZV9kZXZpY2VfZGF0YSIsIm9mZmxpbmVfYWNjZXNzIiwidmVoaWNsZV9jbWRzIiwidmVoaWNsZV9jaGFyZ2luZ19jbWRzIl0sImFtciI6WyJwd2QiLCJtZmEiLCJvdHAiLCJib3RwIl0sImV4cCI6MTczNzU1NTA4NiwiaWF0IjoxNzM3NTI2Mjg2LCJvdV9jb2RlIjoiTkEiLCJsb2NhbGUiOiJlbi1VUyIsImFjY291bnRfdHlwZSI6ImJ1c2luZXNzIiwib3Blbl9zb3VyY2UiOmZhbHNlLCJhY2NvdW50X2lkIjoiN2IwZTRkYzUtZTBmOC00OGU3LWE1NmItMjU4NzRkNmNlNmEyIiwiYXV0aF90aW1lIjoxNzM3NTI2MjgxfQ.DjD77K1MKf-w8ITvz_gpTRo7n_cOiFGoVjLHcd-1YdFC6cJLkxhiqTArXucuZ2jzlBWb2BFrd2lRXGaGvjOoWBglE2a8DVGoglbGLEJYodlvTwTD7FBmtPFCl9I6uKqyRJt9YdJAv6niK54sUqoVmSbA4u38UQM5uvkWgeMwAls3yyEMcodZ3f_cZ3UTxEAcQsrdhg21NqN5HkHBdJxB8z8oTapxtwPRyYSrEqC1On5xPZllZsA3-IM45llXJgvXNkTV_X5AAoUoiJPBASgzvL-q27ztfXwfd_0miSKhZOUdi3HxMgw7qvdLYLa-f4F_xz96bCCpoJzcg-H0q9iRIw"
export VIN="5YJSA1E20JF291673"
curl --cacert cert.pem \
--header "Authorization: Bearer $TESLA_AUTH_TOKEN" \
"https://localhost:4443/api/1/vehicles/$VIN/vehicle_data" \
| jq -r .

export TESLA_AUTH_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InFEc3NoM2FTV0cyT05YTTdLMzFWV0VVRW5BNCJ9.eyJpc3MiOiJodHRwczovL2F1dGgudGVzbGEuY29tL29hdXRoMi92My9udHMiLCJhenAiOiI2OWU1NTgxNC0xNjc5LTQ2ZDMtYTNiNi1hYzcxM2Y3N2YyODciLCJzdWIiOiJhMTZmMDc2MS1kYzQ1LTQwYjktOGY4ZC0xZWFhNTUwZWY3OGMiLCJhdWQiOlsiaHR0cHM6Ly9mbGVldC1hcGkucHJkLm5hLnZuLmNsb3VkLnRlc2xhLmNvbSIsImh0dHBzOi8vZmxlZXQtYXBpLnByZC5ldS52bi5jbG91ZC50ZXNsYS5jb20iLCJodHRwczovL2F1dGgudGVzbGEuY29tL29hdXRoMi92My91c2VyaW5mbyJdLCJzY3AiOlsib3BlbmlkIiwidmVoaWNsZV9kZXZpY2VfZGF0YSIsIm9mZmxpbmVfYWNjZXNzIiwidmVoaWNsZV9jbWRzIiwidmVoaWNsZV9jaGFyZ2luZ19jbWRzIl0sImFtciI6WyJwd2QiLCJtZmEiLCJvdHAiLCJib3RwIl0sImV4cCI6MTczNzU1NTA4NiwiaWF0IjoxNzM3NTI2Mjg2LCJvdV9jb2RlIjoiTkEiLCJsb2NhbGUiOiJlbi1VUyIsImFjY291bnRfdHlwZSI6ImJ1c2luZXNzIiwib3Blbl9zb3VyY2UiOmZhbHNlLCJhY2NvdW50X2lkIjoiN2IwZTRkYzUtZTBmOC00OGU3LWE1NmItMjU4NzRkNmNlNmEyIiwiYXV0aF90aW1lIjoxNzM3NTI2MjgxfQ.DjD77K1MKf-w8ITvz_gpTRo7n_cOiFGoVjLHcd-1YdFC6cJLkxhiqTArXucuZ2jzlBWb2BFrd2lRXGaGvjOoWBglE2a8DVGoglbGLEJYodlvTwTD7FBmtPFCl9I6uKqyRJt9YdJAv6niK54sUqoVmSbA4u38UQM5uvkWgeMwAls3yyEMcodZ3f_cZ3UTxEAcQsrdhg21NqN5HkHBdJxB8z8oTapxtwPRyYSrEqC1On5xPZllZsA3-IM45llXJgvXNkTV_X5AAoUoiJPBASgzvL-q27ztfXwfd_0miSKhZOUdi3HxMgw7qvdLYLa-f4F_xz96bCCpoJzcg-H0q9iRIw"
curl --cacert cert.pem \
    -H "Authorization: Bearer $TESLA_AUTH_TOKEN" \
     -H 'Content-Type: application/json' \
     --data '{
			    "vins": [
			        "5YJSA1E20JF291673"
			    ]
			}' \
      -X POST \
      -i https://localhost:4443/api/1/vehicles/fleet_status


export TESLA_AUTH_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InFEc3NoM2FTV0cyT05YTTdLMzFWV0VVRW5BNCJ9.eyJpc3MiOiJodHRwczovL2F1dGgudGVzbGEuY29tL29hdXRoMi92My9udHMiLCJhenAiOiI2OWU1NTgxNC0xNjc5LTQ2ZDMtYTNiNi1hYzcxM2Y3N2YyODciLCJzdWIiOiJhMTZmMDc2MS1kYzQ1LTQwYjktOGY4ZC0xZWFhNTUwZWY3OGMiLCJhdWQiOlsiaHR0cHM6Ly9mbGVldC1hcGkucHJkLm5hLnZuLmNsb3VkLnRlc2xhLmNvbSIsImh0dHBzOi8vZmxlZXQtYXBpLnByZC5ldS52bi5jbG91ZC50ZXNsYS5jb20iLCJodHRwczovL2F1dGgudGVzbGEuY29tL29hdXRoMi92My91c2VyaW5mbyJdLCJzY3AiOlsib2ZmbGluZV9hY2Nlc3MiLCJvcGVuaWQiLCJ1c2VyX2RhdGEiLCJ2ZWhpY2xlX2RldmljZV9kYXRhIiwidmVoaWNsZV9jbWRzIiwidmVoaWNsZV9jaGFyZ2luZ19jbWRzIl0sImFtciI6WyJwd2QiLCJtZmEiLCJvdHAiLCJib3RwIl0sImV4cCI6MTczNzY4OTYxNSwiaWF0IjoxNzM3NjYwODE1LCJvdV9jb2RlIjoiTkEiLCJsb2NhbGUiOiJlbi1VUyIsImFjY291bnRfdHlwZSI6ImJ1c2luZXNzIiwib3Blbl9zb3VyY2UiOmZhbHNlLCJhY2NvdW50X2lkIjoiN2IwZTRkYzUtZTBmOC00OGU3LWE1NmItMjU4NzRkNmNlNmEyIiwiYXV0aF90aW1lIjoxNzM3NjYwODE1LCJub25jZSI6bnVsbH0.HkT8upf7zZAkxKmu__n49y1dL2Kl9xGKf7r7e_ALsJsKKq3YeP2SAkgE_hub3mGXRZduaDSU-QdYcf6YIOiW1X992PRrh2Xmr1pDYJ1-uDNT1A-au1_68uIjH79NyQCtLIIFdYuA6e0k67stRsGwcsyiMYmcNUbBuLh0tin1Qwn9s2xuO_ebFEdNd5dw02QECkTpbF5-xSLmzPnckDxbUuNwye_ALh6M_U4LFig6k-LSoITzhluNi8WcktLNfsbBn3ij5ntGs0K4c1kzNgEzk0mS48g6oA7k5__6P05v5tlVk25oxGTjXBoeTK6m0AGBO9QI0SJlUwIRVrzK68QzNw"
export VIN="5YJSA1E20JF291673"
curl --cacert cert.pem \
--header "Authorization: Bearer $TESLA_AUTH_TOKEN" \
"https://localhost:4443/api/1/vehicles/$VIN/vehicle_data" \
| jq -r .


{
    "response": {
        "account_id": "7b0e4dc5-e0f8-48e7-a56b-25874d6ce6a2",
        "domain": "t3slaapi.moovetrax.com",
        "name": "MooveTrax-New-API",
        "description": "We would like to develop integration for the new api while keeping the old",
        "client_id": "69e55814-1679-46d3-a3b6-ac713f77f287",
        "ca": null,
        "created_at": "2024-09-12T01:02:16.051Z",
        "updated_at": "2025-01-23T19:53:09.678Z",
        "enterprise_tier": "pay_as_you_go",
        "issuer": null,
        "csr": null,
        "csr_updated_at": null,
        "public_key": "046fbe009f96e70f53c61c5ceb9979e27dd47f12d8a3658b29880da358d8b92b71ae0c2698b5ebb15e46dc4672cb1e38f588f1c32854ac7bb9416e9a2e9c294492"
    }
}

{
    "response": {
        "account_id": "7b0e4dc5-e0f8-48e7-a56b-25874d6ce6a2",
        "domain": "t3slaapi.moovetrax.com",
        "name": "MooveTrax-New-API",
        "description": "We would like to develop integration for the new api while keeping the old",
        "client_id": "69e55814-1679-46d3-a3b6-ac713f77f287",
        "ca": null,
        "created_at": "2024-09-12T01:02:16.051Z",
        "updated_at": "2025-01-23T20:20:31.398Z",
        "enterprise_tier": "pay_as_you_go",
        "issuer": null,
        "csr": null,
        "csr_updated_at": null,
        "public_key": "04b2f349b74180568d487bd8b16a4e5024bad70efa767fe242356e3db12ad61fc7b6c55621a1651d2b0775c9c9031fcbb568ba50167540eeff5cdd3da97501abb7"
    }
}

openssl req -out t3slaapi.moovetrax.com.csr -key private-key.pem -subj /CN=t3slaapi.moovetrax.com/ -new



sudo yum install certbot python3-certbot-nginx -y

sudo certbot --nginx -d fleetapi.moovetrax.com

sudo systemctl enable certbot-renew.timer

sudo systemctl start certbot-renew.timer

sudo certbot renew --dry-run

certbot certonly -d fleetapi.moovetrax.com --csr fleetapi.moovetrax.com.csr

echo | openssl s_client -connect fleetapi.moovetrax.com:8443 -servername fleetapi.moovetrax.com -showcerts 2>/dev/null | awk '/BEGIN CERTIFICATE/,/END CERTIFICATE/ {print}' > ca_cert.pem

openssl req -x509 -nodes -newkey ec \
    -pkeyopt ec_paramgen_curve:secp521r1 \
    -pkeyopt ec_param_enc:named_curve  \
    -subj '/CN=fleetapi.moovetrax.com' \
    -keyout key.pem -out cert.pem -sha256 -days 3650 \
    -addext "subjectAltName = DNS:fleetapi.moovetrax.com" \
    -addext "extendedKeyUsage = serverAuth" \
    -addext "keyUsage = digitalSignature, keyCertSign, keyAgreement" 


openssl req -out fleetapi.moovetrax.com.csr -key private_key.pem -subj /CN=fleetapi.moovetrax.com/ -new


export TESLA_AUTH_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InFEc3NoM2FTV0cyT05YTTdLMzFWV0VVRW5BNCJ9.eyJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMiLCJzdWIiOiIzYTNhOGVlZi1hNTEzLTQwODQtYjM1My1hNTAyYmI5MzZmNTUiLCJpc3MiOiJodHRwczovL2ZsZWV0LWF1dGgudGVzbGEuY29tL29hdXRoMi92My9udHMiLCJhenAiOiI2MGQ5NzkxOC05YjZiLTRjOTItODhlMy1mZjllOTQwMzIzOWYiLCJhdWQiOlsiaHR0cHM6Ly9mbGVldC1hdXRoLnRlc2xhLmNvbS9vYXV0aDIvdjMvY2xpZW50aW5mbyIsImh0dHBzOi8vZmxlZXQtYXBpLnByZC5uYS52bi5jbG91ZC50ZXNsYS5jb20iXSwiZXhwIjoxNzQwMDIxODM4LCJpYXQiOjE3Mzk5OTMwMzgsImFjY291bnRfdHlwZSI6ImJ1c2luZXNzIiwib3Blbl9zb3VyY2UiOmZhbHNlLCJzY3AiOlsidmVoaWNsZV9kZXZpY2VfZGF0YSIsInZlaGljbGVfbG9jYXRpb24iLCJ2ZWhpY2xlX2NtZHMiLCJ2ZWhpY2xlX2NoYXJnaW5nX2NtZHMiLCJvcGVuaWQiLCJvZmZsaW5lX2FjY2VzcyJdfQ.bO-gfdvTyGXA_BnYXK_VUqQ9J6eGmRUenuOEgqrNDk9_K9zmhRt_mfyHw8vQjwMAg_3qrjSL8ZvB6qEiFBNWWjWuRtAcb4JJXkIRbcVQtX31NkTeJGrfsQlm5GCB9YwPkIH2MDwmjDTr6-g9x3Vxy7yhAr91k4ujMbljWsS_NAL-BY2QGnCIZK1awWEFqMPY_25AnJ7U7bc6eHARhv0_1YA6nl5T5SNUKK2FrYiJgkMBdqUQNqLQ_RJdbf45lBaVLTETSqgPA-2kPSKWnCQg_981N53Ok0nBtqyBCIkQa4winw9viE5U7CnRw5D1qOnW5vrVpgGmygF8_XXTUomNDw"
export VIN="7SAYGDEE9RA313640"
curl --cacert cert.pem \
--header "Authorization: Bearer $TESLA_AUTH_TOKEN" \
"https://fleetapi.moovetrax.com:4443/api/1/vehicles/$VIN/vehicle_data" \
| jq -r .

curl -v --cacert /etc/cert/moovetrax-fullchain.crt \
    --cert cert.pem \
    --key key.pem \
    https://fleetapi.moovetrax.com:8443
 