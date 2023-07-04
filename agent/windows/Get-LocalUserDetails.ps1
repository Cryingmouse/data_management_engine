param(
    [string]$UserName
)

function Get-LocalUserDetails {
    param (
        [string]$UserName = ""
    )

    $users = Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True'"

    $userDetails = @{}

    foreach ($user in $users) {
        if ($UserName -eq "" -or $user.Name -eq $UserName) {
            $userProperties = @{
                'Name'               = $user.Name
                'SID'                = $user.SID
                'FullName'           = $user.FullName
                'Description'        = $user.Description
                'Status'             = $user.Status
                'Disabled'           = $user.Disabled
                'PasswordRequired'   = $user.PasswordRequired
                'PasswordExpires'    = $user.PasswordExpires
                'PasswordChangeable' = $user.PasswordChangeable
                'Lockout'            = $user.Lockout
                'PSComputerName'     = $user.PSComputerName
            }

            $userDetails[$user.Name] = $userProperties
        }
    }

    return $userDetails
}

# Call the function to retrieve and output local user details
if ($PSBoundParameters.Count -eq 0) {
    $localUserDetails = Get-LocalUserDetails
}
else {
    $localUserDetails = Get-LocalUserDetails -UserName $UserName
}

$localUserDetails | ConvertTo-Json

