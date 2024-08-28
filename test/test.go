package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const (
	apicURL      = "https://sandboxapicdc.cisco.com/api/node/mo/uni.json"
	apicLoginURL = "https://sandboxapicdc.cisco.com/api/aaaLogin.json"
	apicUsername = "admin"
	apicPassword = ""
)

func GetAPICLoginToken() (string, error) {
	loginBody := map[string]interface{}{
		"aaaUser": map[string]interface{}{
			"attributes": map[string]string{
				"name": apicUsername,
				"pwd":  apicPassword,
			},
		},
	}

	payloadBytes, err := json.Marshal(loginBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login payload: %w", err)
	}

	resp, err := http.Post(apicLoginURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	token, ok := response["imdata"].([]interface{})[0].(map[string]interface{})["aaaLogin"].(map[string]interface{})["attributes"].(map[string]interface{})["token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract token from response")
	}

	return token, nil

}

func PostToAPIC(token string, payload []byte) error {
	client := &http.Client{}

	req, err := http.NewRequest("POST", apicURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Cookie", "APIC-cookie="+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func CheckTerraformPlan() (bool, error) {
	planBin := "plan.bin"
	planJSON := "plan.json"

	cmdBin := exec.Command("terraform", "plan", "-refresh-only", "-out="+planBin)
	if err := cmdBin.Run(); err != nil {
		return false, fmt.Errorf("failed to run terraform plan: %w", err)
	}

	cmdOutJson := exec.Command("terraform", "show", "-json", planBin)
	output, err := cmdOutJson.Output()
	if err != nil {
		return false, fmt.Errorf("failed to show terraform plan: %w", err)
	}

	if err := os.WriteFile(planJSON, output, 0644); err != nil {
		return false, fmt.Errorf("failed to write JSON plan to file: %w", err)
	}

	var plan map[string]interface{}
	if err := json.Unmarshal(output, &plan); err != nil {
		return false, fmt.Errorf("failed to parse terraform plan output: %w", err)
	}

	if changes, ok := plan["planned_values"].(map[string]interface{})["root_module"].(map[string]interface{}); ok && len(changes) == 0 {
		return true, nil
	}

	return false, nil
}
