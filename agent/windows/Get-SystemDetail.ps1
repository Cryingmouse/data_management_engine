function Get-SystemDetails {
    $operatingSystem = Get-WmiObject -Class Win32_OperatingSystem | Select-Object Caption, Version, OSArchitecture, BuildNumber
    $computerName = $env:COMPUTERNAME

    $systemDetail = @{
        "Caption"        = $operatingSystem.Caption
        "Version"        = $operatingSystem.Version
        "OSArchitecture" = $operatingSystem.OSArchitecture
        "BuildNumber"    = $operatingSystem.BuildNumber
        "ComputerName"   = $computerName
    }

    return $systemDetail
}

# 调用函数以检索和输出本地系统详细信息
$systemDetail = Get-SystemDetails

$systemDetail | ConvertTo-Json