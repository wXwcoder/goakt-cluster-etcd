# PowerShell 脚本测试运维 API

Write-Host "=== 测试运维 API ===" -ForegroundColor Green

# 测试健康检查
Write-Host "`n1. 测试健康检查 API" -ForegroundColor Yellow
$healthCheckBody = @{}
$healthCheckJson = $healthCheckBody | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:9090/opspb.OpsService/HealthCheck" -Method Post -Body $healthCheckJson -ContentType "application/json" -Headers @{"Connect-Protocol-Version"="1"}
    Write-Host "健康检查响应: " -ForegroundColor Green -NoNewline
    Write-Host ($response | ConvertTo-Json -Depth 3) -ForegroundColor White
} catch {
    Write-Host "健康检查失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 测试获取集群节点
Write-Host "`n2. 测试获取集群节点 API" -ForegroundColor Yellow
$clusterNodesBody = @{
    includeDetails = $true
}
$clusterNodesJson = $clusterNodesBody | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:9090/opspb.OpsService/GetClusterNodes" -Method Post -Body $clusterNodesJson -ContentType "application/json" -Headers @{"Connect-Protocol-Version"="1"}
    Write-Host "集群节点响应: " -ForegroundColor Green -NoNewline
    Write-Host ($response | ConvertTo-Json -Depth 3) -ForegroundColor White
} catch {
    Write-Host "获取集群节点失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Green