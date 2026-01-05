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

// TestHTTPClient 测试 HTTP 客户端功能
func TestHTTPClient() {
	client := &http.Client{Timeout: 10 * time.Second}

	// 测试健康检查
	testHealthCheck(client)

	// 测试获取集群节点
	testGetClusterNodes(client)

	// 测试获取集群统计信息
	testGetClusterStats(client)

	// 测试获取 Actor 详情
	testGetActorDetails(client)

	// 测试 Account Actor 交互
	testAccountActorInteractions(client)

	// 测试 OpsActor 交互
	testOpsActorInteractions(client)
}

func testHealthCheck(client *http.Client) {
	fmt.Println("=== 测试健康检查 API ===")

	reqBody := map[string]interface{}{}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/HealthCheck", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("状态码: %d\n", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			return
		}
		fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

func testGetClusterNodes(client *http.Client) {
	fmt.Println("=== 测试获取集群节点 API ===")

	reqBody := map[string]interface{}{
		"includeDetails": true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/GetClusterNodes", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

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
		fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

func testGetClusterStats(client *http.Client) {
	fmt.Println("=== 测试获取集群统计信息 API ===")

	reqBody := map[string]interface{}{}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/GetClusterStats", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

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
		fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

func testGetActorDetails(client *http.Client) {
	fmt.Println("=== 测试获取 Actor 详情 API ===")

	reqBody := map[string]interface{}{
		"actorId": "account-001",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/GetActorDetails", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

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
		fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}

func testAccountActorInteractions(client *http.Client) {
	fmt.Println("=== 测试 Account Actor 交互 ===")

	// 创建账户
	createAccount(client, "account-001", 1000.0)

	// 查询账户
	getAccount(client, "account-001")

	// 存款
	creditAccount(client, "account-001", 500.0)

	// 再次查询账户
	getAccount(client, "account-001")

	// 测试另一个账户
	createAccount(client, "account-002", 2000.0)
	getAccount(client, "account-002")
}

func createAccount(client *http.Client, accountID string, balance float64) {
	fmt.Printf("创建账户: %s, 初始余额: %.2f\n", accountID, balance)

	reqBody := map[string]interface{}{
		"create_account": map[string]interface{}{
			"account_id":      accountID,
			"account_balance": balance,
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/samplepb.AccountService/CreateAccount", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
		} else {
			fmt.Printf("账户创建成功: %+v\n", result)
		}
	} else {
		fmt.Printf("账户创建失败，状态码: %d, 响应体: %s\n", resp.StatusCode, string(readBody))
	}
	fmt.Println()
}

func getAccount(client *http.Client, accountID string) {
	fmt.Printf("查询账户: %s\n", accountID)

	reqBody := map[string]interface{}{
		"account_id": accountID,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/samplepb.AccountService/GetAccount", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
		} else {
			fmt.Printf("账户查询成功: %+v\n", result)
		}
	} else {
		fmt.Printf("账户查询失败，状态码: %d, 响应体: %s\n", resp.StatusCode, string(readBody))
	}
	fmt.Println()
}

func creditAccount(client *http.Client, accountID string, amount float64) {
	fmt.Printf("存款到账户: %s, 金额: %.2f\n", accountID, amount)

	reqBody := map[string]interface{}{
		"credit_account": map[string]interface{}{
			"account_id": accountID,
			"balance":    amount,
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/samplepb.AccountService/CreditAccount", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
		} else {
			fmt.Printf("存款成功: %+v\n", result)
		}
	} else {
		fmt.Printf("存款失败，状态码: %d, 响应体: %s\n", resp.StatusCode, string(readBody))
	}
	fmt.Println()
}

func testOpsActorInteractions(client *http.Client) {
	fmt.Println("=== 测试 OpsActor 交互 ===")

	// 测试获取节点详情
	testGetNodeDetails(client)

	// 测试集群统计信息
	testGetClusterStats(client)

	// 测试健康检查
	testHealthCheck(client)
}

func testGetNodeDetails(client *http.Client) {
	fmt.Println("=== 测试获取节点详情 API ===")

	// 先获取集群节点列表
	reqBody := map[string]interface{}{
		"includeDetails": true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/GetClusterNodes", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			return
		}

		// 尝试从响应中获取节点ID
		if clusterInfo, ok := result["clusterInfo"].(map[string]interface{}); ok {
			if nodes, ok := clusterInfo["nodes"].([]interface{}); ok && len(nodes) > 0 {
				if firstNode, ok := nodes[0].(map[string]interface{}); ok {
					if nodeID, ok := firstNode["nodeId"].(string); ok {
						fmt.Printf("获取节点详情: %s\n", nodeID)
						getNodeDetails(client, nodeID)
						return
					}
				}
			}
		}
	}

	fmt.Println("无法获取节点ID，跳过节点详情测试")
}

func getNodeDetails(client *http.Client, nodeID string) {
	reqBody := map[string]interface{}{
		"node_id": nodeID,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/opspb.OpsService/GetNodeDetails", serverAddr), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	readBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("状态码: %d\n", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(readBody, &result); err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
		} else {
			fmt.Printf("节点详情: %+v\n", result)
		}
	} else {
		fmt.Printf("获取节点详情失败，状态码: %d, 响应体: %s\n", resp.StatusCode, string(readBody))
	}
	fmt.Println()
}
