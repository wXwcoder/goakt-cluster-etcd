package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
}

func testHealthCheck(client *http.Client) {
	fmt.Println("=== 测试健康检查 API ===")

	// 创建健康检查请求
	reqBody := map[string]interface{}{}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "http://localhost:13001/opspb.OpsService/HealthCheck", bytes.NewBuffer(jsonBody))
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

	// 创建获取集群节点请求
	reqBody := map[string]interface{}{
		"includeDetails": true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "http://localhost:13001/opspb.OpsService/GetClusterNodes", bytes.NewBuffer(jsonBody))
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
		fmt.Printf("响应内容: %+v\n", result)
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

	req, err := http.NewRequest("POST", "http://localhost:13001/opspb.OpsService/GetActorDetails", bytes.NewBuffer(jsonBody))
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
		fmt.Printf("响应内容: %+v\n", result)
	} else {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
	}
	fmt.Println()
}
