package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var serverAddr = "localhost:14001"

func main() {
	client := &http.Client{Timeout: 10 * time.Second}

	// 测试健康检查
	testHealthCheck(client)

	// 测试获取集群节点
	testGetClusterNodes(client)

	// 测试获取集群统计信息
	testGetClusterStats(client)

	// 测试获取 Actor 详情
	testGetActorDetails(client)

	// 测试AccountService相关接口
	testCreateAccount(client)
	testCreditAccount(client)
	testGetAccount(client)
}

func testHealthCheck(client *http.Client) {
	fmt.Println("=== 测试健康检查 API ===")

	// 创建健康检查请求
	reqBody := map[string]interface{}{}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/HealthCheck", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置 Connect-go 需要的头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("状态码: %d body: %s\n", resp.StatusCode, string(readBody))

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			return
		}
		//fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

func testGetClusterNodes(client *http.Client) {
	fmt.Println("=== 测试获取集群节点 API ===")

	// 创建获取集群节点请求
	reqBody := map[string]interface{}{
		"includeDetails": true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/GetClusterNodes", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置 Connect-go 需要的头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n 响应体: %s\n", resp.StatusCode, string(readBody))

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			return
		}
		//fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

func testGetClusterStats(client *http.Client) {
	fmt.Println("=== 测试获取集群统计信息 API ===")

}

// 测试获取 Actor 详情
func testGetActorDetails(client *http.Client) {
	fmt.Println("=== 测试获取 Actor 详情 API ===")

	// 创建获取 Actor 详情请求
	reqBody := map[string]interface{}{
		"actorId": "account-001",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/GetActorDetails", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置 Connect-go 需要的头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n 响应体: %s\n", resp.StatusCode, string(readBody))

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			return
		}
		//fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

// 测试创建账户
func testCreateAccount(client *http.Client) {
	fmt.Println("=== 测试创建账户 API ===")

	// 创建创建账户请求
	reqBody := map[string]interface{}{
		"create_account": map[string]interface{}{
			"account_id":      "account-001",
			"account_balance": 1000.0,
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/samplepb.AccountService/CreateAccount", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置 Connect-go 需要的头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n 响应体: %s\n", resp.StatusCode, string(readBody))

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			return
		}
		//fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

// 测试信用账户
func testCreditAccount(client *http.Client) {
	fmt.Println("=== 测试信用账户 API ===")

	// 创建信用账户请求
	reqBody := map[string]interface{}{
		"credit_account": map[string]interface{}{
			"account_id": "account-001",
			"balance":    500.0,
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/samplepb.AccountService/CreditAccount", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置 Connect-go 需要的头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n 响应体: %s\n", resp.StatusCode, string(readBody))

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			return
		}
		//fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

// 测试获取账户
func testGetAccount(client *http.Client) {
	fmt.Println("=== 测试获取账户 API ===")

	// 创建获取账户请求
	reqBody := map[string]interface{}{
		"account_id": "account-001",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/samplepb.AccountService/GetAccount", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置 Connect-go 需要的头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n 响应体: %s\n", resp.StatusCode, string(readBody))

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			return
		}
		//fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}
