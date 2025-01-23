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
