function Get-SystemDetails {
    $operatingSystem = Get-WmiObject -Class Win32_OperatingSystem | Select-Object Caption, Version, OSArchitecture, BuildNumber
    $computerName = $env:COMPUTERNAME

    $systemDetails = @{
        "Caption"        = $operatingSystem.Caption
        "Version"        = $operatingSystem.Version
        "OSArchitecture" = $operatingSystem.OSArchitecture
        "BuildNumber"    = $operatingSystem.BuildNumber
        "ComputerName"   = $computerName
    }

    return $systemDetails
}

# 调用函数以检索和输出本地系统详细信息
$systemDetails = Get-SystemDetails

$systemDetails | ConvertTo-Json