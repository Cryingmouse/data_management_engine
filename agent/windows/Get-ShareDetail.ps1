param (
    [String] $ShareNames
)

$sharesDetail = @()

$ShareNamesArray = $ShareNames -split ','

$allShares = Get-SmbShare | Select-Object Name, ShareState, Description, Path

if (-not $ShareNames) {
    foreach ($share in $allShares) {
        $sharesDetail += @{
            'Name'          = $share.Name
            'ShareState'    = $share.ShareState
            'Description'   = $share.Description
            'DirectoryPath' = $share.Path
        }
    }
}
else {
    $shares = $allShares | Where-Object { $_.Name -in $ShareNamesArray }
    foreach ( $share in $shares ) {
        $sharesDetail += @{
            'Name'          = $share.Name
            'ShareState'    = $share.ShareState
            'Description'   = $share.Description
            'DirectoryPath' = $share.Path
        }
    }
}


$sharesDetail | ConvertTo-Json

