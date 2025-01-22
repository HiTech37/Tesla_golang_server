package controller

import (
	"github.com/gin-gonic/gin"
)

func GetDeviceAccesToken(c *gin.Context) {

	// 	url := "https://fleet-api.prd.na.vn.cloud.tesla.com/api/1/partner_accounts"
	// 	method := "POST"

	// 	payload := strings.NewReader(`{
	//     "domain": "moovetrax.com",
	//     "csr": "-----BEGIN CERTIFICATE REQUEST-----\nMIHSMHoCAQAwGDEWMBQGA1UEAwwNbW9vdmV0cmF4LmNvbTBZMBMGByqGSM49AgEG\nCCqGSM49AwEHA0IABNgDPNNdDv5mH8FN3/Je1Onigpz8fimLfY3VtB8N31DAVdbJ\no6BFq4CAz18Ly9NIwe/0FmYyYIPX2AtrYwQG2wOgADAKBggqhkjOPQQDAgNIADBF\nAiEAh8dirlUY0P67UxL1e2mRAZtq2o9NRYVwOAZhIVjvhVgCIFLdDgL3TKzMMF7/\nJA9rVlboWyJ906ujEmwiQ6PL2dLL\n-----END CERTIFICATE REQUEST-----"
	// }`)

	// 	client := &http.Client{}
	// 	req, err := http.NewRequest(method, url, payload)

	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	req.Header.Add("Content-Type", "application/json")
	// 	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InFEc3NoM2FTV0cyT05YTTdLMzFWV0VVRW5BNCJ9.eyJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMiLCJzdWIiOiI3YjBlNGRjNS1lMGY4LTQ4ZTctYTU2Yi0yNTg3NGQ2Y2U2YTIiLCJpc3MiOiJodHRwczovL2ZsZWV0LWF1dGgudGVzbGEuY29tL29hdXRoMi92My9udHMiLCJhenAiOiI2OWU1NTgxNC0xNjc5LTQ2ZDMtYTNiNi1hYzcxM2Y3N2YyODciLCJhdWQiOlsiaHR0cHM6Ly9mbGVldC1hdXRoLnRlc2xhLmNvbS9vYXV0aDIvdjMvY2xpZW50aW5mbyIsImh0dHBzOi8vZmxlZXQtYXBpLnByZC5uYS52bi5jbG91ZC50ZXNsYS5jb20iXSwiZXhwIjoxNzM3MTUwODYwLCJpYXQiOjE3MzcxMjIwNjAsImFjY291bnRfdHlwZSI6ImJ1c2luZXNzIiwib3Blbl9zb3VyY2UiOmZhbHNlLCJzY3AiOlsidXNlcl9kYXRhIiwidmVoaWNsZV9kZXZpY2VfZGF0YSIsInZlaGljbGVfY21kcyIsInZlaGljbGVfY2hhcmdpbmdfY21kcyIsIm9wZW5pZCJdfQ.AIJxhkqN4yQ-c67EsWkGx2zUboagajs4gV9GOrpjyb-WzMupA6ld_jn9mJjQi6wl5tnA4udAQZ08DBP4fjO1LRPPxDglZRzl1UujRiDRBbnpW4n9cS25OBWZQci9RdhjnyAgQFmJNBa-q12HJtYJedqGUCMhud3ia-sQ4NG7pJybth05Z3ZHh6xO2KpV9nUFeyPt1hg1BwwSs0U7H_OzhKreBLkotvaOqHPD5q4fU_DmdAOm1-WIuU1LA9YwwPM4abhc1tBhD5Kj2Jn7sqLl62byKZHF37t3dStPRdjFHEp3kXVOGtF5ohAhsHu4UvvvFkgJB8hR8yJ6fXY7Z3TBNg")

	// 	res, err := client.Do(req)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	defer res.Body.Close()

	// 	body, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println(string(body))

}
