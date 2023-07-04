function Get-LocalUserDetails {
    $users = Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True'"

    $userDetails = @{}

    foreach ($user in $users) {
        $userProperties = @{
            'Name'             = $user.Name
            'SID'              = $user.SID
            'FullName'         = $user.FullName
            'Description'      = $user.Description
            'Status'           = $user.Status
            'Disabled'         = $user.Disabled
            'PasswordRequired' = $user.PasswordRequired
        }

        $userDetails[$user.Name] = $userProperties
    }

    return $userDetails
}

# Call the function to retrieve and output local user details
$localUserDetails = Get-LocalUserDetails
$localUserDetails | ConvertTo-Json
